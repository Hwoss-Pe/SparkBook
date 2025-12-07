import { get, post } from './http'

// 验证码相关接口类型定义
export interface CodeSendRequest {
  biz: string
  phone: string
}

export interface VerifyRequest {
  biz: string
  phone: string
  inputCode: string
}

export interface VerifyResponse {
  answer: boolean
}

// 验证码相关API
export const codeApi = {
  // 发送验证码
  sendCode: (data: CodeSendRequest) => {
    return post('/code/send', data)
  },
  
  // 验证验证码
  verifyCode: (data: VerifyRequest) => {
    return post<VerifyResponse>('/code/verify', data)
  }
}

