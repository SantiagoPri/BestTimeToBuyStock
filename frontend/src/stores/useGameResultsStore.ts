import { defineStore } from 'pinia';
import { ref } from 'vue';

interface GameResults {
  cash: number;
  status: string;
  total_balance: number;
  username: string;
}

export const useGameResultsStore = defineStore('gameResults', () => {
  const results = ref<GameResults | null>(null);

  const setResults = (gameResults: GameResults) => {
    results.value = gameResults;
  };

  const clearResults = () => {
    results.value = null;
  };

  return {
    results,
    setResults,
    clearResults,
  };
});
