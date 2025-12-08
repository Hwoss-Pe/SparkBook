import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { searchApi } from '@/api/search'
import { interactiveApi } from '@/api/interactive'
import { ElMessage } from 'element-plus'

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

interface User {
  id: number;
  name: string;
  avatar: string;
  aboutMe: string;
  followersCount: number;
  articlesCount: number;
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
  
  // 加载更多结果
  const loadMoreResults = async () => {
    try {
      // 这里应该调用API加载更多搜索结果
      // 实际项目中应该根据不同的标签（文章、用户）调用不同的API
      
      // 模拟加载更多
      if (activeTab.value === 'article') {
        const moreArticles: Article[] = [
          {
            id: 7,
            title: `搜索结果：${searchQuery.value} - 编程入门指南`,
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
            title: `搜索结果：${searchQuery.value} - 冥想习惯养成`,
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
      } else if (activeTab.value === 'user') {
        const moreUsers: User[] = [
          {
            id: 5,
            name: `${searchQuery.value}爱好者`,
            avatar: 'https://picsum.photos/id/1074/100/100',
            aboutMe: '热爱分享生活和学习经验',
            followersCount: 2800,
            articlesCount: 45
          },
          {
            id: 6,
            name: `${searchQuery.value}专家`,
            avatar: 'https://picsum.photos/id/1062/100/100',
            aboutMe: '专注于技术和创新领域',
            followersCount: 5600,
            articlesCount: 78
          }
        ]
        
        users.value = [...users.value, ...moreUsers]
      }
      
      // 假设没有更多结果了
      hasMoreResults.value = false
    } catch (error) {
      console.error('加载更多结果失败:', error)
      hasMoreResults.value = false
    }
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
      
      // 处理文章搜索结果
      const searchedArticles = response.article.articles
      
      // 获取文章ID列表
      const articleIds = searchedArticles.map(article => article.id)
      
      // 获取交互数据
      let interactiveData = { intrs: {} }
      if (articleIds.length > 0) {
        interactiveData = await interactiveApi.getInteractiveByIds({
          biz: 'article',
          ids: articleIds
        })
      }
      
      // 构建文章列表
      articles.value = searchedArticles.map(article => {
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
          createTime: new Date().toISOString() // 模拟创建时间
        }
      })
      
      // 处理用户搜索结果
      users.value = response.user.users.map(user => {
        return {
          id: user.id,
          name: user.nickname,
          avatar: `https://picsum.photos/id/${1000 + user.id}/100/100`, // 模拟头像
          aboutMe: user.aboutMe || '这个人很懒，什么都没留下',
          followersCount: Math.floor(Math.random() * 10000), // 模拟粉丝数
          articlesCount: Math.floor(Math.random() * 100) // 模拟文章数
        }
      })
      
      totalResults.value = articles.value.length + users.value.length
      hasMoreResults.value = totalResults.value >= 10 // 如果返回10条或更多，可能还有更多
      
      // 更新URL，但不触发路由变化
      const query = { ...route.query, q: searchQuery.value }
      router.replace({ query })
    } catch (error) {
      console.error('搜索失败:', error)
      
      // 使用模拟数据
      if (activeTab.value === 'article') {
        articles.value = [
          {
            id: 1,
            title: `搜索结果：${searchQuery.value} - 美食制作`,
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
            title: `搜索结果：${searchQuery.value} - 旅行指南`,
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
            title: `搜索结果：${searchQuery.value} - 极简生活`,
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
            title: `搜索结果：${searchQuery.value} - 数字游民`,
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
      }
      
      users.value = [
        {
          id: 1,
          name: `${searchQuery.value}粉丝`,
          avatar: 'https://picsum.photos/id/1027/100/100',
          aboutMe: '热爱美食和旅行的自由灵魂',
          followersCount: 5600,
          articlesCount: 32
        },
        {
          id: 2,
          name: `${searchQuery.value}达人`,
          avatar: 'https://picsum.photos/id/1012/100/100',
          aboutMe: '专注于旅行体验分享',
          followersCount: 8900,
          articlesCount: 67
        },
        {
          id: 3,
          name: `${searchQuery.value}爱好者`,
          avatar: 'https://picsum.photos/id/1005/100/100',
          aboutMe: '生活方式探索者',
          followersCount: 3200,
          articlesCount: 28
        },
        {
          id: 4,
          name: `${searchQuery.value}专家`,
          avatar: 'https://picsum.photos/id/1025/100/100',
          aboutMe: '自由职业者，数字游民',
          followersCount: 7400,
          articlesCount: 51
        }
      ]
      
      totalResults.value = articles.value.length + users.value.length
      hasMoreResults.value = true
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
