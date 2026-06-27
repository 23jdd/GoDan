<template>
  <div>
    <div class="page-title">搜索结果: {{ q }}</div>
    <div class="video-grid">
      <VideoCard v-for="v in videos" :key="v.id" :data="v" />
    </div>
    <div v-if="!videos.length" style="text-align:center;color:#999;padding:60px">暂无结果</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import * as api from '@/api'
import VideoCard from '@/components/VideoCard.vue'

const route = useRoute()
const q = ref(route.query.q || '')
const videos = ref([])

onMounted(async () => {
  if (q.value) {
    const res = await api.searchVideos(q.value, 1)
    videos.value = res.data.list
  }
})
</script>
