<template>
  <MainLayout>
    <div class="follow-container">
      <h1 class="page-title">我的关注</h1>
      
      <div class="follow-content">
        <el-tabs v-model="activeTab" class="follow-tabs">
          <el-tab-pane label="关注的文章" name="articles">
            <div v-if="followedArticles.length > 0" class="article-list">
              <div v-for="article in followedArticles" :key="article.id" class="article-card">
                <div class="article-cover" v-if="article.coverImage" @click="viewArticle(article.id)">
                  <img :src="article.coverImage" :alt="article.title" />
                </div>
                <div class="article-info">
                  <h3 class="article-title" @click="viewArticle(article.id)">{{ article.title }}</h3>
                  <p class="article-abstract">{{ article.abstract }}</p>
                  <div class="article-meta">
                    <div class="author-info">
                      <el-avatar :size="24" :src="article.author.avatar">
                        {{ article.author.name.substring(0, 1) }}
                      </el-avatar>
                      <span class="author-name">{{ article.author.name }}</span>
                    </div>
                    <div class="interaction-info">
                      <span class="interaction-item">
                        <el-icon><View /></el-icon>
                        {{ formatNumber(article.readCount) }}
                      </span>
                      <span class="interaction-item">
                        <el-icon><ThumbsUp /></el-icon>
                        {{ formatNumber(article.likeCount) }}
                      </span>
                      <span class="interaction-item">
                        <el-icon><ChatDotRound /></el-icon>
                        {{ formatNumber(article.commentCount) }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
              
              <div class="load-more" v-if="hasMoreArticles">
                <el-button @click="loadMoreArticles">加载更多</el-button>
              </div>
            </div>
            <div v-else class="empty-state">
              <el-empty description="暂无关注的文章" />
              <el-button type="primary" @click="goToHome">去首页发现更多内容</el-button>
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="关注的用户" name="users">
            <div v-if="followedUsers.length > 0" class="user-list">
              <div v-for="user in followedUsers" :key="user.id" class="user-card">
                <div class="user-avatar">
                  <el-avatar :size="50" :src="user.avatar">
                    {{ user.name.substring(0, 1) }}
                  </el-avatar>
                </div>
                <div class="user-info">
                  <div class="user-name" @click="viewUser(user.id)">{{ user.name }}</div>
                  <div class="user-desc">{{ user.description }}</div>
                  <div class="user-stats">
                    <span>{{ user.articleCount }} 篇文章</span>
                    <span>{{ formatNumber(user.followerCount) }} 粉丝</span>
                  </div>
                </div>
                <div class="user-action">
                  <el-button type="primary" @click="unfollowUser(user.id)">已关注</el-button>
                </div>
              </div>
              
              <div class="load-more" v-if="hasMoreUsers">
                <el-button @click="loadMoreUsers">加载更多</el-button>
              </div>
            </div>
            <div v-else class="empty-state">
              <el-empty description="暂无关注的用户" />
              <el-button type="primary" @click="goToHome">去首页发现更多用户</el-button>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import MainLayout from '@/components/layout/MainLayout.vue'
import { View, ChatDotRound } from '@element-plus/icons-vue'
import { Star as ThumbsUp } from '@element-plus/icons-vue'

const router = useRouter()
const activeTab = ref('articles')

// 模拟数据 - 关注的文章
const followedArticles = ref([
  {
    id: 1,
    title: '如何在家制作完美的提拉米苏',
    abstract: '提拉米苏是一道经典的意大利甜点，本文将分享专业大厨的独家秘方...',
    coverImage: 'https://picsum.photos/id/431/400/300',
    author: {
      id: 101,
      name: '美食达人',
      avatar: 'https://picsum.photos/id/1027/100/100'
    },
    readCount: 12500,
    likeCount: 3200,
    commentCount: 128
  },
  {
    id: 2,
    title: '2025年最值得去的10个小众旅行地',
    abstract: '厌倦了人山人海的热门景点？这些鲜为人知的目的地将带给你全新的旅行体验...',
    coverImage: 'https://picsum.photos/id/1036/400/300',
    author: {
      id: 102,
      name: '旅行笔记',
      avatar: 'https://picsum.photos/id/1012/100/100'
    },
    readCount: 18700,
    likeCount: 5400,
    commentCount: 342
  }
])

// 模拟数据 - 关注的用户
const followedUsers = ref([
  {
    id: 101,
    name: '美食达人',
    avatar: 'https://picsum.photos/id/1027/100/100',
    description: '探索美食世界的专业吃货',
    articleCount: 45,
    followerCount: 12800
  },
  {
    id: 102,
    name: '旅行笔记',
    avatar: 'https://picsum.photos/id/1012/100/100',
    description: '记录世界各地的美景与文化',
    articleCount: 78,
    followerCount: 25600
  },
  {
    id: 103,
    name: '生活家',
    avatar: 'https://picsum.photos/id/1005/100/100',
    description: '分享高品质生活方式',
    articleCount: 32,
    followerCount: 9400
  }
])

const hasMoreArticles = ref(false)
const hasMoreUsers = ref(true)

// 格式化数字，例如1200显示为1.2k
const formatNumber = (num: number): string => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

// 查看文章详情
const viewArticle = (id: number) => {
  router.push(`/article/${id}`)
}

// 查看用户主页
const viewUser = (id: number) => {
  router.push(`/user/${id}`)
}

// 取消关注用户
const unfollowUser = (id: number) => {
  // 这里应该调用取消关注API
  // 模拟取消关注
  followedUsers.value = followedUsers.value.filter(user => user.id !== id)
}

// 加载更多文章
const loadMoreArticles = () => {
  // 这里应该调用API加载更多关注的文章
  // 模拟加载更多
  const moreArticles = [
    {
      id: 3,
      title: '极简主义：如何通过断舍离改变你的生活',
      abstract: '极简主义不仅是一种生活方式，更是一种思维模式。本文将分享如何开始你的极简之旅...',
      coverImage: 'https://picsum.photos/id/106/400/300',
      author: {
        id: 103,
        name: '生活家',
        avatar: 'https://picsum.photos/id/1005/100/100'
      },
      readCount: 9800,
      likeCount: 2100,
      commentCount: 98
    }
  ]
  
  followedArticles.value = [...followedArticles.value, ...moreArticles]
  
  // 假设没有更多文章了
  hasMoreArticles.value = false
}

// 加载更多用户
const loadMoreUsers = () => {
  // 这里应该调用API加载更多关注的用户
  // 模拟加载更多
  const moreUsers = [
    {
      id: 104,
      name: '摄影师小王',
      avatar: 'https://picsum.photos/id/1062/100/100',
      description: '用镜头记录生活的美好瞬间',
      articleCount: 56,
      followerCount: 15300
    }
  ]
  
  followedUsers.value = [...followedUsers.value, ...moreUsers]
  
  // 假设没有更多用户了
  hasMoreUsers.value = false
}

// 前往首页
const goToHome = () => {
  router.push('/')
}

onMounted(() => {
  // 这里应该调用API获取关注的文章和用户
  // 目前使用模拟数据
})
</script>

<style scoped>
.follow-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.page-title {
  font-size: 24px;
  margin-bottom: 20px;
  color: var(--text-primary);
}

.follow-content {
  background-color: var(--bg-primary);
  border-radius: var(--border-radius-md);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}

.follow-tabs {
  margin-bottom: 20px;
}

.article-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.article-card {
  background-color: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s;
}

.article-card:hover {
  transform: translateY(-5px);
}

.article-cover {
  width: 100%;
  height: 180px;
  overflow: hidden;
  cursor: pointer;
}

.article-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.article-cover:hover img {
  transform: scale(1.05);
}

.article-info {
  padding: 15px;
}

.article-title {
  margin: 0 0 10px;
  font-size: 16px;
  font-weight: 600;
  line-height: 1.4;
  cursor: pointer;
}

.article-title:hover {
  color: var(--primary-color);
}

.article-abstract {
  margin: 0 0 15px;
  font-size: 14px;
  color: #666;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.author-name {
  font-size: 14px;
  color: #333;
}

.interaction-info {
  display: flex;
  gap: 12px;
}

.interaction-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #999;
}

.user-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.user-card {
  display: flex;
  align-items: center;
  padding: 15px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.user-avatar {
  margin-right: 15px;
}

.user-info {
  flex: 1;
}

.user-name {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 5px;
  cursor: pointer;
}

.user-name:hover {
  color: var(--primary-color);
}

.user-desc {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.user-stats {
  font-size: 12px;
  color: #999;
}

.user-stats span {
  margin-right: 15px;
}

.user-action {
  margin-left: 15px;
}

.load-more {
  text-align: center;
  margin-top: 30px;
}

.empty-state {
  text-align: center;
  padding: 40px 0;
}

.empty-state .el-button {
  margin-top: 20px;
}
</style>
