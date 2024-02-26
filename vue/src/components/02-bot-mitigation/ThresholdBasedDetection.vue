<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Threshold Based Detection</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem;" @click="showHelp = !showHelp"></i> <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">Threshold Based Detection.</p>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Threshold Based Detection</h5>
    </div>
    <div class="card-body">
    <p>You can configure threshold based detection rules to define occurrence, time period, severity, and trigger policy, etc of the following suspicious behaviors, and thus FortiWeb judges whether the request comes from a human or a bot.</p>
    <ul>
        <li>Crawler</li>
        <li>Vulnerability Scanning</li>
        <li>Slow Attack</li>
        <li>Content Scraping</li>
        <li>Illegal User Scan</li>
    </ul>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      showHelp: false,
    };
  },
};
</script>

<style></style>

