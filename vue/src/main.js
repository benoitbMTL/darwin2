import { createApp } from 'vue'; // Vue 3 import
import App from './App.vue';
import router from './router';

import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import 'bootstrap-icons/font/bootstrap-icons.css';
import 'highlight.js/styles/monokai.css';

import { library } from '@fortawesome/fontawesome-svg-core';
import { faFileImport, faHeartCircleCheck, faPenToSquare } from '@fortawesome/free-solid-svg-icons';
import { faDocker, faGithub } from '@fortawesome/free-brands-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

library.add(faDocker, faFileImport, faGithub, faHeartCircleCheck, faPenToSquare);

import './assets/styles/global.css';

const app = createApp(App); // Create the app using Vue 3 syntax

// Register the FontAwesomeIcon component globally
app.component('font-awesome-icon', FontAwesomeIcon);

app.use(router); // Use the router

app.mount('#app'); // Mount the app to the DOM
