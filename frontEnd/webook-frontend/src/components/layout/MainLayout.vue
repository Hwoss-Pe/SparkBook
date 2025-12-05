<template>
  <div class="main-layout">
    <header class="header">
      <div class="logo-container">
        <h1 class="logo">小微书</h1>
      </div>
      <div class="search-container">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索文章、用户"
          prefix-icon="el-icon-search"
          clearable
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
      <div class="user-container">
        <template v-if="isLoggedIn">
          <el-dropdown trigger="click">
            <div class="user-avatar">
              <el-avatar :size="32" :src="userAvatar">{{ userInitials }}</el-avatar>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="navigateTo('/user/profile')">个人主页</el-dropdown-item>
                <el-dropdown-item @click="navigateTo('/user/settings')">设置</el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
        <template v-else>
          <el-button type="primary" @click="navigateTo('/login')">登录/注册</el-button>
        </template>
      </div>
    </header>
    
    <main class="main-content">
      <div class="sidebar">
        <el-menu
          :default-active="activeMenu"
          class="main-menu"
          router
        >
          <el-menu-item index="/">
            <el-icon><House /></el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="/follow">
            <el-icon><Star /></el-icon>
            <span>关注</span>
          </el-menu-item>
          <el-menu-item index="/hot">
            <el-icon><Histogram /></el-icon>
            <span>热榜</span>
          </el-menu-item>
          <el-menu-item index="/message">
            <el-icon><ChatDotRound /></el-icon>
            <span>消息</span>
          </el-menu-item>
          <el-menu-item index="/create" v-if="isLoggedIn">
            <el-icon><Edit /></el-icon>
            <span>创作中心</span>
          </el-menu-item>
        </el-menu>
      </div>
      
      <div class="content">
        <slot></slot>
      </div>
    </main>
    
    <footer class="footer">
      <p>© 2025 小微书 - 一个类似小红书的内容平台</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue' // 修复导入
import { useRouter } from 'vue-router'
import { Search, House, Star, Histogram, ChatDotRound, Edit } from '@element-plus/icons-vue'

// 模拟用户登录状态
const isLoggedIn = ref(false)
const userAvatar = ref('')
const userName = ref('')

const userInitials = computed(() => {
  return userName.value ? userName.value.substring(0, 1).toUpperCase() : 'U'
})

const router = useRouter()
const searchKeyword = ref('')
const activeMenu = ref('/')

// 处理搜索
const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push({
      path: '/search',
      query: { q: searchKeyword.value }
    })
  }
}

// 页面导航
const navigateTo = (path: string) => {
  router.push(path)
}

// 处理登出
const handleLogout = () => {
  isLoggedIn.value = false
  // 这里应该调用登出API，清除token等
  router.push('/login')
}
</script>

<style scoped>
.main-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.header {
  display: flex;
  align-items: center;
  padding: 0 20px;
  height: 60px;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
}

.logo-container {
  flex: 0 0 200px;
}

.logo {
  margin: 0;
  font-size: 24px;
  color: #ff2442; /* 小红书风格的红色 */
}

.search-container {
  flex: 1;
  max-width: 500px;
  margin: 0 20px;
}

.user-container {
  flex: 0 0 auto;
}

.user-avatar {
  cursor: pointer;
}

.main-content {
  display: flex;
  margin-top: 60px;
  flex: 1;
}

.sidebar {
  flex: 0 0 200px;
  background-color: #fff;
  border-right: 1px solid #eee;
  position: fixed;
  height: calc(100vh - 60px);
  top: 60px;
}

.main-menu {
  border-right: none;
}

.content {
  flex: 1;
  margin-left: 200px;
  padding: 20px;
}

.footer {
  background-color: #f8f8f8;
  padding: 20px;
  text-align: center;
  color: #666;
}
</style>
