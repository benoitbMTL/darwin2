<template>

  <!-- CARD  -->
  <div class="card my-4">

    <!-- CARD HEADER -->
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>HTTP Playground</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>

    <!-- CARD BODY -->
    <div class="card-body">

      <!-- FORM -->
      <form @submit.prevent="sendRequest">
        <!-- Line 1: Target Selection and Action Buttons -->
        <div class="mb-3 d-flex align-items-center">
          <select class="form-select form-select-sm me-3" v-model="selectedTarget" style="width: auto;">
            <option value="DVWA">DVWA</option>
            <option value="Bank">Bank</option>
            <option value="JuiceShop">Juice Shop</option>
            <option value="Petstore">Petstore</option>
            <option value="Speedtest">Speedtest</option>
          </select>

          <button class="btn btn-primary btn-sm me-2" @click="sendHTTP">Run</button>
          <button class="btn btn-secondary btn-sm" @click="resetResult">Reset</button>
        </div>

        <!-- HTTP Method and URL Form -->
        <div class="mb-3 d-flex justify-content-start align-items-center">

          <select id="httpMethod" class="form-select form-select-sm me-2" v-model="requestConfig.method"
            style="width: auto;">
            <option value="GET">GET</option>
            <option value="POST">POST</option>
            <option value="PUT">PUT</option>
            <option value="DELETE">DELETE</option>
            <option value="PATCH">PATCH</option>
            <option value="OPTIONS">OPTIONS</option>
            <option value="HEAD">HEAD</option>
          </select>

          <input type="url" class="form-control form-control-sm" id="requestUrl" v-model="requestConfig.url"
            placeholder="URL">
        </div>

        <!-- Loop Count -->
        <div class="mb-3 d-flex align-items-baseline">
          <label for="loopCount" class="form-label me-3" style="white-space: nowrap;"><strong>Loop
              Count</strong></label>
          <input type="number" class="form-control form-control-sm" id="loopCount" v-model="requestConfig.loopCount"
            min="1" max="9999" placeholder="1" style="width: 80px;">
        </div>

        <!-- Follow Redirects Checkbox -->
        <div class="mb-3 d-flex align-items-center">
          <label for="followRedirects" class="form-check-label me-2"><strong>Follow Redirects</strong></label>
          <input class="form-check-input" type="checkbox" id="followRedirects" v-model="requestConfig.followRedirects"
            checked>
        </div>

        <!-- Data Content for POST, PUT, etc. -->
        <div v-if="['POST', 'PUT', 'PATCH'].includes(requestConfig.method)" class="mb-3">
          <label for="dataContent" class="form-label"><strong>Data Content</strong></label>
          <textarea class="form-control form-control-sm" id="dataContent"
            v-model="requestConfig.dataContent"></textarea>
        </div>

        <!-- User-Agent -->
        <div class="mb-3 d-flex align-items-baseline">
          <label for="userAgent" class="form-label me-3"
            style="white-space: nowrap;"><strong>User-Agent</strong></label>
          <input type="text" class="form-control form-control-sm" id="userAgent" v-model="requestConfig.userAgent"
            placeholder="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36">
        </div>

        <!-- Content-Type Buttons -->
        <div class="mb-3 d-flex align-items-baseline">
          <label class="form-label me-3"><strong>Content-Type</strong></label>
          <div>
            <button type="button" class="btn btn-outline-dark btn-sm me-1"
              :class="{ 'active': requestConfig.contentType === 'text/plain' }"
              @click="setContentType('text/plain')">text/plain</button>
            <button type="button" class="btn btn-outline-dark btn-sm me-1"
              :class="{ 'active': requestConfig.contentType === 'text/html' }"
              @click="setContentType('text/html')">text/html</button>
            <button type="button" class="btn btn-outline-dark btn-sm me-1"
              :class="{ 'active': requestConfig.contentType === 'application/x-www-form-urlencoded' }"
              @click="setContentType('application/x-www-form-urlencoded')">application/x-www-form-urlencoded</button>
            <button type="button" class="btn btn-outline-dark btn-sm me-1"
              :class="{ 'active': requestConfig.contentType === 'application/json' }"
              @click="setContentType('application/json')">application/json</button>
            <button type="button" class="btn btn-outline-dark btn-sm"
              :class="{ 'active': requestConfig.contentType === 'application/xml' }"
              @click="setContentType('application/xml')">application/xml</button>
          </div>
        </div>

        <!-- Cookie -->
        <div class="mb-3 d-flex align-items-baseline">
          <label for="cookie" class="form-label me-3"><strong>Cookie</strong></label>
          <input type="text" class="form-control form-control-sm" id="cookie" v-model="requestConfig.cookie"
            placeholder="<cookie-name>=<cookie-value>">
        </div>

        <!-- X-Forwarded-For -->
        <div class="mb-3 d-flex align-items-baseline">
          <label for="xForwardedFor" class="form-label me-3"
            style="white-space: nowrap;"><strong>X-Forwarded-For</strong></label>
          <input type="text" class="form-control form-control-sm" id="xForwardedFor"
            v-model="requestConfig.xForwardedFor" placeholder="192.168.1.1">
        </div>

        <!-- Referer -->
        <div class="mb-3 d-flex align-items-baseline">
          <label for="referer" class="form-label me-3"><strong>Referer</strong></label>
          <input type="text" class="form-control form-control-sm" id="referer" v-model="requestConfig.referer"
            placeholder="www.google.com">
        </div>
      </form>
    </div>
  </div>

  <!-- Display Request HTTP headers -->
  <div v-if="requestHeaders" class="card mb-3">
    <div class="card-header">
      <h5><i class="bi bi-box-arrow-right me-2"></i>Request Headers</h5>
    </div>
    <div class="card-body">
      <pre
        style="white-space: pre-wrap; word-break: keep-all; border: 1px solid lightgray; padding: 10px">{{ requestHeaders }}</pre>
    </div>
  </div>

  <!-- Display Response HTTP headers -->
  <div v-if="responseHeaders" class="card mb-3">
    <div class="card-header">
      <h5><i class="bi bi-box-arrow-in-left me-2"></i>Response Headers</h5>
    </div>
    <div class="card-body">
      <pre
        style="white-space: pre-wrap; word-break: keep-all; border: 1px solid lightgray; padding: 10px">{{ responseHeaders }}</pre>
    </div>
  </div>

  <!-- Display TXT RESPONSE BODY -->
  <div v-if="txtResponseBody" class="card mb-3">
    <div class="card-header">
      <h5><i class="bi bi-filetype-txt me-2"></i>Response Body</h5>
    </div>
    <div class="card-body">
      <pre
        style="white-space: pre-wrap; word-break: keep-all; border: 1px solid lightgray; padding: 10px">{{ txtResponseBody }}</pre>
    </div>
  </div>

  <!-- Display HTML RESPONSE BODY -->
  <div v-if="htmlResponseBody" class="card mb-3">
    <div class="card-header">
      <h5><i class="bi bi-filetype-html me-2"></i>Response Body</h5>
    </div>
    <div class="card-body">
      <iframe ref="responseIframe" :srcdoc="htmlResponseBody" @load="adjustIframeHeight"
        style="width: 100%; border: 1px solid lightgray"></iframe>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About HTTP Playground</h5>
    </div>
    <div class="card-body">
    </div>
  </div>

</template>

<script>
export default {
  data() {
    return {
      requestConfig: {
        method: 'GET',
        url: '',
        loopCount: 1,
        followRedirects: true,
        dataContent: '', // Data content for POST, PUT, etc.
        userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
        contentType: 'text/html',
        //cookie: "<cookie-name>=<cookie-value>",
        //xForwardedFor: "192.168.1.1",
        referer: "www.google.com",
      },
      selectedTarget: '',
      showHelp: false,
      config: {},

      requestHeaders: '',
      responseHeaders: '',
      txtResponseBody: '',
      htmlResponseBody: '',

    };
  },


  watch: {
    selectedTarget(newVal, oldVal) {
      console.log(`Target changed from ${oldVal} to ${newVal}`);
      // Directly use newVal to get the corresponding URL from the config object
      this.requestConfig.url = this.config[newVal] || '';
      console.log(`URL updated to: ${this.requestConfig.url}`);
    }
  },

  mounted() {
    // Fetch configuration data when the component is mounted
    this.fetchConfig();
  },

  methods: {
    async fetchConfig() {
      console.log("Fetching configuration from /config...");
      try {
        const response = await fetch('/config');
        if (!response.ok) {
          throw new Error('Failed to fetch configuration');
        }
        const data = await response.json();
        console.log("Configuration fetched successfully:", data);

        // Process the fetched configuration
        this.processConfig(data);
      } catch (error) {
        console.error('Error fetching configuration:', error);
      }
    },

    processConfig(data) {
      // Dynamically create the config object based on fetched data
      this.config = {
        DVWA: data.DVWAURL,
        Bank: data.BANKURL,
        JuiceShop: data.JUICESHOPURL,
        Petstore: data.PETSTOREURL,
        Speedtest: data.SPEEDTESTURL,
      };
      console.log("Configuration processed:", this.config);

      // Set default selected target 
      this.selectedTarget = 'DVWA';

    },

    sendHTTP() {
      // Construct the request data
      const requestData = {
        method: this.requestConfig.method,
        url: this.requestConfig.url,
        loopCount: this.requestConfig.loopCount,
        followRedirects: this.requestConfig.followRedirects,
        dataContent: this.requestConfig.dataContent,
        userAgent: this.requestConfig.userAgent,
        contentType: this.requestConfig.contentType,
        cookie: this.requestConfig.cookie,
        xForwardedFor: this.requestConfig.xForwardedFor,
        referer: this.requestConfig.referer,
      };
      console.log("Sending HTTP request with data:", requestData);

      fetch('/api/http-request', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
      })
        .then(response => {
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }
          return response.json();
        })
        .then(data => {
          this.requestHeaders = data.requestHeaders.trim();
          this.responseHeaders = data.responseHeaders.trim();
          this.responseBody = data.responseBody.trim();

          //console.log("Request Headers:", this.requestHeaders);
          //console.log("Response Headers:", this.requestHeaders);
          //console.log("Response Body:", this.responseBody);

          if (typeof this.responseBody === "string") {
            if (this.responseBody.startsWith("<") && this.responseBody.endsWith(">")) {
              console.log("Detected HTML content");
              this.htmlResponseBody = data.responseBody;
              this.txtResponseBody = "";
            } else {
              console.log("Detected text content");
              this.txtResponseBody = data.responseBody;
              this.htmlResponseBody = "";
            }
          } else {
            console.log("Unexpected data type received:", typeof this.responseBody);
          }

        })
        .catch(error => {
          console.error('Error sending HTTP request:', error);
        });
    },

    setContentType(type) {
      this.requestConfig.contentType = type;
    },

    adjustIframeHeight() {
      const iframe = this.$refs.responseIframe;
      if (iframe && iframe.contentWindow && iframe.contentWindow.document.body) {
        iframe.style.height = iframe.contentWindow.document.body.scrollHeight + 30 + "px";
        console.log("Iframe height adjusted.");
      }
    },

    resetResult() {
      this.requestHeaders = "";
      this.responseHeaders = "";
      this.txtResponseBody = "";
      this.htmlResponseBody = "";
      this.requestConfig = {
        method: 'GET',
        url: '',
        followRedirects: true,
        dataContent: '',
        userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
        contentType: 'text/html',
        cookie: '',
        xForwardedFor: '',
        referer: 'www.google.com',
      };
      this.selectedTarget = 'DVWA';
      this.requestConfig.url = this.config['DVWA'];
    },

  },
};
</script>
