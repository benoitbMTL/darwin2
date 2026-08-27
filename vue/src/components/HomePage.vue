<template>
  <div class="container">
    <div class="row mb-4 mt-4">
      <div class="col-2">
        <div class="card h-100">

          <!-- Card content -->
          <div class="card-header"><i class="bi bi-shield-shaded me-2"></i><strong>Web Protection</strong></div>
          <ul class="list-group list-group-flush">
            <li class="list-group-item"><a href="/web-scan">Web Scan</a></li>
            <li class="list-group-item"><a href="/traffic-generation">Traffic Generation</a></li>
            <li class="list-group-item"><a href="/web-attacks">Web Attacks</a></li>
            <li class="list-group-item"><a href="/machine-learning">Machine Learning</a></li>
            <li class="list-group-item"><a href="/cookie-security">Cookie Security</a></li>
            <li class="list-group-item"><a href="/credential-stuffing-defense">Credential Stuffing</a></li>
          </ul>

        </div>
      </div>
      <div class="col-2">
        <div class="card h-100">

          <!-- Card content -->
          <div class="card-header"><i class="bi bi-robot me-2"></i><strong>Bot Mitigation</strong></div>
          <ul class="list-group list-group-flush">
            <li class="list-group-item"><a href="/known-bots">Known Bots</a></li>
            <li class="list-group-item"><a href="/bot-selenium">Selenium Playground</a></li>
            <li class="list-group-item"><a href="/bot-deception">Bot Deception</a></li>
            <li class="list-group-item"><a href="/bot-scraping">Web Scraping</a></li>
            <li class="list-group-item"><a href="/http-request">HTTP Playground</a></li>
          </ul>
        </div>
      </div>
      <div class="col-2">
        <div class="card h-100">

          <!-- Card content -->
          <div class="card-header"><i class="bi bi-braces me-2"></i><strong>API Protection</strong></div>
          <ul class="list-group list-group-flush">
            <li class="list-group-item"><a href="/api-requests">API Requests</a></li>
            <li class="list-group-item"><a href="/api-traffic-generation">API Traffic Generation</a></li>
          </ul>
        </div>
      </div>
      <div class="col-2">
        <div class="card h-100">

          <!-- Card content -->
          <div class="card-header"><i class="bi bi-terminal me-2"></i><strong>REST API</strong></div>
          <ul class="list-group list-group-flush">
            <li class="list-group-item"><a href="/rest-api-management">Policy Management</a></li>
          </ul>
        </div>
      </div>
    </div>

    <div class="row mb-4">
      <div class="col-8">
        <div class="card overflow-hidden">
          <img
            :src="homepageDiagram"
            class="img-fluid w-100"
            alt="Darwin2 traffic flow through FortiWeb to the demonstration applications"
            decoding="async" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import homepageDiagram from "../assets/homepage.png";
</script>


<style>
.list-group-item a {
  color: #000;
  /* Set your desired color, here it's black */
  text-decoration: none;
  /* Removes the underline */
}
</style>
