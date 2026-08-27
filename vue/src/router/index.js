import { createRouter, createWebHistory } from 'vue-router';

import HomePage from '../components/HomePage.vue';

const WebScan = () => import('../components/01-web-protection/WebScan.vue');
const TrafficGeneration = () => import('../components/01-web-protection/TrafficGeneration.vue');
const WebAttacks = () => import('../components/01-web-protection/WebAttacks.vue');
const MachineLearning = () => import('../components/01-web-protection/MachineLearning.vue');
const CookieSecurity = () => import('../components/01-web-protection/CookieSecurity.vue');
const CredentialStuffingDefense = () => import('../components/01-web-protection/CredentialStuffingDefense.vue');

const BotSelenium = () => import('../components/02-bot-mitigation/BotSelenium.vue');
const BotDeception = () => import('../components/02-bot-mitigation/BotDeception.vue');
const KnownBots = () => import('../components/02-bot-mitigation/KnownBots.vue');
const BotScraping = () => import('../components/02-bot-mitigation/BotScraping.vue');
const HTTPRequest = () => import('../components/02-bot-mitigation/HTTPRequest.vue');

const ApiRequests = () => import('../components/03-api-protection/ApiRequests.vue');
const ApiTrafficGeneration = () => import('../components/03-api-protection/ApiTrafficGeneration.vue');

const RestAPIManagement = () => import('../components/04-rest-api/RestAPIManagement.vue');

const ConfigForm = () => import('../components/05-tool/ConfigForm.vue');
const HealthCheck = () => import('../components/05-tool/HealthCheck.vue');
const AppDocker = () => import('../components/05-tool/AppDocker.vue');
const FortiWebBootstrap = () => import('../components/05-tool/FortiWebBootstrap.vue');

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
  { path: '/http-request', component: HTTPRequest },

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
