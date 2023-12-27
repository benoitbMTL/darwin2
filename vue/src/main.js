import { createApp } from 'vue' // Import the createApp function from the Vue framework
import App from './App.vue' // Import the root component App from the current directory
import router from './router'; // Import the router configuration from the router directory
import 'bootstrap/dist/css/bootstrap.min.css'; // Import Bootstrap CSS for styling
import 'bootstrap/dist/js/bootstrap.bundle.min.js'; // Import Bootstrap JavaScript for interactivity (includes Popper.js)
import 'bootstrap-icons/font/bootstrap-icons.css';
import './assets/styles/global.css'; // Import global CSS styles from the assets directory

const app = createApp(App); // Create a new Vue application instance with the root component App
app.use(router); // Use the Vue Router for handling navigation within the app
app.mount('#app'); // Mount the Vue application to the DOM element with the id 'app'