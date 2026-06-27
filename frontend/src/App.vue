<template>
  <div id="app-root">
    <Navbar />
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import Navbar from '@/components/Navbar.vue'
</script>
