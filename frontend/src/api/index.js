import axios from 'axios'
import { message } from 'ant-design-vue'
import router from '@/router'
import {
  getUserProfile,
  getVideoById,
  makePage,
  mockComments,
  mockFolders,
  mockLiveRooms,
  mockReplies,
  mockTimeline,
  mockVideos,
  searchMockVideos,
} from '@/mock/data'

const api = axios.create({ baseURL: '/api/v1', timeout: 15000 })

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  (res) => res.data,
  (err) => {
    const msg = err.response?.data?.message || '请求失败'
    if (err.response?.status === 401) {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('username')
      router.push('/login')
    }
    message.error(msg)
    return Promise.reject(err)
  },
)

export default api

function ok(data) {
  return Promise.resolve({ data })
}

async function tryRequest(request, fallback) {
  try {
    return await request()
  } catch {
    return ok(typeof fallback === 'function' ? fallback() : fallback)
  }
}

export const register = (data) => tryRequest(() => api.post('/user/register', data), { success: true })
export const registerCode = (data) => tryRequest(() => api.post('/user/register/code', data), { success: true })
export const login = (data) =>
  tryRequest(() => api.post('/user/login', data), {
    access_token: 'mock-access-token',
    refresh_token: 'mock-refresh-token',
    user_id: 1,
    username: data.account || 'GoDan 用户',
    avatar: getUserProfile(1).avatar,
  })
export const loginCode = (data) =>
  tryRequest(() => api.post('/user/login/code', data), {
    access_token: 'mock-access-token',
    refresh_token: 'mock-refresh-token',
    user_id: 1,
    username: data.account || 'GoDan 用户',
    avatar: getUserProfile(1).avatar,
  })
export const sendCode = (target) => tryRequest(() => api.post('/user/code/send', { target }), { success: true })
export const refreshToken = (token) => tryRequest(() => api.post('/user/refresh', { refresh_token: token }), { access_token: token })
export const getProfile = (id) => tryRequest(() => api.get(id ? `/user/profile/${id}` : '/user/profile'), getUserProfile(id || 1))
export const updateProfile = (data) => tryRequest(() => api.put('/user/profile', data), { success: true })
export const getFollowers = (id, page) => tryRequest(() => api.get(`/user/${id}/followers`, { params: { page, page_size: 20 } }), makePage([getUserProfile(id)], page, 20))
export const getFollowees = (id, page) => tryRequest(() => api.get(`/user/${id}/followees`, { params: { page, page_size: 20 } }), makePage([getUserProfile(id)], page, 20))
export const follow = (userId) => tryRequest(() => api.post('/user/follow', { user_id: userId }), { success: true })
export const unfollow = (userId) => tryRequest(() => api.post('/user/unfollow', { user_id: userId }), { success: true })

export const initUpload = (data) =>
  tryRequest(() => api.post('/video/upload/init', data), {
    upload_id: `mock-${Date.now()}`,
    chunk_count: Math.max(1, Math.ceil((data.file_size || 1) / (5 << 20))),
  })
export const uploadChunk = (formData) => tryRequest(() => api.post('/video/upload/chunk', formData, { headers: { 'Content-Type': 'multipart/form-data' } }), { success: true })
export const completeUpload = (data) => tryRequest(() => api.post('/video/upload/complete', data), { video_id: 1001 })
export const uploadStatus = (uploadId) => tryRequest(() => api.get('/video/upload/status', { params: { upload_id: uploadId } }), { status: 'done' })
export const uploadAvatar = (formData) => tryRequest(() => api.post('/upload/avatar', formData, { headers: { 'Content-Type': 'multipart/form-data' } }), { url: getUserProfile(1).avatar })
export const uploadCover = (formData) => tryRequest(() => api.post('/video/cover/upload', formData, { headers: { 'Content-Type': 'multipart/form-data' } }), { url: mockVideos[0].cover_url })
export const getVideoDetail = (id) => tryRequest(() => api.get(`/video/${id}`), getVideoById(id))
export const getHotVideos = (page = 1) => tryRequest(() => api.get('/videos/hot', { params: { page, page_size: 12 } }), makePage(mockVideos, page, 12))
export const getCategoryVideos = (categoryId, sort, page) => tryRequest(() => api.get('/videos', { params: { category_id: categoryId, sort, page, page_size: 12 } }), makePage(mockVideos.filter((item) => item.category_id === categoryId), page, 12))
export const searchVideos = (q, page) => tryRequest(() => api.get('/videos/search', { params: { q, page, page_size: 12 } }), makePage(searchMockVideos(q), page, 12))
export const getRelatedVideos = (id, page) => tryRequest(() => api.get(`/video/${id}/related`, { params: { page, page_size: 8 } }), makePage(mockVideos.filter((item) => String(item.id) !== String(id)).slice(0, 8), page, 8))
export const getUserVideos = (id, page) => tryRequest(() => api.get(`/user/${id}/videos`, { params: { page, page_size: 12 } }), makePage(mockVideos.filter((item) => String(item.user_id) === String(id)), page, 12))
export const publishVideo = (id) => tryRequest(() => api.post(`/video/${id}/publish`), { success: true })
export const deleteVideo = (id) => tryRequest(() => api.delete(`/video/${id}`), { success: true })
export const updateVideoCover = (id, cover_url) => tryRequest(() => api.put(`/video/${id}/cover`, { cover_url }), { success: true })

export const likeVideo = (id) => tryRequest(() => api.post(`/video/${id}/like`), { success: true })
export const cancelLike = (id) => tryRequest(() => api.delete(`/video/${id}/like`), { success: true })
export const likeStatus = (id) => tryRequest(() => api.get(`/video/${id}/like/status`), { type: 0 })
export const giveCoin = (id, count) => tryRequest(() => api.post(`/video/${id}/coin`, { count }), { success: true })
export const shareVideo = (id) => tryRequest(() => api.post(`/video/${id}/share`), { link: `/video/${id}` })

export const getComments = (videoId, sort, page) => tryRequest(() => api.get('/comments', { params: { video_id: videoId, sort, page, page_size: 10 } }), makePage(mockComments, page, 10))
export const getReplies = (rootId, page) => tryRequest(() => api.get('/comments/replies', { params: { root_id: rootId, page, page_size: 10 } }), makePage(mockReplies[rootId] || [], page, 10))
export const postComment = (data) => tryRequest(() => api.post('/comment', data), { success: true })
export const deleteComment = (id) => tryRequest(() => api.delete(`/comment/${id}`), { success: true })
export const likeComment = (id) => tryRequest(() => api.post(`/comment/${id}/like`), { success: true })
export const unlikeComment = (id) => tryRequest(() => api.delete(`/comment/${id}/like`), { success: true })

export const getDanmakus = (videoId, start, end) =>
  tryRequest(() => api.get('/danmakus', { params: { video_id: videoId, start, end } }), [
    { text: '这版播放器真顺滑', time: 3, color: '#ffffff' },
    { text: 'phase 12 完成！', time: 8, color: '#00aeec' },
    { text: '右侧推荐区也很像 B 站', time: 15, color: '#ff85ad' },
  ])

export const getFolders = () => tryRequest(() => api.get('/user/folders'), { list: mockFolders })
export const createFolder = (data) => tryRequest(() => api.post('/favorite/folder', data), { success: true })
export const deleteFolder = (id) => tryRequest(() => api.delete(`/favorite/folder/${id}`), { success: true })
export const addToFolder = (data) => tryRequest(() => api.post('/favorite/add', data), { success: true })
export const removeFromFolder = (data) => tryRequest(() => api.post('/favorite/remove', data), { success: true })
export const getFolderItems = (id, page) => tryRequest(() => api.get(`/favorite/folder/${id}/items`, { params: { page, page_size: 12 } }), makePage(mockVideos, page, 12))

export const getTimeline = (page) => tryRequest(() => api.get('/timeline', { params: { page, page_size: 10 } }), makePage(mockTimeline, page, 10))
export const getNotifications = (page) => tryRequest(() => api.get('/notifications', { params: { page, page_size: 10 } }), makePage([], page, 10))
export const getUnreadCount = () => tryRequest(() => api.get('/notifications/unread'), { count: 3 })
export const markRead = (id) => tryRequest(() => api.post(`/notifications/${id}/read`), { success: true })
export const markAllRead = () => tryRequest(() => api.post('/notifications/read-all'), { success: true })

export const getLiveList = (page) => tryRequest(() => api.get('/live/list', { params: { page, page_size: 12 } }), makePage(mockLiveRooms, page, 12))
export const getRoomInfo = (id) => tryRequest(() => api.get(`/live/room/${id}`), mockLiveRooms.find((item) => String(item.id) === String(id)) || mockLiveRooms[0])
export const createRoom = (data) => tryRequest(() => api.post('/live/room', data), { id: 1 })
export const startLive = (id) => tryRequest(() => api.post(`/live/room/${id}/start`), { success: true })
export const stopLive = (id) => tryRequest(() => api.post(`/live/room/${id}/stop`), { success: true })
export const getGiftList = () => tryRequest(() => api.get('/live/gifts'), [{ id: 1, name: '小心心' }])
export const sendGift = (data) => tryRequest(() => api.post('/live/gift', data), { success: true })
export const getGiftRank = (id) => tryRequest(() => api.get(`/live/room/${id}/rank`), [])
