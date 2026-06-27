<template>
  <div class="upload-page">
    <section class="page-header">
      <h1>视频投稿</h1>
      <p>延续 B 站投稿流程的熟悉感，同时保留 GoDan 自己更轻快的视觉表达。</p>
    </section>

    <section class="upload-shell card-surface">
      <div v-if="!uploadId && !videoSubmitted" class="upload-step">
        <a-upload-dragger :before-upload="beforeUpload" :show-upload-list="false" accept="video/*">
          <p class="ant-upload-drag-icon">
            <InboxOutlined />
          </p>
          <p class="upload-title">拖拽视频到这里，或点击选择文件</p>
          <p class="upload-tip">支持 mp4、mov、mkv、webm，当前演示支持分片上传流程。</p>
        </a-upload-dragger>

        <div v-if="file" class="file-summary">
          <strong>{{ file.name }}</strong>
          <span>{{ (file.size / 1024 / 1024).toFixed(1) }} MB</span>
          <a-button type="primary" :loading="loading" @click="initUploadAction">开始上传</a-button>
        </div>
      </div>

      <div v-if="uploadId && !videoSubmitted && !readyToComplete" class="upload-step">
        <h3>正在上传视频</h3>
        <a-progress :percent="progress" />
        <p class="progress-text">已上传分片 {{ uploadedChunks }}/{{ chunkCount }}</p>
      </div>

      <div v-if="readyToComplete && !videoSubmitted" class="upload-step meta-form">
        <h3>填写视频信息</h3>
        <a-input v-model:value="title" size="large" placeholder="视频标题" />
        <a-textarea v-model:value="description" :rows="4" placeholder="简介、亮点、分区说明..." />
        <a-button type="primary" size="large" :loading="loading" @click="completeUploadAction">提交发布</a-button>
      </div>

      <a-result
        v-if="videoSubmitted"
        status="success"
        title="投稿成功"
        sub-title="视频信息已经提交，后续接入审核后即可走真实发布链路。"
      >
        <template #extra>
          <a-button type="primary" @click="$router.push('/')">返回首页</a-button>
        </template>
      </a-result>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { InboxOutlined } from '@ant-design/icons-vue'
import * as api from '@/api'

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

function beforeUpload(rawFile) {
  file.value = rawFile
  return false
}

async function initUploadAction() {
  if (!file.value) return
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
  if (!title.value) return message.warning('请输入标题')
  loading.value = true
  await api.completeUpload({ upload_id: uploadId.value, title: title.value, description: description.value })
  await api.publishVideo(uploadId.value)
  videoSubmitted.value = true
  loading.value = false
}
</script>

<style scoped>
.upload-page {
  max-width: 960px;
  margin: 0 auto;
}

.page-header h1 {
  font-size: 32px;
}

.page-header p {
  margin-top: 10px;
  color: var(--text-secondary);
}

.upload-shell {
  margin-top: 24px;
  padding: 24px;
}

.upload-step h3 {
  margin-bottom: 16px;
  font-size: 24px;
}

.upload-title {
  font-size: 18px;
  font-weight: 700;
}

.upload-tip,
.progress-text {
  margin-top: 10px;
  color: var(--text-secondary);
}

.file-summary {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 18px;
}

.meta-form {
  display: grid;
  gap: 14px;
}
</style>
