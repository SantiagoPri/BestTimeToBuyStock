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

interface GameResults {
  cash: number;
  status: string;
  total_balance: number;
  username: string;
}

export class GameSessionApiRepository implements GameSessionRepository {
  constructor(private readonly httpClient: HttpClient) {}

  async getLeaderboard(): Promise<GameSession[]> {
    return this.httpClient.get<GameSession[]>(endpoints.leaderboard);
  }

  async createSession(request: CreateSessionRequest): Promise<CreateSessionResponse> {
    return this.httpClient.post<CreateSessionResponse>(endpoints.sessions, request);
  }

  async getSessionState(sessionId: string): Promise<GameSession> {
    return this.httpClient.get<GameSession>(endpoints.sessionsState, {
      headers: {
        Authorization: `Bearer ${sessionId}`,
      },
    });
  }

  async buyStocks(sessionId: string, request: TradeRequest): Promise<void> {
    await this.httpClient.post(endpoints.sessionsBuy, request, {
      headers: {
        Authorization: `Bearer ${sessionId}`,
      },
    });
  }

  async sellStocks(sessionId: string, request: TradeRequest): Promise<void> {
    await this.httpClient.post(endpoints.sessionsSell, request, {
      headers: {
        Authorization: `Bearer ${sessionId}`,
      },
    });
  }

  async advanceWeek(sessionId: string): Promise<void> {
    await this.httpClient.post(endpoints.sessionsAdvance, {
      headers: {
        Authorization: `Bearer ${sessionId}`,
      },
    });
  }

  async endSession(sessionId: string): Promise<GameResults> {
    return this.httpClient.post<GameResults>(endpoints.sessionsEnd, undefined, {
      headers: {
        Authorization: `Bearer ${sessionId}`,
      },
    });
  }

  async getWeekData(week: number, sessionId: string): Promise<WeekData> {
    return this.httpClient.get<WeekData>(endpoints.gmWeek(week), {
      headers: {
        Authorization: `Bearer ${sessionId}`,
      },
    });
  }
}
