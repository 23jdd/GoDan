<template>
  <div>
    <section class="page-header">
      <h1>直播广场</h1>
      <p>直播页先完成视觉与信息结构，后续可继续接入真实流媒体与 WebSocket 弹幕。</p>
    </section>

    <div class="live-grid">
      <div v-for="r in rooms" :key="r.id" class="live-card card-surface" @click="$router.push(`/live/${r.id}`)">
        <div class="live-cover">
          <img :src="r.cover_url" />
          <span class="live-badge">直播中</span>
          <span class="live-viewers">{{ r.viewer_count || 0 }} 观看</span>
        </div>
        <div class="live-info">
          <div class="live-title">{{ r.title }}</div>
          <div class="live-user">{{ r.username }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import * as api from '@/api'

const rooms = ref([])

onMounted(async () => {
  const res = await api.getLiveList(1)
  rooms.value = res.data.list
})
</script>

<style scoped>
.page-header {
  margin-bottom: 20px;
}

.page-header h1 {
  font-size: 32px;
}

.page-header p {
  margin-top: 10px;
  color: var(--text-secondary);
}

.live-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.live-card {
  cursor: pointer;
  overflow: hidden;
}

.live-cover {
  position: relative;
  aspect-ratio: 16/9;
  background: #eee;
}

.live-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.live-badge {
  position: absolute;
  left: 8px;
  top: 8px;
  background: #ff4d4f;
  color: #fff;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
}

.live-viewers {
  position: absolute;
  right: 8px;
  top: 8px;
  background: rgba(0,0,0,.6);
  color: #fff;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
}

.live-info {
  padding: 12px 14px 16px;
  background: #fff;
}

.live-title {
  font-size: 15px;
  font-weight: 600;
}

.live-user {
  font-size: 12px;
  color: var(--text-secondary);
  margin-top: 6px;
}
</style>
