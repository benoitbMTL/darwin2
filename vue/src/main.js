import { createApp } from 'vue'; // Vue 3 import
import App from './App.vue';
import router from './router';

import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import 'bootstrap-icons/font/bootstrap-icons.css';

// Import the necessary icons from FontAwesome
import { library } from '@fortawesome/fontawesome-svg-core';
import { faGithub } from '@fortawesome/free-brands-svg-icons/faGithub'; // Correct import for GitHub
import { faDocker } from '@fortawesome/free-brands-svg-icons/faDocker'; // Correct import for Docker
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

// Add the imported icons to the library
library.add(faGithub, faDocker);

import './assets/styles/global.css';

const app = createApp(App); // Create the app using Vue 3 syntax

// Register the FontAwesomeIcon component globally
app.component('font-awesome-icon', FontAwesomeIcon);

app.use(router); // Use the router

app.mount('#app'); // Mount the app to the DOM
