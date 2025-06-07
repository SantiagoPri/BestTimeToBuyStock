<template>
  <div class="min-h-screen bg-gray-900 bg-[url('/images/card-home/background-stars.png')] bg-cover">
    <Toast 
      :visible="showToast"
      :message="toastMessage"
      @close="closeToast"
    />
    <!-- Top Bar -->
    <header class="flex justify-between items-center px-16 py-6">
      <div class="flex items-center gap-4">
        <img src="/images/card-home/drakeHead.png" alt="Avatar" class="w-12 h-12 rounded-full" />
        <div class="flex flex-col">
          <h2 class="text-2xl font-bold text-gray-100">Username</h2>
          <span class="text-gray-400">Week {{ currentWeek }}/5</span>
        </div>
      </div>
      <div class="flex gap-6">
        <div class="bg-gray-900 rounded p-2 flex items-center gap-2 w-44">
          <WalletIcon class="w-6 h-6 text-gray-400" />
          <div>
            <div class="text-sm text-gray-100">Balance</div>
            <div class="text-xs text-gray-400">${{ cash.toFixed(2) }} USD</div>
          </div>
        </div>
        <div class="bg-gray-900 rounded p-2 flex items-center gap-2 w-44">
          <CircleStackIcon class="w-6 h-6 text-gray-400" />
          <div>
            <div class="text-sm text-gray-100">Stocks</div>
            <div class="text-xs text-gray-400">${{ holdingsValue.toFixed(2) }} USD</div>
          </div>
        </div>
        <div class="bg-gray-900 rounded p-2 flex items-center gap-2 w-44">
          <CurrencyDollarIcon class="w-6 h-6 text-gray-400" />
          <div>
            <div class="text-sm text-gray-100">Total</div>
            <div class="text-xs text-green-500">${{ totalBalance.toFixed(2) }} USD</div>
          </div>
        </div>
      </div>
    </header>

    <!-- Market News -->
    <section class="px-16 py-12">
      <h2 class="text-2xl font-medium font-orbitron text-gray-400 mb-8">Market News</h2>
      <div class="grid grid-cols-3 gap-8">
        <div v-for="(headline, i) in marketNews" :key="i" class="bg-gray-800 rounded-lg shadow-sm p-3">
          <p class="text-gray-100">{{ headline }}</p>
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
              <th class="p-4 font-normal text-left">COMPANY</th>
              <th class="p-4 font-normal text-left">TICKER</th>
              <th class="p-4 font-normal text-center">RATINGS</th>
              <th class="p-4 font-normal text-center">PRICE (USD)</th>
              <th class="p-4 font-normal text-center">MOOD</th>
              <th class="p-4 font-normal text-center">YOU OWN</th>
              <th class="p-4 font-normal text-center">ACTIONS</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-700">
            <tr v-for="stock in availableStocks" :key="stock.ticker" class="text-gray-100 hover:bg-gray-700 transition-colors duration-150">
              <td class="p-4 text-left">{{ stock.company }}</td>
              <td class="p-4 text-left">
                <span class="text-green-400">{{ stock.ticker }}</span>
              </td>
              <td class="p-4 text-center">{{ stock.ratings }}</td>
              <td class="p-4 text-center">
                <span class="font-bold">${{ stock.currentPrice.toFixed(2) }}</span>
                <span class="ml-4" :class="{
                  'text-green-400': stock.change > 0,
                  'text-red-400': stock.change < 0,
                  'text-gray-400': stock.change === 0
                }">
                  {{ stock.change > 0 ? '+' : '' }}{{ stock.changePercent }}
                </span>
              </td>
              <td class="p-4 text-center">
                <component 
                  :is="stock.marketSentiment === 'up' ? ArrowUpCircleIcon : stock.marketSentiment === 'down' ? ArrowDownCircleIcon : MinusCircleIcon"
                  class="w-8 h-8 mx-auto"
                  :class="{
                    'text-green-400 animate-bounce': stock.marketSentiment === 'up',
                    'text-red-400 animate-bounce': stock.marketSentiment === 'down',
                    'text-gray-400': stock.marketSentiment === 'neutral'
                  }"
                />
              </td>
              <td class="p-4 text-center">{{ holdings[stock.ticker] || 0 }}</td>
              <td class="p-4 text-center">
                <div class="flex gap-2 justify-center">
                  <button 
                    class="bg-green-500 text-white px-8 py-2 rounded-md font-bold hover:bg-green-600 transition"
                    @click="openTradeModal(stock, 'buy')"
                  >
                    Buy
                  </button>
                  <button 
                    class="border-2 border-white border-opacity-15 text-white px-8 py-2 rounded-md font-bold hover:bg-gray-700 transition"
                    @click="openTradeModal(stock, 'sell')"
                  >
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
      <button 
        class="bg-green-500 text-white px-8 py-3 rounded-md font-bold hover:bg-green-600 transition disabled:opacity-50 disabled:cursor-not-allowed"
        @click="handleNextWeek"
        :disabled="isAdvancing"
      >
        <template v-if="isAdvancing">
          {{ currentWeek === 5 ? 'Finishing...' : 'Advancing...' }}
        </template>
        <template v-else>
          {{ currentWeek === 5 ? 'Finish Game' : 'Next Week' }}
        </template>
      </button>
    </div>

    <TradeModal
      :visible="isTradeModalVisible"
      :type="tradeType"
      :stock="selectedStock"
      :cash="cash"
      :holdings="holdings"
      @confirm="handleTradeConfirm"
      @cancel="closeTradeModal"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSessionStore } from '../stores/useSessionStore'
import { GameSessionService } from '../domain/services/GameSessionService'
import { getMarketSentiment, type Stock } from '../domain/entities/Stock'
import { GameSessionApiRepository } from '../infrastructure/repositories/GameSessionApiRepository'
import { HttpClient } from '../infrastructure/http/HttpClient'
import {
  WalletIcon,
  CircleStackIcon,
  CurrencyDollarIcon,
  ArrowUpCircleIcon,
  ArrowDownCircleIcon,
  MinusCircleIcon
} from '@heroicons/vue/24/outline'
import TradeModal from '../components/TradeModal.vue'
import Toast from '../components/Toast.vue'

const router = useRouter()
const sessionStore = useSessionStore()
const gameSessionService = new GameSessionService(new GameSessionApiRepository(new HttpClient()))

// State
const cash = ref(0)
const holdingsValue = ref(0)
const totalBalance = ref(0)
const currentWeek = ref(1)
const marketNews = ref<string[]>([])
const availableStocks = ref<Stock[]>([])
const holdings = ref<Record<string, number>>({})

// Trade modal state
const isTradeModalVisible = ref(false)
const tradeType = ref<'buy' | 'sell'>('buy')
const selectedStock = ref({
  name: '',
  ticker: '',
  price: 0
})
const isTrading = ref(false)
const tradeError = ref<string | null>(null)
const isAdvancing = ref(false)

// Toast state
const showToast = ref(false)
const toastMessage = ref('')

const closeToast = () => {
  showToast.value = false
}

onMounted(async () => {
  try {
    // Check if session exists
    if (!sessionStore.sessionId) {
      router.push('/')
      return
    }

    // Get session state
    const sessionState = await gameSessionService.getSessionState()
    cash.value = sessionState.cash
    holdingsValue.value = sessionState.holdings_value
    totalBalance.value = sessionState.total_balance
    holdings.value = Object.fromEntries(
      Object.entries(sessionState.metadata.holdings).map(([ticker, info]) => [ticker, info.quantity])
    )

    // Extract week number from status
    const weekMatch = sessionState.status.match(/week(\d+)/)
    if (weekMatch) {
      currentWeek.value = parseInt(weekMatch[1])
      
      // Get week data
      const weekData = await gameSessionService.getWeekData(currentWeek.value)
      marketNews.value = weekData.headlines
      
      availableStocks.value = weekData.stocks.map(stock => ({
        ticker: stock.ticker,
        company: stock.companyName,
        currentPrice: stock.price,
        changePercent: `${(stock.priceChange * 100).toFixed(2)}%`,
        change: stock.priceChange,
        ratings: `${stock.rating_from} -> ${stock.rating_to}`,
        marketSentiment: getMarketSentiment(stock.action)
      }))
    } else {
      throw new Error('Invalid week format')
    }
  } catch (error) {
    console.error('Failed to load game state:', error)
    sessionStore.setSessionId('')
    router.push('/')
  }
})

const openTradeModal = (stock: Stock, type: 'buy' | 'sell') => {
  selectedStock.value = {
    name: stock.company,
    ticker: stock.ticker,
    price: stock.currentPrice
  }
  tradeType.value = type
  isTradeModalVisible.value = true
  tradeError.value = null
}

const closeTradeModal = () => {
  isTradeModalVisible.value = false
  tradeError.value = null
}

const handleTradeConfirm = async ({ ticker, quantity, type }: { ticker: string; quantity: number; type: 'buy' | 'sell' }) => {
  if (isTrading.value) return
  
  isTrading.value = true
  tradeError.value = null
  
  try {
    if (type === 'buy') {
      await gameSessionService.buyStocks({
        ticker,
        quantity
      })
      toastMessage.value = "You've successfully bought it!"
    } else {
      await gameSessionService.sellStocks({
        ticker,
        quantity
      })
      toastMessage.value = "You've successfully sold it!"
    }
    
    // Refresh session state after trade
    const sessionState = await gameSessionService.getSessionState()
    cash.value = sessionState.cash
    holdingsValue.value = sessionState.holdings_value
    totalBalance.value = sessionState.total_balance
    holdings.value = Object.fromEntries(
      Object.entries(sessionState.metadata.holdings).map(([ticker, info]) => [ticker, info.quantity])
    )
    
    closeTradeModal()
    showToast.value = true
    setTimeout(() => {
      showToast.value = false
    }, 3000)
  } catch (error) {
    console.error('Trade failed:', error)
    tradeError.value = 'Transaction failed. Please try again.'
  } finally {
    isTrading.value = false
  }
}

const handleNextWeek = async () => {
  if (isAdvancing.value) return

  isAdvancing.value = true
  try {
    if (currentWeek.value === 5) {
      await gameSessionService.endSession()
      router.push('/results')
    } else {
      await gameSessionService.advanceWeek()
      
      // Reload session state and week data
      const sessionState = await gameSessionService.getSessionState()
      cash.value = sessionState.cash
      holdingsValue.value = sessionState.holdings_value
      totalBalance.value = sessionState.total_balance
      holdings.value = Object.fromEntries(
        Object.entries(sessionState.metadata.holdings).map(([ticker, info]) => [ticker, info.quantity])
      )

      // Extract week number and load new week data
      const weekMatch = sessionState.status.match(/week(\d+)/)
      if (weekMatch) {
        currentWeek.value = parseInt(weekMatch[1])
        const weekData = await gameSessionService.getWeekData(currentWeek.value)
        marketNews.value = weekData.headlines
        availableStocks.value = weekData.stocks.map(stock => ({
          ticker: stock.ticker,
          company: stock.companyName,
          currentPrice: stock.price,
          change: stock.priceChange,
          changePercent: `${stock.priceChange * 100 }%`, 
          ratings: `${stock.rating_from} -> ${stock.rating_to}`,
          marketSentiment: getMarketSentiment(stock.action)
        }))
      }
    }
  } catch (error) {
    console.error('Failed to advance week:', error)
    sessionStore.setSessionId('')
    router.push('/')
  } finally {
    isAdvancing.value = false
  }
}
</script>

<style scoped>
.font-orbitron {
  font-family: 'Orbitron', sans-serif;
}
</style> 