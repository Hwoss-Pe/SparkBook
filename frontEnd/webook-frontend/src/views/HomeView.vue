<template>
  <MainLayout>
    <div class="home-container">
      <div class="banner">
        <el-carousel :interval="4000" type="card" height="300px">
          <el-carousel-item v-for="(item, index) in bannerItems" :key="index">
            <div class="banner-item" :style="{ backgroundImage: `url(${item.image})` }">
              <div class="banner-content">
                <h3>{{ item.title }}</h3>
                <p>{{ item.description }}</p>
              </div>
            </div>
          </el-carousel-item>
        </el-carousel>
      </div>

      <div class="tab-navigation">
        <div class="tab-container">
          <div class="tab active">综合榜</div>
          <div class="tab">日榜</div>
          <div class="tab">周榜</div>
          <div class="tab">月榜</div>
        </div>
      </div>

      <div class="content-container">
        <div class="article-list">
          <h2 class="section-title">推荐内容</h2>
          <div class="waterfall-container">
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
          </div>
          <div class="load-more" v-if="hasMoreArticles">
            <el-button @click="loadMoreArticles">加载更多</el-button>
          </div>
        </div>

        <div class="sidebar-content">
          <div class="hot-ranking">
            <h2 class="section-title">热门榜单</h2>
            <div class="ranking-list">
              <div 
                v-for="(item, index) in hotRankings" 
                :key="item.id" 
                class="ranking-item"
                @click="viewArticle(item.id)"
              >
                <div class="ranking-number" :class="{ 'top-three': index < 3 }">{{ index + 1 }}</div>
                <div class="ranking-content">
                  <h4 class="ranking-title">{{ item.title }}</h4>
                  <div class="ranking-meta">
                    <span>{{ item.author.name }}</span>
                    <span>{{ formatNumber(item.readCount) }} 阅读</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="recommended-authors">
            <h2 class="section-title">推荐作者</h2>
            <div class="author-list">
              <div v-for="author in recommendedAuthors" :key="author.id" class="author-item">
                <el-avatar :size="40" :src="author.avatar">
                  {{ author.name.substring(0, 1) }}
                </el-avatar>
                <div class="author-info">
                  <div class="author-name">{{ author.name }}</div>
                  <div class="author-desc">{{ author.description }}</div>
                </div>
                <el-button size="small" type="primary" plain @click="followAuthor(author.id)">
                  {{ author.isFollowed ? '已关注' : '关注' }}
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import MainLayout from '@/components/layout/MainLayout.vue'
import { View, ChatDotRound } from '@element-plus/icons-vue'
import { Star as ThumbsUp } from '@element-plus/icons-vue'
import useHomeView from '@/scripts/views/HomeView'

const {
  bannerItems,
  articles,
  hotRankings,
  recommendedAuthors,
  hasMoreArticles,
  formatNumber,
  viewArticle,
  loadMoreArticles,
  followAuthor
} = useHomeView()
</script>

<style lang="scss">
@import '@/styles/views/HomeView.scss';
</style>