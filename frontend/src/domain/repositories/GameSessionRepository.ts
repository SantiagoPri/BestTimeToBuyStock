import type { GameSession, CreateSessionRequest, CreateSessionResponse, TradeRequest } from '../entities/GameSession';

export interface GameSessionRepository {
  getLeaderboard(): Promise<GameSession[]>;
  createSession(request: CreateSessionRequest): Promise<CreateSessionResponse>;
  getSessionState(sessionId: string): Promise<GameSession>;
  buyStocks(sessionId: string, request: TradeRequest): Promise<void>;
  sellStocks(sessionId: string, request: TradeRequest): Promise<void>;
  advanceWeek(sessionId: string): Promise<void>;
  endSession(sessionId: string): Promise<void>;
}
