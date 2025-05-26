import { fetchStockData } from './fetchData.js';
import { initializeDatabase, insertStocks } from './db.js';
import { setTimeout } from 'timers/promises';

export async function runPipeline() {
  try {
    await initializeDatabase();

    let nextPage: string | undefined = undefined;
    let totalProcessed = 0;

    do {
      const response = await fetchStockData(nextPage);

      if (response.items.length > 0) {
        await insertStocks(response.items);
        totalProcessed += response.items.length;
        console.log(`Processed ${response.items.length} stocks. Total: ${totalProcessed}`);
      }

      nextPage = response.next_page || undefined;

      if (nextPage) {
        await setTimeout(1000);
      }
    } while (nextPage);

    console.log(`Pipeline completed. Total stocks processed: ${totalProcessed}`);
  } catch (error) {
    if (error instanceof Error) {
      console.error('Pipeline error:', error.message);
    } else {
      console.error('Pipeline error:', error);
    }
    process.exit(1);
  }
}

runPipeline();
