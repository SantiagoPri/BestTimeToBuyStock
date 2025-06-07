import type { GameSession, CreateSessionRequest, CreateSessionResponse, TradeRequest } from '../entities/GameSession';
import type { WeekData } from '../services/GameSessionService';

interface GameResults {
  cash: number;
  status: string;
  total_balance: number;
  username: string;
}

export interface GameSessionRepository {
  getLeaderboard(): Promise<GameSession[]>;
  createSession(request: CreateSessionRequest): Promise<CreateSessionResponse>;
  getSessionState(sessionId: string): Promise<GameSession>;
  buyStocks(sessionId: string, request: TradeRequest): Promise<void>;
  sellStocks(sessionId: string, request: TradeRequest): Promise<void>;
  advanceWeek(sessionId: string): Promise<void>;
  endSession(sessionId: string): Promise<GameResults>;
  getWeekData(week: number, sessionId: string): Promise<WeekData>;
}
