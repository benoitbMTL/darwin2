<template>
  <div class="card my-4">

    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Cookie Security</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem;" @click="showHelp = !showHelp"></i>
    </div>

    <div class="card-body">
      <p class="card-text">Select a user and perform a Cookie Based Attack.</p>

      <div class="d-flex align-items-center mb-3">
        <select class="form-select form-select-sm me-2" v-model="selectedUser" style="width: 250px">
          <option value="admin">admin</option>
          <option value="gordonb">gordonb</option>
          <option value="1337">1337</option>
          <option value="pablo">pablo</option>
          <option value="smithy">smithy</option>
        </select>


        <button type="button" class="btn btn-primary btn-sm me-2" @click="performCookieSecurity">Manipulate
          Cookie</button>

        <button class="btn btn-secondary btn-sm" @click="resetResult">
          Reset
        </button>

      </div>


      <div>
        <div v-if="initialCookie" class="mt-4">
          <h6><i class="bi bi-1-circle-fill me-2"></i>You are now authenticated. Your cookie security level is set to
            <span style="color: red;">low</span>.</h6>
          <iframe ref="responseIframe" :srcdoc="initialCookie"
            style="width: 100%; height: 80px; border: 1px solid lightgray;"></iframe>
        </div>

        <div v-if="modifiedCookie" class="mt-3">
          <h6><i class="bi bi-2-circle-fill me-2"></i>Let's change the cookie security level to <span
              style="color: red;">medium</span>.</h6>
          <iframe ref="responseIframe" :srcdoc="modifiedCookie"
            style="width: 100%; height: 80px; border: 1px solid lightgray;"></iframe>
        </div>

        <div v-if="webPageHTML" class="mt-4">
          <h6><i class="bi bi-3-circle-fill me-2"></i>Let's connect again to the web app with the new crafted cookie.</h6>
          <iframe ref="responseIframe" :srcdoc="webPageHTML" @load="adjustIframeHeight"
            style="width: 100%; border: 1px solid lightgray;"></iframe>
        </div>
      </div>

    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Cookie Security Policy</h5>
    </div>
    <div class="card-body">
      <p>A cookie security policy allows you to configure FortiWeb features that prevent cookie-based attacks and apply
        them in a protection profile. For example, a policy can enable cookie poisoning detection, encrypt the cookies
        issued by a back-end server, and add security attributes to cookies.</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      showHelp: false,
      initialCookie: '',
      modifiedCookie: '',
      webPageHTML: '',
      selectedUser: "admin",
    };
  },

  methods: {
    performCookieSecurity() {
      const formData = new URLSearchParams();
      formData.append('username', this.selectedUser);

      fetch('/cookie-security', {
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
          return response.json();
        })
        .then(data => {
          this.initialCookie = data.initialCookie;
          this.modifiedCookie = data.modifiedCookie;
          this.webPageHTML = data.webPageHTML;
        })
        .catch(error => {
          console.error('Error:', error);
          this.initialCookie = 'Failed to perform cookie security';
          this.modifiedCookie = '';
          this.webPageHTML = '';
        });
    },


    adjustIframeHeight() {
      const iframe = this.$refs.responseIframe;
      if (iframe && iframe.contentWindow && iframe.contentWindow.document.body) {
        iframe.style.height = (iframe.contentWindow.document.body.scrollHeight + 30) + 'px';
      }
    },

    resetResult() {
      this.selectedOption = "admin"; // Reset selected option
      this.initialCookie = "";
      this.modifiedCookie = "";
      this.webPageHTML = "";
    },


  }
};


</script>

<style></style>

