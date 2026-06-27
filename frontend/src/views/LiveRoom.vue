<template>
  <div v-if="room">
    <div class="room-header">
      <div>
        <h2>{{ room.title }}</h2>
        <span class="room-author">{{ room.username }}</span>
      </div>
      <a-tag color="red">{{ room.viewer_count || 0 }} 人观看</a-tag>
    </div>

    <div class="room-body">
      <div class="player-area card-surface">
        <div class="player-placeholder">
          <VideoCameraOutlined class="player-icon" />
          <p>直播播放区域</p>
          <p class="tip">这里预留给 SRS / RTMP / WebRTC 接入。</p>
        </div>
      </div>

      <div class="chat-area card-surface">
        <h4>直播互动</h4>
        <div class="chat-messages" ref="chatBox">
          <div v-for="(m, i) in messages" :key="i" class="chat-msg">
            <span class="chat-user">{{ m.username || 'User' + m.user_id }}</span>
            <span>{{ m.content }}</span>
          </div>
        </div>
        <div class="chat-input">
          <a-input v-model:value="chatText" placeholder="发送一条弹幕..." @pressEnter="sendChat" />
          <a-button type="primary" @click="sendChat">发送</a-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { VideoCameraOutlined } from '@ant-design/icons-vue'
import * as api from '@/api'

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
  messages.value = [
    { user_id: 1, username: '系统提示', content: '欢迎来到直播间。' },
    { user_id: 2, username: '薄荷汽水', content: '这个布局很清爽！' },
  ]
})

function connectWS(roomId) {
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
  ws = new WebSocket(`${proto}//${location.host}/api/v1/live/ws?room_id=${roomId}`)
  ws.onmessage = (e) => {
    try {
      messages.value.push(JSON.parse(e.data))
    } catch {}
    setTimeout(() => {
      if (chatBox.value) chatBox.value.scrollTop = chatBox.value.scrollHeight
    }, 50)
  }
  ws.onerror = () => {}
}

function sendChat() {
  if (!chatText.value.trim()) return
  if (ws && ws.readyState === 1) {
    ws.send(JSON.stringify({
      user_id: Number(localStorage.getItem('user_id') || 0),
      username: localStorage.getItem('username') || 'GoDan 用户',
      content: chatText.value,
      color: '#FFFFFF',
    }))
  }
  messages.value.push({
    user_id: Number(localStorage.getItem('user_id') || 0),
    username: localStorage.getItem('username') || 'GoDan 用户',
    content: chatText.value,
  })
  chatText.value = ''
}

onUnmounted(() => {
  if (ws) ws.close()
})
</script>

<style scoped>
.room-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 18px;
  padding: 16px 0;
}

.room-header h2 {
  font-size: 30px;
}

.room-author {
  color: var(--text-secondary);
  font-size: 14px;
  display: inline-block;
  margin-top: 10px;
}

.room-body {
  display: flex;
  gap: 16px;
}

.player-area {
  flex: 1;
}

.player-placeholder {
  aspect-ratio: 16/9;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #fff;
  background: linear-gradient(135deg, #0f172a, #1e293b);
}

.player-icon {
  font-size: 64px;
}

.player-placeholder .tip {
  font-size: 12px;
  color: #cbd5e1;
  margin-top: 8px;
}

.chat-area {
  width: 320px;
  display: flex;
  flex-direction: column;
  height: 500px;
  padding: 12px;
}

.chat-area h4 {
  padding: 12px;
  border-bottom: 1px solid var(--border);
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.chat-msg {
  font-size: 13px;
  padding: 8px 0;
  display: flex;
  gap: 8px;
}

.chat-user {
  color: var(--primary);
}

.chat-input {
  padding: 8px;
  border-top: 1px solid var(--border);
  display: flex;
  gap: 6px;
}

@media (max-width: 1100px) {
  .room-body {
    flex-direction: column;
  }

  .chat-area {
    width: 100%;
  }
}
</style>
