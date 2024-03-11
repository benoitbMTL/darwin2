<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Web Vulnerability Scanner</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem;" @click="showHelp = !showHelp"></i> <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">Select a Nikto scan option to analyze web server security and identify potential
        vulnerabilities.</p>
      <div class="d-flex align-items-center mb-3">
        <select class="form-select form-select-sm me-2" v-model="selectedOption" style="width: 350px">
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

        <button class="btn btn-secondary btn-sm" @click="resetResult">
          Reset
        </button>
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
      <h5>About Nikto Web Vulnerability Scanner</h5>
    </div>
    <div class="card-body">
      <p><strong>Nikto 2.5</strong> is a comprehensive web server scanner used for security testing. It's an open-source
        tool designed to perform extensive checks against web servers.</p>
  <p>
    <a href="https://github.com/sullo/nikto" class="d-inline-flex align-items-center" target="_blank" rel="noopener noreferrer">
      <i class="bi bi-github bs-icon me-2"></i>https://github.com/sullo/nikto
    </a>
  </p>

   <ul>
      <li><strong>Comprehensive Tests:</strong> Checks over 7,000 files/programs and more than 1250 server types for potential risks and outdated versions.</li>
      <li><strong>Configuration Checks:</strong> Inspects server configurations and identifies installed web servers and software.</li>
      <li><strong>Updates and Plugins:</strong> Nikto's scan items and plugins are regularly updated to detect new vulnerabilities.</li>
      <li><strong>Visibility:</strong> Detectable in logs and by IPS/IDS systems, but offers some anti-IDS techniques.</li>
      <li><strong>Informational Checks:</strong> Reports both security flaws and 'information only' items for webmasters and security engineers.</li>
    </ul>
    <p>The valid tuning options are:</p>
    <ul>
      <li><strong>File Upload:</strong> Exploits which allow a file to be uploaded to the target server.</li>
      <li><strong>Interesting File / Seen in logs:</strong> An unknown but suspicious file or attack that has been seen in web server logs.</li>
      <li><strong>Misconfiguration / Default File:</strong> Default files or files which have been misconfigured. This could include documentation, or a resource which should be password protected.</li>
      <li><strong>Information Disclosure:</strong> A resource that reveals information about the target, such as a file system path or account name.</li>
      <li><strong>Injection (XSS/Script/HTML):</strong> Any form of injection, including cross-site scripting (XSS) or content (HTML) injection. Does not include command injection.</li>
      <li><strong>Remote File Retrieval - Inside Web Root:</strong> Allows remote users to retrieve unauthorized files from within the web server's root directory.</li>
      <li><strong>Denial of Service:</strong> Resource permits a denial of service against the application, web server, or host.</li>
      <li><strong>Remote File Retrieval - Server Wide:</strong> Allows remote users to retrieve unauthorized files from anywhere on the target.</li>
      <li><strong>Command Execution / Remote Shell:</strong> Permits the execution of system commands or spawning of a remote shell.</li>
      <li><strong>SQL Injection:</strong> Any attack allowing SQL execution against a database.</li>
      <li><strong>Authentication Bypass:</strong> Allows clients to access resources they should not have access to.</li>
      <li><strong>Software Identification:</strong> Identifies installed software or programs.</li>
        <li><strong>Remote Source Inclusion:</strong> Allows remote inclusion of source code.</li>
        <li><strong>Web Service:</strong> Checks for issues that can affect APIs, SOAP services, and other web service endpoints.</li>
        <li><strong>Administrative Console:</strong> Vulnerabilities related to administrative consoles.</li>
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
      selectedOption: "All",
      jobResult: "", // Your scan result data
      highlightedCode: "",
      isLoading: false, // Initialize isLoading
      showHelp: false,
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
    runScan() {
      this.isLoading = true; // Set loading state to true
      this.jobResult = ""; // Reset Result
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

      fetch("localhost:8080/web-scan", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ selectedOption: tuningFlag }),
      })
        .then((response) => response.text())
        .then((data) => {
          this.jobResult = data;
          this.isLoading = false; // Set loading state to false when data is received
        })
        .catch((error) => {
          console.error("Error:", error);
          this.isLoading = false; // Set loading state to false on error
        });
    },

    resetResult() {
      this.selectedOption = "All"; // Reset selected option
      this.jobResult = ""; // Clear Result
    },
  },
};
</script>

<style></style>
