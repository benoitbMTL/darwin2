<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Machine Learning & Zero-Day Attacks</h5>
      <i class="bi bi-question-circle-fill bs-icon" @click="showHelp = !showHelp"></i> <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">

      <div class="row justify-content-center">

        <!-- Card #1 (Column 1) -->
        <div class="card col-md-5 me-3 align-items-center">
          <!-- Card #1 content goes here... -->


          <p class="card-text mt-3">Simulate traffic with random samples to build machine learning model.</p>
          <div class="d-flex align-items-center mb-3">
            <button class="btn btn-primary btn-sm" @click="generateTraffic(1)" :disabled="isLoading1">
              <span v-if="isLoading1" class="spinner-border spinner-border-sm me-2" role="status"
                aria-hidden="true"></span>
              <span>{{ isLoading1 ? 'Simulating...' : 'Send 1 Sample' }}</span>
            </button>
            <button class="btn btn-primary btn-sm ms-2" @click="generateTraffic(10)" :disabled="isLoading10">
              <span v-if="isLoading10" class="spinner-border spinner-border-sm me-2" role="status"
                aria-hidden="true"></span>
              <span>{{ isLoading10 ? 'Simulating...' : 'Send 10 Samples' }}</span>
            </button>
            <button class="btn btn-primary btn-sm ms-2" @click="generateTraffic(3000)" :disabled="isLoading3000">
              <span v-if="isLoading3000" class="spinner-border spinner-border-sm me-2" role="status"
                aria-hidden="true"></span>
              <span>{{ isLoading3000 ? 'Simulating...' : 'Send 3000 Samples' }}</span>
            </button>
            <button class="btn btn-secondary btn-sm ms-2" @click="resetResult">
              Reset
            </button>
          </div>
        </div>





        <!-- Card #2 (Column 2) -->
        <div class="card col-md-6 align-items-center">
          <!-- Card #2 content goes here... -->
          <p class="card-text mt-3">
            Select a user and a zero-day, then click "Run" to generate the attack scenario.
          </p>

          <div class="d-flex align-items-center mb-3 flex-wrap">
            <select class="form-select form-select-sm me-2 mb-3" v-model="selectedUser" style="width: 100px">
              <option value="admin">admin</option>
              <option value="gordonb">gordonb</option>
              <option value="1337">1337</option>
              <option value="pablo">pablo</option>
              <option value="smithy">smithy</option>
            </select>

            <select class="form-select form-select-sm me-2 mb-3" v-model="selectedAttackType" style="width: 250px">
              <option value="command_injection">Command Injection</option>
              <option value="sql_injection">SQL Injection</option>
            </select>

            <button class="btn btn-primary btn-sm me-2 mb-3" @click="performAttack">
              Run
            </button>
            <button class="btn btn-secondary btn-sm me-2 mb-3" @click="resetResult">
              Reset
            </button>
          </div>


        </div>


        <div v-if="sendSampleResult" class="mt-3">
          <h6>Simulation Result:</h6>
          <pre class="code-block"><code v-html="highlightedCode"></code></pre>
        </div>



        <div v-if="performAttackResult" class="mt-4 mb-3">
          <h6>{{ currentAttackName }} Result:</h6>
          <iframe ref="attackIframe" :srcdoc="performAttackResult" @load="adjustIframeHeight"
            style="width: 100%; border: 2px solid lightgray;"></iframe>
        </div>



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
        <li>The simulation tool generates random traffic using data from <a
            href="https://api.namefake.com/">https://api.namefake.com/</a>.</li>
        <li>The tool can send 3,000 samples to the server to simulate legitimate traffic, which is used to train
          FortiWeb's Machine Learning (ML).</li>
        <li>While only 400 requests are necessary for the Machine Learning (ML) system to build its initial model,
          ongoing
          traffic generation enables the ML to continue learning and refining its model.
        </li>
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
      isLoading1: false,
      isLoading10: false,
      isLoading3000: false,
      sendSampleResult: '',
      selectedUser: "admin",
      selectedAttackType: "command_injection",
      performAttackResult: '',
      showHelp: false,
      highlightedCode: "",
    };
  },

  watch: {
    sendSampleResult(newVal) {
      if (newVal) {
        this.highlightedCode = hljs.highlightAuto(newVal).value;
      }
    },
  },

  methods: {
    generateTraffic(sampleCount) {
      this.resetResult(); // Reset results before generating new traffic
      console.log(`Starting ML traffic simulation with ${sampleCount} samples...`);
      let isLoadingKey;
      switch (sampleCount) {
        case 1:
          isLoadingKey = 'isLoading1';
          break;
        case 10:
          isLoadingKey = 'isLoading10';
          break;
        case 3000:
          isLoadingKey = 'isLoading3000';
          break;
      }
      this[isLoadingKey] = true;
      this.sendSampleResult = ''; // Reset result

      // Make HTTP POST request to the server
      console.log('Making POST request to server');
      fetch("http://localhost:8080/machine-learning", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ sampleCount: sampleCount })
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
          this.sendSampleResult = data;
          this[isLoadingKey] = false;
        })

        .catch(error => {
          console.error('Error during fetch operation:', error);
          this.sendSampleResult = 'Error: Unable to simulate ML traffic.';
          this[isLoadingKey] = false;
        });
    },




    performAttack() {
      this.resetResult(); // Reset results before generating new traffic

      console.log('Performing attack, selected attack type:', this.selectedAttackType);

      switch (this.selectedAttackType) {
        case "command_injection":
          this.currentAttackName = 'Command Injection';
          break;
        case "sql_injection":
          this.currentAttackName = 'SQL Injection';
          break;
        default:
          this.currentAttackName = '';
      }

      console.log('Current attack name set to:', this.currentAttackName);
      this.sendAttackRequest(this.selectedAttackType);
    },

    sendAttackRequest(attackType) {
      console.log('Sending attack request with type:', attackType);

      const url = 'http://localhost:8080/web-attacks';
      const formData = new URLSearchParams();
      formData.append('type', attackType);
      formData.append('username', this.selectedUser);

      fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: formData
      })
        .then(response => {
          console.log('Received response from server:', response);
          return response.text();
        })
        .then(html => {
          console.log('Attack simulation successful:', html);
          this.performAttackResult = html;
        })
        .catch(error => {
          console.error('Error during attack simulation:', error);
          this.performAttackResult = 'Failed to perform attack';
        });
    },


    adjustIframeHeight() {
      const iframe = this.$refs.attackIframe;
      if (iframe && iframe.contentWindow && iframe.contentWindow.document.body) {
        iframe.style.height = (iframe.contentWindow.document.body.scrollHeight + 30) + 'px';
      }
    },





    resetResult() {
      console.log('Resetting Result');
      this.selectedOption = "All"; // Reset selected option
      this.sendSampleResult = '';
      this.performAttackResult = '';
    },
  },
};
</script>

<style></style>
