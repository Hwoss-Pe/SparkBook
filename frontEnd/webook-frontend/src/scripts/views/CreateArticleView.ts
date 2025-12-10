import { ref, onMounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { articleApi } from '@/api/article'
import { useUserStore } from '@/stores/user'

// 定义类型接口
interface ArticleForm {
  id: number;
  title: string;
  abstract: string;
  content: string;
  coverImage: string;
  visibility: number; // 1: 公开, 2: 仅自己可见
  tags: string[];
}

interface Draft {
  id: number;
  title: string;
  abstract?: string;
  content?: string;
  coverImage?: string;
  visibility?: number;
  allowComment?: boolean;
  tags?: string[];
  updateTime: string;
}

interface UploadFile {
  raw: File;
}

export default function useCreateArticleView() {
  const router = useRouter()
  const route = useRoute()
  const userStore = useUserStore()

  // 是否为编辑模式
  const isEdit = ref(false)

  // 文章表单
  const articleForm = ref<ArticleForm>({
    id: 0,
    title: '',
    abstract: '',
    content: '',
    coverImage: '',
    visibility: 1, // 1: 公开, 2: 仅自己可见
    tags: []
  })

  // 标签输入
  const inputTagVisible = ref(false)
  const inputTag = ref('')
  const tagInputRef = ref<HTMLInputElement | null>(null)

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
  const saveAsDraft = async () => {
    if (!articleForm.value.title && !articleForm.value.content) {
      ElMessage.warning('文章标题和内容不能同时为空')
      return
    }
    
    try {
      const resp = await articleApi.saveArticle({
        id: articleForm.value.id,
        title: articleForm.value.title,
        content: articleForm.value.content,
        status: 1
      })
      ElMessage.success('草稿保存成功')
      const newId = typeof resp === 'number' ? resp : (resp && typeof resp === 'object' ? (resp as any).id : 0)
      if (!articleForm.value.id && newId) {
        articleForm.value.id = newId
      }
      loadDraftList()
    } catch (error) {
      console.error('保存草稿失败:', error)
      ElMessage.error('草稿保存失败，请稍后重试')
    }
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
  const confirmPublish = async () => {
    publishing.value = true
    
    try {
      const resp = await articleApi.publishArticle({
        id: articleForm.value.id,
        title: articleForm.value.title,
        content: articleForm.value.content
      })
      publishing.value = false
      showPublishDialog.value = false
      ElMessage.success('文章发布成功')
      const articleId = typeof resp === 'number' ? resp : (resp && typeof resp === 'object' ? (resp as any).id : 0)

      // 如果选择仅自己可见，则发布后立即撤回到私有（status=3）
      try {
        if (articleForm.value.visibility === 2 && articleId) {
          userStore.initUserState()
          const uid = userStore.user?.id ?? 0
          if (uid) {
            await articleApi.withdrawArticle(articleId, uid)
            ElMessage.success('已设置为仅自己可见')
          }
        }
      } catch (err) {
        console.error('设置仅自己可见失败:', err)
        ElMessage.error('设置仅自己可见失败，请稍后重试')
      }

      // 发布后跳转详情
      router.push(`/article/${articleId}`)
    } catch (error) {
      console.error('发布文章失败:', error)
      ElMessage.error('文章发布失败，请稍后重试')
      publishing.value = false
    }
  }

  // 加载草稿
  const loadDraft = async (draft: Draft) => {
    try {
      // 调用API获取草稿详情
      const articleDetail = await articleApi.getArticleById(draft.id)
      
      articleForm.value = {
        id: articleDetail.id,
        title: articleDetail.title,
        abstract: '',
        content: articleDetail.content,
        coverImage: draft.coverImage || '',
        visibility: draft.visibility || 1,
        tags: draft.tags || []
      }
      
      ElMessage.success('草稿加载成功')
    } catch (error) {
      console.error('加载草稿失败:', error)
      ElMessage.error('草稿加载失败，请稍后重试')
      
      // 使用本地数据
      articleForm.value = {
        id: draft.id,
        title: draft.title,
        abstract: '',
        content: draft.content || '',
        coverImage: draft.coverImage || '',
        visibility: draft.visibility || 1,
        tags: draft.tags || []
      }
    }
  }

  // 加载草稿列表
  const loadDraftList = async () => {
    try {
      userStore.initUserState()
      const uid = userStore.user?.id ?? 101
      const response = await articleApi.getList({
        offset: 0,
        limit: 100
      })
      const rawDrafts = response.filter(article => article.status === 1)
      const details = await Promise.all(rawDrafts.map(a => articleApi.getArticleById(a.id).catch(() => null)))
      const drafts = details
        .filter(d => d)
        .map(d => ({
          id: d!.id,
          title: d!.title,
          updateTime: d!.utime
        }))
        .sort((a, b) => new Date(b.updateTime).getTime() - new Date(a.updateTime).getTime())
        .slice(0, 5)
      draftList.value = drafts
    } catch (error) {
      console.error('获取草稿列表失败:', error)
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
  }

  // 删除草稿
  const onDeleteDraft = async (id: number) => {
    try {
      userStore.initUserState()
      const uid = userStore.user?.id ?? 0
      if (!uid) {
        ElMessage.error('请先登录再操作')
        return
      }
      await articleApi.deleteDraft(id, uid)
      ElMessage.success('草稿已删除')
      // 如果当前编辑的草稿被删除，清空表单
      if (articleForm.value.id === id) {
        articleForm.value = {
          id: 0,
          title: '',
          abstract: '',
          content: '',
          coverImage: '',
          visibility: 1,
          tags: []
        }
      }
      loadDraftList()
    } catch (error) {
      console.error('删除草稿失败:', error)
      ElMessage.error('删除草稿失败，请稍后重试')
    }
  }

  // 草稿下拉菜单命令
  const onDraftCommand = (payload: any) => {
    try {
      if (!payload || typeof payload !== 'object') return
      const { type, id } = payload as { type: 'edit' | 'delete'; id: number }
      if (type === 'edit') {
        const d = draftList.value.find(x => x.id === id)
        if (d) {
          // 复用现有的加载草稿逻辑
          loadDraft(d)
        }
      } else if (type === 'delete') {
        onDeleteDraft(id)
      }
    } catch (e) {
      console.error('处理草稿菜单命令失败:', e)
    }
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
    const draftId = route.query.draft
    if (articleId) {
      isEdit.value = true
      
      // 调用API获取文章详情
      articleApi.getArticleById(Number(articleId))
        .then(articleDetail => {
          articleForm.value = {
            id: articleDetail.id,
            title: articleDetail.title,
            abstract: articleDetail.abstract || '',
            content: articleDetail.content,
            coverImage: `https://picsum.photos/id/${400 + articleDetail.id}/800/400`, // 实际应该从文章信息中获取
            visibility: 1,
            tags: ['美食', '甜点', '意大利菜'] // 实际应该从文章信息中获取
          }
        })
        .catch(error => {
          console.error('获取文章详情失败:', error)
          
          // 使用模拟数据
          articleForm.value = {
            id: Number(articleId),
            title: '如何在家制作完美的提拉米苏',
            abstract: '提拉米苏是一道经典的意大利甜点，本文将分享专业大厨的独家秘方...',
            content: '提拉米苏（Tiramisu）是一道经典的意大利甜点，起源于威尼托地区。这道甜点的名字在意大利语中意为"带我走"或"振奋我"，这也反映了它独特的风味和口感。\n\n传统的提拉米苏由几个关键成分组成：手指饼干（Ladyfingers）、浓缩咖啡、马斯卡彭奶酪、蛋黄、糖和可可粉。制作过程相对简单，但要做出完美的提拉米苏，需要注意一些关键步骤。\n\n以下是制作完美提拉米苏的详细步骤：\n\n1. 准备材料：\n   - 250克马斯卡彭奶酪\n   - 3个鸡蛋（分离蛋黄和蛋白）\n   - 75克细砂糖\n   - 200克手指饼干\n   - 300毫升浓缩咖啡（冷却）\n   - 2汤匙朗姆酒（可选）\n   - 可可粉适量\n\n2. 制作奶油层：\n   - 将蛋黄和一半的糖放入碗中，用电动打蛋器打至颜色变浅且体积增大。\n   - 加入马斯卡彭奶酪，继续搅拌至均匀。\n   - 在另一个碗中，将蛋白打发至起泡，然后逐渐加入剩余的糖，继续打至形成硬性发泡。\n   - 将打发的蛋白轻轻折叠进奶酪混合物中，注意保持蓬松感。\n\n3. 组装提拉米苏：\n   - 将咖啡和朗姆酒（如果使用）混合在一起。\n   - 快速将手指饼干浸入咖啡混合物中（每一面约1秒），然后排列在容器底部。\n   - 铺上一半的奶油混合物。\n   - 重复上述步骤，形成第二层饼干和奶油。\n   - 最后，在顶部撒上一层可可粉。\n\n4. 冷藏：\n   - 将提拉米苏放入冰箱冷藏至少4小时，最好是过夜，让风味充分融合。\n\n制作提拉米苏的关键在于：\n\n- 不要将手指饼干在咖啡中浸泡太久，否则会变得太软。\n- 确保马斯卡彭奶酪是室温的，这样更容易混合。\n- 轻轻折叠蛋白，保持混合物的蓬松感。\n- 给予足够的冷藏时间，让风味充分发展。\n\n按照这些步骤，你就能在家中制作出口感丰富、层次分明的完美提拉米苏。无论是作为家庭聚餐的甜点，还是招待客人的精致甜品，提拉米苏都是一个绝佳的选择。',
            coverImage: 'https://picsum.photos/id/431/800/600',
            visibility: 1,
            tags: ['美食', '甜点', '意大利菜']
          }
        })
    } else if (draftId) {
      isEdit.value = true
      articleApi.getArticleById(Number(draftId))
        .then(articleDetail => {
          articleForm.value = {
            id: articleDetail.id,
            title: articleDetail.title,
            abstract: '',
            content: articleDetail.content,
            coverImage: articleDetail.coverImage || '',
            visibility: 1,
            tags: []
          }
        })
        .catch(error => {
          console.error('获取草稿详情失败:', error)
          // 保底：只填充ID，其他留空由用户编辑
          articleForm.value = {
            id: Number(draftId),
            title: '',
            abstract: '',
            content: '',
            coverImage: '',
            visibility: 1,
            tags: []
          }
        })
    }
    
    // 加载草稿列表
    loadDraftList()
  })

  return {
    isEdit,
    articleForm,
    inputTagVisible,
    inputTag,
    tagInputRef,
    draftList,
    showPublishDialog,
    publishing,
    showTagInput,
    addTag,
    removeTag,
    handleCoverChange,
    saveAsDraft,
    publishArticle,
    confirmPublish,
    loadDraft,
    onDeleteDraft,
    onDraftCommand,
    formatDate
  }
}
