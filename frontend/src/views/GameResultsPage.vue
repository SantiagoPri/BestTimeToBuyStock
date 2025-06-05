<script setup lang="ts">
import { computed } from 'vue';
import { WalletIcon, PresentationChartLineIcon, CurrencyDollarIcon, ChartBarIcon, NewspaperIcon, ScaleIcon, ClockIcon } from '@heroicons/vue/24/outline';
import SummaryStat from '../components/SummaryStat.vue';
import ResultCard from '../components/ResultCard.vue';

interface Props {
  status: 'win' | 'lose' | 'neutral';
  initialBalance: number;
  finalBalance: number;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: 'navigate', path: string): void;
}>();

const difference = computed(() => props.finalBalance - props.initialBalance);

interface StatusConfig {
  bgColor: string;
  borderColor: string;
  message: string;
  image: string;
}

const statusConfig = computed<StatusConfig>(() => {
  switch (props.status) {
    case 'win':
      return {
        bgColor: 'bg-emerald-700/10',
        borderColor: 'border-emerald-500',
        message: 'Well done, trader! Even the dragon is impressed with those profits.',
        image: '/images/card-home/drake-winner.png'
      };
    case 'neutral':
      return {
        bgColor: 'bg-yellow-800/10',
        borderColor: 'border-yellow-700',
        message: 'Not bad! You managed to stay in the game. Keep learning and improving!',
        image: '/images/card-home/drake-neutro.png'
      };
    case 'lose':
      return {
        bgColor: 'bg-red-900/10',
        borderColor: 'border-red-900',
        message: 'The market can be tough, but every loss is a lesson. Try again!',
        image: '/images/card-home/drake-loser.png'
      };
    default:
      return {
        bgColor: 'bg-gray-800/10',
        borderColor: 'border-gray-700',
        message: 'Game completed!',
        image: '/images/card-home/drake-neutro.png'
      };
  }
});

const differenceColor = computed(() => {
  if (difference.value > 0) return 'text-emerald-500';
  if (difference.value < 0) return 'text-red-500';
  return 'text-yellow-600';
});

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(value);
};

const tips = [
  {
    icon: NewspaperIcon,
    description: 'Study the news\nMarket news gives you clues about which stocks might go up or down.'
  },
  {
    icon: ChartBarIcon,
    description: 'Diversify\nDon\'t put all your eggs in one basket. Invest in different sectors.'
  },
  {
    icon: ClockIcon,
    description: 'Timing is key\nWatch weekly trends to decide better. Your timing will define your success!'
  },
  {
    icon: ScaleIcon,
    description: 'Take calculated risks\nDon\'t be too conservative, but don\'t bet everything on one card either.'
  }
];
</script>

<template>
  <div class="min-h-screen bg-[url('/images/card-home/background-stars.png')] bg-cover bg-center bg-gray-900">
    <!-- Header with background -->
    <div class="w-full h-[456px] bg-[url('/images/card-home/header-background.png')] bg-cover bg-center bg-black/60">
      <!-- Logo -->
      <div class="container mx-auto px-16 pt-6">
        <div class="w-[120px] h-[120px] rounded-[187.5px] bg-[#1C2431] flex items-center justify-center">
          <img 
            src="/images/card-home/drake-gm-1.png" 
            alt="Dragon Game Master Logo" 
            class="w-[119px] h-[93px] object-contain mt-[15px]"
          />
        </div>
      </div>
      
      <div class="container mx-auto px-16 pt-12">
        <!-- Header Text -->
        <div class="text-center mb-12">
          <h1 class="font-orbitron text-3xl text-white mb-4">The Game Master Has Spoken!</h1>
          <p class="text-gray-400 text-lg">Let's see how you did...</p>
        </div>
      </div>
    </div>

    <!-- Content Section -->
    <div class="container mx-auto px-16 pt-16">
      <!-- Result Card -->
      <div 
        class="max-w-4xl w-full p-12 rounded-lg border mb-16 mx-auto"
        :class="[statusConfig.bgColor, statusConfig.borderColor]"
      >
        <p class="text-2xl text-white text-center font-bold mb-6">{{ statusConfig.message }}</p>
        <img 
          :src="statusConfig.image" 
          :alt="'Dragon ' + status" 
          class="w-48 h-48 mx-auto mb-6"
        />
        <p class="text-base text-white text-center">
          Keep making smart moves. The Top 10 traders are waiting for you. Show the market what you're made of!
        </p>
      </div>

      <!-- Summary Section -->
      <h2 class="font-orbitron text-2xl text-white mb-8 text-center">Your investment summary</h2>
      <div class="flex gap-6 flex-wrap justify-center mb-16">
        <SummaryStat
          :icon="WalletIcon"
          label="Initial Balance"
          :value="formatCurrency(initialBalance)"
        />
        <SummaryStat
          :icon="PresentationChartLineIcon"
          label="Final Balance"
          :value="formatCurrency(finalBalance)"
        />
        <SummaryStat
          :icon="CurrencyDollarIcon"
          :label="difference >= 0 ? 'Gains' : 'Loss'"
          :value="formatCurrency(Math.abs(difference))"
          :value-color="differenceColor"
        />
      </div>

      <!-- Tips Section -->
      <h2 class="font-orbitron text-2xl text-white mb-8 text-center">Tips for your next game</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 max-w-4xl mx-auto mb-16">
        <ResultCard
          v-for="tip in tips"
          :key="tip.description"
          :icon="tip.icon"
          :description="tip.description"
        />
      </div>

      <!-- Coin GIF and Action Button -->
      <div class="flex flex-col items-center gap-6 py-12">
        <img 
          src="/images/card-home/coin-gift.gif" 
          alt="Coin animation" 
          class="w-24 h-24 mb-4"
        />
        <button
          class="px-8 py-3 border-2 border-white/15 rounded-md text-white hover:bg-white/5 transition-colors"
          @click="emit('navigate', '/settings')"
        >
          Back home
        </button>
      </div>
    </div>
  </div>
</template> 