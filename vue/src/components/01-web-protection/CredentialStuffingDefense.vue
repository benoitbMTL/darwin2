<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Credential Stuffing Defense</h5>
      <i class="bi bi-question-circle-fill" style="font-size: 1.5rem;" @click="showHelp = !showHelp"></i>
    </div>
    <div class="card-body">
      <p class="card-text">Authentication with a stolen username / password.</p>

      <div class="d-flex align-items-center mb-3 flex-wrap">
        <select class="form-select form-select-sm me-2" v-model="selectedUser" style="width: 250px">
          <option value="pklangdon4@msn.com">pklangdon4@msn.com</option>
          <option value="muldersstan@gmail.com">muldersstan@gmail.com</option>
          <option value="forsternp2@aol.com">forsternp2@aol.com</option>
          <option value="cragsy@msn.com">cragsy@msn.com</option>
          <option value="bjrehdorf@hotmail.com">bjrehdorf@hotmail.com</option>
          <option value="baz2709@icloud.com">baz2709@icloud.com</option>
          <option value="amysiura@ymail.com">amysiura@ymail.com</option>
          <option value="jond714@gmail.com">jond714@gmail.com</option>
          <option value="josefahorenstein87@hotmail.com">josefahorenstein87@hotmail.com</option>
          <option value="bizotic6@gmail.com">bizotic6@gmail.com</option>
        </select>

        <button class="btn btn-primary btn-sm me-2" @click="performUserAuth">
          Run
        </button>

        <button class="btn btn-secondary btn-sm" @click="resetResult">
          Reset
        </button>
      </div>


      <div v-if="jobResult" class="mt-4 mb-3">
        <h6>Authentication Result:</h6>
        <iframe ref="attackIframe" :srcdoc="jobResult" @load="adjustIframeHeight"
          style="width: 100%; border: 1px solid lightgray;"></iframe>
      </div>

    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Credential Stuffing Defense</h5>
    </div>
    <div class="card-body">
      <p>
      <ul>
        <li>The Credential Stuffing Defense database is designed to safeguard against credential stuffing attacks.</li>
        <li>When activated, FortiWeb verifies the username and password from login requests by comparing them with the
          information in this database.</li>
        <li>This process helps identify whether the username/password pair has been compromised. If a compromise is
          detected, the login attempt is blocked and recorded.</li>
      </ul>
      </p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      selectedUser: "pklangdon4@msn.com",
      showHelp: false,
      jobResult: '',
    };
  },

  methods: {
    performUserAuth() {
      const formData = new URLSearchParams();
      formData.append('username', this.selectedUser);

      fetch("http://localhost:8080//user-auth", {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: formData
      })
        .then(response => {
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }
          return response.text();
        })
        .then(html => {
          this.jobResult = html;
        })
        .catch(error => {
          console.error('Error:', error);
          this.jobResult = 'Failed to perform attack';
        });
    },


    adjustIframeHeight() {
      const iframe = this.$refs.attackIframe;
      if (iframe && iframe.contentWindow && iframe.contentWindow.document.body) {
        iframe.style.height = (iframe.contentWindow.document.body.scrollHeight + 30) + 'px';
      }
    },

    resetResult() {
      this.selectedOption = "All"; // Reset selected option
      this.jobResult = ""; // Clear Result
    },


  }
};
</script>

<style></style>

