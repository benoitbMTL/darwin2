<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>REST API</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem" @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>

    <div class="card-body">
      <p>
        This tool provides two primary sets of API tasks for quick onboarding and decommissioning of the application.<br>After completing these tasks, you can
        verify the application's accessibility at <a href="http://speedtest.corp.fabriclab.ca" target="_blank">http://speedtest.corp.fabriclab.ca</a>.
      </p>

      <div class="container">
        <div class="row">
          <div class="col-md-6">
            <div class="card">
              <div class="card-header">
                <button class="btn btn-primary btn-sm me-2" @click="createPolicy" :disabled="createLoading">
                  <span v-if="createLoading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                  <span>{{ createLoading ? "Creating..." : "Create" }}</span>
                </button>
                <button class="btn btn-secondary btn-sm" @click="resetResult">Reset</button>
              </div>

              <ul class="list-group list-group-flush">
                <li v-for="(task, index) in tasks" :key="index" class="list-group-item d-flex justify-content-between">
                  <span>{{ task.description }}</span>
                  <span :class="['badge', 'rounded-pill', task.colorClass]">{{ task.statusText }}</span>
                </li>
              </ul>
            </div>
          </div>
          <div class="col-md-6">
            <div class="card">
              <div class="card-header">
                <button class="btn btn-primary btn-sm me-2" @click="deletePolicy" :disabled="deleteLoading">
                  <span v-if="deleteLoading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                  <span>{{ deleteLoading ? "Deleting..." : "Delete" }}</span>
                </button>
                <button class="btn btn-secondary btn-sm" @click="resetDeleteResult">Reset</button>
              </div>

              <ul class="list-group list-group-flush">
                <li v-for="(deleteTask, index) in deleteTasks" :key="index" class="list-group-item d-flex justify-content-between">
                  <span>{{ deleteTask.description }}</span>
                  <span :class="['badge', 'rounded-pill', deleteTask.colorClass]">{{ deleteTask.statusText }}</span>
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
        FortiWeb provides an API that uses Representational State Transfer (RESTful API) design principles to access and modify the settings of FortiWeb
        applications.
      </p>
      <p>
        You can use the RESTful API to control FortiWeb and seamlessly integrate FortiWeb to other systems. With a secure and programmable management style by
        HTTPS+ authentication, FortiWeb RESTful API provides enough convenience for those who hope to integrate FortiWeb with other configurations. Furthermore,
        FortiWeb RESTful API can help realize comprehensive management on all FortiWeb features.
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

      tasks: [
        { id: "createNewVirtualIP", description: "Create new Virtual IP", status: "incomplete", colorClass: "bg-secondary", statusText: "Incomplete" },
        { id: "createNewServerPool", description: "Create new Server Pool", status: "incomplete", colorClass: "bg-secondary", statusText: "Incomplete" },
        { id: "createNewMemberPool", description: "Create new Member Pool", status: "incomplete", colorClass: "bg-secondary", statusText: "Incomplete" },
        { id: "createNewVirtualServer", description: "Create new Virtual Server", status: "incomplete", colorClass: "bg-secondary", statusText: "Incomplete" },
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
        { id: "cloneInlineProtection", description: "Clone Inline Protection", status: "incomplete", colorClass: "bg-secondary", statusText: "Incomplete" },
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
        { id: "createNewPolicy", description: "Create new Policy", status: "incomplete", colorClass: "bg-secondary", statusText: "Incomplete" },
      ],

      deleteTasks: [
        { id: "deletePolicy", description: "Delete Policy", status: "incomplete", colorClass: "bg-secondary", statusText: "Incomplete" },
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
        { id: "deleteVirtualServer", description: "Delete Virtual Server", status: "incomplete", colorClass: "bg-secondary", statusText: "Incomplete" },
        { id: "deleteServerPool", description: "Delete Server Pool", status: "incomplete", colorClass: "bg-secondary", statusText: "Incomplete" },
        { id: "deleteVirtualIP", description: "Delete Virtual IP", status: "incomplete", colorClass: "bg-secondary", statusText: "Incomplete" },
      ],
    };
  },

  methods: {
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

    createPolicy() {
      this.resetResult();
      this.resetDeleteResult();
      this.createLoading = true;
      console.log("Creating policy...");

      fetch("http://localhost:8080/create-policy", {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
      })
        .then((response) => response.json())
        .then((data) => {
          console.log("Results", data);
          this.createLoading = false;
          this.jobResult = data;
          data.forEach((status) => {
            this.updateTaskStatus(status.taskId, status.status);
          });
        })
        .catch((error) => {
          console.error("Error creating policy:", error);
          this.createLoading = false;
        });
    },

    deletePolicy() {
      this.resetResult();
      this.resetDeleteResult();
      this.deleteLoading = true;
      console.log("Deleting policy...");

      fetch("http://localhost:8080/delete-policy", {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
      })
        .then((response) => response.json())
        .then((data) => {
          console.log("Policy deleted successfully:", data);
          this.deleteLoading = false;
          this.jobResult = data;
          data.forEach((status) => {
            this.updateDeleteTaskStatus(status.taskId, status.status);
          });
        })
        .catch((error) => {
          console.error("Error deleting policy:", error);
          this.deleteLoading = false;
        });
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
