import { createRouter, createWebHistory } from 'vue-router';
import SettingsPage from '../../views/SettingsPage.vue';
import Home from '../../components/Home.vue';
import WeekGamePage from '../../views/WeekGamePage.vue';
import GameResultsPage from '../../views/GameResultsPage.vue';

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
    path: '/week',
    name: 'week',
    component: WeekGamePage,
  },
  {
    path: '/results',
    name: 'results',
    component: GameResultsPage,
    props: true, // This allows passing route params as props
  },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 };
  },
});
