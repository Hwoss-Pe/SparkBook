<template>
  <MainLayout>
    <div class="hot-container">
      <div class="hot-header">
        <div class="header-content">
          <div class="title-section">
            <h1 class="hot-title">热门榜单</h1>
            <p class="hot-subtitle">总榜与官方标签榜单</p>
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
            :class="{ active: activeTab === 'overall' }"
            @click="changeTab('overall')"
          >
            总榜
          </div>
          <div
            v-for="tag in officialTags"
            :key="tag"
            class="tab"
            :class="{ active: activeTab === tag }"
            @click="changeTab(tag)"
          >
            {{ tag }}
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
            <div class="article-meta">
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
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import MainLayout from '@/components/layout/MainLayout.vue'
import { View, ChatDotRound } from '@element-plus/icons-vue'
import useHotView from '@/scripts/views/HotView'

const {
  activeTab,
  officialTags,
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
