<template>
  <div id="health" class="card my-4">
    <div class="card-header d-flex flex-wrap justify-content-between align-items-center gap-2">
      <div>
        <h5 class="mb-0">Health Check</h5>
        <small v-if="report" class="text-body-secondary">Last check: {{ formatDate(report.checkedAt) }}</small>
      </div>
      <div class="d-flex gap-2">
        <button v-if="isRunning" class="btn btn-outline-danger btn-sm" @click="cancel">Cancel</button>
        <button class="btn btn-primary btn-sm text-nowrap" :disabled="isRunning" @click="run">
          <span v-if="isRunning" class="spinner-border spinner-border-sm me-2" aria-hidden="true"></span>
          {{ report ? "Run again" : "Run health check" }}
        </button>
        <button v-if="report || error" class="btn btn-secondary btn-sm text-nowrap" @click="clear">
          Clear results
        </button>
      </div>
    </div>

    <div class="card-body">
      <p class="text-body-secondary">Check application endpoints and FortiWeb management/API connectivity.</p>

      <div v-if="job && isRunning" class="mb-3">
        <div class="d-flex justify-content-between small mb-1">
          <span>{{ job.message || "Waiting for a worker" }}</span>
          <span>{{ job.current }} / {{ job.total }}</span>
        </div>
        <div class="progress" role="progressbar" :aria-valuenow="progress" aria-valuemin="0" aria-valuemax="100">
          <div class="progress-bar progress-bar-striped progress-bar-animated" :style="{ width: `${progress}%` }">
            {{ progress }}%
          </div>
        </div>
      </div>

      <div v-if="error" class="alert alert-danger py-2">{{ error }}</div>

      <template v-if="report">
        <div class="d-flex flex-wrap gap-2 mb-3">
          <span class="badge bg-success fs-6">{{ report.summary.up }} Up</span>
          <span class="badge bg-danger fs-6">{{ report.summary.down }} Down</span>
          <span v-if="report.summary.notConfigured" class="badge bg-secondary fs-6">
            {{ report.summary.notConfigured }} Not configured
          </span>
          <span class="badge bg-light text-dark border fs-6">{{ report.summary.total }} Total</span>
          <span class="ms-auto small text-body-secondary align-self-center">
            Completed in {{ formatDuration(report.durationMilliseconds) }}
          </span>
        </div>

        <div class="table-responsive">
          <table class="table table-sm align-middle mb-0">
            <thead>
              <tr><th>Resource</th><th>Status</th><th>HTTP</th><th>Duration</th><th>Message</th></tr>
            </thead>
            <tbody>
              <tr v-for="result in report.results" :key="result.id">
                <td><strong>{{ result.name }}</strong><small class="d-block text-body-secondary">{{ result.url || "—" }}</small></td>
                <td><span class="badge" :class="statusClass(result.status)">{{ statusLabel(result.status) }}</span></td>
                <td>{{ result.httpCode || "—" }}</td>
                <td>{{ formatDuration(result.durationMilliseconds) }}</td>
                <td class="health-message">{{ result.message || "—" }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </template>

      <div v-else-if="!isRunning" class="text-center text-body-secondary py-4">
        No health check results yet.
      </div>
    </div>
  </div>
</template>

<script>
import { cancelJob, getJob, isTerminalJob, startJob } from "../../services/jobs";

export default {
  name: "HealthPanel",
  data() {
    return { job: null, report: null, error: "", timer: null };
  },
  computed: {
    isRunning() {
      return this.job && !isTerminalJob(this.job);
    },
    progress() {
      if (!this.job?.total) return 0;
      return Math.min(100, Math.round((this.job.current / this.job.total) * 100));
    },
  },
  beforeUnmount() {
    this.stopPolling();
  },
  methods: {
    async run() {
      this.clear();
      try {
        this.job = await startJob("health-check", {
          headers: { Accept: "application/json" },
        });
        this.timer = window.setInterval(this.poll, 400);
        await this.poll();
      } catch (error) {
        this.error = error.message;
      }
    },
    async poll() {
      if (!this.job?.id) return;
      const jobID = this.job.id;
      try {
        const job = await getJob(jobID);
        if (this.job?.id !== jobID) return;
        this.job = job;
        if (isTerminalJob(job)) {
          this.stopPolling();
          if (job.status === "succeeded" && job.result) {
            this.report = JSON.parse(job.result);
          } else if (job.status === "cancelled") {
            this.error = "Health check cancelled.";
          } else {
            this.error = job.error || "Health check failed.";
          }
        }
      } catch (error) {
        this.stopPolling();
        this.error = error.message;
      }
    },
    async cancel() {
      if (!this.job?.id) return;
      try {
        this.job = await cancelJob(this.job.id);
      } catch (error) {
        this.error = error.message;
      }
    },
    clear() {
      this.stopPolling();
      this.job = null;
      this.report = null;
      this.error = "";
    },
    stopPolling() {
      if (this.timer) window.clearInterval(this.timer);
      this.timer = null;
    },
    statusClass(status) {
      return { up: "bg-success", down: "bg-danger", "not-configured": "bg-secondary" }[status] || "bg-secondary";
    },
    statusLabel(status) {
      return { up: "Up", down: "Down", "not-configured": "Not configured" }[status] || status;
    },
    formatDuration(milliseconds) {
      if (milliseconds < 1000) return `${milliseconds} ms`;
      return `${(milliseconds / 1000).toFixed(1)} s`;
    },
    formatDate(value) {
      return new Intl.DateTimeFormat(undefined, { dateStyle: "short", timeStyle: "medium" }).format(new Date(value));
    },
  },
};
</script>

<style scoped>
.health-message { max-width: 24rem; overflow-wrap: anywhere; }
</style>
