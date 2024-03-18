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


        <form @submit.prevent="runCustomSelenium">
          <div class="mb-3">
            <label for="loopCount" class="form-label">Loop Count</label>
            <input type="number" class="form-control" id="loopCount" v-model.number="loopCount" min="1">
          </div>

          <div class="mb-3">
            <label class="form-label">Select Actions:</label>
            <div class="form-check">
              <input class="form-check-input" type="checkbox" value="clickAllProducts" id="clickAllProducts"
                v-model="selectedActions">
              <label class="form-check-label" for="clickAllProducts">Click All Products</label>
            </div>
            <div class="form-check">
              <input class="form-check-input" type="checkbox" value="createAccountLoginLogout"
                id="createAccountLoginLogout" v-model="selectedActions">
              <label class="form-check-label" for="createAccountLoginLogout">Create Account and Login/Logout</label>
            </div>
            <div class="form-check">
              <input class="form-check-input" type="checkbox" value="createAccountLoginAddressPaymentLogout"
                id="createAccountLoginAddressPaymentLogout" v-model="selectedActions">
              <label class="form-check-label" for="createAccountLoginAddressPaymentLogout">Create Full Account and Login/Logout</label>
            </div>
            <div class="form-check">
              <input class="form-check-input" type="checkbox" value="fullUserExperience" id="fullUserExperience"
                v-model="selectedActions">
              <label class="form-check-label" for="fullUserExperience">Full User Experience (Account, Address, Payment,
                Cart, Checkout)</label>
            </div>
          </div>

          <button type="submit" class="btn btn-primary">Run Test</button>
        </form>


        <button class="btn btn-secondary btn-sm me-2" @click="resetResult">Reset</button>
      </div>

      <div v-if="jobResult" class="mt-4 mb-3">
        <h6>Result:</h6>
        <iframe ref="seleniumIframe" :srcdoc="jobResult" @load="adjustIframeHeight"
          style="width: 100%; border: 1px solid lightgray"></iframe>
      </div>
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
    };
  },

  methods: {
    runCustomSelenium() {
      const formData = new URLSearchParams();
      formData.append("actions", JSON.stringify(this.selectedActions));
      formData.append("loopCount", this.loopCount.toString());

      // Print user's choices to the console
      console.log("User selected actions:", this.selectedActions);
      console.log("User selected loop count:", this.loopCount);

      fetch("/selenium", {
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
    },
  },
};
</script>

<style></style>
