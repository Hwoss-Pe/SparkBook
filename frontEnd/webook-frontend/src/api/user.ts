import { get, post } from './http'

// 用户相关接口类型定义
export interface User {
  id: number
  email: string
  nickname: string
  phone: string
  aboutMe: string
  birthday?: string
  wechatInfo?: {
    openId: string
    unionId: string
  }
}

export interface LoginRequest {
  email: string
  password: string
}

export interface SignupRequest {
  email: string
  password: string
  confirmPassword: string
}

export interface LoginResponse {
  user: User
  token: string
  refreshToken: string
}

export interface ProfileResponse {
  id?: number
  email?: string
  nickname?: string
  phone?: string
  aboutMe?: string
  birthday?: string
  // 后端实际返回的字段名（大写开头）
  Email?: string
  Nickname?: string
  Phone?: string
  AboutMe?: string
  Birthday?: string
}

// 用户相关API
export const userApi = {
  // 登录
  login: (data: LoginRequest) => {
    return post<LoginResponse>('/users/login', data)
  },
  
  // 注册
  signup: (data: SignupRequest) => {
    return post('/users/signup', data)
  },
  
  // 获取用户信息
  getProfile: () => {
    console.log('=== 发起个人信息请求 ===')
    console.log('请求URL: /users/profile')
    console.log('请求方法: GET')
    console.log('当前时间:', new Date().toISOString())
    
    return get<ProfileResponse>('/users/profile').then(response => {
      console.log('=== 个人信息请求成功 ===')
      console.log('响应数据:', response)
      console.log('响应数据类型:', typeof response)
      console.log('响应数据详细:', JSON.stringify(response, null, 2))
      return response
    }).catch(error => {
      console.error('=== 个人信息请求失败 ===')
      console.error('请求错误:', error)
      console.error('错误详情:', JSON.stringify(error, null, 2))
      throw error
    })
  },
  
  // 更新用户信息
  updateProfile: (user: Partial<User>) => {
    return post('/users/edit', user)
  },
  
  // 发送短信验证码
  sendSmsCode: (phone: string) => {
    return post('/users/login_sms/code/send', { phone })
  },
  
  // 手机号登录/注册
  loginByPhone: (phone: string, code: string) => {
    return post<LoginResponse>('/users/login_sms', { phone, code })
  },
  
  // 登出
  logout: () => {
    return post('/users/logout', {})
  },
  
  // 刷新 token
  refreshToken: () => {
    return post('/users/refresh_token', {})
  }
}

