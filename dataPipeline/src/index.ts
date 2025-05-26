import { fetchStockData } from './fetchData.js';
import { initializeDatabase, insertStocks } from './db.js';
import dotenv from 'dotenv';

dotenv.config();

async function main() {
  try {
    console.log('Initializing database...');
    await initializeDatabase();

    console.log('Fetching stock data...');
    const stocks = await fetchStockData();

    console.log(`Retrieved ${stocks.length} stocks. Inserting into database...`);
    await insertStocks(stocks);

    console.log('Data pipeline completed successfully!');
  } catch (error) {
    console.error('Error in data pipeline:', error instanceof Error ? error.message : 'Unknown error');
    process.exit(1);
  }
}

// Run the pipeline
main();
