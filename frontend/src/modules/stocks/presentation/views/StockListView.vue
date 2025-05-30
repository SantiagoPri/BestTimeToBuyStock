<script setup lang="ts">
import { onMounted, computed } from 'vue';
import { useStockStore } from '../../store/useStockStore';
import StockTable from '../components/StockTable.vue';

const stockStore = useStockStore();

const currentPage = computed(() => stockStore.currentPage);
const totalPages = computed(() => Math.ceil(stockStore.total / stockStore.limit));

const goToNextPage = async () => {
  if (currentPage.value * stockStore.limit < stockStore.total) {
    await stockStore.fetchStocks(currentPage.value + 1);
  }
};

const goToPreviousPage = async () => {
  if (currentPage.value > 1) {
    await stockStore.fetchStocks(currentPage.value - 1);
  }
};

onMounted(async () => {
  await stockStore.fetchStocks(1);
});
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <h1 class="text-3xl font-bold text-center text-gray-900 mb-8">
      Stock Market Overview
    </h1>

    <StockTable
      :stocks="stockStore.stocks"
      :loading="stockStore.loading"
    />

    <div class="mt-6 flex flex-col items-center space-y-4">
      <div class="flex items-center gap-2">
        <button
          @click="goToPreviousPage"
          :disabled="currentPage === 1"
          class="px-4 py-2 text-sm font-medium rounded-md
                 transition-colors duration-150 ease-in-out
                 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2
                 disabled:bg-gray-300 disabled:cursor-not-allowed
                 enabled:bg-indigo-500 enabled:hover:bg-indigo-600 enabled:text-white"
        >
          ← Previous
        </button>
        <button
          @click="goToNextPage"
          :disabled="currentPage * stockStore.limit >= stockStore.total"
          class="px-4 py-2 text-sm font-medium rounded-md
                 transition-colors duration-150 ease-in-out
                 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2
                 disabled:bg-gray-300 disabled:cursor-not-allowed
                 enabled:bg-indigo-500 enabled:hover:bg-indigo-600 enabled:text-white"
        >
          Next →
        </button>
      </div>

      <div class="text-sm text-gray-600">
        Page {{ currentPage }} of {{ totalPages }}
      </div>
    </div>
  </div>
</template> 