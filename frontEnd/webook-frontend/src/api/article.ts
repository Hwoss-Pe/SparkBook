import { get, post } from './http'

// 作者信息
export interface Author {
  id: number
  name: string
  avatar: string
}

// 文章相关接口类型定义
export interface Article {
  id: number
  title: string
  content: string
  abstract: string
  coverImage: string
  author: Author
  status: number
  ctime: string
  utime: string
  readCnt?: number
  likeCnt?: number
  collectCnt?: number
  tags?: string[]
}

// 推荐文章（首页使用）
export interface ArticlePub {
  id: number
  title: string
  abstract: string
  coverImage: string
  author: Author
  ctime: string
  utime: string
  readCnt: number
  likeCnt: number
  collectCnt: number
  liked?: boolean
  collected?: boolean
}

export interface ArticleDetail extends Article {
  // 扩展字段
  readCnt: number
  likeCnt: number
  collectCnt: number
  liked: boolean
  collected: boolean
}

export interface AuthorStats {
  publishedCount: number
  draftCount: number
  totalReadCount: number
  totalLikeCount: number
  followingCount: number
  followerCount: number
}

// 列表接口返回的是文章数组
export type ListResponse = Article[]

export interface ListRequest {
  offset: number
  limit: number
}

export interface PublishedListRequest {
  offset: number
  limit: number
}

// 文章相关API
export const articleApi = {
  // 获取文章列表（作者视角）
  getList: (params: ListRequest) => {
    return post<ListResponse>('/articles/list', params)
  },
  // 根据官方标签获取文章列表（匿名可访问）
  getArticlesByOfficialTag: (params: { tag: string; offset?: number; limit?: number }) => {
    return get<ArticlePub[]>('/articles/pub/tag/articles', params)
  },
  
  // 获取推荐文章列表（首页使用，按时间排序）
  getRecommendList: (params: PublishedListRequest) => {
    return post<ArticlePub[]>('/articles/pub/list', params)
  },
  // 获取关注作者的文章列表（需要登录）
  getFollowingList: (params: PublishedListRequest) => {
    return get<ArticlePub[]>('/articles/pub/following/list', params)
  },
  // 按作者获取已发布文章列表（读者视角，匿名可访问）
  getAuthorPublishedList: (authorId: number, params: PublishedListRequest) => {
    return get<ArticlePub[]>(`/articles/pub/author/${authorId}/list`, params)
  },
  
  // 获取用户收藏的文章列表（需要登录）
  getCollectedList: (params: PublishedListRequest) => {
    return get<ArticlePub[]>(`/articles/pub/collected/list`, params)
  },

  // AI 自动生成摘要和标题（兼容旧调用，但推荐使用 generateAI）
  generateSummary: (content: string) => {
    return post<{ title: string; abstract: string }>('/articles/generate', { content, type: 'generate' }, { timeout: 120000 })
  },

  // AI 智能助手通用接口
  generateAI: (data: { content: string; type: 'generate' | 'polish' | 'tag'; instruction?: string }) => {
    return post<{ 
      title?: string; 
      abstract?: string; 
      content?: string; 
      tags?: string[] 
    }>('/articles/generate', data, { timeout: 120000 })
  },
  
  // 获取文章详情（作者视角）
  getArticleById: (id: number) => {
    return get<ArticleDetail>(`/articles/detail/${id}`)
  },
  
  // 获取已发布文章详情（读者视角）
  getPublishedArticleById: (id: number) => {
    return get<ArticleDetail>(`/articles/pub/${id}`)
  },
  
  // 保存文章（草稿）
  saveArticle: (article: Partial<Article>) => {
    return post<number>('/articles/edit', article)
  },
  
  // 发布文章
  publishArticle: (article: Partial<Article>) => {
    return post<number>('/articles/publish', article)
  },
  
  // 撤回文章
  withdrawArticle: (id: number, uid: number) => {
    return post('/articles/withdraw', { id, uid })
  },
  
  unpublishArticle: (id: number, uid: number) => {
    return post('/articles/unpublish', { id, uid })
  },
  
  deleteDraft: (id: number, uid: number) => {
    return post('/articles/unpublish', { id, uid })
  },
  
  // 点赞
  like: (id: number) => {
    return post('/articles/pub/like', { id, like: true })
  },
  cancelLike: (id: number) => {
    return post('/articles/pub/cancelLike', { id, like: false })
  },
  
  // 收藏
  collect: (id: number, cid: number) => {
    return post('/articles/pub/collect', { id, cid })
  },
  cancelCollect: (id: number, cid: number) => {
    return post('/articles/pub/cancelCollect', { id, cid })
  },

  getAuthorStats: (authorId: number) => {
    return get<AuthorStats>(`/articles/author/${authorId}/stats`)
  },
  // 获取官方标签（uid=0）
  getOfficialTags: () => {
    return get<string[]>(`/articles/tags/official`)
  },
}
