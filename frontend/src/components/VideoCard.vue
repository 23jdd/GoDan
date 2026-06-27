<template>
  <div class="video-card card-surface" @click="$router.push(`/video/${data.id}`)">
    <div class="cover-wrap">
      <img :src="data.cover_url" class="cover" />
      <div class="cover-overlay">
        <span class="pill">{{ data.category }}</span>
        <span class="duration">{{ formatDuration(data.duration) }}</span>
      </div>
    </div>

    <div class="info">
      <h3 class="title">{{ data.title }}</h3>
      <p class="author">{{ data.author?.username || data.username || 'GoDan 创作者' }}</p>
      <div class="meta">
        <span>{{ formatCount(data.play_count) }} 播放</span>
        <span>{{ formatCount(data.danmaku_count || data.like_count) }} 弹幕</span>
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

function formatCount(value) {
  if (!value) return '0'
  if (value >= 10000) return `${(value / 10000).toFixed(1)}万`
  return String(value)
}
</script>

<style scoped>
.video-card {
  cursor: pointer;
  overflow: hidden;
  transition: transform 0.22s ease, box-shadow 0.22s ease;
}

.video-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 28px 48px rgba(17, 24, 39, 0.12);
}

.cover-wrap {
  position: relative;
  overflow: hidden;
  background: #dbeafe;
}

.cover {
  width: 100%;
  aspect-ratio: 16 / 9;
  display: block;
  object-fit: cover;
}

.cover-overlay {
  position: absolute;
  inset: auto 0 0 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 12px;
  background: linear-gradient(180deg, transparent, rgba(15, 23, 42, 0.75));
}

.pill,
.duration {
  color: #fff;
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.18);
}

.info {
  padding: 14px;
}

.title {
  font-size: 16px;
  line-height: 1.45;
  min-height: 46px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.author {
  margin-top: 8px;
  font-size: 13px;
  color: var(--text-secondary);
}

.meta {
  display: flex;
  gap: 12px;
  margin-top: 10px;
  color: var(--text-muted);
  font-size: 12px;
}
</style>
