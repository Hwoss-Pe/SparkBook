import { ref, onMounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { articleApi, ArticleDetail } from '@/api/article'
import { interactiveApi } from '@/api/interactive'
import { commentApi, Comment } from '@/api/comment'

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

  // 提交评论
  const submitComment = async () => {
    if (!commentContent.value.trim()) return
    
    try {
      const currentUserId = 101 // 实际应该从用户状态中获取
      
      // 调用创建评论API
      await commentApi.createComment({
        comment: {
          uid: currentUserId,
          biz: 'article',
          bizid: article.value.id,
          content: commentContent.value
        }
      })
      
      // 创建本地评论对象
      const newComment: CommentData = {
        id: Date.now(),
        content: commentContent.value,
        createTime: new Date(),
        likeCount: 0,
        isLiked: false,
        user: {
          id: currentUserId,
          name: '当前用户',
          avatar: ''
        },
        replies: []
      }
      
      comments.value.unshift(newComment)
      article.value.commentCount++
      commentContent.value = ''
      ElMessage.success('评论发表成功')
    } catch (error) {
      console.error('发表评论失败:', error)
      ElMessage.error('评论发表失败，请稍后重试')
    }
  }

  // 点赞评论
  const likeComment = async (comment: CommentData | ReplyData) => {
    try {
      const currentUserId = 101 // 实际应该从用户状态中获取
      
      // 这里应该调用点赞评论API
      comment.isLiked = !comment.isLiked
      comment.likeCount += comment.isLiked ? 1 : -1
    } catch (error) {
      console.error('点赞评论失败:', error)
      ElMessage.error('操作失败，请稍后重试')
    }
  }

  // 回复评论
  const replyToComment = (comment: CommentData | ReplyData, parentComment: CommentData | null = null) => {
    // 设置评论框的内容为回复格式
    commentContent.value = `回复 @${comment.user.name}：`
    scrollToComments()
  }

  // 加载更多评论
  const loadMoreComments = async () => {
    try {
      // 调用API加载更多评论
      const response = await commentApi.getCommentList({
        biz: 'article',
        bizid: article.value.id,
        min_id: comments.value.length > 0 ? comments.value[comments.value.length - 1].id : 0,
        limit: 10
      })
      
      // 将API返回的评论转换为本地格式
      const moreComments: CommentData[] = response.comments.map(comment => ({
        id: comment.id,
        content: comment.content,
        createTime: new Date(comment.ctime),
        likeCount: 5, // 实际应该从交互API获取
        isLiked: false,
        user: {
          id: comment.uid,
          name: `用户${comment.uid}`,
          avatar: `https://picsum.photos/id/${1000 + comment.uid}/100/100`
        },
        replies: [] // 实际应该从API获取回复
      }))
      
      comments.value = [...comments.value, ...moreComments]
      
      // 如果返回的评论数量小于请求的数量，说明没有更多评论了
      if (response.comments.length < 10) {
        hasMoreComments.value = false
      }
    } catch (error) {
      console.error('加载更多评论失败:', error)
      ElMessage.error('加载评论失败，请稍后重试')
      
      // 使用模拟数据
      const moreComments = [
        {
          id: Date.now() - 1000,
          content: '这是加载的更多评论',
          createTime: new Date(Date.now() - 3600000),
          likeCount: 5,
          isLiked: false,
          user: {
            id: 888,
            name: '评论用户',
            avatar: 'https://picsum.photos/id/1005/100/100'
          },
          replies: []
        }
      ]
      
      comments.value = [...comments.value, ...moreComments]
      hasMoreComments.value = false
    }
  }

  // 获取文章详情
  const fetchArticleDetail = async () => {
    loading.value = true
    error.value = ''
    
    try {
      // 调用API获取文章详情
      const articleDetail = await articleApi.getPublishedArticleById(articleId.value)
      
      // 调用API获取交互数据
      const currentUserId = 101 // 实际应该从用户状态中获取
      const interactiveData = await interactiveApi.getInteractive({
        biz: 'article',
        biz_id: articleId.value,
        uid: currentUserId
      })
      
      // 增加阅读计数
      await interactiveApi.incrReadCnt('article', articleId.value)
      
      // 构建文章数据
      article.value = {
        id: articleDetail.id,
        title: articleDetail.title,
        content: articleDetail.content,
        coverImage: `https://picsum.photos/id/${400 + articleDetail.id}/800/400`, // 实际应该从文章信息中获取
        publishTime: new Date(articleDetail.ctime),
        readCount: interactiveData.intr.read_cnt,
        likeCount: interactiveData.intr.like_cnt,
        commentCount: 128, // 实际应该从评论API获取
        favoriteCount: interactiveData.intr.collect_cnt,
        isLiked: interactiveData.intr.liked,
        isFavorited: interactiveData.intr.collected,
        isFollowed: false, // 实际应该从关注API获取
        author: {
          id: articleDetail.author.id,
          name: articleDetail.author.name,
          avatar: `https://picsum.photos/id/${1000 + articleDetail.author.id}/100/100` // 实际应该从用户信息中获取
        },
        tags: ['美食', '甜点', '意大利菜', '烘焙'] // 实际应该从文章信息中获取
      }
      
      // 获取评论
      const commentsResponse = await commentApi.getCommentList({
        biz: 'article',
        bizid: articleId.value,
        min_id: 0,
        limit: 10
      })
      
      // 将API返回的评论转换为本地格式
      comments.value = commentsResponse.comments.map(comment => ({
        id: comment.id,
        content: comment.content,
        createTime: new Date(comment.ctime),
        likeCount: 42, // 实际应该从交互API获取
        isLiked: false,
        user: {
          id: comment.uid,
          name: `用户${comment.uid}`,
          avatar: `https://picsum.photos/id/${1000 + comment.uid}/100/100`
        },
        replies: [] // 实际应该从API获取回复
      }))
      
      // 如果返回的评论数量等于请求的数量，说明可能还有更多评论
      hasMoreComments.value = commentsResponse.comments.length === 10
      
      // 获取相关文章
      // 实际项目中应该有相关文章的API
      relatedArticles.value = [
        {
          id: 10,
          title: '意大利经典甜点大全：从提拉米苏到奶油冰淇淋',
          coverImage: 'https://picsum.photos/id/432/200/200',
          author: { name: '美食达人' },
          readCount: 9800
        },
        {
          id: 11,
          title: '在家制作专业级咖啡的秘诀',
          coverImage: 'https://picsum.photos/id/766/200/200',
          author: { name: '咖啡师' },
          readCount: 8500
        },
        {
          id: 12,
          title: '10种最适合新手的简易甜点食谱',
          coverImage: 'https://picsum.photos/id/835/200/200',
          author: { name: '烘焙教练' },
          readCount: 12000
        }
      ]
      
      loading.value = false
    } catch (error) {
      console.error('Failed to load article:', error)
      loading.value = false
      error.value = '文章加载失败，请稍后再试'
      
      // 使用模拟数据
      setTimeout(() => {
        article.value = {
          id: articleId.value,
          title: '如何在家制作完美的提拉米苏',
          content: `
            <p>提拉米苏是一道经典的意大利甜点，口感醇厚，香气四溢。今天我将分享如何在家中制作一份完美的提拉米苏。</p>
            <h2>食材准备</h2>
            <ul>
              <li>马斯卡彭奶酪 - 250克</li>
              <li>手指饼干 - 200克</li>
              <li>浓缩咖啡 - 1杯</li>
              <li>鸡蛋 - 3个</li>
              <li>白砂糖 - 80克</li>
              <li>可可粉 - 适量</li>
              <li>朗姆酒 - 2勺（可选）</li>
            </ul>
            <h2>制作步骤</h2>
            <p>1. 将鸡蛋分离，蛋黄和蛋白分开备用。</p>
            <p>2. 将蛋黄与40克白砂糖混合，搅拌至颜色变浅，质地变得浓稠。</p>
            <p>3. 加入马斯卡彭奶酪，继续搅拌至完全融合。</p>
            <p>4. 在另一个干净的碗中，将蛋白打发至起泡，逐渐加入剩余的白砂糖，打至形成硬性发泡。</p>
            <p>5. 将打发的蛋白轻轻折入奶酪混合物中，注意保持蓬松感。</p>
            <p>6. 将浓缩咖啡倒入浅盘中，如果喜欢可以加入朗姆酒。</p>
            <p>7. 将手指饼干浸入咖啡中约1-2秒，然后排列在容器底部。</p>
            <p>8. 在饼干层上均匀涂抹一半的奶酪混合物。</p>
            <p>9. 重复一次饼干层和奶酪层。</p>
            <p>10. 在最上层撒上可可粉。</p>
            <p>11. 将提拉米苏放入冰箱冷藏至少4小时，最好隔夜。</p>
            <h2>小贴士</h2>
            <p>- 手指饼干浸泡咖啡的时间不要太长，否则会变得太湿软。</p>
            <p>- 使用室温的马斯卡彭奶酪更容易搅拌均匀。</p>
            <p>- 如果没有马斯卡彭奶酪，可以用奶油奶酪加上少量的鲜奶油代替。</p>
            <p>希望这个配方能帮助你在家中制作出美味的提拉米苏！</p>
          `,
          coverImage: 'https://picsum.photos/id/431/800/400',
          publishTime: new Date(Date.now() - 86400000 * 3),
          readCount: 12500,
          likeCount: 3200,
          commentCount: 128,
          favoriteCount: 980,
          isLiked: false,
          isFavorited: false,
          isFollowed: false,
          author: {
            id: 101,
            name: '美食达人',
            avatar: 'https://picsum.photos/id/1027/100/100'
          },
          tags: ['美食', '甜点', '意大利菜', '烘焙']
        }
        
        // 模拟评论数据
        comments.value = [
          {
            id: 1001,
            content: '按照这个配方做了一次，非常成功！家人都很喜欢，谢谢分享！',
            createTime: new Date(Date.now() - 86400000),
            likeCount: 42,
            isLiked: false,
            user: {
              id: 201,
              name: '烘焙爱好者',
              avatar: 'https://picsum.photos/id/1005/100/100'
            },
            replies: [
              {
                id: 2001,
                content: '很高兴你喜欢这个配方！你可以尝试加入一点香草精，会有不一样的风味。',
                createTime: new Date(Date.now() - 86400000 + 3600000),
                likeCount: 15,
                isLiked: false,
                user: {
                  id: 101,
                  name: '美食达人',
                  avatar: 'https://picsum.photos/id/1027/100/100'
                },
                replyTo: {
                  id: 201,
                  name: '烘焙爱好者'
                }
              }
            ]
          },
          {
            id: 1002,
            content: '请问马斯卡彭奶酪可以用什么替代吗？我这边买不到。',
            createTime: new Date(Date.now() - 86400000 * 2),
            likeCount: 18,
            isLiked: false,
            user: {
              id: 202,
              name: '甜点新手',
              avatar: 'https://picsum.photos/id/1012/100/100'
            },
            replies: [
              {
                id: 2002,
                content: '可以用奶油奶酪混合一些鲜奶油来替代，口感会略有不同，但也很美味！',
                createTime: new Date(Date.now() - 86400000 * 2 + 7200000),
                likeCount: 10,
                isLiked: false,
                user: {
                  id: 101,
                  name: '美食达人',
                  avatar: 'https://picsum.photos/id/1027/100/100'
                },
                replyTo: {
                  id: 202,
                  name: '甜点新手'
                }
              }
            ]
          }
        ]
        
        // 模拟相关文章
        relatedArticles.value = [
          {
            id: 10,
            title: '意大利经典甜点大全：从提拉米苏到奶油冰淇淋',
            coverImage: 'https://picsum.photos/id/432/200/200',
            author: { name: '美食达人' },
            readCount: 9800
          },
          {
            id: 11,
            title: '在家制作专业级咖啡的秘诀',
            coverImage: 'https://picsum.photos/id/766/200/200',
            author: { name: '咖啡师' },
            readCount: 8500
          },
          {
            id: 12,
            title: '10种最适合新手的简易甜点食谱',
            coverImage: 'https://picsum.photos/id/835/200/200',
            author: { name: '烘焙教练' },
            readCount: 12000
          }
        ]
        
        hasMoreComments.value = true
        loading.value = false
      }, 1000)
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
