<template>
  <MainLayout>
    <div class="profile-container">
      <div class="profile-header">
        <div class="profile-cover">
          <img src="https://picsum.photos/id/1039/1200/300" alt="Cover" />
        </div>
        
        <div class="profile-info">
          <div class="profile-avatar">
            <el-avatar :size="100" :src="userProfile.avatar">
              {{ userProfile.nickname ? userProfile.nickname.substring(0, 1) : '匿' }}
            </el-avatar>
          </div>
          
          <div class="profile-details">
            <h1 class="profile-name">{{ userProfile.nickname }}</h1>
            <p class="profile-bio">{{ userProfile.aboutMe || '这个人很懒，还没有填写个人简介' }}</p>
            <div class="profile-extra" v-if="userProfile.email || userProfile.birthday">
              <span v-if="userProfile.email">邮箱：{{ userProfile.email }}</span>
              <span v-if="userProfile.birthday" style="margin-left: 16px;">生日：{{ userProfile.birthday }}</span>
            </div>
            
            <div class="profile-stats">
              <div class="stat-item">
                <div class="stat-value">{{ userProfile.articleCount }}</div>
                <div class="stat-label">文章</div>
              </div>
              <div class="stat-item clickable" @click="showFollowers">
                <div class="stat-value">{{ formatNumber(userProfile.followerCount) }}</div>
                <div class="stat-label">粉丝</div>
              </div>
              <div class="stat-item clickable" @click="showFollowing">
                <div class="stat-value">{{ formatNumber(userProfile.followingCount) }}</div>
                <div class="stat-label">关注</div>
              </div>
            </div>
          </div>
          
          <div class="profile-actions" v-if="!isCurrentUser">
            <el-button
              type="primary"
              :plain="userProfile.isFollowing"
              @click="toggleFollow"
            >
              {{ userProfile.isFollowing ? '已关注' : '关注' }}
            </el-button>
          </div>
          
          
        </div>
      </div>
      
      <div class="profile-content">
        <el-tabs v-model="activeTab" class="profile-tabs">
          <el-tab-pane label="文章" name="articles">
            <div v-if="userArticles.length > 0" class="article-list">
              <div v-for="article in userArticles" :key="article.id" class="article-card">
                <div class="article-cover" v-if="article.coverImage" @click="viewArticle(article.id)">
                  <img :src="article.coverImage" :alt="article.title" />
                </div>
                <div class="article-info">
                  <h3 class="article-title" @click="viewArticle(article.id)">{{ article.title }}</h3>
                  <p class="article-abstract">{{ article.abstract }}</p>
                  <div class="article-meta">
                    <div class="article-time">{{ formatDate(article.createTime) }}</div>
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
                    <svg class="icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
                    </svg>
                    {{ formatNumber(article.collectCount) }}
                  </span>
                    </div>
                  </div>
                </div>
              </div>
              
              <div class="load-more" v-if="hasMoreArticles">
                <el-button @click="loadMoreArticles">加载更多</el-button>
              </div>
            </div>
            <div v-else class="empty-state">
              <el-empty description="暂无文章" />
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="收藏" name="collections" v-if="isCurrentUser">
            <div v-if="userCollections.length > 0" class="article-list">
              <div v-for="article in userCollections" :key="article.id" class="article-card">
                <div class="article-cover" v-if="article.coverImage" @click="viewArticle(article.id)">
                  <img :src="article.coverImage" :alt="article.title" />
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
                        <svg class="icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
                        </svg>
                        {{ formatNumber(article.collectCount) }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
              
              <div class="load-more" v-if="hasMoreCollections">
                <el-button @click="loadMoreCollections">加载更多</el-button>
              </div>
            </div>
            <div v-else class="empty-state">
              <el-empty description="暂无收藏" />
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>
    
    <!-- 关注/粉丝弹窗 -->
    <el-dialog
      v-model="showFollowDialog"
      :title="followDialogTitle"
      width="500px"
      @closed="onFollowDialogClosed"
      @close="onFollowDialogClosed"
    >
      <div class="follow-dialog-content">
        <div v-if="followDialogUsers.length > 0" class="user-list">
          <div v-for="user in followDialogUsers" :key="user.id" class="user-item" @click="viewUser(user.id)">
            <div class="user-avatar">
              <el-avatar :size="40" :src="user.avatar">
                {{ user.name ? user.name.substring(0, 1) : '匿' }}
              </el-avatar>
            </div>
            <div class="user-info">
              <div class="user-name">{{ user.name }}</div>
              <div class="user-desc">{{ user.description }}</div>
            </div>
            <div class="user-action" v-if="followDialogTitle === '关注' && isCurrentUser">
              <el-button
                size="small"
                :type="user.isFollowing ? 'info' : 'primary'"
                @click.stop="toggleFollowUser(user)"
              >
                {{ user.isFollowing ? '已关注' : '关注' }}
              </el-button>
            </div>
          </div>
        </div>
        <div v-else class="empty-state">
          <el-empty :description="followDialogEmptyText" />
        </div>
      </div>
    </el-dialog>
    
    
  </MainLayout>
</template>

<script setup lang="ts">
import MainLayout from '@/components/layout/MainLayout.vue'
import { View } from '@element-plus/icons-vue'
import useUserProfileView from '@/scripts/views/UserProfileView'

const {
  userProfile,
  isCurrentUser,
  activeTab,
  userArticles,
  hasMoreArticles,
  userCollections,
  hasMoreCollections,
  showFollowDialog,
  followDialogTitle,
  followDialogUsers,
  followDialogEmptyText,
  currentUserId,
  formatNumber,
  formatDate,
  viewArticle,
  viewUser,
  toggleFollow,
  toggleFollowUser,
  showFollowers,
  showFollowing,
  loadMoreArticles,
  loadMoreCollections
} = useUserProfileView()
</script>

<style lang="scss">
@import '@/styles/views/UserProfileView.scss';
</style>
