<script setup lang="ts">
import type { FunctionalComponent, SVGAttributes } from 'vue';

interface Props {
  icon: FunctionalComponent<SVGAttributes>;
  description: string;
}

defineProps<Props>();
</script>

<template>
  <div class="bg-gray-900 rounded p-4 flex flex-col gap-2">
    <div class="w-6 h-6 text-white">
      <component :is="icon" />
    </div>
    <p class="text-sm text-white whitespace-pre-line">{{ description }}</p>
  </div>
</template> 