import { get, post } from './http'

// 热榜相关接口类型定义
export interface RankingArticle {
  id: number
  title: string
  abstract: string
  coverImage: string
  author: {
    id: number
    name: string
    avatar: string
  }
  ctime: string
  utime: string
  readCnt: number
  likeCnt: number
  collectCnt: number
}

export interface RankingRequest {
  offset?: number
  limit?: number
}

// 实际返回为后端 Result 的 data 字段，Axios 拦截器已直接返回数组
export type RankingResponse = RankingArticle[]

// 热榜相关API
export const rankingApi = {
  // 获取热榜文章
  getRanking: (params: RankingRequest = { offset: 0, limit: 10 }) => {
    return get<RankingResponse>('/articles/pub/ranking', params)
  },
  
  // 手动触发热榜计算
  triggerRanking: () => {
    return post<Record<string, any>>('/articles/pub/ranking/trigger', {})
  }
}
