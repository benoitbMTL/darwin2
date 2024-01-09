<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Machine Learning</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">Machine Learning Options.</p>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Machine Learning</h5>
    </div>
    <div class="card-body">
      <p>Help content goes here...</p>
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
      const url = "http://localhost:8080/selenium";
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
