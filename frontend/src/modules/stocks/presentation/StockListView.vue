<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import { useStockStore } from '../store/useStockStore';
import { usePagination } from '@/shared/composables/usePagination';
import type { Stock } from '../domain/models/Stock';

const store = useStockStore();
const searchQuery = ref('');
const { currentPage, itemsPerPage } = usePagination();

const fetchStocksData = async () => {
  await store.fetchStocks({
    page: currentPage.value,
    limit: itemsPerPage.value,
    search: searchQuery.value,
  });
};

watch([currentPage, itemsPerPage, searchQuery], () => {
  fetchStocksData();
});

onMounted(() => {
  fetchStocksData();
});

const formatPrice = (price: number) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
  }).format(price);
};

const getChangeColor = (change: number) => {
  return change >= 0 ? 'text-green-600' : 'text-red-600';
};
</script>

<template>
  <div class="p-6">
    <div class="mb-6">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search by company or ticker..."
        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
      />
    </div>

    <div v-if="store.loading" class="text-center py-8">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
    </div>

    <div v-else-if="store.error" class="text-center py-8 text-red-600">
      {{ store.error }}
    </div>

    <div v-else>
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Ticker</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Company</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Price</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Change</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Recommendation
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="stock in store.stocks" :key="stock.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ stock.ticker }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ stock.company }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ formatPrice(stock.currentPrice) }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm" :class="getChangeColor(stock.change)">
                {{ stock.change > 0 ? '+' : '' }}{{ stock.changePercent.toFixed(2) }}%
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <span
                  class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                  :class="{
                    'bg-green-100 text-green-800': stock.recommendation === 'BUY',
                    'bg-red-100 text-red-800': stock.recommendation === 'SELL',
                    'bg-yellow-100 text-yellow-800': stock.recommendation === 'HOLD',
                  }"
                >
                  {{ stock.recommendation }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="mt-4 flex items-center justify-between">
        <div class="flex-1 flex justify-between sm:hidden">
          <button
            @click="currentPage--"
            :disabled="currentPage === 1"
            class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
            :class="{ 'opacity-50 cursor-not-allowed': currentPage === 1 }"
          >
            Previous
          </button>
          <button
            @click="currentPage++"
            :disabled="currentPage === store.pagination.totalPages"
            class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
            :class="{ 'opacity-50 cursor-not-allowed': currentPage === store.pagination.totalPages }"
          >
            Next
          </button>
        </div>
        <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
          <div>
            <p class="text-sm text-gray-700">
              Showing
              <span class="font-medium">{{ (currentPage - 1) * itemsPerPage + 1 }}</span>
              to
              <span class="font-medium">
                {{ Math.min(currentPage * itemsPerPage, store.pagination.total) }}
              </span>
              of
              <span class="font-medium">{{ store.pagination.total }}</span>
              results
            </p>
          </div>
          <div>
            <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
              <button
                v-for="page in store.pagination.totalPages"
                :key="page"
                @click="currentPage = page"
                :class="[
                  page === currentPage
                    ? 'z-10 bg-primary-50 border-primary-500 text-primary-600'
                    : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50',
                  'relative inline-flex items-center px-4 py-2 border text-sm font-medium',
                ]"
              >
                {{ page }}
              </button>
            </nav>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
