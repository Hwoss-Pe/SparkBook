import { ref, onMounted, nextTick, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { articleApi, type ArticleDetail } from '@/api/article'
import { commentApi, type Comment } from '@/api/comment'
import { useUserStore } from '@/stores/user'
import { resolveStaticUrl } from '@/api/http'

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
  isLikeAnimating?: boolean;
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
  isLikeAnimating?: boolean;
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
  const replyTarget = ref<{ rootId: number; parentId: number; userName: string } | null>(null)
  const commentsSection = ref<HTMLElement | null>(null)
  const relatedArticles = ref<RelatedArticle[]>([])
  const minCommentId = ref(0)
  const pageSize = 10
  const userStore = useUserStore()
  const replyMaxId = ref<Record<number, number>>({})
  const hasMoreRepliesMap = ref<Record<number, boolean>>({})
  const repliesFirstLoadDone = ref<Record<number, boolean>>({})
  const threadModal = ref<{ visible: boolean; rootId: number }>({ visible: false, rootId: 0 })

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
    const d = typeof date === 'string' ? new Date(date.replace(' ', 'T')) : new Date(date)
    if (isNaN(d.getTime()) || d.getFullYear() < 1971) {
      return '-'
    }
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
      if (article.value.isLiked) {
        await articleApi.cancelLike(article.value.id)
        article.value.isLiked = false
        article.value.likeCount = Math.max(0, article.value.likeCount - 1)
      } else {
        await articleApi.like(article.value.id)
        article.value.isLiked = true
        article.value.likeCount += 1
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
      const defaultCid = 1
      if (article.value.isFavorited) {
        await articleApi.cancelCollect(article.value.id, defaultCid)
        article.value.isFavorited = false
        article.value.favoriteCount = Math.max(0, article.value.favoriteCount - 1)
      } else {
        await articleApi.collect(article.value.id, defaultCid)
        article.value.isFavorited = true
        article.value.favoriteCount += 1
      }

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

  const toReplyData = (r: Comment): ReplyData => ({
    id: r.id,
    content: r.content,
    createTime: typeof r.ctime === 'string' ? new Date(r.ctime.replace(' ', 'T')) : new Date(r.ctime),
    likeCount: 0,
    isLiked: false,
    isLikeAnimating: false,
    user: {
      id: r.uid,
      name: r.name || `用户#${r.uid}`,
      avatar: r.avatar ? resolveStaticUrl(r.avatar) : `https://picsum.photos/id/${1000 + r.uid}/100/100`
    },
    replyTo: r.parent_comment?.id
      ? { id: r.parent_comment.id, name: r.parent_comment?.name || (r.parent_comment?.uid ? `用户#${r.parent_comment.uid}` : `用户#${r.parent_comment.id}`) }
      : undefined
  })

  const flattenReplies = (items: Comment[]): ReplyData[] => {
    const out: ReplyData[] = []
    const stack: Comment[] = [...items]
    while (stack.length) {
      const cur = stack.shift() as Comment
      out.push(toReplyData(cur))
      if (cur.children && cur.children.length) {
        stack.push(...cur.children)
      }
    }
    return out
  }

  const threadOrder = (replies: ReplyData[], rootId: number): ReplyData[] => {
    const seconds = replies.filter(r => r.replyTo?.id === rootId)
      .sort((a, b) => b.createTime.getTime() - a.createTime.getTime())
    const rest = replies.filter(r => r.replyTo?.id !== rootId)
    const childrenMap = new Map<number, ReplyData[]>()
    for (const r of rest) {
      const pid = r.replyTo?.id
      if (pid == null) continue
      const arr = childrenMap.get(pid) || []
      arr.push(r)
      childrenMap.set(pid, arr)
    }
    for (const [k, arr] of childrenMap.entries()) {
      arr.sort((a, b) => b.createTime.getTime() - a.createTime.getTime())
      childrenMap.set(k, arr)
    }
    const out: ReplyData[] = []
    const emitChildren = (pid: number) => {
      const list = childrenMap.get(pid) || []
      for (const child of list) {
        out.push(child)
        emitChildren(child.id)
      }
    }
    for (const s of seconds) {
      out.push(s)
      emitChildren(s.id)
    }
    const emitted = new Set(out.map(r => r.id))
    const remain = replies.filter(r => !emitted.has(r.id))
      .sort((a, b) => b.createTime.getTime() - a.createTime.getTime())
    return [...out, ...remain]
  }

  const calcDepth = (r: ReplyData, rootId: number, parentOf: Map<number, number>): number => {
    let depth = 2
    let pid = r.replyTo?.id
    const guard = new Set<number>()
    while (pid && pid !== rootId && !guard.has(pid)) {
      guard.add(pid)
      depth += 1
      pid = parentOf.get(pid)
    }
    return depth
  }

  const getThreadView = (rootId: number): Array<{ reply: ReplyData; depth: number }> => {
    const c = comments.value.find(x => x.id === rootId)
    if (!c) return []
    const flat = (c.replies || []).slice().sort((a, b) => b.createTime.getTime() - a.createTime.getTime())
    return flat.map(r => ({ reply: r, depth: 2 }))
  }

  const openThreadModal = async (rootId: number) => {
    threadModal.value.visible = true
    threadModal.value.rootId = rootId
    if (!repliesFirstLoadDone.value[rootId]) {
      await loadMoreRepliesFor(rootId)
    }
  }

  const closeThreadModal = () => {
    threadModal.value.visible = false
  }

  const mapComment = (c: Comment): CommentData => {
    const replies: ReplyData[] = threadOrder(flattenReplies(c.children || []), c.id)
    return {
      id: c.id,
      content: c.content,
      createTime: typeof c.ctime === 'string' ? new Date(c.ctime.replace(' ', 'T')) : new Date(c.ctime),
      likeCount: 0,
      isLiked: false,
      isLikeAnimating: false,
      user: {
        id: c.uid,
        name: c.name || `用户#${c.uid}`,
        avatar: c.avatar ? resolveStaticUrl(c.avatar) : `https://picsum.photos/id/${1000 + c.uid}/100/100`
      },
      replies
    }
  }

  const fetchComments = async (initial = false) => {
    try {
      const res = await commentApi.getCommentList({
        biz: 'article',
        bizid: articleId.value,
        min_id: initial ? 0 : minCommentId.value,
        limit: pageSize
      })
      const list = (res.comments || []).map(mapComment)
      if (initial) {
        comments.value = list
        // 初始化每个根评论的回复分页游标为当前已展示的最小ID
        list.forEach(c => {
          const direct = (res.comments || []).find(rc => rc.id === c.id)?.children || []
          replyMaxId.value[c.id] = Number.MAX_SAFE_INTEGER
          hasMoreRepliesMap.value[c.id] = direct.length >= 3
          repliesFirstLoadDone.value[c.id] = false
        })
      } else {
        comments.value = [...comments.value, ...list]
        list.forEach(c => {
          if (!(c.id in replyMaxId.value)) {
            const direct = (res.comments || []).find(rc => rc.id === c.id)?.children || []
            replyMaxId.value[c.id] = Number.MAX_SAFE_INTEGER
            hasMoreRepliesMap.value[c.id] = direct.length >= 3
            repliesFirstLoadDone.value[c.id] = false
          }
        })
      }
      if (list.length > 0) {
        minCommentId.value = list[list.length - 1].id
      }
      hasMoreComments.value = list.length === pageSize
      article.value.commentCount = comments.value.length
    } catch (error) {
      console.error('加载评论失败:', error)
      ElMessage.error('加载评论失败')
    }
  }

  const submitComment = async () => {
    try {
      const content = commentContent.value.trim()
      if (!content) return
      if (!userStore.user?.id) {
        ElMessage.error('请先登录后再发表评论')
        return
      }
      const payload: any = {
        comment: {
          uid: userStore.user.id,
          biz: 'article',
          bizid: articleId.value,
          content
        }
      }
      if (replyTarget.value) {
        payload.comment.root_comment = { id: replyTarget.value.rootId }
        payload.comment.parent_comment = { id: replyTarget.value.parentId }
      }
      await commentApi.createComment(payload)
      commentContent.value = ''
      replyTarget.value = null
      await fetchComments(true)
      ElMessage.success('评论已发布')
    } catch (error) {
      console.error('发表评论失败:', error)
      ElMessage.error('发表评论失败')
    }
  }

  // 点赞评论（暂时禁用）
  const likeComment = async (comment: CommentData | ReplyData) => {
    comment.isLikeAnimating = true
    if (comment.isLiked) {
      comment.isLiked = false
      comment.likeCount = Math.max(0, comment.likeCount - 1)
    } else {
      comment.isLiked = true
      comment.likeCount += 1
    }
    setTimeout(() => {
      comment.isLikeAnimating = false
    }, 300)
  }

  // 回复评论（暂时禁用）
  const replyToComment = (comment: CommentData | ReplyData, parentComment: CommentData | null = null) => {
    if ('replies' in comment) {
      replyTarget.value = { rootId: comment.id, parentId: comment.id, userName: comment.user.name || `用户#${comment.user.id}` }
    } else if (parentComment) {
      replyTarget.value = { rootId: parentComment.id, parentId: comment.id, userName: comment.user.name || `用户#${comment.user.id}` }
    }
    scrollToComments()
  }

  const loadMoreComments = async () => {
    await fetchComments(false)
  }

  const mapReply = (r: Comment): ReplyData => toReplyData(r)

  const loadMoreRepliesFor = async (rootId: number) => {
    try {
      const isFirstLoad = !repliesFirstLoadDone.value[rootId]
      const maxId = isFirstLoad ? Number.MAX_SAFE_INTEGER : (replyMaxId.value[rootId] ?? Number.MAX_SAFE_INTEGER)
      const pageLimit = 10
      const res = await commentApi.getMoreReplies({ rid: rootId, max_id: maxId, limit: pageLimit })
      const fetched = flattenReplies(res.replies || [])
      repliesFirstLoadDone.value[rootId] = true
      const idx = comments.value.findIndex(c => c.id === rootId)
      if (idx < 0) return
      const existingIds = new Set(comments.value[idx].replies.map(r => r.id))
      const add = fetched.filter(r => !existingIds.has(r.id))
      if (add.length > 0) {
        comments.value[idx].replies = [...comments.value[idx].replies, ...add]
        // 更新游标为当前已加载的最小ID，确保下一页只拿更旧的数据
        const currentMin = replyMaxId.value[rootId] ?? Number.MAX_SAFE_INTEGER
        const rootIds = (res.replies || []).map(r => r.id)
        const nextMin = rootIds.length > 0 ? Math.min(...rootIds) : currentMin
        replyMaxId.value[rootId] = Math.min(currentMin, nextMin)
        hasMoreRepliesMap.value[rootId] = (res.replies || []).length === pageLimit
      } else {
        // 没有更多数据，隐藏按钮
        const currentMin = replyMaxId.value[rootId] ?? Number.MAX_SAFE_INTEGER
        const rootIds = (res.replies || []).map(r => r.id)
        const nextMin = rootIds.length > 0 ? Math.min(...rootIds) : currentMin
        replyMaxId.value[rootId] = Math.min(currentMin, nextMin)
        hasMoreRepliesMap.value[rootId] = (res.replies || []).length === pageLimit
      }
    } catch (error) {
      console.error('加载更多回复失败:', error)
      ElMessage.error('加载更多回复失败')
    }
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
        coverImage: articleDetail.coverImage ? resolveStaticUrl(articleDetail.coverImage) : '',
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
          avatar: articleDetail.author.avatar ? resolveStaticUrl(articleDetail.author.avatar) : 'https://picsum.photos/seed/avatar/100/100'
        },
        tags: (articleDetail as any).tags || []
      }

      await fetchComments(true)

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
    replyTarget,
    loadMoreRepliesFor,
    hasMoreRepliesMap,
    threadModal,
    getThreadView,
    openThreadModal,
    closeThreadModal,
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
