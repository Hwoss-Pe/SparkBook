<template>
  <MainLayout>
    <div class="message-container">
      <h1 class="page-title">消息中心</h1>
      
      <div class="message-content">
        <el-tabs v-model="activeTab" class="message-tabs">
          <el-tab-pane label="互动消息" name="interaction">
            <div v-if="interactionMessages.length > 0" class="message-list">
              <div v-for="message in interactionMessages" :key="message.id" class="message-item">
                <div class="message-avatar">
                  <el-avatar :size="40" :src="message.sender.avatar">
                    {{ message.sender.name.substring(0, 1) }}
                  </el-avatar>
                </div>
                <div class="message-body">
                  <div class="message-header">
                    <span class="sender-name">{{ message.sender.name }}</span>
                    <span class="message-time">{{ formatTime(message.time) }}</span>
                  </div>
                  <div class="message-content-text" v-html="message.content"></div>
                  <div class="message-target" v-if="message.target" @click="navigateToTarget(message)">
                    <div class="target-content">
                      <div class="target-title">{{ message.target.title }}</div>
                      <div class="target-preview">{{ message.target.preview }}</div>
                    </div>
                  </div>
                </div>
              </div>
              
              <div class="load-more" v-if="hasMoreInteraction">
                <el-button @click="loadMoreInteraction">加载更多</el-button>
              </div>
            </div>
            <div v-else class="empty-state">
              <el-empty description="暂无互动消息" />
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="关注消息" name="follow">
            <div v-if="followMessages.length > 0" class="message-list">
              <div v-for="message in followMessages" :key="message.id" class="message-item">
                <div class="message-avatar">
                  <el-avatar :size="40" :src="message.sender.avatar">
                    {{ message.sender.name.substring(0, 1) }}
                  </el-avatar>
                </div>
                <div class="message-body">
                  <div class="message-header">
                    <span class="sender-name">{{ message.sender.name }}</span>
                    <span class="message-time">{{ formatTime(message.time) }}</span>
                  </div>
                  <div class="message-content-text">{{ message.content }}</div>
                  <div class="message-actions" v-if="!message.isFollowing">
                    <el-button size="small" type="primary" @click="followUser(message.sender.id)">关注</el-button>
                  </div>
                </div>
              </div>
              
              <div class="load-more" v-if="hasMoreFollow">
                <el-button @click="loadMoreFollow">加载更多</el-button>
              </div>
            </div>
            <div v-else class="empty-state">
              <el-empty description="暂无关注消息" />
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="系统消息" name="system">
            <div v-if="systemMessages.length > 0" class="message-list">
              <div v-for="message in systemMessages" :key="message.id" class="message-item system-message">
                <div class="message-avatar">
                  <el-avatar :size="40" icon="el-icon-bell">
                    系统
                  </el-avatar>
                </div>
                <div class="message-body">
                  <div class="message-header">
                    <span class="sender-name">系统通知</span>
                    <span class="message-time">{{ formatTime(message.time) }}</span>
                  </div>
                  <div class="message-content-text">{{ message.content }}</div>
                </div>
              </div>
              
              <div class="load-more" v-if="hasMoreSystem">
                <el-button @click="loadMoreSystem">加载更多</el-button>
              </div>
            </div>
            <div v-else class="empty-state">
              <el-empty description="暂无系统消息" />
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
import { Bell } from '@element-plus/icons-vue'

const router = useRouter()
const activeTab = ref('interaction')

// 模拟互动消息数据
const interactionMessages = ref([
  {
    id: 1,
    sender: {
      id: 101,
      name: '美食达人',
      avatar: 'https://picsum.photos/id/1027/100/100'
    },
    type: 'like',
    content: '赞了你的文章',
    time: new Date(Date.now() - 30 * 60 * 1000), // 30分钟前
    target: {
      type: 'article',
      id: 1,
      title: '如何在家制作完美的提拉米苏',
      preview: '提拉米苏是一道经典的意大利甜点，本文将分享专业大厨的独家秘方...'
    }
  },
  {
    id: 2,
    sender: {
      id: 102,
      name: '旅行笔记',
      avatar: 'https://picsum.photos/id/1012/100/100'
    },
    type: 'comment',
    content: '评论了你的文章：<span class="comment-text">这篇文章写得太棒了，非常实用的建议！</span>',
    time: new Date(Date.now() - 2 * 60 * 60 * 1000), // 2小时前
    target: {
      type: 'article',
      id: 2,
      title: '2025年最值得去的10个小众旅行地',
      preview: '厌倦了人山人海的热门景点？这些鲜为人知的目的地将带给你全新的旅行体验...'
    }
  },
  {
    id: 3,
    sender: {
      id: 103,
      name: '生活家',
      avatar: 'https://picsum.photos/id/1005/100/100'
    },
    type: 'collect',
    content: '收藏了你的文章',
    time: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000), // 1天前
    target: {
      type: 'article',
      id: 3,
      title: '极简主义：如何通过断舍离改变你的生活',
      preview: '极简主义不仅是一种生活方式，更是一种思维模式。本文将分享如何开始你的极简之旅...'
    }
  }
])

// 模拟关注消息数据
const followMessages = ref([
  {
    id: 1,
    sender: {
      id: 104,
      name: '摄影师小王',
      avatar: 'https://picsum.photos/id/1062/100/100'
    },
    content: '关注了你',
    time: new Date(Date.now() - 5 * 60 * 60 * 1000), // 5小时前
    isFollowing: false
  },
  {
    id: 2,
    sender: {
      id: 105,
      name: '健身达人',
      avatar: 'https://picsum.photos/id/1025/100/100'
    },
    content: '关注了你',
    time: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000), // 2天前
    isFollowing: true
  }
])

// 模拟系统消息数据
const systemMessages = ref([
  {
    id: 1,
    content: '您的文章《如何在家制作完美的提拉米苏》已被推荐到首页',
    time: new Date(Date.now() - 12 * 60 * 60 * 1000) // 12小时前
  },
  {
    id: 2,
    content: '您的账号已完成实名认证',
    time: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000) // 5天前
  }
])

const hasMoreInteraction = ref(false)
const hasMoreFollow = ref(false)
const hasMoreSystem = ref(false)

// 格式化时间
const formatTime = (date: Date): string => {
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  // 小于1分钟
  if (diff < 60 * 1000) {
    return '刚刚'
  }
  
  // 小于1小时
  if (diff < 60 * 60 * 1000) {
    return Math.floor(diff / (60 * 1000)) + '分钟前'
  }
  
  // 小于24小时
  if (diff < 24 * 60 * 60 * 1000) {
    return Math.floor(diff / (60 * 60 * 1000)) + '小时前'
  }
  
  // 小于30天
  if (diff < 30 * 24 * 60 * 60 * 1000) {
    return Math.floor(diff / (24 * 60 * 60 * 1000)) + '天前'
  }
  
  // 大于30天
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${year}-${month < 10 ? '0' + month : month}-${day < 10 ? '0' + day : day}`
}

// 导航到目标
const navigateToTarget = (message: any) => {
  if (message.target.type === 'article') {
    router.push(`/article/${message.target.id}`)
  }
}

// 关注用户
const followUser = (userId: number) => {
  // 这里应该调用关注API
  // 模拟关注
  const message = followMessages.value.find(msg => msg.sender.id === userId)
  if (message) {
    message.isFollowing = true
  }
}

// 加载更多互动消息
const loadMoreInteraction = () => {
  // 这里应该调用API加载更多互动消息
  // 模拟加载更多
  hasMoreInteraction.value = false
}

// 加载更多关注消息
const loadMoreFollow = () => {
  // 这里应该调用API加载更多关注消息
  // 模拟加载更多
  hasMoreFollow.value = false
}

// 加载更多系统消息
const loadMoreSystem = () => {
  // 这里应该调用API加载更多系统消息
  // 模拟加载更多
  hasMoreSystem.value = false
}

onMounted(() => {
  // 这里应该调用API获取消息数据
  // 目前使用模拟数据
})
</script>

<style scoped>
.message-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.page-title {
  font-size: 24px;
  margin-bottom: 20px;
  color: var(--text-primary);
}

.message-content {
  background-color: var(--bg-primary);
  border-radius: var(--border-radius-md);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}

.message-tabs {
  margin-bottom: 20px;
}

.message-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.message-item {
  display: flex;
  padding: 15px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.message-avatar {
  margin-right: 15px;
  flex-shrink: 0;
}

.message-body {
  flex: 1;
}

.message-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.sender-name {
  font-weight: 600;
  font-size: 16px;
  color: #333;
}

.message-time {
  font-size: 12px;
  color: #999;
}

.message-content-text {
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
  line-height: 1.5;
}

.comment-text {
  color: #333;
  font-style: italic;
}

.message-target {
  background-color: #f5f5f5;
  border-radius: 4px;
  padding: 10px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.message-target:hover {
  background-color: #f0f0f0;
}

.target-title {
  font-weight: 600;
  font-size: 14px;
  color: #333;
  margin-bottom: 5px;
}

.target-preview {
  font-size: 12px;
  color: #666;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.message-actions {
  margin-top: 10px;
}

.system-message .message-avatar {
  background-color: #f0f0f0;
}

.load-more {
  text-align: center;
  margin-top: 20px;
}

.empty-state {
  text-align: center;
  padding: 40px 0;
}
</style>
