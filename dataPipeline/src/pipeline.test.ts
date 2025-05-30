import { jest, describe, expect, it, beforeEach } from '@jest/globals';
import { fetchStockData } from './fetchData.js';
import { initializeDatabase, insertStocks } from './db.js';
import { setTimeout } from 'timers/promises';
import type { Stock } from './fetchData.js';

jest.mock('./fetchData.js');
jest.mock('./db.js');
jest.mock('timers/promises');

interface ApiResponse {
  items: Stock[];
  next_page: string | null;
}

type FetchStockDataFn = typeof fetchStockData;
type InitializeDatabaseFn = typeof initializeDatabase;
type InsertStocksFn = typeof insertStocks;
type SetTimeoutFn = typeof setTimeout;

describe('Pipeline', () => {
  const mockStocks: Stock[] = [
    {
      ticker: 'INSG',
      target_from: '$11.00',
      target_to: '$8.00',
      company: 'Inseego',
      action: 'target lowered by',
      brokerage: 'Stifel Nicolaus',
      rating_from: 'Hold',
      rating_to: 'Hold',
      time: '2025-04-18T00:30:07.845283435Z',
    },
    {
      ticker: 'MODG',
      target_from: '$12.00',
      target_to: '$7.00',
      company: 'Topgolf Callaway Brands',
      action: 'target lowered by',
      brokerage: 'Truist Financial',
      rating_from: 'Buy',
      rating_to: 'Buy',
      time: '2025-04-15T00:30:13.351058975Z',
    },
  ];

  beforeEach(() => {
    jest.clearAllMocks();
    (fetchStockData as jest.MockedFunction<FetchStockDataFn>).mockReset();
    (initializeDatabase as jest.MockedFunction<InitializeDatabaseFn>).mockReset();
    (insertStocks as jest.MockedFunction<InsertStocksFn>).mockReset();
    (setTimeout as jest.MockedFunction<SetTimeoutFn>).mockReset();
  });

  it('should process all pages of stock data', async () => {
    // Mock first page with next_page token
    const mockFetchStockData = fetchStockData as jest.MockedFunction<FetchStockDataFn>;
    mockFetchStockData.mockResolvedValueOnce({
      items: mockStocks.slice(0, 1),
      next_page: 'page2',
    });

    // Mock second page with no next_page token
    mockFetchStockData.mockResolvedValueOnce({
      items: mockStocks.slice(1),
      next_page: null,
    });

    const mockInitializeDatabase = initializeDatabase as jest.MockedFunction<InitializeDatabaseFn>;
    const mockInsertStocks = insertStocks as jest.MockedFunction<InsertStocksFn>;
    const mockSetTimeout = setTimeout as jest.MockedFunction<SetTimeoutFn>;

    mockInitializeDatabase.mockResolvedValue(undefined);
    mockInsertStocks.mockResolvedValue(undefined);
    mockSetTimeout.mockResolvedValue(undefined);

    // Import the runPipeline function dynamically to avoid hoisting issues with mocks
    const { runPipeline } = await import('./pipeline.js');
    await runPipeline();

    expect(mockInitializeDatabase).toHaveBeenCalledTimes(1);
    expect(mockFetchStockData).toHaveBeenCalledTimes(2);
    expect(mockFetchStockData).toHaveBeenNthCalledWith(1, undefined);
    expect(mockFetchStockData).toHaveBeenNthCalledWith(2, 'page2');
    expect(mockInsertStocks).toHaveBeenCalledTimes(2);
    expect(mockSetTimeout).toHaveBeenCalledWith(1000);
  });

  it('should handle API errors', async () => {
    const error = new Error('API error');
    const mockFetchStockData = fetchStockData as jest.MockedFunction<FetchStockDataFn>;
    const mockInitializeDatabase = initializeDatabase as jest.MockedFunction<InitializeDatabaseFn>;

    mockFetchStockData.mockRejectedValue(error);
    mockInitializeDatabase.mockResolvedValue(undefined);

    const mockExit = jest.spyOn(process, 'exit').mockImplementation(() => undefined as never);
    const mockConsoleError = jest.spyOn(console, 'error').mockImplementation(() => {});

    const { runPipeline } = await import('./pipeline.js');
    await runPipeline();

    expect(mockConsoleError).toHaveBeenCalledWith('Pipeline error:', error.message);
    expect(mockExit).toHaveBeenCalledWith(1);

    mockExit.mockRestore();
    mockConsoleError.mockRestore();
  });

  it('should handle database errors', async () => {
    const error = new Error('Database error');
    const mockFetchStockData = fetchStockData as jest.MockedFunction<FetchStockDataFn>;
    const mockInitializeDatabase = initializeDatabase as jest.MockedFunction<InitializeDatabaseFn>;

    mockFetchStockData.mockResolvedValue({
      items: mockStocks,
      next_page: null,
    });
    mockInitializeDatabase.mockRejectedValue(error);

    const mockExit = jest.spyOn(process, 'exit').mockImplementation(() => undefined as never);
    const mockConsoleError = jest.spyOn(console, 'error').mockImplementation(() => {});

    const { runPipeline } = await import('./pipeline.js');
    await runPipeline();

    expect(mockConsoleError).toHaveBeenCalledWith('Pipeline error:', error.message);
    expect(mockExit).toHaveBeenCalledWith(1);

    mockExit.mockRestore();
    mockConsoleError.mockRestore();
  });
});
