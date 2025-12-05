import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

// 创建axios实例
const service: AxiosInstance = axios.create({
  baseURL: '/api', // API的base_url
  timeout: 15000, // 请求超时时间
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    // 在发送请求之前做些什么
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['x-jwt-token'] = token
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
    const res = response.data
    
    // 如果是文件下载等二进制数据，直接返回
    if (response.config.responseType === 'blob' || response.config.responseType === 'arraybuffer') {
      return response
    }
    
    // 根据后端API的响应结构进行处理
    // 假设后端返回的数据结构为 { code: number, data: any, msg: string }
    if (res.code === 0) {
      return res.data
    } else {
      ElMessage.error(res.msg || '请求失败')
      
      // 处理特定错误码
      if (res.code === 401) {
        // 未授权，需要重新登录
        localStorage.removeItem('token')
        localStorage.removeItem('refreshToken')
        router.push('/login')
      }
      
      return Promise.reject(new Error(res.msg || '请求失败'))
    }
  },
  (error) => {
    // 处理HTTP错误
    let message = '网络错误，请稍后重试'
    
    if (error.response) {
      switch (error.response.status) {
        case 400:
          message = '请求参数错误'
          break
        case 401:
          message = '未授权，请重新登录'
          localStorage.removeItem('token')
          localStorage.removeItem('refreshToken')
          router.push('/login')
          break
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
    
    ElMessage.error(message)
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
