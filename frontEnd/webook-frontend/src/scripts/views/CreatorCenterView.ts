import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { articleApi, type Article } from '@/api/article'

// 定义接口
interface CreatorStats {
  publishedCount: number
  draftCount: number
  totalViews: number
  totalLikes: number
}

interface PublishedArticle {
  id: number
  title: string
  abstract: string
  coverImage: string
  publishTime: string
  readCount: number
  likeCount: number
  commentCount: number
  status: number
}

interface DraftArticle {
  id: number
  title: string
  abstract: string
  coverImage: string
  updateTime: string
  wordCount: number
  status: number
}

export default function useCreatorCenterView() {
  const router = useRouter()
  const userStore = useUserStore()
  
  // 响应式数据
  const stats = ref<CreatorStats>({
    publishedCount: 0,
    draftCount: 0,
    totalViews: 0,
    totalLikes: 0
  })
  
  const activeTab = ref('published')
  const searchKeyword = ref('')
  const sortBy = ref('publish_time')
  
  const publishedArticles = ref<PublishedArticle[]>([])
  const draftArticles = ref<DraftArticle[]>([])
  const allPublishedArticles = ref<PublishedArticle[]>([])
  
  // 弹窗状态
  const showWithdrawDialog = ref(false)
  const showDeleteDialog = ref(false)
  const withdrawing = ref(false)
  const deleting = ref(false)
  const currentArticleId = ref<number | null>(null)
  
  // 计算属性 - 过滤后的已发布文章
  const filteredPublishedArticles = computed(() => {
    let articles = [...allPublishedArticles.value]
    
    // 搜索过滤
    if (searchKeyword.value) {
      articles = articles.filter(article =>
        article.title.toLowerCase().includes(searchKeyword.value.toLowerCase())
      )
    }
    
    // 排序
    articles.sort((a, b) => {
      switch (sortBy.value) {
        case 'read_count':
          return b.readCount - a.readCount
        case 'like_count':
          return b.likeCount - a.likeCount
        case 'publish_time':
        default:
          return new Date(b.publishTime).getTime() - new Date(a.publishTime).getTime()
      }
    })
    
    return articles
  })
  
  // 格式化数字
  const formatNumber = (num: number): string => {
    if (num >= 10000) {
      return (num / 10000).toFixed(1) + 'w'
    } else if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'k'
    }
    return num.toString()
  }
  
  // 格式化日期
  const formatDate = (dateString: string): string => {
    const date = new Date(dateString)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))
    
    if (days === 0) {
      const hours = Math.floor(diff / (1000 * 60 * 60))
      if (hours === 0) {
        const minutes = Math.floor(diff / (1000 * 60))
        return minutes <= 0 ? '刚刚' : `${minutes}分钟前`
      }
      return `${hours}小时前`
    } else if (days < 7) {
      return `${days}天前`
    } else {
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      })
    }
  }
  
  // 标签页切换
  const handleTabChange = (tabName: string) => {
    activeTab.value = tabName
    if (tabName === 'published') {
      loadPublishedArticles()
    } else if (tabName === 'drafts') {
      loadDraftArticles()
    }
  }
  
  // 搜索处理
  const handleSearch = () => {
    publishedArticles.value = filteredPublishedArticles.value
  }
  
  // 排序处理
  const handleSort = () => {
    publishedArticles.value = filteredPublishedArticles.value
  }
  
  // 创建新文章
  const createNewArticle = () => {
    router.push('/create/article')
  }
  
  // 查看文章
  const viewArticle = (id: number) => {
    router.push(`/article/${id}`)
  }
  
  // 编辑文章
  const editArticle = (id: number) => {
    router.push(`/create/article?id=${id}`)
  }
  
  // 编辑草稿
  const editDraft = (id: number) => {
    router.push(`/create/article?draft=${id}`)
  }
  
  // 撤回文章
  const withdrawArticle = (id: number) => {
    currentArticleId.value = id
    showWithdrawDialog.value = true
  }
  
  // 确认撤回
  const confirmWithdraw = async () => {
    if (!currentArticleId.value) return
    
    withdrawing.value = true
    try {
      await articleApi.withdrawArticle(currentArticleId.value, userStore.user?.id || 0)
      
      // 从已发布列表中移除
      publishedArticles.value = publishedArticles.value.filter(
        article => article.id !== currentArticleId.value
      )
      allPublishedArticles.value = allPublishedArticles.value.filter(
        article => article.id !== currentArticleId.value
      )
      
      // 更新统计数据
      stats.value.publishedCount--
      stats.value.draftCount++
      
      ElMessage.success('文章已撤回')
      showWithdrawDialog.value = false
      
      // 重新加载草稿列表
      if (activeTab.value === 'drafts') {
        loadDraftArticles()
      }
    } catch (error) {
      console.error('撤回文章失败:', error)
      ElMessage.error('撤回失败，请重试')
    } finally {
      withdrawing.value = false
    }
  }
  
  // 发布草稿
  const publishDraft = async (id: number) => {
    try {
      await articleApi.publishArticle({ id })
      
      // 从草稿列表中移除
      draftArticles.value = draftArticles.value.filter(draft => draft.id !== id)
      
      // 更新统计数据
      stats.value.draftCount--
      stats.value.publishedCount++
      
      ElMessage.success('草稿已发布')
      
      // 重新加载已发布列表
      if (activeTab.value === 'published') {
        loadPublishedArticles()
      }
    } catch (error) {
      console.error('发布草稿失败:', error)
      ElMessage.error('发布失败，请重试')
    }
  }
  
  // 复制草稿
  const duplicateDraft = async (id: number) => {
    try {
      const draft = draftArticles.value.find(d => d.id === id)
      if (!draft) return
      
      // 创建新草稿
      const newDraft = {
        title: `${draft.title} - 副本`,
        abstract: draft.abstract,
        content: '', // 需要从API获取完整内容
        status: 0
      }
      
      const response = await articleApi.saveArticle(newDraft)
      ElMessage.success('草稿已复制')
      
      // 重新加载草稿列表
      loadDraftArticles()
    } catch (error) {
      console.error('复制草稿失败:', error)
      ElMessage.error('复制失败，请重试')
    }
  }
  
  // 删除草稿
  const deleteDraft = (id: number) => {
    currentArticleId.value = id
    showDeleteDialog.value = true
  }
  
  // 确认删除
  const confirmDelete = async () => {
    if (!currentArticleId.value) return
    
    deleting.value = true
    try {
      // 这里应该调用删除API
      // await articleApi.deleteArticle(currentArticleId.value)
      
      // 从草稿列表中移除
      draftArticles.value = draftArticles.value.filter(
        draft => draft.id !== currentArticleId.value
      )
      
      // 更新统计数据
      stats.value.draftCount--
      
      ElMessage.success('草稿已删除')
      showDeleteDialog.value = false
    } catch (error) {
      console.error('删除草稿失败:', error)
      ElMessage.error('删除失败，请重试')
    } finally {
      deleting.value = false
    }
  }
  
  // 获取统计数据
  const loadStats = async () => {
    try {
      // 这里应该调用统计API
      // 暂时使用模拟数据
      stats.value = {
        publishedCount: 15,
        draftCount: 3,
        totalViews: 12580,
        totalLikes: 1024
      }
    } catch (error) {
      console.error('获取统计数据失败:', error)
    }
  }
  
  // 加载已发布文章
  const loadPublishedArticles = async () => {
    try {
      // 这里应该调用API获取已发布文章
      // 暂时使用模拟数据
      const mockArticles: PublishedArticle[] = [
        {
          id: 1,
          title: '如何在家制作完美的提拉米苏',
          abstract: '提拉米苏是一道经典的意大利甜点，本文将分享专业大厨的独家秘方...',
          coverImage: 'https://picsum.photos/id/431/300/200',
          publishTime: '2024-12-06T10:30:00',
          readCount: 1250,
          likeCount: 89,
          commentCount: 23,
          status: 1
        },
        {
          id: 2,
          title: '极简主义生活方式探索',
          abstract: '在这个物质丰富的时代，越来越多的人开始追求极简主义生活方式...',
          coverImage: 'https://picsum.photos/id/432/300/200',
          publishTime: '2024-12-05T14:20:00',
          readCount: 890,
          likeCount: 67,
          commentCount: 15,
          status: 1
        },
        {
          id: 3,
          title: '2025年旅行计划制定指南',
          abstract: '新的一年即将到来，是时候开始规划你的旅行计划了...',
          coverImage: 'https://picsum.photos/id/433/300/200',
          publishTime: '2024-12-04T09:15:00',
          readCount: 2100,
          likeCount: 156,
          commentCount: 42,
          status: 1
        }
      ]
      
      allPublishedArticles.value = mockArticles
      publishedArticles.value = filteredPublishedArticles.value
    } catch (error) {
      console.error('获取已发布文章失败:', error)
      ElMessage.error('获取文章列表失败')
    }
  }
  
  // 加载草稿文章
  const loadDraftArticles = async () => {
    try {
      // 这里应该调用API获取草稿
      // 暂时使用模拟数据
      const mockDrafts: DraftArticle[] = [
        {
          id: 4,
          title: '春季护肤心得分享',
          abstract: '春天来了，皮肤护理也要跟着季节调整...',
          coverImage: 'https://picsum.photos/id/434/300/200',
          updateTime: '2024-12-07T16:45:00',
          wordCount: 1200,
          status: 0
        },
        {
          id: 5,
          title: '无标题草稿',
          abstract: '',
          coverImage: '',
          updateTime: '2024-12-06T20:30:00',
          wordCount: 350,
          status: 0
        },
        {
          id: 6,
          title: '健身新手入门指南',
          abstract: '对于刚开始健身的朋友来说，制定合适的计划很重要...',
          coverImage: 'https://picsum.photos/id/435/300/200',
          updateTime: '2024-12-05T11:20:00',
          wordCount: 800,
          status: 0
        }
      ]
      
      draftArticles.value = mockDrafts
    } catch (error) {
      console.error('获取草稿失败:', error)
      ElMessage.error('获取草稿列表失败')
    }
  }
  
  // 初始化
  onMounted(() => {
    loadStats()
    loadPublishedArticles()
  })
  
  return {
    stats,
    activeTab,
    searchKeyword,
    sortBy,
    publishedArticles,
    draftArticles,
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
    confirmDelete
  }
}
