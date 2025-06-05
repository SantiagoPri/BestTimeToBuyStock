<script setup lang="ts">
import { ref, computed } from 'vue'

interface Stock {
  name: string
  ticker: string
  price: number
}

interface Props {
  visible: boolean
  type: 'buy' | 'sell'
  stock: Stock
  cash: number
  holdings: Record<string, number>
}

const props = defineProps<Props>()
const emit = defineEmits<{
  confirm: [{ ticker: string; quantity: number; type: 'buy' | 'sell' }]
  cancel: []
}>()

const shares = ref<number>(0)
const isLoading = ref(false)
const error = ref<string | null>(null)

const total = computed(() => {
  return shares.value * props.stock.price
})

const hasError = computed(() => {
  if (!shares.value || shares.value <= 0) {
    return 'Please enter a valid number of shares'
  }
  
  if (props.type === 'buy' && total.value > props.cash) {
    return 'Insufficient funds for this purchase'
  }
  
  if (props.type === 'sell') {
    const currentHoldings = props.holdings[props.stock.ticker] || 0
    if (shares.value > currentHoldings) {
      return `You only own ${currentHoldings} shares`
    }
  }
  
  return null
})

const handleConfirm = async () => {
  if (hasError.value || isLoading.value) return
  
  isLoading.value = true
  error.value = null
  
  try {
    emit('confirm', {
      ticker: props.stock.ticker,
      quantity: shares.value,
      type: props.type
    })
  } catch (e) {
    error.value = 'Transaction failed. Please try again.'
  } finally {
    isLoading.value = false
  }
}

const handleCancel = () => {
  shares.value = 0
  error.value = null
  emit('cancel')
}

const formatUSD = (value: number) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(value)
}

const handleInput = (event: Event) => {
  const input = event.target as HTMLInputElement
  input.value = input.value.replace(/\D/g, '')
  error.value = null
}
</script>

<template>
  <Teleport to="body">
    <div
      v-if="visible"
      class="fixed inset-0 bg-[#1C1C1C]/90 flex items-center justify-center z-50"
      @click="handleCancel"
    >
      <div
        class="bg-[#111827] rounded-xl shadow-lg w-[400px] p-6 text-white"
        @click.stop
      >
        <h2 class="text-[20px] font-bold text-white">
          {{ type === 'buy' ? 'Buy' : 'Sell' }} {{ stock.name }}
        </h2>
        <div class="h-1 w-12 bg-green-400 mt-1 mb-4"></div>

        <div class="space-y-4">
          <h3 class="text-lg font-semibold text-gray-300">Your Trade Recap</h3>

          <div class="bg-[#1F2937] rounded-lg p-4 space-y-3">
            <div class="flex justify-between">
              <span class="text-gray-400">Ticker:</span>
              <span class="text-gray-200">{{ stock.ticker }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-400">Current price:</span>
              <span class="text-gray-200">{{ formatUSD(stock.price) }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-400">Your balance:</span>
              <span class="text-gray-200">{{ formatUSD(cash) }}</span>
            </div>
            <div v-if="type === 'sell'" class="flex justify-between">
              <span class="text-gray-400">You own:</span>
              <span class="text-gray-200">{{ holdings[stock.ticker] || 0 }} shares</span>
            </div>
          </div>

          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-400">
              Number of shares to {{ type === 'buy' ? 'buy' : 'sell' }}
            </label>
            <input
              type="number"
              v-model.number="shares"
              min="0"
              step="1"
              class="bg-[#1F2937] text-white border border-gray-500 rounded-md p-2 w-full"
              :class="{ 'border-red-500': hasError }"
              @input="handleInput"
              :disabled="isLoading"
            />
            <p class="text-sm text-gray-400">No decimals allowed</p>
          </div>

          <div class="bg-[#1F2937] rounded-lg p-4 flex justify-between font-semibold">
            <span class="text-gray-400">Total:</span>
            <span class="text-gray-200">{{ formatUSD(total) }}</span>
          </div>

          <div v-if="hasError || error" class="bg-red-900/20 border border-red-800 rounded-lg p-3">
            <p class="text-red-400 text-sm">{{ hasError || error }}</p>
          </div>

          <div class="flex gap-3 mt-6">
            <button
              class="flex-1 bg-green-500 hover:bg-green-600 text-white font-semibold py-2 px-4 rounded-md disabled:opacity-50 disabled:cursor-not-allowed"
              @click="handleConfirm"
              :disabled="!!hasError || isLoading"
            >
              <template v-if="isLoading">
                <span class="inline-block animate-spin mr-2">⟳</span>
                Processing...
              </template>
              <template v-else>
                Confirm {{ type === 'buy' ? 'Purchase' : 'Sale' }}
              </template>
            </button>
            <button
              class="flex-1 border border-gray-400 text-white font-medium py-2 px-4 rounded-md hover:bg-[#1F2937] disabled:opacity-50 disabled:cursor-not-allowed"
              @click="handleCancel"
              :disabled="isLoading"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template> 