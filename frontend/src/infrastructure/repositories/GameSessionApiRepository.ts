import type {
  GameSession,
  CreateSessionRequest,
  CreateSessionResponse,
  TradeRequest,
} from '../../domain/entities/GameSession';
import type { GameSessionRepository } from '../../domain/repositories/GameSessionRepository';
import type { WeekData } from '../../domain/services/GameSessionService';
import { HttpClient } from '../http/HttpClient';
import { endpoints } from '../api/config';

export class GameSessionApiRepository implements GameSessionRepository {
  constructor(private readonly httpClient: HttpClient) {}

  async getLeaderboard(): Promise<GameSession[]> {
    return this.httpClient.get<GameSession[]>(endpoints.leaderboard);
  }

  async createSession(request: CreateSessionRequest): Promise<CreateSessionResponse> {
    return this.httpClient.post<CreateSessionResponse>(endpoints.sessions, request);
  }

  async getSessionState(sessionId: string): Promise<GameSession> {
    return this.httpClient.get<GameSession>(endpoints.sessionsState);
  }

  async buyStocks(sessionId: string, request: TradeRequest): Promise<void> {
    await this.httpClient.post(endpoints.sessionsBuy, request);
  }

  async sellStocks(sessionId: string, request: TradeRequest): Promise<void> {
    await this.httpClient.post(endpoints.sessionsSell, request);
  }

  async advanceWeek(sessionId: string): Promise<void> {
    await this.httpClient.post(endpoints.sessionsAdvance);
  }

  async endSession(sessionId: string): Promise<void> {
    await this.httpClient.post(endpoints.sessionsEnd);
  }

  async getWeekData(week: number, sessionId: string): Promise<WeekData> {
    return this.httpClient.get<WeekData>(endpoints.gmWeek(week), {
      headers: {
        Authorization: `Bearer ${sessionId}`,
      },
    });
  }
}
