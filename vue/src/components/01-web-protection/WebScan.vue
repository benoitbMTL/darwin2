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
          <option value="Interesting File">Interesting File</option>
          <option value="Misconfiguration / Default File">
            Misconfiguration / Default File
          </option>
          <option value="Information Disclosure">Information Disclosure</option>
          <option value="Injection">Injection (XSS/Script/HTML)</option>
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

        <button class="btn btn-primary btn-sm me-2" @click="runScan">
          Run
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
      this.scanResult = ""; // Réinitialiser le résultat du scan
      fetch("http://localhost:8080/web-scan", {
        // Assurez-vous que l'URL correspond à votre route API
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ selectedOption: this.selectedOption }),
      })
        .then((response) => response.text())
        .then((data) => {
          this.scanResult = data; // Mettre à jour le résultat du scan
        })
        .catch((error) => console.error("Error:", error));
    },
    resetScan() {
      this.selectedOption = "All"; // Reset selected option
      this.scanResult = ""; // Clear scan result
    },
  },
};
</script>
