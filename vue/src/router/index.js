import { createRouter, createWebHistory } from 'vue-router'
import WebProtection from '../views/WebProtection.vue'
import BotMitigation from '../views/BotMitigation.vue'
import ApiProtection from '../views/ApiProtection.vue'
import RestApi from '../views/RestApi.vue'
import Configuration from '../views/Configuration.vue'
import HealthCheck from '../views/HealthCheck.vue'

const routes = [
  {
    path: '/',
    name: 'web-protection',
    component: WebProtection
  },
  {
    path: '/bot-mitigation',
    name: 'bot-mitigation',
    component: BotMitigation
  },
  {
    path: '/api-protection',
    name: 'api-protection',
    component: ApiProtection
  },
  {
    path: '/rest-api',
    name: 'rest-api',
    component: RestApi
  },
  {
    path: '/health-check',
    name: 'health-check',
    component: HealthCheck
  },
  {
    path: '/configuration',
    name: 'configuration',
    component: Configuration
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
