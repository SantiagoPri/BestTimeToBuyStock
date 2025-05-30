import type { StockPagination } from '../domain/models/Stock';
import { stocksApi, type GetStocksParams } from '../infrastructure/stocksApi';

export class StockService {
  async fetchStocks(params: GetStocksParams = {}): Promise<StockPagination> {
    try {
      return await stocksApi.getStocks(params);
    } catch (error) {
      console.error('Error fetching stocks:', error);
      throw new Error('Failed to fetch stocks');
    }
  }

  async fetchStockById(id: string) {
    try {
      return await stocksApi.getStockById(id);
    } catch (error) {
      console.error('Error fetching stock:', error);
      throw new Error('Failed to fetch stock');
    }
  }
}

export const stockService = new StockService();
