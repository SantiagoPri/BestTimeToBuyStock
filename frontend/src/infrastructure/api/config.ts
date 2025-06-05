export const API_BASE_URL = 'http://localhost:8080/api';

export const endpoints = {
  categories: '/categories',
  gmWeek: (week: number) => `/game/week/${week}`,
  leaderboard: '/leaderboard',
  sessions: '/session/start',
  sessionsAdvance: '/session/advance',
  sessionsBuy: '/session/buy',
  sessionsEnd: '/session/end',
  sessionsSell: '/session/sell',
  sessionsState: '/session/state',
  stocks: '/stocks',
  stockByParam: (param: string) => `/stocks/${param}`,
} as const;
