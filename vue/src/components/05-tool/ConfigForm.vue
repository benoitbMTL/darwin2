<template>
  <div class="card my-4">

    <!-- HEADER -->
    <div class="card-header d-flex justify-content-between align-items-center">
      <!-- Static title on the left -->
      <h5>Demo Tool Configuration</h5>

      <!-- Container for the alert message and configuration name -->
      <div class="d-flex align-items-center">
        <!-- Alert Message -->
        <div v-if="showAlert" class="alert alert-success alert-dismissible fade show p-1 mb-0 me-3" role="alert"
          style="font-size: 0.875rem">
          <i class="bi bi-check-circle me-1"></i> {{ alertMessage }}
        </div>

        <!-- Dynamically displayed configuration name on the right -->
        <span v-if="config.NAME" style="color: red;">
          Active Configuration: {{ config.NAME }}
        </span>
      </div>
    </div>

    <!-- BODY -->
    <div class="card-body">
      <div class="container">


        <!-- BUTTONS -->
        <div class="card mb-3">
          <div class="card-body">
            <!-- Use row and cols to align buttons -->
            <div class="row justify-content-between">
              <!-- Left aligned buttons -->
              <div class="col-auto">
                <button type="button" class="btn btn-success btn-sm me-2" @click="applyConfigLocal">
                  <i class="bi bi-arrow-up-square"></i> Apply
                </button>
                <button type="button" class="btn btn-primary btn-sm me-2" @click="promptRenameConfig">
                  <i class="bi bi-pencil"></i> Rename
                </button>
                <button @click="cloneConfigLocal" class="btn btn-primary btn-sm me-2">
                  <i class="bi bi-copy"></i> Clone
                </button>
                <button type="button" class="btn btn-primary btn-sm me-2" @click="triggerFileInput">
                  <i class="bi bi-box-arrow-in-down-right"></i> Import
                </button>
                <input type="file" ref="fileInput" style="display: none" @change="importConfig" />
                <button type="button" class="btn btn-danger btn-sm me-2" @click="deleteConfigLocal">
                  <i class="bi bi-x-square"></i> Delete
                </button>
              </div>

              <!-- Middle aligned buttons -->
              <div class="col-auto">

              </div>

              <!-- Right aligned buttons -->
              <div class="col-auto">
                <button type="button" class="btn btn-secondary btn-sm" @click="resetConfig">
                  <i class="bi bi-arrow-clockwise"></i> Reset to Default
                </button>
              </div>
            </div>
          </div>
        </div>




        <div class="row">

          <!-- COL 1 -->
          <div class="col-md-4">
            <div class="card">
              <div class="card-header">
                <strong>Configuration Profiles</strong>
              </div>

              <div class="card-body">
                <ul class="list-group">
                  <li v-for="(configName, index) in configs" :key="index"
                    class="list-group-item d-flex justify-content-between align-items-center"
                    :class="{ active: selectedConfig === configName }" @click="selectConfig(configName)">
                    {{ configName }}
                    <i v-if="configName === currentConfigName" class="bi bi-arrow-right-circle" style="color: red;"></i>
                  </li>
                </ul>
              </div>
            </div>
          </div>

          <!-- COL 2 -->
          <div class="col-md-8">
            <form @submit.prevent="saveConfig">
              <div class="card">


                <!-- Menu -->
                <div class="card-header d-flex justify-content-between">

                  <!-- Navigation links on the left -->
                  <ul class="nav nav-tabs card-header-tabs" role="button">
                    <li class="nav-conf-item">
                      <a class="nav-link" :class="{ active: activeTab === 'applications' }"
                        @click="activeTab = 'applications'">Applications</a>
                    </li>
                    <li class="nav-conf-item">
                      <a class="nav-link" :class="{ active: activeTab === 'restApi' }"
                        @click="activeTab = 'restApi'">REST API</a>
                    </li>
                    <li class="nav-conf-item">
                      <a class="nav-link" :class="{ active: activeTab === 'misc' }"
                        @click="activeTab = 'misc'">Miscellaneous</a>
                    </li>
                  </ul>

                  <!-- Buttons on the right -->
                  <div>
                    <button @click="saveConfig" type="button" class="btn btn-success btn-sm me-2">
                      <i class="bi bi-floppy"></i> Save
                    </button>
                    <button @click="exportConfig" class="btn btn-primary btn-sm me-2">
                      <i class="bi bi-box-arrow-up-right"></i> Export
                    </button>
                  </div>

                </div>

                <!-- Applications Section -->
                <div class="card-body" v-if="activeTab === 'applications'">
                  <!-- Application Form Fields -->
                  <div class="mb-3">
                    <label for="dvwaUrl" class="form-label">DVWA URL</label>
                    <input type="text" class="form-control" id="dvwaUrl" v-model="config.DVWAURL" />
                  </div>
                  <div class="mb-3">
                    <label for="bankUrl" class="form-label">Bank URL</label>
                    <input type="text" class="form-control" id="bankUrl" v-model="config.BANKURL" />
                  </div>
                  <div class="mb-3">
                    <label for="juiceShopUrl" class="form-label">Juice Shop URL</label>
                    <input type="text" class="form-control" id="juiceShopUrl" v-model="config.JUICESHOPURL" />
                  </div>
                  <div class="mb-3">
                    <label for="petstoreUrl" class="form-label">Petstore URL</label>
                    <input type="text" class="form-control" id="petstoreUrl" v-model="config.PETSTOREURL" />
                  </div>
                  <div class="mb-3">
                    <label for="speedtestUrl" class="form-label">Speedtest URL</label>
                    <input type="text" class="form-control" id="speedtestUrl" v-model="config.SPEEDTESTURL" />
                  </div>
                </div>

                <!-- REST API Section -->
                <div class="card-body" v-if="activeTab === 'restApi'">
                  <!-- REST API Form Fields -->

                  <div class="mb-3">
                    <label for="usernameApi" class="form-label">API Username</label>
                    <input type="text" class="form-control" id="usernameApi" v-model="config.USERNAMEAPI" />
                  </div>

                  <div class="mb-3">
                    <label for="passwordApi" class="form-label">API Password</label>
                    <input type="password" class="form-control" id="passwordApi" v-model="config.PASSWORDAPI" />
                  </div>

                  <div class="mb-3">
                    <label for="vdomApi" class="form-label">VDOM API</label>
                    <input type="text" class="form-control" id="vdomApi" v-model="config.VDOMAPI" />
                  </div>

                  <div class="mb-3">
                    <label for="fwbMgtIp" class="form-label">FortiWeb Management IP/FQDN</label>
                    <input type="text" class="form-control" id="fwbMgtIp" v-model="config.FWBMGTIP" />
                  </div>

                  <div class="mb-3">
                    <label for="fwbMgtPort" class="form-label">FortiWeb Management Port</label>
                    <input type="text" class="form-control" id="fwbMgtPort" v-model="config.FWBMGTPORT" />
                  </div>

                  <div class="mb-3">
                    <label for="mlPolicy" class="form-label">Machine Learning Policy</label>
                    <input type="text" class="form-control" id="mlPolicy" v-model="config.MLPOLICY" />
                  </div>
                </div>

                <!-- Misc Section -->
                <div class="card-body" v-if="activeTab === 'misc'">
                  <!-- Misc Form Fields -->
                  <div class="mb-3">
                    <label for="userAgent" class="form-label">User Agent</label>
                    <input type="text" class="form-control" id="userAgent" v-model="config.USERAGENT" />
                  </div>

                  <div class="mb-3">
                    <label for="fabricLabStory" class="form-label">Fabric Lab Story (Leave empty if the Demo Tool is not
                      running
                      inside the Fabric Lab)</label>
                    <input type="text" class="form-control" id="fabricLabStory" v-model="config.FABRICLABSTORY" />
                  </div>
                </div>
              </div>
            </form>
          </div> <!-- COL -->
        </div> <!-- Row -->
      </div> <!-- Container -->

    </div> <!-- Card Body -->
  </div> <!-- Main Card -->
</template>

<script>
export default {
  data() {
    return {
      currentConfigName: '',
      showBorders: true,
      activeTab: "applications", // Default active tab
      configs: [], // List of saved configuration names
      selectedConfig: null, // Currently selected configuration
      backupName: "", // Name for the new backup

      showAlert: false,
      alertMessage: "",

      config: {
        NAME: "",
        DVWAURL: "",
        BANKURL: "",
        JUICESHOPURL: "",
        PETSTOREURL: "",
        SPEEDTESTURL: "",
        USERNAMEAPI: "",
        PASSWORDAPI: "",
        VDOMAPI: "",
        FWBMGTIP: "",
        FWBMGTPORT: "",
        MLPOLICY: "",
        USERAGENT: "",
        FABRICLABSTORY: "",
      },
    };
  },
  methods: {
    ///////////////////////////////////////////////////////////////////////////////////
    /// FETCH CONFIG, FETCH LIST
    ///////////////////////////////////////////////////////////////////////////////////

    fetchConfig() {
      console.log("Fetching configuration from /config");
      fetch("/config")
        .then((response) => {
          console.log("HTTP return code for /config:", response.status);
          if (!response.ok) {
            throw new Error(`Network response was not ok, status: ${response.status}`);
          }
          return response.json();
        })
        .then((data) => {
          console.log("Configuration successfully fetched:", data);
          this.config = data;
          // Update currentConfigName with the name of the currently active configuration
          this.currentConfigName = data.NAME;
        })
        .catch((error) => {
          console.error("Error fetching updated configuration:", error);
        });
    },


    fetchConfigsList() {
      console.log("Fetching configurations list from /list-configs");
      fetch("/list-configs")
        .then((response) => {
          console.log("HTTP return code for /list-configs:", response.status);
          if (!response.ok) {
            throw new Error(`Network response was not ok, status: ${response.status}`);
          }
          return response.json();
        })
        .then((data) => {
          console.log("Configurations list successfully fetched:", data);
          // Sort the configurations list from 0 to 9 and then from A to Z
          const sortedData = data.sort((a, b) => a.localeCompare(b, 'en', { numeric: true }));
          this.configs = sortedData;
        })
        .catch((error) => {
          console.error("Error fetching configurations list:", error);
        });
    },

    ///////////////////////////////////////////////////////////////////////////////////
    /// RENAME
    ///////////////////////////////////////////////////////////////////////////////////

    promptRenameConfig() {
      if (!this.selectedConfig) {
        alert("Please select a configuration to rename.");
        return;
      }

      const newName = prompt("Enter a new name for the configuration:");
      if (!newName || newName.trim() === "") {
        alert("Renaming aborted. A new name is required.");
        return;
      }

      this.renameConfig(this.selectedConfig, newName.trim());
    },

    renameConfig(oldName, newName) {
      fetch("/rename-config", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ oldName, newName }),
      })
        .then(response => {
          if (!response.ok) {
            return response.json().then(error => Promise.reject(new Error(error.message)));
          }
          return response.json();
        })
        .then(() => {
          this.showAlert = true;
          this.alertMessage = `Configuration '${oldName}' renamed to '${newName}'.`;
          setTimeout(() => {
            this.showAlert = false;
          }, 6000);
          this.fetchConfigsList();
          this.fetchConfig();
        })
        .catch(error => {
          console.error("Rename error:", error);
          this.showAlert = true;
          this.alertMessage = error.message; // Display the parsed error message
          setTimeout(() => {
            this.showAlert = false;
          }, 6000);
        });
    },

    ///////////////////////////////////////////////////////////////////////////////////
    /// SAVE
    ///////////////////////////////////////////////////////////////////////////////////

    saveConfig() {
      console.log("Saving configuration");

      if (this.config.NAME === "Default") {
        // Ask the user for a new name
        const newName = prompt("The 'Default' configuration cannot be overwritten. Please enter a new name for your configuration:");
        if (!newName || newName.trim() === "") {
          alert("Saving aborted. A new name is required.");
          return;
        }
        this.config.NAME = newName.trim();
      }

      fetch("/save-config", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(this.config),
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.json();
        })
        .then((data) => {
          this.showAlert = true;
          this.alertMessage = "Configuration saved successfully.";
          setTimeout(() => {
            this.showAlert = false;
          }, 6000);
          console.log("Success:", data);
          this.fetchConfig();
          this.fetchConfigsList();
        })
        .catch((error) => {
          this.showAlert = true;
          this.alertMessage = "Error saving configuration.";
          setTimeout(() => {
            this.showAlert = false;
          }, 6000);
          console.error("Error:", error);
        });
    },

    ///////////////////////////////////////////////////////////////////////////////////
    /// RESET
    ///////////////////////////////////////////////////////////////////////////////////

    resetConfig() {
      fetch("/reset-config")
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.json();
        })
        .then((data) => {
          this.showAlert = true;
          this.alertMessage = "Configuration reset to default.";
          setTimeout(() => {
            this.showAlert = false;
          }, 6000);
          this.config = data;
          console.log("Configuration saved successfully:", data);
          this.fetchConfig();
          this.fetchConfigsList();
        })
        .catch((error) => {
          this.showAlert = true;
          this.alertMessage = "Error resetting configuration.";
          setTimeout(() => {
            this.showAlert = false;
          }, 6000);
          console.error("Reset error:", error);
        });
    },

    ///////////////////////////////////////////////////////////////////////////////////
    /// EXPORT
    ///////////////////////////////////////////////////////////////////////////////////

    exportConfig() {
      console.log("Exporting configuration");
      fetch("/config")
        .then((response) => {
          if (!response.ok) {
            throw new Error(`Network response was not ok, status: ${response.status}`);
          }
          return response.json();
        })
        .then((config) => {
          // Sanitize the configuration name for use in the filename
          const safeName = config.NAME.replace(/[^a-zA-Z0-9]+/g, "_");
          const filename = `fwb_demo_tool_conf_${safeName}.json`;

          // Create a blob from the configuration JSON and trigger the download
          const blob = new Blob([JSON.stringify(config, null, 2)], { type: "application/json" });
          const url = URL.createObjectURL(blob);
          const a = document.createElement("a");
          a.href = url;
          a.download = filename;
          document.body.appendChild(a);
          a.click();
          window.URL.revokeObjectURL(url);

          this.showAlert = true;
          this.alertMessage = "Configuration exported successfully";
          setTimeout(() => (this.showAlert = false), 6000);
        })
        .catch((error) => {
          console.error("Error during export:", error);
          this.showAlert = true;
          this.alertMessage = "Error during export";
          setTimeout(() => (this.showAlert = false), 6000);
        });
    },

    triggerFileInput() {
      this.$refs.fileInput.click();
    },

    ///////////////////////////////////////////////////////////////////////////////////
    /// IMPORT
    ///////////////////////////////////////////////////////////////////////////////////

    importConfig(e) {
      const file = e.target.files[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
          try {
            const config = JSON.parse(e.target.result);
            fetch("/import", {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify(config),
            })
              .then((response) => {
                if (!response.ok) {
                  throw new Error("Network response was not ok");
                }
                return response.json();
              })
              .then((data) => {
                console.log("Success:", data);
                this.showAlert = true;
                this.alertMessage = "Configuration imported successfully.";
                setTimeout(() => (this.showAlert = false), 6000);
                this.fetchConfig();
                this.fetchConfigsList();
              })
              .catch((error) => {
                console.error("Error during import:", error);
                this.showAlert = true;
                this.alertMessage = "Error importing configuration.";
                setTimeout(() => (this.showAlert = false), 6000);
              });
          } catch (error) {
            console.error("Error parsing file:", error);
            this.showAlert = true;
            this.alertMessage = "Error parsing configuration file.";
            setTimeout(() => (this.showAlert = false), 6000);
          }
        };
        reader.readAsText(file);
      }
    },

    ///////////////////////////////////////////////////////////////////////////////////
    /// SELECT
    ///////////////////////////////////////////////////////////////////////////////////

    selectConfig(configName) {
      this.selectedConfig = configName;
      // Vous pouvez également ajouter ici une logique pour charger les détails
      // de la configuration sélectionnée si nécessaire
      // Par exemple, charger la configuration du serveur et mettre à jour `this.config`
    },

    ///////////////////////////////////////////////////////////////////////////////////
    /// CLONE
    ///////////////////////////////////////////////////////////////////////////////////

    cloneConfigLocal() {
      // Ensure a configuration has been selected for cloning
      if (!this.selectedConfig) {
        alert("Please select a configuration to clone.");
        return;
      }

      // Prompt the user for a new configuration name
      const newName = prompt("Please enter a name for the new configuration:");

      // Check if a name was provided
      if (!newName || newName.trim() === "") {
        alert("Cloning aborted. A new name is required.");
        return;
      }

      // Prepare the request data
      const requestData = {
        sourceName: this.selectedConfig,
        newName: newName.trim(),
      };

      // Proceed to clone the selected configuration under the new name
      fetch("/clone-config", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(requestData),
      })
        .then(response => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.json();
        })
        .then(() => {
          this.showAlert = true;
          this.alertMessage = `Configuration '${requestData.sourceName}' cloned to '${requestData.newName}' and set as current.`;
          setTimeout(() => {
            this.showAlert = false;
          }, 6000);

          // After successful cloning, update the current configuration name
          this.currentConfigName = requestData.newName;
          this.selectedConfig = requestData.newName; // Update the selected configuration to the new clone

          // Refresh the configurations list and fetch the details of the new active configuration
          this.fetchConfig();
          this.fetchConfigsList();
        })
        .catch(error => {
          console.error("Clone error:", error);
          this.showAlert = true;
          this.alertMessage = "Error during cloning.";
          setTimeout(() => {
            this.showAlert = false;
          }, 6000);
        });
    },

    ///////////////////////////////////////////////////////////////////////////////////
    /// APPLY
    ///////////////////////////////////////////////////////////////////////////////////

    applyConfigLocal() {
      // Check if a configuration has been selected
      if (!this.selectedConfig) {
        alert("Please select a configuration to restore.");
        return;
      }

      const data = {
        name: this.selectedConfig,
      };

      fetch("/apply-config", {
        method: "POST", // Use POST method for sending data to the server
        headers: {
          "Content-Type": "application/json", // Indicate that we're sending JSON data
        },
        body: JSON.stringify(data), // Convert the JavaScript object to a JSON string
      })
        .then((response) => {
          if (!response.ok) {
            // If the server responds with a status code that indicates an error,
            // throw an error to be caught in the catch block.
            throw new Error("Failed to restore configuration.");
          }
          return response.json(); // Parse the JSON response body
        })
        .then((data) => {
          // Handle the successful configuration restoration
          console.log("Configuration applied successfully:", data);
          this.showAlert = true;
          this.alertMessage = "Configuration applied successfully.";
          setTimeout(() => { this.showAlert = false; }, 6000);
          this.fetchConfig();
          this.fetchConfigsList();
        })
        .catch((error) => {
          // Handle any errors that occurred during the fetch operation
          console.error("Restore error:", error);
          this.showAlert = true; // Show error alert message
          this.alertMessage = "Error during restoration.";
          setTimeout(() => { this.showAlert = false; }, 6000);
        });
    },

    ///////////////////////////////////////////////////////////////////////////////////
    /// DELETE
    ///////////////////////////////////////////////////////////////////////////////////

    deleteConfigLocal() {
      // Check if a configuration has been selected for deletion
      if (!this.selectedConfig) {
        alert("Please select a configuration to delete.");
        return;
      }

      // Prevent deletion if the selected configuration is "Default"
      if (this.selectedConfig === "Default") {
        alert("The 'Default' configuration cannot be deleted.");
        return;
      }

      // Prepare the data to be sent to the server. The structure of this data
      // might vary depending on your backend requirements. Here, we're assuming
      // the backend needs the name of the configuration to be deleted.
      const data = {
        name: this.selectedConfig,
      };

      // Send a POST request to the "/delete-local" endpoint with the data of the configuration to be deleted.
      fetch("/delete-local", {
        method: "POST", // Use POST method for sending data to the server
        headers: {
          "Content-Type": "application/json", // Indicate that we're sending JSON data
        },
        body: JSON.stringify(data), // Convert the JavaScript object to a JSON string
      })
        .then((response) => {
          if (!response.ok) {
            // If the server responds with a status code that indicates an error,
            // parse the response to get the error message.
            return response.json().then((errorData) => {
              throw new Error(errorData.message);
            });
          }
          return response.json(); // Parse the JSON response body
        })
        .then(() => {
          // Remove the deleted configuration from the 'configs' array
          this.configs = this.configs.filter(
            (config) => config !== this.selectedConfig
          );
          this.selectedConfig = null; // Reset the selected configuration

          // Handle the successful deletion of the configuration
          console.log("Configuration deleted successfully:", this.selectedConfig);
          this.showAlert = true;
          this.alertMessage = "Configuration deleted successfully.";
          setTimeout(() => { this.showAlert = false; }, 6000);
        })
        .catch((error) => {
          // Handle any errors that occurred during the fetch operation
          console.error("Delete error:", error);
          this.showAlert = true; // Show error alert message
          // Use the error message from the server response
          this.alertMessage = error.message;
          setTimeout(() => { this.showAlert = false; }, 6000);
        });
    },
  },


  ///////////////////////////////////////////////////////////////////////////////////
  /// MOUNT
  ///////////////////////////////////////////////////////////////////////////////////

  mounted() {
    console.log("Fetching config");
    this.fetchConfig(); // Load config to the form
    this.fetchConfigsList(); // Load config list
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
