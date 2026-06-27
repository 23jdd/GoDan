<template>
  <header class="navbar">
    <div class="nav-inner">
      <div class="nav-left">
        <router-link to="/" class="logo">GoDan</router-link>
        <router-link to="/live" class="nav-link">直播</router-link>
        <router-link to="/timeline" class="nav-link">动态</router-link>
      </div>
      <div class="nav-center">
        <el-input v-model="keyword" placeholder="搜索视频..." class="search-input" @keyup.enter="doSearch" clearable>
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
      </div>
      <div class="nav-right">
        <el-button type="primary" size="small" @click="$router.push('/upload')">
          <el-icon><Upload /></el-icon> 投稿
        </el-button>
        <el-dropdown v-if="isLogin" trigger="click">
          <el-avatar :size="36" :src="avatar" class="avatar-cursor">{{ username?.[0] }}</el-avatar>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="$router.push(`/user/${userId}`)">个人中心</el-dropdown-item>
              <el-dropdown-item @click="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button v-else size="small" @click="$router.push('/login')">登录</el-button>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Upload } from '@element-plus/icons-vue'

const router = useRouter()
const keyword = ref('')
const username = ref('')
const avatar = ref('')
const userId = ref(0)

const isLogin = computed(() => !!localStorage.getItem('access_token'))

onMounted(() => {
  username.value = localStorage.getItem('username') || ''
  avatar.value = localStorage.getItem('avatar') || ''
  userId.value = +localStorage.getItem('user_id') || 0
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
.navbar { position: sticky; top: 0; z-index: 100; background: #fff; border-bottom: 1px solid var(--border); height: 64px; }
.nav-inner { max-width: 1400px; margin: 0 auto; display: flex; align-items: center; height: 100%; padding: 0 16px; gap: 24px; }
.nav-left { display: flex; align-items: center; gap: 20px; }
.logo { font-size: 24px; font-weight: 700; color: var(--primary); }
.nav-link { font-size: 14px; color: var(--text-secondary); }
.nav-link:hover { color: var(--primary); }
.nav-center { flex: 1; max-width: 480px; }
.nav-right { display: flex; align-items: center; gap: 12px; }
.avatar-cursor { cursor: pointer; }
</style>
