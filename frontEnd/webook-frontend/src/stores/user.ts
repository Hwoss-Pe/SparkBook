import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { User } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  // 用户信息
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const isLoggedIn = ref(false)

  // 初始化用户状态
  function initUserState() {
    const storedToken = localStorage.getItem('token')
    const storedUser = localStorage.getItem('user')
    const storedRefreshToken = localStorage.getItem('refreshToken')

    if (storedToken && storedUser) {
      token.value = storedToken
      refreshToken.value = storedRefreshToken
      user.value = JSON.parse(storedUser)
      isLoggedIn.value = true
    }
  }

  // 设置用户信息
  function setUser(userData: User, userToken: string, userRefreshToken: string) {
    user.value = userData
    token.value = userToken
    refreshToken.value = userRefreshToken
    isLoggedIn.value = true

    // 保存到本地存储
    localStorage.setItem('token', userToken)
    localStorage.setItem('refreshToken', userRefreshToken)
    localStorage.setItem('user', JSON.stringify(userData))
  }

  // 清除用户信息
  function clearUser() {
    user.value = null
    token.value = null
    refreshToken.value = null
    isLoggedIn.value = false

    // 清除本地存储
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('user')
  }

  return { user, token, refreshToken, isLoggedIn, initUserState, setUser, clearUser }
})
