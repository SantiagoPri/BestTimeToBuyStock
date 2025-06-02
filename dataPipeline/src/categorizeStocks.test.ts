import { jest, describe, expect, it, beforeEach } from '@jest/globals';
import pg from 'pg';
import axios from 'axios';

// Mock the modules
jest.mock('pg', () => ({
  Pool: jest.fn(),
}));
jest.mock('axios');
jest.mock('dotenv', () => ({
  config: jest.fn(),
}));

// Import after mocking
import './categorizeStocks.js';

interface MockQueryResult {
  rows?: Array<{ id: number; ticker: string; company: string }>;
  rowCount?: number;
}

type QueryFn = (sql: string, values?: any[]) => Promise<MockQueryResult>;

interface MockPool {
  query: jest.MockedFunction<QueryFn>;
  end: jest.MockedFunction<() => Promise<void>>;
}

interface MockApiResponse {
  data: {
    choices: Array<{
      message: {
        content: string;
      };
    }>;
  };
}

type ApiPostFn = (url: string, data: unknown, config: unknown) => Promise<MockApiResponse>;

describe('Stock Categorization', () => {
  const mockPool: MockPool = {
    query: jest.fn() as jest.MockedFunction<QueryFn>,
    end: jest.fn() as jest.MockedFunction<() => Promise<void>>,
  };

  beforeEach(() => {
    jest.clearAllMocks();
    process.env.DB_URL = 'mock-db-url';
    process.env.OPENROUTER_API_KEY = 'mock-api-key';
    (pg.Pool as unknown as jest.Mock).mockImplementation(() => mockPool);
    Object.assign(axios, { post: jest.fn() as jest.MockedFunction<ApiPostFn> });
  });

  it('should update valid categories', async () => {
    // Mock database responses
    const mockRows = [
      { id: 1, ticker: 'AAPL', company: 'Apple Inc' },
      { id: 2, ticker: 'JPM', company: 'JPMorgan Chase' },
    ];
    mockPool.query.mockResolvedValueOnce({ rows: mockRows }).mockResolvedValue({ rowCount: 1 }); // For UPDATE queries

    // Mock OpenRouter API responses
    const mockApiResponse = (category: string): MockApiResponse => ({
      data: { choices: [{ message: { content: category } }] },
    });
    const axiosPostMock = axios.post as jest.MockedFunction<ApiPostFn>;
    axiosPostMock.mockResolvedValueOnce(mockApiResponse('Tech')).mockResolvedValueOnce(mockApiResponse('Finance'));

    await new Promise(process.nextTick); // Wait for main() to complete

    // Verify database queries
    expect(mockPool.query).toHaveBeenCalledWith('SELECT id, company, ticker FROM stocks WHERE category = $1', [
      'Unclassified',
    ]);

    expect(mockPool.query).toHaveBeenCalledWith('UPDATE stocks SET category = $1 WHERE id = $2', ['Tech', 1]);

    expect(mockPool.query).toHaveBeenCalledWith('UPDATE stocks SET category = $1 WHERE id = $2', ['Finance', 2]);

    // Verify API calls
    expect(axiosPostMock).toHaveBeenCalledTimes(2);
    expect(axiosPostMock).toHaveBeenCalledWith(
      'https://openrouter.ai/api/v1/chat/completions',
      expect.objectContaining({
        model: 'mistralai/mistral-7b-instruct',
        messages: expect.arrayContaining([
          expect.objectContaining({
            content: expect.stringContaining('Apple Inc'),
          }),
        ]),
      }),
      expect.any(Object)
    );
  });

  it('should handle invalid categories', async () => {
    // Mock database responses
    const mockRows = [{ id: 1, ticker: 'TEST', company: 'Test Company' }];
    mockPool.query.mockResolvedValueOnce({ rows: mockRows }).mockResolvedValue({ rowCount: 1 }); // For UPDATE queries

    // Mock OpenRouter API response with invalid category
    const axiosPostMock = axios.post as jest.MockedFunction<ApiPostFn>;
    axiosPostMock.mockResolvedValueOnce({
      data: { choices: [{ message: { content: 'Invalid Category' } }] },
    });

    const consoleSpy = jest.spyOn(console, 'warn');

    await new Promise(process.nextTick);

    // Verify warning was logged
    expect(consoleSpy).toHaveBeenCalledWith(expect.stringContaining('Invalid category for TEST'));

    // Verify no update was attempted
    expect(mockPool.query).not.toHaveBeenCalledWith('UPDATE stocks SET category = $1 WHERE id = $2', expect.any(Array));
  });

  it('should handle API errors', async () => {
    // Mock database responses
    const mockRows = [{ id: 1, ticker: 'ERROR', company: 'Error Test Company' }];
    mockPool.query.mockResolvedValueOnce({ rows: mockRows }).mockResolvedValue({ rowCount: 1 }); // For UPDATE queries

    // Mock API error
    const apiError = new Error('API Error');
    const axiosPostMock = axios.post as jest.MockedFunction<ApiPostFn>;
    axiosPostMock.mockRejectedValueOnce(apiError);

    const consoleSpy = jest.spyOn(console, 'error');

    await new Promise(process.nextTick); // Wait for main() to complete

    // Verify error was logged
    expect(consoleSpy).toHaveBeenCalledWith(expect.stringContaining('Error processing ERROR'), 'API Error');

    // Verify no update was attempted
    expect(mockPool.query).not.toHaveBeenCalledWith('UPDATE stocks SET category = $1 WHERE id = $2', expect.any(Array));
  });
});
