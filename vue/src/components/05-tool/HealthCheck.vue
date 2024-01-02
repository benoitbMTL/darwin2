<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Running Applications in Docker</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem;" @click="showHelp = !showHelp"></i> <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">Machine Learning Options.</p>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Docker</h5>
    </div>
    <div class="card-body">
      <p>Help content goes here...</p>
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