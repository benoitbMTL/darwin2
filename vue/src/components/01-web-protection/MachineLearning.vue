<template>
  <div class="card my-4">

    <!-- HEADER -->
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Machine Learning Anomaly Detection</h5>



      <div class="d-flex align-items-center">
        <div v-if="showResetMLMessage" class="me-2">
          <div class="alert alert-success alert-dismissible fade show p-1 me-2 mb-0" role="alert"
            style="font-size: 0.875rem">
            <i class="bi bi-check-circle me-1"></i> {{ resetMLMessage }}
          </div>
        </div>
        <div class="me-2">
          <button type="button" class="btn btn-warning btn-sm" @click="resetMachineLearning">
            Reset Machine Learning
          </button>
        </div>

        <div>
          <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
          <!-- Bootstrap icon for help -->
        </div>
      </div>
    </div>

    <!-- BODY -->
    <div class="card-body">
      <div class="container">

        <p>
          This tool will help protect our
          <a :href="bankingDynamicUrl" target="_blank">banking form</a>
          with machine learning, thereby blocking zero-day attacks and reducing
          the attack surface.
        </p>

        <div class="row">

          <!-- COL 1 -->
          <div class="col-md-6">
            <!-- CARD -->
            <div class="card">
              <div class="card-header">

                <strong>Build Model</strong>
              </div>

              <!-- CARD BODY -->
              <div class="card-body">

                <p>
                  Simulate legitimate traffic with random samples to build a ML model.
                </p>
                <div class="d-flex mb-3"> <!-- TEST -->

                  <button class="btn btn-primary btn-sm" @click="generateTraffic(1)" :disabled="isLoading1">
                    <span v-if="isLoading1" class="spinner-border spinner-border-sm me-2" role="status"
                      aria-hidden="true"></span>
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
                  <button class="btn btn-secondary btn-sm ms-2" @click="resetResult">
                    Reset
                  </button>
                </div>

              </div>
            </div>
          </div> <!-- COL -->

          <!-- COL 2 -->
          <div class="col-md-6">
            <div class="card">
              <div class="card-header">

                <strong>Evaluate Model</strong>
              </div>
              <div class="card-body">

                <p>
                  Select a zero-day attack scenario and evaluate ML anomaly detection.
                </p>
                <div class="d-flex"> <!-- TEST -->

                  <!-- Menu -->
                  <select class="form-select form-select-sm me-2 mb-3" v-model="selectedAttackType"
                    style="width: 350px">
                    <option value="zero_day_sqli_1">
                      Zero Day SQL Injection: A 'DIV' B
                    </option>
                    <option value="zero_day_sqli_2">
                      Zero Day SQL Injection: A '^' B
                    </option>
                    <option value="zero_day_sqli_3">
                      Zero Day SQL Injection: 3)+1+(0
                    </option>
                    <option value="zero_day_sqli_4">
                      Zero Day SQL Injection: 3||1
                    </option>
                    <option value="zero_day_remote_exploit_1">
                      Zero Day Remote Exploits: %X%X%X%X%X
                    </option>
                    <option value="zero_day_remote_exploit_2">
                      Zero Day Remote Exploits: %p%p%p%p%p
                    </option>
                    <option value="zero_day_command_injection_1">
                      Zero Day Command Injection: /???/l?
                    </option>
                    <option value="zero_day_command_injection_2">
                      Zero Day Command Injection: var1=l var2=s
                    </option>
                    <option value="zero_day_xss_1">
                      Zero Day Cross Site Scripting: window['ale'+'rt'](1)
                    </option>
                    <option value="zero_day_xss_2">
                      Zero Day Cross Site Scripting: ___=1?'ert(123)
                    </option>
                  </select>

                  <button class="btn btn-primary btn-sm me-2 mb-3" @click="performAttack">
                    Run
                  </button>
                  <button class="btn btn-secondary btn-sm me-2 mb-3" @click="resetResult">
                    Reset
                  </button>
                </div>

              </div>
            </div>



          </div> <!-- COL -->
        </div> <!-- Row -->


        <div v-if="sendSampleResult" class="mt-3">
          <h6>Simulation Result:</h6>
          <pre class="code-block"><code v-html="highlightedCode"></code></pre>
        </div>

        <div v-if="performAttackResult" class="mt-4 mb-3">
          <h6>{{ currentAttackName }} Result:</h6>
          <iframe ref="attackIframe" :srcdoc="performAttackResult" @load="adjustIframeHeight"
            style="width: 100%; border: 1px solid lightgray"></iframe>
        </div>


      </div> <!-- Container -->

    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Machine Learning Traffic Simulation</h5>
    </div>
    <div class="card-body">
      <ul>
        <li>
          The simulation tool generates random traffic using data from
          <a href="https://api.namefake.com/">https://api.namefake.com/</a>.
        </li>
        <li>
          The tool sends random samples to the server to simulate legitimate
          traffic, which is used to train FortiWeb's Machine Learning (ML).
        </li>
        <li>
          While only 400 requests are necessary for the Machine Learning (ML)
          system to build its initial model, ongoing traffic generation enables
          the ML to continue learning and refining its model.
        </li>
        <li>
          Press the "Reset Machine Learning" button to delete all existing
          learning results and start the demo from scratch.
        </li>
      </ul>
      <p>
        The following Machine Learning configuration provides an optimized setup
        for demonstrations.
      </p>

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
      highlightedCode: "",
      isLoading1: false,
      isLoading10: false,
      isLoading500: false,
      sendSampleResult: "",
      selectedAttackType: "zero_day_sqli_1",
      performAttackResult: "",
      showHelp: false,
      config: {
        BANKURL: "",
      },
      resetMLMessage: "", // To store the response message
      showResetMLMessage: false, // To control the visibility of the response message

      configSnippet: `config waf machine-learning-policy
  edit 1
    set sample-limit-by-ip 0
    set ip-expire-cnts 1
    set ip-expire-intval 1
    set svm-type extended
  next
end`,
    };
  },

  mounted() {
    this.fetchConfig(); // Fetch config on component mount
    this.highlightCode();
  },
  
  computed: {
    bankingDynamicUrl() {
      if (this.config.FABRICLABSTORY) {
        return `https://bank.${this.config.FABRICLABSTORY}.fabriclab.ca`;
      } else {
        return this.config.BANKURL;
      }
    }
  },

  watch: {
    sendSampleResult(newVal) {
      if (newVal) {
        this.highlightedCode = hljs.highlightAuto(newVal).value;
      }
    },
  },

  methods: {
    highlightCode() {
      // Use Highlight.js to apply syntax highlighting to the config snippet
      this.highlightedCode = hljs.highlightAuto(this.configSnippet).value;
    },

    fetchConfig() {
      // Fetch config from server
      fetch("/config")
        .then((response) => response.json())
        .then((data) => {
          this.config = data; // Update config with fetched data
        })
        .catch((error) => {
          console.error("Error fetching config:", error);
        });
    },

    resetMachineLearning() {
      this.resetMLMessage = "";
      this.showResetMLMessage = false;

      fetch("/reset-machine-learning", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then((response) => {
          if (!response.ok) {
            // Assuming the server responds with a JSON payload in case of errors
            return response.json().then(error => Promise.reject(error));
          }
          return response.text();
        })
        .then((text) => {
          this.resetMLMessage = text;
          this.showResetMLMessage = true;

          setTimeout(() => {
            this.showResetMLMessage = false;
          }, 6000);
        })
        .catch((error) => {
          console.error("Error:", error);
          // Display the error message from the server or a default message
          this.resetMLMessage = error.message || "An error occurred during the reset process.";
          this.showResetMLMessage = true;

          setTimeout(() => {
            this.showResetMLMessage = false;
          }, 6000);
        });
    },

    generateTraffic(sampleCount) {
      this.resetResult(); // Reset results before generating new traffic
      console.log(
        `Starting ML traffic simulation with ${sampleCount} samples...`
      );
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
      this.sendSampleResult = ""; // Reset result

      // Make HTTP POST request to the server
      console.log("Making POST request to server");
      fetch("/machine-learning", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ sampleCount: sampleCount }),
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
          this.sendSampleResult = data;
          this[isLoadingKey] = false;
        })

        .catch((error) => {
          console.error("Error during fetch operation:", error);
          this.sendSampleResult = "Error: Unable to simulate ML traffic.";
          this[isLoadingKey] = false;
        });
    },

    performAttack() {
      this.resetResult(); // Reset results before generating new traffic

      console.log(
        "Performing attack, selected attack type:",
        this.selectedAttackType
      );

      switch (this.selectedAttackType) {
        case "zero_day_sqli_1":
        case "zero_day_sqli_2":
        case "zero_day_sqli_3":
        case "zero_day_sqli_4":
          this.currentAttackName = "Zero Day SQL Injection";
          break;
        case "zero_day_remote_exploit_1":
        case "zero_day_remote_exploit_2":
          this.currentAttackName = "Zero Day Remote Exploits";
          break;
        case "zero_day_command_injection_1":
        case "zero_day_command_injection_2":
          this.currentAttackName = "Zero Day Command Injection";
          break;
        case "zero_day_xss_1":
        case "zero_day_xss_2":
          this.currentAttackName = "Zero Day Cross Site Scripting";
          break;
        default:
          this.currentAttackName = "";
      }

      console.log("Current attack name set to:", this.currentAttackName);
      this.sendAttackRequest(this.selectedAttackType);
    },

    sendAttackRequest(attackType) {
      console.log("Sending attack request with type:", attackType);

      const url = "/web-attacks";
      const formData = new URLSearchParams();
      formData.append("type", attackType);

      fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        body: formData,
      })
        .then((response) => {
          console.log("Received response from server:", response);
          return response.text();
        })
        .then((html) => {
          console.log("Attack simulation successful:", html);
          this.performAttackResult = html;
        })
        .catch((error) => {
          console.error("Error during attack simulation:", error);
          this.performAttackResult = "Failed to perform attack";
        });
    },

    adjustIframeHeight() {
      const iframe = this.$refs.attackIframe;
      if (
        iframe &&
        iframe.contentWindow &&
        iframe.contentWindow.document.body
      ) {
        iframe.style.height =
          iframe.contentWindow.document.body.scrollHeight + 30 + "px";
      }
    },

    resetResult() {
      console.log("Resetting Result");
      this.selectedOption = "All"; // Reset selected option
      this.sendSampleResult = "";
      this.performAttackResult = "";
    },
  },
};
</script>

<style></style>
