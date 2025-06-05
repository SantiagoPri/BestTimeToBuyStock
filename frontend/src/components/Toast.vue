<template>
  <Transition
    enter-active-class="transition ease-out duration-300"
    enter-from-class="transform opacity-0 scale-95"
    enter-to-class="transform opacity-100 scale-100"
    leave-active-class="transition ease-in duration-200"
    leave-from-class="transform opacity-100 scale-100"
    leave-to-class="transform opacity-0 scale-95"
  >
    <div 
      v-if="visible"
      class="fixed top-4 right-4 w-[280px] z-50"
    >
      <div class="bg-[#E8FAF0] rounded-2xl p-3 flex items-center gap-4 relative">
        <!-- Success Icon -->
        <div class="w-6 h-6 flex-shrink-0">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2C6.49 2 2 6.49 2 12C2 17.51 6.49 22 12 22C17.51 22 22 17.51 22 12C22 6.49 17.51 2 12 2ZM16.78 9.7L11.11 15.37C10.97 15.51 10.78 15.59 10.58 15.59C10.38 15.59 10.19 15.51 10.05 15.37L7.22 12.54C6.93 12.25 6.93 11.77 7.22 11.48C7.51 11.19 7.99 11.19 8.28 11.48L10.58 13.78L15.72 8.64C16.01 8.35 16.49 8.35 16.78 8.64C17.07 8.93 17.07 9.4 16.78 9.7Z" fill="#17C964"/>
          </svg>
        </div>
        
        <!-- Message -->
        <span class="text-sm font-normal text-gray-900 font-noto">{{ message }}</span>
        
        <!-- Close Button -->
        <button 
          @click="$emit('close')"
          class="absolute -top-1 -right-1 bg-[#F4F4F5] border border-[#D4D4D8] rounded-full p-0.5"
        >
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M6 1.5C3.525 1.5 1.5 3.525 1.5 6C1.5 8.475 3.525 10.5 6 10.5C8.475 10.5 10.5 8.475 10.5 6C10.5 3.525 8.475 1.5 6 1.5ZM7.6875 7.6875C7.5 7.875 7.1875 7.875 7 7.6875L6 6.6875L5 7.6875C4.8125 7.875 4.5 7.875 4.3125 7.6875C4.125 7.5 4.125 7.1875 4.3125 7L5.3125 6L4.3125 5C4.125 4.8125 4.125 4.5 4.3125 4.3125C4.5 4.125 4.8125 4.125 5 4.3125L6 5.3125L7 4.3125C7.1875 4.125 7.5 4.125 7.6875 4.3125C7.875 4.5 7.875 4.8125 7.6875 5L6.6875 6L7.6875 7C7.875 7.1875 7.875 7.5 7.6875 7.6875Z" fill="#1F2937"/>
          </svg>
        </button>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
defineProps<{
  visible: boolean
  message: string
}>()

defineEmits<{
  (e: 'close'): void
}>()
</script>

<style scoped>
.font-noto {
  font-family: 'Noto Sans', sans-serif;
}
</style> 