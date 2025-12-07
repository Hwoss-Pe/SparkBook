import { get, post } from './http'

// 交互相关接口类型定义
export interface Interactive {
  biz: string
  biz_id: number
  read_cnt: number
  like_cnt: number
  collect_cnt: number
  liked: boolean
  collected: boolean
}

export interface GetInteractiveRequest {
  biz: string
  biz_id: number
  uid: number
}

export interface GetInteractiveResponse {
  intr: Interactive
}

export interface GetInteractiveByIdsRequest {
  biz: string
  ids: number[]
}

export interface GetInteractiveByIdsResponse {
  intrs: Record<number, Interactive>
}

// 交互相关API
export const interactiveApi = {
  // 获取交互数据
  getInteractive: (params: GetInteractiveRequest) => {
    return get<GetInteractiveResponse>('/interactive/get', params)
  },
  
  // 批量获取交互数据
  getInteractiveByIds: (params: GetInteractiveByIdsRequest) => {
    return get<GetInteractiveByIdsResponse>('/interactive/batch', params)
  },
  
  // 增加阅读计数
  incrReadCnt: (biz: string, biz_id: number) => {
    return post('/interactive/read', { biz, biz_id })
  },
  
  // 点赞
  like: (biz: string, biz_id: number, uid: number) => {
    return post('/interactive/like', { biz, biz_id, uid })
  },
  
  // 取消点赞
  cancelLike: (biz: string, biz_id: number, uid: number) => {
    return post('/interactive/like/cancel', { biz, biz_id, uid })
  },
  
  // 收藏
  collect: (biz: string, biz_id: number, uid: number, cid: number) => {
    return post('/interactive/collect', { biz, biz_id, uid, cid })
  }
}

