import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { articleApi } from '@/api/article'

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

export default function useFollowView() {
  const router = useRouter()
  
  // 当前选中的标签
  const activeTab = ref('recommend')
  const limit = 10
  const offsetRecommend = ref(0)
  const offsetFollowing = ref(0)
  
  // 文章列表
  const articles = ref<Article[]>([])
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
  
  // 切换标签
  const changeTab = (tab: string) => {
    activeTab.value = tab
    fetchArticles(true)
  }
  
  // 加载更多文章
  const loadMoreArticles = async () => {
    try {
      const isFollowing = activeTab.value === 'following'
      const currentOffset = isFollowing ? offsetFollowing.value : offsetRecommend.value
      const resp = isFollowing
        ? await articleApi.getFollowingList({ offset: currentOffset, limit })
        : await articleApi.getRecommendList({ offset: currentOffset, limit })

      const more: Article[] = resp.map((article: any) => ({
        id: article.id,
        title: article.title,
        abstract: article.abstract,
        coverImage: article.coverImage,
        author: {
          id: article.author?.id,
          name: article.author?.name,
          avatar: article.author?.avatar
        },
        readCount: article.readCnt || 0,
        likeCount: article.likeCnt || 0,
        commentCount: 0,
        createTime: article.ctime
      }))

      articles.value = [...articles.value, ...more]
      if (isFollowing) {
        offsetFollowing.value += more.length
      } else {
        offsetRecommend.value += more.length
      }
      hasMoreArticles.value = more.length === limit
    } catch (error) {
      console.error('加载更多文章失败:', error)
      hasMoreArticles.value = false
    }
  }
  
  // 获取文章列表
  const fetchArticles = async (reset = false) => {
    try {
      const isFollowing = activeTab.value === 'following'
      if (reset) {
        if (isFollowing) {
          offsetFollowing.value = 0
        } else {
          offsetRecommend.value = 0
        }
        articles.value = []
        hasMoreArticles.value = true
      }

      const currentOffset = isFollowing ? offsetFollowing.value : offsetRecommend.value
      const resp = isFollowing
        ? await articleApi.getFollowingList({ offset: currentOffset, limit })
        : await articleApi.getRecommendList({ offset: currentOffset, limit })

      const list: Article[] = resp.map((article: any) => ({
        id: article.id,
        title: article.title,
        abstract: article.abstract,
        coverImage: article.coverImage,
        author: {
          id: article.author?.id,
          name: article.author?.name,
          avatar: article.author?.avatar
        },
        readCount: article.readCnt || 0,
        likeCount: article.likeCnt || 0,
        commentCount: 0,
        createTime: article.ctime
      }))

      articles.value = list
      if (isFollowing) {
        offsetFollowing.value += list.length
      } else {
        offsetRecommend.value += list.length
      }
      hasMoreArticles.value = list.length === limit
    } catch (error) {
      console.error('获取文章失败:', error)
      
      // 使用模拟数据
      articles.value = [
        {
          id: 1,
          title: '如何在家制作完美的提拉米苏',
          abstract: '提拉米苏是一道经典的意大利甜点，本文将分享专业大厨的独家秘方...',
          coverImage: 'https://picsum.photos/id/431/400/300',
          author: {
            id: 101,
            name: '美食达人',
            avatar: 'https://picsum.photos/id/1027/100/100'
          },
          readCount: 12500,
          likeCount: 3200,
          commentCount: 128,
          createTime: '2024-11-28T14:20:00'
        },
        {
          id: 2,
          title: '2025年最值得去的10个小众旅行地',
          abstract: '厌倦了人山人海的热门景点？这些鲜为人知的目的地将带给你全新的旅行体验...',
          coverImage: 'https://picsum.photos/id/1036/400/300',
          author: {
            id: 102,
            name: '旅行笔记',
            avatar: 'https://picsum.photos/id/1012/100/100'
          },
          readCount: 18700,
          likeCount: 5400,
          commentCount: 342,
          createTime: '2024-11-25T09:15:00'
        },
        {
          id: 3,
          title: '极简主义：如何通过断舍离改变你的生活',
          abstract: '极简主义不仅是一种生活方式，更是一种思维模式。本文将分享如何开始你的极简之旅...',
          coverImage: 'https://picsum.photos/id/106/400/300',
          author: {
            id: 103,
            name: '生活家',
            avatar: 'https://picsum.photos/id/1005/100/100'
          },
          readCount: 9800,
          likeCount: 2100,
          commentCount: 98,
          createTime: '2024-11-20T16:40:00'
        },
        {
          id: 4,
          title: '数字游民：如何边旅行边工作',
          abstract: '远程工作正在改变我们的生活和工作方式，本文分享如何成为一名成功的数字游民...',
          coverImage: 'https://picsum.photos/id/1081/400/300',
          author: {
            id: 104,
            name: '自由职业者',
            avatar: 'https://picsum.photos/id/1025/100/100'
          },
          readCount: 15300,
          likeCount: 4200,
          commentCount: 215,
          createTime: '2024-11-15T08:45:00'
        }
      ]
      
      hasMoreArticles.value = true
    }
  }
  
  onMounted(() => {
    fetchArticles(true)
  })
  
  return {
    activeTab,
    articles,
    hasMoreArticles,
    formatNumber,
    viewArticle,
    changeTab,
    loadMoreArticles
  }
}
