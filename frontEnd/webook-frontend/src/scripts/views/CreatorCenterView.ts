import { ref, reactive, onMounted, onUnmounted, nextTick, computed } from 'vue'
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
  collectCount: number
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
  const publishedOffset = ref(0)
  const publishedLimit = 10
  const hasMorePublished = ref(true)
  const loadingPublished = ref(false)
  const publishedLoadMoreRef = ref<HTMLElement | null>(null)
  let publishedObserver: IntersectionObserver | null = null
  
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
      loadPublishedArticles(true)
      setupPublishedObserver()
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
      publishedArticles.value = publishedArticles.value.map(a =>
        a.id === currentArticleId.value ? { ...a, status: 3 } : a
      )
      allPublishedArticles.value = allPublishedArticles.value.map(a =>
        a.id === currentArticleId.value ? { ...a, status: 3 } : a
      )
      ElMessage.success('已设为仅自己可见')
      showWithdrawDialog.value = false
    } catch (error) {
      console.error('撤回文章失败:', error)
      ElMessage.error('撤回失败，请重试')
    } finally {
      withdrawing.value = false
    }
  }

  const unpublishArticle = async (id: number) => {
    try {
      await ElMessageBox.confirm('撤回发布后文章将变为草稿，确认继续？', '撤回发布', {
        type: 'warning'
      })
      await articleApi.unpublishArticle(id, userStore.user?.id || 0)
      publishedArticles.value = publishedArticles.value.filter(a => a.id !== id)
      allPublishedArticles.value = allPublishedArticles.value.filter(a => a.id !== id)
      stats.value.publishedCount--
      stats.value.draftCount++
      ElMessage.success('已撤回为草稿')
      if (activeTab.value === 'drafts') {
        loadDraftArticles()
      }
    } catch (error: any) {
      if (error !== 'cancel') {
        console.error('撤回发布失败:', error)
        ElMessage.error('撤回失败，请重试')
      }
    }
  }
  
  // 发布草稿
  const publishDraft = async (id: number) => {
    try {
      const detail = await articleApi.getArticleById(id)
      await articleApi.publishArticle({ id, title: detail.title, content: detail.content })
      draftArticles.value = draftArticles.value.filter(draft => draft.id !== id)
      stats.value.draftCount--
      stats.value.publishedCount++
      ElMessage.success('草稿已发布')
      if (activeTab.value === 'published') {
        loadPublishedArticles(true)
      }
    } catch (error) {
      console.error('发布草稿失败:', error)
      ElMessage.error('发布失败，请重试')
    }
  }
  
  // 复制草稿
  const duplicateDraft = async (id: number) => {
    try {
      const detail = await articleApi.getArticleById(id)
      const response = await articleApi.saveArticle({
        title: `${detail.title} - 副本`,
        abstract: detail.abstract,
        content: detail.content,
        coverImage: detail.coverImage,
        status: 1
      })
      ElMessage.success('草稿已复制')
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
      userStore.initUserState()
      const uid = userStore.user?.id || 0
      await articleApi.deleteDraft(currentArticleId.value, uid)
      
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
  
  const loadStats = async () => {
    try {
      userStore.initUserState()
      const uid = userStore.user?.id || 0
      const s = await articleApi.getAuthorStats(uid)
      stats.value = {
        publishedCount: s.publishedCount || 0,
        draftCount: s.draftCount || 0,
        totalViews: s.totalReadCount || 0,
        totalLikes: s.totalLikeCount || 0
      }
    } catch (error) {
      console.error('获取统计数据失败:', error)
    }
  }
  
  // 加载已发布文章
  const loadPublishedArticles = async (reset = false) => {
    try {
      userStore.initUserState()
      const uid = userStore.user?.id || 0
      const currentOffset = reset ? 0 : publishedOffset.value
      const resp = await articleApi.getList({ offset: currentOffset, limit: publishedLimit })
      const pubs = resp.filter(a => a.status === 2 || a.status === 3)
      const list: PublishedArticle[] = pubs.map(a => ({
        id: a.id,
        title: a.title,
        abstract: a.abstract,
        coverImage: a.coverImage || '',
        publishTime: a.ctime,
        readCount: a.readCnt || 0,
        likeCount: a.likeCnt || 0,
        collectCount: a.collectCnt || 0,
        status: a.status
      }))
      allPublishedArticles.value = reset ? list : [...allPublishedArticles.value, ...list]
      publishedArticles.value = filteredPublishedArticles.value
      publishedOffset.value = currentOffset + resp.length
      hasMorePublished.value = resp.length >= publishedLimit
    } catch (error) {
      console.error('获取已发布文章失败:', error)
      ElMessage.error('获取文章列表失败')
    }
  }

  const loadMorePublishedArticles = async () => {
    if (loadingPublished.value || !hasMorePublished.value || activeTab.value !== 'published') return
    loadingPublished.value = true
    await loadPublishedArticles(false)
    loadingPublished.value = false
  }

  const setupPublishedObserver = () => {
    if (publishedObserver) publishedObserver.disconnect()
    publishedObserver = new IntersectionObserver((entries) => {
      if (entries.some(e => e.isIntersecting)) {
        loadMorePublishedArticles()
      }
    })
    nextTick(() => {
      if (publishedLoadMoreRef.value) {
        publishedObserver?.observe(publishedLoadMoreRef.value)
      }
    })
  }
  
  // 加载草稿文章
  const loadDraftArticles = async () => {
    try {
      userStore.initUserState()
      const uid = userStore.user?.id || 0
      const resp = await articleApi.getList({ offset: 0, limit: 100 })
      const rawDrafts = resp.filter(a => a.status === 1)
      const details = await Promise.all(rawDrafts.map(a => articleApi.getArticleById(a.id).catch(() => null)))
      const valid = details.filter(d => d)
      draftArticles.value = valid.map(d => ({
        id: d!.id,
        title: d!.title,
        abstract: '',
        coverImage: d!.coverImage || '',
        updateTime: d!.utime,
        wordCount: (d!.content || '').length,
        status: 1
      }))
    } catch (error) {
      console.error('获取草稿失败:', error)
      ElMessage.error('获取草稿列表失败')
    }
  }
  
  // 初始化
  onMounted(() => {
    loadStats()
    loadPublishedArticles(true)
    setupPublishedObserver()
  })

  onUnmounted(() => {
    publishedObserver?.disconnect()
  })
  
  return {
    stats,
    activeTab,
    searchKeyword,
    sortBy,
    publishedArticles,
    draftArticles,
    hasMorePublished,
    loadingPublished,
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
  }
}
