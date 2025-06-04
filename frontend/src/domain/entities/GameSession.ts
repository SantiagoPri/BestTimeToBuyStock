export type GameSessionStatus = 'starting' | 'week1' | 'week2' | 'week3' | 'week4' | 'week5' | 'finished' | 'expired';

export interface HoldingInfo {
  quantity: number;
  total_spent: number;
}

export interface SessionMetadata {
  holdings: Record<string, HoldingInfo>;
}

export interface GameSession {
  session_id: string;
  username: string;
  cash: number;
  holdings_value: number;
  total_balance: number;
  status: GameSessionStatus;
  metadata: SessionMetadata;
  created_at: string;
  updated_at: string;
}

export interface CreateSessionRequest {
  username: string;
  categories: string[];
}

export interface CreateSessionResponse {
  sessionId: string;
}

export interface TradeRequest {
  ticker: string;
  quantity: number;
}
