import pg from 'pg';
import dotenv from 'dotenv';

dotenv.config();

interface StockSnapshot {
  stock_id: number;
  week: number;
  rating_from: string;
  rating_to: string;
  target_from: number;
  target_to: number;
  price: number;
  action: string;
  market_sentiment: string;
  signal_strength: number;
}

interface WeekPrediction {
  week: number;
  market_sentiment: string;
  signal_strength: number;
}

function getRandomFloat(min: number, max: number): number {
  return Math.round((Math.random() * (max - min) + min) * 100) / 100;
}

async function getMarketPredictions(
  company: string,
  ticker: string,
  price: number,
  rating: string
): Promise<WeekPrediction[]> {
  const prompt = `You are an API that only returns JSON. Do not explain anything.
Given the stock ${company} (${ticker}), with a current price of $${price} and a rating of ${rating}, simulate the market behavior for the next 4 weeks.

For each week, return:
- week: number (2 to 5)
- market_sentiment: "positive", "neutral" or "negative"
- signal_strength: number (between 0.0 and 1.0)

Respond ONLY with a valid JSON array like:
[
  { "week": 2, "market_sentiment": "positive", "signal_strength": 0.82 },
  { "week": 3, "market_sentiment": "neutral", "signal_strength": 0.45 },
  { "week": 4, "market_sentiment": "negative", "signal_strength": 0.72 },
  { "week": 5, "market_sentiment": "positive", "signal_strength": 0.66 }
]`;

  const response = await fetch('https://openrouter.ai/api/v1/chat/completions', {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${process.env.OPENROUTER_API_KEY}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      model: process.env.OPENROUTER_MODEL,
      messages: [{ role: 'user', content: prompt }],
    }),
  });

  if (!response.ok) {
    throw new Error(`OpenRouter API error: ${response.statusText}`);
  }

  const data = await response.json();
  let predictions: WeekPrediction[];

  try {
    const content = data.choices[0].message.content;
    predictions = JSON.parse(content.trim());
    console.log({ predictions });
  } catch (error) {
    console.error('Failed to parse AI response:', error);
    throw new Error('Invalid AI response format');
  }

  return predictions;
}

function calculateNewPrice(basePrice: number, sentiment: string, signalStrength: number): number {
  let change: number;

  switch (sentiment.toLowerCase()) {
    case 'positive':
      change = signalStrength * getRandomFloat(0.05, 0.1);
      break;
    case 'negative':
      change = -signalStrength * getRandomFloat(0.05, 0.1);
      break;
    default: // neutral
      change = getRandomFloat(-0.01, 0.01);
  }

  return Math.round(basePrice * (1 + change) * 100) / 100;
}

async function main() {
  if (!process.env.DB_URL || !process.env.OPENROUTER_API_KEY) {
    console.error('Required environment variables are not set');
    process.exit(1);
  }

  const pool = new pg.Pool({
    connectionString: process.env.DB_URL,
    ssl: { rejectUnauthorized: false },
  });

  try {
    const week1Data = await pool.query(`
      SELECT s.id as stock_id, s.company, s.ticker, ss.*
      FROM stocks s
      JOIN stock_snapshots ss ON s.id = ss.stock_id
      WHERE ss.week = 1
    `);

    const totalStocks = week1Data.rows.length;
    let processedStocks = 0;
    let successfulStocks = 0;
    let failedStocks = 0;

    console.log(`Found ${totalStocks} stocks to process\n`);

    for (const baseSnapshot of week1Data.rows) {
      processedStocks++;
      console.log(
        `Processing stock ${processedStocks}/${totalStocks}: ${baseSnapshot.ticker} (${baseSnapshot.company})...`
      );

      try {
        const predictions = await getMarketPredictions(
          baseSnapshot.company,
          baseSnapshot.ticker,
          baseSnapshot.price,
          baseSnapshot.rating_to
        );

        let currentPrice = baseSnapshot.price;

        for (const prediction of predictions) {
          const newPrice = calculateNewPrice(currentPrice, prediction.market_sentiment, prediction.signal_strength);

          const priceChange = newPrice - currentPrice;
          const targetAdjustment = priceChange > 0 ? 1.05 : priceChange < 0 ? 0.95 : 1;

          const snapshot: StockSnapshot = {
            stock_id: baseSnapshot.stock_id,
            week: prediction.week,
            rating_from: baseSnapshot.rating_from,
            rating_to: baseSnapshot.rating_to,
            target_from: baseSnapshot.target_from * targetAdjustment,
            target_to: baseSnapshot.target_to * targetAdjustment,
            price: newPrice,
            action: baseSnapshot.action,
            market_sentiment: prediction.market_sentiment,
            signal_strength: prediction.signal_strength,
          };

          await pool.query(
            `
            INSERT INTO stock_snapshots (
              stock_id, week, rating_from, rating_to,
              target_from, target_to, price, action,
              market_sentiment, signal_strength, created_at
            )
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
            ON CONFLICT (stock_id, week) DO UPDATE SET
              rating_from = EXCLUDED.rating_from,
              rating_to = EXCLUDED.rating_to,
              target_from = EXCLUDED.target_from,
              target_to = EXCLUDED.target_to,
              price = EXCLUDED.price,
              action = EXCLUDED.action,
              market_sentiment = EXCLUDED.market_sentiment,
              signal_strength = EXCLUDED.signal_strength,
              created_at = EXCLUDED.created_at
          `,
            [
              snapshot.stock_id,
              snapshot.week,
              snapshot.rating_from,
              snapshot.rating_to,
              snapshot.target_from,
              snapshot.target_to,
              snapshot.price,
              snapshot.action,
              snapshot.market_sentiment,
              snapshot.signal_strength,
            ]
          );

          currentPrice = newPrice;
          console.log(
            `  Week ${prediction.week}: ${prediction.market_sentiment} (${prediction.signal_strength.toFixed(
              2
            )}) -> $${newPrice}`
          );
        }

        successfulStocks++;
      } catch (error) {
        console.error(`Failed to process ${baseSnapshot.ticker}:`, error instanceof Error ? error.message : error);
        failedStocks++;
        continue;
      }
    }

    console.log('\nSimulation complete!');
    console.log(`Successfully processed: ${successfulStocks}/${totalStocks} stocks`);
    if (failedStocks > 0) {
      console.log(`Failed to process: ${failedStocks} stocks`);
    }
  } catch (error) {
    console.error('Database error:', error instanceof Error ? error.message : error);
  } finally {
    await pool.end();
  }
}

main().catch(console.error);
