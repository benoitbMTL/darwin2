<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>REST API</h5>
      <i
        class="bi bi-question-circle-fill bs-icon"
        style="font-size: 1.5rem"
        @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>

    <div class="card-body">
    <p>
      This tool provides two sets of API tasks for quick onboarding
      and decommissioning of the <strong>Speedtest</strong> application.<br /><br />When completing these
      tasks, you can verify the application's accessibility at
      <a :href="speedtestDynamicUrl" target="_blank">Speedtest</a>.
    </p>


      <div class="container">
        <div class="row">
          <div class="col-md-6">
            <div class="card">
              <div class="card-header">
                <button
                  class="btn btn-primary btn-sm me-2"
                  @click="createPolicy"
                  :disabled="createLoading">
                  <span
                    v-if="createLoading"
                    class="spinner-border spinner-border-sm me-2"
                    role="status"
                    aria-hidden="true"></span>
                  <span>{{ createLoading ? "Creating..." : "Create" }}</span>
                </button>
                <button class="btn btn-secondary btn-sm" @click="resetResult">
                  Reset
                </button>
              </div>

              <ul class="list-group list-group-flush">
                <li
                  v-for="(task, index) in tasks"
                  :key="index"
                  class="list-group-item d-flex justify-content-between">
                  <span>{{ task.description }}</span>
                  <span :class="['badge', 'rounded-pill', task.colorClass]">{{
                    task.statusText
                  }}</span>
                </li>
              </ul>
            </div>
          </div>
          <div class="col-md-6">
            <div class="card">
              <div class="card-header">
                <button
                  class="btn btn-primary btn-sm me-2"
                  @click="deletePolicy"
                  :disabled="deleteLoading">
                  <span
                    v-if="deleteLoading"
                    class="spinner-border spinner-border-sm me-2"
                    role="status"
                    aria-hidden="true"></span>
                  <span>{{ deleteLoading ? "Deleting..." : "Delete" }}</span>
                </button>
                <button
                  class="btn btn-secondary btn-sm"
                  @click="resetDeleteResult">
                  Reset
                </button>
              </div>

              <ul class="list-group list-group-flush">
                <li
                  v-for="(deleteTask, index) in deleteTasks"
                  :key="index"
                  class="list-group-item d-flex justify-content-between">
                  <span>{{ deleteTask.description }}</span>
                  <span
                    :class="['badge', 'rounded-pill', deleteTask.colorClass]"
                    >{{ deleteTask.statusText }}</span
                  >
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About FortiWeb Rest API</h5>
    </div>
    <div class="card-body">
      <p>
        FortiWeb provides an API that uses Representational State Transfer
        (RESTful API) design principles to access and modify the settings of
        FortiWeb applications.
      </p>
      <p>
        You can use the RESTful API to control FortiWeb and seamlessly integrate
        FortiWeb to other systems. With a secure and programmable management
        style by HTTPS+ authentication, FortiWeb RESTful API provides enough
        convenience for those who hope to integrate FortiWeb with other
        configurations. Furthermore, FortiWeb RESTful API can help realize
        comprehensive management on all FortiWeb features.
      </p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      jobResult: [], // Initialize as an empty array
      createLoading: false,
      deleteLoading: false,
      showHelp: false,
      config: {
        SPEEDTESTURL: "",
      },

      tasks: [
        {
          id: "createNewVirtualIP",
          description: "Create new Virtual IP",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
        {
          id: "createNewServerPool",
          description: "Create new Server Pool",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
        {
          id: "createNewMemberPool",
          description: "Create new Member Pool",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
        {
          id: "createNewVirtualServer",
          description: "Create new Virtual Server",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
        {
          id: "assignVIPToVirtualServer",
          description: "Assign Virtual IP to Virtual Server",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
        {
          id: "cloneSignatureProtection",
          description: "Clone Signature Protection",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
        {
          id: "cloneInlineProtection",
          description: "Clone Inline Protection",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
        {
          id: "createNewXForwardedForRule",
          description: "Create new X-Forwarded-For Rule",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
        {
          id: "configureProtectionProfile",
          description: "Configure Protection Profile",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
        {
          id: "createNewPolicy",
          description: "Create new Policy",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
      ],

      deleteTasks: [
        {
          id: "deletePolicy",
          description: "Delete Policy",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
        {
          id: "deleteInlineProtection",
          description: "Delete Inline Protection Profile",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
        {
          id: "deleteXForwardedForRule",
          description: "Delete X-Forwarded-For Rule",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
        {
          id: "deleteSignatureProtection",
          description: "Delete Signature Protection",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
        {
          id: "deleteVirtualServer",
          description: "Delete Virtual Server",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
        {
          id: "deleteServerPool",
          description: "Delete Server Pool",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
        {
          id: "deleteVirtualIP",
          description: "Delete Virtual IP",
          status: "incomplete",
          colorClass: "bg-secondary",
          statusText: "Incomplete",
        },
      ],
    };
  },

  mounted() {
    this.fetchConfig(); // Fetch config on component mount
  },

computed: {
  speedtestDynamicUrl() {
    if (this.config.FABRICLABSTORY) {
      return `https://speedtest.${this.config.FABRICLABSTORY}.fabriclab.ca`;
    } else {
      return this.config.SPEEDTESTURL;
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

    updateTaskStatus(taskId, status) {
      let task = this.tasks.find((t) => t.id === taskId);
      if (task) {
        task.status = status;
        if (status === "success") {
          task.colorClass = "bg-success";
          task.statusText = "Done";
        } else if (status === "failure") {
          task.colorClass = "bg-danger";
          task.statusText = "Failed";
        }
      }
    },

    updateDeleteTaskStatus(taskId, status) {
      let task = this.deleteTasks.find((t) => t.id === taskId);
      if (task) {
        task.status = status;
        if (status === "success") {
          task.colorClass = "bg-success";
          task.statusText = "Done";
        } else if (status === "failure") {
          task.colorClass = "bg-danger";
          task.statusText = "Failed";
        }
      }
    },

    async createPolicy() {
      this.resetResult();
      this.resetDeleteResult();
      this.createLoading = true;

      const endpoints = [
        { url: "/create-virtual-ip", taskId: "createNewVirtualIP" },
        { url: "/create-server-pool", taskId: "createNewServerPool" },
        { url: "/create-member-pool", taskId: "createNewMemberPool" },
        { url: "/create-virtual-server", taskId: "createNewVirtualServer" },
        {
          url: "/assign-vip-to-virtual-server",
          taskId: "assignVIPToVirtualServer",
        },
        {
          url: "/clone-signature-protection",
          taskId: "cloneSignatureProtection",
        },
        { url: "/clone-inline-protection", taskId: "cloneInlineProtection" },
        {
          url: "/create-x-forwarded-for-rule",
          taskId: "createNewXForwardedForRule",
        },
        {
          url: "/configure-protection-profile",
          taskId: "configureProtectionProfile",
        },
        { url: "/create-policy", taskId: "createNewPolicy" },
      ];

      for (const endpoint of endpoints) {
        try {
          const response = await fetch(`${endpoint.url}`, {
            method: "POST",
            headers: {
              "Content-Type": "application/x-www-form-urlencoded",
            },
          });

          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }

          const result = await response.json();

          // Print the response values
          // console.log(`Response from ${endpoint.url}`);
          // console.log("TaskID:", result.TaskID);
          // console.log("Description:", result.Description);
          // console.log("Status:", result.Status);
          // console.log("Message:", result.Message);

          this.jobResult.push(result);
          this.updateTaskStatus(endpoint.taskId, result.Status);
        } catch (error) {
          console.error(`Error in ${endpoint.taskId}:`, error);
          this.updateTaskStatus(endpoint.taskId, "failure");
        }
      }

      this.createLoading = false;
    },

    async deletePolicy() {
      this.resetResult();
      this.resetDeleteResult();
      this.deleteLoading = true;

      const endpoints = [
        { url: "/delete-policy", taskId: "deletePolicy" },
        { url: "/delete-inline-protection", taskId: "deleteInlineProtection" },
        {
          url: "/delete-x-forwarded-for-rule",
          taskId: "deleteXForwardedForRule",
        },
        {
          url: "/delete-signature-protection",
          taskId: "deleteSignatureProtection",
        },
        { url: "/delete-virtual-server", taskId: "deleteVirtualServer" },
        { url: "/delete-server-pool", taskId: "deleteServerPool" },
        { url: "/delete-virtual-ip", taskId: "deleteVirtualIP" },
      ];

      for (const endpoint of endpoints) {
        try {
          const response = await fetch(`${endpoint.url}`, {
            method: "POST",
            headers: {
              "Content-Type": "application/x-www-form-urlencoded",
            },
          });

          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }

          const result = await response.json();

          // Print the response values
          // console.log(`Response from ${endpoint.url}`);
          // console.log("TaskID:", result.TaskID);
          // console.log("Description:", result.Description);
          // console.log("Status:", result.Status);
          // console.log("Message:", result.Message);

          this.jobResult.push(result);
          this.updateDeleteTaskStatus(endpoint.taskId, result.Status);
        } catch (error) {
          console.error(`Error in ${endpoint.taskId}:`, error);
          this.updateDeleteTaskStatus(endpoint.taskId, "failure");
        }
      }

      this.deleteLoading = false;
    },

    resetResult() {
      this.tasks = this.tasks.map((task) => ({
        ...task,
        status: "incomplete",
        colorClass: "bg-secondary",
        statusText: "Incomplete",
      }));
    },

    resetDeleteResult() {
      this.deleteTasks = this.deleteTasks.map((task) => ({
        ...task,
        status: "incomplete",
        colorClass: "bg-secondary",
        statusText: "Incomplete",
      }));
    },
  },
};
</script>

<style></style>
