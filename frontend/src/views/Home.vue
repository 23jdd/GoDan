<template>
  <div>
    <div class="page-title">热门推荐</div>
    <div class="video-grid">
      <VideoCard v-for="v in videos" :key="v.id" :data="v" />
    </div>
    <div class="flex-center" style="margin-top:24px">
      <el-pagination
        v-model:current-page="page"
        :total="total"
        :page-size="12"
        layout="prev, pager, next"
        @current-change="load"
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

onMounted(() => load())

async function load() {
  const res = await api.getHotVideos(page.value)
  videos.value = res.data.list
  total.value = res.data.total
}
</script>
