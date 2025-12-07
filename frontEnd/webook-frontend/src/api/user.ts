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
  user: User
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
    return get<ProfileResponse>('/users/profile')
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

