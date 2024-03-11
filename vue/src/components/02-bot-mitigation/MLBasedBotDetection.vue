<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Machine Learning Based Bot Detection</h5>
      <i
        class="bi bi-question-circle-fill bs-icon"
        style="font-size: 1.5rem"
        @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">
        Web scraping is a technique that uses a program to extract data from a
        site.
      </p>

      <div class="d-flex align-items-center mb-3 flex-wrap">
        <button class="btn btn-primary btn-sm me-2" @click="runScrapWithApi">
          Web scraping
        </button>
        <button class="btn btn-secondary btn-sm me-2" @click="resetResult">
          Reset
        </button>
      </div>

      <div v-if="jobResult" class="mt-4 mb-3">
        <h6>Result:</h6>
        <iframe
          ref="juiceShopIframe"
          :srcdoc="jobResult"
          @load="adjustIframeHeight"
          style="width: 100%; border: 1px solid lightgray"></iframe>
      </div>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Machine Learning Based Bot Detection</h5>
    </div>
    <div class="card-body">
      <p>
        The primary goal of web scraping is to gather content from websites in a
        way that preserves the original structure of the data. Therefore, this
        technique is often used for competitive analysis, especially on
        e-commerce websites.
      </p>
      <p>
        The machine learning-based bot detection model enhances existing
        detection methods by identifying advanced bots that may otherwise remain
        unnoticed. It analyzes user behavior across thirteen dimensions, such as
        the frequency of HTTP requests, the use of irregular HTTP versions, and
        whether the requests call for JSON/XML resources.
      </p>
      <p>
        This AI-driven approach eliminates the need for setting and testing
        various thresholds to pinpoint abnormal behavior, which can be a
        cumbersome process with traditional methods. For instance, determining
        the number of user-initiated HTTP requests that should be deemed
        suspicious often requires extensive trial and error along with constant
        monitoring of attack logs.
      </p>
      <p>
        Utilizing the bot detection model simplifies the process. FortiWeb
        employs a Support Vector Machine (SVM) algorithm to create a model that
        learns the normal traffic patterns of users. Incoming traffic is then
        assessed against these patterns; a mismatch leads to the classification
        of the traffic as anomalous. The model also self-adjusts to significant
        changes in user behavior resulting from updates to your application.
      </p>
      <p>
        Tests indicate that the bot detection model excels in identifying
        crawlers and scrapers, evaluating traffic from multiple perspectives to
        improve accuracy and reduce false positives.
      </p>

      <p>
        The following Machine Learning configuration provides an optimized setup
        for demonstrations.
      </p>

      <pre class="code-block"><code v-html="highlightedCode"></code></pre>
    </div>
  </div>
</template>

<script>
import hljs from "highlight.js";
import "highlight.js/styles/monokai.css"; // Monokai theme for Highlight.js

export default {
  data() {
    return {
      highlightedCode: "",
      showHelp: false,
      jobResult: null,
      configSnippet: `
config waf bot-detection-policy
  edit 1
    set sampling-count 10
    set sampling-count-per-client 1
    set sampling-time-per-vector 1
    set training-accuracy 20.0
    set cross-validation 20.0
    set testing-accuracy 20.0
    set anomaly-count 1
    set bot-confirmation disable
    set space-clustering disable
    set clustering-normalization disable
next
end
          `,
    };
  },

  mounted() {
    this.highlightCode(); // Call this method to apply syntax highlighting
  },

  methods: {
    fetchConfig() {
      return fetch("localhost:8080/config", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.json();
        })
        .catch((error) => {
          console.error("Failed to fetch config:", error);
          throw error; // Propagate error to caller
        });
    },

    highlightCode() {
      // Use Highlight.js to apply syntax highlighting to the config snippet
      this.highlightedCode = hljs.highlightAuto(this.configSnippet).value;
    },

    juiceShop() {
      this.fetchConfig()
        .then((config) => {
          const url = config.JUICESHOPURL;

          fetch(url, {
            method: "GET", // Assuming the juice shop allows GET requests to fetch content
          })
            .then((response) => response.text())
            .then((html) => {
              this.jobResult = html; // Set the HTML content to jobResult
              this.adjustIframeHeight(); // Adjust the iframe height after setting the content
            })
            .catch((error) => {
              console.error("Error:", error);
              this.jobResult = "Failed to retrieve content from Juice Shop";
            });
        })
        .catch((error) => {
          console.error("Config fetch error:", error);
          this.jobResult = "Failed to load configuration for Juice Shop URL";
        });
    },

    runScrapWithApi() {
      const url = "localhost:8080/bot-scraper-api";

      fetch(url, {
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
          this.jobResult = "Failed to send Bot";
        });
    },

    adjustIframeHeight() {
      const iframe = this.$refs.juiceShopIframe;
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
      this.jobResult = ""; // Clear Result
    },
  },
};
</script>

<style></style>
