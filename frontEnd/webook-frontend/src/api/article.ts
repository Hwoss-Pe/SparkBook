import { get, post } from './http'

// 作者信息
export interface Author {
  id: number
  name: string
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
}

export interface ArticleDetail extends Article {
  // 扩展字段
  readCnt: number
  likeCnt: number
  collectCnt: number
  liked: boolean
  collected: boolean
}

export interface ListResponse {
  articles: Article[]
}

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
  getList: (params: ListRequest & { author: number }) => {
    return get<ListResponse>('/articles/list', params)
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
  
  // 点赞
  like: (id: number, like: boolean) => {
    return post('/articles/pub/like', { id, like })
  },
  
  // 收藏
  collect: (id: number, cid: number) => {
    return post('/articles/pub/collect', { id, cid })
  }
}

