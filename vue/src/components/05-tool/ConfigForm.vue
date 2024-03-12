<template>

  <h5 class="my-4 d-flex justify-content-between align-items-center">
    Demo Tool Configuration
    <span v-if="config.Name" style="color: red;">Configuration Name: {{ config.Name }}</span>
  </h5>

  <form @submit.prevent="saveConfig">
    <div class="card my-4">
      <!-- Menu -->
      <div class="card-header d-flex justify-content-between align-items-center">
        <div class="d-flex align-items-center">
          <ul class="nav nav-tabs card-header-tabs" role="button">
            <li class="nav-conf-item">
              <a class="nav-link" :class="{ active: activeTab === 'applications' }"
                @click="activeTab = 'applications'">Applications</a>
            </li>
            <li class="nav-conf-item">
              <a class="nav-link" :class="{ active: activeTab === 'restApi' }" @click="activeTab = 'restApi'">REST
                API</a>
            </li>
            <li class="nav-conf-item">
              <a class="nav-link" :class="{ active: activeTab === 'misc' }"
                @click="activeTab = 'misc'">Miscellaneous</a>
            </li>
            <li class="nav-conf-item">
              <a class="nav-link" :class="{ active: activeTab === 'backupRestore' }"
                @click="activeTab = 'backupRestore'">Backup & Restore</a>
            </li>
          </ul>
        </div>

        <div class="d-flex align-items-center">
          <!-- Alert Message -->
          <div v-if="showAlertSaveReset" class="alert alert-success alert-dismissible fade show p-1 me-2 mb-0"
            role="alert" style="font-size: 0.875rem">
            <i class="bi bi-check-circle me-1"></i> {{ alertMessageSaveReset }}
          </div>

          <!-- Buttons -->
          <div>
            <button type="submit" class="btn btn-primary btn-sm me-2">
              Save
            </button>
            <button type="button" class="btn btn-secondary btn-sm" @click="resetConfig">
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
          <label for="fabricLabStory" class="form-label">Fabric Lab Story (Leave empty if the Demo Tool is not running
            inside the Fabric Lab)</label>
          <input type="text" class="form-control" id="fabricLabStory" v-model="config.FABRICLABSTORY" />
        </div>
      </div>

      <!-- Backup & Restore Section -->
      <div class="card-body" v-if="activeTab === 'backupRestore'">
        <!-- Backup & Restore Buttons -->

        <div class="card my-4">
          <div class="card-header">
            <h5>Export & Import Configuration</h5>
          </div>

          <div class="card-body d-flex align-items-end">
            <button @click="exportConfig" class="btn btn-primary btn-sm me-2">
              Export
            </button>

            <button type="button" class="btn btn-primary btn-sm me-2" @click="triggerFileInput">
              Import
            </button>
            <input type="file" ref="fileInput" style="display: none" @change="importConfig" />

            <!-- Alert Message -->
            <div v-if="showAlertFileExport" class="alert alert-success alert-dismissible fade show" role="alert"
              style="font-size: 0.875rem; padding: .25rem 1rem; line-height: 1.5;">
              <i class="bi bi-check-circle me-1"></i>
              {{ alertMessageFileExport }}
            </div>

          </div>
        </div>

        <div class="card my-4">
          <div class="card-header">
            <h5>Local Configurations</h5>
          </div>

          <!-- Card Body -->
          <div class="card-body">
            <!-- Line 1: Buttons and Alert -->
            <div class="d-flex align-items-end mb-3">
              <button @click="backupConfigLocal" class="btn btn-primary btn-sm me-2">
                Backup
              </button>
              <button type="button" class="btn btn-success btn-sm me-2" @click="restoreConfigLocal">
                Restore
              </button>
              <button type="button" class="btn btn-danger btn-sm me-2" @click="deleteConfigLocal">
                Delete
              </button>
              <div v-if="showAlertLocalBackup" class="alert alert-success alert-dismissible fade show" role="alert"
                style="font-size: 0.875rem; padding: .25rem 1rem; line-height: 1.5;">
                <i class="bi bi-check-circle me-1"></i>
                {{ alertMessageLocalBackup }}
              </div>
            </div>

            <!-- Line 2: List -->
            <div class="row">
              <div class="col-12 col-md-3">
                <ul class="list-group">
                  <li v-for="(configName, index) in configs" :key="index" class="list-group-item"
                    :class="{ active: selectedConfig === configName }" @click="selectConfig(configName)">
                    {{ configName }}
                  </li>
                </ul>
              </div>
            </div>
          </div>



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
      configs: [], // List of saved configuration names
      selectedConfig: null, // Currently selected configuration
      backupName: "", // Name for the new backup

      showAlertSaveReset: false,
      alertMessageSaveReset: "",
      showAlertExportImport: false,
      alertMessageExportImport: "",
      showAlertLocalBackup: false,
      alertMessageLocalBackup: "",

      config: {
        Name: "",
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


    ////////////////////////////////////////////////////////
    /// SAVE / RESET
    ///////////////////////////////////////////////////////////////////////////////////

    saveConfig() {

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
          this.showAlertSaveReset = true;
          this.alertMessageSaveReset = "Configuration saved successfully.";
          setTimeout(() => {
            this.showAlertSaveReset = false;
          }, 5000);
          console.log("Success:", data);




          this.fetchConfigsList();




        })
        .catch((error) => {
          this.showAlertSaveReset = true;
          this.alertMessageSaveReset = "Error saving configuration.";
          setTimeout(() => {
            this.showAlertSaveReset = false;
          }, 5000);
          console.error("Error:", error);
        });
    },

    resetConfig() {
      fetch("/reset-config")
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.json();
        })
        .then((data) => {
          this.showAlertSaveReset = true;
          this.alertMessageSaveReset = "Configuration reset to default.";
          setTimeout(() => {
            this.showAlertSaveReset = false;
          }, 5000);
          this.config = data;
          console.log("Configuration saved successfully:", data);
          this.fetchConfig();
          this.fetchConfigsList();
        })
        .catch((error) => {
          this.showAlertSaveReset = true;
          this.alertMessageSaveReset = "Error resetting configuration.";
          setTimeout(() => {
            this.showAlertSaveReset = false;
          }, 5000);
          console.error("Reset error:", error);
        });
    },

    ///////////////////////////////////////////////////////////////////////////////////
    /// EXPORT, IMPORT
    ///////////////////////////////////////////////////////////////////////////////////

    exportConfig() {
      fetch("/export", {
        method: "GET",
      })
        .then((response) => response.blob())
        .then((blob) => {
          const url = window.URL.createObjectURL(blob);
          const a = document.createElement("a");
          a.style.display = "none";
          a.href = url;
          a.download = "config_backup.json";
          document.body.appendChild(a);
          a.click();
          window.URL.revokeObjectURL(url);
          this.showAlertFileExport = true;
          this.alertMessageFileExport = "Configuration backed up successfully";
          setTimeout(() => (this.showAlertFileExport = false), 5000);
        })
        .catch((error) => {
          console.error("Error:", error);
          this.showAlertFileExport = true;
          this.alertMessageFileExport = "Error during backup";
          setTimeout(() => (this.showAlertFileExport = false), 5000);
        });
    },

    triggerFileInput() {
      this.$refs.fileInput.click();
    },

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
                this.showAlertFileExport = true;
                this.alertMessageFileExport = "Configuration restored successfully.";
                setTimeout(() => (this.showAlertFileExport = false), 5000);
                this.fetchConfig();
                this.fetchConfigsList();
              })
              .catch((error) => {
                console.error("Error during restore:", error);
                this.showAlertFileExport = true;
                this.alertMessage = "Error restoring configuration.";
                setTimeout(() => (this.showAlertFileExport = false), 5000);
              });
          } catch (error) {
            console.error("Error parsing file:", error);
            this.showAlertFileExport = true;
            this.alertMessageFileExport = "Error parsing configuration file.";
            setTimeout(() => (this.showAlertFileExport = false), 5000);
          }
        };
        reader.readAsText(file);
      }
    },

    ///////////////////////////////////////////////////////////////////////////////////
    /// BACKUP / RESTORE / DELETE
    ///////////////////////////////////////////////////////////////////////////////////

    selectConfig(configName) {
      this.selectedConfig = configName;
      // Vous pouvez également ajouter ici une logique pour charger les détails
      // de la configuration sélectionnée si nécessaire
      // Par exemple, charger la configuration du serveur et mettre à jour `this.config`
    },

    backupConfigLocal() {
      // Check if the backup name is provided
      if (!this.backupName) {
        alert("Please provide a name for the backup.");
        return;
      }

      // Prepare the data to be sent to the server. The structure of this data
      // might vary depending on your backend requirements. Here, we're assuming
      // the backend needs the name of the backup.
      const data = {
        name: this.backupName,
      };

      // Send a POST request to the "/backup-local" endpoint with the backup data.
      // Use the Fetch API for this purpose.
      fetch("/backup-local", {
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
            throw new Error("Failed to backup configuration.");
          }
          return response.json(); // Parse the JSON response body
        })
        .then((data) => {
          // Handle the successful backup operation
          console.log("Backup successful:", data);
          this.configs.push(this.backupName); // Add the new backup name to the list of configs
          this.selectedConfig = this.backupName; // Optionally, select the new backup
          this.backupName = ""; // Reset the backup name input for future backups
          this.showAlertLocalBackup = true; // Show success alert message
          this.alertMessageLocalBackup =
            "Configuration backed up successfully.";
        })
        .catch((error) => {
          // Handle any errors that occurred during the fetch operation
          console.error("Backup error:", error);
          this.showAlertLocalBackup = true; // Show error alert message
          this.alertMessageLocalBackup = "Error during backup.";
        });
    },

    restoreConfigLocal() {
      // Check if a configuration has been selected
      if (!this.selectedConfig) {
        alert("Please select a configuration to restore.");
        return;
      }

      // Prepare the data to be sent to the server. The structure of this data
      // might vary depending on your backend requirements. Here, we're assuming
      // the backend needs the name of the configuration to be restored.
      const data = {
        name: this.selectedConfig,
      };

      // Send a POST request to the "/restore-local" endpoint with the data of the configuration to be restored.
      fetch("/restore-local", {
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
          console.log("Configuration restored successfully:", data);
          this.showAlertLocalBackup = true; // Show success alert message
          this.alertMessageLocalBackup = "Configuration restored successfully.";
          this.fetchConfig();
          this.fetchConfigsList();
        })
        .catch((error) => {
          // Handle any errors that occurred during the fetch operation
          console.error("Restore error:", error);
          this.showAlertLocalBackup = true; // Show error alert message
          this.alertMessageLocalBackup = "Error during restoration.";
        });
    },

    deleteConfigLocal() {
      // Check if a configuration has been selected for deletion
      if (!this.selectedConfig) {
        alert("Please select a configuration to delete.");
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
          // Handle the successful deletion of the configuration
          console.log("Configuration deleted successfully:", this.selectedConfig);
          this.showAlertLocalBackup = true;
          this.alertMessageLocalBackup = "Configuration deleted successfully.";

          // Remove the deleted configuration from the 'configs' array
          this.configs = this.configs.filter(
            (config) => config !== this.selectedConfig
          );
          this.selectedConfig = null; // Reset the selected configuration
        })
        .catch((error) => {
          // Handle any errors that occurred during the fetch operation
          console.error("Delete error:", error);
          this.showAlertLocalBackup = true; // Show error alert message
          // Use the error message from the server response
          this.alertMessageLocalBackup = error.message;
        });
    },
  },

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
