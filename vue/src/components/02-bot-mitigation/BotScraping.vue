<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Web scraping</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">
        By initiating a web scraping process, you will extract data from the <a :href="juiceShopDynamicUrl"
          target="_blank">Juice Shop</a> website, providing structured information without the need for manual
        navigation.
      </p>

      <div class="d-flex align-items-center mb-3 flex-wrap">

        <button class="btn btn-primary btn-sm me-2" @click="runScrapWithApi">
          Scrape Juice Shop
        </button>
        <button class="btn btn-secondary btn-sm me-2" @click="resetResult">
          Reset
        </button>
      </div>

      <div v-if="jobResult" class="mt-4 mb-3">
        <h6>Result:</h6>
        <iframe ref="juiceShopIframe" :srcdoc="jobResult" @load="adjustIframeHeight"
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
      <p>You can download the Bot Detection ML dat file here: <a href="/juiceshop-ml-bot.dat" download="juiceshop-ml-bot.dat">juiceshop-ml-bot.dat</a></p>
      <p>
        The following Machine Learning configuration provides an optimized setup for demonstrations.
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
      config: {
        JUICESHOPURL: "",
      },
      configSnippet: `config waf bot-detection-policy
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
end`,
    };
  },

  mounted() {
    this.fetchConfig(); // Fetch config on component mount
    this.highlightCode();
  },

  computed: {
    juiceShopDynamicUrl() {
      if (this.config.FABRICLABSTORY) {
        return `https://juiceshop.${this.config.FABRICLABSTORY}.fabriclab.ca`;
      } else {
        return this.config.JUICESHOPURL;
      }
    }
  },

  watch: {
    sendSampleResult(newVal) {
      if (newVal) {
        this.highlightedCode = hljs.highlightAuto(newVal).value;
      }
    },
  },


  methods: {

    highlightCode() {
      // Use Highlight.js to apply syntax highlighting to the config snippet
      this.highlightedCode = hljs.highlightAuto(this.configSnippet).value;
    },

    fetchConfig() {
      // Fetch config from server
      fetch("/config")
        .then((response) => response.json())
        .then((data) => {
          this.config = data; // Update config with fetched data
          console.log("Config fetched: ", this.config);
        })
        .catch((error) => {
          console.error("Error fetching config:", error);
        });
    },

    runScrapWithApi() {
      const url = "/bot-scraper-api";

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
