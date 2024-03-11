<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <div>
        <h5>API Traffic generation</h5>
      </div>


      <div class="d-flex align-items-center">
        <div v-if="showResetMLMessage" class="me-2">
          <div class="alert alert-success alert-dismissible fade show p-1 me-2 mb-0" role="alert"
            style="font-size: 0.875rem">
            <i class="bi bi-check-circle me-1"></i> {{ resetMLMessage }}
          </div>
        </div>
        <div class="me-2">
          <button type="button" class="btn btn-warning btn-sm" @click="resetApiMachineLearning">
            Reset Machine Learning
          </button>
        </div>



        <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
        <!-- Bootstrap icon for help -->
      </div>
    </div>

    <div class="card-body">
      <p class="card-text">Launch a random traffic simulation towards the Petstore API to build FortiWeb's ML model.</p>

      <div class="d-flex align-items-center mb-3">
        <button class="btn btn-primary btn-sm" @click="generateTraffic(1)" :disabled="isLoading1">
          <span v-if="isLoading1" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
          <span>{{ isLoading1 ? "Simulating..." : "Send 1 Sample" }}</span>
        </button>
        <button class="btn btn-primary btn-sm ms-2" @click="generateTraffic(10)" :disabled="isLoading10">
          <span v-if="isLoading10" class="spinner-border spinner-border-sm me-2" role="status"
            aria-hidden="true"></span>
          <span>{{ isLoading10 ? "Simulating..." : "Send 10 Samples" }}</span>
        </button>
        <button class="btn btn-primary btn-sm ms-2" @click="generateTraffic(500)" :disabled="isLoading500">
          <span v-if="isLoading500" class="spinner-border spinner-border-sm me-2" role="status"
            aria-hidden="true"></span>
          <span>{{ isLoading500 ? "Simulating..." : "Send 500 Samples" }}</span>
        </button>
        <button class="btn btn-secondary btn-sm ms-2" @click="resetResult">Reset</button>
      </div>

      <div v-if="jobResult" class="mt-3">
        <h6>Scan Result:</h6>
        <pre class="code-block"><code v-html="highlightedCode"></code></pre>
      </div>
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
import hljs from "highlight.js";
import "highlight.js/styles/monokai.css"; // Monokai theme for Highlight.js

export default {
  data() {
    return {
      jobResult: "",
      highlightedCode: "",
      isLoading1: false,
      isLoading10: false,
      isLoading500: false,
      showHelp: false,
      sendSampleResult: "",

      resetMLMessage: "", // To store the response message
      showResetMLMessage: false, // To control the visibility of the response message

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
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.text();
        })
        .then((text) => {
          this.resetMLMessage = text; // Set the response message
          this.showResetMLMessage = true; // Show message

          setTimeout(() => {
            this.showResetMLMessage = false; // Hide message after 15 seconds
          }, 5001);
        })
        .catch((error) => console.error("Error:", error));
    },




    generateTraffic(sampleCount) {
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

      // Constructing form data
      const formData = new URLSearchParams();
      formData.append("count", sampleCount);

      // Make HTTP POST request to the server
      console.log("Making POST request to server");
      fetch("/api-traffic-generation", {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        body: formData,
      })
        .then((response) => {
          console.log("Received response from server:", response);
          if (!response.ok) {
            console.error("Network response was not ok", response);
            throw new Error("Network response was not ok");
          }
          return response.text();
        })

        .then((data) => {
          console.log("ML traffic simulation successful:", data);
          this.jobResult = data; // Update this line
          this[isLoadingKey] = false;
        })

        .catch((error) => {
          console.error("Error during fetch operation:", error);
          this.jobResult = "Error: Unable to simulate ML traffic."; // Update this line
          this[isLoadingKey] = false;
        });
    },

    resetResult() {
      this.jobResult = ""; // Clear Result
    },
  },
};
</script>

<style></style>
