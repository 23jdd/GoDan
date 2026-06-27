<template>
  <div class="search-page">
    <section class="page-header">
      <h1>搜索结果</h1>
      <p>关键词 “{{ q }}” 共找到 {{ videos.length }} 个内容结果。</p>
    </section>

    <div v-if="videos.length" class="video-grid">
      <VideoCard v-for="v in videos" :key="v.id" :data="v" />
    </div>

    <a-empty v-else description="暂时没有匹配结果，换个关键词试试。" class="empty-state" />
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import * as api from '@/api'
import VideoCard from '@/components/VideoCard.vue'

const route = useRoute()
const q = ref(route.query.q || '')
const videos = ref([])

onMounted(load)
watch(() => route.query.q, load)

async function load() {
  q.value = route.query.q || ''
  if (!q.value) {
    videos.value = []
    return
  }
  const res = await api.searchVideos(q.value, 1)
  videos.value = res.data.list
}
</script>

<style scoped>
.page-header h1 {
  font-size: 32px;
}

.page-header p {
  margin-top: 10px;
  color: var(--text-secondary);
}

.empty-state {
  padding: 72px 0;
}
</style>
