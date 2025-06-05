<template>
  <div class="min-h-screen bg-gray-900 bg-[url('/images/card-home/background-stars.png')] bg-cover">
    <!-- Top Bar -->
    <header class="flex justify-between items-center px-16 py-6">
      <div class="flex items-center gap-4">
        <img src="/images/card-home/drakeHead.png" alt="Avatar" class="w-12 h-12 rounded-full" />
        <div class="flex flex-col">
          <h2 class="text-2xl font-bold text-gray-100">Username</h2>
          <span class="text-gray-400">Week 1/5</span>
        </div>
      </div>
      <div class="flex gap-6">
        <div class="bg-gray-900 rounded p-2 flex items-center gap-2 w-44">
          <WalletIcon class="w-6 h-6 text-gray-400" />
          <div>
            <div class="text-sm text-gray-100">Balance</div>
            <div class="text-xs text-gray-400">$10,000 USD</div>
          </div>
        </div>
        <div class="bg-gray-900 rounded p-2 flex items-center gap-2 w-44">
          <CircleStackIcon class="w-6 h-6 text-gray-400" />
          <div>
            <div class="text-sm text-gray-100">Stocks</div>
            <div class="text-xs text-gray-400">$0.00 USD</div>
          </div>
        </div>
        <div class="bg-gray-900 rounded p-2 flex items-center gap-2 w-44">
          <CurrencyDollarIcon class="w-6 h-6 text-gray-400" />
          <div>
            <div class="text-sm text-gray-100">Total</div>
            <div class="text-xs text-green-500">$10,000 USD</div>
          </div>
        </div>
      </div>
    </header>

    <!-- Market News -->
    <section class="px-16 py-12">
      <h2 class="text-2xl font-medium font-orbitron text-gray-400 mb-8">Market News</h2>
      <div class="grid grid-cols-3 gap-8">
        <div v-for="i in 3" :key="i" class="bg-gray-800 rounded-lg shadow-sm p-3">
          <p class="text-gray-100">Lorem ipsum dolor sit amet consectetur. Leo blandit vehicula velit non ut.</p>
        </div>
      </div>
    </section>

    <!-- Game Master's Tip -->
    <section class="px-16 py-12">
      <div class="border-t border-gray-700"></div>
      <h2 class="text-2xl font-medium font-orbitron text-gray-400 my-8">Game Master's Tip for This Week</h2>
      <div class="bg-blue-900 bg-opacity-50 rounded-lg p-6">
        <div class="flex justify-between items-center mb-4 px-3">
          <h3 class="text-2xl font-semibold text-gray-100">Visa Inc</h3>
          <span class="text-2xl font-bold text-gray-100">$234.75</span>
        </div>
        <div class="border-b border-gray-700 mb-4"></div>
        <div class="px-3 mb-4">
          <p class="text-gray-100">Ticker: V</p>
        </div>
        <div class="border-b border-gray-700 mb-4"></div>
        <div class="px-3">
          <p class="text-gray-100">Another description related this recommendation</p>
        </div>
      </div>
    </section>

    <!-- Available Stocks Table -->
    <section class="px-16 py-12">
      <div class="border-t border-gray-700"></div>
      <h2 class="text-2xl font-medium font-orbitron text-gray-400 my-8">Available stocks</h2>
      <div class="bg-gray-800 rounded-lg p-2">
        <table class="w-full">
          <thead>
            <tr class="text-left text-gray-100">
              <th class="p-4 font-normal">NAME</th>
              <th class="p-4 font-normal">TICKER</th>
              <th class="p-4 font-normal">COMPANY</th>
              <th class="p-4 font-normal">PRICE (USD)</th>
              <th class="p-4 font-normal">CHANGE</th>
              <th class="p-4 font-normal">YOU OWN</th>
              <th class="p-4 font-normal">ACTIONS</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-700">
            <tr v-for="(stock, index) in stocks" :key="index" class="text-gray-100 hover:bg-gray-700 transition-colors duration-150">
              <td class="p-4">{{ stock.name }}</td>
              <td class="p-4">
                <span class="text-green-400">{{ stock.ticker }}</span>
              </td>
              <td class="p-4">{{ stock.company }}</td>
              <td class="p-4 font-bold">${{ stock.price }}</td>
              <td class="p-4">
                <component 
                  :is="stock.trend === 'up' ? ArrowUpCircleIcon : stock.trend === 'down' ? ArrowDownCircleIcon : MinusCircleIcon"
                  class="w-6 h-6"
                  :class="{
                    'text-green-500': stock.trend === 'up',
                    'text-red-500': stock.trend === 'down',
                    'text-gray-500': stock.trend === 'neutral'
                  }"
                />
              </td>
              <td class="p-4">{{ stock.owned }}</td>
              <td class="p-4">
                <div class="flex gap-2">
                  <button class="bg-green-500 text-white px-8 py-2 rounded-md font-bold hover:bg-green-600 transition">
                    Buy
                  </button>
                  <button class="border-2 border-white border-opacity-15 text-white px-8 py-2 rounded-md font-bold hover:bg-gray-700 transition">
                    Sell
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- Next Week Button -->
    <div class="px-16 pb-16">
      <button class="bg-green-500 text-white px-8 py-3 rounded-md font-bold hover:bg-green-600 transition">
        Go to next week
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  WalletIcon,
  CircleStackIcon,
  CurrencyDollarIcon,
  ArrowUpCircleIcon,
  ArrowDownCircleIcon,
  MinusCircleIcon
} from '@heroicons/vue/24/outline'

// Mock data for stocks table
const stocks = ref([
  {
    name: 'JPMorgan Chase',
    ticker: 'JPM',
    company: 'JPMorgan Chase & Co.',
    price: '148.25',
    trend: 'up',
    owned: 0
  },
  {
    name: 'JPMorgan Chase',
    ticker: 'JPM',
    company: 'JPMorgan Chase & Co.',
    price: '148.25',
    trend: 'down',
    owned: 0
  },
  {
    name: 'JPMorgan Chase',
    ticker: 'JPM',
    company: 'JPMorgan Chase & Co.',
    price: '148.25',
    trend: 'neutral',
    owned: 0
  }
])
</script>

<style scoped>
.font-orbitron {
  font-family: 'Orbitron', sans-serif;
}
</style> 