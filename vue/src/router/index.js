import { createRouter, createWebHistory } from 'vue-router';

// Import components from each folder
import HomePage from '../components/HomePage.vue';

import WebScan from '../components/01-web-protection/WebScan.vue';
import TrafficGeneration from '../components/01-web-protection/TrafficGeneration.vue';
import WebAttacks from '../components/01-web-protection/WebAttacks.vue';
import MachineLearning from '../components/01-web-protection/MachineLearning.vue';
import CookieSecurity from '../components/01-web-protection/CookieSecurity.vue';
import CredentialStuffingDefense from '../components/01-web-protection/CredentialStuffingDefense.vue';

import BotSelenium from '../components/02-bot-mitigation/BotSelenium.vue';
import BotDeception from '../components/02-bot-mitigation/BotDeception.vue';
import KnownBots from '../components/02-bot-mitigation/KnownBots.vue';
import BotScraping from '../components/02-bot-mitigation/BotScraping.vue';

import ApiRequests from '../components/03-api-protection/ApiRequests.vue';
import ApiTrafficGeneration from '../components/03-api-protection/ApiTrafficGeneration.vue';

import RestAPIManagement from '../components/04-rest-api/RestAPIManagement.vue';

import ConfigForm from '../components/05-tool/ConfigForm.vue';
import HealthCheck from '../components/05-tool/HealthCheck.vue';
import AppDocker from '../components/05-tool/AppDocker.vue';
import FortiWebBootstrap from '../components/05-tool/FortiWebBootstrap.vue';

const routes = [
  { path: '/', component: HomePage },

  { path: '/web-scan', component: WebScan },
  { path: '/traffic-generation', component: TrafficGeneration },
  { path: '/web-attacks', component: WebAttacks },
  { path: '/machine-learning', component: MachineLearning },
  { path: '/cookie-security', component: CookieSecurity },
  { path: '/credential-stuffing-defense', component: CredentialStuffingDefense },

  { path: '/bot-deception', component: BotDeception },
  { path: '/bot-selenium', component: BotSelenium },
  { path: '/known-bots', component: KnownBots },
  { path: '/bot-scraping', component: BotScraping },

  { path: '/api-requests', component: ApiRequests },
  { path: '/api-traffic-generation', component: ApiTrafficGeneration },

  { path: '/rest-api-management', component: RestAPIManagement },

  { path: '/configuration', component: ConfigForm },
  { path: '/health-check', component: HealthCheck },
  { path: '/app-docker', component: AppDocker },
  { path: '/bootstrap', component: FortiWebBootstrap },

];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
