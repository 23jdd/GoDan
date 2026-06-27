<template>
  <div class="login-page">
    <el-card style="width:400px">
      <h2 style="text-align:center;margin-bottom:24px">登录 GoDan</h2>
      <el-tabs v-model="mode" @tab-click="clearMsg">
        <el-tab-pane label="密码登录" name="pwd">
          <el-input v-model="account" placeholder="邮箱 / 手机号" style="margin-bottom:12px" />
          <el-input v-model="password" type="password" placeholder="密码" show-password @keyup.enter="loginPwd" />
          <el-button type="primary" style="width:100%;margin-top:16px" @click="loginPwd">登录</el-button>
        </el-tab-pane>
        <el-tab-pane label="验证码登录" name="code">
          <el-input v-model="account" placeholder="邮箱 / 手机号" style="margin-bottom:12px" />
          <div style="display:flex;gap:8px">
            <el-input v-model="code" placeholder="验证码" style="flex:1" />
            <el-button :disabled="sendDisabled" @click="sendVerificationCode">{{ sendText }}</el-button>
          </div>
          <el-button type="primary" style="width:100%;margin-top:16px" @click="loginByCode">登录</el-button>
        </el-tab-pane>
        <el-tab-pane label="注册" name="register">
          <el-input v-model="username" placeholder="用户名" style="margin-bottom:12px" />
          <el-input v-model="account" placeholder="邮箱" style="margin-bottom:12px" />
          <el-input v-model="password" type="password" placeholder="密码" show-password style="margin-bottom:12px" />
          <div style="display:flex;gap:8px">
            <el-input v-model="code" placeholder="验证码" style="flex:1" />
            <el-button :disabled="sendDisabled" @click="sendVerificationCode">{{ sendText }}</el-button>
          </div>
          <el-button type="primary" style="width:100%;margin-top:16px" @click="doRegister">注册</el-button>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import * as api from '@/api'
import { ElMessage } from 'element-plus'

const router = useRouter()
const mode = ref('pwd')
const account = ref('')
const password = ref('')
const username = ref('')
const code = ref('')
const sendText = ref('发送验证码')
const sendDisabled = ref(false)

function clearMsg() {}

async function loginPwd() {
  try {
    const res = await api.login({ account: account.value, password: password.value })
    saveLogin(res.data)
  } catch {}
}

async function loginByCode() {
  try {
    const res = await api.loginCode({ account: account.value, code: code.value })
    saveLogin(res.data)
  } catch {}
}

async function doRegister() {
  try {
    await api.registerCode({ username: username.value, email: account.value, password: password.value, code: code.value })
    ElMessage.success('注册成功，请登录')
    mode.value = 'pwd'
  } catch {}
}

function saveLogin(data) {
  localStorage.setItem('access_token', data.access_token)
  localStorage.setItem('refresh_token', data.refresh_token)
  localStorage.setItem('user_id', data.user_id)
  localStorage.setItem('username', data.username)
  ElMessage.success('登录成功')
  router.push('/')
  setTimeout(() => location.reload(), 200)
}

async function sendVerificationCode() {
  if (!account.value) return ElMessage.warning('请输入邮箱')
  sendDisabled.value = true
  try {
    await api.sendCode(account.value)
    ElMessage.success('验证码已发送')
    let s = 60
    sendText.value = `${s}s`
    const t = setInterval(() => { s--; sendText.value = `${s}s`; if (s <= 0) { clearInterval(t); sendText.value = '发送验证码'; sendDisabled.value = false } }, 1000)
  } catch { sendDisabled.value = false }
}
</script>

<style scoped>
.login-page { display: flex; justify-content: center; align-items: center; min-height: calc(100vh - 120px); }
</style>
