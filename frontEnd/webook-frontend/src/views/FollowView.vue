<template>
  <MainLayout>
    <div class="follow-container">
      <div class="follow-header">
        <h1 class="follow-title">我的关注</h1>
        <p class="follow-subtitle">查看你关注的作者的最新内容</p>
      </div>
      
      <div class="tab-container">
        <div class="tabs">
          <div 
            class="tab" 
            :class="{ active: activeTab === 'recommend' }"
            @click="changeTab('recommend')"
          >
            推荐
          </div>
          <div 
            class="tab" 
            :class="{ active: activeTab === 'following' }"
            @click="changeTab('following')"
          >
            关注
          </div>
        </div>
      </div>
      
      <div class="follow-content">
        <template v-if="articles.length > 0">
          <div v-for="article in articles" :key="article.id" class="article-card">
            <div class="article-cover" @click="viewArticle(article.id)">
              <img v-if="article.coverImage" :src="article.coverImage" :alt="article.title" />
              <div v-else class="cover-placeholder"></div>
            </div>
            <div class="article-info">
              <h3 class="article-title" @click="viewArticle(article.id)">{{ article.title }}</h3>
              <p class="article-abstract">{{ article.abstract }}</p>
              <div class="article-meta">
                <div class="author-info">
                  <el-avatar :size="24" :src="article.author.avatar">
                    {{ article.author.name ? article.author.name.substring(0, 1) : '匿' }}
                  </el-avatar>
                  <span class="author-name">{{ article.author.name || '匿名用户' }}</span>
                </div>
                <div class="interaction-info">
                  <span class="interaction-item">
                    <el-icon><View /></el-icon>
                    {{ formatNumber(article.readCount) }}
                  </span>
                  <span class="interaction-item">
                    <svg class="icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"></path>
                    </svg>
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
        </template>
        
        <div v-else class="empty-state">
          <el-empty description="暂无关注内容">
            <el-button type="primary" @click="$router.push('/')">去首页发现更多</el-button>
          </el-empty>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import MainLayout from '@/components/layout/MainLayout.vue'
import { View, ChatDotRound } from '@element-plus/icons-vue'
import useFollowView from '@/scripts/views/FollowView'

const {
  activeTab,
  articles,
  hasMoreArticles,
  formatNumber,
  viewArticle,
  changeTab,
  loadMoreArticles
} = useFollowView()
</script>

<style lang="scss">
@import '@/styles/views/FollowView.scss';
</style>
