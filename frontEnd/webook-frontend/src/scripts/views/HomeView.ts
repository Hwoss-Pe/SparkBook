import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

export default function useHomeView() {
  const router = useRouter()

  // 模拟Banner数据
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

  // 模拟文章数据
  const articles = ref([
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
      commentCount: 128
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
      commentCount: 342
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
      commentCount: 98
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
      commentCount: 215
    },
    {
      id: 5,
      title: '家庭花园：从零开始的种植指南',
      abstract: '无论你是有着一片后院，还是只有一个小阳台，都可以打造自己的绿色天地...',
      coverImage: 'https://picsum.photos/id/145/400/300',
      author: {
        id: 105,
        name: '园艺爱好者',
        avatar: 'https://picsum.photos/id/1074/100/100'
      },
      readCount: 7600,
      likeCount: 1800,
      commentCount: 76
    },
    {
      id: 6,
      title: '如何提高你的摄影技巧：从入门到精通',
      abstract: '无需昂贵的设备，掌握这些基本技巧，你也能拍出令人惊艳的照片...',
      coverImage: 'https://picsum.photos/id/250/400/300',
      author: {
        id: 106,
        name: '摄影师小王',
        avatar: 'https://picsum.photos/id/1062/100/100'
      },
      readCount: 21000,
      likeCount: 6300,
      commentCount: 430
    }
  ])

  // 模拟热榜数据
  const hotRankings = ref([
    {
      id: 101,
      title: '年轻人为什么都在做副业？',
      author: { name: '经济观察家' },
      readCount: 125000
    },
    {
      id: 102,
      title: '这些小众景点比网红打卡地更值得去',
      author: { name: '旅行达人' },
      readCount: 98000
    },
    {
      id: 103,
      title: '如何在30天内养成一个新习惯',
      author: { name: '生活改造师' },
      readCount: 87000
    },
    {
      id: 104,
      title: '2025年最值得学习的5个技能',
      author: { name: '职场导师' },
      readCount: 76000
    },
    {
      id: 105,
      title: '这样做，让你的朋友圈不再单调',
      author: { name: '社交达人' },
      readCount: 65000
    },
    {
      id: 106,
      title: '低成本高品质的家居改造指南',
      author: { name: '家居设计师' },
      readCount: 54000
    },
    {
      id: 107,
      title: '如何科学健身：避开这些常见误区',
      author: { name: '健身教练' },
      readCount: 43000
    },
    {
      id: 108,
      title: '读懂财报的10个关键指标',
      author: { name: '投资顾问' },
      readCount: 32000
    }
  ])

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

  // 加载更多文章
  const loadMoreArticles = () => {
    // 这里应该调用API加载更多文章
    // 模拟加载更多
    const moreArticles = [
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
        commentCount: 320
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
        commentCount: 145
      }
    ]
    
    articles.value = [...articles.value, ...moreArticles]
    
    // 假设没有更多文章了
    hasMoreArticles.value = false
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
    // 这里应该调用API获取首页数据
    // 目前使用模拟数据
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
