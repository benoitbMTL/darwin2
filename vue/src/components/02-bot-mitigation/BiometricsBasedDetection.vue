<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Biometrics-Based Detection</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>

    <div class="card-body">
      <div class="mb-3">
        <label class="form-label">Select Actions:</label>
        <div class="form-check">
          <input class="form-check-input" type="checkbox" value="clickAllProducts" id="clickAllProducts"
            v-model="selectedActions">
          <label class="form-check-label" for="clickAllProducts">Click All Products</label>
        </div>
        <div class="form-check">
          <input class="form-check-input" type="checkbox" value="createAccount" id="createAccount"
            v-model="selectedActions">
          <label class="form-check-label" for="createAccount">Create Account</label>
        </div>
        <div class="form-check">
          <input class="form-check-input" type="checkbox" value="login" id="login" v-model="selectedActions"
            :disabled="!selectedActions.includes('createAccount')">
          <label class="form-check-label" for="login">Login</label>
        </div>
        <div class="form-check">
          <input class="form-check-input" type="checkbox" value="addNewAddress" id="addNewAddress"
            v-model="selectedActions" :disabled="!selectedActions.includes('login')">
          <label class="form-check-label" for="addNewAddress">Add New Address</label>
        </div>
        <div class="form-check">
          <input class="form-check-input" type="checkbox" value="addNewPayment" id="addNewPayment"
            v-model="selectedActions" :disabled="!selectedActions.includes('login')">
          <label class="form-check-label" for="addNewPayment">Add New Payment</label>
        </div>
        <div class="form-check">
          <input class="form-check-input" type="checkbox" value="addItemsToShoppingCart" id="addItemsToShoppingCart"
            v-model="selectedActions" :disabled="!selectedActions.includes('login')">
          <label class="form-check-label" for="addItemsToShoppingCart">Add Items to Shopping Cart</label>
        </div>
        <div class="form-check">
          <input class="form-check-input" type="checkbox" value="checkoutShoppingCart" id="checkoutShoppingCart"
            v-model="selectedActions"
            :disabled="!selectedActions.includes('addNewAddress') || !selectedActions.includes('addNewPayment')">
          <label class="form-check-label" for="checkoutShoppingCart">Checkout Shopping Cart</label>
        </div>
        <div class="form-check">
          <input class="form-check-input" type="checkbox" value="logout" id="logout" v-model="selectedActions"
            :disabled="!selectedActions.includes('login')">
          <label class="form-check-label" for="logout">Logout</label>
        </div>
      </div>

      <div>
        <button class="btn btn-primary me-2" @click="runCustomSelenium">Run Actions</button>
        <button class="btn btn-secondary" @click="resetResult">Reset</button>
      </div>

    </div>

    <div v-if="jobResult" class="mt-4 mb-3">
      <h6>Result:</h6>
      <iframe ref="seleniumIframe" :srcdoc="jobResult" @load="adjustIframeHeight"
        style="width: 100%; border: 1px solid lightgray"></iframe>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Biometrics-Based Detection</h5>
    </div>
    <div class="card-body">
      <p>By checking the client events such as mouse movement, keyboard, screen touch, and scroll, etc in specified
        period, FortiWeb judges whether the request comes from a human or from a bot. You can configure the biometrics
        based detection rule to define the client event, collection period, and the request URL, etc.</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      showHelp: false,
      jobResult: null,
      selectedActions: [],
      loopCount: 1, // Define loopCount with a default value
    };
  },

  methods: {
    runCustomSelenium() {

      const payload = {
        actions: this.selectedActions,
        loopCount: this.loopCount,
      };

      fetch("/selenium", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      })
        .then((response) => response.text())
        .then((html) => {
          this.jobResult = html;
        })
        .catch((error) => {
          console.error("Error:", error);
          this.jobResult = "Failed to run test";
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
      this.selectedActions = []; // Reset selected actions
      this.loopCount = 1; // Reset loop count to default
    },
  },
};
</script>

<style></style>
