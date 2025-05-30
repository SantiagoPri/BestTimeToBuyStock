import axios from 'axios';
import type { Stock, StockPagination } from '../domain/models/Stock';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:3000',
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface GetStocksParams {
  page?: number;
  limit?: number;
  search?: string;
}

export const stocksApi = {
  getStocks: async (params: GetStocksParams = {}): Promise<StockPagination> => {
    const { data } = await api.get<StockPagination>('/stocks', {
      params: {
        page: params.page || 1,
        limit: params.limit || 10,
        search: params.search,
      },
    });
    return data;
  },

  getStockById: async (id: string): Promise<Stock> => {
    const { data } = await api.get<Stock>(`/stocks/${id}`);
    return data;
  },
};
