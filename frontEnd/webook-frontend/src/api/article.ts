import { get, post } from './http'

// 文章相关接口类型定义
export interface Article {
  id: number
  title: string
  content: string
  abstract: string
  author: {
    id: number
    name: string
  }
  status: number
  ctime: string
  utime: string
}

export interface ArticleDetail extends Article {
  // 扩展字段
}

export interface ListResponse {
  articles: Article[]
}

export interface ListRequest {
  offset: number
  limit: number
}

export interface PublishedListRequest {
  start_time?: string
  offset: number
  limit: number
}

// 文章相关API
export const articleApi = {
  // 获取文章列表（作者视角）
  getList: (params: ListRequest & { author: number }) => {
    return get<ListResponse>('/article/list', params)
  },
  
  // 获取已发布文章列表
  getPublishedList: (params: PublishedListRequest) => {
    return get<ListResponse>('/article/list/published', params)
  },
  
  // 获取文章详情
  getArticleById: (id: number) => {
    return get<ArticleDetail>(`/article/${id}`)
  },
  
  // 获取已发布文章详情
  getPublishedArticleById: (id: number, uid?: number) => {
    return get<ArticleDetail>(`/article/published/${id}`, { uid })
  },
  
  // 保存文章（草稿）
  saveArticle: (article: Partial<Article>) => {
    return post<{ id: number }>('/article/save', article)
  },
  
  // 发布文章
  publishArticle: (article: Partial<Article>) => {
    return post<{ id: number }>('/article/publish', article)
  },
  
  // 撤回文章
  withdrawArticle: (id: number, uid: number) => {
    return post('/article/withdraw', { id, uid })
  }
}
