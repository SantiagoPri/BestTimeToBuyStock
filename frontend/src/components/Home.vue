<script setup lang="ts">
import HomeHeader from './HomeHeader.vue';
import GameInstructions from './GameInstructions.vue';
import LeaderboardCard from './LeaderboardCard.vue';
</script>

<template>
  <div class="min-h-screen bg-[#E1E1E1]">
    <!-- Header Section -->
    <HomeHeader />
    
    <!-- Main Content -->
    <main class="px-16 py-6 bg-[url('/images/card-home/background-stars.png')] bg-cover bg-center bg-black/20">
      <div class="flex justify-center gap-4 max-w-[1280px] mx-auto">
        <!-- Game Instructions -->
        <GameInstructions class="flex-1 max-w-[814px]" />
        
        <!-- Leaderboard -->
        <LeaderboardCard />
      </div>
    </main>
  </div>
</template>

<style>
@font-face {
  font-family: 'Orbitron';
  src: url('/fonts/Orbitron-Regular.woff2') format('woff2');
  font-weight: 400;
  font-style: normal;
  font-display: swap;
}

@font-face {
  font-family: 'Orbitron';
  src: url('/fonts/Orbitron-Bold.woff2') format('woff2');
  font-weight: 700;
  font-style: normal;
  font-display: swap;
}
</style> 