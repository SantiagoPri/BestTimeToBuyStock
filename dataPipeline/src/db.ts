import pg from 'pg';
import dotenv from 'dotenv';
import type { Stock } from './fetchData.js';

dotenv.config();

const { Pool } = pg;

const pool = new Pool({
  connectionString: process.env.DB_URL,
  ssl: {
    rejectUnauthorized: false,
  },
});

export async function initializeDatabase(): Promise<void> {
  const client = await pool.connect();
  try {
    await client.query(`
      CREATE TABLE IF NOT EXISTS stocks (
        id SERIAL PRIMARY KEY,
        ticker TEXT UNIQUE NOT NULL,
        company TEXT NOT NULL,
        category TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
      )
    `);

    await client.query(`
      CREATE TABLE IF NOT EXISTS stock_snapshots (
        id SERIAL PRIMARY KEY,
        stock_id INTEGER REFERENCES stocks(id),
        week INTEGER NOT NULL,
        rating_from TEXT NOT NULL,
        rating_to TEXT NOT NULL,
        target_from NUMERIC NOT NULL,
        target_to NUMERIC NOT NULL,
        price NUMERIC NOT NULL,
        action TEXT NOT NULL,
        news_title TEXT NOT NULL,
        news_summary TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(stock_id, week)
      )
    `);

    await client.query(`
      CREATE INDEX IF NOT EXISTS idx_stocks_ticker ON stocks(ticker);
      CREATE INDEX IF NOT EXISTS idx_stock_snapshots_stock_id_week ON stock_snapshots(stock_id, week);
    `);
  } finally {
    client.release();
  }
}

function generateRandomPrice(targetFrom: number, targetTo: number): number {
  const avgTarget = (targetFrom + targetTo) / 2;
  const randomFactor = 0.85 + Math.random() * 0.25; // Random between 0.85 and 1.1
  return avgTarget * randomFactor;
}

function parsePrice(priceStr: string): number {
  return parseFloat(priceStr.replace('$', ''));
}

export async function insertStocks(stocks: Stock[]): Promise<void> {
  const client = await pool.connect();
  try {
    // Begin transaction
    await client.query('BEGIN');

    for (const stock of stocks) {
      // Try to insert the stock, if it doesn't exist
      const stockResult = await client.query(
        `
        INSERT INTO stocks (ticker, company, category)
        VALUES ($1, $2, $3)
        ON CONFLICT (ticker) DO UPDATE 
        SET company = EXCLUDED.company
        RETURNING id
        `,
        [stock.ticker, stock.company, 'Unclassified']
      );

      const stockId = stockResult.rows[0].id;
      const targetFrom = parsePrice(stock.target_from);
      const targetTo = parsePrice(stock.target_to);
      const price = generateRandomPrice(targetFrom, targetTo);
      const newsTitle = `Recommendation by ${stock.brokerage}`;
      const newsSummary = `${stock.brokerage} rated ${stock.company} as ${stock.rating_to} with a price target of ${stock.target_to}.`;

      await client.query(
        `
        INSERT INTO stock_snapshots (
          stock_id, week, rating_from, rating_to, target_from, target_to,
          price, action, news_title, news_summary
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        ON CONFLICT (stock_id, week) DO UPDATE
        SET rating_from = EXCLUDED.rating_from,
            rating_to = EXCLUDED.rating_to,
            target_from = EXCLUDED.target_from,
            target_to = EXCLUDED.target_to,
            price = EXCLUDED.price,
            action = EXCLUDED.action,
            news_title = EXCLUDED.news_title,
            news_summary = EXCLUDED.news_summary
        `,
        [
          stockId,
          1,
          stock.rating_from,
          stock.rating_to,
          targetFrom,
          targetTo,
          price,
          stock.action,
          newsTitle,
          newsSummary,
        ]
      );
    }

    await client.query('COMMIT');
  } catch (error) {
    await client.query('ROLLBACK');
    if (error instanceof Error) {
      throw new Error(`Failed to insert stock data: ${error.message}`);
    }
    throw new Error('Failed to insert stock data: Unknown error');
  } finally {
    client.release();
  }
}
