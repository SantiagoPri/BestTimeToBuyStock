import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useSessionStore = defineStore('session', () => {
  const sessionId = ref<string | null>(null);

  const setSessionId = (id: string) => {
    sessionId.value = id;
    sessionStorage.setItem('sessionId', id);
  };

  // Initialize from sessionStorage if available
  if (typeof window !== 'undefined') {
    const stored = sessionStorage.getItem('sessionId');
    if (stored) {
      sessionId.value = stored;
    }
  }

  return {
    sessionId,
    setSessionId,
  };
});
