<script setup lang="ts">
interface Stock {
  Ticker: string;
  Company: string;
  TargetFrom: string;
  TargetTo: string;
  Action: string;
  Brokerage: string;
  RatingFrom: string;
  RatingTo: string;
  Time: string;
}

interface Props {
  stocks: Stock[];
  loading?: boolean;
}

const props = defineProps<Props>();

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
};
</script>

<template>
  <div class="w-full overflow-x-auto">
    <table class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50 sticky top-0">
        <tr>
          <th
            v-for="header in ['Ticker', 'Company', 'Target Price', 'Rating', 'Action', 'Brokerage', 'Date']"
            :key="header"
            class="px-4 py-3 text-sm font-semibold text-gray-700 text-left border-b"
          >
            {{ header }}
          </th>
        </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
        <template v-if="loading">
          <tr>
            <td
              colspan="7"
              class="px-4 py-3 text-sm text-gray-500 text-left"
            >
              Loading...
            </td>
          </tr>
        </template>
        <template v-else-if="stocks.length === 0">
          <tr>
            <td
              colspan="7"
              class="px-4 py-3 text-sm text-gray-500 text-left"
            >
              No stocks found
            </td>
          </tr>
        </template>
        <template v-else>
          <tr
            v-for="(stock, index) in stocks"
            :key="stock.Ticker + stock.Time"
            :class="{
              'bg-gray-50': index % 2 === 0,
              'bg-white': index % 2 === 1
            }"
          >
            <td class="px-4 py-3 text-sm text-left font-medium text-indigo-600">
              {{ stock.Ticker }}
            </td>
            <td class="px-4 py-3 text-sm text-left">
              {{ stock.Company }}
            </td>
            <td class="px-4 py-3 text-sm text-left whitespace-nowrap">
              {{ stock.TargetFrom === stock.TargetTo 
                ? stock.TargetFrom 
                : `${stock.TargetFrom} → ${stock.TargetTo}` }}
            </td>
            <td class="px-4 py-3 text-sm text-left whitespace-nowrap">
              {{ stock.RatingFrom === stock.RatingTo 
                ? stock.RatingFrom 
                : `${stock.RatingFrom} → ${stock.RatingTo}` }}
            </td>
            <td class="px-4 py-3 text-sm text-left">
              {{ stock.Action }}
            </td>
            <td class="px-4 py-3 text-sm text-left">
              {{ stock.Brokerage }}
            </td>
            <td class="px-4 py-3 text-sm text-left whitespace-nowrap">
              {{ formatDate(stock.Time) }}
            </td>
          </tr>
        </template>
      </tbody>
    </table>
  </div>
</template> 