<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Biometrics-Based Detection</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>

    <div class="card-body">
      <p class="card-text">
        Text.
      </p>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Biometrics-Based Detection</h5>
    </div>
    <div class="card-body">
      <p>By checking the client events such as mouse movement, keyboard, screen touch, and scroll, etc in specified
        period, FortiWeb judges whether the request comes from a human or from a bot. You can configure the
        biometrics
        based detection rule to define the client event, collection period, and the request URL, etc.</p>
    </div>
  </div>

</template>

<script>
export default {
  data() {
    return {
      pageSource: "",
      jobResult: "",
      showHelp: false,
      highlightedCode: "",
    };
  },

  methods: {
    viewPageSource() {
      fetch("/bot-page-source", {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
      })
        .then((response) => response.text())
        .then((html) => {
          // Process and update pageSource and highlightedCode with last 10 lines
          const lines = html.split("\n");
          const lastTenLines = lines.slice(-10).join("\n");
          this.highlightedCode = this.escapeHtml(lastTenLines);
          this.pageSource = true;
        })
        .catch((error) => {
          console.error("Error:", error);
        });
    },

    performBotDeception() {
      // Ensure this method matches the button's @click assignment
      fetch("/bot-biometric", {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
      })
        .then((response) => response.text())
        .then((html) => {
          this.jobResult = html;
        })
        .catch((error) => {
          console.error("Error:", error);
          this.jobResult = "Failed to perform Bot Deception";
        });
    },

    adjustIframeHeight() {
      const iframe = this.$refs.botIframe;
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
      this.jobResult = "";
      this.pageSource = "";
      this.highlightedCode = "";
    },

    escapeHtml(unsafe) {
      return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
    },
  },
};
</script>

<style></style>
