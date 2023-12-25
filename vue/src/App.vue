<template>
  <div id="app">
    <div class="title-bar">
      <h1>FortiWeb Demo Tool</h1>
    </div>
    <n-menu :options="menuOptions" v-model:value="menuValue" />
    <router-view />
  </div>
</template>

<script>
import { NMenu } from 'naive-ui'

export default {
  components: {
    NMenu
  },
  data() {
    return {
      menuValue: 'web-protection',
      menuOptions: [
        { label: 'Web Protection', key: 'web-protection' },
        { label: 'Bot Mitigation', key: 'bot-mitigation' },
        { label: 'API Protection', key: 'api-protection' },
        { label: 'REST API', key: 'rest-api' },
        { label: 'Health Check', key: 'health-check' },
        { label: 'Configuration', key: 'configuration' }
      ]
    }
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.title-bar {
  text-align: center;
  padding: 20px;
}

nav a {
  font-weight: bold;
  color: #2c3e50;
}

nav a.router-link-exact-active {
  color: #42b983;
}
</style>
