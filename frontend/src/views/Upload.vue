<template>
  <div>
    <div class="page-title">视频投稿</div>
    <el-card style="max-width:700px;margin:0 auto">
      <!-- Select file -->
      <div v-if="!uploadId && !videoSubmitted" class="upload-step">
        <h3>选择视频文件</h3>
        <el-upload drag :auto-upload="false" :on-change="onFileSelected" accept="video/*">
          <el-icon size="48"><UploadFilled /></el-icon>
          <div>拖拽或点击选择视频文件</div>
          <div class="tip">支持 mp4, avi, mov, mkv, flv, webm</div>
        </el-upload>
        <div v-if="file" style="margin-top:16px">
          <p>文件: {{ file.name }} ({{ (file.size / 1024 / 1024).toFixed(1) }} MB)</p>
          <el-button type="primary" @click="initUploadAction" :loading="loading">开始上传</el-button>
        </div>
      </div>

      <!-- Upload progress -->
      <div v-if="uploadId && !videoSubmitted" class="upload-step">
        <h3>正在上传...</h3>
        <el-progress :percentage="progress" :text-inside="true" :stroke-width="20" />
        <p style="text-align:center;margin-top:8px;color:#999">分片 {{ uploadedChunks }}/{{ chunkCount }}</p>
      </div>

      <!-- Complete -->
      <div v-if="videoSubmitted" class="upload-step">
        <el-result icon="success" title="投稿成功" sub-title="视频已提交，审核通过后即可展示">
          <template #extra>
            <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
          </template>
        </el-result>
      </div>

      <!-- Metadata form (after upload) -->
      <div v-if="readyToComplete" class="upload-step">
        <h3>填写视频信息</h3>
        <el-input v-model="title" placeholder="视频标题" maxlength="200" style="margin-bottom:12px" />
        <el-input v-model="description" type="textarea" :rows="3" placeholder="视频简介" style="margin-bottom:12px" />
        <el-button type="primary" @click="completeUploadAction" :loading="loading">提交发布</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import * as api from '@/api'
import { UploadFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const CHUNK_SIZE = 5 << 20
const file = ref(null)
const uploadId = ref('')
const chunkCount = ref(0)
const uploadedChunks = ref(0)
const progress = ref(0)
const readyToComplete = ref(false)
const videoSubmitted = ref(false)
const loading = ref(false)
const title = ref('')
const description = ref('')

function onFileSelected(f) { file.value = f.raw }

async function initUploadAction() {
  loading.value = true
  const res = await api.initUpload({ filename: file.value.name, file_size: file.value.size })
  uploadId.value = res.data.upload_id
  chunkCount.value = res.data.chunk_count
  loading.value = false
  await uploadAllChunks()
}

async function uploadAllChunks() {
  for (let i = 0; i < chunkCount.value; i++) {
    const start = i * CHUNK_SIZE
    const end = Math.min(start + CHUNK_SIZE, file.value.size)
    const chunk = file.value.slice(start, end)
    const fd = new FormData()
    fd.append('upload_id', uploadId.value)
    fd.append('chunk_index', i)
    fd.append('file', chunk)
    await api.uploadChunk(fd)
    uploadedChunks.value = i + 1
    progress.value = Math.round(((i + 1) / chunkCount.value) * 100)
  }
  readyToComplete.value = true
}

async function completeUploadAction() {
  if (!title.value) return ElMessage.warning('请输入标题')
  loading.value = true
  await api.completeUpload({ upload_id: uploadId.value, title: title.value, description: description.value })
  await api.publishVideo(uploadId.value) // try auto-publish (admin review in prod)
  videoSubmitted.value = true
}
</script>

<style scoped>
.upload-step h3 { margin-bottom: 16px; }
.tip { font-size: 12px; color: #ccc; }
</style>
