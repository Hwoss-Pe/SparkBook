import { ref, onMounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { articleApi, type ArticleDetail } from '@/api/article'
import { interactiveApi } from '@/api/interactive'
import { commentApi, type Comment } from '@/api/comment'

// 定义类型接口
interface ArticleData {
  id: number;
  title: string;
  content: string;
  coverImage: string;
  publishTime: Date;
  readCount: number;
  likeCount: number;
  commentCount: number;
  favoriteCount: number;
  isLiked: boolean;
  isFavorited: boolean;
  isFollowed: boolean;
  author: {
    id: number;
    name: string;
    avatar: string;
  };
  tags: string[];
}

interface CommentData {
  id: number;
  content: string;
  createTime: Date;
  likeCount: number;
  isLiked: boolean;
  user: {
    id: number;
    name: string;
    avatar: string;
  };
  replies: ReplyData[];
}

interface ReplyData {
  id: number;
  content: string;
  createTime: Date;
  likeCount: number;
  isLiked: boolean;
  user: {
    id: number;
    name: string;
    avatar: string;
  };
  replyTo?: {
    id: number;
    name: string;
  };
}

interface RelatedArticle {
  id: number;
  title: string;
  coverImage?: string;
  author: {
    name: string;
  };
  readCount: number;
}

export default function useArticleDetailView() {
  const router = useRouter()
  const route = useRoute()
  const articleId = ref(Number(route.params.id))

  const loading = ref(true)
  const error = ref('')
  const article = ref<ArticleData>({
    id: 0,
    title: '',
    content: '',
    coverImage: '',
    publishTime: new Date(),
    readCount: 0,
    likeCount: 0,
    commentCount: 0,
    favoriteCount: 0,
    isLiked: false,
    isFavorited: false,
    isFollowed: false,
    author: {
      id: 0,
      name: '',
      avatar: ''
    },
    tags: []
  })

  const comments = ref<CommentData[]>([])
  const hasMoreComments = ref(false)
  const commentContent = ref('')
  const commentsSection = ref<HTMLElement | null>(null)
  const relatedArticles = ref<RelatedArticle[]>([])

  // 格式化数字，例如1200显示为1.2k
  const formatNumber = (num: number): string => {
    if (num >= 10000) {
      return (num / 10000).toFixed(1) + 'w'
    } else if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'k'
    }
    return num.toString()
  }

  // 格式化日期时间
  const formatDateTime = (date: Date | string): string => {
    const d = new Date(date)
    return d.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  // 返回上一页
  const goBack = () => {
    router.back()
  }

  // 查看作者主页
  const viewAuthor = (authorId: number) => {
    router.push(`/user/${authorId}`)
  }

  // 查看文章详情
  const viewArticle = (id: number) => {
    router.push(`/article/${id}`)
  }

  // 关注/取消关注作者
  const toggleFollow = async () => {
    try {
      // 这里应该调用关注/取消关注API
      article.value.isFollowed = !article.value.isFollowed
      ElMessage.success(article.value.isFollowed ? '已关注作者' : '已取消关注')
    } catch (error) {
      console.error('关注操作失败:', error)
      ElMessage.error('操作失败，请稍后重试')
    }
  }

  // 点赞/取消点赞
  const toggleLike = async () => {
    try {
      const currentUserId = 101 // 实际应该从用户状态中获取
      
      if (article.value.isLiked) {
        // 取消点赞
        await interactiveApi.cancelLike('article', article.value.id, currentUserId)
        article.value.isLiked = false
        article.value.likeCount--
      } else {
        // 点赞
        await interactiveApi.like('article', article.value.id, currentUserId)
        article.value.isLiked = true
        article.value.likeCount++
      }
      
      ElMessage.success(article.value.isLiked ? '已点赞' : '已取消点赞')
    } catch (error) {
      console.error('点赞操作失败:', error)
      ElMessage.error('操作失败，请稍后重试')
    }
  }

  // 收藏/取消收藏
  const toggleFavorite = async () => {
    try {
      const currentUserId = 101 // 实际应该从用户状态中获取
      
      // 这里应该调用收藏/取消收藏API
      if (!article.value.isFavorited) {
        await interactiveApi.collect('article', article.value.id, currentUserId, 1) // 1是默认收藏夹ID
      }
      
      article.value.isFavorited = !article.value.isFavorited
      article.value.favoriteCount += article.value.isFavorited ? 1 : -1
      
      ElMessage.success(article.value.isFavorited ? '已收藏' : '已取消收藏')
    } catch (error) {
      console.error('收藏操作失败:', error)
      ElMessage.error('操作失败，请稍后重试')
    }
  }

  // 滚动到评论区
  const scrollToComments = () => {
    nextTick(() => {
      if (commentsSection.value) {
        commentsSection.value.scrollIntoView({ behavior: 'smooth' })
      }
    })
  }

  // 提交评论（暂时禁用）
  const submitComment = async () => {
    // 暂时不支持评论功能
    ElMessage.info('评论功能暂未开放')
  }

  // 点赞评论（暂时禁用）
  const likeComment = async (comment: CommentData | ReplyData) => {
    // 暂时不支持评论功能
    ElMessage.info('评论功能暂未开放')
  }

  // 回复评论（暂时禁用）
  const replyToComment = (comment: CommentData | ReplyData, parentComment: CommentData | null = null) => {
    // 暂时不支持评论功能
    ElMessage.info('评论功能暂未开放')
  }

  // 加载更多评论（暂时禁用）
  const loadMoreComments = async () => {
    // 暂时不加载评论数据
    ElMessage.info('评论功能暂未开放')
  }

  // 获取文章详情
  const fetchArticleDetail = async () => {
    loading.value = true
    error.value = ''
    
    try {
      // 调用API获取文章详情（包含交互数据）
      // 后端会自动处理阅读计数增加，前端无需调用
      const articleDetail = await articleApi.getPublishedArticleById(articleId.value)
      
      // 构建文章数据（直接使用API返回的交互数据）
      article.value = {
        id: articleDetail.id,
        title: articleDetail.title,
        content: articleDetail.content,
        coverImage: articleDetail.coverImage || '', // 使用API返回的封面图片
        publishTime: new Date(articleDetail.ctime),
        readCount: articleDetail.readCnt,
        likeCount: articleDetail.likeCnt,
        commentCount: 0, // 暂时不显示评论数量
        favoriteCount: articleDetail.collectCnt,
        isLiked: articleDetail.liked,
        isFavorited: articleDetail.collected,
        isFollowed: false, // 实际应该从关注API获取
        author: {
          id: articleDetail.author.id,
          name: articleDetail.author.name,
          avatar: '' // 移除假的头像，使用默认头像
        },
        tags: [] // 移除假标签，等待后端支持标签功能
      }
      
      // 暂时不加载评论数据
      comments.value = []
      hasMoreComments.value = false
      
      // 获取相关文章
      // TODO: 实际项目中应该有相关文章的API
      relatedArticles.value = [] // 暂时为空，等待相关文章API
      
      loading.value = false
    } catch (err) {
      console.error('Failed to load article:', err)
      loading.value = false
      error.value = '文章加载失败，请稍后再试'
    }
  }

  onMounted(() => {
    fetchArticleDetail()
  })

  return {
    articleId,
    loading,
    error,
    article,
    comments,
    hasMoreComments,
    commentContent,
    commentsSection,
    relatedArticles,
    formatNumber,
    formatDateTime,
    goBack,
    viewAuthor,
    viewArticle,
    toggleFollow,
    toggleLike,
    toggleFavorite,
    scrollToComments,
    submitComment,
    likeComment,
    replyToComment,
    loadMoreComments
  }
}
