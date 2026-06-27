<template>
  <div id="app-root" class="app-shell">
    <Navbar />
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import Navbar from '@/components/Navbar.vue'
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
}
</style>
