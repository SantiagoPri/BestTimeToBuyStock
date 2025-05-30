import { defineStore } from 'pinia';
import type { Stock, StockPagination } from '../domain/models/Stock';
import { stockService } from '../application/fetchStocks';
import type { GetStocksParams } from '../infrastructure/stocksApi';

interface StockState {
  stocks: Stock[];
  selectedStock: Stock | null;
  loading: boolean;
  error: string | null;
  pagination: {
    total: number;
    page: number;
    limit: number;
    totalPages: number;
  };
}

export const useStockStore = defineStore('stock', {
  state: (): StockState => ({
    stocks: [],
    selectedStock: null,
    loading: false,
    error: null,
    pagination: {
      total: 0,
      page: 1,
      limit: 10,
      totalPages: 0,
    },
  }),

  getters: {
    isLoading: (state) => state.loading,
    hasError: (state) => state.error !== null,
  },

  actions: {
    async fetchStocks(params: GetStocksParams = {}) {
      this.loading = true;
      this.error = null;

      try {
        const response = await stockService.fetchStocks(params);
        this.stocks = response.items;
        this.pagination = {
          total: response.total,
          page: response.page,
          limit: response.limit,
          totalPages: response.totalPages,
        };
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'An error occurred';
      } finally {
        this.loading = false;
      }
    },

    async fetchStockById(id: string) {
      this.loading = true;
      this.error = null;

      try {
        this.selectedStock = await stockService.fetchStockById(id);
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'An error occurred';
      } finally {
        this.loading = false;
      }
    },
  },
});
