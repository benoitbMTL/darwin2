<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <div>
        <h5>API Traffic generation</h5>
      </div>


      <div class="d-flex align-items-center">
        <div v-if="showResetMLMessage" class="me-2">
          <div class="alert alert-dismissible fade show p-1 me-2 mb-0"
            :class="resetMLSuccess ? 'alert-success' : 'alert-danger'" role="alert"
            style="font-size: 0.875rem">
            <i class="bi me-1" :class="resetMLSuccess ? 'bi-check-circle' : 'bi-exclamation-triangle'"></i>
            {{ resetMLMessage }}
          </div>
        </div>
        <div class="me-2">
          <button type="button" class="btn btn-warning btn-sm" @click="resetApiMachineLearning">
            Reset API Machine Learning
          </button>
        </div>



        <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
        <!-- Bootstrap icon for help -->
      </div>
    </div>

    <div class="card-body">
      <p class="card-text">Launch a random traffic simulation towards the Petstore API to build FortiWeb's ML model.</p>

      <div class="d-flex align-items-center mb-3">
        <button class="btn btn-primary btn-sm" @click="generateTraffic(1)" :disabled="isTrafficLoading">
          <span v-if="isLoading1" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
          <span>{{ isLoading1 ? "Simulating..." : "Send 1 Sample" }}</span>
        </button>
        <button class="btn btn-primary btn-sm ms-2" @click="generateTraffic(10)" :disabled="isTrafficLoading">
          <span v-if="isLoading10" class="spinner-border spinner-border-sm me-2" role="status"
            aria-hidden="true"></span>
          <span>{{ isLoading10 ? "Simulating..." : "Send 10 Samples" }}</span>
        </button>
        <button class="btn btn-primary btn-sm ms-2" @click="generateTraffic(500)" :disabled="isTrafficLoading">
          <span v-if="isLoading500" class="spinner-border spinner-border-sm me-2" role="status"
            aria-hidden="true"></span>
          <span>{{ isLoading500 ? "Simulating..." : "Send 500 Samples" }}</span>
        </button>
        <button class="btn btn-secondary btn-sm ms-2" @click="resetResult">Clear display</button>
      </div>

      <div v-if="jobResult" class="alert alert-danger py-2">{{ jobResult }}</div>
      <JobMonitor ref="trafficJobMonitor" :job-id="activeJobId" job-type="api-traffic-generation"
        @finished="finishTrafficJob" />
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About API Traffic generation</h5>
    </div>
    <div class="card-body">
      <ul>
        <li>
          The machine learning based API Protection learns the REST API data structure from user traffic samples and
          then builds a mathematical model to screen
          out malicious API requests.
        </li>
        <li>
          It analyzes the method, URL, and endpoint data of the API request samples to generate an API data structure
          file for your application. This model
          describes the API data schema model of endpoint data. If the incoming API request violates the data structure,
          it will be detected as an attack.
        </li>
        <li>API Protection supports JSON request body.</li>
      </ul>

      <p>The following Machine Learning configuration provides an optimized setup for demonstrations.</p>
      <pre class="code-block"><code v-html="highlightedCode"></code></pre>

    </div>
  </div>
</template>

<script>
import hljs from "../../utils/highlight";
import JobMonitor from "../jobs/JobMonitor.vue";
import { startJob } from "../../services/jobs";

export default {
  components: { JobMonitor },
  data() {
    return {
      jobResult: "",
      activeJobId: "",
      highlightedCode: "",
      isLoading1: false,
      isLoading10: false,
      isLoading500: false,
      showHelp: false,
      sendSampleResult: "",

      resetMLMessage: "", // To store the response message
      showResetMLMessage: false, // To control the visibility of the response message
      resetMLSuccess: true,

      configSnippet:
        `config waf api-learning-policy
  edit 1
    set start-training-cnt 400
    set url-replacer-policy PETSTORE_REPLACER
    set action-mlapi alert_deny
    set schema-property maxLength minLength 
    set de-duplication-all disable
    set sample-limit-by-ip 0
  next
end`,

    };
  },

  watch: {
    jobResult(newVal) {
      if (newVal) {
        this.highlightedCode = hljs.highlightAuto(newVal).value;
      }
    },
  },

  mounted() {
    this.highlightCode(); // Call this method to apply syntax highlighting
  },

  computed: {
    isTrafficLoading() {
      return this.isLoading1 || this.isLoading10 || this.isLoading500;
    },
  },

  methods: {

    highlightCode() {
      // Use Highlight.js to apply syntax highlighting to the config snippet
      this.highlightedCode = hljs.highlightAuto(this.configSnippet).value;
    },


    // Reset API Machine Learning
    resetApiMachineLearning() {
      this.resetMLMessage = ""; // Reset message
      this.showResetMLMessage = false; // Hide message initially

      fetch("/reset-api-machine-learning", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then(async (response) => {
          const text = await response.text();
          this.resetMLSuccess = response.ok;
          this.resetMLMessage = text || `Reset failed with HTTP status ${response.status}`;
          this.showResetMLMessage = true;
        })
        .catch((error) => {
          this.resetMLSuccess = false;
          this.resetMLMessage = `Unable to contact the backend: ${error.message}`;
          this.showResetMLMessage = true;
        })
        .finally(() => {
          setTimeout(() => {
            this.showResetMLMessage = false;
          }, 5001);
        });
    },




    async generateTraffic(sampleCount) {
      this.resetResult();
      console.log(`Starting ML traffic simulation with ${sampleCount} samples...`);
      let isLoadingKey;

      switch (sampleCount) {
        case 1:
          isLoadingKey = "isLoading1";
          break;
        case 10:
          isLoadingKey = "isLoading10";
          break;
        case 500:
          isLoadingKey = "isLoading500";
          break;
      }
      this[isLoadingKey] = true;

      const formData = new URLSearchParams();
      formData.append("count", sampleCount);
      try {
        const job = await startJob("api-traffic-generation", {
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
          body: formData,
        });
        this.activeJobId = job.id;
      } catch (error) {
        this.jobResult = `Error: ${error.message}`;
        this[isLoadingKey] = false;
      }
    },

    finishTrafficJob() {
      this.isLoading1 = false;
      this.isLoading10 = false;
      this.isLoading500 = false;
    },

    resetResult() {
      this.jobResult = ""; // Clear Result
      this.activeJobId = "";
      this.finishTrafficJob();
      this.$refs.trafficJobMonitor?.clear();
    },
  },
};
</script>

<style></style>
