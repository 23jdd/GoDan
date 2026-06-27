import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'Home', component: () => import('@/views/Home.vue') },
  { path: '/login', name: 'Login', component: () => import('@/views/Login.vue') },
  { path: '/video/:id', name: 'VideoDetail', component: () => import('@/views/VideoDetail.vue') },
  { path: '/upload', name: 'Upload', component: () => import('@/views/Upload.vue'), meta: { requiresAuth: true } },
  { path: '/user/:id', name: 'UserProfile', component: () => import('@/views/UserProfile.vue') },
  { path: '/search', name: 'Search', component: () => import('@/views/Search.vue') },
  { path: '/live', name: 'LiveList', component: () => import('@/views/LiveList.vue') },
  { path: '/live/:id', name: 'LiveRoom', component: () => import('@/views/LiveRoom.vue') },
  { path: '/timeline', name: 'Timeline', component: () => import('@/views/Timeline.vue'), meta: { requiresAuth: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('access_token')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
