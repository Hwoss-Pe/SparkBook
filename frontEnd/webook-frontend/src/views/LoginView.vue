<template>
  <div class="login-container">
    <div class="login-box">
      <div class="logo">
        <h1>小微书</h1>
        <p>探索、分享、创造</p>
      </div>
      
      <el-tabs v-model="activeTab" class="login-tabs">
        <el-tab-pane label="手机号登录" name="phone">
          <el-form :model="phoneForm" :rules="phoneRules" ref="phoneFormRef" @submit.prevent="handlePhoneLogin">
            <el-form-item prop="phone">
              <el-input v-model="phoneForm.phone" placeholder="请输入手机号" prefix-icon="el-icon-mobile">
                <template #prefix>
                  <el-icon><Iphone /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item prop="code">
              <div class="code-input">
                <el-input v-model="phoneForm.code" placeholder="请输入验证码" prefix-icon="el-icon-key">
                  <template #prefix>
                    <el-icon><Key /></el-icon>
                  </template>
                </el-input>
                <el-button 
                  type="primary" 
                  :disabled="codeSending || countdown > 0" 
                  @click="handleSendCode"
                >
                  {{ countdown > 0 ? `${countdown}秒后重新获取` : '获取验证码' }}
                </el-button>
              </div>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" class="login-button" @click="handlePhoneLogin" :loading="phoneLoading">
                登录/注册
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="邮箱密码登录" name="email">
          <el-form :model="emailForm" :rules="emailRules" ref="emailFormRef" @submit.prevent="handleEmailLogin">
            <el-form-item prop="email">
              <el-input v-model="emailForm.email" placeholder="请输入邮箱" prefix-icon="el-icon-message">
                <template #prefix>
                  <el-icon><Message /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item prop="password">
              <el-input v-model="emailForm.password" type="password" placeholder="请输入密码" prefix-icon="el-icon-lock" show-password>
                <template #prefix>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item>
              <div class="form-actions">
                <el-checkbox v-model="rememberMe">记住我</el-checkbox>
                <el-link type="primary" @click="handleForgotPassword">忘记密码？</el-link>
              </div>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" class="login-button" @click="handleEmailLogin" :loading="emailLoading">
                登录
              </el-button>
            </el-form-item>
            
            <el-form-item>
              <div class="register-link">
                <span>还没有账号？</span>
                <el-link type="primary" @click="activeTab = 'register'">立即注册</el-link>
              </div>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="注册" name="register">
          <el-form :model="registerForm" :rules="registerRules" ref="registerFormRef" @submit.prevent="handleRegister">
            <el-form-item prop="email">
              <el-input v-model="registerForm.email" placeholder="请输入邮箱" prefix-icon="el-icon-message">
                <template #prefix>
                  <el-icon><Message /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item prop="nickname">
              <el-input v-model="registerForm.nickname" placeholder="请输入昵称" prefix-icon="el-icon-user">
                <template #prefix>
                  <el-icon><User /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item prop="phone">
              <el-input v-model="registerForm.phone" placeholder="请输入手机号" prefix-icon="el-icon-mobile">
                <template #prefix>
                  <el-icon><Iphone /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item prop="password">
              <el-input v-model="registerForm.password" type="password" placeholder="请输入密码" prefix-icon="el-icon-lock" show-password>
                <template #prefix>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item prop="confirmPassword">
              <el-input v-model="registerForm.confirmPassword" type="password" placeholder="请确认密码" prefix-icon="el-icon-lock" show-password>
                <template #prefix>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" class="login-button" @click="handleRegister" :loading="registerLoading">
                注册
              </el-button>
            </el-form-item>
            
            <el-form-item>
              <div class="register-link">
                <span>已有账号？</span>
                <el-link type="primary" @click="activeTab = 'email'">立即登录</el-link>
              </div>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
      
      <div class="other-login">
        <div class="divider">
          <span>其他登录方式</span>
        </div>
        <div class="social-login">
          <el-button class="wechat-btn" @click="handleWechatLogin">
            <el-icon><Promotion /></el-icon>
            微信登录
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Iphone, Key, Message, Lock, User, Promotion } from '@element-plus/icons-vue'
import { userApi, codeApi } from '@/api'

const router = useRouter()
const route = useRoute()

// 表单引用
const phoneFormRef = ref()
const emailFormRef = ref()
const registerFormRef = ref()

// 当前激活的标签页
const activeTab = ref('phone')

// 记住我
const rememberMe = ref(false)

// 手机登录表单
const phoneForm = reactive({
  phone: '',
  code: ''
})

// 邮箱登录表单
const emailForm = reactive({
  email: '',
  password: ''
})

// 注册表单
const registerForm = reactive({
  email: '',
  nickname: '',
  phone: '',
  password: '',
  confirmPassword: ''
})

// 表单校验规则
const phoneRules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { pattern: /^\d{6}$/, message: '验证码为6位数字', trigger: 'blur' }
  ]
}

const emailRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码不能少于6个字符', trigger: 'blur' }
  ]
}

const registerRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '昵称长度在2到20个字符之间', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码不能少于6个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: any) => {
        if (value !== registerForm.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 加载状态
const phoneLoading = ref(false)
const emailLoading = ref(false)
const registerLoading = ref(false)
const codeSending = ref(false)
const countdown = ref(0)

// 发送验证码
const handleSendCode = async () => {
  try {
    await phoneFormRef.value.validateField('phone')
    
    codeSending.value = true
    
    // 调用发送验证码API
    await codeApi.sendCode({
      biz: 'login',
      phone: phoneForm.phone
    })
    
    ElMessage.success('验证码已发送，请注意查收')
    
    // 开始倒计时
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error) {
    console.error('发送验证码失败:', error)
  } finally {
    codeSending.value = false
  }
}

// 手机号登录
const handlePhoneLogin = async () => {
  try {
    await phoneFormRef.value.validate()
    
    phoneLoading.value = true
    
    // 调用手机号登录API
    const res = await userApi.loginByPhone(phoneForm.phone, phoneForm.code)
    
    // 保存登录状态
    localStorage.setItem('token', res.token)
    localStorage.setItem('refreshToken', res.refreshToken)
    localStorage.setItem('user', JSON.stringify(res.user))
    
    ElMessage.success('登录成功')
    
    // 跳转到首页或重定向页面
    const redirectUrl = route.query.redirect as string || '/'
    router.push(redirectUrl)
  } catch (error) {
    console.error('手机号登录失败:', error)
  } finally {
    phoneLoading.value = false
  }
}

// 邮箱登录
const handleEmailLogin = async () => {
  try {
    await emailFormRef.value.validate()
    
    emailLoading.value = true
    
    // 调用邮箱登录API
    const res = await userApi.login({
      email: emailForm.email,
      password: emailForm.password
    })
    
    // 保存登录状态
    localStorage.setItem('token', res.token)
    localStorage.setItem('refreshToken', res.refreshToken)
    localStorage.setItem('user', JSON.stringify(res.user))
    
    // 记住我
    if (rememberMe.value) {
      localStorage.setItem('remember', 'true')
      localStorage.setItem('email', emailForm.email)
    } else {
      localStorage.removeItem('remember')
      localStorage.removeItem('email')
    }
    
    ElMessage.success('登录成功')
    
    // 跳转到首页或重定向页面
    const redirectUrl = route.query.redirect as string || '/'
    router.push(redirectUrl)
  } catch (error) {
    console.error('邮箱登录失败:', error)
  } finally {
    emailLoading.value = false
  }
}

// 注册
const handleRegister = async () => {
  try {
    await registerFormRef.value.validate()
    
    registerLoading.value = true
    
    // 调用注册API
    await userApi.signup({
      user: {
        email: registerForm.email,
        nickname: registerForm.nickname,
        password: registerForm.password,
        phone: registerForm.phone
      }
    })
    
    ElMessage.success('注册成功，请登录')
    
    // 切换到登录标签页
    activeTab.value = 'email'
    emailForm.email = registerForm.email
  } catch (error) {
    console.error('注册失败:', error)
  } finally {
    registerLoading.value = false
  }
}

// 忘记密码
const handleForgotPassword = () => {
  ElMessage.info('请通过手机号登录或联系客服重置密码')
}

// 微信登录
const handleWechatLogin = async () => {
  // 生成随机state
  const state = Math.random().toString(36).substring(2)
  localStorage.setItem('wechat_state', state)
  
  // 跳转到微信授权页面
  // 实际项目中应该调用后端接口获取授权URL
  window.location.href = `/api/oauth2/auth?state=${state}`
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.login-box {
  width: 400px;
  padding: 30px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.logo {
  text-align: center;
  margin-bottom: 30px;
}

.logo h1 {
  margin: 0;
  font-size: 28px;
  color: #ff2442;
}

.logo p {
  margin: 10px 0 0;
  color: #666;
}

.login-tabs {
  margin-bottom: 20px;
}

.code-input {
  display: flex;
  gap: 10px;
}

.login-button {
  width: 100%;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.register-link {
  text-align: center;
  font-size: 14px;
}

.other-login {
  margin-top: 30px;
}

.divider {
  display: flex;
  align-items: center;
  margin: 20px 0;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background-color: #ddd;
}

.divider span {
  padding: 0 10px;
  color: #999;
  font-size: 14px;
}

.social-login {
  display: flex;
  justify-content: center;
}

.wechat-btn {
  width: 100%;
  background-color: #07c160;
  border-color: #07c160;
  color: white;
}

.wechat-btn:hover {
  background-color: #06ad56;
  border-color: #06ad56;
}
</style>
