import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { userApi } from '@/api/user'
import { codeApi } from '@/api/code'
import { useUserStore } from '@/stores/user'

// 定义类型接口
interface EmailForm {
  email: string;
  password: string;
}

interface PhoneForm {
  phone: string;
  code: string;
}

interface SignupForm {
  email: string;
  password: string;
  confirmPassword: string;
}

export default function useLoginView() {
  const router = useRouter()
  const route = useRoute()
  const userStore = useUserStore()

  // 标签页
  const activeTab = ref('email')

  // 邮箱登录表单
  const emailFormRef = ref<FormInstance>()
  const emailForm = reactive<EmailForm>({
    email: '',
    password: ''
  })
  const emailRules: FormRules = {
    email: [
      { required: true, message: '请输入邮箱', trigger: 'blur' },
      { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
    ],
    password: [
      { required: true, message: '请输入密码', trigger: 'blur' },
      { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
    ]
  }
  const emailLoading = ref(false)
  const rememberMe = ref(false)

  // 手机号登录表单
  const phoneFormRef = ref<FormInstance>()
  const phoneForm = reactive<PhoneForm>({
    phone: '',
    code: ''
  })
  const phoneRules: FormRules = {
    phone: [
      { required: true, message: '请输入手机号', trigger: 'blur' },
      { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
    ],
    code: [
      { required: true, message: '请输入验证码', trigger: 'blur' },
      { min: 6, max: 6, message: '验证码长度为6位', trigger: 'blur' }
    ]
  }
  const phoneLoading = ref(false)
  const codeSent = ref(false)
  const countdown = ref(0)
  let timer: number | null = null

  // 登录成功后的通用处理
  const handleLoginSuccess = async () => {
    try {
      console.log('=== 开始获取用户个人信息 ===')
      console.log('当前localStorage中的token:', localStorage.getItem('token'))
      
      // 获取用户信息
      const userProfile = await userApi.getProfile()
      console.log('=== 个人信息接口响应数据 ===')
      console.log('userProfile 原始响应:', userProfile)
      console.log('userProfile 类型:', typeof userProfile)
      console.log('userProfile 详细内容:', JSON.stringify(userProfile, null, 2))
      
      const token = localStorage.getItem('token')
      const refreshToken = localStorage.getItem('refreshToken')
      
      if (token && userProfile) {
        // 后端返回的字段名是大写开头，需要转换为小写开头
        // 并且需要从token中解析用户ID
        let userId = 0
        try {
          // 解析JWT token获取用户ID
          if (token) {
            const parts = token.split('.')
            if (parts.length === 3 && parts[1]) {
              const tokenPayload = JSON.parse(atob(parts[1]))
              userId = tokenPayload.Id || tokenPayload.id || 0
              console.log('从token解析的用户ID:', userId)
            }
          }
        } catch (error) {
          console.error('解析token失败:', error)
        }
        
        const userData = {
          id: userId,
          email: userProfile.Email || userProfile.email || '',
          nickname: userProfile.Nickname || userProfile.nickname || '',
          phone: userProfile.Phone || userProfile.phone || '',
          aboutMe: userProfile.AboutMe || userProfile.aboutMe || '',
          birthday: userProfile.Birthday || userProfile.birthday || '',
          avatar: userProfile.Avatar || userProfile.avatar || ''
        }
        console.log('=== 处理后的用户数据 ===')
        console.log('userData:', JSON.stringify(userData, null, 2))
        
        userStore.setUser(userData, token, refreshToken || '')
        console.log('用户状态已设置到store:', userData)
        console.log('=== 用户信息获取完成 ===')
      } else {
        console.warn('token 或 userProfile 为空:')
        console.warn('token:', token)
        console.warn('userProfile:', userProfile)
      }
    } catch (error) {
      console.error('=== 获取用户信息失败 ===')
      console.error('错误详情:', error)
      console.error('错误类型:', typeof error)
      if (error instanceof Error) {
        console.error('错误消息:', error.message)
        console.error('错误堆栈:', error.stack)
      }
      // 即使获取用户信息失败，也不影响登录流程
    }
  }

  // 注册表单
  const signupFormRef = ref<FormInstance>()
  const signupForm = reactive<SignupForm>({
    email: '',
    password: '',
    confirmPassword: ''
  })
  const signupRules: FormRules = {
    email: [
      { required: true, message: '请输入邮箱', trigger: 'blur' },
      { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
    ],
    password: [
      { required: true, message: '请输入密码', trigger: 'blur' },
      { min: 8, message: '密码长度不能小于8位', trigger: 'blur' },
      { 
        pattern: /^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$/,
        message: '密码必须包含字母、数字和特殊字符',
        trigger: 'blur'
      }
    ],
    confirmPassword: [
      { required: true, message: '请确认密码', trigger: 'blur' },
      {
        validator: (rule: any, value: string, callback: any) => {
          if (value !== signupForm.password) {
            callback(new Error('两次输入的密码不一致'))
          } else {
            callback()
          }
        },
        trigger: 'blur'
      }
    ]
  }
  const signupLoading = ref(false)
  const showSignup = ref(false)

  // 切换到注册
  const toggleSignup = () => {
    showSignup.value = !showSignup.value
  }

  // 发送验证码
  const sendVerificationCode = async () => {
    try {
      // 验证手机号
      await phoneFormRef.value?.validateField('phone')
      
      // 调用发送验证码API
      await codeApi.sendCode({
        biz: 'login',
        phone: phoneForm.phone
      })
      
      // 开始倒计时
      codeSent.value = true
      countdown.value = 60
      timer = window.setInterval(() => {
        countdown.value--
        if (countdown.value <= 0) {
          clearInterval(timer!)
          timer = null
          codeSent.value = false
        }
      }, 1000)
      
      ElMessage.success('验证码已发送')
    } catch (error) {
      console.error('发送验证码失败:', error)
      ElMessage.error('验证码发送失败，请稍后重试')
    }
  }

  // 手机号登录
  const handlePhoneLogin = async () => {
    try {
      await phoneFormRef.value?.validate()
      
      phoneLoading.value = true
      
      // 调用手机号登录API（token 会在响应拦截器中自动从响应头获取并保存）
      await userApi.loginByPhone(phoneForm.phone, phoneForm.code)
      
      // 登录成功后处理用户状态
      await handleLoginSuccess()
      
      ElMessage.success('登录成功')
      
      // 跳转到首页或重定向页面
      const redirectUrl = route.query.redirect as string || '/'
      router.push(redirectUrl)
    } catch (error) {
      console.error('手机号登录失败:', error)
      ElMessage.error('登录失败，请检查手机号和验证码是否正确')
    } finally {
      phoneLoading.value = false
    }
  }

  // 邮箱登录
  const handleEmailLogin = async () => {
    try {
      await emailFormRef.value?.validate()
      
      emailLoading.value = true
      
      // 调用邮箱登录API（token 会在响应拦截器中自动从响应头获取并保存）
      await userApi.login({
        email: emailForm.email,
        password: emailForm.password
      })
      
      // 登录成功后处理用户状态
      await handleLoginSuccess()
      
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
      ElMessage.error('登录失败，请检查邮箱和密码是否正确')
    } finally {
      emailLoading.value = false
    }
  }

  // 注册
  const handleSignup = async () => {
    try {
      await signupFormRef.value?.validate()
      
      signupLoading.value = true
      
      // 保存邮箱用于自动填充
      const registeredEmail = signupForm.email
      
      // 调用注册API，只发送邮箱和密码
      await userApi.signup({
        email: signupForm.email,
        password: signupForm.password,
        confirmPassword: signupForm.confirmPassword
      })
      
      ElMessage.success('注册成功，请登录')
      
      // 重置表单并切换到登录
      signupForm.email = ''
      signupForm.password = ''
      signupForm.confirmPassword = ''
      
      showSignup.value = false
      
      // 自动填充邮箱
      emailForm.email = registeredEmail
    } catch (error) {
      console.error('注册失败:', error)
      ElMessage.error('注册失败，请稍后重试')
    } finally {
      signupLoading.value = false
    }
  }

  // 微信登录
  const handleWechatLogin = () => {
    ElMessage.info('微信登录功能暂未实现')
  }

  // QQ登录
  const handleQQLogin = () => {
    ElMessage.info('QQ登录功能暂未实现')
  }

  // 微博登录
  const handleWeiboLogin = () => {
    ElMessage.info('微博登录功能暂未实现')
  }

  onMounted(() => {
    // 检查是否记住了邮箱
    const remembered = localStorage.getItem('remember')
    if (remembered === 'true') {
      rememberMe.value = true
      emailForm.email = localStorage.getItem('email') || ''
    }
  })

  return {
    activeTab,
    emailFormRef,
    emailForm,
    emailRules,
    emailLoading,
    rememberMe,
    phoneFormRef,
    phoneForm,
    phoneRules,
    phoneLoading,
    codeSent,
    countdown,
    signupFormRef,
    signupForm,
    signupRules,
    signupLoading,
    showSignup,
    toggleSignup,
    sendVerificationCode,
    handlePhoneLogin,
    handleEmailLogin,
    handleSignup,
    handleWechatLogin,
    handleQQLogin,
    handleWeiboLogin
  }
}
