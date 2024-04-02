<template>
  <nav class="navbar fixed-top navbar-expand-lg bg-body-tertiary" data-bs-theme="dark">
    <div class="container">
      <router-link to="/" class="navbar-brand mb-0 h1">
        <img src="../assets/logo-fortiweb-white.svg" alt="FortiWeb Logo" height="30" class="me-2" />
        FortiWeb Demo Tool
      </router-link>

      <button
        class="navbar-toggler"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#navbarNav"
        aria-controls="navbarNav"
        aria-expanded="false"
        aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarNav">
        <ul class="navbar-nav me-auto mb-2 mb-lg-0">
          <!-- Web Protection Dropdown -->
          <li class="nav-item dropdown">
            <a class="nav-link dropdown-toggle" href="#" id="webProtectionDropdown" role="button" data-bs-toggle="dropdown" aria-expanded="false">
              Web Protection
            </a>
            <ul class="dropdown-menu" aria-labelledby="webProtectionDropdown">
              <li><router-link class="dropdown-item" to="/web-scan">Web Scan</router-link></li>
              <li><router-link class="dropdown-item" to="/traffic-generation">Traffic Generation</router-link></li>
              <li><router-link class="dropdown-item" to="/web-attacks">Web Attacks</router-link></li>
              <li><router-link class="dropdown-item" to="/machine-learning">Machine Learning</router-link></li>
              <li><router-link class="dropdown-item" to="/cookie-security">Cookie Security</router-link></li>
              <li><router-link class="dropdown-item" to="/credential-stuffing-defense">Credential Stuffing</router-link></li>
            </ul>
          </li>

          <!-- Bot Mitigation Dropdown -->
          <li class="nav-item dropdown">
            <a class="nav-link dropdown-toggle" href="#" id="botMitigationDropdown" role="button" data-bs-toggle="dropdown" aria-expanded="false">
              Bot Mitigation
            </a>
            <ul class="dropdown-menu" aria-labelledby="botMitigationDropdown">
              <li><router-link class="dropdown-item" to="/known-bots">Known Bots</router-link></li>
              <li><router-link class="dropdown-item" to="/bot-selenium">Selenium Playground</router-link></li>
              <li><router-link class="dropdown-item" to="/bot-deception">Bot Deception</router-link></li>
              <li><router-link class="dropdown-item" to="/bot-scraping">Web Scraping</router-link></li>
              <li><router-link class="dropdown-item" to="/http-request">HTTP Playground</router-link></li>
            </ul>
          </li>

          <!-- API Protection Dropdown -->
          <li class="nav-item dropdown">
            <a class="nav-link dropdown-toggle" href="#" id="apiProtectionDropdown" role="button" data-bs-toggle="dropdown" aria-expanded="false">
              API Protection
            </a>
            <ul class="dropdown-menu" aria-labelledby="apiProtectionDropdown">
              <li><router-link class="dropdown-item" to="/api-requests">API Requests</router-link></li>
              <li><router-link class="dropdown-item" to="/api-traffic-generation">API Traffic Generation</router-link></li>
            </ul>
          </li>

          <!-- REST API Dropdown -->
          <li class="nav-item dropdown">
            <a class="nav-link dropdown-toggle" href="#" id="restApiDropdown" role="button" data-bs-toggle="dropdown" aria-expanded="false"> REST API </a>
            <ul class="dropdown-menu" aria-labelledby="restApiDropdown">
              <li><router-link class="dropdown-item" to="/rest-api-management">Policy Management</router-link></li>
            </ul>
          </li>

          <!-- System Dropdown -->
          <li class="nav-item dropdown">
            <a class="nav-link dropdown-toggle" href="#" id="systemDropdown" role="button" data-bs-toggle="dropdown" aria-expanded="false"> System </a>
            <ul class="dropdown-menu" aria-labelledby="systemDropdown">
              <li>
                <router-link class="dropdown-item" to="/configuration"><font-awesome-icon :icon="['fas', 'pen-to-square']" class="me-3" />Configuration</router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/health-check"><font-awesome-icon :icon="['fas', 'heart-circle-check']" class="me-3" />Health Check</router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/app-docker"><font-awesome-icon :icon="['fab', 'docker']" class="me-3" />Docker Images</router-link>
              </li>
              <li>
                <router-link class="dropdown-item" to="/bootstrap"><font-awesome-icon :icon="['fas', 'file-import']" class="me-3" />FortiWeb Bootstrap</router-link>
              </li>
            </ul>
          </li>
        </ul>
      </div>
    </div>
  </nav>

  <div class="pt-3"></div>
</template>

<script>
export default {
  data() {
    return {
      currentSection: null,
    };
  },
};
</script>
