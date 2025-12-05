import { get } from './http'

// 搜索相关接口类型定义
export interface User {
  id: number
  nickname: string
  aboutMe: string
  // 其他用户字段
}

export interface Article {
  id: number
  title: string
  abstract: string
  author: {
    id: number
    name: string
  }
  // 其他文章字段
}

export interface SearchRequest {
  expression: string
  uid: number
}

export interface SearchResponse {
  user: {
    users: User[]
  }
  article: {
    articles: Article[]
  }
}

// 搜索相关API
export const searchApi = {
  // 搜索
  search: (params: SearchRequest) => {
    return get<SearchResponse>('/search', params)
  }
}
