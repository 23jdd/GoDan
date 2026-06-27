<template>
  <div>
    <section class="page-header">
      <h1>关注动态</h1>
      <p>关注的创作者最近做了什么，一眼看清。</p>
    </section>

    <div v-for="a in list" :key="a.id" class="activity-item card-surface">
      <a-avatar :size="44" :src="a.avatar">{{ a.username?.[0] }}</a-avatar>
      <div class="activity-body">
        <span class="act-user">{{ a.username }}</span>
        <span class="act-action">{{ actionText(a.type) }}</span>
        <p class="act-content">{{ a.content }}</p>
        <span class="act-time">{{ a.created_at?.slice(0,10) }}</span>
      </div>
    </div>

    <a-empty v-if="!list.length" description="多关注一些 UP，动态流就会出现在这里。" class="empty-state" />
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import * as api from '@/api'

const list = ref([])

onMounted(async () => {
  const r = await api.getTimeline(1)
  list.value = r.data.list
})

function actionText(type) {
  const map = { 1: '发布了视频', 2: '点赞了视频', 3: '给视频投了币', 4: '收藏了内容', 5: '分享了视频' }
  return map[type] || '更新了动态'
}
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

.activity-item {
  display: flex;
  gap: 14px;
  padding: 18px;
  margin-bottom: 12px;
}

.activity-body {
  font-size: 14px;
}

.act-user {
  font-weight: 500;
  color: var(--primary);
}

.act-action {
  color: var(--text-secondary);
  margin-left: 6px;
}

.act-content {
  margin-top: 10px;
  color: #334155;
}

.act-time {
  font-size: 12px;
  color: #94a3b8;
  display: block;
  margin-top: 8px;
}

.empty-state {
  padding: 72px 0;
}
</style>
