import { Pool } from 'pg';
import dotenv from 'dotenv';

dotenv.config();

async function syncStockRatingsCategory() {
  const pool = new Pool({
    connectionString: process.env.DB_URL,
    ssl: {
      rejectUnauthorized: false,
    },
  });

  try {
    const ratingsResult = await pool.query<{ id: number; ticker: string }>(`
      SELECT id, ticker
      FROM stock_ratings;
    `);

    const totalRecords = ratingsResult.rows.length;
    console.log(`Starting migration of categories for ${totalRecords} stock ratings...`);

    let processedCount = 0;
    let updatedCount = 0;
    let skippedCount = 0;

    for (const { id, ticker } of ratingsResult.rows) {
      processedCount++;

      const stockResult = await pool.query<{ category: string }>(
        `
        SELECT category
        FROM stocks
        WHERE ticker = $1;
      `,
        [ticker]
      );

      if (stockResult.rowCount === 1) {
        const category = stockResult.rows[0].category;
        await pool.query(
          `
          UPDATE stock_ratings
          SET category = $1
          WHERE id = $2;
        `,
          [category, id]
        );
        updatedCount++;

        // Log progress every 100 records or when it's the last record
        if (processedCount % 100 === 0 || processedCount === totalRecords) {
          console.log(
            `Progress: ${processedCount}/${totalRecords} (${Math.round(
              (processedCount / totalRecords) * 100
            )}%) - Updated: ${updatedCount}, Skipped: ${skippedCount}`
          );
        }
      } else {
        console.warn(`No matching stock found for ticker "${ticker}" (rating_id=${id}), skipping update.`);
        skippedCount++;
      }
    }

    console.log('\nMigration Summary:');
    console.log(`Total Processed: ${processedCount}`);
    console.log(`Successfully Updated: ${updatedCount}`);
    console.log(`Skipped: ${skippedCount}`);
    console.log('Category migration completed.');
  } catch (error) {
    console.error('Error while syncing categories:', error);
  } finally {
    await pool.end();
  }
}

syncStockRatingsCategory();
