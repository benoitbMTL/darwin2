<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Bot Deception</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">Bot Deception</p>


    <p>Inserts invisible links in HTML responses to distinguish between regular clients and malicious bots like web crawlers.</p>

    <h5>Testing Bot Deception</h5>
    <p>Browse to <a href="https://dvwa.canadaeast.cloudapp.azure.com/" target="_blank" rel="noopener noreferrer">https://dvwa.canadaeast.cloudapp.azure.com/</a>, right click and select “View Page Source”.</p>

    <p>There’s a hidden link. Malicious bots like web crawler may request the link…</p>

    <p>Click on it or Browse to <a href="https://dvwa.canadaeast.cloudapp.azure.com/fake_url.php" target="_blank" rel="noopener noreferrer">https://dvwa.canadaeast.cloudapp.azure.com/fake_url.php</a></p>



    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Bot Deception</h5>
    </div>
    <div class="card-body">
      <p>To prevent bot deception, you can configure the bot deception policy to insert link in HTML type response page. For regular clients, the link is invisible, while for malicious bots like web crawler, they may request the resources which the invisible link points at.</p>
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
      const url = "/selenium";
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
