<template>
  <div v-if="video" class="detail-layout">
    <section class="main-column">
      <div class="player-shell card-surface">
        <ArtPlayerPlayer :url="video.video_url" :danmakus="danmakus" />
      </div>

      <div class="video-info card-surface">
        <div class="video-head">
          <div>
            <h1>{{ video.title }}</h1>
            <div class="video-meta">
              <span>{{ formatCount(video.play_count) }} 播放</span>
              <span>{{ formatCount(video.danmaku_count) }} 弹幕</span>
              <span>{{ video.created_at }}</span>
              <a-tag color="pink">{{ video.category }}</a-tag>
            </div>
          </div>
          <a-space wrap>
            <a-button :type="liked ? 'primary' : 'default'" @click="toggleLike">
              <template #icon><LikeOutlined /></template>
              点赞 {{ formatCount(video.like_count) }}
            </a-button>
            <a-button @click="coinDialog = true">
              <template #icon><GiftOutlined /></template>
              投币 {{ formatCount(video.coin_count) }}
            </a-button>
            <a-button @click="favDialog = true">
              <template #icon><StarOutlined /></template>
              收藏 {{ formatCount(video.fav_count) }}
            </a-button>
            <a-button @click="doShare">
              <template #icon><ShareAltOutlined /></template>
              分享 {{ formatCount(video.share_count) }}
            </a-button>
          </a-space>
        </div>

        <div class="author-panel" @click="$router.push(`/user/${video.user_id}`)">
          <a-avatar :size="54" :src="video.author?.avatar">{{ video.author?.username?.[0] }}</a-avatar>
          <div class="author-main">
            <strong>{{ video.author?.username }}</strong>
            <p>{{ video.description || '这个 UP 还没有填写简介。' }}</p>
          </div>
          <a-button type="primary">关注</a-button>
        </div>
      </div>

      <CommentSection :video-id="video.id" />
    </section>

    <aside class="side-column">
      <section class="playlist card-surface">
        <div class="side-head">
          <h3>相关推荐</h3>
          <span>同分区与同标签混合推荐</span>
        </div>
        <div class="related-list">
          <VideoCard v-for="v in related" :key="v.id" :data="v" />
        </div>
      </section>
    </aside>

    <a-modal v-model:open="coinDialog" title="给这个视频投币" @ok="doCoin" ok-text="确认投币" cancel-text="取消">
      <a-radio-group v-model:value="coinCount" button-style="solid">
        <a-radio-button :value="1">1 枚</a-radio-button>
        <a-radio-button :value="2">2 枚</a-radio-button>
      </a-radio-group>
    </a-modal>

    <a-modal v-model:open="favDialog" title="收藏到" :footer="null">
      <div class="folder-list">
        <div v-for="f in folders" :key="f.id" class="folder-item">
          <span>{{ f.name }}</span>
          <a-button type="primary" ghost @click="doFav(f.id)">收藏</a-button>
        </div>
      </div>
      <a-input
        v-model:value="newFolderName"
        placeholder="创建新的收藏夹"
        class="new-folder-input"
        @pressEnter="createNewFolder"
      />
      <a-button block @click="createNewFolder">创建收藏夹</a-button>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { GiftOutlined, LikeOutlined, ShareAltOutlined, StarOutlined } from '@ant-design/icons-vue'
import * as api from '@/api'
import ArtPlayerPlayer from '@/components/ArtPlayerPlayer.vue'
import CommentSection from '@/components/CommentSection.vue'
import VideoCard from '@/components/VideoCard.vue'

const route = useRoute()
const video = ref(null)
const liked = ref(false)
const related = ref([])
const danmakus = ref([])
const coinDialog = ref(false)
const coinCount = ref(1)
const favDialog = ref(false)
const folders = ref([])
const newFolderName = ref('')

onMounted(async () => {
  const id = route.params.id
  const res = await api.getVideoDetail(id)
  video.value = res.data
  await Promise.all([loadFolders(), loadRelated(), loadDanmakus()])
  const l = await api.likeStatus(id)
  liked.value = l.data.type === 1
})

async function loadFolders() {
  const r = await api.getFolders()
  folders.value = r.data.list
}

async function loadRelated() {
  const r = await api.getRelatedVideos(route.params.id, 1)
  related.value = r.data.list
}

async function loadDanmakus() {
  const r = await api.getDanmakus(route.params.id, 0, 120)
  danmakus.value = r.data
}

async function toggleLike() {
  if (liked.value) {
    await api.cancelLike(video.value.id)
    video.value.like_count--
    liked.value = false
  } else {
    await api.likeVideo(video.value.id)
    video.value.like_count++
    liked.value = true
  }
}

async function doCoin() {
  await api.giveCoin(video.value.id, coinCount.value)
  video.value.coin_count += coinCount.value
  coinDialog.value = false
  message.success('投币成功')
}

async function doFav(folderId) {
  await api.addToFolder({ folder_id: folderId, video_id: video.value.id })
  video.value.fav_count++
  favDialog.value = false
  message.success('收藏成功')
}

async function createNewFolder() {
  if (!newFolderName.value) return
  await api.createFolder({ name: newFolderName.value })
  newFolderName.value = ''
  await loadFolders()
  message.success('收藏夹已创建')
}

async function doShare() {
  const r = await api.shareVideo(video.value.id)
  video.value.share_count++
  await navigator.clipboard.writeText(location.origin + r.data.link)
  message.success('分享链接已复制')
}

function formatCount(value) {
  if (!value) return '0'
  if (value >= 10000) return `${(value / 10000).toFixed(1)}万`
  return String(value)
}
</script>

<style scoped>
.detail-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 24px;
}

.main-column {
  min-width: 0;
}

.player-shell {
  overflow: hidden;
  padding: 0;
  background: #0f172a;
}

.video-info {
  margin-top: 18px;
  padding: 24px;
}

.video-head {
  display: flex;
  justify-content: space-between;
  gap: 18px;
}

.video-head h1 {
  font-size: 30px;
  line-height: 1.2;
}

.video-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 14px;
  color: var(--text-secondary);
}

.author-panel {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  margin-top: 22px;
  border-radius: 24px;
  background: linear-gradient(180deg, #fefbff, #f8fbff);
  cursor: pointer;
}

.author-main {
  flex: 1;
}

.author-main strong {
  font-size: 18px;
}

.author-main p {
  margin-top: 8px;
  color: #64748b;
  line-height: 1.7;
}

.playlist {
  padding: 20px;
}

.side-head h3 {
  font-size: 20px;
}

.side-head span {
  display: block;
  margin-top: 8px;
  color: var(--text-secondary);
  font-size: 13px;
}

.related-list {
  display: grid;
  gap: 14px;
  margin-top: 18px;
}

.folder-list {
  display: grid;
  gap: 10px;
  margin-bottom: 16px;
}

.folder-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 14px;
  border-radius: 16px;
  background: #f8fafc;
}

.new-folder-input {
  margin-bottom: 12px;
}

@media (max-width: 1200px) {
  .detail-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .video-head {
    flex-direction: column;
  }

  .author-panel {
    align-items: start;
    flex-wrap: wrap;
  }
}
</style>
