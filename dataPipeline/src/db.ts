import pg from 'pg';
import dotenv from 'dotenv';
import type { Stock } from './fetchData.js';

dotenv.config();

const { Pool } = pg;

// Create a connection pool
const pool = new Pool({
  connectionString: process.env.DB_URL,
  ssl: {
    rejectUnauthorized: false, // Required for CockroachDB Cloud
  },
});

export async function initializeDatabase(): Promise<void> {
  const client = await pool.connect();
  try {
    await client.query(`
      CREATE TABLE IF NOT EXISTS stock_data (
        symbol TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        price DECIMAL NOT NULL,
        volume BIGINT NOT NULL,
        percent_change DECIMAL NOT NULL,
        market_cap DECIMAL NOT NULL,
        sector TEXT NOT NULL,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
      )
    `);
  } finally {
    client.release();
  }
}

export async function insertStocks(stocks: Stock[]): Promise<void> {
  const client = await pool.connect();
  try {
    // Begin transaction
    await client.query('BEGIN');

    for (const stock of stocks) {
      await client.query(
        `
        INSERT INTO stock_data (symbol, name, price, volume, percent_change, market_cap, sector)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (symbol) 
        DO UPDATE SET 
          name = EXCLUDED.name,
          price = EXCLUDED.price,
          volume = EXCLUDED.volume,
          percent_change = EXCLUDED.percent_change,
          market_cap = EXCLUDED.market_cap,
          sector = EXCLUDED.sector,
          updated_at = CURRENT_TIMESTAMP
        `,
        [stock.symbol, stock.name, stock.price, stock.volume, stock.percent_change, stock.market_cap, stock.sector]
      );
    }

    // Commit transaction
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
