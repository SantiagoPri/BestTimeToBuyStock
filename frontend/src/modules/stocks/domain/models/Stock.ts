export type StockRecommendation = 'BUY' | 'SELL' | 'HOLD';

export interface Stock {
  id: string;
  ticker: string;
  company: string;
  currentPrice: number;
  previousClose: number;
  change: number;
  changePercent: number;
  recommendation: StockRecommendation;
  updatedAt: string;
}

export interface StockPagination {
  items: Stock[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}
