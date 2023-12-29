<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Web Attacks</h5>
      <i class="bi bi-question-circle-fill bs-icon" @click="showHelp = !showHelp"></i> <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">
        Select a user from the list and generate an attack scenario.
      </p>

      <div class="d-flex align-items-center mb-3">
        <select class="form-select form-select-sm me-2" v-model="selectedUser" style="width: 250px">
          <option value="admin">admin</option>
          <option value="gordonb">gordonb</option>
          <option value="1337">1337</option>
          <option value="pablo">pablo</option>
          <option value="smithy">smithy</option>
        </select>

        <button class="btn btn-primary btn-sm me-2" @click="performCommandInjection">
          Command Injection
        </button>
        <button class="btn btn-primary btn-sm me-2" @click="performSQLInjection">
          SQL Injection
        </button>
        <button class="btn btn-primary btn-sm me-2" @click="performCrossSiteScripting">
          Cross-site Scripting
        </button>
        <button class="btn btn-warning btn-sm me-2" @click="performZeroDayCommandInjection">
          Zero Day Command Injection
        </button>
        <button class="btn btn-warning btn-sm me-2" @click="performZeroDayCrossSiteScripting">
          Zero Day Cross-site Scripting
        </button>
        <button class="btn btn-secondary btn-sm" @click="resetResult">
          Reset
        </button>
      </div>

      <div v-if="jobResult" class="mt-4">
        <iframe ref="attackIframe" :srcdoc="jobResult" @load="adjustIframeHeight"
          style="width: 100%; border: 1px solid lightgray;"></iframe>
      </div>

    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Web Attacks</h5>
    </div>
    <div class="card-body">

      
      <table style="border-collapse: collapse; width: 100%;">
    <tr>
      <th style="border-bottom: 1px solid #ddd; padding: 8px;">Attack Name</th>
      <th style="border-bottom: 1px solid #ddd; padding: 8px;">Description</th>
      <th style="border-bottom: 1px solid #ddd; padding: 8px;">Example</th>
    </tr>
    <tr>
      <td style="border-bottom: 1px solid #ddd; padding: 8px;">Command Injection</td>
      <td style="border-bottom: 1px solid #ddd; padding: 8px;">Injects arbitrary commands into a vulnerable application for execution by the server.</td>
      <td style="border-bottom: 1px solid #ddd; padding: 8px; font-family:Courier;">;more /etc/passwd</td>
    </tr>
    <tr>
      <td style="border-bottom: 1px solid #ddd; padding: 8px;">SQL Injection</td>
      <td style="border-bottom: 1px solid #ddd; padding: 8px;">Exploits vulnerabilities in SQL query execution.</td>
      <td style="border-bottom: 1px solid #ddd; padding: 8px; font-family:Courier;">'OR 1=1#</td>
    </tr>
    <tr>
      <td style="border-bottom: 1px solid #ddd; padding: 8px;">XSS</td>
      <td style="border-bottom: 1px solid #ddd; padding: 8px;">Injects malicious scripts into web pages viewed by other users.</td>
      <td style="border-bottom: 1px solid #ddd; padding: 8px; font-family:Courier;">&lt;script&gt;alert("XSS-hack-attempt")&lt;/script&gt;</td>
    </tr>
    <tr>
      <td style="border-bottom: 1px solid #ddd; padding: 8px;">Zero Day SQL Injection</td>
      <td style="border-bottom: 1px solid #ddd; padding: 8px;">Exploits zero-day vulnerabilities in SQL query processing.</td>
      <td style="border-bottom: 1px solid #ddd; padding: 8px; font-family:Courier;">A%20'DIV'%20B%20-%20A%20'DIV%20B</td>
    </tr>
    <tr>
      <td style="border-bottom: 1px solid #ddd; padding: 8px;">Zero Day Command Injection</td>
      <td style="border-bottom: 1px solid #ddd; padding: 8px;">Injects commands exploiting zero-day vulnerabilities in server command processing.</td>
      <td style="border-bottom: 1px solid #ddd; padding: 8px; font-family:Courier;">/???/1? - /???/1?</td>
    </tr>
  </table>




      </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      selectedUser: "admin",
      jobResult: '',
      showHelp: false,
    };
  },
  methods: {
    performCommandInjection() {
      this.sendAttackRequest("command_injection");
    },
    performSQLInjection() {
      this.sendAttackRequest("sql_injection");
    },
    performCrossSiteScripting() {
      this.sendAttackRequest("xss");
    },
    performZeroDayCommandInjection() {
      this.sendAttackRequest("zero_day_sql_injection");
    },
    performZeroDayCrossSiteScripting() {
      this.sendAttackRequest("zero_day_command_injection");
    },
    sendAttackRequest(attackType) {
      const url = 'http://localhost:8080/web-attacks';
      const formData = new URLSearchParams();
      formData.append('type', attackType);
      formData.append('username', this.selectedUser);

      fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: formData
      })
        .then(response => response.text())
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