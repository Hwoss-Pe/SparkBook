import { get } from './http'

// 热榜相关接口类型定义
export interface RankingArticle {
  id: number
  title: string
  status: number
  content: string
  author: {
    id: number
    name: string
  }
  ctime: string
  utime: string
}

export interface TopNResponse {
  articles: RankingArticle[]
}

// 热榜相关API
export const rankingApi = {
  // 获取热榜文章
  getTopN: () => {
    return get<TopNResponse>('/ranking/top')
  },
  
  // 计算热榜（通常是后台任务，前端一般不调用）
  rankTopN: () => {
    return get('/ranking/calculate')
  }
}

