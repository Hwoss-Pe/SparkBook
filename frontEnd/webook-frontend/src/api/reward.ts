import { get, post } from './http'

// 打赏相关接口类型定义
export enum RewardStatus {
  Unknown = 0,
  Init = 1,
  Payed = 2,
  Failed = 3
}

export interface PreRewardRequest {
  biz: string
  biz_id: number
  biz_name: string
  target_uid: number // 被打赏的人
  uid: number // 打赏的人
  amt: number // 打赏金额
}

export interface PreRewardResponse {
  code_url: string
  rid: number
}

export interface GetRewardRequest {
  rid: number
  uid: number
}

export interface GetRewardResponse {
  status: RewardStatus
}

// 打赏相关API
export const rewardApi = {
  // 预打赏
  preReward: (data: PreRewardRequest) => {
    return post<PreRewardResponse>('/reward/pre', data)
  },
  
  // 获取打赏状态
  getReward: (params: GetRewardRequest) => {
    return get<GetRewardResponse>('/reward/get', params)
  }
}
