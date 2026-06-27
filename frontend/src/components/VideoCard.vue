<template>
  <div class="video-card" @click="$router.push(`/video/${data.id}`)">
    <div class="cover-wrap">
      <img :src="data.cover_url || '/placeholder.png'" class="cover" />
      <span class="duration">{{ formatDuration(data.duration) }}</span>
    </div>
    <div class="info">
      <div class="title">{{ data.title }}</div>
      <div class="meta">
        <span>{{ data.play_count || 0 }} 播放 · {{ data.like_count || 0 }} 点赞</span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({ data: Object })

function formatDuration(s) {
  if (!s) return '00:00'
  const m = Math.floor(s / 60)
  const sec = s % 60
  return `${String(m).padStart(2,'0')}:${String(sec).padStart(2,'0')}`
}
</script>

<style scoped>
.video-card { cursor: pointer; transition: transform .2s; }
.video-card:hover { transform: translateY(-4px); }
.cover-wrap { position: relative; border-radius: 6px; overflow: hidden; background: #eee; }
.cover { width: 100%; aspect-ratio: 16/9; display: block; object-fit: cover; }
.duration { position: absolute; right: 6px; bottom: 6px; background: rgba(0,0,0,.7); color: #fff; font-size: 11px; padding: 2px 6px; border-radius: 4px; }
.info { padding: 8px 0; }
.title { font-size: 14px; font-weight: 500; line-height: 1.4; height: 40px; overflow: hidden; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; }
.meta { font-size: 12px; color: var(--text-secondary); margin-top: 6px; }
</style>
