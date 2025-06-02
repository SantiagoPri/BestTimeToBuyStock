import pg from 'pg';
import axios from 'axios';
import dotenv from 'dotenv';

dotenv.config();

const VALID_CATEGORIES = [
  'Tech',
  'Healthcare',
  'Finance',
  'Energy',
  'Consumer',
  'Industrial',
  'Telecom',
  'Real Estate',
  'Utilities',
  'Materials',
] as const;

async function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function chunkArray<T>(arr: T[], size: number): T[][] {
  const chunks: T[][] = [];
  for (let i = 0; i < arr.length; i += size) {
    chunks.push(arr.slice(i, i + size));
  }
  return chunks;
}

async function getCategoriesFromAI(companies: string[]): Promise<Record<string, string>> {
  const prompt = `
  You are an assistant that must output a JSON array of exactly ${
    companies.length
  } objects. Each object must have two keys: "company" (the company name exactly as given) and "category" (exactly one of these labels: Tech, Healthcare, Finance, Energy, Consumer, Industrial, Telecom, Real Estate, Utilities, Materials). If a company does not fit any of those 10 categories, set its "category" to "Others". Do NOT output anything other than valid JSON—no explanations, no extra punctuation.
  Companies:
  ${companies.map((c, i) => `${i + 1}. ${c}`).join('\n')}
  Answer:
  `;

  const response = await axios.post(
    'https://openrouter.ai/api/v1/chat/completions',
    {
      model: process.env.OPENROUTER_MODEL,
      messages: [{ role: 'user', content: prompt.trim() }],
    },
    {
      headers: {
        Authorization: `Bearer ${process.env.OPENROUTER_API_KEY}`,
        'Content-Type': 'application/json',
      },
    }
  );

  let raw = response.data.choices[0].message.content;
  // Trim whitespace and remove any wrapping backticks or code fences
  raw = raw
    .trim()
    .replace(/^```json\s*/, '')
    .replace(/\s*```$/, '');
  let parsed: Array<{ company: string; category: string }>;
  try {
    parsed = JSON.parse(raw);
  } catch {
    // If direct JSON.parse fails, try stripping leading/trailing newlines/spaces again
    raw = raw.replace(/^[\s\r\n]+|[\s\r\n]+$/g, '');
    parsed = JSON.parse(raw);
  }

  const resultMap: Record<string, string> = {};
  for (const item of parsed) {
    let cat = item.category.trim();
    if (!VALID_CATEGORIES.includes(cat as any)) {
      cat = 'Others';
    }
    resultMap[item.company] = cat;
  }
  return resultMap;
}

async function updateStocksInDatabase(
  pool: pg.Pool,
  batchToUpdate: Array<{ id: number; category: string }>
): Promise<void> {
  try {
    const valuesClauses: string[] = [];
    const params: (number | string)[] = [];

    batchToUpdate.forEach((item, idx) => {
      const baseIdx = idx * 2 + 1;
      valuesClauses.push(`($${baseIdx}::INT, $${baseIdx + 1}::STRING)`);
      params.push(item.id, item.category);
    });

    console.log('Debug - Params:', params);
    console.log('Debug - Values clauses:', valuesClauses);

    const cte = `
      WITH new_values (id, category) AS (
        VALUES ${valuesClauses.join(', ')}
      )
      UPDATE stocks
      SET category = new_values.category
      FROM new_values
      WHERE stocks.id = new_values.id;
    `;

    await pool.query(cte, params);
  } catch (error) {
    console.error('Database update error:');
    console.error('- Error message:', error instanceof Error ? error.message : error);
    throw error; // Re-throw to be handled by the caller
  }
}

async function processBatch(
  pool: pg.Pool,
  batch: Array<{ id: number; company: string; ticker: string }>
): Promise<{ updated: number; skipped: number }> {
  let batchUpdated = 0;
  let batchSkipped = 0;

  try {
    // 1. Extract company names
    const companyNames = batch.map((s) => s.company);
    console.log('Processing companies:', companyNames);

    // 2. Call batch AI
    let categoryMap: Record<string, string>;
    try {
      categoryMap = await getCategoriesFromAI(companyNames);
      console.log('Received categories:', categoryMap);
    } catch (error) {
      console.error('AI categorization error:');
      console.error('- Error message:', error instanceof Error ? error.message : error);
      throw error;
    }

    // 3. Build list for update
    const batchToUpdate: Array<{ id: number; category: string }> = [];

    for (const stock of batch) {
      const predicted = categoryMap[stock.company];
      if (predicted) {
        batchToUpdate.push({ id: stock.id, category: predicted });
      } else {
        batchToUpdate.push({ id: stock.id, category: 'Others' });
      }
      batchUpdated++;
    }

    // 4. Update database
    if (batchToUpdate.length > 0) {
      await updateStocksInDatabase(pool, batchToUpdate);
      console.log(`Successfully updated ${batchToUpdate.length} stocks in database`);
    }

    return { updated: batchUpdated, skipped: batchSkipped };
  } catch (error) {
    console.error('Batch processing error:');
    console.error(
      '- Batch:',
      batch.map((s) => `${s.ticker} (${s.company})`)
    );
    console.error('- Error message:', error instanceof Error ? error.message : error);
    return { updated: 0, skipped: batch.length };
  }
}

async function main() {
  if (!process.env.DB_URL) {
    console.error('DB_URL environment variable is not set');
    process.exit(1);
  }

  if (!process.env.OPENROUTER_API_KEY) {
    console.error('OPENROUTER_API_KEY environment variable is not set');
    process.exit(1);
  }

  const pool = new pg.Pool({
    connectionString: process.env.DB_URL,
    ssl: {
      rejectUnauthorized: false,
    },
    connectionTimeoutMillis: 5000,
    idleTimeoutMillis: 10000,
  });

  pool.on('error', (err) => {
    console.error('Error in pg Pool:', err);
  });

  try {
    // Get all unclassified stocks
    const result = await pool.query('SELECT id, company, ticker FROM stocks WHERE category = $1', ['Unclassified']);

    const stocks = result.rows;
    console.log(`Found ${stocks.length} unclassified stocks`);

    let updatedCount = 0;
    let skippedCount = 0;

    const batches = chunkArray(stocks, 10);

    for (const batch of batches) {
      console.log(`\nProcessing batch of ${batch.length} stocks...`);

      const { updated, skipped } = await processBatch(pool, batch);
      updatedCount += updated;
      skippedCount += skipped;

      // Small delay between batches
      await sleep(500);
    }

    console.log('\nProcessing complete!');
    console.log(`Updated: ${updatedCount} stocks`);
    console.log(`Skipped: ${skippedCount} stocks`);
  } catch (error) {
    console.error('Database error:', error instanceof Error ? error.message : error);
  } finally {
    await pool.end();
  }
}

main().catch((error) => {
  console.error('Fatal error:', error instanceof Error ? error.message : error);
  process.exit(1);
});
