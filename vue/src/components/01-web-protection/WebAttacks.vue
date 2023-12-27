<template>
  <div class="card my-4">
    <h5 class="card-header">Web Attacks</h5>
    <div class="card-body">
      <p class="card-text">
        Select a user from the list and generate an attack scenario.
      </p>

      <div class="d-flex align-items-center mb-3">
        <select
          class="form-select form-select-sm me-2"
          v-model="selectedUser"
          style="width: 250px"
        >
          <option value="admin">admin</option>
          <option value="gordonb">gordonb</option>
          <option value="1337">1337</option>
          <option value="pablo">pablo</option>
          <option value="smithy">smithy</option>
        </select>

        <button
          class="btn btn-primary btn-sm me-2"
          @click="performCommandInjection"
        >
          Command Injection
        </button>
        <button
          class="btn btn-primary btn-sm me-2"
          @click="performSQLInjection"
        >
          SQL Injection
        </button>
        <button
          class="btn btn-primary btn-sm me-2"
          @click="performCrossSiteScripting"
        >
          Cross-site Scripting
        </button>
        <button
          class="btn btn-warning btn-sm me-2"
          @click="performZeroDayCommandInjection"
        >
          Zero Day Command Injection
        </button>
        <button
          class="btn btn-warning btn-sm me-2"
          @click="performZeroDayCrossSiteScripting"
        >
          Zero Day Cross-site Scripting
        </button>
        <button class="btn btn-secondary btn-sm" @click="resetGenericAttack">
          Reset
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      selectedUser: "admin",
    };
  },
  methods: {
    performCommandInjection() {
      // Logique de Command Injection...
    },
    // Définissez d'autres méthodes ici...
    resetGenericAttack() {
      // Logique de réinitialisation...
    },
  },
};
</script>

<style></style>
