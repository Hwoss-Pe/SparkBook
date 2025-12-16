import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

// 创建axios实例
const resolveBaseURL = () => {
  const envBase = (import.meta as any).env?.VITE_API_BASE as string | undefined
  const defaultBase = new URL('/api', window.location.origin).toString()
  let base = envBase || defaultBase
  try {
    const u = new URL(base, window.location.origin)
    if (window.location.protocol === 'https:' && u.protocol !== 'https:') {
      u.protocol = 'https:'
      base = u.toString()
    }
    return base
  } catch {
    return defaultBase
  }
}

const service: AxiosInstance = axios.create({
  baseURL: resolveBaseURL(),
  timeout: 15000,
})

// 是否正在刷新 token
let isRefreshing = false
// 重试队列，每一项是一个函数
let requests: any[] = []

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    // 在发送请求之前做些什么
    const token = localStorage.getItem('token')
    if (window.location.protocol !== 'https:') {
      console.warn('当前非HTTPS环境，仍将附带Authorization头用于开发调试')
    }
    console.log('=== HTTP 请求拦截器 ===')
    console.log('发送请求:', config.url)
    console.log('请求方法:', config.method?.toUpperCase())
    console.log('请求数据:', config.data)
    console.log('请求参数:', config.params)
    console.log('localStorage 中的 token:', token ? token.substring(0, 50) + '...' : 'null')
    if (config.data instanceof FormData) {
      config.headers = config.headers || {}
      config.headers['Content-Type'] = 'multipart/form-data'
    }
    
    // 特别关注个人信息接口
    if (config.url?.includes('/users/profile')) {
      console.log('🔍 这是个人信息接口请求')
      console.log('完整token:', token)
    }
    
    // 如果是刷新 token 的请求，使用 refresh token
    if (config.url?.includes('/users/refresh_token')) {
      const refreshToken = localStorage.getItem('refreshToken')
      if (refreshToken) {
        config.headers['Authorization'] = `Bearer ${refreshToken}`
        console.log('刷新Token请求，使用 refreshToken')
      }
    } else if (token) {
      // 普通请求使用 access token
      // 后端从 Authorization 头提取 token，格式为 "Bearer token"
      config.headers['Authorization'] = `Bearer ${token}`
      console.log('已添加 Authorization 头:', `Bearer ${token.substring(0, 20)}...`)
    } else {
      console.log('没有 token，未添加 Authorization 头')
    }
    return config
  },
  (error) => {
    // 对请求错误做些什么
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse) => {
    console.log('=== HTTP 响应拦截器 ===')
    console.log('响应URL:', response.config.url)
    console.log('响应状态:', response.status)
    console.log('响应状态文本:', response.statusText)
    
    // 特别关注个人信息接口
    if (response.config.url?.includes('/users/profile')) {
      console.log('🔍 这是个人信息接口响应')
      console.log('个人信息响应数据:', response.data)
      console.log('个人信息响应数据类型:', typeof response.data)
      console.log('个人信息响应详细:', JSON.stringify(response.data, null, 2))
    }
    
    // 从响应头获取 token 并保存
    // 注意：响应头的 key 会被浏览器转为小写
    const jwtToken = response.headers['x-jwt-token']
    const refreshToken = response.headers['x-refresh-token']
    
    console.log('响应头:', response.headers)
    console.log('x-jwt-token:', jwtToken)
    console.log('x-refresh-token:', refreshToken)
    
    if (jwtToken) {
      localStorage.setItem('token', jwtToken)
      console.log('token 已保存到 localStorage')
    }
    if (refreshToken) {
      localStorage.setItem('refreshToken', refreshToken)
    }
    
    const res = response.data
    
    // 如果是文件下载等二进制数据，直接返回
    if (response.config.responseType === 'blob' || response.config.responseType === 'arraybuffer') {
      return response
    }
    
    // 根据后端API的响应结构进行处理
    // 判断响应是否成功
    if (res.msg === '登录成功' || res.code === 0 || response.status === 200) {
      // 对于个人信息接口，直接返回数据
      if (response.config.url?.includes('/users/profile')) {
        console.log('✅ 个人信息接口响应成功，直接返回数据')
        return res
      }
      return res.data || res
    } else {
      ElMessage.error(res.msg || res.Msg || '请求失败')
      
      // 处理特定错误码
      if (res.code === 401) {
        // 未授权，需要重新登录
        localStorage.removeItem('token')
        localStorage.removeItem('refreshToken')
        router.push('/login')
      }
      
      return Promise.reject(new Error(res.msg || res.Msg || '请求失败'))
    }
  },
  (error) => {
    // 处理HTTP错误
    let message = '网络错误，请稍后重试'
    
    if (error.response) {
      const config = error.config
      switch (error.response.status) {
        case 400:
          message = '请求参数错误'
          break
        case 401:
          // 如果是刷新 token 的请求失败，直接登出
          if (config.url?.includes('/users/refresh_token')) {
            message = '登录已过期，请重新登录'
            localStorage.removeItem('token')
            localStorage.removeItem('refreshToken')
            router.push('/login')
            break
          }

          // 如果不是刷新请求，尝试刷新 token
          if (!isRefreshing) {
            isRefreshing = true
            // 尝试刷新 token
            return service.post('/users/refresh_token', {})
              .then(() => {
                // 刷新成功，重试队列中的请求
                requests.forEach(cb => cb())
                requests = []
                
                // 重试当前请求
                config.headers['Authorization'] = `Bearer ${localStorage.getItem('token')}`
                return service(config)
              })
              .catch(refreshErr => {
                console.error('Refresh token failed:', refreshErr)
                // 刷新失败，清空队列并登出
                requests = []
                localStorage.removeItem('token')
                localStorage.removeItem('refreshToken')
                router.push('/login')
                return Promise.reject(refreshErr)
              })
              .finally(() => {
                isRefreshing = false
              })
          } else {
            // 正在刷新，将请求加入队列
            return new Promise((resolve) => {
              requests.push(() => {
                config.headers['Authorization'] = `Bearer ${localStorage.getItem('token')}`
                resolve(service(config))
              })
            })
          }
        case 403:
          message = '拒绝访问'
          break
        case 404:
          message = '请求的资源不存在'
          break
        case 500:
          message = '服务器内部错误'
          break
        default:
          message = `请求失败: ${error.response.status}`
      }
    } else if (error.request) {
      message = '服务器无响应'
    }
    
    // 只有在非 401 错误，或者 401 且是刷新 token 失败时才提示错误
    if (error.response?.status !== 401 || (error.response?.status === 401 && error.config?.url?.includes('/users/refresh_token'))) {
      ElMessage.error(message)
    }
    
    console.error('Response error:', error)
    return Promise.reject(error)
  }
)

// 封装GET请求
export const get = <T>(url: string, params?: any, config?: AxiosRequestConfig): Promise<T> => {
  return service.get(url, { params, ...config })
}

// 封装POST请求
export const post = <T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
  return service.post(url, data, config)
}

// 封装PUT请求
export const put = <T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
  return service.put(url, data, config)
}

// 封装DELETE请求
export const del = <T>(url: string, params?: any, config?: AxiosRequestConfig): Promise<T> => {
  return service.delete(url, { params, ...config })
}

export default service

export const resolveStaticUrl = (url: string): string => {
  if (!url) return ''
  if (/^https?:\/\//i.test(url)) return url
  try {
    const apiBase = resolveBaseURL()
    const backendOrigin = apiBase.replace(/\/?api\/?$/, '')
    return new URL(url, backendOrigin).toString()
  } catch {
    return url
  }
}
