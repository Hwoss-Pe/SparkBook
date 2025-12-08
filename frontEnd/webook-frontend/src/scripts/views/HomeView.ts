import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { articleApi } from '@/api/article'
import type { ArticlePub } from '@/api/article'

// 首页文章展示的数据结构
interface HomeArticle {
  id: number
  title: string
  abstract: string
  coverImage: string
  author: {
    id: number
    name: string
  }
  readCount: number
  likeCount: number
  collectCount: number
}

export default function useHomeView() {
  const router = useRouter()

  // Banner数据
  const bannerItems = ref([
    {
      title: '探索美食新世界',
      description: '跟随我们的美食专家，发现城市里的隐藏美食',
      image: 'https://picsum.photos/id/1080/800/450'
    },
    {
      title: '旅行的意义',
      description: '走过山川湖海，感受不一样的人生',
      image: 'https://picsum.photos/id/1036/800/450'
    },
    {
      title: '生活方式指南',
      description: '如何打造舒适且有品质的生活空间',
      image: 'https://picsum.photos/id/106/800/450'
    }
  ])

  // 文章列表
  const articles = ref<HomeArticle[]>([])
  // 当前分页
  const currentOffset = ref(0)
  const pageSize = 10

  // 热榜数据
  const hotRankings = ref<Array<{
    id: number
    title: string
    author: { name: string }
    readCount: number
  }>>([])

  // 获取热榜数据
  const fetchHotRankings = async () => {
    try {
      const { rankingApi } = await import('@/api/ranking')
      const response = await rankingApi.getRanking({ offset: 0, limit: 8 })
      
      // if (response.code === 0 && response.data) {
        hotRankings.value = response.map((article: { id: any; title: any; author: { name: any }; readCnt: any }) => ({
          id: article.id,
          title: article.title,
          author: { name: article.author.name || '匿名用户' },
          readCount: article.readCnt || 0
        }))
      // }
    } catch (error) {
      console.error('获取热榜数据失败:', error)
      // 如果获取失败，保持空数组
      hotRankings.value = []
    }
  }

  // 模拟推荐作者数据
  const recommendedAuthors = ref([
    {
      id: 201,
      name: '美食探店家',
      avatar: 'https://picsum.photos/id/1062/100/100',
      description: '探索城市里的美食秘境',
      isFollowed: false
    },
    {
      id: 202,
      name: '旅行摄影师',
      avatar: 'https://picsum.photos/id/1074/100/100',
      description: '用镜头记录世界的美',
      isFollowed: true
    },
    {
      id: 203,
      name: '生活方式指导',
      avatar: 'https://picsum.photos/id/1025/100/100',
      description: '让生活更有品质的小贴士',
      isFollowed: false
    }
  ])

  const hasMoreArticles = ref(true)

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

  // 将API返回的数据转换为首页展示的数据结构
  const convertToHomeArticle = (article: ArticlePub): HomeArticle => {
    return {
      id: article.id,
      title: article.title,
      abstract: article.abstract,
      coverImage: article.coverImage || `https://picsum.photos/id/${400 + article.id}/400/300`,
      author: {
        id: article.author?.id || 0,
        name: article.author?.name || '匿名用户'
      },
      readCount: article.readCnt || 0,
      likeCount: article.likeCnt || 0,
      collectCount: article.collectCnt || 0
    }
  }

  // 获取推荐文章列表
  const fetchArticles = async (isLoadMore = false) => {
    try {
      const offset = isLoadMore ? currentOffset.value : 0
      const res = await articleApi.getRecommendList({
        offset,
        limit: pageSize
      })
      
      const newArticles = (res || []).map(convertToHomeArticle)
      
      if (isLoadMore) {
        articles.value = [...articles.value, ...newArticles]
      } else {
        articles.value = newArticles
      }
      
      currentOffset.value = offset + newArticles.length
      hasMoreArticles.value = newArticles.length >= pageSize
    } catch (error) {
      console.error('获取推荐文章失败:', error)
      hasMoreArticles.value = false
    }
  }

  // 加载更多文章
  const loadMoreArticles = () => {
    fetchArticles(true)
  }

  // 关注作者
  const followAuthor = (authorId: number) => {
    // 这里应该调用关注API
    // 模拟关注/取消关注
    const author = recommendedAuthors.value.find(a => a.id === authorId)
    if (author) {
      author.isFollowed = !author.isFollowed
    }
  }

  onMounted(() => {
    // 获取推荐文章列表
    fetchArticles()
    // 获取热榜数据
    fetchHotRankings()
  })

  return {
    bannerItems,
    articles,
    hotRankings,
    recommendedAuthors,
    hasMoreArticles,
    formatNumber,
    viewArticle,
    loadMoreArticles,
    followAuthor
  }
}
