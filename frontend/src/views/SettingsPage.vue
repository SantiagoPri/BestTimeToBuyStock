<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useSessionStore } from '../stores/useSessionStore'
import { GameSessionService } from '../domain/services/GameSessionService'
import { GameSessionApiRepository } from '../infrastructure/repositories/GameSessionApiRepository'
import { HttpClient } from '../infrastructure/http/HttpClient'
import {
  BeakerIcon,
  BoltIcon,
  BuildingLibraryIcon,
  BuildingOffice2Icon,
  ChartBarIcon,
  ComputerDesktopIcon,
  CursorArrowRippleIcon,
  GlobeAltIcon,
  HomeModernIcon,
  ShoppingBagIcon,
  SparklesIcon,
  WrenchScrewdriverIcon
} from '@heroicons/vue/24/outline'
import LoadingModal from '../components/LoadingModal.vue'

// Initialize game service
const httpClient = new HttpClient()
const repository = new GameSessionApiRepository(httpClient)
const gameService = new GameSessionService(repository)

interface StockCategory {
  id: number
  name: string
  description: string
  icon: any
  selected: boolean
}

const playerName = ref('')
const selectedCategories = ref<number[]>([])
const showLoading = ref(false)

const router = useRouter()
const sessionStore = useSessionStore()

const stockCategories = ref<StockCategory[]>([
  { id: 1, name: 'Trending', description: 'Most popular stock of the moment', icon: ChartBarIcon, selected: false },
  { id: 2, name: 'Recent', description: 'Latest market additions', icon: SparklesIcon, selected: false },
  { id: 3, name: 'Tech', description: 'Technology and innovation', icon: ComputerDesktopIcon, selected: false },
  { id: 4, name: 'Healthcare', description: 'Health and biotechnology', icon: BeakerIcon, selected: false },
  { id: 5, name: 'Finance', description: 'Banking and financial services', icon: BuildingLibraryIcon, selected: false },
  { id: 6, name: 'Energy', description: 'Energy and renewables', icon: BoltIcon, selected: false },
  { id: 7, name: 'Consumer', description: 'Consumer goods', icon: ShoppingBagIcon, selected: false },
  { id: 8, name: 'Industrials', description: 'Manufacturing and industry', icon: BuildingOffice2Icon, selected: false },
  { id: 9, name: 'Telecom', description: 'Telecommunications and media', icon: GlobeAltIcon, selected: false },
  { id: 10, name: 'Real Estate', description: 'Commercial buildings', icon: HomeModernIcon, selected: false },
  { id: 11, name: 'Utilities', description: 'Public utilities', icon: CursorArrowRippleIcon, selected: false },
  { id: 12, name: 'Materials', description: 'Raw materials', icon: WrenchScrewdriverIcon, selected: false }
])

const toggleCategory = (categoryId: number) => {
  const index = selectedCategories.value.indexOf(categoryId)
  if (index === -1 && selectedCategories.value.length < 3) {
    selectedCategories.value.push(categoryId)
    stockCategories.value.find(cat => cat.id === categoryId)!.selected = true
  } else if (index !== -1) {
    selectedCategories.value.splice(index, 1)
    stockCategories.value.find(cat => cat.id === categoryId)!.selected = false
  }
}

const startPlaying = async () => {
  if (playerName.value.trim() && selectedCategories.value.length === 3) {
    showLoading.value = true
    try {
      const response = await gameService.createSession({
        username: playerName.value,
        categories: selectedCategories.value.map(id => id.toString())
      })
      
      sessionStore.setSessionId(response.sessionId)
      await router.push('/week')
    } catch (error) {
      console.error('Error creating session:', error)
      showLoading.value = false
    }
  }
}


</script>

<template>
  <div class="min-h-screen bg-[#444444] bg-cover bg-center" style="background-image: url('/images/card-home/background-stars.png')">
    <LoadingModal v-if="showLoading" />
    <div class="container mx-auto px-4 py-8 md:py-16">
      <!-- Header -->
      <header class="flex items-center gap-6 mb-16">
        <div class="relative w-24 h-24 md:w-32 md:h-32 bg-[#1C2431] rounded-full overflow-hidden">
          <img src="/images/card-home/drake-gm-1.png" alt="Drake GM" class="absolute bottom-4 w-full h-auto">
        </div>
        <h1 class="text-4xl md:text-6xl font-orbitron text-white">Best time to buy stock</h1>
      </header>

      <div class="flex flex-col gap-8">
        <!-- Player Name Section -->
        <div class="bg-[#1F2937] border border-[#374151] rounded-lg p-6 shadow-sm">
          <div class="flex items-center gap-4 mb-6">
            <div class="w-12 h-12 rounded-lg overflow-hidden">
              <img src="/images/card-home/drakeHead.png" alt="Player Avatar" class="w-full h-full object-cover">
            </div>
            <h2 class="text-xl font-bold text-white">Choose your player name</h2>
          </div>
          <div class="bg-[#374151] rounded-lg p-3">
            <input
              v-model="playerName"
              type="text"
              placeholder="Enter your name..."
              class="w-full bg-transparent text-white placeholder-gray-400 focus:outline-none"
            >
          </div>
        </div>

        <!-- Stock Categories Section -->
        <div class="bg-[#1F2937] border border-[#374151] rounded-lg p-6 shadow-sm">
          <div class="flex items-center gap-4 mb-6">
            <div class="w-12 h-12 rounded-lg overflow-hidden">
              <img src="/images/card-home/bookmark-2.png" alt="Categories" class="w-full h-full object-cover">
            </div>
            <h2 class="text-xl font-bold text-white">Select 3 stock categories</h2>
          </div>
          <p class="text-gray-400 mb-6">We'll use your picks to show you the most relevant news and stock tips during the game.</p>
          
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            <button
              v-for="category in stockCategories"
              :key="category.id"
              @click="toggleCategory(category.id)"
              :class="[
                'flex items-center gap-2 p-4 rounded-lg transition-colors border',
                category.selected
                  ? 'bg-[#21393B] border-[#34C759] text-white'
                  : 'bg-[#374151] border-[#374151] text-white hover:bg-[#4B5563]'
              ]"
            >
              <component
                :is="category.icon"
                class="w-8 h-8"
                :class="category.selected ? 'text-[#34C759]' : 'text-[#34C759]'"
              />
              <div class="text-left flex-1">
                <div class="flex items-center justify-between">
                  <h3 class="font-bold text-[15.25px]">{{ category.name }}</h3>
                  <div v-if="category.selected" class="w-6 h-6 flex items-center justify-center">
                    <svg width="16" height="11" viewBox="0 0 16 11" fill="none" xmlns="http://www.w3.org/2000/svg">
                      <path d="M1 5L6 10L15 1" stroke="#34C759" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                  </div>
                </div>
                <p class="text-sm text-gray-400">{{ category.description }}</p>
              </div>
            </button>
          </div>
        </div>
      </div>

      <!-- Start Button -->
      <div class="mt-8 flex flex-col items-center gap-2">
        <p class="text-gray-400 text-sm">
          {{ selectedCategories.length }}/3 categories selected
        </p>
        <button
          @click="startPlaying"
          :disabled="!playerName.trim() || selectedCategories.length !== 3"
          :class="[
            'h-10 px-8 rounded-md font-normal text-[17.296875px] leading-[1.62] transition-colors',
            selectedCategories.length === 3 && playerName.trim()
              ? 'text-[#E1E1E1] bg-[#22C55E] hover:bg-[#1EA34D]'
              : 'text-[#E1E1E1] bg-gray-500 cursor-not-allowed'
          ]"
        >
          Start playing
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.font-orbitron {
  font-family: 'Orbitron', sans-serif;
}
</style> 