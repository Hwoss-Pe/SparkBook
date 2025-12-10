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
  
  // 获取推荐文章列表（首页使用，按时间排序）
  getRecommendList: (params: PublishedListRequest) => {
    return post<ArticlePub[]>('/articles/pub/list', params)
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
  }
}
