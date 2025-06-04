<script setup lang="ts">
import { ref, onMounted } from 'vue';
import type { GameSession } from '../domain/entities/GameSession';
import { GameSessionService } from '../domain/services/GameSessionService';
import { GameSessionApiRepository } from '../infrastructure/repositories/GameSessionApiRepository';
import { HttpClient } from '../infrastructure/http/HttpClient';

const httpClient = new HttpClient();
const repository = new GameSessionApiRepository(httpClient);
const gameService = new GameSessionService(repository);

const players = ref<GameSession[]>([]);

onMounted(async () => {
  try {
    players.value = await gameService.getLeaderboard();
  } catch (error) {
    console.error('Failed to fetch leaderboard:', error);
  }
});

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
        Top 10 Traders
      </h2>
    </div>

    <!-- Players List -->
    <div class="divide-y divide-[#374151]">
      <div
        v-for="(player, index) in players"
        :key="player.session_id"
        :class="[
          'flex items-center p-4 transition-colors',
          index < 3 ? 'bg-[rgba(27,31,38,0.72)]' : 'bg-[#374151] hover:bg-[#374151]'
        ]"
      >
        <!-- Rank Circle -->
        <div class="relative w-12 h-12">
          <div 
            v-if="index < 3"
            class="absolute inset-0 flex items-center"
          >
            <img 
              :src="`/images/card-home/star-${index + 1}.svg`" 
              :alt="`Rank ${index + 1} star`" 
              class="w-12 h-12"
            />
          </div>
          <div 
            v-else
            class="w-8 h-8 rounded-full bg-[#4B5563] flex items-center justify-center mx-auto my-2"
          >
            <span class="font-bold text-[#E1E1E1] text-[16px] leading-[1.5]">{{ index + 1 }}</span>
          </div>
        </div>

        <!-- Player Info -->
        <div class="flex-1 ml-2">
          <span :class="[
            'font-bold',
            index === 0 ? 'text-[14.875px]' :
            index === 1 ? 'text-[15.125px]' :
            index === 2 ? 'text-[15.25px]' :
            'text-[15px]',
            'text-[#E1E1E1]'
          ]">{{ player.username }}</span>
        </div>

        <!-- Profit Info -->
        <div class="text-right">
          <div class="text-[#4ADE80] font-bold text-[15.625px] leading-[1.536]">
            +${{ player.total_balance.toLocaleString() }}
          </div>
          <div class="text-[#9CA3AF] text-[14px] leading-[1.43]">
            +{{ ((player.total_balance - 10000) / 100).toFixed(1) }}%
          </div>
        </div>
      </div>
    </div>
  </div>
</template> 