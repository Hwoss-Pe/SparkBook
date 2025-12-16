import { get, post } from './http'

export type MessageCategory = 'interaction' | 'follow' | 'system'

export interface NotificationSender {
  id: number
  name?: string
  avatar?: string
}

export interface NotificationTarget {
  type?: string
  id?: number
  title?: string
  preview?: string
}

export interface NotificationItem {
  id: number
  category: MessageCategory
  content: string
  time: string
  sender?: NotificationSender
  target?: NotificationTarget
  status: 'unread' | 'read'
}

export interface UnreadCountsResponse {
  interaction: number
  follow: number
  system: number
  total: number
}

export interface ListRequest {
  type: MessageCategory
  offset?: number
  limit?: number
}

export type ListResponse = NotificationItem[]

export interface MarkReadRequest {
  ids?: number[]
  type?: MessageCategory
}

export const notificationApi = {
  getUnreadCounts: () => {
    return get<UnreadCountsResponse>('/notifications/unread_counts')
  },
  list: (params: ListRequest) => {
    return get<ListResponse>('/notifications', params)
  },
  markRead: (data: MarkReadRequest) => {
    return post<Record<string, any>>('/notifications/mark_read', data)
  }
}

