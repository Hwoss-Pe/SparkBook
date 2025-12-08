import { get, post } from './http'

// 评论相关接口类型定义
export interface Comment {
  id: number
  uid: number
  biz: string
  bizid: number
  content: string
  root_comment?: Comment
  parent_comment?: Comment
  children?: Comment[]
  ctime: string
  utime: string
}

export interface CommentListRequest {
  biz: string
  bizid: number
  min_id: number
  limit: number
}

export interface CommentListResponse {
  comments: Comment[]
}

export interface CreateCommentRequest {
  comment: {
    uid: number
    biz: string
    bizid: number
    content: string
    root_comment?: {
      id: number
    }
    parent_comment?: {
      id: number
    }
  }
}

export interface GetMoreRepliesRequest {
  rid: number
  max_id: number
  limit: number
}

export interface GetMoreRepliesResponse {
  replies: Comment[]
}

// 评论相关API
export const commentApi = {
  // 获取评论列表
  getCommentList: (params: CommentListRequest) => {
    return get<CommentListResponse>('/comment/list', params)
  },

  // 创建评论
  createComment: (data: CreateCommentRequest) => {
    return post('/comment/create', data)
  },

  // 删除评论
  deleteComment: (id: number) => {
    return post('/comment/delete', { id })
  },

  // 获取更多回复
  getMoreReplies: (params: GetMoreRepliesRequest) => {
    return get<GetMoreRepliesResponse>('/comment/replies', params)
  }
}
