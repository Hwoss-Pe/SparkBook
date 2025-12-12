import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { rankingApi, type RankingArticle } from '@/api/ranking'
import { articleApi, type ArticlePub } from '@/api/article'
import { ElMessage } from 'element-plus'
import { resolveStaticUrl } from '@/api/http'

// 定义类型接口
interface Article {
  id: number;
  title: string;
  abstract: string;
  coverImage: string;
  author: {
    id: number;
    name: string;
    avatar: string;
  };
  readCount: number;
  likeCount: number;
  commentCount: number;
  createTime: string;
}

export default function useHotView() {
  const router = useRouter()
  
  // 当前选中的标签
  const activeTab = ref('overall')
  const officialTags = ref<string[]>([])
  
  // 热门文章列表
  const hotArticles = ref<Article[]>([])
  const hasMoreArticles = ref(true)
  const isTriggering = ref(false)
  
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
  
  // 切换标签
  const changeTab = (tab: string) => {
    activeTab.value = tab
    fetchHotArticles()
  }
  
  // 加载更多文章
  const loadMoreArticles = async () => {
    try {
      const currentLength = hotArticles.value.length
      if (activeTab.value === 'overall') {
        const data = await rankingApi.getRanking({ offset: currentLength, limit: 10 })
        if (Array.isArray(data) && data.length > 0) {
          const moreArticles: Article[] = data.map((article: RankingArticle) => ({
            id: article.id,
            title: article.title,
            abstract: article.abstract || '',
            coverImage: article.coverImage ? resolveStaticUrl(article.coverImage) : `https://picsum.photos/id/${400 + article.id}/400/300`,
            author: {
              id: article.author.id,
              name: article.author.name || '匿名用户',
              avatar: article.author.avatar ? resolveStaticUrl(article.author.avatar) : `https://picsum.photos/id/${1000 + article.author.id}/100/100`
            },
            readCount: article.readCnt || 0,
            likeCount: article.likeCnt || 0,
            commentCount: article.collectCnt || 0,
            createTime: article.ctime
          }))
          hotArticles.value = [...hotArticles.value, ...moreArticles]
          hasMoreArticles.value = data.length >= 10
        } else {
          hasMoreArticles.value = false
        }
      } else {
        const data = await articleApi.getArticlesByOfficialTag({ tag: activeTab.value, offset: currentLength, limit: 10 })
        if (Array.isArray(data) && data.length > 0) {
          const moreArticles: Article[] = data.map((article: ArticlePub) => ({
            id: article.id,
            title: article.title,
            abstract: article.abstract || '',
            coverImage: article.coverImage ? resolveStaticUrl(article.coverImage) : `https://picsum.photos/id/${400 + article.id}/400/300`,
            author: {
              id: article.author.id,
              name: article.author.name || '匿名用户',
              avatar: article.author.avatar ? resolveStaticUrl(article.author.avatar) : `https://picsum.photos/id/${1000 + article.author.id}/100/100`
            },
            readCount: article.readCnt || 0,
            likeCount: article.likeCnt || 0,
            commentCount: article.collectCnt || 0,
            createTime: article.ctime
          }))
          hotArticles.value = [...hotArticles.value, ...moreArticles]
          hasMoreArticles.value = data.length >= 10
        } else {
          hasMoreArticles.value = false
        }
      }
    } catch (error) {
      console.error('加载更多文章失败:', error)
      ElMessage.error('加载更多文章失败')
      hasMoreArticles.value = false
    }
  }
  
  // 获取热门文章
  const fetchHotArticles = async () => {
    try {
      if (activeTab.value === 'overall') {
        const data = await rankingApi.getRanking({ offset: 0, limit: 20 })
        const articles: Article[] = (Array.isArray(data) ? data : []).map((article: RankingArticle) => ({
          id: article.id,
          title: article.title,
          abstract: article.abstract || '',
          coverImage: article.coverImage ? resolveStaticUrl(article.coverImage) : `https://picsum.photos/id/${400 + article.id}/400/300`,
          author: {
            id: article.author.id,
            name: article.author.name || '匿名用户',
            avatar: article.author.avatar ? resolveStaticUrl(article.author.avatar) : `https://picsum.photos/id/${1000 + article.author.id}/100/100`
          },
          readCount: article.readCnt || 0,
          likeCount: article.likeCnt || 0,
          commentCount: article.collectCnt || 0,
          createTime: article.ctime
        }))
        hotArticles.value = articles
        hasMoreArticles.value = articles.length >= 20
      } else {
        const data = await articleApi.getArticlesByOfficialTag({ tag: activeTab.value, offset: 0, limit: 20 })
        const articles: Article[] = (Array.isArray(data) ? data : []).map((article: ArticlePub) => ({
          id: article.id,
          title: article.title,
          abstract: article.abstract || '',
          coverImage: article.coverImage ? resolveStaticUrl(article.coverImage) : `https://picsum.photos/id/${400 + article.id}/400/300`,
          author: {
            id: article.author.id,
            name: article.author.name || '匿名用户',
            avatar: article.author.avatar ? resolveStaticUrl(article.author.avatar) : `https://picsum.photos/id/${1000 + article.author.id}/100/100`
          },
          readCount: article.readCnt || 0,
          likeCount: article.likeCnt || 0,
          commentCount: article.collectCnt || 0,
          createTime: article.ctime || ''
        }))
        hotArticles.value = articles
        hasMoreArticles.value = articles.length >= 20
      }
    } catch (error) {
      console.error('获取热门文章失败:', error)
      ElMessage.error('获取热榜数据失败，请稍后重试')
      
      // 如果API调用失败，显示空列表
      hotArticles.value = []
      hasMoreArticles.value = false
    }
  }

  // 手动触发热榜计算
  const triggerRankingCalculation = async () => {
    try {
      isTriggering.value = true
      await rankingApi.triggerRanking()
      ElMessage.success('热榜计算已触发，正在重新计算中...')
      setTimeout(() => {
        fetchHotArticles()
      }, 2000)
    } catch (error) {
      console.error('触发热榜计算失败:', error)
      ElMessage.error('触发热榜计算失败，请稍后重试')
    } finally {
      isTriggering.value = false
    }
  }
  
  onMounted(async () => {
    try {
      const tags = await articleApi.getOfficialTags()
      officialTags.value = Array.isArray(tags) ? tags : []
    } catch {
      officialTags.value = []
    }
    fetchHotArticles()
  })
  
  return {
    activeTab,
    officialTags,
    hotArticles,
    hasMoreArticles,
    isTriggering,
    formatNumber,
    viewArticle,
    changeTab,
    loadMoreArticles,
    triggerRankingCalculation
  }
}
