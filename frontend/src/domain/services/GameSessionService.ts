import type { GameSession, CreateSessionRequest, CreateSessionResponse, TradeRequest } from '../entities/GameSession';
import type { GameSessionRepository } from '../repositories/GameSessionRepository';

export class GameSessionService {
  constructor(private readonly repository: GameSessionRepository) {}

  async getLeaderboard(): Promise<GameSession[]> {
    return this.repository.getLeaderboard();
  }

  async createSession(request: CreateSessionRequest): Promise<CreateSessionResponse> {
    const response = await this.repository.createSession(request);
    localStorage.setItem('sessionId', response.sessionId);
    return response;
  }

  async getSessionState(): Promise<GameSession> {
    const sessionId = localStorage.getItem('sessionId');
    if (!sessionId) {
      throw new Error('No active session');
    }
    return this.repository.getSessionState(sessionId);
  }

  async buyStocks(request: TradeRequest): Promise<void> {
    const sessionId = localStorage.getItem('sessionId');
    if (!sessionId) {
      throw new Error('No active session');
    }
    await this.repository.buyStocks(sessionId, request);
  }

  async sellStocks(request: TradeRequest): Promise<void> {
    const sessionId = localStorage.getItem('sessionId');
    if (!sessionId) {
      throw new Error('No active session');
    }
    await this.repository.sellStocks(sessionId, request);
  }

  async advanceWeek(): Promise<void> {
    const sessionId = localStorage.getItem('sessionId');
    if (!sessionId) {
      throw new Error('No active session');
    }
    await this.repository.advanceWeek(sessionId);
  }

  async endSession(): Promise<void> {
    const sessionId = localStorage.getItem('sessionId');
    if (!sessionId) {
      throw new Error('No active session');
    }
    await this.repository.endSession(sessionId);
    localStorage.removeItem('sessionId');
  }
}
