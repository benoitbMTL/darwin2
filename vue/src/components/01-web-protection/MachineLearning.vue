<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Machine Learning Traffic Simulation</h5>
      <i class="bi bi-question-circle-fill bs-icon" @click="showHelp = !showHelp"></i> <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">Simulate machine learning data traffic for analysis and testing.</p>
      <div class="d-flex align-items-center mb-3">
        <button class="btn btn-primary btn-sm" @click="generateTraffic" :disabled="isLoading">
          <span v-if="isLoading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
          <span>{{ isLoading ? 'Simulating...' : 'Simulate Traffic' }}</span>
        </button>
        <button class="btn btn-secondary btn-sm ms-2" @click="resetResult">
          Reset
        </button>
      </div>
      <div v-if="jobResult" class="mt-3">
        <h6>Simulation Result:</h6>
        <pre class="code-block"><code v-html="highlightedCode"></code></pre>
      </div>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Machine Learning Traffic Simulation</h5>
    </div>
    <div class="card-body">
      <ul>
          <li>The simulation tool generates random traffic using data from <a href="https://api.namefake.com/">https://api.namefake.com/</a>.</li>
          <li>The tool sends 3,000 requests to the server to simulate legitimate traffic for training FortiWeb's Machine Learning (ML) algorithms.</li>
          <li>While only 400 requests are needed for the ML to build its initial mathematical model, continued traffic generation allows the ML to further learn and refine its model.</li>
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
      console.log('Starting ML traffic simulation...');
      this.isLoading = true;
      this.jobResult = ''; // Reset result

      // Make HTTP POST request to the server
      console.log('Making POST request to server');
      fetch("http://localhost:8080/machine-learning", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then(response => {
          console.log('Received response from server:', response);
          if (!response.ok) {
            console.error('Network response was not ok', response);
            throw new Error('Network response was not ok');
          }
          return response.text();
        })
        .then(data => {
          console.log('ML traffic simulation successful:', data);
          this.jobResult = data;
          this.isLoading = false;
        })
        .catch(error => {
          console.error('Error during fetch operation:', error);
          this.jobResult = 'Error: Unable to simulate ML traffic.';
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

<style>
.bs-icon {
  cursor: pointer;
  font-size: 1.5em;
}
</style>
