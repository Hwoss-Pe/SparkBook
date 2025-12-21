<template>
  <MainLayout>
    <div class="create-container">
      <div class="create-header">
        <h1 class="page-title">{{ isEdit ? '编辑文章' : '创建文章' }}</h1>
        <div class="header-actions">
          <el-button @click="saveAsDraft">保存草稿</el-button>
          <el-button type="primary" @click="publishArticle">发布文章</el-button>
        </div>
      </div>
      
      <div class="create-content">
        <div class="editor-container">
          <el-form :model="articleForm" label-position="top">
            <el-form-item label="文章标题" prop="title">
              <el-input
                v-model="articleForm.title"
                placeholder="请输入文章标题"
                maxlength="100"
                show-word-limit
              />
            </el-form-item>
            
            
            
            <el-form-item label="封面图" prop="coverImage">
              <el-upload
                class="cover-uploader"
                action="#"
                :show-file-list="false"
                :auto-upload="false"
                :on-change="handleCoverChange"
              >
                <img v-if="articleForm.coverImage" :src="articleForm.coverImage" class="cover-image" />
                <div v-else class="cover-placeholder">
                  <el-icon><Plus /></el-icon>
                  <div class="placeholder-text">点击上传封面图</div>
                </div>
              </el-upload>
            </el-form-item>
            
            <el-form-item label="文章内容" prop="content">
              <!-- 这里应该集成富文本编辑器，如TinyMCE、CKEditor等 -->
              <!-- 为了简化示例，这里使用textarea代替 -->
              <el-input
                v-model="articleForm.content"
                type="textarea"
                placeholder="请输入文章内容"
                :rows="15"
              />
            </el-form-item>

            <div class="ai-generate-bar" style="text-align: right; margin-bottom: 18px; display: flex; justify-content: flex-end; gap: 10px;">
              <el-dropdown split-button type="warning" plain :loading="aiPolishLoading" @click="handleAIPolish('优化表达')" @command="handleAIPolish">
                ✒️ AI 润色
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="修复语法错误">修复语法错误</el-dropdown-item>
                    <el-dropdown-item command="扩写这段内容">扩写内容</el-dropdown-item>
                    <el-dropdown-item command="使用更专业的语气">更专业的语气</el-dropdown-item>
                    <el-dropdown-item command="使用更轻松幽默的语气">更幽默的风格</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>

              <el-button 
                type="success" 
                plain 
                :loading="aiTagLoading"
                @click="handleAITag"
              >
                🏷️ AI 自动标签
              </el-button>

              <el-button 
                type="primary" 
                plain 
                :loading="aiLoading"
                @click="handleAIGenerate"
              >
                ✨ AI 生成标题摘要
              </el-button>
            </div>

            <el-form-item label="文章摘要" prop="abstract">
              <el-input
                v-model="articleForm.abstract"
                type="textarea"
                placeholder="请输入或使用AI生成摘要"
                :rows="3"
                maxlength="200"
                show-word-limit
              />
            </el-form-item>
          </el-form>
        </div>
        
        <div class="editor-sidebar">
          <div class="sidebar-section">
            <h3 class="section-title">文章设置</h3>
            <div class="section-content">
              <el-form :model="articleForm" label-position="top">
                <el-form-item label="可见性">
                  <el-radio-group v-model="articleForm.visibility">
                    <el-radio :label="1">公开</el-radio>
                    <el-radio :label="2">仅自己可见</el-radio>
                  </el-radio-group>
                </el-form-item>
                
                
                
                <el-form-item label="标签">
                  <el-tag
                    v-for="tag in articleForm.tags"
                    :key="tag"
                    closable
                    @close="removeTag(tag)"
                    class="article-tag"
                  >
                    {{ tag }}
                  </el-tag>
                  <el-input
                    v-if="inputTagVisible"
                    ref="tagInputRef"
                    v-model="inputTag"
                    class="tag-input"
                    size="small"
                    @keyup.enter="addTag"
                    @blur="addTag"
                  />
                  <el-button v-else class="button-new-tag" size="small" @click="showTagInput">
                    + 新标签
                  </el-button>
                  <div class="official-tags">
                    <span class="official-label">官方标签：</span>
                    <el-tag
                      v-for="name in officialTags"
                      :key="name"
                      class="article-tag"
                      :type="articleForm.tags.includes(name) ? 'success' : 'info'"
                      @click="toggleOfficialTag(name)"
                    >
                      {{ name }}
                    </el-tag>
                  </div>
                </el-form-item>
              </el-form>
            </div>
          </div>
          
          <div class="sidebar-section">
            <h3 class="section-title">草稿箱</h3>
            <div class="section-content">
              <div v-if="draftList.length > 0" class="draft-list">
                <div
                  v-for="draft in draftList"
                  :key="draft.id"
                  class="draft-item"
                >
                  <div class="draft-title" @click="loadDraft(draft)">{{ draft.title || '无标题草稿' }}</div>
                  <div class="draft-time">{{ formatDate(draft.updateTime) }}</div>
                  <div class="draft-actions">
                  </div>
                </div>
              </div>
              <div v-else class="empty-draft">
                <el-empty description="暂无草稿" :image-size="80" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 确认发布弹窗 -->
    <el-dialog
      v-model="showPublishDialog"
      title="发布文章"
      width="400px"
    >
      <p>确定要发布这篇文章吗？</p>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showPublishDialog = false">取消</el-button>
          <el-button type="primary" @click="confirmPublish" :loading="publishing">
            确定发布
          </el-button>
        </span>
      </template>
    </el-dialog>
  </MainLayout>
</template>

<script setup lang="ts">
import MainLayout from '@/components/layout/MainLayout.vue'
import { Plus, More } from '@element-plus/icons-vue'
import useCreateArticleView from '@/scripts/views/CreateArticleView'

const {
  isEdit,
  articleForm,
  inputTagVisible,
  inputTag,
  tagInputRef,
  draftList,
  showPublishDialog,
  publishing,
  officialTags,
  showTagInput,
  addTag,
  toggleOfficialTag,
  removeTag,
  handleCoverChange,
  saveAsDraft,
  publishArticle,
  confirmPublish,
  loadDraft,
  onDeleteDraft,
  formatDate,
  aiLoading,
  handleAIGenerate,
  aiPolishLoading,
  aiTagLoading,
  handleAITag,
  handleAIPolish
} = useCreateArticleView()
</script>

<style lang="scss" scoped>
@import '@/styles/views/CreateArticleView.scss';
</style>
