import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { notificationApi, type MessageCategory, type NotificationItem, type UnreadCountsResponse } from '@/api/notification'

export const useNotificationStore = defineStore('notification', () => {
  const unread = ref<UnreadCountsResponse>({ interaction: 0, follow: 0, system: 0, total: 0 })
  const interaction = ref<NotificationItem[]>([])
  const follow = ref<NotificationItem[]>([])
  const system = ref<NotificationItem[]>([])
  const offsets = ref<Record<MessageCategory, number>>({ interaction: 0, follow: 0, system: 0 })
  const hasMore = ref<Record<MessageCategory, boolean>>({ interaction: true, follow: true, system: true })

  const totalUnread = computed(() => unread.value.total)

  async function fetchUnreadCounts() {
    const res = await notificationApi.getUnreadCounts()
    unread.value = res
  }

  async function fetchMessages(type: MessageCategory, reset = false, limit = 10) {
    if (reset) {
      offsets.value[type] = 0
      hasMore.value[type] = true
      if (type === 'interaction') interaction.value = []
      else if (type === 'follow') follow.value = []
      else system.value = []
    }
    const resp = await notificationApi.list({ type, offset: offsets.value[type], limit })
    if (type === 'interaction') interaction.value = [...interaction.value, ...resp]
    else if (type === 'follow') follow.value = [...follow.value, ...resp]
    else system.value = [...system.value, ...resp]
    offsets.value[type] += resp.length
    hasMore.value[type] = resp.length === limit
  }

  async function markReadByCategory(type: MessageCategory) {
    await notificationApi.markRead({ type })
    await fetchUnreadCounts()
  }

  async function markReadByIds(ids: number[]) {
    await notificationApi.markRead({ ids })
    await fetchUnreadCounts()
  }

  return {
    unread,
    interaction,
    follow,
    system,
    hasMore,
    totalUnread,
    fetchUnreadCounts,
    fetchMessages,
    markReadByCategory,
    markReadByIds
  }
})

