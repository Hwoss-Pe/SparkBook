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
                    <el-avatar :size="32" :src="article.author.avatar" @click="viewUser(article.author.id)">
                      {{ article.author.name ? article.author.name.substring(0, 1) : '匿' }}
                    </el-avatar>
                    <span class="author-name" @click="viewUser(article.author.id)">{{ article.author.name || '匿名用户' }}</span>
                  </div>
                  <div class="interaction-info">
                    <span class="interaction-item">
                      <svg class="icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                        <circle cx="12" cy="12" r="3"></circle>
                      </svg>
                      {{ formatNumber(article.readCount) }}
                    </span>
                    <span class="interaction-item like" :class="{ active: article.isLiked, animating: article.isLikeAnimating }" @click="toggleArticleLike(article)">
                      <svg class="icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"></path>
                      </svg>
                      {{ formatNumber(article.likeCount) }}
                    </span>
                    <span class="interaction-item collect" :class="{ active: article.isFavorited, animating: article.isFavAnimating }" @click="toggleArticleFavorite(article)">
                      <svg class="icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
                      </svg>
                      {{ formatNumber(article.collectCount) }}
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
                <el-avatar :size="40" :src="author.avatar" @click="viewUser(author.id)">
                  {{ author.name ? author.name.substring(0, 1) : '匿' }}
                </el-avatar>
                <div class="author-info">
                  <div class="author-name" @click="viewUser(author.id)">{{ author.name }}</div>
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
import useHomeView from '@/scripts/views/HomeView'

const {
  bannerItems,
  articles,
  hotRankings,
  recommendedAuthors,
  hasMoreArticles,
  formatNumber,
  viewArticle,
  viewUser,
  loadMoreArticles,
  followAuthor,
  toggleArticleLike,
  toggleArticleFavorite
} = useHomeView()
</script>

<style lang="scss" scoped>
@import '@/styles/views/HomeView.scss';
</style>
