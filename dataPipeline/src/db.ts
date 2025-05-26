import pg from 'pg';
import dotenv from 'dotenv';
import type { Stock } from './fetchData.js';

dotenv.config();

const { Pool } = pg;

// Create a connection pool
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
      CREATE TABLE IF NOT EXISTS stock_ratings (
        id SERIAL PRIMARY KEY,
        ticker TEXT NOT NULL,
        target_from TEXT NOT NULL,
        target_to TEXT NOT NULL,
        company TEXT NOT NULL,
        action TEXT NOT NULL,
        brokerage TEXT NOT NULL,
        rating_from TEXT NOT NULL,
        rating_to TEXT NOT NULL,
        time TIMESTAMP NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
      )
    `);

    // Create an index on ticker for faster lookups
    await client.query(`
      CREATE INDEX IF NOT EXISTS idx_stock_ratings_ticker ON stock_ratings(ticker)
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
        INSERT INTO stock_ratings (
          ticker, target_from, target_to, company, action, 
          brokerage, rating_from, rating_to, time
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        `,
        [
          stock.ticker,
          stock.target_from,
          stock.target_to,
          stock.company,
          stock.action,
          stock.brokerage,
          stock.rating_from,
          stock.rating_to,
          stock.time,
        ]
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
