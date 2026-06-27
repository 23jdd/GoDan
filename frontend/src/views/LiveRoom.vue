<template>
  <div v-if="room">
    <div class="room-header">
      <h2>{{ room.title }}</h2>
      <span class="room-author">{{ room.username }}</span>
      <span class="room-viewers">{{ room.viewer_count || 0 }} 人观看</span>
    </div>
    <div class="room-body">
      <div class="player-area">
        <div class="player-placeholder">
          <el-icon size="64"><VideoPlay /></el-icon>
          <p>直播播放区域</p>
          <p class="tip">(需接入 SRS/RTMP 流媒体服务)</p>
        </div>
      </div>
      <div class="chat-area">
        <h4>直播弹幕</h4>
        <div class="chat-messages" ref="chatBox">
          <div v-for="(m, i) in messages" :key="i" class="chat-msg">
            <span class="chat-user">{{ m.username || 'User' + m.user_id }}</span>: {{ m.content }}
          </div>
        </div>
        <div class="chat-input">
          <el-input v-model="chatText" placeholder="发送弹幕..." @keyup.enter="sendChat" />
          <el-button size="small" @click="sendChat">发送</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import * as api from '@/api'
import { VideoPlay } from '@element-plus/icons-vue'

const route = useRoute()
const room = ref(null)
const messages = ref([])
const chatText = ref('')
const chatBox = ref(null)
let ws = null

onMounted(async () => {
  const id = route.params.id
  room.value = (await api.getRoomInfo(id)).data
  connectWS(id)
})

function connectWS(roomId) {
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
  ws = new WebSocket(`${proto}//${location.host}/api/v1/live/ws?room_id=${roomId}`)
  ws.onmessage = e => {
    try { messages.value.push(JSON.parse(e.data)) } catch {}
    setTimeout(() => { if (chatBox.value) chatBox.value.scrollTop = chatBox.value.scrollHeight }, 50)
  }
}

function sendChat() {
  if (!chatText.value.trim() || !ws) return
  ws.send(JSON.stringify({
    user_id: +localStorage.getItem('user_id') || 0,
    username: localStorage.getItem('username') || '用户',
    content: chatText.value,
    color: '#FFFFFF'
  }))
  chatText.value = ''
}

onUnmounted(() => { if (ws) ws.close() })
</script>

<style scoped>
.room-header { padding: 16px 0; }
.room-author { color: var(--text-secondary); font-size: 14px; margin: 0 16px; }
.room-viewers { color: #ff4d4f; font-size: 13px; }
.room-body { display: flex; gap: 16px; }
.player-area { flex: 1; }
.player-placeholder { background: #000; aspect-ratio: 16/9; border-radius: 8px; display: flex; flex-direction: column; align-items: center; justify-content: center; color: #fff; }
.player-placeholder .tip { font-size: 12px; color: #999; margin-top: 8px; }
.chat-area { width: 320px; background: #fff; border-radius: 8px; display: flex; flex-direction: column; height: 500px; }
.chat-area h4 { padding: 12px; border-bottom: 1px solid var(--border); }
.chat-messages { flex: 1; overflow-y: auto; padding: 8px; }
.chat-msg { font-size: 13px; padding: 2px 0; }
.chat-user { color: var(--primary); }
.chat-input { padding: 8px; border-top: 1px solid var(--border); display: flex; gap: 6px; }
</style>
