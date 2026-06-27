<template>
  <div class="home-page">
    <section class="hero-banner">
      <div class="hero-copy">
        <span class="hero-tag">Phase 12</span>
        <h1>GoDan 前端已切到 Vue3 + Vite + Ant Design Vue + ArtPlayer</h1>
        <p>
          以 B 站式信息流为视觉基底，重做了首页、视频详情、评论区、搜索、投稿、直播和用户主页。
        </p>
        <div class="hero-actions">
          <a-button type="primary" size="large" @click="$router.push('/video/1')">查看演示详情</a-button>
          <a-button size="large" @click="$router.push('/upload')">进入投稿页</a-button>
        </div>
      </div>
      <div class="hero-panel card-surface">
        <div class="hero-panel-grid">
          <div>
            <strong>12</strong>
            <span>个演示视频</span>
          </div>
          <div>
            <strong>3</strong>
            <span>个直播房间</span>
          </div>
          <div>
            <strong>Antd</strong>
            <span>统一组件基座</span>
          </div>
          <div>
            <strong>ArtPlayer</strong>
            <span>播放器体验升级</span>
          </div>
        </div>
      </div>
    </section>

    <section class="section-head">
      <div>
        <h2>热门推荐</h2>
        <p>延续 B 站首页的轻快节奏，但保留 GoDan 自己的品牌感。</p>
      </div>
      <a-segmented v-model:value="feedMode" :options="feedOptions" />
    </section>

    <div class="video-grid">
      <VideoCard v-for="v in videos" :key="v.id" :data="v" />
    </div>

    <div class="pagination-wrap">
      <a-pagination
        v-model:current="page"
        :total="total"
        :page-size="12"
        show-less-items
        @change="load"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import * as api from '@/api'
import VideoCard from '@/components/VideoCard.vue'

const videos = ref([])
const page = ref(1)
const total = ref(0)
const feedMode = ref('hot')
const feedOptions = [
  { label: '综合', value: 'hot' },
  { label: '最新', value: 'latest' },
  { label: '动画', value: 'anime' },
]

onMounted(load)

async function load() {
  const res = await api.getHotVideos(page.value)
  videos.value = res.data.list
  total.value = res.data.total
}
</script>

<style scoped>
.home-page {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.hero-banner {
  display: grid;
  grid-template-columns: minmax(0, 1.45fr) minmax(320px, 0.9fr);
  gap: 24px;
  padding: 28px;
  border-radius: 32px;
  background:
    radial-gradient(circle at top left, rgba(251, 114, 153, 0.26), transparent 32%),
    radial-gradient(circle at right center, rgba(71, 187, 255, 0.2), transparent 28%),
    linear-gradient(135deg, #fff6fb 0%, #ffffff 42%, #f3fbff 100%);
  border: 1px solid rgba(251, 114, 153, 0.12);
}

.hero-copy h1 {
  margin-top: 14px;
  font-size: clamp(34px, 4vw, 54px);
  line-height: 1.05;
  letter-spacing: -0.04em;
}

.hero-copy p {
  max-width: 680px;
  margin-top: 18px;
  color: #475569;
  font-size: 17px;
  line-height: 1.8;
}

.hero-tag {
  display: inline-flex;
  padding: 8px 14px;
  border-radius: 999px;
  background: rgba(251, 114, 153, 0.12);
  color: #d83a72;
  font-weight: 700;
}

.hero-actions {
  display: flex;
  gap: 12px;
  margin-top: 28px;
}

.hero-panel {
  padding: 22px;
}

.hero-panel-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
  height: 100%;
}

.hero-panel-grid > div {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  min-height: 140px;
  padding: 18px;
  border-radius: 24px;
  background: linear-gradient(180deg, #ffffff, #f8fbff);
  border: 1px solid rgba(148, 163, 184, 0.14);
}

.hero-panel-grid strong {
  font-size: 26px;
}

.hero-panel-grid span {
  margin-top: 8px;
  color: var(--text-secondary);
}

.section-head {
  display: flex;
  justify-content: space-between;
  align-items: end;
  gap: 16px;
}

.section-head h2 {
  font-size: 28px;
}

.section-head p {
  margin-top: 8px;
  color: var(--text-secondary);
}

.pagination-wrap {
  display: flex;
  justify-content: center;
  margin-top: 8px;
}

@media (max-width: 1024px) {
  .hero-banner {
    grid-template-columns: 1fr;
  }

  .section-head {
    flex-direction: column;
    align-items: start;
  }
}

@media (max-width: 640px) {
  .hero-banner {
    padding: 20px;
    border-radius: 24px;
  }

  .hero-actions {
    flex-direction: column;
  }
}
</style>
