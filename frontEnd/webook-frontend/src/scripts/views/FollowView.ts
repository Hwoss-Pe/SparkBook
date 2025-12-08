import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { followApi } from '@/api/follow'
import { articleApi } from '@/api/article'
import { interactiveApi } from '@/api/interactive'

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
    fetchArticles()
  }
  
  // 加载更多文章
  const loadMoreArticles = async () => {
    try {
      // 这里应该调用API加载更多文章
      // 实际项目中应该根据不同的标签（推荐、关注）调用不同的API
      
      // 模拟加载更多
      const moreArticles: Article[] = [
        {
          id: 7,
          title: '零基础学习编程：从何开始？',
          abstract: '想学编程但不知道从何入手？本文为你提供清晰的学习路径...',
          coverImage: 'https://picsum.photos/id/0/400/300',
          author: {
            id: 107,
            name: '编程教练',
            avatar: 'https://picsum.photos/id/1/100/100'
          },
          readCount: 18500,
          likeCount: 5100,
          commentCount: 320,
          createTime: '2024-11-10T10:30:00'
        },
        {
          id: 8,
          title: '每天10分钟，21天养成冥想习惯',
          abstract: '冥想不仅能减轻压力，还能提高注意力和创造力...',
          coverImage: 'https://picsum.photos/id/1029/400/300',
          author: {
            id: 108,
            name: '心灵导师',
            avatar: 'https://picsum.photos/id/1002/100/100'
          },
          readCount: 12300,
          likeCount: 3600,
          commentCount: 145,
          createTime: '2024-11-05T15:45:00'
        }
      ]
      
      articles.value = [...articles.value, ...moreArticles]
      
      // 假设没有更多文章了
      hasMoreArticles.value = false
    } catch (error) {
      console.error('加载更多文章失败:', error)
      hasMoreArticles.value = false
    }
  }
  
  // 获取文章列表
  const fetchArticles = async () => {
    try {
      const currentUserId = 101 // 当前用户ID，实际应该从用户状态中获取
      
      // 获取关注列表
      const followResponse = await followApi.getFollowee({
        follower: currentUserId,
        offset: 0,
        limit: 100
      })
      
      // 提取关注的用户ID
      const followeeIds = followResponse.follow_relations.map(relation => relation.followee)
      
      if (followeeIds.length === 0) {
        articles.value = []
        hasMoreArticles.value = false
        return
      }
      
      // 获取关注用户的文章
      // 实际项目中应该有专门的API来获取关注用户的文章
      // 这里模拟一下，假设调用了文章列表API并过滤出关注用户的文章
      const articlesResponse = await articleApi.getPublishedList({
        offset: 0,
        limit: 10
      })
      
      // 过滤出关注用户的文章
      const followedArticles = articlesResponse.articles.filter(article => 
        followeeIds.includes(article.author.id)
      )
      
      // 获取交互数据
      const articleIds = followedArticles.map(article => article.id)
      const interactiveData = await interactiveApi.getInteractiveByIds({
        biz: 'article',
        ids: articleIds
      })
      
      // 构建文章列表
      const articleList: Article[] = followedArticles.map(article => {
        const interactive = interactiveData.intrs[article.id] || {
          read_cnt: 0,
          like_cnt: 0,
          collect_cnt: 0
        }
        
        return {
          id: article.id,
          title: article.title,
          abstract: article.abstract,
          coverImage: `https://picsum.photos/id/${400 + article.id}/400/300`, // 模拟封面图
          author: {
            id: article.author.id,
            name: article.author.name,
            avatar: `https://picsum.photos/id/${1000 + article.author.id}/100/100` // 模拟头像
          },
          readCount: interactive.read_cnt,
          likeCount: interactive.like_cnt,
          commentCount: 0, // 评论数需要从评论API获取
          createTime: article.ctime
        }
      })
      
      articles.value = articleList
      hasMoreArticles.value = articleList.length >= 10 // 如果返回10条或更多，可能还有更多
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
    fetchArticles()
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
