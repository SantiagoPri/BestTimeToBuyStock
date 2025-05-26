import axios from 'axios';
import dotenv from 'dotenv';

dotenv.config();

export interface Stock {
  ticker: string;
  target_from: string;
  target_to: string;
  company: string;
  action: string;
  brokerage: string;
  rating_from: string;
  rating_to: string;
  time: string;
}

interface ApiResponse {
  items: Stock[];
  next_page: string | null;
}

export async function fetchStockData(nextPage?: string): Promise<ApiResponse> {
  const apiUrl = process.env.API_URL;
  const token = process.env.API_TOKEN;
  const url = nextPage ? `${apiUrl}?next_page=${nextPage}` : apiUrl;

  if (!token) {
    throw new Error('API_TOKEN environment variable is not set');
  }

  try {
    const response = await axios.get<ApiResponse>(url!, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    if (error instanceof Error) {
      console.log({ url });
      throw new Error(`Failed to fetch stock data: ${error.message}`);
    }
    throw new Error('Failed to fetch stock data: Unknown error');
  }
}
