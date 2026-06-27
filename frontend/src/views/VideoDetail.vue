<template>
  <div v-if="video" class="detail">
    <div class="player-section">
      <video :src="video.video_url" controls class="video-player" />
    </div>
    <div class="video-info">
      <h1>{{ video.title }}</h1>
      <div class="info-bar">
        <div class="author" @click="$router.push(`/user/${video.user_id}`)">
          <el-avatar :size="40" :src="video.author?.avatar">{{ video.author?.username?.[0] }}</el-avatar>
          <span class="author-name">{{ video.author?.username }}</span>
        </div>
        <div class="actions">
          <el-button :type="liked ? 'primary' : 'default'" size="small" @click="toggleLike">
            <el-icon><CaretTop /></el-icon> {{ video.like_count || 0 }}
          </el-button>
          <el-button size="small" @click="coinDialog = true">
            <el-icon><Coin /></el-icon> {{ video.coin_count || 0 }}
          </el-button>
          <el-button size="small" @click="favDialog = true">
            <el-icon><Star /></el-icon> {{ video.fav_count || 0 }}
          </el-button>
          <el-button size="small" @click="doShare">
            <el-icon><Share /></el-icon> {{ video.share_count || 0 }}
          </el-button>
        </div>
      </div>
      <div class="desc">{{ video.description || '暂无简介' }}</div>
    </div>
    <CommentSection :videoId="video.id" />

    <el-dialog v-model="coinDialog" title="投币" width="300px">
      <div style="text-align:center">
        <el-radio-group v-model="coinCount">
          <el-radio :value="1">1 枚</el-radio>
          <el-radio :value="2">2 枚</el-radio>
        </el-radio-group>
      </div>
      <template #footer>
        <el-button @click="coinDialog = false">取消</el-button>
        <el-button type="primary" @click="doCoin">投币</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="favDialog" title="收藏到" width="400px">
      <div v-for="f in folders" :key="f.id" style="display:flex;align-items:center;gap:8px;padding:8px">
        <span style="flex:1">{{ f.name }}</span>
        <el-button size="small" @click="doFav(f.id)">收藏</el-button>
      </div>
      <div v-if="!folders.length" style="text-align:center;color:#999">暂无收藏夹，请先创建</div>
      <el-input v-model="newFolderName" placeholder="创建新收藏夹" style="margin-top:12px" />
      <el-button size="small" style="margin-top:8px" @click="createNewFolder">创建</el-button>
    </el-dialog>

    <h3 style="margin-top:24px">相关推荐</h3>
    <div class="video-grid" style="margin-top:10px">
      <VideoCard v-for="v in related" :key="v.id" :data="v" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import * as api from '@/api'
import VideoCard from '@/components/VideoCard.vue'
import CommentSection from '@/components/CommentSection.vue'
import { CaretTop, Coin, Star, Share } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const route = useRoute()
const video = ref(null)
const liked = ref(false)
const related = ref([])
const coinDialog = ref(false)
const coinCount = ref(1)
const favDialog = ref(false)
const folders = ref([])
const newFolderName = ref('')

onMounted(async () => {
  const id = route.params.id
  const res = await api.getVideoDetail(id)
  video.value = res.data
  loadFolders()
  loadRelated()
  try { const l = await api.likeStatus(id); liked.value = l.data.type === 1 } catch {}
})

async function loadFolders() {
  try { const r = await api.getFolders(); folders.value = r.data.list } catch {}
}

async function loadRelated() {
  try { const r = await api.getRelatedVideos(route.params.id, 1); related.value = r.data.list } catch {}
}

async function toggleLike() {
  if (liked.value) { await api.cancelLike(video.value.id); video.value.like_count--; liked.value = false }
  else { await api.likeVideo(video.value.id); video.value.like_count++; liked.value = true }
}

async function doCoin() {
  await api.giveCoin(video.value.id, coinCount.value)
  video.value.coin_count += coinCount.value
  coinDialog.value = false
  ElMessage.success('投币成功')
}

async function doFav(folderId) {
  await api.addToFolder({ folder_id: folderId, video_id: video.value.id })
  video.value.fav_count++
  favDialog.value = false
  ElMessage.success('收藏成功')
}

async function createNewFolder() {
  if (!newFolderName.value) return
  await api.createFolder({ name: newFolderName.value })
  newFolderName.value = ''
  loadFolders()
}

async function doShare() {
  const r = await api.shareVideo(video.value.id)
  video.value.share_count++
  await navigator.clipboard.writeText(location.origin + r.data.link)
  ElMessage.success('链接已复制')
}
</script>

<style scoped>
.player-section { background: #000; border-radius: 8px; overflow: hidden; }
.video-player { width: 100%; max-height: 560px; }
.video-info { padding: 16px 0; }
.video-info h1 { font-size: 20px; margin-bottom: 12px; }
.info-bar { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.author { display: flex; align-items: center; gap: 10px; cursor: pointer; }
.author-name { font-weight: 500; }
.actions { display: flex; gap: 8px; }
.desc { background: #f6f7f8; padding: 12px; border-radius: 6px; font-size: 14px; color: var(--text-secondary); white-space: pre-wrap; }
</style>
