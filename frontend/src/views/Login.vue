<template>
  <div class="login-page">
    <section class="login-hero">
      <span class="hero-pill">创作社区登录</span>
      <h1>把 GoDan 的创作、观看、互动体验串起来。</h1>
      <p>这版登录页也同步切到 Ant Design Vue，保留轻盈的 B 站气质和更统一的视觉节奏。</p>
    </section>

    <section class="login-card card-surface">
      <h2>欢迎回来</h2>
      <a-tabs v-model:activeKey="mode">
        <a-tab-pane key="pwd" tab="密码登录">
          <a-space direction="vertical" size="middle" style="width: 100%">
            <a-input v-model:value="account" size="large" placeholder="邮箱 / 手机号" />
            <a-input-password v-model:value="password" size="large" placeholder="密码" @pressEnter="loginPwd" />
            <a-button type="primary" size="large" block @click="loginPwd">登录</a-button>
          </a-space>
        </a-tab-pane>

        <a-tab-pane key="code" tab="验证码登录">
          <a-space direction="vertical" size="middle" style="width: 100%">
            <a-input v-model:value="account" size="large" placeholder="邮箱 / 手机号" />
            <div class="inline-row">
              <a-input v-model:value="code" size="large" placeholder="验证码" />
              <a-button size="large" :disabled="sendDisabled" @click="sendVerificationCode">{{ sendText }}</a-button>
            </div>
            <a-button type="primary" size="large" block @click="loginByCode">登录</a-button>
          </a-space>
        </a-tab-pane>

        <a-tab-pane key="register" tab="注册">
          <a-space direction="vertical" size="middle" style="width: 100%">
            <a-input v-model:value="username" size="large" placeholder="用户名" />
            <a-input v-model:value="account" size="large" placeholder="邮箱" />
            <a-input-password v-model:value="password" size="large" placeholder="密码" />
            <div class="inline-row">
              <a-input v-model:value="code" size="large" placeholder="验证码" />
              <a-button size="large" :disabled="sendDisabled" @click="sendVerificationCode">{{ sendText }}</a-button>
            </div>
            <a-button type="primary" size="large" block @click="doRegister">注册</a-button>
          </a-space>
        </a-tab-pane>
      </a-tabs>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import * as api from '@/api'

const router = useRouter()
const mode = ref('pwd')
const account = ref('')
const password = ref('')
const username = ref('')
const code = ref('')
const sendText = ref('发送验证码')
const sendDisabled = ref(false)

async function loginPwd() {
  const res = await api.login({ account: account.value, password: password.value })
  saveLogin(res.data)
}

async function loginByCode() {
  const res = await api.loginCode({ account: account.value, code: code.value })
  saveLogin(res.data)
}

async function doRegister() {
  await api.registerCode({ username: username.value, email: account.value, password: password.value, code: code.value })
  message.success('注册成功，请登录')
  mode.value = 'pwd'
}

function saveLogin(data) {
  localStorage.setItem('access_token', data.access_token)
  localStorage.setItem('refresh_token', data.refresh_token)
  localStorage.setItem('user_id', data.user_id)
  localStorage.setItem('username', data.username)
  if (data.avatar) localStorage.setItem('avatar', data.avatar)
  message.success('登录成功')
  router.push('/')
  setTimeout(() => location.reload(), 200)
}

async function sendVerificationCode() {
  if (!account.value) return message.warning('请先输入邮箱或手机号')
  sendDisabled.value = true
  await api.sendCode(account.value)
  message.success('验证码已发送')
  let s = 60
  sendText.value = `${s}s`
  const t = setInterval(() => {
    s--
    sendText.value = `${s}s`
    if (s <= 0) {
      clearInterval(t)
      sendText.value = '发送验证码'
      sendDisabled.value = false
    }
  }, 1000)
}
</script>

<style scoped>
.login-page {
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) minmax(340px, 460px);
  gap: 32px;
  align-items: center;
  min-height: calc(100vh - 150px);
}

.login-hero {
  padding: 24px;
}

.hero-pill {
  display: inline-flex;
  padding: 8px 14px;
  border-radius: 999px;
  background: rgba(71, 187, 255, 0.12);
  color: #0077b6;
  font-weight: 700;
}

.login-hero h1 {
  margin-top: 20px;
  font-size: clamp(34px, 4vw, 56px);
  line-height: 1.05;
  letter-spacing: -0.04em;
}

.login-hero p {
  max-width: 600px;
  margin-top: 20px;
  color: #64748b;
  line-height: 1.85;
  font-size: 16px;
}

.login-card {
  padding: 26px;
}

.login-card h2 {
  margin-bottom: 18px;
  font-size: 28px;
}

.inline-row {
  display: grid;
  grid-template-columns: 1fr 120px;
  gap: 10px;
}

@media (max-width: 960px) {
  .login-page {
    grid-template-columns: 1fr;
  }

  .login-hero {
    padding: 0;
  }
}
</style>
