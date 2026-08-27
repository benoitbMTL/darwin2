<template>
  <div class="mt-3">
    <div v-if="error" class="alert alert-danger py-2">{{ error }}</div>

    <div v-if="job" class="card border-secondary">
      <div class="card-header d-flex justify-content-between align-items-center">
        <div>
          <strong>{{ job.name }}</strong>
          <span class="badge ms-2" :class="statusClass(job.status)">{{ statusLabel(job.status) }}</span>
        </div>
        <button v-if="isActive" class="btn btn-outline-danger btn-sm" :disabled="job.cancellationPending"
          @click="cancel">
          {{ job.cancellationPending ? "Cancelling..." : "Cancel" }}
        </button>
      </div>
      <div class="card-body">
        <div class="d-flex justify-content-between small mb-1">
          <span>{{ job.message || "Waiting for a worker" }}</span>
          <span v-if="job.total > 0">{{ job.current }} / {{ job.total }} {{ job.unit }}</span>
        </div>
        <div class="progress mb-2" role="progressbar" :aria-valuenow="progress" aria-valuemin="0"
          aria-valuemax="100">
          <div class="progress-bar progress-bar-striped" :class="{ 'progress-bar-animated': isActive }"
            :style="{ width: `${progress}%` }">{{ progress }}%</div>
        </div>
        <div class="d-flex flex-wrap gap-3 text-body-secondary small">
          <span><i class="bi bi-stopwatch me-1"></i>{{ elapsed }}</span>
          <span><i class="bi bi-speedometer2 me-1"></i>{{ rate }}</span>
          <span class="font-monospace">{{ job.id }}</span>
        </div>

        <ul v-if="job.steps?.length" class="list-group list-group-flush mt-3">
          <li v-for="step in job.steps" :key="step.id"
            class="list-group-item px-0 d-flex justify-content-between align-items-start">
            <span>{{ step.label }}<small v-if="step.message" class="d-block text-danger">{{ step.message }}</small></span>
            <span class="badge" :class="statusClass(step.status)">{{ statusLabel(step.status) }}</span>
          </li>
        </ul>

        <div v-if="job.partialErrors?.length" class="alert alert-warning py-2 mt-3 mb-0">
          <strong>{{ job.partialErrors.length }} partial error(s)</strong>
          <ul class="mb-0 mt-1">
            <li v-for="(partialError, index) in job.partialErrors" :key="index">
              <span v-if="partialError.item">{{ partialError.item }}: </span>{{ partialError.message }}
            </li>
          </ul>
        </div>

        <div v-if="job.error" class="alert alert-danger py-2 mt-3 mb-0">{{ job.error }}</div>
        <div v-if="job.result" class="mt-3">
          <h6>Result</h6>
          <pre class="code-block job-result"><code>{{ job.result }}</code></pre>
        </div>
      </div>
    </div>

    <details v-if="history.length" class="mt-3">
      <summary class="small text-body-secondary">Recent jobs ({{ history.length }})</summary>
      <div class="table-responsive mt-2">
        <table class="table table-sm align-middle mb-0">
          <thead><tr><th>Job</th><th>Status</th><th>Progress</th><th>Duration</th><th>Started</th></tr></thead>
          <tbody>
            <tr v-for="item in history" :key="item.id" role="button" @click="selectHistory(item)">
              <td>{{ item.name }}</td>
              <td><span class="badge" :class="statusClass(item.status)">{{ statusLabel(item.status) }}</span></td>
              <td>{{ item.current }} / {{ item.total }}</td>
              <td>{{ formatDuration(item.elapsedMilliseconds) }}</td>
              <td>{{ formatDate(item.createdAt) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </details>
  </div>
</template>

<script>
import { cancelJob, getJob, isTerminalJob, listJobs } from "../../services/jobs";

export default {
  name: "JobMonitor",
  props: {
    jobId: { type: String, default: "" },
    jobType: { type: String, default: "" },
    historyLimit: { type: Number, default: 5 },
  },
  emits: ["updated", "finished"],
  data() {
    return { job: null, history: [], timer: null, error: "", finishedJobId: "" };
  },
  computed: {
    isActive() {
      return this.job && !isTerminalJob(this.job);
    },
    progress() {
      if (!this.job?.total) return this.job?.status === "succeeded" ? 100 : 0;
      return Math.min(100, Math.round((this.job.current / this.job.total) * 100));
    },
    elapsed() {
      return this.formatDuration(this.job?.elapsedMilliseconds || 0);
    },
    rate() {
      const value = this.job?.ratePerSecond || 0;
      return `${value.toFixed(value >= 10 ? 1 : 2)} ${this.job?.unit || "items"}/s`;
    },
  },
  watch: {
    jobId: {
      immediate: true,
      handler(id) {
        this.stopPolling();
        this.finishedJobId = "";
        if (id) {
          this.loadJob();
          this.timer = window.setInterval(this.loadJob, 750);
        } else {
          this.job = null;
        }
        this.loadHistory();
      },
    },
  },
  beforeUnmount() {
    this.stopPolling();
  },
  methods: {
    async loadJob() {
      if (!this.jobId) return;
      const requestedID = this.jobId;
      try {
        const job = await getJob(requestedID);
        if (requestedID !== this.jobId) return;
        this.job = job;
        this.error = "";
        this.$emit("updated", this.job);
        if (isTerminalJob(this.job)) {
          this.stopPolling();
          await this.loadHistory();
          if (this.finishedJobId !== this.job.id) {
            this.finishedJobId = this.job.id;
            this.$emit("finished", this.job);
          }
        }
      } catch (error) {
        this.error = error.message;
        this.stopPolling();
      }
    },
    async loadHistory() {
      try {
        this.history = await listJobs(this.jobType, this.historyLimit);
      } catch (error) {
        this.error = error.message;
      }
    },
    async cancel() {
      try {
        this.job = await cancelJob(this.job.id);
      } catch (error) {
        this.error = error.message;
      }
    },
    selectHistory(item) {
      this.job = item;
    },
    clear() {
      this.stopPolling();
      this.job = null;
      this.error = "";
      this.finishedJobId = "";
    },
    stopPolling() {
      if (this.timer) window.clearInterval(this.timer);
      this.timer = null;
    },
    statusClass(status) {
      return {
        queued: "bg-secondary",
        running: "bg-primary",
        succeeded: "bg-success",
        failed: "bg-danger",
        cancelled: "bg-warning text-dark",
      }[status] || "bg-secondary";
    },
    statusLabel(status) {
      return {
        queued: "Queued",
        running: "Running",
        succeeded: "Completed",
        failed: "Failed",
        cancelled: "Cancelled",
      }[status] || status;
    },
    formatDuration(milliseconds) {
      const seconds = Math.max(0, Math.floor(milliseconds / 1000));
      const minutes = Math.floor(seconds / 60);
      const remainder = seconds % 60;
      return minutes ? `${minutes}m ${remainder}s` : `${remainder}s`;
    },
    formatDate(value) {
      return new Intl.DateTimeFormat(undefined, { dateStyle: "short", timeStyle: "medium" }).format(new Date(value));
    },
  },
};
</script>

<style scoped>
.job-result { max-height: 24rem; overflow: auto; white-space: pre-wrap; }
summary, tbody tr { cursor: pointer; }
</style>
