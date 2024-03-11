<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Web Attacks</h5>
      <i
        class="bi bi-question-circle-fill bs-icon"
        style="font-size: 1.5rem"
        @click="showHelp = !showHelp"></i>
      <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">
        Select a user and an attack type from the list, then click "Run" to
        generate the attack scenario.
      </p>

      <div class="d-flex align-items-center mb-3 flex-wrap">
        <select
          class="form-select form-select-sm me-2"
          v-model="selectedUser"
          style="width: 250px">
          <option value="admin">admin</option>
          <option value="gordonb">gordonb</option>
          <option value="1337">1337</option>
          <option value="pablo">pablo</option>
          <option value="smithy">smithy</option>
        </select>

        <select
          class="form-select form-select-sm me-2"
          v-model="selectedAttackType"
          style="width: 250px">
          <option value="command_injection">Command Injection</option>
          <option value="sql_injection">SQL Injection</option>
          <option value="xss">Cross-site Scripting</option>
          <option value="os_command_injection">
            OS Command Injection Attacks
          </option>
          <option value="coldfusion_injection">Coldfusion Injection</option>
          <option value="ldap_injection">LDAP Injection</option>
          <option value="session_fixation">Session Fixation</option>
          <option value="file_injection">File Injection</option>
          <option value="php_injection">PHP Injection</option>
          <option value="ssi_injection">SSI Injection</option>
          <option value="updf_xss">UPDF XSS</option>
          <option value="email_injection">Email Injection</option>
          <option value="http_response_splitting">
            HTTP Response Splitting
          </option>
          <option value="rfi_injection">RFI Injection</option>
          <option value="lfi_injection">LFI Injection</option>
          <option value="src_disclosure">SRC Disclosure</option>
          <option value="java_method_injection">Java Method Injection</option>
          <option value="directory_traversal">Directory Traversal</option>
          <option value="format_string_attack">Format String Attack</option>
          <option value="xpath_xquery_injection">Xpath/XQuery Injection</option>
          <option value="xslt_injection">XSLT Injection</option>
          <option value="trojans">Trojans</option>
        </select>

        <button class="btn btn-primary btn-sm me-2" @click="performAttack">
          Run
        </button>

        <button class="btn btn-secondary btn-sm me-2" @click="resetResult">
          Reset
        </button>
      </div>

      <div v-if="jobResult" class="mt-4 mb-3">
        <h6>{{ currentAttackName }} Result:</h6>
        <iframe
          ref="attackIframe"
          :srcdoc="jobResult"
          @load="adjustIframeHeight"
          style="width: 100%; border: 1px solid lightgray"></iframe>
      </div>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Web Attacks</h5>
    </div>
    <div class="card-body">
      <table style="border-collapse: collapse; width: 100%">
        <tr>
          <th style="border-bottom: 1px solid #ddd; padding: 8px">
            Attack Name
          </th>
          <th style="border-bottom: 1px solid #ddd; padding: 8px">
            Description
          </th>
          <th style="border-bottom: 1px solid #ddd; padding: 8px">Example</th>
        </tr>
        <tr>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            Command Injection
          </td>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            Injects arbitrary commands into a vulnerable application for
            execution by the server.
          </td>
          <td
            style="
              border-bottom: 1px solid #ddd;
              padding: 8px;
              font-family: Courier;
            ">
            ;more /etc/passwd
          </td>
        </tr>
        <tr>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            SQL Injection
          </td>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            Exploits vulnerabilities in SQL query execution.
          </td>
          <td
            style="
              border-bottom: 1px solid #ddd;
              padding: 8px;
              font-family: Courier;
            ">
            'OR 1=1#
          </td>
        </tr>
        <tr>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            Cross-site Scripting
          </td>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            Injects malicious scripts into web pages viewed by other users.
          </td>
          <td
            style="
              border-bottom: 1px solid #ddd;
              padding: 8px;
              font-family: Courier;
            ">
            &lt;script&gt;alert("XSS-hack-attempt")&lt;/script&gt;
          </td>
        </tr>
        <tr>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            OS Command Injection
          </td>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            Execution of arbitrary commands via OS command injection.
          </td>
          <td
            style="
              border-bottom: 1px solid #ddd;
              padding: 8px;
              font-family: Courier;
            ">
            ;cc evil.c
          </td>
        </tr>
        <tr>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            ColdFusion Injection
          </td>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            Manipulates ColdFusion data through unauthorized registry access.
          </td>
          <td
            style="
              border-bottom: 1px solid #ddd;
              padding: 8px;
              font-family: Courier;
            ">
            &lt;CFNEWINTERNALREGISTRY ACTION="Set"
            BRANCH="HKEY_LOCAL_MACHINE\Software\Allaire\ColdFusion\CurrentVersion\Server"
            NAME="test" TYPE="String" VALUE="0"&gt;
          </td>
        </tr>
        <tr>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            LDAP Injection
          </td>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            Exploits web applications that construct LDAP statements from user
            input.
          </td>
          <td
            style="
              border-bottom: 1px solid #ddd;
              padding: 8px;
              font-family: Courier;
            ">
            *)(uid=*))(|(uid=*
          </td>
        </tr>
        <tr>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            Session Fixation
          </td>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            Manipulates session IDs to hijack a user session.
          </td>
          <td
            style="
              border-bottom: 1px solid #ddd;
              padding: 8px;
              font-family: Courier;
            ">
            &lt;script&gt;document.cookie="sessionid=1234; Expires=Friday,
            1-Jan-2010 00:00:00 GMT";&lt;/script&gt;
          </td>
        </tr>
        <tr>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            File Injection
          </td>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            Injects malicious files into a server or application.
          </td>
          <td
            style="
              border-bottom: 1px solid #ddd;
              padding: 8px;
              font-family: Courier;
            ">
            C:/boot.ini
          </td>
        </tr>
        <tr>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            PHP Injection
          </td>
          <td style="border-bottom: 1px solid #ddd; padding: 8px">
            Injects PHP code into an application, leading to unauthorized
            command execution.
          </td>
          <td
            style="
              border-bottom: 1px solid #ddd;
              padding: 8px;
              font-family: Courier;
            ">
            abc;$_SESSION[authuser]=1
          </td>
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
      selectedAttackType: "command_injection",
      jobResult: "",
      showHelp: false,
    };
  },
  methods: {
    performAttack() {
      switch (this.selectedAttackType) {
        case "command_injection":
          this.currentAttackName = "Command Injection";
          break;
        case "sql_injection":
          this.currentAttackName = "SQL Injection";
          break;
        case "xss":
          this.currentAttackName = "Cross-site Scripting";
          break;
        // New cases
        case "os_command_injection":
          this.currentAttackName = "OS Command Injection Attacks";
          break;
        case "coldfusion_injection":
          this.currentAttackName = "Coldfusion Injection";
          break;
        case "ldap_injection":
          this.currentAttackName = "LDAP Injection";
          break;
        case "session_fixation":
          this.currentAttackName = "Session Fixation";
          break;
        case "file_injection":
          this.currentAttackName = "File Injection";
          break;
        case "php_injection":
          this.currentAttackName = "PHP Injection";
          break;
        case "ssi_injection":
          this.currentAttackName = "SSI Injection";
          break;
        case "updf_xss":
          this.currentAttackName = "UPDF XSS";
          break;
        case "email_injection":
          this.currentAttackName = "Email Injection";
          break;
        case "http_response_splitting":
          this.currentAttackName = "HTTP Response Splitting";
          break;
        case "rfi_injection":
          this.currentAttackName = "RFI Injection";
          break;
        case "lfi_injection":
          this.currentAttackName = "LFI Injection";
          break;
        case "src_disclosure":
          this.currentAttackName = "SRC Disclosure";
          break;
        case "java_method_injection":
          this.currentAttackName = "Java Method Injection";
          break;
        case "directory_traversal":
          this.currentAttackName = "Directory Traversal";
          break;
        case "format_string_attack":
          this.currentAttackName = "Format String Attack";
          break;
        case "xpath_xquery_injection":
          this.currentAttackName = "Xpath/XQuery Injection";
          break;
        case "xslt_injection":
          this.currentAttackName = "XSLT Injection";
          break;
        case "trojans":
          this.currentAttackName = "Trojans";
          break;
        default:
          this.currentAttackName = "";
      }
      this.sendAttackRequest(this.selectedAttackType);
    },

    sendAttackRequest(attackType) {
      const url = "localhost:8080/web-attacks";
      const formData = new URLSearchParams();
      formData.append("type", attackType);
      formData.append("username", this.selectedUser);

      fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        body: formData,
      })
        .then((response) => response.text())
        .then((html) => {
          this.jobResult = html;
        })
        .catch((error) => {
          console.error("Error:", error);
          this.jobResult = "Failed to perform attack";
        });
    },

    adjustIframeHeight() {
      const iframe = this.$refs.attackIframe;
      if (
        iframe &&
        iframe.contentWindow &&
        iframe.contentWindow.document.body
      ) {
        iframe.style.height =
          iframe.contentWindow.document.body.scrollHeight + 30 + "px";
      }
    },

    resetResult() {
      this.selectedOption = "All"; // Reset selected option
      this.jobResult = ""; // Clear Result
    },
  },
};
</script>

<style></style>
