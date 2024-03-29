<template>

  <div class="card-container" style="display: flex; flex-wrap: wrap; justify-content: flex-start">

    <div class="card" style="margin: 15px">
      <div class="card-header"><i class="bi bi-shield-shaded me-2"></i><strong>Web Protection</strong></div>
      <ul class="list-group list-group-flush">
        <li class="list-group-item"><a href="/web-scan">Web Scan</a></li>
        <li class="list-group-item"><a href="/traffic-generation">Traffic Generation</a></li>
        <li class="list-group-item"><a href="/web-attacks">Web Attacks</a></li>
        <li class="list-group-item"><a href="/machine-learning">Machine Learning</a></li>
        <li class="list-group-item"><a href="/cookie-security">Cookie Security</a></li>
        <li class="list-group-item"><a href="/credential-stuffing-defense">Credential Stuffing Defense</a></li>
      </ul>
    </div>

    <div class="card" style="margin: 15px">
      <div class="card-header"><i class="bi bi-robot me-2"></i><strong>Bot Mitigation</strong></div>
      <ul class="list-group list-group-flush">
        <li class="list-group-item"><a href="/known-bots">Known Bots</a></li>
        <li class="list-group-item"><a href="/bot-selenium">Selenium Sandbox</a></li>
        <li class="list-group-item"><a href="/bot-deception">Bot Deception</a></li>
        <li class="list-group-item"><a href="/bot-scraping">Web Scraping</a></li>
      </ul>
    </div>

    <div class="card" style="margin: 15px">
      <div class="card-header"><i class="bi bi-braces me-2"></i><strong>API Protection</strong></div>
      <ul class="list-group list-group-flush">
        <li class="list-group-item"><a href="/api-requests">API Requests</a></li>
        <li class="list-group-item"><a href="/api-traffic-generation">API Traffic Generation</a></li>
      </ul>
    </div>

    <div class="card" style="margin: 15px">
      <div class="card-header"><i class="bi bi-terminal me-2"></i><strong>REST API Management</strong></div>
      <ul class="list-group list-group-flush">
        <li class="list-group-item"><a href="/rest-api-management">Policy Management</a></li>
      </ul>
    </div>
  </div>
</template>

<script></script>

<style>
.card-container {
  display: flex;
  justify-content: flex-start;
  /* This will space out your cards evenly */
  flex-wrap: wrap;
  /* This allows cards to wrap onto the next line if there's not enough space */
}

.list-group-item a {
  color: #000;
  /* Set your desired color, here it's black */
  text-decoration: none;
  /* Removes the underline */
}
</style>
