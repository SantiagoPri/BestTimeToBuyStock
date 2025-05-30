import { createRouter, createWebHistory } from 'vue-router';
import StockListView from '../../modules/stocks/presentation/StockListView.vue';

const routes = [
  {
    path: '/',
    redirect: '/stocks',
  },
  {
    path: '/stocks',
    name: 'stocks',
    component: StockListView,
  },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});
