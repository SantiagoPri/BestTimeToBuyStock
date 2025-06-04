<script setup lang="ts">
interface LeaderboardPlayer {
  rank: number;
  username: string;
  profit: number;
  profitPercentage: number;
}

interface Props {
  players?: LeaderboardPlayer[];
  title?: string;
}

defineProps<Props>();

const defaultPlayers: LeaderboardPlayer[] = [
  { rank: 1, username: 'DragonTrader', profit: 8750, profitPercentage: 87.5 },
  { rank: 2, username: 'StockMaster99', profit: 7240, profitPercentage: 72.4 },
  { rank: 3, username: 'WolfOfWallSt', profit: 6890, profitPercentage: 68.9 },
  { rank: 4, username: 'BullRunner', profit: 5630, profitPercentage: 56.3 },
  { rank: 5, username: 'MarketNinja', profit: 4920, profitPercentage: 49.2 },
  { rank: 6, username: 'GoldDigger', profit: 4150, profitPercentage: 41.5 },
  { rank: 7, username: 'CryptoKing', profit: 3780, profitPercentage: 37.8 },
  { rank: 8, username: 'TechInvestor', profit: 3240, profitPercentage: 32.4 },
  { rank: 9, username: 'SmartMoney', profit: 2890, profitPercentage: 28.9 },
  { rank: 10, username: 'RiskyBiz', profit: 2150, profitPercentage: 21.5 }
];

const getRankStyles = (rank: number) => {
  switch (rank) {
    case 1:
      return {
        background: '#E3AC00',
        inner: '#F9C900'
      };
    case 2:
      return {
        background: '#74787C',
        inner: '#9CA3AF'
      };
    case 3:
      return {
        background: '#AB7521',
        inner: '#CA8A04'
      };
    default:
      return {
        background: '#4B5563',
        inner: '#4B5563'
      };
  }
};
</script>

<template>
  <div class="bg-[#1F2937] rounded-lg border border-[#374151] shadow-[0px_1px_2px_0px_rgba(0,0,0,0.05)] w-[328px]">
    <!-- Header -->
    <div class="flex items-center gap-4 p-6 mb-[-12px]">
      <div class="w-12 h-12 flex items-center justify-center">
        <img src="/images/card-home/trophy-cup.png" alt="Trophy" class="w-11 h-12" />
      </div>
      <h2 class="font-orbitron font-bold text-[22.125px] text-[#E1E1E1] leading-[1.45] tracking-[-2.71%]">
        {{ title || 'Top 10 Traders' }}
      </h2>
    </div>

    <!-- Players List -->
    <div class="divide-y divide-[#374151]">
      <div
        v-for="player in (players || defaultPlayers)"
        :key="player.username"
        :class="[
          'flex items-center p-4 transition-colors',
          player.rank <= 3 ? 'bg-[rgba(27,31,38,0.72)]' : 'bg-[#374151] hover:bg-[#374151]'
        ]"
      >
        <!-- Rank Circle -->
        <div class="relative w-12 h-12">
          <div 
            v-if="player.rank <= 3"
            class="absolute inset-0 flex items-center"
          >
            <img 
              :src="`/images/card-home/star-${player.rank}.svg`" 
              :alt="`Rank ${player.rank} star`" 
              class="w-12 h-12"
            />
          </div>
          <div 
            v-else
            class="w-8 h-8 rounded-full bg-[#4B5563] flex items-center justify-center mx-auto my-2"
          >
            <span class="font-bold text-[#E1E1E1] text-[16px] leading-[1.5]">{{ player.rank }}</span>
          </div>
        </div>

        <!-- Player Info -->
        <div class="flex-1 ml-2">
          <span :class="[
            'font-bold',
            player.rank === 1 ? 'text-[14.875px]' :
            player.rank === 2 ? 'text-[15.125px]' :
            player.rank === 3 ? 'text-[15.25px]' :
            'text-[15px]',
            'text-[#E1E1E1]'
          ]">{{ player.username }}</span>
        </div>

        <!-- Profit Info -->
        <div class="text-right">
          <div class="text-[#4ADE80] font-bold text-[15.625px] leading-[1.536]">
            +${{ player.profit.toLocaleString() }}
          </div>
          <div class="text-[#9CA3AF] text-[14px] leading-[1.43]">
            +{{ player.profitPercentage.toFixed(1) }}%
          </div>
        </div>
      </div>
    </div>
  </div>
</template> 