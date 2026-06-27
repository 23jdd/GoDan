<template>
  <div>
    <div class="page-title">直播列表</div>
    <div class="live-grid">
      <div v-for="r in rooms" :key="r.id" class="live-card" @click="$router.push(`/live/${r.id}`)">
        <div class="live-cover">
          <img :src="r.cover_url || '/placeholder.png'" />
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
import { ref, onMounted } from 'vue'
import * as api from '@/api'

const rooms = ref([])

onMounted(async () => {
  const res = await api.getLiveList(1)
  rooms.value = res.data.list
})
</script>

<style scoped>
.live-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 16px; }
.live-card { cursor: pointer; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 4px rgba(0,0,0,0.1); }
.live-cover { position: relative; aspect-ratio: 16/9; background: #eee; }
.live-cover img { width: 100%; height: 100%; object-fit: cover; }
.live-badge { position: absolute; left: 8px; top: 8px; background: #ff4d4f; color: #fff; font-size: 12px; padding: 2px 8px; border-radius: 4px; }
.live-viewers { position: absolute; right: 8px; top: 8px; background: rgba(0,0,0,.6); color: #fff; font-size: 12px; padding: 2px 8px; border-radius: 4px; }
.live-info { padding: 10px; background: #fff; }
.live-title { font-size: 14px; font-weight: 500; }
.live-user { font-size: 12px; color: var(--text-secondary); margin-top: 4px; }
</style>
