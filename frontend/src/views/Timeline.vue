<template>
  <div>
    <div class="page-title">关注动态</div>
    <div v-for="a in list" :key="a.id" class="activity-item">
      <el-avatar :size="36" :src="a.avatar">{{ a.username?.[0] }}</el-avatar>
      <div class="activity-body">
        <span class="act-user">{{ a.username }}</span>
        <span class="act-action">{{ actionText(a.type) }}</span>
        <span class="act-time">{{ a.created_at?.slice(0,10) }}</span>
      </div>
    </div>
    <div v-if="!list.length" style="text-align:center;color:#999;padding:60px">
      关注更多用户，动态流就会出现在这里
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import * as api from '@/api'

const list = ref([])

onMounted(async () => {
  try { const r = await api.getTimeline(1); list.value = r.data.list } catch {}
})

function actionText(type) {
  const map = { 1: '投稿了视频', 2: '点赞了视频', 3: '给视频投了币', 4: '收藏了视频', 5: '分享了视频' }
  return map[type] || '动态'
}
</script>

<style scoped>
.activity-item { display: flex; gap: 12px; padding: 12px; background: #fff; border-radius: 8px; margin-bottom: 8px; }
.activity-body { font-size: 14px; }
.act-user { font-weight: 500; color: var(--primary); }
.act-action { color: var(--text-secondary); margin: 0 4px; }
.act-time { font-size: 12px; color: #ccc; display: block; margin-top: 4px; }
</style>
