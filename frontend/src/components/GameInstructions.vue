<script setup lang="ts">
interface Instruction {
  step: number;
  title: string;
  description: string;
}

defineProps<{
  instructions?: Instruction[];
}>();

const defaultInstructions: Instruction[] = [
  {
    step: 1,
    title: "Start playing with $10,000 in virtual money",
    description: "It's all simulated, no real money, no risk. Just learn and play!"
  },
  {
    step: 2,
    title: "You have 5 weeks to play, 5 rounds",
    description: "Each week is a new game level. make your moves, react to the market, and become a top trader."
  },
  {
    step: 3,
    title: "Invest your full balance",
    description: "Don't hold back, take calculated risks with the Game Master recommendations and make every dollar work for you!"
  },
  {
    step: 4,
    title: "Big Gains or Big Lessons?",
    description: "Made the right calls? Let's see how much you racked up from your virtual winnings."
  }
];
</script>

<template>
  <div class="bg-transparent rounded-lg border border-[#374151] p-6 shadow-[0px_1px_2px_0px_rgba(0,0,0,0.05)]">
    <h2 class="text-[32px] font-orbitron font-bold text-white mb-12 leading-[1.43] tracking-[-2.69%]">Your Trading Journey Begins:</h2>
    
    <div class="space-y-16">
      <div 
        v-for="instruction in (instructions || defaultInstructions)" 
        :key="instruction.step"
        class="flex gap-6"
      >
        <!-- Step Number -->
        <div class="relative flex-shrink-0 w-12 h-12">
          <img src="/images/card-home/group-203.svg" alt="Number background" class="absolute inset-0 w-full h-full" />
          <span class="absolute inset-0 flex items-center justify-center text-[16px] font-bold text-[#22C55E] leading-[1.5]">
            {{ instruction.step }}
          </span>
        </div>
        
        <!-- Instruction Content -->
        <div class="space-y-1">
          <h3 class="font-bold text-[15.875px] leading-[1.51] text-white">{{ instruction.title }}</h3>
          <p class="text-[#9CA3AF] text-[15.875px] leading-[1.51]">{{ instruction.description }}</p>
        </div>
      </div>
    </div>

    <!-- Call to Action -->
    <div class="mt-12 text-center space-y-6">
      <h3 class="text-[17.578125px] font-bold text-[#E1E1E1] leading-[1.59]">Are you ready?</h3>
      <img src="/images/card-home/coin-gift.gif" alt="Coin Gift" class="w-[100px] h-[100px] mx-auto" />
      <p class="text-[15.75px] text-[#D1D5DB] leading-[1.52]">Master the markets and climb the leaderboard!</p>
      <button class="w-fit mx-auto bg-[#22C55E] text-[#E1E1E1] py-3 px-8 rounded-[6px] hover:bg-emerald-600 transition-colors text-[17.296875px] leading-[1.62] font-normal">
        Start
      </button>
    </div>
  </div>
</template> 