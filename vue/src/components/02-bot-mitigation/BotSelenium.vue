<template>
  <div class="card my-4">

    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Selenium Sandbox</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>

    <div class="card-body">
      <label class="form-label">Automated Web Interactions on Juice Shop with Selenium.</label>

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

            </ul>
          </div>
        </div>


        <div> <!-- Result List -->
          <h6>Result</h6> <!-- Title for List 3 -->
          <div class="card" style="width: 18rem;">
            <ul class="list-group">
              <li class="list-group-item" v-for="(result, index) in results" :key="index">
                {{ result }}
              </li>
            </ul>
          </div>
        </div>

      </div>

      <div class="mt-3">
        <button class="btn btn-primary btn-sm me-2" @click="runCustomSelenium">Run Actions</button>
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
    };
  },

  methods: {
    async runCustomSelenium() {
      this.results = []; // Reset results before running
      for (let i = 1; i <= this.loopCount; i++) {
        const payload = {
          actions: this.selectedActions,
          isHeadless: this.isHeadless,
        };

        try {
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
    },
  },
};
</script>

<style></style>
