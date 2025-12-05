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
            
            <el-form-item label="文章摘要" prop="abstract">
              <el-input
                v-model="articleForm.abstract"
                type="textarea"
                placeholder="请输入文章摘要"
                maxlength="200"
                show-word-limit
                :rows="3"
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
                
                <el-form-item label="允许评论">
                  <el-switch v-model="articleForm.allowComment" />
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
                  @click="loadDraft(draft)"
                >
                  <div class="draft-title">{{ draft.title || '无标题草稿' }}</div>
                  <div class="draft-time">{{ formatDate(draft.updateTime) }}</div>
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
import { ref, onMounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import MainLayout from '@/components/layout/MainLayout.vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()

// 是否为编辑模式
const isEdit = ref(false)

interface ArticleForm {
  id: number;
  title: string;
  abstract: string;
  content: string;
  coverImage: string;
  visibility: number; // 1: 公开, 2: 仅自己可见
  allowComment: boolean;
  tags: string[];
}

// 文章表单
const articleForm = ref<ArticleForm>({
  id: 0,
  title: '',
  abstract: '',
  content: '',
  coverImage: '',
  visibility: 1, // 1: 公开, 2: 仅自己可见
  allowComment: true,
  tags: []
})

// 标签输入
const inputTagVisible = ref(false)
const inputTag = ref('')
const tagInputRef = ref<HTMLInputElement | null>(null)

interface Draft {
  id: number;
  title: string;
  updateTime: string;
  abstract?: string;
  content?: string;
  coverImage?: string;
  visibility?: number;
  allowComment?: boolean;
  tags?: string[];
}

// 草稿列表
const draftList = ref<Draft[]>([])

// 发布确认弹窗
const showPublishDialog = ref(false)
const publishing = ref(false)

// 显示标签输入框
const showTagInput = () => {
  inputTagVisible.value = true
  nextTick(() => {
    if (tagInputRef.value) {
      tagInputRef.value.focus()
    }
  })
}

// 添加标签
const addTag = () => {
  if (inputTag.value) {
    if (articleForm.value.tags.length < 5) {
      if (!articleForm.value.tags.includes(inputTag.value)) {
        articleForm.value.tags.push(inputTag.value)
      } else {
        ElMessage.warning('标签已存在')
      }
    } else {
      ElMessage.warning('最多添加5个标签')
    }
  }
  inputTagVisible.value = false
  inputTag.value = ''
}

// 移除标签
const removeTag = (tag: string) => {
  articleForm.value.tags = articleForm.value.tags.filter(t => t !== tag)
}

interface UploadFile {
  raw: File;
}

// 处理封面图变更
const handleCoverChange = (file: UploadFile) => {
  // 这里应该上传图片到服务器
  // 模拟上传成功
  const reader = new FileReader()
  reader.readAsDataURL(file.raw)
  reader.onload = () => {
    articleForm.value.coverImage = reader.result as string
  }
}

// 保存为草稿
const saveAsDraft = () => {
  if (!articleForm.value.title && !articleForm.value.content) {
    ElMessage.warning('文章标题和内容不能同时为空')
    return
  }
  
  // 这里应该调用API保存草稿
  // 模拟保存成功
  ElMessage.success('草稿保存成功')
  
  // 如果是新文章，保存后应该获取ID并更新表单
  if (!isEdit.value && !articleForm.value.id) {
    articleForm.value.id = Date.now() // 模拟生成ID
  }
  
  // 更新草稿列表
  loadDraftList()
}

// 发布文章
const publishArticle = () => {
  if (!articleForm.value.title) {
    ElMessage.warning('请输入文章标题')
    return
  }
  
  if (!articleForm.value.content) {
    ElMessage.warning('请输入文章内容')
    return
  }
  
  if (!articleForm.value.coverImage) {
    ElMessage.warning('请上传封面图')
    return
  }
  
  showPublishDialog.value = true
}

// 确认发布
const confirmPublish = () => {
  publishing.value = true
  
  // 这里应该调用API发布文章
  // 模拟发布成功
  setTimeout(() => {
    publishing.value = false
    showPublishDialog.value = false
    ElMessage.success('文章发布成功')
    
    // 发布成功后跳转到文章详情页
    router.push(`/article/${articleForm.value.id}`)
  }, 1000)
}

// 加载草稿
const loadDraft = (draft: Draft) => {
  // 这里应该调用API获取草稿详情
  // 模拟获取草稿详情
  articleForm.value = {
    id: draft.id,
    title: draft.title,
    abstract: draft.abstract || '',
    content: draft.content || '',
    coverImage: draft.coverImage || '',
    visibility: draft.visibility || 1,
    allowComment: draft.allowComment !== false,
    tags: draft.tags || []
  }
  
  ElMessage.success('草稿加载成功')
}

// 加载草稿列表
const loadDraftList = () => {
  // 这里应该调用API获取草稿列表
  // 模拟草稿列表
  draftList.value = [
    {
      id: 1,
      title: '极简主义生活方式探索',
      updateTime: '2024-11-28T14:20:00'
    },
    {
      id: 2,
      title: '2025年旅行计划',
      updateTime: '2024-11-25T09:15:00'
    },
    {
      id: 3,
      title: '',
      updateTime: '2024-11-20T16:40:00'
    }
  ]
}

// 格式化日期
const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  // 检查是否为编辑模式
  const articleId = route.query.id
  if (articleId) {
    isEdit.value = true
    
    // 这里应该调用API获取文章详情
    // 模拟获取文章详情
    articleForm.value = {
      id: Number(articleId),
      title: '如何在家制作完美的提拉米苏',
      abstract: '提拉米苏是一道经典的意大利甜点，本文将分享专业大厨的独家秘方...',
      content: '提拉米苏（Tiramisu）是一道经典的意大利甜点，起源于威尼托地区。这道甜点的名字在意大利语中意为"带我走"或"振奋我"，这也反映了它独特的风味和口感。\n\n传统的提拉米苏由几个关键成分组成：手指饼干（Ladyfingers）、浓缩咖啡、马斯卡彭奶酪、蛋黄、糖和可可粉。制作过程相对简单，但要做出完美的提拉米苏，需要注意一些关键步骤。\n\n以下是制作完美提拉米苏的详细步骤：\n\n1. 准备材料：\n   - 250克马斯卡彭奶酪\n   - 3个鸡蛋（分离蛋黄和蛋白）\n   - 75克细砂糖\n   - 200克手指饼干\n   - 300毫升浓缩咖啡（冷却）\n   - 2汤匙朗姆酒（可选）\n   - 可可粉适量\n\n2. 制作奶油层：\n   - 将蛋黄和一半的糖放入碗中，用电动打蛋器打至颜色变浅且体积增大。\n   - 加入马斯卡彭奶酪，继续搅拌至均匀。\n   - 在另一个碗中，将蛋白打发至起泡，然后逐渐加入剩余的糖，继续打至形成硬性发泡。\n   - 将打发的蛋白轻轻折叠进奶酪混合物中，注意保持蓬松感。\n\n3. 组装提拉米苏：\n   - 将咖啡和朗姆酒（如果使用）混合在一起。\n   - 快速将手指饼干浸入咖啡混合物中（每一面约1秒），然后排列在容器底部。\n   - 铺上一半的奶油混合物。\n   - 重复上述步骤，形成第二层饼干和奶油。\n   - 最后，在顶部撒上一层可可粉。\n\n4. 冷藏：\n   - 将提拉米苏放入冰箱冷藏至少4小时，最好是过夜，让风味充分融合。\n\n制作提拉米苏的关键在于：\n\n- 不要将手指饼干在咖啡中浸泡太久，否则会变得太软。\n- 确保马斯卡彭奶酪是室温的，这样更容易混合。\n- 轻轻折叠蛋白，保持混合物的蓬松感。\n- 给予足够的冷藏时间，让风味充分发展。\n\n按照这些步骤，你就能在家中制作出口感丰富、层次分明的完美提拉米苏。无论是作为家庭聚餐的甜点，还是招待客人的精致甜品，提拉米苏都是一个绝佳的选择。',
      coverImage: 'https://picsum.photos/id/431/800/600',
      visibility: 1,
      allowComment: true,
      tags: ['美食', '甜点', '意大利菜']
    }
  }
  
  // 加载草稿列表
  loadDraftList()
})
</script>

<style scoped>
.create-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.create-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  font-size: 24px;
  margin: 0;
  color: var(--text-primary);
}

.header-actions {
  display: flex;
  gap: 10px;
}

.create-content {
  display: flex;
  gap: 20px;
}

.editor-container {
  flex: 1;
  background-color: var(--bg-primary);
  border-radius: var(--border-radius-md);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}

.editor-sidebar {
  width: 300px;
  flex-shrink: 0;
}

.sidebar-section {
  background-color: var(--bg-primary);
  border-radius: var(--border-radius-md);
  padding: 20px;
  box-shadow: var(--shadow-sm);
  margin-bottom: 20px;
}

.section-title {
  font-size: 18px;
  margin: 0 0 15px;
  color: var(--text-primary);
}

.cover-uploader {
  width: 100%;
  height: 200px;
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  display: flex;
  justify-content: center;
  align-items: center;
}

.cover-uploader:hover {
  border-color: var(--primary-color);
}

.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cover-placeholder {
  text-align: center;
  color: #8c939d;
}

.placeholder-text {
  margin-top: 8px;
}

.article-tag {
  margin-right: 10px;
  margin-bottom: 10px;
}

.tag-input {
  width: 100px;
  margin-right: 10px;
  vertical-align: bottom;
}

.draft-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.draft-item {
  padding: 10px;
  border-radius: 4px;
  background-color: #f5f5f5;
  cursor: pointer;
  transition: background-color 0.3s;
}

.draft-item:hover {
  background-color: #f0f0f0;
}

.draft-title {
  font-weight: 600;
  margin-bottom: 5px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.draft-time {
  font-size: 12px;
  color: #999;
}

.empty-draft {
  text-align: center;
  padding: 20px 0;
}

@media (max-width: 768px) {
  .create-content {
    flex-direction: column;
  }
  
  .editor-sidebar {
    width: 100%;
  }
}
</style>
