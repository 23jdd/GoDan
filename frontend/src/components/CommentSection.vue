<template>
  <section class="comment-section card-surface">
    <div class="comment-head">
      <div>
        <h3>评论区</h3>
        <p>{{ total }} 条热评与讨论</p>
      </div>
      <a-segmented v-model:value="sortMode" :options="sortOptions" />
    </div>

    <div class="comment-input">
      <a-textarea v-model:value="content" :rows="3" placeholder="发一条友善的评论，分享你的看法吧。" />
      <div class="comment-submit">
        <span>支持一级评论与楼中楼回复</span>
        <a-button type="primary" @click="submit">发布评论</a-button>
      </div>
    </div>

    <div v-for="c in list" :key="c.id" class="comment-item">
      <a-avatar :size="42" :src="c.avatar">{{ c.username?.[0] }}</a-avatar>
      <div class="comment-body">
        <div class="comment-user">{{ c.username }}</div>
        <div class="comment-text">{{ c.content }}</div>
        <div class="comment-actions">
          <span>{{ c.created_at?.slice(0, 10) }}</span>
          <button class="action-btn" @click="doLike(c)">
            <LikeOutlined />
            {{ c.like_count || 0 }}
          </button>
          <button class="action-btn" @click="openReply(c)">回复</button>
        </div>

        <div v-if="c.reply_count > 0" class="reply-box">
          <button class="reply-toggle" @click="loadReplies(c)">
            {{ showReplies[c.id] ? '收起回复' : `展开 ${c.reply_count} 条回复` }}
          </button>
          <div v-if="showReplies[c.id]" class="reply-list">
            <div v-for="r in replies[c.id]" :key="r.id" class="reply-item">
              <span class="reply-user">{{ r.username }}</span>
              <span>{{ r.content }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <a-modal v-model:open="replyVisible" title="回复评论" @ok="submitReply" ok-text="发送回复" cancel-text="取消">
      <a-textarea v-model:value="replyContent" :rows="4" placeholder="写下你的回复..." />
    </a-modal>
  </section>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { LikeOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import * as api from '@/api'

const props = defineProps({ videoId: [String, Number] })
const list = ref([])
const total = ref(0)
const content = ref('')
const replies = ref({})
const showReplies = ref({})
const replyVisible = ref(false)
const replyTo = ref(null)
const replyContent = ref('')
const sortMode = ref('hot')
const sortOptions = [
  { label: '最热', value: 'hot' },
  { label: '最新', value: 'latest' },
]

onMounted(loadComments)
watch(sortMode, loadComments)

async function loadComments() {
  const res = await api.getComments(props.videoId, sortMode.value, 1)
  list.value = res.data.list
  total.value = res.data.total
}

async function submit() {
  if (!content.value.trim()) return
  await api.postComment({ video_id: props.videoId, content: content.value })
  message.success('评论已发布')
  content.value = ''
  loadComments()
}

async function loadReplies(c) {
  if (showReplies.value[c.id]) {
    showReplies.value[c.id] = false
    return
  }
  const res = await api.getReplies(c.id, 1)
  replies.value[c.id] = res.data.list
  showReplies.value[c.id] = true
}

async function doLike(c) {
  await api.likeComment(c.id)
  c.like_count++
}

function openReply(c) {
  replyTo.value = c
  replyVisible.value = true
}

async function submitReply() {
  if (!replyContent.value.trim() || !replyTo.value) return
  await api.postComment({
    video_id: props.videoId,
    content: replyContent.value,
    parent_id: replyTo.value.id,
    root_id: replyTo.value.root_id || replyTo.value.id,
  })
  replyVisible.value = false
  replyContent.value = ''
  message.success('回复已发送')
  loadComments()
}
</script>

<style scoped>
.comment-section {
  margin-top: 28px;
  padding: 24px;
}

.comment-head {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
}

.comment-head h3 {
  font-size: 22px;
}

.comment-head p {
  margin-top: 6px;
  color: var(--text-secondary);
}

.comment-input {
  margin: 20px 0 28px;
}

.comment-submit {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
  color: var(--text-muted);
  font-size: 13px;
}

.comment-item {
  display: flex;
  gap: 14px;
  padding: 18px 0;
  border-top: 1px solid rgba(15, 23, 42, 0.06);
}

.comment-body {
  flex: 1;
}

.comment-user {
  font-weight: 700;
}

.comment-text {
  margin-top: 8px;
  line-height: 1.75;
  color: #334155;
}

.comment-actions {
  display: flex;
  gap: 16px;
  margin-top: 12px;
  color: var(--text-secondary);
  font-size: 13px;
}

.action-btn,
.reply-toggle {
  border: none;
  background: transparent;
  color: inherit;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.reply-box {
  margin-top: 12px;
}

.reply-list {
  margin-top: 10px;
  padding: 14px;
  border-radius: 18px;
  background: #f8fafc;
}

.reply-item + .reply-item {
  margin-top: 10px;
}

.reply-user {
  font-weight: 600;
  color: #fb7299;
  margin-right: 8px;
}
</style>
