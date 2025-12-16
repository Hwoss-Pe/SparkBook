<template>
  <MainLayout>
    <div class="creator-center-container">
      <div class="creator-header">
        <h1 class="page-title">创作者中心</h1>
        <div class="header-actions">
          <el-button type="primary" @click="createNewArticle">
            <el-icon><Edit /></el-icon>
            写文章
          </el-button>
        </div>
      </div>
      
      <div class="creator-content">
        <div class="stats-overview">
          <div class="stat-card">
            <div class="stat-icon published">
              <el-icon><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.publishedCount }}</div>
              <div class="stat-label">已发布</div>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-icon draft">
              <el-icon><EditPen /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.draftCount }}</div>
              <div class="stat-label">草稿</div>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-icon views">
              <el-icon><View /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ formatNumber(stats.totalViews) }}</div>
              <div class="stat-label">总阅读</div>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-icon likes">
              <el-icon><Star /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ formatNumber(stats.totalLikes) }}</div>
              <div class="stat-label">总点赞</div>
            </div>
          </div>
        </div>
        
        <div class="content-tabs">
          <el-tabs v-model="activeTab" @tab-change="handleTabChange">
            <el-tab-pane label="已发布" name="published">
              <div class="article-list">
                <div class="list-header">
                  <div class="search-bar">
                    <el-input
                      v-model="searchKeyword"
                      placeholder="搜索文章标题"
                      prefix-icon="Search"
                      clearable
                      @input="handleSearch"
                    />
                  </div>
                  <div class="sort-options">
                    <el-select v-model="sortBy" placeholder="排序方式" @change="handleSort">
                      <el-option label="最新发布" value="publish_time" />
                      <el-option label="最多阅读" value="read_count" />
                      <el-option label="最多点赞" value="like_count" />
                    </el-select>
                  </div>
                </div>
                
                <div v-if="publishedArticles.length > 0" class="articles">
                  <div
                    v-for="article in publishedArticles"
                    :key="article.id"
                    class="article-item"
                  >
                    <div class="article-cover" v-if="article.coverImage">
                      <img :src="article.coverImage" :alt="article.title" />
                    </div>
                    <div class="article-info">
                      <h3 class="article-title" @click="viewArticle(article.id)">
                        {{ article.title }}
                      </h3>
                      <div class="article-flags">
                        <el-tag v-if="article.status === 3" type="warning" size="small">仅自己可见</el-tag>
                      </div>
                      <p class="article-abstract">{{ article.abstract }}</p>
                      <div class="article-meta">
                        <span class="publish-time">{{ formatDate(article.publishTime) }}</span>
                        <div class="article-stats">
                          <span class="stat-item">
                            <el-icon><View /></el-icon>
                            {{ formatNumber(article.readCount) }}
                          </span>
                          <span class="stat-item">
                            <el-icon><Star /></el-icon>
                            {{ formatNumber(article.likeCount) }}
                          </span>
                          <span class="stat-item">
                            <el-icon><ChatDotRound /></el-icon>
                            {{ formatNumber(article.collectCount) }}
                          </span>
                        </div>
                      </div>
                    </div>
                    <div class="article-actions">
                      <el-dropdown trigger="click">
                        <el-button type="text">
                          <el-icon><MoreFilled /></el-icon>
                        </el-button>
                        <template #dropdown>
                          <el-dropdown-menu>
                            <el-dropdown-item @click="editArticle(article.id)">编辑</el-dropdown-item>
                            <el-dropdown-item @click="viewArticle(article.id)">查看</el-dropdown-item>
                            <el-dropdown-item v-if="article.status !== 3" @click="withdrawArticle(article.id)" divided>
                              仅自己可见
                            </el-dropdown-item>
                            <el-dropdown-item @click="unpublishArticle(article.id)" divided>
                              撤回发布（变为草稿）
                            </el-dropdown-item>
                          </el-dropdown-menu>
                        </template>
                      </el-dropdown>
                    </div>
                  </div>
                  <div ref="publishedLoadMoreRef" style="height: 1px;"></div>
                </div>
                <div v-else class="empty-state">
                  <el-empty description="暂无已发布文章">
                    <el-button type="primary" @click="createNewArticle">写第一篇文章</el-button>
                  </el-empty>
                </div>
              </div>
            </el-tab-pane>
            
            <el-tab-pane label="草稿" name="drafts">
              <div class="draft-list">
                <div v-if="draftArticles.length > 0" class="articles">
                  <div
                    v-for="draft in draftArticles"
                    :key="draft.id"
                    class="article-item draft-item"
                  >
                    <div class="article-cover" v-if="draft.coverImage">
                      <img :src="draft.coverImage" :alt="draft.title" />
                    </div>
                    <div class="article-info">
                      <h3 class="article-title" @click="editDraft(draft.id)">
                        {{ draft.title || '无标题草稿' }}
                      </h3>
                      
                      <div class="article-meta">
                        <span class="update-time">{{ formatDate(draft.updateTime) }}</span>
                        <span class="word-count">{{ draft.wordCount || 0 }} 字</span>
                      </div>
                    </div>
                    <div class="article-actions">
                      <el-button type="primary" size="small" @click="editDraft(draft.id)">
                        继续编辑
                      </el-button>
                      <el-dropdown trigger="click">
                        <el-button type="text">
                          <el-icon><MoreFilled /></el-icon>
                        </el-button>
                        <template #dropdown>
                          <el-dropdown-menu>
                            <el-dropdown-item @click="publishDraft(draft.id)">发布</el-dropdown-item>
                            <el-dropdown-item @click="duplicateDraft(draft.id)">复制</el-dropdown-item>
                            <el-dropdown-item @click="deleteDraft(draft.id)" divided>
                              删除
                            </el-dropdown-item>
                          </el-dropdown-menu>
                        </template>
                      </el-dropdown>
                    </div>
                  </div>
                </div>
                <div v-else class="empty-state">
                  <el-empty description="暂无草稿">
                    <el-button type="primary" @click="createNewArticle">开始创作</el-button>
                  </el-empty>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>
    </div>
    
    <!-- 撤回确认弹窗 -->
    <el-dialog
      v-model="showWithdrawDialog"
      title="仅自己可见"
      width="400px"
    >
      <p>设置仅自己可见后文章只有自己可见噢。</p>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showWithdrawDialog = false">取消</el-button>
          <el-button type="primary" @click="confirmWithdraw" :loading="withdrawing">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 删除确认弹窗 -->
    <el-dialog
      v-model="showDeleteDialog"
      title="删除草稿"
      width="400px"
    >
      <p>确定要删除这篇草稿吗？删除后无法恢复。</p>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showDeleteDialog = false">取消</el-button>
          <el-button type="danger" @click="confirmDelete" :loading="deleting">
            确定删除
          </el-button>
        </span>
      </template>
    </el-dialog>
  </MainLayout>
</template>

<script setup lang="ts">
import MainLayout from '@/components/layout/MainLayout.vue'
import {
  Edit,
  Document,
  EditPen,
  View,
  Star,
  ChatDotRound,
  MoreFilled
} from '@element-plus/icons-vue'
import useCreatorCenterView from '@/scripts/views/CreatorCenterView'

const {
  stats,
  activeTab,
  searchKeyword,
  sortBy,
  publishedArticles,
  draftArticles,
  publishedLoadMoreRef,
  showWithdrawDialog,
  showDeleteDialog,
  withdrawing,
  deleting,
  formatNumber,
  formatDate,
  handleTabChange,
  handleSearch,
  handleSort,
  createNewArticle,
  viewArticle,
  editArticle,
  editDraft,
  withdrawArticle,
  publishDraft,
  duplicateDraft,
  deleteDraft,
  confirmWithdraw,
  confirmDelete,
  unpublishArticle
} = useCreatorCenterView()
</script>

<style lang="scss" scoped>
@import '@/styles/views/CreatorCenterView.scss';
</style>
