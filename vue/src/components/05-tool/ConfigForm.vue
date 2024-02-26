<template>
  <h5 class="my-4">Demo Tool Configuration</h5>

  <form @submit.prevent="submitForm">
    <div class="card my-4">
      <!-- Menu -->
      <div
        class="card-header d-flex justify-content-between align-items-center">
        <div class="d-flex align-items-center">
          <ul class="nav nav-tabs card-header-tabs" role="button">
            <li class="nav-conf-item">
              <a
                class="nav-link"
                :class="{ active: activeTab === 'applications' }"
                @click="activeTab = 'applications'"
                >Applications</a
              >
            </li>
            <li class="nav-conf-item">
              <a
                class="nav-link"
                :class="{ active: activeTab === 'restApi' }"
                @click="activeTab = 'restApi'"
                >REST API</a
              >
            </li>
            <li class="nav-conf-item">
              <a
                class="nav-link"
                :class="{ active: activeTab === 'misc' }"
                @click="activeTab = 'misc'"
                >Miscellaneous</a
              >
            </li>
          </ul>
        </div>

        <div class="d-flex align-items-center">
          <!-- Alert Message -->
          <div
            v-if="showAlert"
            class="alert alert-success alert-dismissible fade show p-1 me-2 mb-0"
            role="alert"
            style="font-size: 0.875rem">
            <i class="bi bi-check-circle me-1"></i> {{ alertMessage }}
          </div>

          <!-- Buttons -->
          <div>
            <button type="submit" class="btn btn-primary btn-sm me-2">
              Save Configuration
            </button>
            <button
              type="button"
              class="btn btn-secondary btn-sm"
              @click="resetConfig">
              Reset to Default
            </button>
          </div>
        </div>
      </div>

      <!-- Applications Section -->
      <div class="card-body" v-if="activeTab === 'applications'">
        <!-- Application Form Fields -->
        <div class="mb-3">
          <label for="dvwaUrl" class="form-label">DVWA URL</label>
          <input
            type="text"
            class="form-control"
            id="dvwaUrl"
            v-model="config.DVWAURL" />
        </div>
        <div class="mb-3">
          <label for="dvwaHost" class="form-label">DVWA Host</label>
          <input
            type="text"
            class="form-control"
            id="dvwaHost"
            v-model="config.DVWAHOST" />
        </div>
        <div class="mb-3">
          <label for="juiceShopUrl" class="form-label">Juice Shop URL</label>
          <input
            type="text"
            class="form-control"
            id="juiceShopUrl"
            v-model="config.JUICESHOPURL" />
        </div>
        <div class="mb-3">
          <label for="bankUrl" class="form-label">Bank URL</label>
          <input
            type="text"
            class="form-control"
            id="bankUrl"
            v-model="config.BANKURL" />
        </div>
        <div class="mb-3">
          <label for="speedtestUrl" class="form-label">Speedtest URL</label>
          <input
            type="text"
            class="form-control"
            id="speedtestUrl"
            v-model="config.SPEEDTESTURL" />
        </div>
        <div class="mb-3">
          <label for="petstoreUrl" class="form-label">Petstore URL</label>
          <input
            type="text"
            class="form-control"
            id="petstoreUrl"
            v-model="config.PETSTOREURL" />
        </div>
      </div>

      <!-- REST API Section -->
      <div class="card-body" v-if="activeTab === 'restApi'">
        <!-- REST API Form Fields -->
        <div class="mb-3">
          <label for="usernameApi" class="form-label">Username API</label>
          <input
            type="text"
            class="form-control"
            id="usernameApi"
            v-model="config.USERNAMEAPI" />
        </div>
        <div class="mb-3">
          <label for="passwordApi" class="form-label">Password API</label>
          <input
            type="text"
            class="form-control"
            id="passwordApi"
            v-model="config.PASSWORDAPI" />
        </div>
        <div class="mb-3">
          <label for="vdomApi" class="form-label">VDOM API</label>
          <input
            type="text"
            class="form-control"
            id="vdomApi"
            v-model="config.VDOMAPI" />
        </div>

        <div class="mb-3">
          <label for="fwbMgtIp" class="form-label"
            >FortiWeb Management IP</label
          >
          <input
            type="text"
            class="form-control"
            id="fwbMgtIp"
            v-model="config.FWBMGTIP" />
        </div>

        <div class="mb-3">
          <label for="fwbMgtPort" class="form-label"
            >FortiWeb Management PORT</label
          >
          <input
            type="text"
            class="form-control"
            id="fwbMgtPort"
            v-model="config.FWBMGTPORT" />
        </div>

        <div class="mb-3">
          <label for="mlPolicy" class="form-label"
            >Machine Learning Policy</label
          >
          <input
            type="text"
            class="form-control"
            id="mlPolicy"
            v-model="config.MLPOLICY" />
        </div>
      </div>

      <!-- Misc Section -->
      <div class="card-body" v-if="activeTab === 'misc'">
        <!-- Misc Form Fields -->
        <div class="mb-3">
          <label for="userAgent" class="form-label">User Agent</label>
          <input
            type="text"
            class="form-control"
            id="userAgent"
            v-model="config.USERAGENT" />
        </div>
      </div>
    </div>
  </form>
</template>

<script>
export default {
  data() {
    return {
      activeTab: "applications", // Default active tab
      showAlert: false,
      alertMessage: "",
      config: {
        DVWAURL: "",
        DVWAHOST: "",
        JUICESHOPURL: "",
        BANKURL: "",
        SPEEDTESTURL: "",
        PETSTOREURL: "",
        USERNAMEAPI: "",
        PASSWORDAPI: "",
        VDOMAPI: "",
        FWBMGTIP: "",
        FWBMGTPORT: "",
        MLPOLICY: "",
        USERAGENT: "",
      },
    };
  },
  methods: {
    submitForm() {
      // Implement API call to update configuration
      fetch(`${process.env.VUE_APP_BACKEND_URL}/config`, { 
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(this.config),
      })
        .then((response) => response.json())
        .then((data) => {
          this.showAlert = true;
          this.alertMessage = "Configuration saved successfully.";
          // Reset showAlert after some time if needed
          setTimeout(() => (this.showAlert = false), 5000);
          console.log("Success:", data);
        })
        .catch((error) => {
          this.showAlert = true;
          this.alertMessage = "Error saving configuration.";
          // Reset showAlert after some time if needed
          setTimeout(() => (this.showAlert = false), 5000);
          console.error("Error:", error);
        });
    },
    resetConfig() {
      fetch(`${process.env.VUE_APP_BACKEND_URL}/reset`)
        .then((response) => response.json())
        .then((data) => {
          this.showAlert = true;
          this.alertMessage = "Configuration reset to default.";
          // Reset showAlert after some time if needed
          setTimeout(() => (this.showAlert = false), 5000);
          this.config = data;
          console.log("Success:", data);
        })
        .catch((error) => {
          this.showAlert = true;
          this.alertMessage = "Error resetting configuration.";
          // Reset showAlert after some time if needed
          setTimeout(() => (this.showAlert = false), 5000);
          console.error("Error:", error);
        });
    },
  },
  mounted() {
    // Fetch current configuration from the Go backend
console.log(`Making a GET request to ${process.env.VUE_APP_BACKEND_URL}/config`);
    fetch(`${process.env.VUE_APP_BACKEND_URL}/config`) 
      .then((response) => {
        console.log("HTTP return code:", response.status); // Print HTTP return code

        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        return response.json();
      })
      .then((data) => {
        this.config = data; // Update config with fetched data
        console.log("Received config values:", data); // Print the config values
      })
      .catch((error) => {
        console.error("Fetch error:", error);
      });
  },
};
</script>

<style>
.nav-conf-item a {
  color: #000;
  /* Set your desired color, here it's black */
  text-decoration: none;
  /* Removes the underline */
}
</style>
