import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

const api = axios.create({ baseURL: '/api/v1', timeout: 15000 })

api.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  res => res.data,
  err => {
    const msg = err.response?.data?.message || '请求失败'
    if (err.response?.status === 401) {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('username')
      router.push('/login')
    }
    ElMessage.error(msg)
    return Promise.reject(err)
  }
)

export default api

// User
export const register = data => api.post('/user/register', data)
export const registerCode = data => api.post('/user/register/code', data)
export const login = data => api.post('/user/login', data)
export const loginCode = data => api.post('/user/login/code', data)
export const sendCode = target => api.post('/user/code/send', { target })
export const refreshToken = token => api.post('/user/refresh', { refresh_token: token })
export const getProfile = id => api.get(id ? `/user/profile/${id}` : '/user/profile')
export const updateProfile = data => api.put('/user/profile', data)
export const getFollowers = (id, page) => api.get(`/user/${id}/followers`, { params: { page, page_size: 20 } })
export const getFollowees = (id, page) => api.get(`/user/${id}/followees`, { params: { page, page_size: 20 } })
export const follow = userId => api.post('/user/follow', { user_id: userId })
export const unfollow = userId => api.post('/user/unfollow', { user_id: userId })

// Video
export const initUpload = data => api.post('/video/upload/init', data)
export const uploadChunk = formData => api.post('/video/upload/chunk', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
export const completeUpload = data => api.post('/video/upload/complete', data)
export const uploadStatus = uploadId => api.get('/video/upload/status', { params: { upload_id: uploadId } })
export const uploadAvatar = formData => api.post('/upload/avatar', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
export const uploadCover = formData => api.post('/video/cover/upload', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
export const getVideoDetail = id => api.get(`/video/${id}`)
export const getHotVideos = (page = 1) => api.get('/videos/hot', { params: { page, page_size: 12 } })
export const getCategoryVideos = (categoryId, sort, page) => api.get('/videos', { params: { category_id: categoryId, sort, page, page_size: 12 } })
export const searchVideos = (q, page) => api.get('/videos/search', { params: { q, page, page_size: 12 } })
export const getRelatedVideos = (id, page) => api.get(`/video/${id}/related`, { params: { page, page_size: 8 } })
export const getUserVideos = (id, page) => api.get(`/user/${id}/videos`, { params: { page, page_size: 12 } })
export const publishVideo = id => api.post(`/video/${id}/publish`)
export const deleteVideo = id => api.delete(`/video/${id}`)
export const updateVideoCover = (id, cover_url) => api.put(`/video/${id}/cover`, { cover_url })

// Interaction
export const likeVideo = id => api.post(`/video/${id}/like`)
export const cancelLike = id => api.delete(`/video/${id}/like`)
export const likeStatus = id => api.get(`/video/${id}/like/status`)
export const giveCoin = (id, count) => api.post(`/video/${id}/coin`, { count })
export const shareVideo = id => api.post(`/video/${id}/share`)

// Comments
export const getComments = (videoId, sort, page) => api.get('/comments', { params: { video_id: videoId, sort, page, page_size: 10 } })
export const getReplies = (rootId, page) => api.get('/comments/replies', { params: { root_id: rootId, page, page_size: 10 } })
export const postComment = data => api.post('/comment', data)
export const deleteComment = id => api.delete(`/comment/${id}`)
export const likeComment = id => api.post(`/comment/${id}/like`)
export const unlikeComment = id => api.delete(`/comment/${id}/like`)

// Danmaku
export const getDanmakus = (videoId, start, end) => api.get('/danmakus', { params: { video_id: videoId, start, end } })

// Favorite
export const getFolders = () => api.get('/user/folders')
export const createFolder = data => api.post('/favorite/folder', data)
export const deleteFolder = id => api.delete(`/favorite/folder/${id}`)
export const addToFolder = data => api.post('/favorite/add', data)
export const removeFromFolder = data => api.post('/favorite/remove', data)
export const getFolderItems = (id, page) => api.get(`/favorite/folder/${id}/items`, { params: { page, page_size: 12 } })

// Timeline & Notifications
export const getTimeline = page => api.get('/timeline', { params: { page, page_size: 10 } })
export const getNotifications = page => api.get('/notifications', { params: { page, page_size: 10 } })
export const getUnreadCount = () => api.get('/notifications/unread')
export const markRead = id => api.post(`/notifications/${id}/read`)
export const markAllRead = () => api.post('/notifications/read-all')

// Live
export const getLiveList = page => api.get('/live/list', { params: { page, page_size: 12 } })
export const getRoomInfo = id => api.get(`/live/room/${id}`)
export const createRoom = data => api.post('/live/room', data)
export const startLive = id => api.post(`/live/room/${id}/start`)
export const stopLive = id => api.post(`/live/room/${id}/stop`)
export const getGiftList = () => api.get('/live/gifts')
export const sendGift = data => api.post('/live/gift', data)
export const getGiftRank = id => api.get(`/live/room/${id}/rank`)
