<template>
  <MainLayout>
    <div class="message-container">
      <h1 class="page-title">消息中心</h1>
      
      <div class="message-content">
        <el-tabs v-model="activeTab" class="message-tabs">
          <el-tab-pane name="interaction">
            <template #label>
              互动消息
              <el-badge v-if="unread.interaction > 0" :value="unread.interaction" class="tab-badge" />
            </template>
            <div v-if="interactionMessages.length > 0" class="message-list">
              <div v-for="message in interactionMessages" :key="message.id" class="message-item clickable" @click="navigateToTarget(message)">
                <div class="message-avatar">
                  <el-avatar :size="40" :src="message.sender.avatar">
                    {{ message.sender.name ? message.sender.name.substring(0, 1) : '匿' }}
                  </el-avatar>
                </div>
                <div class="message-body">
                  <div class="message-header">
                    <span class="sender-name">{{ message.sender.name }}</span>
                    <span class="message-time">{{ formatTime(message.time) }}</span>
                  </div>
                  <div class="message-content-text" v-html="message.content"></div>
                  <div class="message-target" v-if="message.target && (message.target.title || message.target.preview)">
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
          
          <el-tab-pane name="follow">
            <template #label>
              关注消息
              <el-badge v-if="unread.follow > 0" :value="unread.follow" class="tab-badge" />
            </template>
            <div v-if="followMessages.length > 0" class="message-list">
              <div v-for="message in followMessages" :key="message.id" class="message-item clickable" @click="navigateToUser(message)">
                <div class="message-avatar">
                  <el-avatar :size="40" :src="message.sender.avatar">
                    {{ message.sender.name ? message.sender.name.substring(0, 1) : '匿' }}
                  </el-avatar>
                </div>
                <div class="message-body">
                  <div class="message-header">
                    <span class="sender-name">{{ message.sender.name }}</span>
                    <span class="message-time">{{ formatTime(message.time) }}</span>
                  </div>
                  <div class="message-content-text">{{ message.content }}</div>
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
          
          <el-tab-pane name="system">
            <template #label>
              系统消息
              <el-badge v-if="unread.system > 0" :value="unread.system" class="tab-badge" />
            </template>
            <div v-if="systemMessages.length > 0" class="message-list">
              <div v-for="message in systemMessages" :key="message.id" class="message-item system-message">
                <div class="system-icon">
                  <el-icon><Bell /></el-icon>
                </div>
                <div class="message-body">
                  <div class="message-header system-header">
                    <span class="sender-name">系统通知</span>
                    <span class="message-time">{{ formatTime(message.time) }}</span>
                  </div>
                  <div class="system-content">
                    <div class="system-text">{{ message.content }}</div>
                  </div>
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
import { ref, onMounted, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import MainLayout from '@/components/layout/MainLayout.vue'
import { Bell } from '@element-plus/icons-vue'
import { useNotificationStore } from '@/stores/notification'

const router = useRouter()
const activeTab = ref('interaction')
const store = useNotificationStore()
const unread = store.unread
const interactionMessages = computed(() => store.interaction)
const followMessages = computed(() => store.follow)
const systemMessages = computed(() => store.system)
const hasMoreInteraction = computed(() => store.hasMore.interaction)
const hasMoreFollow = computed(() => store.hasMore.follow)
const hasMoreSystem = computed(() => store.hasMore.system)

 

// 格式化时间
const formatTime = (dateLike: any): string => {
  const date = typeof dateLike === 'string' ? new Date(dateLike) : dateLike
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
  if (message && message.target && message.target.type === 'article' && message.target.id) {
    router.push(`/article/${message.target.id}`)
  }
}

// 跳转到对方个人主页（关注消息）
const navigateToUser = (message: any) => {
  if (message && message.sender && message.sender.id) {
    router.push(`/user/${message.sender.id}`)
  }
}

// 关注用户
// 关注消息不在此进行操作，点击卡片直接跳转个人主页

const loadMoreInteraction = () => { store.fetchMessages('interaction', false, 10) }
const loadMoreFollow = () => { store.fetchMessages('follow', false, 10) }
const loadMoreSystem = () => { store.fetchMessages('system', false, 10) }

onMounted(async () => {
  await store.fetchUnreadCounts()
  await store.fetchMessages('interaction', true, 10)
})

watch(activeTab, async (tab) => {
  if (tab === 'interaction' || tab === 'follow' || tab === 'system') {
    await store.markReadByCategory(tab as any)
    await store.fetchMessages(tab as any, true, 10)
  }
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

.message-item.clickable {
  cursor: pointer;
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

.system-message {
  background: linear-gradient(135deg, #f8fbff 0%, #ffffff 60%);
  border-left: 4px solid #409EFF;
}

.system-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rgba(64, 158, 255, 0.15);
  color: #409EFF;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 15px;
}

.system-header .sender-name {
  color: #409EFF;
  font-weight: 700;
}

.system-content {
  background-color: rgba(64, 158, 255, 0.08);
  border: 1px solid rgba(64, 158, 255, 0.15);
  border-radius: 6px;
  padding: 10px 12px;
}

.system-text {
  color: #3a3a3a;
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
