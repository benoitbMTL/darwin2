<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Biometrics-Based Detection</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">Choose a Bot action.</p>

      <div class="d-flex align-items-center mb-3 flex-wrap">
        <button class="btn btn-primary btn-sm me-2" @click="runSelenimum('chrome')">Buy Products</button>
        <button class="btn btn-primary btn-sm me-2" @click="runScrapWithApi('foobar')">Scrap Products</button>
        <button class="btn btn-secondary btn-sm me-2" @click="resetResult">Reset</button>
      </div>

      <div v-if="jobResult" class="mt-4 mb-3">
        <h6>Result:</h6>
        <iframe ref="seleniumIframe" :srcdoc="jobResult" @load="adjustIframeHeight" style="width: 100%; border: 1px solid lightgray"></iframe>
      </div>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Biometrics-Based Detection</h5>
    </div>
    <div class="card-body">
      <p>By checking the client events such as mouse movement, keyboard, screen touch, and scroll, etc in specified period, FortiWeb judges whether the request comes from a human or from a bot. You can configure the biometrics based detection rule to define the client event, collection period, and the request URL, etc.</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      showHelp: false,
      jobResult: null, 
    };
  },

  methods: {
    runSelenimum(browserName) {
      const url = "localhost:8080/selenium";
      const formData = new URLSearchParams();
      formData.append("name", browserName);

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


        runScrapWithApi(foobar) {
      const url = "localhost:8080/bot-scraper-api";
      const formData = new URLSearchParams();
      formData.append("name", foobar);

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
      const iframe = this.$refs.seleniumIframe;
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
