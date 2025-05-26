import axios from 'axios';
import dotenv from 'dotenv';

dotenv.config();

export interface Stock {
  symbol: string;
  name: string;
  price: number;
  volume: number;
  percent_change: number;
  market_cap: number;
  sector: string;
}

export async function fetchStockData(): Promise<Stock[]> {
  try {
    const apiUrl = process.env.API_URL || 'https://api.example.com/stocks';
    const response = await axios.get(apiUrl);

    // Assuming the API returns an array of stock data
    // Map the response data to our Stock interface
    const stocks: Stock[] = response.data.map((item: any) => ({
      symbol: item.symbol,
      name: item.name,
      price: Number(item.price),
      volume: Number(item.volume),
      percent_change: Number(item.percent_change),
      market_cap: Number(item.market_cap),
      sector: item.sector,
    }));

    return stocks;
  } catch (error) {
    if (error instanceof Error) {
      throw new Error(`Failed to fetch stock data: ${error.message}`);
    }
    throw new Error('Failed to fetch stock data: Unknown error');
  }
}
