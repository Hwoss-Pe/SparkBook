import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { searchApi, type SearchResponse } from '@/api/search'
import { ElMessage } from 'element-plus'
import { resolveStaticUrl } from '@/api/http'

// 定义类型接口（与UI展示所需的最小字段）
interface Article {
  id: number
  title: string
  abstract: string
}

interface User {
  id: number
  name: string
  avatar: string
  aboutMe?: string
}

export default function useSearchView() {
  const router = useRouter()
  const route = useRoute()
  
  // 搜索关键词
  const searchQuery = ref('')
  
  // 当前选中的标签
  const activeTab = ref('article')
  
  // 搜索结果
  const articles = ref<Article[]>([])
  const users = ref<User[]>([])
  const hasMoreResults = ref(false)
  const isSearching = ref(false)
  const totalResults = ref(0)

  const normalizeAvatarUrl = (u: string | undefined): string => {
    if (!u) return ''
    const v = u.replace(/`/g, '').trim()
    return resolveStaticUrl(v)
  }
  
  // 格式化数字，例如1200显示为1.2k
  const formatNumber = (num: number): string => {
    if (num >= 10000) {
      return (num / 10000).toFixed(1) + 'w'
    } else if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'k'
    }
    return num.toString()
  }
  
  // 查看文章详情
  const viewArticle = (id: number) => {
    router.push(`/article/${id}`)
  }
  
  // 查看用户主页
  const viewUser = (id: number) => {
    router.push(`/user/${id}`)
  }
  
  // 切换标签
  const changeTab = (tab: string) => {
    activeTab.value = tab
  }
  
  // 加载更多结果（暂不实现分页，使用后端真实数据）
  const loadMoreResults = async () => {
    hasMoreResults.value = false
  }
  
  // 执行搜索
  const performSearch = async () => {
    if (!searchQuery.value.trim()) {
      ElMessage.warning('请输入搜索关键词')
      return
    }
    
    try {
      isSearching.value = true
      
      // 调用搜索API
      const currentUserId = 101 // 当前用户ID，实际应该从用户状态中获取
      const response = await searchApi.search({
        expression: searchQuery.value,
        uid: currentUserId
      })
      
      // 处理文章搜索结果，使用真实字段并生成摘要
      const searchedArticles = (response as SearchResponse).article?.articles ?? []
      const makeAbstract = (content: string): string => {
        const text = (content || '').replace(/[#`*_>\-]+/g, ' ').replace(/\s+/g, ' ').trim()
        return text.length > 120 ? text.slice(0, 120) + '…' : text
      }
      articles.value = searchedArticles.map(a => ({
        id: a.id,
        title: a.title,
        abstract: makeAbstract(a.content)
      }))
      
      // 处理用户搜索结果，使用真实字段
      users.value = ((response as SearchResponse).user?.users ?? []).map(u => ({
        id: u.id,
        name: u.nickname,
        avatar: normalizeAvatarUrl(u.avatar) || '',
        aboutMe: ''
      }))
      
      totalResults.value = articles.value.length + users.value.length
      hasMoreResults.value = false
      
      // 更新URL，但不触发路由变化
      const query = { ...route.query, q: searchQuery.value }
      router.replace({ query })
    } catch (error) {
      console.error('搜索失败:', error)
      articles.value = []
      users.value = []
      totalResults.value = 0
      hasMoreResults.value = false
    } finally {
      isSearching.value = false
    }
  }
  
  // 监听路由变化
  watch(() => route.query.q, (newQuery) => {
    if (newQuery) {
      searchQuery.value = newQuery as string
      performSearch()
    }
  }, { immediate: true })
  
  onMounted(() => {
    // 如果URL中有搜索参数，执行搜索
    if (route.query.q) {
      searchQuery.value = route.query.q as string
      performSearch()
    }
  })
  
  return {
    searchQuery,
    activeTab,
    articles,
    users,
    hasMoreResults,
    isSearching,
    totalResults,
    formatNumber,
    viewArticle,
    viewUser,
    changeTab,
    loadMoreResults,
    performSearch
  }
}
