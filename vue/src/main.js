import { createApp } from 'vue'; // Vue 3 import
import App from './App.vue';
import router from './router';

import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import 'bootstrap-icons/font/bootstrap-icons.css';

// Import the entire sets of free icons from FontAwesome
import { library } from '@fortawesome/fontawesome-svg-core';
import { fas } from '@fortawesome/free-solid-svg-icons';  // Free solid icons
import { far } from '@fortawesome/free-regular-svg-icons'; // Free regular icons
import { fab } from '@fortawesome/free-brands-svg-icons'; // Free brand icons
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

// Add the imported icons to the library
library.add(fas, far, fab);

// Register the FontAwesomeIcon component globally
app.component('font-awesome-icon', FontAwesomeIcon);

import './assets/styles/global.css';

const app = createApp(App); // Create the app using Vue 3 syntax

app.use(router); // Use the router

app.mount('#app'); // Mount the app to the DOM
