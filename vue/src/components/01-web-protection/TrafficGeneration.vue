<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Traffic Generation</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem;" @click="showHelp = !showHelp"></i> <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">Generate random cyber attacks from various public IP addresses.</p>
      <div class="d-flex align-items-center mb-3">

        <button class="btn btn-primary btn-sm" @click="generateTraffic" :disabled="isLoading">
          <span v-if="isLoading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
          <span>{{ isLoading ? 'Generating...' : 'Generate Traffic' }}</span>
        </button>

        <button class="btn btn-secondary btn-sm ms-2" @click="resetResult">
          Reset
        </button>
      </div>

      <div v-if="jobResult" class="mt-3">
        <h6>Traffic Result:</h6>
        <pre class="code-block"><code v-html="highlightedCode"></code></pre>
      </div>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Traffic Generation</h5>
    </div>
    <div class="card-body">
      <ul>
        <li>The Traffic Generator simulates cyber attacks using randomly generated public IP addresses.</li>
        <li>Nikto is utilized by the Traffic Generator to generate random attacks.</li>
        <li>Launch the Traffic Generator before a demonstration to populate FortiWeb logs and FortiView dashboards.</li>
      </ul>

    </div>
  </div>
</template>

<script>
import hljs from "highlight.js";
import "highlight.js/styles/monokai.css"; // Monokai theme for Highlight.js

export default {
  data() {
    return {
      isLoading: false,
      jobResult: '',
      showHelp: false,
      highlightedCode: "",
    };
  },

  watch: {
    jobResult(newVal) {
      if (newVal) {
        this.highlightedCode = hljs.highlightAuto(newVal).value;
      }
    },
  },
  methods: {
    generateTraffic() {
      console.log('Starting traffic generation...');
      this.isLoading = true;
      this.jobResult = ''; // Reset windows result

      // Make HTTP POST request to the server
      fetch("localhost:8080/traffic-generation", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then(response => {
          console.log('Received response from server');
          if (!response.ok) {
            console.error('Network response was not ok');
            throw new Error('Network response was not ok');
          }
          return response.text();
        })
        .then(data => {
          console.log('Traffic generation successful:', data);
          this.jobResult = data;
          this.isLoading = false;
        })
        .catch(error => {
          console.error('Error during fetch operation:', error);
          this.jobResult = 'Error: Unable to generate traffic.';
          this.isLoading = false;
        });
    },
    resetResult() {
      console.log('Resetting Result');
      this.jobResult = '';
    },
  },
};
</script>

<style></style>

