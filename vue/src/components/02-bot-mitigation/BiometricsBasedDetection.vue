<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Biometrics-Based Detection</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">Selenium.</p>

      <div class="d-flex align-items-center mb-3 flex-wrap">
        <button class="btn btn-primary btn-sm me-2" @click="runSelenimum('chrome')">Chrome/Selenium</button>
        <button class="btn btn-primary btn-sm me-2" @click="runSelenimum('firefox')">Firefox/Selenium</button>
        <button class="btn btn-primary btn-sm me-2" @click="runColly('foobar')">Colly</button>
        <button class="btn btn-primary btn-sm me-2" @click="runScrapWithApi('foobar')">Scrap with API</button>
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
      <h5>About Selenium</h5>
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
      showHelp: false,
      jobResult: null, 
    };
  },

  methods: {
    runSelenimum(browserName) {
      const url = "http://localhost:8080/selenium";
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


        runColly(foobar) {
      const url = "http://localhost:8080/colly";
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


        runScrapWithApi(foobar) {
      const url = "http://localhost:8080/bot-scraper-api";
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
