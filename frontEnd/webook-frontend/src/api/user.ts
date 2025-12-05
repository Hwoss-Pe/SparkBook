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
  user: {
    email: string
    nickname: string
    password: string
    phone: string
  }
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
    return post<LoginResponse>('/user/login', data)
  },
  
  // 注册
  signup: (data: SignupRequest) => {
    return post('/user/signup', data)
  },
  
  // 获取用户信息
  getProfile: (id: number) => {
    return get<ProfileResponse>(`/user/profile/${id}`)
  },
  
  // 更新用户信息
  updateProfile: (user: Partial<User>) => {
    return post('/user/profile/update', { user })
  },
  
  // 手机号登录/注册
  loginByPhone: (phone: string, code: string) => {
    return post<LoginResponse>('/user/login/phone', { phone, code })
  },
  
  // 微信登录
  loginByWechat: (code: string, state: string) => {
    return post<LoginResponse>('/user/login/wechat', { code, state })
  }
}
