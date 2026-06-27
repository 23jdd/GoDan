<template>
  <div class="comment-section">
    <h4>评论 ({{ total }})</h4>
    <div class="comment-input">
      <el-input v-model="content" type="textarea" :rows="2" placeholder="发一条友善的评论..." />
      <el-button type="primary" size="small" @click="submit" style="margin-top:8px">发表</el-button>
    </div>
    <div v-for="c in list" :key="c.id" class="comment-item">
      <div class="comment-avatar">
        <el-avatar :size="32" :src="c.avatar">{{ c.username?.[0] }}</el-avatar>
      </div>
      <div class="comment-body">
        <div class="comment-user">{{ c.username }}</div>
        <div class="comment-text">{{ c.content }}</div>
        <div class="comment-actions">
          <span>{{ c.created_at?.slice(0,10) }}</span>
          <span class="action" @click="doLike(c)"><el-icon><CaretTop /></el-icon> {{ c.like_count || 0 }}</span>
          <span class="action" @click="replyTo = c">回复</span>
        </div>
        <div v-if="c.reply_count > 0">
          <span class="action" @click="loadReplies(c)">{{ showReplies[c.id] ? '收起' : `展开 ${c.reply_count} 条回复` }}</span>
          <div v-if="showReplies[c.id]" v-for="r in replies[c.id]" :key="r.id" class="reply-item">
            <span class="reply-user">{{ r.username }}</span>: {{ r.content }}
          </div>
        </div>
      </div>
    </div>
    <!-- Reply dialog -->
    <el-dialog v-model="replyVisible" title="回复评论" width="400px">
      <el-input v-model="replyContent" type="textarea" placeholder="回复..." />
      <template #footer>
        <el-button @click="replyVisible = false">取消</el-button>
        <el-button type="primary" @click="submitReply">回复</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { CaretTop } from '@element-plus/icons-vue'
import * as api from '@/api'
const list = ref([])
const total = ref(0)
const content = ref('')
const replies = ref({})
const showReplies = ref({})
const replyVisible = ref(false)
const replyTo = ref(null)
const replyContent = ref('')

onMounted(() => loadComments())

async function loadComments() {
  try { const res = await api.getComments(props.videoId, 'hot', 1); list.value = res.data.list; total.value = res.data.total } catch {}
}

async function submit() {
  if (!content.value.trim()) return
  try { await api.postComment({ video_id: props.videoId, content: content.value }); content.value = ''; loadComments() } catch {}
}

async function loadReplies(c) {
  if (showReplies.value[c.id]) { showReplies.value[c.id] = false; return }
  try { const res = await api.getReplies(c.id, 1); replies.value[c.id] = res.data.list; showReplies.value[c.id] = true } catch {}
}

async function doLike(c) {
  try { await api.likeComment(c.id); c.like_count++ } catch {}
}

async function submitReply() {
  if (!replyContent.value.trim() || !replyTo.value) return
  try {
    await api.postComment({
      video_id: props.videoId, content: replyContent.value,
      parent_id: replyTo.value.id, root_id: replyTo.value.root_id || replyTo.value.id
    })
    replyVisible.value = false; replyContent.value = ''; loadComments()
  } catch {}
}
</script>

<style scoped>
.comment-section { margin-top: 24px; }
.comment-input { margin: 12px 0 20px; }
.comment-item { display: flex; gap: 12px; padding: 12px 0; border-bottom: 1px solid var(--border); }
.comment-body { flex: 1; }
.comment-user { font-size: 13px; font-weight: 500; color: var(--primary); }
.comment-text { font-size: 14px; margin: 4px 0; line-height: 1.5; }
.comment-actions { display: flex; gap: 16px; font-size: 12px; color: var(--text-secondary); }
.action { cursor: pointer; display: flex; align-items: center; gap: 2px; }
.action:hover { color: var(--primary); }
.reply-item { padding: 6px 0; font-size: 13px; }
.reply-user { color: var(--primary); }
</style>
