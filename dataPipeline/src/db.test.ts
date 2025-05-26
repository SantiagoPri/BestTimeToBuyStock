import { jest, describe, expect, it, beforeEach } from '@jest/globals';
import pg from 'pg';
import { insertStocks, initializeDatabase } from './db.js';
import type { Stock } from './fetchData.js';

type MockClient = {
  query: jest.Mock;
  release: jest.Mock;
};

jest.mock('pg', () => ({
  Pool: jest.fn().mockImplementation(() => ({
    connect: jest.fn().mockResolvedValue({
      query: jest.fn(),
      release: jest.fn(),
    } as MockClient),
  })),
}));

describe('Database operations', () => {
  const mockClient: MockClient = {
    query: jest.fn(),
    release: jest.fn(),
  };

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
    {
      ticker: 'ETR',
      target_from: '$88.00',
      target_to: '$91.00',
      company: 'Entergy',
      action: 'target raised by',
      brokerage: 'Barclays',
      rating_from: 'Overweight',
      rating_to: 'Overweight',
      time: '2025-05-02T00:30:07.507252285Z',
    },
  ];

  beforeEach(() => {
    jest.clearAllMocks();
    const poolMock = pg.Pool as unknown as jest.Mock;
    poolMock.mockImplementation(() => ({
      connect: jest.fn().mockResolvedValue(mockClient),
    }));
  });

  describe('initializeDatabase', () => {
    it('should create tables and indexes', async () => {
      await initializeDatabase();

      expect(mockClient.query).toHaveBeenCalledTimes(2);
      expect(mockClient.query.mock.calls[0][0]).toContain('CREATE TABLE IF NOT EXISTS stock_ratings');
      expect(mockClient.query.mock.calls[1][0]).toContain('CREATE INDEX IF NOT EXISTS idx_stock_ratings_ticker');
      expect(mockClient.release).toHaveBeenCalled();
    });

    it('should handle database errors', async () => {
      const error = new Error('Database connection failed');
      mockClient.query.mockRejectedValueOnce(error);

      await expect(initializeDatabase()).rejects.toThrow('Database connection failed');
      expect(mockClient.release).toHaveBeenCalled();
    });
  });

  describe('insertStocks', () => {
    it('should insert stocks successfully', async () => {
      mockClient.query.mockResolvedValue({ rowCount: 1 });

      await insertStocks(mockStocks);

      expect(mockClient.query).toHaveBeenCalledWith('BEGIN');
      expect(mockClient.query).toHaveBeenCalledWith(
        expect.stringContaining('INSERT INTO stock_ratings'),
        expect.arrayContaining([mockStocks[0].ticker])
      );
      expect(mockClient.query).toHaveBeenCalledWith('COMMIT');
      expect(mockClient.release).toHaveBeenCalled();
    });

    it('should rollback transaction on error', async () => {
      const error = new Error('Insert failed');
      mockClient.query.mockImplementation((query: string) => {
        if (query === 'BEGIN') return Promise.resolve({});
        if (query.includes('INSERT')) return Promise.reject(error);
        return Promise.resolve({});
      });

      await expect(insertStocks(mockStocks)).rejects.toThrow('Failed to insert stock data: Insert failed');
      expect(mockClient.query).toHaveBeenCalledWith('ROLLBACK');
      expect(mockClient.release).toHaveBeenCalled();
    });
  });
});
