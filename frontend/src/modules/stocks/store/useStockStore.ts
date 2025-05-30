import { defineStore } from 'pinia';
import axios from 'axios';

interface Stock {
  Ticker: string;
  Company: string;
  TargetFrom: string;
  TargetTo: string;
  Action: string;
  Brokerage: string;
  RatingFrom: string;
  RatingTo: string;
  Time: string;
}

interface StockApiResponse {
  currentPage: number;
  limit: number;
  stocks: Stock[];
  total: number;
}

interface StockState {
  stocks: Stock[];
  loading: boolean;
  error: string | null;
  pagination: {
    total: number;
    page: number;
    limit: number;
  };
}

interface GetStocksParams {
  page?: number;
  limit?: number;
  search?: string;
}

export const useStockStore = defineStore('stocks', {
  state: (): StockState => ({
    stocks: [],
    loading: false,
    error: null,
    pagination: {
      total: 0,
      page: 0,
      limit: 10,
    },
  }),

  getters: {
    currentPage: (state) => state.pagination.page,
    total: (state) => state.pagination.total,
    limit: (state) => state.pagination.limit,
  },

  actions: {
    async fetchStocks(page: number) {
      this.loading = true;
      this.error = null;

      try {
        const response = await axios.get<StockApiResponse>(`${import.meta.env.VITE_API_URL}/stocks`, {
          params: {
            page,
            limit: this.pagination.limit,
          },
        });

        console.log('API Response:', response.data); // Debug log
        this.stocks = response.data.stocks || [];
        this.pagination = {
          total: response.data.total,
          page: response.data.currentPage,
          limit: response.data.limit,
        };
      } catch (error) {
        console.error('API Error:', error); // Debug log
        this.error = error instanceof Error ? error.message : 'An error occurred';
        this.stocks = [];
        this.pagination.total = 0;
      } finally {
        this.loading = false;
      }
    },
  },
});
