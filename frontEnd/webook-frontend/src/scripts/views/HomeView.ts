import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { articleApi } from '@/api/article'
import { resolveStaticUrl } from '@/api/http'
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
    avatar: string
  }
  readCount: number
  likeCount: number
  collectCount: number
  isLiked?: boolean
  isFavorited?: boolean
  isLikeAnimating?: boolean
  isFavAnimating?: boolean
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
        hotRankings.value = response.map((article: { id: any; title: any; author: { name: any }; readCnt: any;}) => ({
          id: article.id,
          title: article.title,
          author: { name: article.author.name || '匿名用户' },
          avatar:{ name: article.author.name || '匿名用户' },
          readCount: article.readCnt || 0
        }))
      // }
    } catch (error) {
      console.error('获取热榜数据失败:', error)
      // 如果获取失败，保持空数组
      hotRankings.value = []
    }
  }

  // 推荐作者数据（真实接口）
  const recommendedAuthors = ref<Array<{
    id: number
    name: string
    avatar?: string
    description?: string
    isFollowed: boolean
  }>>([])

  // 获取推荐作者
  const fetchRecommendedAuthors = async () => {
    try {
      const { userApi } = await import('@/api/user')
      const authors = await userApi.getRecommendAuthors(6)
      recommendedAuthors.value = (authors || []).map(a => ({
        id: a.id,
        name: a.name || '匿名作者',
        avatar: a.avatar ? resolveStaticUrl(a.avatar) : 'https://picsum.photos/seed/avatar/100/100',
        description: (() => {
          const desc = a.description || ''
          return desc.length > 10 ? desc.slice(0, 10) + '...' : desc
        })(),
        isFollowed: false
      }))
    } catch (error) {
      console.error('获取推荐作者失败:', error)
      recommendedAuthors.value = []
    }
  }

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

  // 查看用户主页
  const viewUser = (id: number) => {
    router.push(`/user/${id}`)
  }

  // 将API返回的数据转换为首页展示的数据结构
  const convertToHomeArticle = (article: ArticlePub): HomeArticle => {
    return {
      id: article.id,
      title: article.title,
      abstract: article.abstract,
      coverImage: article.coverImage ? resolveStaticUrl(article.coverImage) : `https://picsum.photos/id/${400 + article.id}/400/300`,
      author: {
        id: article.author?.id || 0,
        name: article.author?.name || '匿名用户',
        avatar: article.author?.avatar ? resolveStaticUrl(article.author.avatar) : 'https://picsum.photos/seed/avatar/100/100'
      },
      readCount: article.readCnt || 0,
      likeCount: article.likeCnt || 0,
      collectCount: article.collectCnt || 0,
      isLiked: !!article.liked,
      isFavorited: !!article.collected
    }
  }

  let isMounted = true
  onUnmounted(() => {
    isMounted = false
  })

  // 获取推荐文章列表
  const fetchArticles = async (isLoadMore = false) => {
    if (!isMounted) return
    try {
      const offset = isLoadMore ? currentOffset.value : 0
      const res = await articleApi.getRecommendList({
        offset,
        limit: pageSize
      })
      
      if (!isMounted) return

      const newArticles = (res || []).map(convertToHomeArticle)
      
      if (isLoadMore) {
        // 过滤重复文章
        const existingIds = new Set(articles.value.map(a => a.id))
        const uniqueNewArticles = newArticles.filter(a => !existingIds.has(a.id))
        articles.value = [...articles.value, ...uniqueNewArticles]
      } else {
        // 确保新加载的列表也没有重复（后端可能返回重复数据）
        const uniqueNewArticles: HomeArticle[] = []
        const seenIds = new Set<number>()
        for (const article of newArticles) {
          if (!seenIds.has(article.id)) {
            uniqueNewArticles.push(article)
            seenIds.add(article.id)
          }
        }
        articles.value = uniqueNewArticles
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

  const toggleArticleLike = async (article: HomeArticle) => {
    const prevLiked = !!article.isLiked
    const prevCnt = article.likeCount || 0
    // 乐观更新：先改UI
    article.isLiked = !prevLiked
    article.likeCount = Math.max(0, prevCnt + (article.isLiked ? 1 : -1))
    article.isLikeAnimating = true
    setTimeout(() => { article.isLikeAnimating = false }, 300)

    try {
      if (article.isLiked) {
        await articleApi.like(article.id)
      } else {
        await articleApi.cancelLike(article.id)
      }
    } catch (error) {
      // 回滚
      article.isLiked = prevLiked
      article.likeCount = prevCnt
      console.error('首页点赞操作失败:', error)
    }
  }

  const toggleArticleFavorite = async (article: HomeArticle) => {
    try {
      const defaultCid = 1
      if (article.isFavorited) {
        await articleApi.cancelCollect(article.id, defaultCid)
        article.isFavorited = false
        article.collectCount = Math.max(0, (article.collectCount || 0) - 1)
      } else {
        await articleApi.collect(article.id, defaultCid)
        article.isFavorited = true
        article.collectCount = (article.collectCount || 0) + 1
      }
      article.isFavAnimating = true
      setTimeout(() => { article.isFavAnimating = false }, 300)
    } catch (error) {
      console.error('首页收藏操作失败:', error)
    }
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
    // 获取推荐作者
    fetchRecommendedAuthors()
  })

  onMounted(() => {
    // 获取推荐文章列表
    fetchArticles()
    // 获取热榜数据
    fetchHotRankings()
    // 获取推荐作者
    fetchRecommendedAuthors()
  })

  return {
    bannerItems,
    articles,
    hotRankings,
    recommendedAuthors,
    hasMoreArticles,
    formatNumber,
    viewArticle,
    viewUser,
    loadMoreArticles,
    followAuthor,
    toggleArticleLike,
    toggleArticleFavorite
  }
}
