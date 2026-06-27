<template>
  <div v-if="profile">
    <div class="profile-header">
      <el-avatar :size="80" :src="profile.avatar">{{ profile.username?.[0] }}</el-avatar>
      <div class="profile-info">
        <h2>{{ profile.username }}</h2>
        <p>{{ profile.bio }}</p>
        <div class="stats">
          <span>{{ profile.follower_count || 0 }} 粉丝</span>
          <span>{{ profile.followee_count || 0 }} 关注</span>
          <span>{{ profile.video_count || 0 }} 视频</span>
        </div>
        <el-button v-if="!isSelf" size="small" :type="followed ? 'default' : 'primary'" @click="toggleFollow">
          {{ followed ? '已关注' : '关注' }}
        </el-button>
      </div>
    </div>
    <h3 style="margin-top:24px">投稿视频</h3>
    <div class="video-grid">
      <VideoCard v-for="v in videos" :key="v.id" :data="v" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import * as api from '@/api'
import VideoCard from '@/components/VideoCard.vue'

const route = useRoute()
const profile = ref(null)
const videos = ref([])
const followed = ref(false)
const isSelf = computed(() => +route.params.id === +localStorage.getItem('user_id'))

onMounted(async () => {
  const id = route.params.id
  profile.value = (await api.getProfile(id)).data
  videos.value = (await api.getUserVideos(id, 1)).data.list
  if (!isSelf.value) {
    try { const r = await api.getFollowers(localStorage.getItem('user_id'), 1); followed.value = true } catch { followed.value = false }
  }
})

async function toggleFollow() {
  if (followed.value) { await api.unfollow(route.params.id); followed.value = false }
  else { await api.follow(route.params.id); followed.value = true }
}
</script>

<style scoped>
.profile-header { display: flex; gap: 24px; align-items: center; padding: 24px; background: #fff; border-radius: 8px; }
.profile-info h2 { margin-bottom: 8px; }
.profile-info p { color: var(--text-secondary); font-size: 14px; margin-bottom: 8px; }
.stats { display: flex; gap: 20px; margin-bottom: 12px; font-size: 14px; }
.stats span { color: var(--text-secondary); }
</style>
