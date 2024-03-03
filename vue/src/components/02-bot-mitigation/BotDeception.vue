<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Bot Deception</h5>
      <i
        class="bi bi-question-circle-fill bs-icon"
        style="font-size: 1.5rem"
        @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">
        Bot Deception inserts invisible links in HTML responses to distinguish
        between regular clients and malicious bots like web crawlers.
      </p>

      <button class="btn btn-primary btn-sm me-2" @click="viewPageSource">
        View Page Source AAAAA
      </button>
      <button class="btn btn-primary btn-sm me-2" @click="performBotDeception">
        Run Deception
      </button>
      <button class="btn btn-secondary btn-sm me-2" @click="resetResult">
        Reset
      </button>
    </div>

    <div v-if="jobResult" class="mt-4 mb-3">
      <h6>Bot Deception Result:</h6>
      <iframe
        ref="attackIframe"
        :srcdoc="jobResult"
        @load="adjustIframeHeight"
        style="width: 100%; border: 1px solid lightgray"></iframe>
    </div>

    <div v-if="pageSource" class="mt-3">
      <h6>
        There’s a hidden link. Malicious bots like web crawler may request the
        link:
      </h6>
      <pre class="code-block"><code v-html="highlightedCode"></code></pre>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Bot Deception</h5>
    </div>
    <div class="card-body">
      <p>
        To prevent bot deception, you can configure the bot deception policy to
        insert link in HTML type response page. For regular clients, the link is
        invisible, while for malicious bots like web crawler, they may request
        the resources which the invisible link points at.
      </p>
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
        })
        .catch((error) => {
          console.error("Error:", error);
        });
    },

    performBotDeception() {
      // Ensure this method matches the button's @click assignment
      fetch("/bot-deception", {
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
