<template>
  <div class="card my-4">

    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Selenium Playground</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>

    <div class="card-body">
      <label class="form-label">Automated Web Interactions on <a :href="juiceShopDynamicUrl" target="_blank">Juice
          Shop</a> web application with Selenium.</label>
      <div class="d-flex justify-content-start gap-3 mt-3">

        <div> <!-- Actions List -->
          <h6>Actions</h6> <!-- Title for List 1 -->
          <div class="card" style="width: 18rem;">
            <ul class="list-group list-group-flush">
              <li class="list-group-item">
                <div class="form-check">
                  <input class="form-check-input" type="checkbox" value="clickAllProducts" id="clickAllProducts"
                    v-model="selectedActions">
                  <label class="form-check-label" for="clickAllProducts">View All Products</label>
                </div>
              </li>
              <li class="list-group-item">
                <div class="form-check">
                  <input class="form-check-input" type="checkbox" value="createAccount" id="createAccount"
                    v-model="selectedActions">
                  <label class="form-check-label" for="createAccount">Create Account</label>
                </div>
              </li>
              <li class="list-group-item">
                <div class="form-check">
                  <input class="form-check-input" type="checkbox" value="login" id="login" v-model="selectedActions"
                    :disabled="!selectedActions.includes('createAccount')">
                  <label class="form-check-label" for="login">Login</label>
                </div>
              </li>
              <li class="list-group-item">
                <div class="form-check">
                  <input class="form-check-input" type="checkbox" value="addNewAddress" id="addNewAddress"
                    v-model="selectedActions" :disabled="!selectedActions.includes('login')">
                  <label class="form-check-label" for="addNewAddress">Add New Address</label>
                </div>
              </li>
              <li class="list-group-item">
                <div class="form-check">
                  <input class="form-check-input" type="checkbox" value="addNewPayment" id="addNewPayment"
                    v-model="selectedActions" :disabled="!selectedActions.includes('login')">
                  <label class="form-check-label" for="addNewPayment">Add New Payment</label>
                </div>
              </li>
              <li class="list-group-item">
                <div class="form-check">
                  <input class="form-check-input" type="checkbox" value="addItemsToShoppingCart"
                    id="addItemsToShoppingCart" v-model="selectedActions"
                    :disabled="!selectedActions.includes('login')">
                  <label class="form-check-label" for="addItemsToShoppingCart">Add Items to Shopping Cart</label>
                </div>
              </li>
              <li class="list-group-item">
                <div class="form-check">
                  <input class="form-check-input" type="checkbox" value="checkoutShoppingCart" id="checkoutShoppingCart"
                    v-model="selectedActions"
                    :disabled="!selectedActions.includes('addNewAddress') || !selectedActions.includes('addNewPayment')">
                  <label class="form-check-label" for="checkoutShoppingCart">Checkout Shopping Cart</label>
                </div>
              </li>
              <li class="list-group-item">
                <div class="form-check">
                  <input class="form-check-input" type="checkbox" value="logout" id="logout" v-model="selectedActions"
                    :disabled="!selectedActions.includes('login')">
                  <label class="form-check-label" for="logout">Logout</label>
                </div>
              </li>
            </ul>
          </div>
        </div>


        <div> <!-- Options List -->
          <h6>Options</h6> <!-- Title for List 2 -->
          <div class="card" style="width: 18rem;">
            <ul class="list-group list-group-flush">
              <li class="list-group-item d-flex align-items-center">
                Number of loops
                <input type="number" class="form-control form-control-sm ms-2" v-model="loopCount" min="1"
                  style="width: 75px;" placeholder="1">
              </li>
              <li class="list-group-item">
                <div class="form-check">
                  <input class="form-check-input" type="checkbox" id="headless" v-model="isHeadless">
                  <label class="form-check-label" for="headless">Headless</label>
                </div>
              </li>
              <li class="list-group-item d-flex align-items-center">
                Speed
                <select class="form-select form-select-sm ms-2" v-model="selectedSpeed" style="width: 100px;">
                  <option value="1">1</option>
                  <option value="2">2</option>
                  <option value="3">3</option>
                  <option value="4">4</option>
                  <option value="5">5</option>
                  <option value="6">6</option>
                  <option value="7">7</option>
                  <option value="8">8</option>
                  <option value="9">9</option>
                  <option value="10">10</option>
                  <option value="random">Random</option>
                </select>
              </li>
            </ul>
          </div>
        </div>


        <div v-if="results.length > 0"> <!-- Show this portion only if there's at least one result -->
          <h6>Result</h6> <!-- Title for List 3 -->
          <div class="card" style="width: 24rem;">
            <ul class="list-group">
              <li class="list-group-item" v-for="(result, index) in results" :key="index">
                {{ result }}
              </li>
            </ul>
          </div>
        </div>

      </div>

      <div class="mt-3">
        <button class="btn btn-primary btn-sm me-2" @click="runCustomSelenium" :disabled="isRunning">
          <span v-if="isRunning" class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
          <span v-if="!isRunning">Run Actions</span>
          <span v-if="isRunning"> Running...</span>
        </button>
        <button class="btn btn-secondary btn-sm me-2" @click="resetResult">Reset</button>
      </div>

    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Selenium</h5>
    </div>
    <div class="card-body">
      <p>Selenium is a powerful and widely used open-source framework for automating web browsers. It enables testers
        and developers to write scripts in various programming languages like Python, Java, C#, Ruby, and JavaScript to
        control browser actions automatically, simulate user behavior, and interact with web pages. This can include
        tasks like clicking buttons, entering text, navigating through links, and much more, making it a key tool for
        testing web applications to ensure they work as expected across different environments and browsers. Selenium
        supports major browsers like Chrome, Firefox, Safari, and Internet Explorer, and can be run on Windows, Linux,
        and macOS. Its capabilities make it an essential component of quality assurance processes, allowing for both
        functional and regression testing, and it's also used for web scraping and automated task execution.</p>
    </div>
  </div>

</template>

<script>
export default {
  data() {
    return {
      showHelp: false,
      results: [],
      selectedActions: [],
      loopCount: 1, // Default value 
      isHeadless: false, // Default value 
      isRunning: false, // Add this line
      selectedSpeed: "5", // Default value
      config: {
        JUICESHOPURL: "",
      },
    };
  },

  mounted() {
    this.fetchConfig(); // Fetch config on component mount
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


  methods: {
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

    async runCustomSelenium() {
      this.isRunning = true; // Start running
      this.results = []; // Reset results before running
      for (let i = 1; i <= this.loopCount; i++) {
        const payload = {
          actions: this.selectedActions,
          loopCount: this.loopCount,
          headless: this.isHeadless,
          speed: this.selectedSpeed,
        };

        try {
          console.log("Sending payload:", JSON.stringify(payload)); // Print the payload to the console
          const response = await fetch("/selenium", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(payload),
          });

          if (!response.ok) {
            throw new Error('Network response was not ok');
          }
          const data = await response.json();
          this.results.push(`Loop ${i}: ${data}`); // Push the result of each loop with a message
        } catch (error) {
          console.error("Error:", error);
          this.results.push(`Loop ${i}: Failed to run test`);
        }
      }
      this.isRunning = false; // Indicate the process has finished
    },

    adjustIframeHeight() {
      const iframe = this.$refs.seleniumIframe;
      if (iframe && iframe.contentWindow && iframe.contentWindow.document.body) {
        iframe.style.height = iframe.contentWindow.document.body.scrollHeight + 30 + "px";
      }
    },

    resetResult() {
      this.results = [];
      this.selectedActions = [];
      this.loopCount = 1;
      this.isHeadless = false;
      this.selectedSpeed = "5";
    },
  },
};
</script>

<style></style>
