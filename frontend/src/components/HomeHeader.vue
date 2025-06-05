<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';

defineProps<{
  title?: string;
}>();

const router = useRouter();
const isLoggedIn = ref(false); // This would come from your auth store in a real app

const goToSettings = () => {
  router.push('/settings');
};
</script>

<template>
  <header
    class="relative min-h-[654px] bg-cover bg-center px-16 py-6"
    style="background-image: url('/images/card-home/header-background.png')"
  >
    <div class="flex justify-between items-start w-full">
      <!-- Logo -->
      <img src="/images/card-home/drake-gm-1.png" alt="Logo" class="w-[120px] h-[120px] rounded-full bg-[#1C2431]" />
    </div>

    <!-- Main centered block -->
    <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-full max-w-7xl px-4">
      <!-- Title -->
      <h1 class="font-orbitron font-normal text-[64px] text-white leading-[70.4px] tracking-[-1.28px] mb-6 whitespace-nowrap">
        {{ title || 'Best time to buy stock' }}
      </h1>

      <!-- Call to action -->
      <div v-if="!isLoggedIn" class="flex flex-col items-start gap-2">
        <p class="font-bold text-[14.875px] text-[#22C55E]">Already know the game?</p>
        <button
          @click="goToSettings"
          class="px-8 py-3 h-[40px] border-2 border-[rgba(255,255,255,0.15)] rounded-[6px] text-[#E1E1E1] hover:bg-gray-700/20 transition-colors text-[17.3px]"
        >
          Play now!
        </button>
      </div>
    </div>
  </header>
</template>
