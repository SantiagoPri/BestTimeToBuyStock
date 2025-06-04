export const API_BASE_URL = 'http://localhost:8080/api';

export const endpoints = {
  categories: '/categories',
  gmWeek: (week: number) => `/gm/week/${week}`,
  leaderboard: '/leaderboard',
  sessions: '/sessions',
  sessionsAdvance: '/sessions/advance',
  sessionsBuy: '/sessions/buy',
  sessionsEnd: '/sessions/end',
  sessionsSell: '/sessions/sell',
  sessionsState: '/sessions/state',
  stocks: '/stocks',
  stockByParam: (param: string) => `/stocks/${param}`,
} as const;
