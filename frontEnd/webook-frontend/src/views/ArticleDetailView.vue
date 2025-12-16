<template>
  <MainLayout>
    <div class="article-detail-container">
      <div v-if="loading" class="loading-state">
        <el-skeleton :rows="15" animated />
      </div>
      
      <div v-else-if="error" class="error-state">
        <el-empty :description="error" />
        <el-button type="primary" @click="goBack">返回</el-button>
      </div>
      
      <template v-else>
        <div class="article-header">
          <h1 class="article-title">{{ article.title }}</h1>
          
          <div class="article-meta">
            <div class="author-info" @click="viewAuthor(article.author.id)">
              <el-avatar :size="40" :src="article.author.avatar">
                {{ article.author.name ? article.author.name.substring(0, 1) : '匿' }}
              </el-avatar>
              <div class="author-details">
                <div class="author-name">{{ article.author.name || '匿名用户' }}</div>
                <div class="publish-time">{{ formatDateTime(article.publishTime) }}</div>
              </div>
            </div>
            
            <div class="action-buttons">
              <el-button
                :type="article.isFollowed ? 'success' : 'primary'"
                size="small" 
                @click="toggleFollow"
              >
                {{ article.isFollowed ? '已关注' : '关注作者' }}
              </el-button>
            </div>
          </div>
        </div>
        
        <div class="article-content">
          <div v-if="article.coverImage" class="article-cover">
            <img :src="article.coverImage" :alt="article.title" />
          </div>
          
          <div class="article-text markdown-body" v-html="renderedContent"></div>
          
          <div class="article-tags">
            <el-tag 
              v-for="tag in article.tags" 
              :key="tag" 
              size="small" 
              effect="light" 
              class="article-tag"
            >
              {{ tag }}
            </el-tag>
          </div>
        </div>
        
        <div class="article-actions">
          <div class="action-item" @click="toggleLike">
            <svg class="icon-svg" :class="{ active: article.isLiked }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"></path>
            </svg>
            <span>{{ formatNumber(article.likeCount) }}</span>
          </div>
          <div class="action-item" @click="toggleFavorite">
            <el-icon :class="{ active: article.isFavorited }"><Star /></el-icon>
            <span>{{ formatNumber(article.favoriteCount) }}</span>
          </div>
          <div class="action-item">
            <el-icon><View /></el-icon>
            <span>{{ formatNumber(article.readCount) }}</span>
          </div>
          <div class="action-item" @click="scrollToComments">
            <el-icon><ChatDotRound /></el-icon>
            <span>{{ formatNumber(article.commentCount) }}</span>
          </div>
          <div class="action-item">
            <el-icon><Share /></el-icon>
            <span>分享</span>
          </div>
        </div>
        
        <div class="article-comments" ref="commentsSection">
          <h2 class="section-title">评论 ({{ article.commentCount }})</h2>
          
          <div class="comment-input">
            <el-input
              v-model="commentContent"
              type="textarea"
              :rows="3"
              placeholder="写下你的评论..."
              maxlength="500"
              show-word-limit
            />
            <div class="comment-submit">
              <el-button type="primary" @click="submitComment" :disabled="!commentContent.trim()">
                发表评论
              </el-button>
            </div>
          </div>
          
          <div class="comment-list" v-if="comments.length > 0">
            <div v-for="comment in comments" :key="comment.id" class="comment-item">
              <div class="comment-user">
                <el-avatar :size="32" :src="comment.user.avatar">
                  {{ comment.user.name ? comment.user.name.substring(0, 1) : '匿' }}
                </el-avatar>
                <div class="comment-user-info">
                  <div class="comment-user-name">{{ comment.user.name }}</div>
                  <div class="comment-time">{{ formatDateTime(comment.createTime) }}</div>
                </div>
              </div>
              <div class="comment-content">{{ comment.content }}</div>
              <div class="comment-actions">
                <span class="comment-like" @click="likeComment(comment)">
                  <svg class="icon-svg" :class="{ active: comment.isLiked, animating: comment.isLikeAnimating }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"></path>
                  </svg>
                  <span>{{ formatNumber(comment.likeCount) }}</span>
                </span>
                <span class="comment-reply" @click="replyToComment(comment)">回复</span>
              </div>
              
              <div v-if="comment.replies && comment.replies.length > 0" class="comment-replies">
                <div v-for="reply in comment.replies" :key="reply.id" class="reply-item">
                  <div class="comment-user">
                    <el-avatar :size="24" :src="reply.user.avatar">
                      {{ reply.user.name ? reply.user.name.substring(0, 1) : '匿' }}
                    </el-avatar>
                    <div class="comment-user-info">
                      <div class="comment-user-name">
                        {{ reply.user.name }}
                        <span class="reply-to" v-if="reply.replyTo">
                          回复 {{ reply.replyTo.name }}
                        </span>
                      </div>
                      <div class="comment-time">{{ formatDateTime(reply.createTime) }}</div>
                    </div>
                  </div>
                  <div class="comment-content">{{ reply.content }}</div>
                  <div class="comment-actions">
                    <span class="comment-like" @click="likeComment(reply)">
                      <svg class="icon-svg" :class="{ active: reply.isLiked, animating: reply.isLikeAnimating }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"></path>
                      </svg>
                      <span>{{ formatNumber(reply.likeCount) }}</span>
                    </span>
                    <span class="comment-reply" @click="replyToComment(reply, comment)">回复</span>
                  </div>
                </div>
                <div class="more-replies" v-if="hasMoreRepliesMap[comment.id]">
                  <el-button size="small" text @click="openThreadModal(comment.id)">查看更多回复</el-button>
                </div>
              </div>
            </div>
          </div>
          
          <div v-else class="empty-comments">
            <el-empty description="暂无评论，来抢沙发吧！" />
          </div>
          
          <div class="load-more" v-if="hasMoreComments">
            <el-button @click="loadMoreComments">加载更多评论</el-button>
          </div>
        </div>
        
        <div class="related-articles">
          <h2 class="section-title">相关推荐</h2>
          <div class="related-list">
            <div 
              v-for="item in relatedArticles" 
              :key="item.id" 
              class="related-item"
              @click="viewArticle(item.id)"
            >
              <div class="related-cover" v-if="item.coverImage">
                <img :src="item.coverImage" :alt="item.title" />
              </div>
              <div class="related-info">
                <h3 class="related-title">{{ item.title }}</h3>
                <div class="related-meta">
                  <span>{{ item.author.name }}</span>
                  <span>{{ formatNumber(item.readCount) }} 阅读</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </MainLayout>
  <el-dialog v-model="threadModal.visible" title="楼中楼回复" width="640px" class="thread-dialog" @close="closeThreadModal">
    <div class="thread-list" v-if="threadModal.visible">
      <div v-for="item in getThreadView(threadModal.rootId)" :key="item.reply.id" class="thread-item" :style="{ marginLeft: (item.depth - 2) * 20 + 'px' }">
        <div class="comment-user">
          <el-avatar :size="24" :src="item.reply.user.avatar">
            {{ item.reply.user.name ? item.reply.user.name.substring(0, 1) : '匿' }}
          </el-avatar>
          <div class="comment-user-info">
            <div class="comment-user-name">
              {{ item.reply.user.name }}
              <span class="reply-to" v-if="item.reply.replyTo">回复 {{ item.reply.replyTo.name }}</span>
            </div>
            <div class="comment-time">{{ formatDateTime(item.reply.createTime) }}</div>
          </div>
        </div>
        <div class="comment-content">{{ item.reply.content }}</div>
        <div class="comment-actions">
          <span class="comment-like" @click="likeComment(item.reply)">
            <svg class="icon-svg" :class="{ active: item.reply.isLiked, animating: item.reply.isLikeAnimating }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"></path>
            </svg>
            <span>{{ formatNumber(item.reply.likeCount) }}</span>
          </span>
          <span class="comment-reply" @click="replyToComment(item.reply)">回复</span>
        </div>
      </div>
      <div class="load-more" v-if="hasMoreRepliesMap[threadModal.rootId]">
        <el-button @click="loadMoreRepliesFor(threadModal.rootId)">加载更多回复</el-button>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import MainLayout from '@/components/layout/MainLayout.vue'
import { View, ChatDotRound, Star, Share } from '@element-plus/icons-vue'
import useArticleDetailView from '@/scripts/views/ArticleDetailView'

const {
  articleId,
  loading,
  error,
  article,
  renderedContent,
  comments,
  hasMoreComments,
  commentContent,
  commentsSection,
  relatedArticles,
  formatNumber,
  formatDateTime,
  goBack,
  viewAuthor,
  viewArticle,
  toggleFollow,
  toggleLike,
  toggleFavorite,
  scrollToComments,
  submitComment,
  likeComment,
  replyToComment,
  hasMoreRepliesMap,
  loadMoreRepliesFor,
  threadModal,
  getThreadView,
  openThreadModal,
  closeThreadModal,
  loadMoreComments
} = useArticleDetailView()
</script>

<style>
@import 'github-markdown-css/github-markdown.css';
</style>

<style lang="scss" scoped>
@import '@/styles/views/ArticleDetailView.scss';
</style>
