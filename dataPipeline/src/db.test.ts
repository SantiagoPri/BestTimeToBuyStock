import { jest, describe, expect, it, beforeEach } from '@jest/globals';
import pg from 'pg';
import { insertStocks, initializeDatabase } from './db.js';
import type { Stock } from './fetchData.js';

interface MockQueryResult {
  rows?: Array<{ id: number }>;
  rowCount?: number;
}

type QueryFn = (query: string, ...args: any[]) => Promise<MockQueryResult>;

interface MockClient {
  query: jest.MockedFunction<QueryFn>;
  release: jest.MockedFunction<() => void>;
}

const createMockClient = (): MockClient => ({
  query: jest.fn<QueryFn>(),
  release: jest.fn(),
});

// Mock pg.Pool
const mockConnect = jest.fn();
const mockPool = jest.fn(() => ({
  connect: mockConnect,
}));

jest.mock('pg', () => ({
  Pool: mockPool,
}));

describe('Database operations', () => {
  const mockClient = createMockClient();

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
    mockConnect.mockReturnValue(mockClient);
    mockClient.query.mockImplementation(
      (query: string, ...args: any[]): Promise<MockQueryResult> => Promise.resolve({})
    );
  });

  describe('initializeDatabase', () => {
    it('should create tables and indexes', async () => {
      await initializeDatabase();

      expect(mockClient.query).toHaveBeenCalledTimes(3);
      expect(mockClient.query.mock.calls[0][0]).toContain('CREATE TABLE IF NOT EXISTS stocks');
      expect(mockClient.query.mock.calls[1][0]).toContain('CREATE TABLE IF NOT EXISTS stock_snapshots');
      expect(mockClient.query.mock.calls[2][0]).toContain('CREATE INDEX IF NOT EXISTS');
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
    beforeEach(() => {
      // Mock successful stock insertion with returning ID
      mockClient.query.mockImplementation((query: string, ...args: any[]): Promise<MockQueryResult> => {
        if (query === 'BEGIN') return Promise.resolve({});
        if (query === 'COMMIT') return Promise.resolve({});
        if (query.includes('INSERT INTO stocks')) {
          return Promise.resolve({ rows: [{ id: 1 }] });
        }
        return Promise.resolve({ rowCount: 1 });
      });
    });

    it('should insert new stocks and snapshots successfully', async () => {
      await insertStocks(mockStocks);

      // Check transaction management
      expect(mockClient.query).toHaveBeenCalledWith('BEGIN');
      expect(mockClient.query).toHaveBeenCalledWith('COMMIT');

      // Verify stock insertion
      const stockInsertCalls = mockClient.query.mock.calls.filter(
        (call: unknown[]) => typeof call[0] === 'string' && call[0].includes('INSERT INTO stocks')
      );
      expect(stockInsertCalls).toHaveLength(2);
      expect(stockInsertCalls[0][1]).toEqual([mockStocks[0].ticker, mockStocks[0].company, 'Unclassified']);

      // Verify snapshot insertion
      const snapshotInsertCalls = mockClient.query.mock.calls.filter(
        (call: unknown[]) => typeof call[0] === 'string' && call[0].includes('INSERT INTO stock_snapshots')
      );
      expect(snapshotInsertCalls).toHaveLength(2);

      // Check price calculation for first stock
      const firstSnapshotParams = snapshotInsertCalls[0][1] as unknown[];
      const targetFrom = parseFloat(mockStocks[0].target_from.replace('$', ''));
      const targetTo = parseFloat(mockStocks[0].target_to.replace('$', ''));
      const price = firstSnapshotParams[6] as number;
      const avgTarget = (targetFrom + targetTo) / 2;

      expect(price).toBeGreaterThanOrEqual(avgTarget * 0.85);
      expect(price).toBeLessThanOrEqual(avgTarget * 1.1);
      expect(firstSnapshotParams[1]).toBe(1); // week
      expect(firstSnapshotParams[8]).toBe(`Recommendation by ${mockStocks[0].brokerage}`); // news_title
    });

    it('should handle duplicate stocks correctly', async () => {
      // Mock the first stock already existing
      const duplicateStock = { ...mockStocks[0] };
      mockClient.query.mockImplementation((query: string, ...args: any[]): Promise<MockQueryResult> => {
        if (query === 'BEGIN') return Promise.resolve({});
        if (query === 'COMMIT') return Promise.resolve({});
        if (query.includes('INSERT INTO stocks')) {
          return Promise.resolve({ rows: [{ id: 1 }] });
        }
        return Promise.resolve({ rowCount: 1 });
      });

      await insertStocks([duplicateStock]);
      await insertStocks([duplicateStock]);

      const stockInsertCalls = mockClient.query.mock.calls.filter(
        (call: unknown[]) => typeof call[0] === 'string' && call[0].includes('INSERT INTO stocks')
      );

      // Should still try to insert twice (with ON CONFLICT DO UPDATE)
      expect(stockInsertCalls).toHaveLength(2);
      // Both calls should have the same parameters
      expect(stockInsertCalls[0][1]).toEqual(stockInsertCalls[1][1]);
    });

    it('should rollback transaction on error', async () => {
      const error = new Error('Insert failed');
      mockClient.query.mockImplementation((query: string, ...args: any[]): Promise<MockQueryResult> => {
        if (query === 'BEGIN') return Promise.resolve({});
        if (query.includes('INSERT INTO stocks')) {
          return Promise.reject(error);
        }
        return Promise.resolve({});
      });

      await expect(insertStocks(mockStocks)).rejects.toThrow('Failed to insert stock data: Insert failed');
      expect(mockClient.query).toHaveBeenCalledWith('ROLLBACK');
      expect(mockClient.release).toHaveBeenCalled();
    });
  });
});
