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
                  Clear display
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
            <div v-if="createError" class="alert alert-danger py-2 mt-2">{{ createError }}</div>
            <JobMonitor ref="createJobMonitor" :job-id="activeCreateJobId" job-type="fortiweb-create"
              @updated="syncCreateJob" @finished="createLoading = false" />
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
                  Clear display
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
            <div v-if="deleteError" class="alert alert-danger py-2 mt-2">{{ deleteError }}</div>
            <JobMonitor ref="deleteJobMonitor" :job-id="activeDeleteJobId" job-type="fortiweb-delete"
              @updated="syncDeleteJob" @finished="deleteLoading = false" />
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
import JobMonitor from "../jobs/JobMonitor.vue";
import { startJob } from "../../services/jobs";

export default {
  components: { JobMonitor },
  data() {
    return {
      jobResult: [], // Initialize as an empty array
      createLoading: false,
      deleteLoading: false,
      activeCreateJobId: "",
      activeDeleteJobId: "",
      createError: "",
      deleteError: "",
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
      if (task) this.applyStatus(task, status);
    },

    updateDeleteTaskStatus(taskId, status) {
      let task = this.deleteTasks.find((t) => t.id === taskId);
      if (task) this.applyStatus(task, status);
    },

    applyStatus(task, status) {
      const normalized = { success: "succeeded", failure: "failed" }[status] || status;
      task.status = normalized;
      const display = {
        queued: ["bg-secondary", "Queued"],
        running: ["bg-primary", "Running"],
        succeeded: ["bg-success", "Done"],
        failed: ["bg-danger", "Failed"],
        cancelled: ["bg-warning text-dark", "Cancelled"],
      }[normalized] || ["bg-secondary", "Incomplete"];
      [task.colorClass, task.statusText] = display;
    },

    syncCreateJob(job) {
      for (const step of job.steps || []) this.updateTaskStatus(step.id, step.status);
    },

    syncDeleteJob(job) {
      for (const step of job.steps || []) this.updateDeleteTaskStatus(step.id, step.status);
    },

    async createPolicy() {
      this.resetResult();
      this.resetDeleteResult();
      this.createLoading = true;
      this.createError = "";

      try {
        const job = await startJob("fortiweb-create", {
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
        });
        this.activeCreateJobId = job.id;
      } catch (error) {
        this.createLoading = false;
        this.createError = error.message;
        console.error("Unable to start FortiWeb creation job:", error);
      }
    },

    async deletePolicy() {
      this.resetResult();
      this.resetDeleteResult();
      this.deleteLoading = true;
      this.deleteError = "";

      try {
        const job = await startJob("fortiweb-delete", {
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
        });
        this.activeDeleteJobId = job.id;
      } catch (error) {
        this.deleteLoading = false;
        this.deleteError = error.message;
        console.error("Unable to start FortiWeb deletion job:", error);
      }
    },

    resetResult() {
      this.tasks = this.tasks.map((task) => ({
        ...task,
        status: "incomplete",
        colorClass: "bg-secondary",
        statusText: "Incomplete",
      }));
      this.activeCreateJobId = "";
      this.createLoading = false;
      this.createError = "";
      this.$refs.createJobMonitor?.clear();
    },

    resetDeleteResult() {
      this.deleteTasks = this.deleteTasks.map((task) => ({
        ...task,
        status: "incomplete",
        colorClass: "bg-secondary",
        statusText: "Incomplete",
      }));
      this.activeDeleteJobId = "";
      this.deleteLoading = false;
      this.deleteError = "";
      this.$refs.deleteJobMonitor?.clear();
    },
  },
};
</script>

<style></style>
