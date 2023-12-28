<template>
  <div class="card my-4">
    <h5 class="card-header">Web Vulnerability Scanner</h5>
    <div class="card-body">
      <p class="card-text">Choose a scan test to perform on your target.</p>
      <div class="d-flex align-items-center mb-3">
        <select
          class="form-select form-select-sm me-2"
          v-model="selectedOption"
          style="width: 350px"
        >
          <option value="All">All</option>
          <option value="Interesting File / Seen in logs">
            Interesting File / Seen in logs
          </option>
          <option value="Misconfiguration / Default File">
            Misconfiguration / Default File
          </option>
          <option value="Information Disclosure">Information Disclosure</option>
          <option value="Injection (XSS/Script/HTML)">
            Injection (XSS/Script/HTML)
          </option>
          <option value="Remote File Retrieval - Inside Web Root">
            Remote File Retrieval - Inside Web Root
          </option>
          <option value="Denial of Service">Denial of Service</option>
          <option value="Remote File Retrieval - Server Wide">
            Remote File Retrieval - Server Wide
          </option>
          <option value="Command Execution / Remote Shell">
            Command Execution / Remote Shell
          </option>
          <option value="SQL Injection">SQL Injection</option>
          <option value="File Upload">File Upload</option>
          <option value="Authentication Bypass">Authentication Bypass</option>
          <option value="Software Identification">
            Software Identification
          </option>
          <option value="Remote Source Inclusion">
            Remote Source Inclusion
          </option>
          <option value="WebService">WebService</option>
          <option value="Administrative Console">Administrative Console</option>
        </select>

        <button class="btn btn-primary btn-sm me-2" @click="runScan" :disabled="isLoading">
          <span v-if="isLoading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
          <span>{{ isLoading ? 'Scanning...' : 'Scan' }}</span>
        </button>

        <button class="btn btn-secondary btn-sm" @click="resetScan">
          Reset
        </button>
      </div>

      <div v-if="scanResult" class="mt-3">
        <h6>Scan Result:</h6>
        <pre class="code-block"><code v-html="highlightedCode"></code></pre>
      </div>
    </div>
  </div>
</template>

<script>
import hljs from "highlight.js";
import "highlight.js/styles/monokai.css"; // Monokai theme for Highlight.js

export default {
  data() {
    return {
      selectedOption: "All",
      scanResult: "", // Your scan result data
      highlightedCode: "",
      isLoading: false, // Initialize isLoading
    };
  },
  watch: {
    scanResult(newVal) {
      if (newVal) {
        this.highlightedCode = hljs.highlightAuto(newVal).value;
      }
    },
  },
  methods: {
    runScan() {
      this.isLoading = true; // Set loading state to true
      this.scanResult = ""; // Reset scan result
      const niktoTuningFlags = {
        All: "1234567890abcde",
        "Interesting File / Seen in logs": "1",
        "Misconfiguration / Default File": "2",
        "Information Disclosure": "3",
        "Injection (XSS/Script/HTML)": "4",
        "Remote File Retrieval - Inside Web Root": "5",
        "Denial of Service": "6",
        "Remote File Retrieval - Server Wide": "7",
        "Command Execution / Remote Shell": "8",
        "SQL Injection": "9",
        "File Upload": "0",
        "Authentication Bypass": "a",
        "Software Identification": "b",
        "Remote Source Inclusion": "c",
        WebService: "d",
        "Administrative Console": "e",
      };

      const tuningFlag = niktoTuningFlags[this.selectedOption] || "";

      fetch("http://localhost:8080/web-scan", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ selectedOption: tuningFlag }),
      })
        .then((response) => response.text())
        .then((data) => {
          this.scanResult = data;
          this.isLoading = false; // Set loading state to false when data is received
        })
        .catch((error) => {
          console.error("Error:", error);
          this.isLoading = false; // Set loading state to false on error
        });
    },

    resetScan() {
      this.selectedOption = "All"; // Reset selected option
      this.scanResult = ""; // Clear scan result
    },
  },
};
</script>
