<script setup lang="ts">
import type { FunctionalComponent, SVGAttributes } from 'vue';

interface Props {
  icon: FunctionalComponent<SVGAttributes>;
  label: string;
  value: string;
  valueColor?: string;
}

defineProps<Props>();
</script>

<template>
  <div class="flex items-center gap-2 bg-gray-800 rounded p-2 w-[178px]">
    <div class="w-6 h-6 text-white">
      <component :is="icon" />
    </div>
    <div class="flex flex-col">
      <span class="text-sm text-white">{{ label }}</span>
      <span class="text-xs" :class="valueColor || 'text-gray-400'">{{ value }}</span>
    </div>
  </div>
</template> 