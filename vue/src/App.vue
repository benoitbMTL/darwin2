<template>
  <div id="app">
    <div class="title-bar">
      <h1>FortiWeb Demo Tool</h1>
    </div>
    <n-menu :options="menuOptions" :value="menuValue" @update:value="handleMenuSelect" />
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
  },
  methods: {
    handleMenuSelect(key) {
      this.$router.push({ name: key });
    }
  }
}
</script>

<style>

</style>
