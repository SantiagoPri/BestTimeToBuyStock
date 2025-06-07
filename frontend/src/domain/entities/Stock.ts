export type StockRecommendation = 'BUY' | 'SELL' | 'HOLD';

export interface Stock {
  ticker: string;
  company: string;
  currentPrice: number;
  changePercent: string;
  ratings: string;
  marketSentiment: 'up' | 'down' | 'neutral';
  change: number;
  // recommendation: StockRecommendation;
}

export interface StockPagination {
  items: Stock[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

export function getMarketSentiment(action: string) {
  switch (action.toLowerCase()) {
    case 'upgraded':
    case 'target raised':
      return 'up';
    case 'downgraded':
    case 'target lowered':
      return 'down';
    case 'reiterated':
    case 'neutral':
      return 'neutral';
    default:
      console.log('Unknown market sentiment action:', action);
      return 'neutral';
  }
}
