export const API_BASE_URL = import.meta.env.API_URL || 'http://localhost:8080';
console.log('API_BASE_URL:', API_BASE_URL);

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
