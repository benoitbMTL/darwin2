<template>
  <div id="app" class="container mt-5">
    <NavBar />
    <router-view></router-view> <!-- Vue Router will render the component here based on the current route -->

  </div>
</template>

<script>
import NavBar from "./components/NavBar.vue";


export default {
  components: {
    NavBar,
  },
};
</script>

<style></style>
