import { get } from './http'

// 搜索相关接口类型定义（与后端真实返回保持一致）
export interface User {
  id: number
  email: string
  nickname: string
  phone: string
  avatar?: string
}

export interface Article {
  id: number
  title: string
  status: number
  content: string
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
