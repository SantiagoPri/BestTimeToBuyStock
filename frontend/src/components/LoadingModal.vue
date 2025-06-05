<script setup lang="ts">
// No props or logic needed for this simple modal
</script>

<template>
  <Teleport to="body">
    <div class="fixed inset-0 bg-[#1C1C1C]/90 flex items-center justify-center z-50">
      <div class="bg-black rounded-xl p-6 shadow-lg max-w-md w-full mx-4">
        <p class="text-white text-[24px] leading-tight text-center mb-6 font-orbitron">
          The Game Master is preparing<br>
          your game board...<br>
          the dragon's on fire!
        </p>
        <div class="flex justify-center mb-6">
          <img 
            src="/images/card-home/loadingModal.gif" 
            alt="Loading..."
            class="w-[180px] h-[101px] object-contain"
          />
        </div>
        <div class="flex justify-center">
          <img 
            src="/images/card-home/dancingDrake.gif" 
            alt="Dancing Drake"
            class="w-32 h-32 object-contain"
          />
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.font-orbitron {
  font-family: 'Orbitron', sans-serif;
}
</style> 