import { computed, ref } from 'vue';

export interface PaginationOptions {
  initialPage?: number;
  initialLimit?: number;
}

export function usePagination(options: PaginationOptions = {}) {
  const currentPage = ref(options.initialPage || 1);
  const itemsPerPage = ref(options.initialLimit || 10);

  const offset = computed(() => {
    return (currentPage.value - 1) * itemsPerPage.value;
  });

  function goToPage(page: number) {
    currentPage.value = page;
  }

  function nextPage() {
    currentPage.value++;
  }

  function previousPage() {
    if (currentPage.value > 1) {
      currentPage.value--;
    }
  }

  function setItemsPerPage(limit: number) {
    itemsPerPage.value = limit;
    currentPage.value = 1; // Reset to first page when changing items per page
  }

  return {
    currentPage,
    itemsPerPage,
    offset,
    goToPage,
    nextPage,
    previousPage,
    setItemsPerPage,
  };
}
