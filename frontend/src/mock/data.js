const videoSource = 'https://interactive-examples.mdn.mozilla.net/media/cc0-videos/flower.mp4'

function svgDataUrl(title, bg = '#89d5ff', fg = '#ffffff') {
  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" width="640" height="360" viewBox="0 0 640 360">
      <defs>
        <linearGradient id="g" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stop-color="${bg}" />
          <stop offset="100%" stop-color="#ffffff" stop-opacity="0.35" />
        </linearGradient>
      </defs>
      <rect width="640" height="360" rx="28" fill="url(#g)" />
      <circle cx="544" cy="96" r="54" fill="rgba(255,255,255,0.28)" />
      <circle cx="108" cy="286" r="76" fill="rgba(255,255,255,0.18)" />
      <text x="48" y="286" font-size="34" font-family="Segoe UI, Arial" fill="${fg}" font-weight="700">${title}</text>
    </svg>
  `
  return `data:image/svg+xml;charset=UTF-8,${encodeURIComponent(svg)}`
}

const authors = [
  { id: 1, username: '阿澜同学', avatar: svgDataUrl('A', '#ffb4d1') },
  { id: 2, username: '木星实验室', avatar: svgDataUrl('M', '#9fd8ff') },
  { id: 3, username: '像素研究所', avatar: svgDataUrl('P', '#b9f0c5') },
  { id: 4, username: '风铃动画社', avatar: svgDataUrl('F', '#ffd79f') },
]

const categories = ['动画', '游戏', '知识', '音乐', '科技', '生活']

function createVideo(id, title, category, author, accent, extra = {}) {
  return {
    id,
    user_id: author.id,
    title,
    description: `${title} 的完整演示页，包含简介、互动区、相关推荐与播放器能力，方便直接对接后端接口。`,
    cover_url: svgDataUrl(title.slice(0, 10), accent),
    video_url: videoSource,
    duration: 60 + id * 17,
    play_count: 16000 + id * 930,
    like_count: 2200 + id * 108,
    coin_count: 360 + id * 13,
    fav_count: 1200 + id * 42,
    share_count: 300 + id * 9,
    danmaku_count: 900 + id * 31,
    category_id: categories.indexOf(category) + 1,
    category,
    tags: [category, 'Vue3', 'GoDan'],
    created_at: `2026-06-${String((id % 20) + 1).padStart(2, '0')} 19:20:00`,
    author,
    ...extra,
  }
}

export const mockVideos = [
  createVideo(1, 'GoDan 首页视觉升级实录', '科技', authors[1], '#74d2ff', { description: '从全局布局到卡片动效，完整搭建 B 站风格首页。' }),
  createVideo(2, '这套播放器交互也太像 B 站了', '动画', authors[0], '#ff9cc3'),
  createVideo(3, 'ArtPlayer + Vue3 实战整合', '知识', authors[2], '#7fdac0'),
  createVideo(4, '上传页怎么做得更有氛围感', '生活', authors[3], '#ffbc69'),
  createVideo(5, '搜索结果页重构：从可用到好看', '科技', authors[1], '#91c9ff'),
  createVideo(6, '登录页也能有作品感', '音乐', authors[0], '#c8a4ff'),
  createVideo(7, '动态流布局改造笔记', '生活', authors[2], '#96ebb4'),
  createVideo(8, '直播页原型搭建过程分享', '游戏', authors[3], '#ffb485'),
  createVideo(9, 'UP 主主页的内容组织方式', '知识', authors[1], '#8ed7ff'),
  createVideo(10, '评论区交互体验优化', '动画', authors[0], '#ffafbf'),
  createVideo(11, '推荐区信息密度怎么拿捏', '科技', authors[2], '#82d8c8'),
  createVideo(12, 'GoDan phase 12 完成演示', '科技', authors[3], '#ffd1a1'),
]

export const mockFolders = [
  { id: 1, name: '稍后再看' },
  { id: 2, name: '前端灵感' },
  { id: 3, name: '播放器拆解' },
]

export const mockComments = [
  {
    id: 1,
    username: '薄荷汽水',
    avatar: svgDataUrl('薄', '#9cd7ff'),
    content: '这版首页的信息密度很舒服，像 B 站但没有生硬照搬。',
    created_at: '2026-06-24 13:00:00',
    like_count: 54,
    reply_count: 2,
  },
  {
    id: 2,
    username: '夜航星',
    avatar: svgDataUrl('夜', '#c5f1a8'),
    content: '播放器和右侧推荐并排的比例拿捏得不错，移动端也能收得住。',
    created_at: '2026-06-25 20:10:00',
    like_count: 29,
    reply_count: 1,
  },
]

export const mockReplies = {
  1: [
    { id: 11, username: 'GoDan_dev', content: '谢谢，后面会继续把分区页和评论楼中楼补完整。' },
    { id: 12, username: '青柠', content: '期待接入真实弹幕数据。' },
  ],
  2: [
    { id: 21, username: '阿澜同学', content: '详情页的栅格我也调了几轮，终于顺眼了。' },
  ],
}

export const mockLiveRooms = [
  { id: 1, title: '深夜前端改版直播', username: '木星实验室', cover_url: svgDataUrl('LIVE', '#89d5ff'), viewer_count: 4210 },
  { id: 2, title: '动画分镜拆解会', username: '风铃动画社', cover_url: svgDataUrl('ANIME', '#ffb4d1'), viewer_count: 2866 },
  { id: 3, title: 'Go 项目架构答疑', username: '像素研究所', cover_url: svgDataUrl('GO', '#a4ebbe'), viewer_count: 1988 },
]

export const mockTimeline = [
  { id: 1, username: '阿澜同学', avatar: svgDataUrl('A', '#ffb4d1'), type: 1, created_at: '2026-06-26 18:30:00', content: '发布了新视频《登录页也能有作品感》' },
  { id: 2, username: '木星实验室', avatar: svgDataUrl('M', '#9fd8ff'), type: 2, created_at: '2026-06-26 16:10:00', content: '点赞了《ArtPlayer + Vue3 实战整合》' },
  { id: 3, username: '像素研究所', avatar: svgDataUrl('P', '#b9f0c5'), type: 4, created_at: '2026-06-25 22:18:00', content: '收藏了一个播放器灵感合集' },
]

export function makePage(list, page = 1, pageSize = 12) {
  const start = (page - 1) * pageSize
  return {
    list: list.slice(start, start + pageSize),
    total: list.length,
  }
}

export function getVideoById(id) {
  const current = mockVideos.find((item) => String(item.id) === String(id)) || mockVideos[0]
  return {
    ...current,
    author: current.author,
  }
}

export function searchMockVideos(keyword) {
  if (!keyword) return mockVideos
  const key = keyword.toLowerCase()
  return mockVideos.filter((item) =>
    [item.title, item.description, item.category, item.author.username].join(' ').toLowerCase().includes(key),
  )
}

export function getUserProfile(id) {
  const author = authors.find((item) => String(item.id) === String(id)) || authors[0]
  const videos = mockVideos.filter((item) => item.user_id === author.id)
  return {
    id: author.id,
    username: author.username,
    avatar: author.avatar,
    bio: '分享播放器设计、Vue3 页面实现与视频社区产品灵感。',
    follower_count: 12000 + author.id * 312,
    followee_count: 120 + author.id * 17,
    video_count: videos.length,
  }
}
