import { get, post } from './http'

// 关注相关接口类型定义
export interface FollowRelation {
  id: number
  follower: number // 关注者
  followee: number // 被关注者
  name?: string
  avatar?: string
  about_me?: string
}

export interface FollowStatic {
  followers: number // 粉丝数
  followees: number // 关注数
}

export interface GetFolloweeRequest {
  follower: number
  offset: number
  limit: number
}

export interface GetFolloweeResponse {
  follow_relations: FollowRelation[]
}

export interface GetFollowerRequest {
  followee: number
  offset: number
  limit: number
}

export interface GetFollowerResponse {
  follow_relations: FollowRelation[]
}

export interface FollowInfoRequest {
  follower: number
  followee: number
}

export interface FollowInfoResponse {
  follow_relation: FollowRelation
}

export interface GetFollowStaticRequest {
  followee: number
}

export interface GetFollowStaticResponse {
  followStatic: FollowStatic
}

// 关注相关API
export const followApi = {
  // 关注
  follow: (followee: number, follower: number) => {
    return post('/follow', { followee, follower })
  },
  
  // 取消关注
  cancelFollow: (followee: number, follower: number) => {
    return post('/follow/cancel', { followee, follower })
  },
  
  // 获取关注列表
  getFollowee: (params: GetFolloweeRequest) => {
    return get<GetFolloweeResponse>('/follow/followee', params)
  },
  
  // 获取粉丝列表
  getFollower: (params: GetFollowerRequest) => {
    return get<GetFollowerResponse>('/follow/follower', params)
  },
  
  // 获取关注信息
  getFollowInfo: async (params: FollowInfoRequest) => {
    try {
      return await get<FollowInfoResponse>('/follow/info', params)
    } catch (e) {
      return { follow_relation: undefined as any }
    }
  },
  
  // 获取关注统计
  getFollowStatics: (params: GetFollowStaticRequest) => {
    return get<GetFollowStaticResponse>('/follow/statics', params)
  }
}

