<template>
  <MainLayout>
    <div class="hot-container">
      <div class="hot-header">
        <div class="header-content">
          <div class="title-section">
            <h1 class="hot-title">热门榜单</h1>
            <p class="hot-subtitle">发现最受欢迎的内容</p>
          </div>
          <div class="action-section">
            <el-button type="primary" @click="triggerRankingCalculation" :loading="isTriggering">
              重新计算热榜
            </el-button>
          </div>
        </div>
      </div>
      
      <div class="tab-container">
        <div class="tabs">
          <div 
            class="tab" 
            :class="{ active: activeTab === 'daily' }"
            @click="changeTab('daily')"
          >
            日榜
          </div>
          <div 
            class="tab" 
            :class="{ active: activeTab === 'weekly' }"
            @click="changeTab('weekly')"
          >
            周榜
          </div>
          <div 
            class="tab" 
            :class="{ active: activeTab === 'monthly' }"
            @click="changeTab('monthly')"
          >
            月榜
          </div>
        </div>
      </div>
      
      <div class="hot-content">
        <div v-for="(article, index) in hotArticles" :key="article.id" class="article-card">
          <div class="article-cover" v-if="article.coverImage" @click="viewArticle(article.id)">
            <img :src="article.coverImage" :alt="article.title" />
            <span class="ranking-number" :class="{ 'top-three': index < 3 }">{{ index + 1 }}</span>
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
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import MainLayout from '@/components/layout/MainLayout.vue'
import { View, ChatDotRound } from '@element-plus/icons-vue'
import { Star as ThumbsUp } from '@element-plus/icons-vue'
import useHotView from '@/scripts/views/HotView'

const {
  activeTab,
  hotArticles,
  hasMoreArticles,
  isTriggering,
  formatNumber,
  viewArticle,
  changeTab,
  loadMoreArticles,
  triggerRankingCalculation
} = useHotView()
</script>

<style lang="scss">
@import '@/styles/views/HotView.scss';
</style>