<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Running Health Check</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>

    <div class="card-body">
      <p class="card-text"></p>

      <button class="btn btn-primary btn-sm me-2" @click="runHealthCheck" :disabled="isLoading">
        <span v-if="isLoading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
        <span>{{ isLoading ? "Checking..." : "Check" }}</span>
      </button>

      <button class="btn btn-secondary btn-sm" @click="resetResult">Reset</button>

      <div v-if="jobResult" class="mt-4">
        <h6>Result:</h6>
        <iframe ref="healthCheckIframe" :srcdoc="jobResult" @load="adjustIframeHeight" style="width: 100%; border: 0px solid lightgray"></iframe>
      </div>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Health Check</h5>
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
      jobResult: "", // Your scan result data
      isLoading: false, // Initialize isLoading
      showHelp: false,
    };
  },

  methods: {
    runHealthCheck() {
      this.isLoading = true; // Set loading state to true
      this.jobResult = ""; // Reset Result
      fetch("http://localhost:8080/health-check", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then((response) => response.text())
        .then((data) => {
                  console.log("data:", data)
          this.jobResult = data;
          this.isLoading = false; // Set loading state to false when data is received
        })
        .catch((error) => {
          console.error("Error:", error);
          this.isLoading = false; // Set loading state to false on error
        });
    },

    adjustIframeHeight() {
      const iframe = this.$refs.healthCheckIframe;
      if (iframe && iframe.contentWindow && iframe.contentWindow.document.body) {
        console.log("test")
        iframe.style.height = iframe.contentWindow.document.body.scrollHeight + 30 + "px";
      }
    },

    resetResult() {
      this.jobResult = ""; // Clear Result
    },
  },
};
</script>

<style></style>
