<template>
  <div class="card my-4">
    <div class="card-header d-flex justify-content-between align-items-center">
      <h5>FortiWeb Bootstrap</h5>
      <i class="bi bi-question-circle-fill bs-icon" @click="showHelp = !showHelp"></i>
    </div>
    <div class="card-body position-relative">

      <pre><code ref="codeBlock">{{ fileContent }}</code></pre>
    </div>
  </div>

  <!-- Help Card -->
  <div v-if="showHelp" class="card bg-light mb-3">
    <div class="card-header">
      <h5>About Docker</h5>
    </div>
    <div class="card-body">
      <p>Help content goes here...</p>
    </div>
  </div>
</template>

<script>
import hljs from 'highlight.js';
import 'highlight.js/styles/monokai.css';

export default {
  data() {
    return {
      showHelp: false,
      fileContent: '', // Variable pour stocker le contenu du fichier
    };
  },
  mounted() {
    this.highlightCode();
    fetch('/bootstrap.txt') // Assurez-vous que le chemin est correct
      .then(response => response.text())
      .then(text => {
        this.fileContent = text; // Stocke le contenu du fichier dans la variable
      })
      .catch(error => console.error('Error loading the file:', error));
  },
  methods: {
    highlightCode() {
      document.querySelectorAll('pre code').forEach((block) => {
        hljs.highlightBlock(block);
      });
    },
  }
};
</script>

<style>
.position-relative {
  position: relative;
}

.bs-icon {
  cursor: pointer;
  font-size: 1.5em;
}
</style>
