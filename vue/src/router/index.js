import { createRouter, createWebHistory } from 'vue-router';

// Import components from each folder
import WebScan from '../components/01-web-protection/WebScan.vue';
import TrafficGeneration from '../components/01-web-protection/TrafficGeneration.vue';
import WebAttacks from '../components/01-web-protection/WebAttacks.vue';
import MachineLearning from '../components/01-web-protection/MachineLearning.vue';
import CookieSecurity from '../components/01-web-protection/CookieSecurity.vue';
import CredentialStuffingDefense from '../components/01-web-protection/CredentialStuffingDefense.vue';

import BiometricsBasedDetection from '../components/02-bot-mitigation/BiometricsBasedDetection.vue';
import BotDeception from '../components/02-bot-mitigation/BotDeception.vue';
import KnownBots from '../components/02-bot-mitigation/KnownBots.vue';
import MLBasedBotDetection from '../components/02-bot-mitigation/MLBasedBotDetection.vue';
import ThresholdBasedDetection from '../components/02-bot-mitigation/ThresholdBasedDetection.vue';

import ApiRequests from '../components/03-api-protection/ApiRequests.vue';
import MLBasedAPIProtection from '../components/03-api-protection/MLBasedAPIProtection.vue';

import CreateNewApplicationPolicy from '../components/04-rest-api/CreateNewApplicationPolicy.vue';
import DeleteApplicationPolicy from '../components/04-rest-api/DeleteApplicationPolicy.vue';

import AppConfiguration from '../components/05-tool/AppConfiguration.vue';
import HealthCheck from '../components/05-tool/HealthCheck.vue';

const routes = [
  { path: '/web-scan', component: WebScan },
  { path: '/traffic-generation', component: TrafficGeneration },
  { path: '/web-attacks', component: WebAttacks },
  { path: '/machine-learning', component: MachineLearning },
  { path: '/cookie-security', component: CookieSecurity },
  { path: '/credential-stuffing-defense', component: CredentialStuffingDefense },

  { path: '/biometrics-based-detection', component: BiometricsBasedDetection },
  { path: '/bot-deception', component: BotDeception },
  { path: '/known-bots', component: KnownBots },
  { path: '/ml-based-bot-detection', component: MLBasedBotDetection },
  { path: '/threshold-based-detection', component: ThresholdBasedDetection },

  { path: '/api-requests', component: ApiRequests },
  { path: '/ml-based-api-protection', component: MLBasedAPIProtection },

  { path: '/create-new-application-policy', component: CreateNewApplicationPolicy },
  { path: '/delete-application-policy', component: DeleteApplicationPolicy },

  { path: '/configuration', component: AppConfiguration },
  { path: '/health-check', component: HealthCheck },

];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
