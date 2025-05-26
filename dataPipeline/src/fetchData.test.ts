/// <reference types="jest" />

import { jest, describe, expect, it, beforeEach } from '@jest/globals';
import axios from 'axios';
import { fetchStockData, Stock } from './fetchData.js';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('fetchStockData', () => {
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
    process.env.API_URL = 'https://api.example.com/stocks';
    process.env.API_TOKEN = 'test-token';
    jest.clearAllMocks();
  });

  it('should fetch stock data successfully', async () => {
    const mockResponse = {
      items: mockStocks,
      next_page: 'next-page-token',
    };

    mockedAxios.get.mockResolvedValueOnce({ data: mockResponse });

    const result = await fetchStockData();

    expect(mockedAxios.get).toHaveBeenCalledWith('https://api.example.com/stocks', {
      headers: {
        Authorization: 'Bearer test-token',
      },
    });
    expect(result).toEqual(mockResponse);
  });

  it('should handle pagination correctly', async () => {
    const mockResponse = {
      items: mockStocks,
      next_page: null,
    };

    mockedAxios.get.mockResolvedValueOnce({ data: mockResponse });

    const result = await fetchStockData('page-token');

    expect(mockedAxios.get).toHaveBeenCalledWith('https://api.example.com/stocks?next_page=page-token', {
      headers: {
        Authorization: 'Bearer test-token',
      },
    });
    expect(result).toEqual(mockResponse);
  });

  it('should throw error when API_TOKEN is not set', async () => {
    process.env.API_TOKEN = '';

    await expect(fetchStockData()).rejects.toThrow('API_TOKEN environment variable is not set');
  });

  it('should handle API error', async () => {
    mockedAxios.get.mockRejectedValueOnce(new Error('Network error'));

    await expect(fetchStockData()).rejects.toThrow('Failed to fetch stock data: Network error');
  });
});
