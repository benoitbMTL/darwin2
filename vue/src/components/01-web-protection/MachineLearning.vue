<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>MachineLearning</h5>
      <i class="bi bi-question-circle-fill bs-icon" @click="showHelp = !showHelp"></i> <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">Machine Learning Options.</p>
    </div>
  </div>
  
  <!-- Help Card -->
  <div v-if="showHelp" class="card mt-3">
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

<style>
.bs-icon {
  cursor: pointer;
  font-size: 1.5em;
}
</style>
