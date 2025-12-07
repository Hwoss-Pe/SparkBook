<template>
  <MainLayout>
    <div class="search-container">
      <div class="search-header">
        <h1 class="search-title">搜索</h1>
        <p class="search-subtitle">发现你感兴趣的内容和创作者</p>
      </div>
      
      <div class="search-box">
        <el-input
          v-model="searchQuery"
          placeholder="输入关键词搜索"
          class="search-input"
          clearable
          @keyup.enter="performSearch"
        >
          <template #append>
            <el-button :icon="Search" @click="performSearch" :loading="isSearching">
              搜索
            </el-button>
          </template>
        </el-input>
      </div>
      
      <div v-if="totalResults > 0" class="search-stats">
        找到 {{ totalResults }} 个与 "{{ searchQuery }}" 相关的结果
      </div>
      
      <div class="tab-container" v-if="totalResults > 0">
        <div class="tabs">
          <div 
            class="tab" 
            :class="{ active: activeTab === 'article' }"
            @click="changeTab('article')"
          >
            文章 ({{ articles.length }})
          </div>
          <div 
            class="tab" 
            :class="{ active: activeTab === 'user' }"
            @click="changeTab('user')"
          >
            用户 ({{ users.length }})
          </div>
        </div>
      </div>
      
      <div class="search-content" v-if="activeTab === 'article' && articles.length > 0">
        <div v-for="article in articles" :key="article.id" class="article-card">
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
        
        <div class="load-more" v-if="hasMoreResults">
          <el-button @click="loadMoreResults">加载更多</el-button>
        </div>
      </div>
      
      <div class="search-content" v-else-if="activeTab === 'user' && users.length > 0">
        <div v-for="user in users" :key="user.id" class="article-card">
          <div class="article-info">
            <div class="author-info" style="margin-bottom: 15px;">
              <el-avatar :size="50" :src="user.avatar">
                {{ user.name.substring(0, 1) }}
              </el-avatar>
              <div>
                <h3 class="article-title" @click="viewUser(user.id)">{{ user.name }}</h3>
                <p class="article-abstract">{{ user.aboutMe }}</p>
              </div>
            </div>
            <div class="article-meta">
              <div class="interaction-info">
                <span class="interaction-item">
                  <el-icon><User /></el-icon>
                  {{ formatNumber(user.followersCount) }} 粉丝
                </span>
                <span class="interaction-item">
                  <el-icon><Document /></el-icon>
                  {{ formatNumber(user.articlesCount) }} 文章
                </span>
              </div>
              <el-button size="small" @click="viewUser(user.id)">查看主页</el-button>
            </div>
          </div>
        </div>
        
        <div class="load-more" v-if="hasMoreResults">
          <el-button @click="loadMoreResults">加载更多</el-button>
        </div>
      </div>
      
      <div v-else-if="searchQuery && !isSearching" class="empty-state">
        <el-empty :description="`没有找到与 '${searchQuery}' 相关的${activeTab === 'article' ? '文章' : '用户'}`">
          <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
        </el-empty>
      </div>
      
      <div v-else-if="!searchQuery" class="empty-state">
        <el-empty description="请输入关键词进行搜索">
          <p>你可以搜索文章标题、内容或用户名</p>
        </el-empty>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import MainLayout from '@/components/layout/MainLayout.vue'
import { View, ChatDotRound, Search, User, Document } from '@element-plus/icons-vue'
import { Star as ThumbsUp } from '@element-plus/icons-vue'
import useSearchView from '@/scripts/views/SearchView'

const {
  searchQuery,
  activeTab,
  articles,
  users,
  hasMoreResults,
  isSearching,
  totalResults,
  formatNumber,
  viewArticle,
  viewUser,
  changeTab,
  loadMoreResults,
  performSearch
} = useSearchView()
</script>

<style lang="scss">
@import '@/styles/views/SearchView.scss';
</style>