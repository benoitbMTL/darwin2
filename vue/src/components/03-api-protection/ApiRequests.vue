<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <div>
        <h5>API Requests</h5>
      </div>

      <div class="d-flex align-items-center">
        <div v-if="showResetMLMessage" class="me-2">
          <div class="alert alert-success alert-dismissible fade show p-1 me-2 mb-0" role="alert"
            style="font-size: 0.875rem">
            <i class="bi bi-check-circle me-1"></i> {{ resetMLMessage }}
          </div>
        </div>
        <div class="me-2">
          <button type="button" class="btn btn-warning btn-sm" @click="resetApiMachineLearning">
            Reset Machine Learning
          </button>
        </div>

        <div>
          <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
          <!-- Bootstrap icon for help -->
        </div>
      </div>
    </div>

    <div class="card-body">
      <div class="d-flex justify-content-between align-items-center mb-3">
        <h6 class="card-text mb-0">Swagger Petstore - OpenAPI 3.0</h6>
        <button class="btn btn-secondary btn-sm" @click="resetResultAndList">Reset</button>
      </div>

      <!-- API GET -->
      <div class="d-flex mb-3">
        <button type="button" class="btn btn-primary btn-sm me-3"
          style="width: 80px; background-color: #64b0fc; border-color: #64b0fc; font-weight: bold"
          @click="performGetRequest()">
          GET
        </button>

        <!-- GET list -->
        <select id="get-pet-list" class="form-select form-select-sm me-3" v-model="selectedGetOption"
          style="width: 300px">
          <option v-for="(option, index) in getList" :value="option.value" :key="index">
            {{ option.text }}
          </option>
        </select>
      </div>

      <!-- API POST -->
      <div class="mb-3 d-flex align-items-center">
        <button type="button" class="btn btn-success btn-sm me-3"
          style="width: 80px; background-color: #4ecc91; border-color: #4ecc91; font-weight: bold"
          @click="performPostRequest()">
          POST
        </button>

        <!-- POST list -->
        <select id="post-pet-list" class="form-select form-select-sm me-3" v-model="selectedPostOption"
          style="width: 300px">
          <option v-for="(option, index) in postList" :value="option.value" :key="index">
            {{ option.text }}
          </option>
        </select>
      </div>

      <!-- API PUT -->
      <div class="mb-3 d-flex align-items-center">
        <button type="button" class="btn btn-success btn-sm me-3"
          style="width: 80px; background-color: #faa03c; border-color: #faa03c; font-weight: bold"
          @click="performPutRequest()">
          PUT
        </button>

        <!-- PUT list -->
        <select id="put-pet-list" class="form-select form-select-sm me-3" v-model="selectedPutOption"
          style="width: 300px">
          <option v-for="(option, index) in putList" :value="option.value" :key="index">
            {{ option.text }}
          </option>
        </select>
      </div>

      <!-- API DELETE -->
      <div class="mb-4 d-flex align-items-center">
        <button type="button" class="btn btn-success btn-sm me-3"
          style="width: 80px; background-color: #f73c43; border-color: #f73c43; font-weight: bold"
          @click="performDeleteRequest()">
          DELETE
        </button>

        <!-- DELETE list -->
        <select id="delete-pet-list" class="form-select form-select-sm me-3" v-model="selectedDeleteOption"
          style="width: 300px">
          <option v-for="(option, index) in deleteList" :value="option.value" :key="index">
            {{ option.text }}
          </option>
        </select>
      </div>

      <div>
        <!-- Display CURL -->
        <div v-if="curlCommand" class="mt-4">
          <h6><i class="bi bi-terminal me-2"></i>Curl</h6>
          <pre
            style="white-space: pre-wrap; word-break: keep-all; border: 1px solid lightgray; padding: 10px">{{ curlCommand }}</pre>
        </div>

        <!-- Display URL -->
        <div v-if="requestURL" class="mt-4">
          <h6><i class="bi bi-link-45deg me-2"></i>Request URL</h6>
          <pre
            style="white-space: pre-wrap; word-break: keep-all; border: 1px solid lightgray; padding: 10px">{{ requestURL }}</pre>
        </div>

        <!-- Display JSON RESPONSE BODY -->
        <div v-if="jsonResponseBody" class="mt-3">
          <h6><i class="bi bi-filetype-json me-2"></i>Response Body</h6>
          <pre
            style="white-space: pre-wrap; word-break: keep-all; border: 1px solid lightgray; padding: 10px">{{ jsonResponseBody }}</pre>
        </div>

        <!-- Display HTML RESPONSE BODY -->
        <div v-if="htmlResponseBody" class="mt-4">
          <h6><i class="bi bi-filetype-html me-2"></i>Request Result</h6>
          <iframe ref="responseIframe" :srcdoc="htmlResponseBody" @load="adjustIframeHeight"
            style="width: 100%; border: 1px solid lightgray"></iframe>
        </div>
      </div>

    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About API Protection</h5>
    </div>
    <div class="card-body">
      <ul>
        <li><strong>FortiWeb Protection</strong>: Analyzes API calls and applies WAF policies to defend against
          malicious traffic and JSON code exploits.</li>
        <li><strong>JSON Validation</strong>: Configures JSON protection to ensure request contents are free from
          potential attacks.</li>
        <li>
          <strong>Machine Learning-based API Protection</strong>: Learns REST API structures from traffic to create
          models that identify and block malicious
          requests.
        </li>
        <li><strong>Attack Detection</strong>: Identifies attacks by comparing incoming API requests against the
          defined
          API data schema model.</li>
      </ul>

      <div class="card">
        <div class="card-header">API Protection CLI settings</div>
        <div class="card-body">
          <pre class="font-monospace">
config waf api-learning-policy
    set action-mlapi alert_deny
    set schema-property maxLength minLength
    set data-format date-time date time email hostname ipv4 ipv6
    set de-duplication-all disable
    set sample-limit-by-ip 0
    set svm-type extended
    set action-anomaly alert_deny
  next
end</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      selectedGetOption: "",
      selectedPostOption: "",
      selectedPutOption: "",
      selectedDeleteOption: "",

      curlCommand: "",
      requestURL: "",
      jsonResponseBody: "",
      htmlResponseBody: "",

      showHelp: false,
      getList: [],
      postList: [],
      putList: [],
      deleteList: [],

      resetMLMessage: "", // To store the response message
      showResetMLMessage: false, // To control the visibility of the response message

    };
  },

  created() {
    /* GET list */
    this.getList = [
      { value: "findByStatus?status=available", text: "Status Available" },
      { value: "findByStatus?status=sold", text: "Status Sold" },
      { value: "findByStatus?status=pending", text: "Status Pending" },
      { value: "findByStatus?", text: "Empty Status" },
      { value: "findByStatus?status=ABCDEFGHIJKL", text: "Very Long Status" },
      { value: "findByStatus?status=A", text: "Very Short Status" },
      {
        value: "findByStatus?status=;cmd.exe",
        text: "Status with Command Injection",
      },
      {
        value: "findByStatus?status=xx& var1=l var2=s;$var1$var2",
        text: "Status with Zero-Day",
      },
      {
        value: "findByStatus?status=sold&status=pending",
        text: "Duplicate Status",
      },
    ];

    // Set the first item of the GET list as the default option
    if (this.getList.length > 0) {
      this.selectedGetOption = this.getList[0].value;
    }

    /* POST list */
    this.postList = [
      {
        value: {
          id: 999,
          name: "FortiPet",
          category: { id: 1, name: "Dogs" },
          photoUrls: ["fortipet.png"],
          tags: [{ id: 0, name: "Cute" }],
          status: "available",
        },
        text: "Add new pet FortiPet [id:999]",
      },
      {
        value: {
          id: 999,
          name: "FortiPet",
          category: { id: 1, name: "Dogs" },
          photoUrls: ["fortipet.png"],
          tags: [{ id: 0, name: "Cute" }],
          status: "/bin/ls",
        },
        text: "New Pet with Command Injection",
      },
      {
        value: {
          id: 999,
          name: "FortiPet",
          category: { id: 1, name: "Dogs" },
          photoUrls: ["fortipet.png"],
          tags: [{ id: 0, name: "Cute" }],
          status: "<script>alert(123)<\\/script>",
        },
        text: "New Pet with Cross-Site-Scripting",
      },
      {
        value: {
          id: 999,
          name: "FortiPet",
          category: { id: 1, name: "Dogs" },
          photoUrls: ["fortipet.png"],
          tags: [{ id: 0, name: "Cute" }],
          status: "xx& var1=l var2=s;$var1$var2",
        },
        text: "New Pet with Zero-Day",
      },
      {
        value: {
          id: 999,
          name: "FortiPet",
          category: { id: 1, name: "Dogs" },
          photoUrls: ["fortipet.png"],
          tags: [{ id: 0, name: "Cute" }],
          status: "X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*",
        },
        text: "New Pet with Malware",
      },
    ];

    // Set the first item of the POST list as the default option
    if (this.postList.length > 0) {
      this.selectedPostOption = this.postList[0].value;
    }

    /* PUT list */
    this.putList = [
      {
        value: {
          id: 999,
          name: "FortiPet",
          category: { id: 1, name: "Dogs" },
          photoUrls: ["fortipet.png"],
          tags: [{ id: 0, name: "Cute" }],
          status: "sold",
        },
        text: "Modify FortiPet [id:999]",
      },
      {
        value: {
          id: 999,
          name: "FortiPet",
          category: { id: 1, name: "Dogs" },
          photoUrls: ["fortipet.png"],
          tags: [{ id: 0, name: "Cute" }],
          status: "ls;;cmd.exe",
        },
        text: "Modify FortiPet with Command Injection",
      },
      {
        value: {
          id: 999,
          name: "FortiPet",
          category: { id: 1, name: "Dogs" },
          photoUrls: ["fortipet.png"],
          tags: [{ id: 0, name: "Cute" }],
          status: "<script>alert(123)<\\/script>",
        },
        text: "Modify FortiPet with Cross-Site-Scripting",
      },
      {
        value: {
          id: 999,
          name: "FortiPet",
          category: { id: 1, name: "Dogs" },
          photoUrls: ["fortipet.png"],
          tags: [{ id: 0, name: "Cute" }],
          status: "xx& var1=l var2=s;$var1$var2",
        },
        text: "Modify FortiPet with Zero-Day",
      },
      {
        value: {
          id: 999,
          name: "FortiPet",
          category: { id: 1, name: "Dogs" },
          photoUrls: ["fortipet.png"],
          tags: [{ id: 0, name: "Cute" }],
          status: "X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*",
        },
        text: "Modify FortiPet with Malware",
      },
    ];

    // Set the first item of the PUT list as the default option
    if (this.putList.length > 0) {
      this.selectedPutOption = this.putList[0].value;
    }

    /* DELETE list */
    this.deleteList = [
      { value: "999", text: "Delete FortiPet [id:999]" },
      { value: "cmd.exe", text: "Delete FortiPet with Command Injection" },
    ];

    // Set the first item of the DELETE list as the default option
    if (this.deleteList.length > 0) {
      this.selectedDeleteOption = this.deleteList[0].value;
    }
  },

  methods: {


    // Reset API Machine Learning
    resetApiMachineLearning() {
      this.resetMLMessage = ""; // Reset message
      this.showResetMLMessage = false; // Hide message initially

      fetch("/reset-api-machine-learning", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.text();
        })
        .then((text) => {
          this.resetMLMessage = text; // Set the response message
          this.showResetMLMessage = true; // Show message

          setTimeout(() => {
            this.showResetMLMessage = false; // Hide message after 15 seconds
          }, 5001);
        })
        .catch((error) => console.error("Error:", error));
    },




    // Method to send a GET request
    performGetRequest() {
      const url = "/api-get";
      const body = JSON.stringify({ option: this.selectedGetOption });

      // Debug: Print the URL and body to the console
      // console.log("Sending GET request to:", url);
      // console.log("Request body:", body);

      this.sendApiRequest(url, body);
    },

    // Method to send a POST request
    performPostRequest() {
      const url = "/api-post";
      const body = JSON.stringify({ option: this.selectedPostOption });

      this.sendApiRequest(url, body);
    },

    // Method to send a PUT request
    performPutRequest() {
      const url = "/api-put";
      const body = JSON.stringify({ option: this.selectedPutOption });

      this.sendApiRequest(url, body);
    },

    // Method to send a DELETE request
    performDeleteRequest() {
      const url = "/api-delete";
      const body = JSON.stringify({ option: this.selectedDeleteOption });

      this.sendApiRequest(url, body);
    },

    // Method to send API request and handle response
    sendApiRequest(url, body) {
      this.resetResult();
      // console.log("fetch url", url);
      // console.log("fetch body", body);

      fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: body,
      })
        .then((response) => {
          // Checking if the response is ok
          if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
          }
          return response.json();
        })

        .then((data) => {
          this.curlCommand = data.curlCommand;
          this.requestURL = data.url;
          this.responseBody = data.data;

          // console.log("Curl Command:", this.curlCommand);
          // console.log("Request URL Command:", this.requestURL);
          // console.log("responseBody:", this.responseBody);

          // Checking the type of the response this.responseBody
          if (typeof this.responseBody === "string") {
            // Check if the string is HTML
            if (this.responseBody.startsWith("<") && this.responseBody.endsWith(">")) {
              this.htmlResponseBody = this.responseBody;
              this.jsonResponseBody = "";
            } else {
              // It's plain text
              this.jsonResponseBody = this.responseBody;
              this.htmlResponseBody = "";
            }
          } else if (typeof this.responseBody === "object") {
            // It's a JSON object
            this.jsonResponseBody = JSON.stringify(this.responseBody, null, 2);
            this.htmlResponseBody = "";
          } else {
            console.log("Unexpected data type received:", typeof this.responseBody);
          }
        })
        .catch((error) => {
          console.error("Error in sendApiRequest:", error);
          this.jsonResponseBody = "";
          this.htmlResponseBody = "";
        });
    },

    adjustIframeHeight() {
      const iframe = this.$refs.responseIframe;
      if (iframe && iframe.contentWindow && iframe.contentWindow.document.body) {
        iframe.style.height = iframe.contentWindow.document.body.scrollHeight + 30 + "px";
      }
    },

    resetResultAndList() {
      this.selectedGetOption = this.getList.length > 0 ? this.getList[0].value : "";
      this.selectedPostOption = this.postList.length > 0 ? this.postList[0].value : "";
      this.selectedPutOption = this.putList.length > 0 ? this.putList[0].value : "";
      this.selectedDeleteOption = this.deleteList.length > 0 ? this.deleteList[0].value : "";
      this.curlCommand = "";
      this.requestURL = "";
      this.jsonResponseBody = "";
      this.htmlResponseBody = "";
    },

    resetResult() {
      this.curlCommand = "";
      this.requestURL = "";
      this.jsonResponseBody = "";
      this.htmlResponseBody = "";
    },
  },
};
</script>
