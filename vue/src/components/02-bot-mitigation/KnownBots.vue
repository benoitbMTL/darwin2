<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Known Bots</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">Select a Bot.</p>

      <div class="d-flex align-items-center mb-3 flex-wrap">
        <button class="btn btn-primary btn-sm me-2" @click="sendBot('DoS')">DoS</button>
        <button class="btn btn-primary btn-sm me-2" @click="sendBot('Spam')">Spam</button>
        <button class="btn btn-primary btn-sm me-2" @click="sendBot('Trojan')">Trojan</button>
        <button class="btn btn-primary btn-sm me-2" @click="sendBot('Scanner')">Scanner</button>
        <button class="btn btn-primary btn-sm me-2" @click="sendBot('Crawler')">Crawler</button>
        <button class="btn btn-primary btn-sm me-2" @click="sendBot('SearchEngine')">Search Engine</button>

        <button class="btn btn-secondary btn-sm me-2" @click="resetResult">Reset</button>
      </div>

      <div v-if="jobResult" class="mt-4 mb-3">
        <h6>{{ currentBotName }} Result:</h6>
        <iframe ref="botIframe" :srcdoc="jobResult" @load="adjustIframeHeight" style="width: 100%; border: 1px solid lightgray"></iframe>
      </div>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Known Bots</h5>
    </div>
    <div class="card-body">
      <ul>
        <li>This tool sends Bot from a specific family.</li>
        <li>
          Known Bots feature protects websites, mobile applications, and APIs from malicious bots such as DoS, Spam, Crawler, etc... without affecting the flow
          of critical traffic.
        </li>
        <li>This feature identifies and manages a wide range of attacks from automated tools no matter where these applications or APIs are deployed.</li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      jobResult: "",
      showHelp: false,
    };
  },

  methods: {
    sendBot(botName) {
      const url = "http://localhost:8080/known-bots";
      const formData = new URLSearchParams();
      formData.append("name", botName);

      fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        body: formData.toString(),
      })
        .then((response) => response.text())
        .then((html) => {
          this.jobResult = html;
        })
        .catch((error) => {
          console.error("Error:", error);
          this.jobResult = "Failed to send Bot";
        });
    },

    adjustIframeHeight() {
      const iframe = this.$refs.botIframe;
      if (iframe && iframe.contentWindow && iframe.contentWindow.document.body) {
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
