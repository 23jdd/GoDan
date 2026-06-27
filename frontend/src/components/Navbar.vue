<template>
  <header class="navbar">
    <div class="nav-inner">
      <div class="nav-left">
        <router-link to="/" class="logo">GoDan</router-link>
        <nav class="nav-links">
          <router-link v-for="item in links" :key="item.to" :to="item.to" class="nav-link">
            {{ item.label }}
          </router-link>
        </nav>
      </div>

      <div class="nav-center">
        <a-input
          v-model:value="keyword"
          class="search-input"
          size="large"
          placeholder="搜索视频、番剧、UP 主"
          @pressEnter="doSearch"
        >
          <template #prefix>
            <SearchOutlined />
          </template>
        </a-input>
      </div>

      <div class="nav-right">
        <a-button type="primary" size="large" class="publish-btn" @click="$router.push('/upload')">
          <template #icon>
            <PlusOutlined />
          </template>
          投稿
        </a-button>

        <a-dropdown v-if="isLogin" placement="bottomRight">
          <div class="user-chip">
            <a-avatar :size="42" :src="avatar">{{ username?.[0] }}</a-avatar>
            <div class="user-meta">
              <strong>{{ username }}</strong>
              <span>创作中心</span>
            </div>
          </div>
          <template #overlay>
            <a-menu>
              <a-menu-item key="profile" @click="$router.push(`/user/${userId}`)">个人中心</a-menu-item>
              <a-menu-item key="timeline" @click="$router.push('/timeline')">我的动态</a-menu-item>
              <a-menu-item key="logout" @click="logout">退出登录</a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>

        <a-button v-else size="large" @click="$router.push('/login')">登录</a-button>
      </div>
    </div>
  </header>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { PlusOutlined, SearchOutlined } from '@ant-design/icons-vue'

const router = useRouter()
const keyword = ref('')
const username = ref('')
const avatar = ref('')
const userId = ref(0)

const links = [
  { to: '/', label: '首页' },
  { to: '/live', label: '直播' },
  { to: '/timeline', label: '动态' },
]

const isLogin = computed(() => !!localStorage.getItem('access_token'))

onMounted(() => {
  username.value = localStorage.getItem('username') || ''
  avatar.value = localStorage.getItem('avatar') || ''
  userId.value = Number(localStorage.getItem('user_id') || 0)
})

function doSearch() {
  if (keyword.value.trim()) router.push(`/search?q=${encodeURIComponent(keyword.value.trim())}`)
}

function logout() {
  localStorage.clear()
  username.value = ''
  router.push('/')
}
</script>

<style scoped>
.navbar {
  position: sticky;
  top: 0;
  z-index: 100;
  backdrop-filter: blur(18px);
  background: rgba(255, 255, 255, 0.88);
  border-bottom: 1px solid rgba(17, 24, 39, 0.06);
}

.nav-inner {
  max-width: 1440px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  gap: 24px;
  min-height: 74px;
  padding: 14px 24px;
}

.nav-left {
  display: flex;
  align-items: center;
  gap: 24px;
}

.logo {
  font-size: 30px;
  font-weight: 800;
  letter-spacing: -0.03em;
  color: #fb7299;
}

.nav-links {
  display: flex;
  gap: 10px;
}

.nav-link {
  padding: 10px 14px;
  border-radius: 999px;
  color: var(--text-secondary);
  transition: all 0.2s ease;
}

.nav-link.router-link-exact-active,
.nav-link:hover {
  color: #111827;
  background: rgba(251, 114, 153, 0.1);
}

.nav-center {
  flex: 1;
  display: flex;
  justify-content: center;
}

.search-input {
  width: min(100%, 640px);
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 14px;
}

.publish-btn {
  background: linear-gradient(135deg, #fb7299, #ff8fab);
  border: none;
  box-shadow: 0 16px 30px rgba(251, 114, 153, 0.22);
}

.user-chip {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.user-meta {
  display: flex;
  flex-direction: column;
  line-height: 1.1;
}

.user-meta strong {
  font-size: 14px;
}

.user-meta span {
  font-size: 12px;
  color: var(--text-secondary);
}

@media (max-width: 960px) {
  .nav-inner {
    flex-wrap: wrap;
    gap: 14px;
    padding: 14px 16px;
  }

  .nav-left,
  .nav-center,
  .nav-right {
    width: 100%;
  }

  .nav-center {
    order: 3;
  }

  .nav-right {
    justify-content: flex-end;
  }
}
</style>
