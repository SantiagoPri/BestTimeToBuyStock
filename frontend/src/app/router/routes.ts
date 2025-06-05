import { createRouter, createWebHistory } from 'vue-router';
import StockListView from '../../modules/stocks/presentation/views/StockListView.vue';
import SettingsPage from '../../views/SettingsPage.vue';
import Home from '../../components/Home.vue';

const routes = [
  {
    path: '/',
    redirect: '/home',
  },
  {
    path: '/home',
    name: 'home',
    component: Home,
  },
  {
    path: '/settings',
    name: 'settings',
    component: SettingsPage,
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
