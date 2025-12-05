<template>
  <MainLayout>
    <div class="search-container">
      <div class="search-header">
        <div class="search-input-container">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索文章、用户"
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
            <template #append>
              <el-button @click="handleSearch">搜索</el-button>
            </template>
          </el-input>
        </div>
      </div>
      
      <div class="search-content">
        <template v-if="!searchKeyword">
          <!-- 未搜索时显示历史记录和热门搜索 -->
          <div class="search-section">
            <div class="section-header">
              <h3>搜索历史</h3>
              <el-button type="text" @click="clearHistory">清除</el-button>
            </div>
            <div class="history-tags" v-if="searchHistory.length > 0">
              <el-tag
                v-for="item in searchHistory"
                :key="item"
                @click="searchKeyword = item; handleSearch()"
                class="history-tag"
              >
                {{ item }}
              </el-tag>
            </div>
            <div v-else class="empty-history">
              <el-empty description="暂无搜索历史" :image-size="100" />
            </div>
          </div>
          
          <div class="search-section">
            <div class="section-header">
              <h3>热门搜索</h3>
            </div>
            <div class="hot-search-list">
              <div
                v-for="(item, index) in hotSearches"
                :key="index"
                class="hot-search-item"
                @click="searchKeyword = item; handleSearch()"
              >
                <span class="hot-search-rank" :class="{ 'top-three': index < 3 }">{{ index + 1 }}</span>
                <span class="hot-search-keyword">{{ item }}</span>
              </div>
            </div>
          </div>
        </template>
        
        <template v-else-if="isSearching">
          <!-- 搜索中显示加载状态 -->
          <div class="search-loading">
            <el-skeleton :rows="10" animated />
          </div>
        </template>
        
        <template v-else>
          <!-- 搜索结果 -->
          <div class="search-tabs">
            <el-tabs v-model="activeTab">
              <el-tab-pane label="全部" name="all">
                <!-- 用户结果 -->
                <div class="search-section" v-if="searchResults.users.length > 0">
                  <div class="section-header">
                    <h3>用户</h3>
                    <el-button type="text" @click="activeTab = 'users'">查看全部</el-button>
                  </div>
                  <div class="user-results">
                    <div v-for="user in searchResults.users.slice(0, 3)" :key="user.id" class="user-result-item">
                      <el-avatar :size="50" :src="user.avatar">{{ user.name.substring(0, 1) }}</el-avatar>
                      <div class="user-info">
                        <div class="user-name">{{ user.name }}</div>
                        <div class="user-desc">{{ user.description }}</div>
                      </div>
                      <el-button
                        size="small"
                        :type="user.isFollowing ? 'info' : 'primary'"
                        @click.stop="toggleFollow(user)"
                      >
                        {{ user.isFollowing ? '已关注' : '关注' }}
                      </el-button>
                    </div>
                  </div>
                </div>
                
                <!-- 文章结果 -->
                <div class="search-section">
                  <div class="section-header">
                    <h3>文章</h3>
                    <el-button type="text" @click="activeTab = 'articles'">查看全部</el-button>
                  </div>
                  <div class="article-results">
                    <div v-for="article in searchResults.articles.slice(0, 5)" :key="article.id" class="article-result-item" @click="viewArticle(article.id)">
                      <div class="article-info">
                        <h3 class="article-title">{{ article.title }}</h3>
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
                          </div>
                        </div>
                      </div>
                      <div class="article-cover" v-if="article.coverImage">
                        <img :src="article.coverImage" :alt="article.title" />
                      </div>
                    </div>
                  </div>
                </div>
              </el-tab-pane>
              
              <el-tab-pane label="用户" name="users">
                <div class="user-results">
                  <div v-for="user in searchResults.users" :key="user.id" class="user-result-item">
                    <el-avatar :size="50" :src="user.avatar">{{ user.name.substring(0, 1) }}</el-avatar>
                    <div class="user-info">
                      <div class="user-name">{{ user.name }}</div>
                      <div class="user-desc">{{ user.description }}</div>
                    </div>
                    <el-button
                      size="small"
                      :type="user.isFollowing ? 'info' : 'primary'"
                      @click.stop="toggleFollow(user)"
                    >
                      {{ user.isFollowing ? '已关注' : '关注' }}
                    </el-button>
                  </div>
                </div>
                
                <div class="load-more" v-if="hasMoreUsers">
                  <el-button @click="loadMoreUsers">加载更多</el-button>
                </div>
                
                <div v-if="searchResults.users.length === 0" class="empty-results">
                  <el-empty description="未找到相关用户" />
                </div>
              </el-tab-pane>
              
              <el-tab-pane label="文章" name="articles">
                <div class="article-results">
                  <div v-for="article in searchResults.articles" :key="article.id" class="article-result-item" @click="viewArticle(article.id)">
                    <div class="article-info">
                      <h3 class="article-title">{{ article.title }}</h3>
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
                        </div>
                      </div>
                    </div>
                    <div class="article-cover" v-if="article.coverImage">
                      <img :src="article.coverImage" :alt="article.title" />
                    </div>
                  </div>
                </div>
                
                <div class="load-more" v-if="hasMoreArticles">
                  <el-button @click="loadMoreArticles">加载更多</el-button>
                </div>
                
                <div v-if="searchResults.articles.length === 0" class="empty-results">
                  <el-empty description="未找到相关文章" />
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
        </template>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import MainLayout from '@/components/layout/MainLayout.vue'
import { Search, View } from '@element-plus/icons-vue'
import { Star as ThumbsUp } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

const searchKeyword = ref('')
const isSearching = ref(false)
const activeTab = ref('all')

// 模拟搜索历史
const searchHistory = ref(['极简主义', '旅行攻略', '健康饮食'])

// 模拟热门搜索
const hotSearches = ref([
  '年轻人副业',
  '小众旅行地',
  '养成好习惯',
  '2025年技能',
  '朋友圈创意',
  '低成本装修',
  '科学健身',
  '财报分析',
  '摄影技巧',
  '手机摄影'
])

// 模拟搜索结果
const searchResults = ref({
  users: [
    {
      id: 101,
      name: '美食达人',
      avatar: 'https://picsum.photos/id/1027/100/100',
      description: '探索美食世界的专业吃货',
      isFollowing: false
    },
    {
      id: 102,
      name: '旅行笔记',
      avatar: 'https://picsum.photos/id/1012/100/100',
      description: '记录世界各地的美景与文化',
      isFollowing: true
    },
    {
      id: 103,
      name: '生活家',
      avatar: 'https://picsum.photos/id/1005/100/100',
      description: '分享高品质生活方式',
      isFollowing: false
    },
    {
      id: 104,
      name: '摄影师小王',
      avatar: 'https://picsum.photos/id/1062/100/100',
      description: '用镜头记录生活的美好瞬间',
      isFollowing: false
    }
  ],
  articles: [
    {
      id: 1,
      title: '如何在家制作完美的提拉米苏',
      abstract: '提拉米苏是一道经典的意大利甜点，本文将分享专业大厨的独家秘方...',
      coverImage: 'https://picsum.photos/id/431/200/200',
      author: {
        id: 101,
        name: '美食达人',
        avatar: 'https://picsum.photos/id/1027/100/100'
      },
      readCount: 12500,
      likeCount: 3200
    },
    {
      id: 2,
      title: '2025年最值得去的10个小众旅行地',
      abstract: '厌倦了人山人海的热门景点？这些鲜为人知的目的地将带给你全新的旅行体验...',
      coverImage: 'https://picsum.photos/id/1036/200/200',
      author: {
        id: 102,
        name: '旅行笔记',
        avatar: 'https://picsum.photos/id/1012/100/100'
      },
      readCount: 18700,
      likeCount: 5400
    },
    {
      id: 3,
      title: '极简主义：如何通过断舍离改变你的生活',
      abstract: '极简主义不仅是一种生活方式，更是一种思维模式。本文将分享如何开始你的极简之旅...',
      coverImage: 'https://picsum.photos/id/106/200/200',
      author: {
        id: 103,
        name: '生活家',
        avatar: 'https://picsum.photos/id/1005/100/100'
      },
      readCount: 9800,
      likeCount: 2100
    },
    {
      id: 4,
      title: '数字游民：如何边旅行边工作',
      abstract: '远程工作正在改变我们的生活和工作方式，本文分享如何成为一名成功的数字游民...',
      coverImage: 'https://picsum.photos/id/1081/200/200',
      author: {
        id: 104,
        name: '自由职业者',
        avatar: 'https://picsum.photos/id/1025/100/100'
      },
      readCount: 15300,
      likeCount: 4200
    },
    {
      id: 5,
      title: '家庭花园：从零开始的种植指南',
      abstract: '无论你是有着一片后院，还是只有一个小阳台，都可以打造自己的绿色天地...',
      coverImage: 'https://picsum.photos/id/145/200/200',
      author: {
        id: 105,
        name: '园艺爱好者',
        avatar: 'https://picsum.photos/id/1074/100/100'
      },
      readCount: 7600,
      likeCount: 1800
    },
    {
      id: 6,
      title: '如何提高你的摄影技巧：从入门到精通',
      abstract: '无需昂贵的设备，掌握这些基本技巧，你也能拍出令人惊艳的照片...',
      coverImage: 'https://picsum.photos/id/250/200/200',
      author: {
        id: 106,
        name: '摄影师小王',
        avatar: 'https://picsum.photos/id/1062/100/100'
      },
      readCount: 21000,
      likeCount: 6300
    }
  ]
})

const hasMoreUsers = ref(false)
const hasMoreArticles = ref(true)

// 格式化数字，例如1200显示为1.2k
const formatNumber = (num: number): string => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

// 处理搜索
const handleSearch = () => {
  if (!searchKeyword.value.trim()) return
  
  isSearching.value = true
  
  // 添加到搜索历史
  if (!searchHistory.value.includes(searchKeyword.value)) {
    searchHistory.value.unshift(searchKeyword.value)
    if (searchHistory.value.length > 10) {
      searchHistory.value.pop()
    }
    // 保存到本地存储
    localStorage.setItem('searchHistory', JSON.stringify(searchHistory.value))
  }
  
  // 更新URL
  router.push({
    path: '/search',
    query: { q: searchKeyword.value }
  })
  
  // 这里应该调用搜索API
  // 模拟搜索延迟
  setTimeout(() => {
    // 实际项目中这里应该根据搜索关键词过滤结果
    isSearching.value = false
  }, 500)
}

// 清除搜索历史
const clearHistory = () => {
  searchHistory.value = []
  localStorage.removeItem('searchHistory')
}

// 查看文章详情
const viewArticle = (id: number) => {
  router.push(`/article/${id}`)
}

// 关注/取消关注用户
const toggleFollow = (user: any) => {
  user.isFollowing = !user.isFollowing
  // 这里应该调用关注/取消关注API
}

// 加载更多用户
const loadMoreUsers = () => {
  // 这里应该调用API加载更多用户
  // 模拟加载更多
  hasMoreUsers.value = false
}

// 加载更多文章
const loadMoreArticles = () => {
  // 这里应该调用API加载更多文章
  // 模拟加载更多
  const moreArticles = [
    {
      id: 7,
      title: '零基础学习编程：从何开始？',
      abstract: '想学编程但不知道从何入手？本文为你提供清晰的学习路径...',
      coverImage: 'https://picsum.photos/id/0/200/200',
      author: {
        id: 107,
        name: '编程教练',
        avatar: 'https://picsum.photos/id/1/100/100'
      },
      readCount: 18500,
      likeCount: 5100
    },
    {
      id: 8,
      title: '每天10分钟，21天养成冥想习惯',
      abstract: '冥想不仅能减轻压力，还能提高注意力和创造力...',
      coverImage: 'https://picsum.photos/id/1029/200/200',
      author: {
        id: 108,
        name: '心灵导师',
        avatar: 'https://picsum.photos/id/1002/100/100'
      },
      readCount: 12300,
      likeCount: 3600
    }
  ]
  
  searchResults.value.articles = [...searchResults.value.articles, ...moreArticles]
  
  // 假设没有更多文章了
  hasMoreArticles.value = false
}

// 监听路由变化
watch(() => route.query.q, (newQuery) => {
  if (newQuery) {
    searchKeyword.value = newQuery as string
    handleSearch()
  }
}, { immediate: true })

onMounted(() => {
  // 从本地存储加载搜索历史
  const storedHistory = localStorage.getItem('searchHistory')
  if (storedHistory) {
    searchHistory.value = JSON.parse(storedHistory)
  }
  
  // 如果URL中有查询参数，执行搜索
  if (route.query.q) {
    searchKeyword.value = route.query.q as string
    handleSearch()
  }
})
</script>

<style scoped>
.search-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.search-header {
  margin-bottom: 20px;
}

.search-input-container {
  max-width: 600px;
  margin: 0 auto;
}

.search-input {
  width: 100%;
}

.search-content {
  background-color: var(--bg-primary);
  border-radius: var(--border-radius-md);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}

.search-section {
  margin-bottom: 30px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.section-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--text-primary);
}

.history-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.history-tag {
  cursor: pointer;
}

.hot-search-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
}

.hot-search-item {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #f5f5f5;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.hot-search-item:hover {
  background-color: #f0f0f0;
}

.hot-search-rank {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #e0e0e0;
  color: #666;
  border-radius: 50%;
  font-weight: bold;
  font-size: 14px;
  margin-right: 10px;
}

.hot-search-rank.top-three {
  background-color: var(--primary-color);
  color: white;
}

.hot-search-keyword {
  font-size: 14px;
  color: #333;
}

.search-tabs {
  margin-top: 20px;
}

.user-results {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.user-result-item {
  display: flex;
  align-items: center;
  padding: 15px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.user-info {
  flex: 1;
  margin-left: 15px;
}

.user-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 5px;
}

.user-desc {
  font-size: 14px;
  color: #666;
}

.article-results {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.article-result-item {
  display: flex;
  padding: 15px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  cursor: pointer;
  transition: transform 0.3s;
}

.article-result-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.article-info {
  flex: 1;
}

.article-title {
  margin: 0 0 10px;
  font-size: 18px;
  font-weight: 600;
  line-height: 1.4;
  color: #333;
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

.article-cover {
  width: 120px;
  height: 120px;
  overflow: hidden;
  border-radius: 8px;
  margin-left: 20px;
  flex-shrink: 0;
}

.article-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.load-more {
  text-align: center;
  margin-top: 20px;
}

.empty-history,
.empty-results {
  text-align: center;
  padding: 20px 0;
}

.search-loading {
  padding: 20px 0;
}
</style>
