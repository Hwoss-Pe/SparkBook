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

export interface RankingResponse {
  map(arg0: (article: { id: any; title: any; author: { name: any }; readCnt: any }) => { avatar:any;id: any; title: any; author: { name: any }; readCount: any }): { id: number; title: string; author: { name: string }; readCount: number }[] | { id: number; title: string; author: { name: string }; readCount: number }[]
  code: number
  msg: string
  data: RankingArticle[]
}

// 热榜相关API
export const rankingApi = {
  // 获取热榜文章
  getRanking: (params: RankingRequest = { offset: 0, limit: 10 }) => {
    return get<RankingResponse>('/articles/pub/ranking', params)
  },
  
  // 手动触发热榜计算
  triggerRanking: () => {
    return post<{ code: number; msg: string; data: any }>('/articles/pub/ranking/trigger', {})
  }
}

