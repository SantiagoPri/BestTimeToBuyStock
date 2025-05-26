import { jest, describe, expect, it, beforeEach } from '@jest/globals';
import { fetchStockData } from './fetchData.js';
import { initializeDatabase, insertStocks } from './db.js';
import { setTimeout } from 'timers/promises';
import type { Stock } from './fetchData.js';

jest.mock('./fetchData.js');
jest.mock('./db.js');
jest.mock('timers/promises');

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
    (fetchStockData as unknown as jest.Mock).mockReset();
    (initializeDatabase as unknown as jest.Mock).mockReset();
    (insertStocks as unknown as jest.Mock).mockReset();
    (setTimeout as unknown as jest.Mock).mockReset();
  });

  it('should process all pages of stock data', async () => {
    // Mock first page with next_page token
    (fetchStockData as unknown as jest.Mock).mockResolvedValueOnce({
      items: mockStocks.slice(0, 1),
      next_page: 'page2',
    });

    // Mock second page with no next_page token
    (fetchStockData as unknown as jest.Mock).mockResolvedValueOnce({
      items: mockStocks.slice(1),
      next_page: null,
    });

    (initializeDatabase as unknown as jest.Mock).mockResolvedValue(undefined);
    (insertStocks as unknown as jest.Mock).mockResolvedValue(undefined);
    (setTimeout as unknown as jest.Mock).mockResolvedValue(undefined);

    // Import the runPipeline function dynamically to avoid hoisting issues with mocks
    const { runPipeline } = await import('./pipeline.js');
    await runPipeline();

    expect(initializeDatabase).toHaveBeenCalledTimes(1);
    expect(fetchStockData).toHaveBeenCalledTimes(2);
    expect(fetchStockData).toHaveBeenNthCalledWith(1, undefined);
    expect(fetchStockData).toHaveBeenNthCalledWith(2, 'page2');
    expect(insertStocks).toHaveBeenCalledTimes(2);
    expect(setTimeout).toHaveBeenCalledWith(1000);
  });

  it('should handle API errors', async () => {
    const error = new Error('API error');
    (fetchStockData as unknown as jest.Mock).mockRejectedValue(error);
    (initializeDatabase as unknown as jest.Mock).mockResolvedValue(undefined);

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
    (fetchStockData as unknown as jest.Mock).mockResolvedValue({
      items: mockStocks,
      next_page: null,
    });
    (initializeDatabase as unknown as jest.Mock).mockRejectedValue(error);

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
