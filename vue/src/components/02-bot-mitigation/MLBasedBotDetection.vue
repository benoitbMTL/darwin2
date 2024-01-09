<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>Machine Learning</h5>
      <i class="bi bi-question-circle-fill bs-icon" style="font-size: 1.5rem;" @click="showHelp = !showHelp"></i> <!-- Bootstrap icon for help -->
    </div>
    <div class="card-body">
      <p class="card-text">Machine Learning Options.</p>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>Brief description of Bot techniques</h5>
    </div>
    <div class="card-body">
      <p>
        
<ul>
    <li><strong>Crawler:</strong> A crawler, also known as a spider or bot, automatically browses the web to collect data, index content for search engines, or for web archiving.</li>
    <li><strong>Vulnerability Scanning:</strong> This is a systematic process that identifies and classifies vulnerabilities in computers, networks, and software systems. It's often used in cybersecurity to assess the security of a system.</li>
    <li><strong>Slow Attack:</strong> Typically referred to as Slowloris or similar attacks, this technique involves sending partial requests to a web server, holding connections open for as long as possible. It's a form of attack aiming to exhaust server resources.</li>
    <li><strong>Content Scraping:</strong> The process of extracting useful information from websites. This can be done for various reasons, including data gathering, price comparison, or republishing information on other sites.</li>
    <li><strong>Illegal User Scan:</strong> This involves scanning a network or system for the presence of unauthorized users or intruders. It's a security measure to detect potential breaches or unauthorized access.</li>
    <li><strong>Brute Force Attack:</strong> A trial-and-error method used to decode encrypted data such as passwords or Data Encryption Standard (DES) keys. This attack method involves systematically checking all possible keys until the correct key is found.</li>
</ul>

      </p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      showHelp: false,
    };
  },
};
</script>

<style></style>

