<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Traffic Generation</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem;" @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>


    <div class="card-body">
      <p class="card-text">Select your target and generate various web attacks from random public IP addresses.</p>

      <div class="d-flex align-items-center mb-3">
        <select class="form-select form-select-sm me-2" v-model="selectedTarget" style="width: auto;">
          <option value="DVWA">dvwa</option>
          <option value="Bank">Bank</option>
          <option value="JuiceShop">Juice Shop</option>
          <option value="Petstore">Petstore</option>
          <option value="Speedtest">Speedtest</option>
        </select>

        <!-- Generate Traffic Button -->
        <button class="btn btn-primary btn-sm" @click="generateTraffic" :disabled="isLoading">
          <span v-if="isLoading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
          <span>{{ isLoading ? 'Generating...' : 'Generate Traffic' }}</span>
        </button>

        <!-- Clear display without cancelling the background job -->
        <button class="btn btn-secondary btn-sm ms-2" @click="resetResult">Clear display</button>
      </div>

      <div v-if="jobResult" class="alert alert-danger py-2">{{ jobResult }}</div>
      <JobMonitor ref="jobMonitor" :job-id="activeJobId" job-type="traffic-generation"
        @finished="isLoading = false" />
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Traffic Generation</h5>
    </div>
    <div class="card-body">
      <ul>
        <li>The Traffic Generator simulates cyber attacks using randomly generated public IP addresses.</li>
        <li>Nikto is utilized by the Traffic Generator to generate random attacks.</li>
        <li>Launch the Traffic Generator before a demonstration to populate FortiWeb logs and FortiView dashboards.</li>
      </ul>

    </div>
  </div>
</template>

<script>
import hljs from "../../utils/highlight";
import JobMonitor from "../jobs/JobMonitor.vue";
import { startJob } from "../../services/jobs";

export default {
  components: { JobMonitor },
  data() {
    return {
      isLoading: false,
      jobResult: '',
      activeJobId: '',
      showHelp: false,
      highlightedCode: "",
      selectedTarget: 'DVWA', // Default selection
    };
  },

  watch: {
    jobResult(newVal) {
      if (newVal) {
        this.highlightedCode = hljs.highlightAuto(newVal).value;
      }
    },
  },
  methods: {
    async generateTraffic() {
      console.log('Starting traffic generation...');
      console.log('Selected target:', this.selectedTarget); // Debug log
      this.isLoading = true;
      this.jobResult = ''; // Reset job result

      try {
        const job = await startJob("traffic-generation", {
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ target: this.selectedTarget }),
        });
        this.activeJobId = job.id;
      } catch (error) {
        this.jobResult = `Error: ${error.message}`;
        this.isLoading = false;
      }
    },
    resetResult() {
      console.log('Resetting Result');
      this.jobResult = '';
      this.selectedTarget = "DVWA"; // Reset selected option
      this.activeJobId = '';
      this.isLoading = false;
      this.$refs.jobMonitor?.clear();
    },
  },
};
</script>

<style></style>
