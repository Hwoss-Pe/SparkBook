import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, type FormInstance, type UploadFile } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { userApi, type User } from '@/api/user'
import { articleApi } from '@/api/article'

export default function useMyProfileView() {
  const router = useRouter()
  const userStore = useUserStore()
  
  // 响应式数据
  const userInfo = ref<User>({
    id: 0,
    email: '',
    nickname: '',
    phone: '',
    aboutMe: '',
    birthday: '',
    avatar: ''
  })
  
  const stats = ref({
    articleCount: 0,
    draftCount: 0,
    followerCount: 0,
    followingCount: 0
  })
  
  const showEditDialog = ref(false)
  const saving = ref(false)
  const editFormRef = ref<FormInstance>()
  
  // 编辑表单
  const editForm = reactive({
    nickname: '',
    email: '',
    phone: '',
    aboutMe: '',
    birthday: '',
    avatar: ''
  })
  
  // 表单验证规则
  const editRules = {
    nickname: [
      { required: true, message: '请输入昵称', trigger: 'blur' },
      { min: 2, max: 20, message: '昵称长度在 2 到 20 个字符', trigger: 'blur' }
    ],
    phone: [
      { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
    ]
  }
  
  // 页面导航
  const navigateTo = (path: string) => {
    router.push(path)
  }

  // 跳转到自己的个人主页
  const navigateToMyProfile = () => {
    console.log("跳转个人主页")
    let uid = userStore.user?.id || userInfo.value.id
    if (!uid) {
      userStore.initUserState()
      uid = userStore.user?.id || 0
    }
    if (uid) {
      router.push(`/user/${uid}`)
    } else {
      router.push({ path: '/login', query: { redirect: '/my' } })
    }
  }

  // 跳转到个人主页的“收藏”标签并直接展示收藏
  const navigateToMyCollections = () => {
    let uid = userStore.user?.id || userInfo.value.id
    if (!uid) {
      userStore.initUserState()
      uid = userStore.user?.id || 0
    }
    if (uid) {
      router.push({ path: `/user/${uid}`, query: { tab: 'collections' } })
    } else {
      router.push({ path: '/login', query: { redirect: '/my' } })
    }
  }
  
  // 处理头像上传
  const handleAvatarChange = (file: UploadFile) => {
    const isImage = file.raw?.type?.startsWith('image/')
    if (!isImage) {
      ElMessage.error('只能上传图片文件!')
      return
    }
    
    const isLt2M = file.raw!.size / 1024 / 1024 < 2
    if (!isLt2M) {
      ElMessage.error('图片大小不能超过 2MB!')
      return
    }
    
    // 创建预览URL
    const reader = new FileReader()
    reader.onload = (e) => {
      editForm.avatar = e.target?.result as string
    }
    reader.readAsDataURL(file.raw!)
  }
  
  // 重置编辑表单
  const resetEditForm = () => {
    Object.assign(editForm, {
      nickname: userInfo.value.nickname || '',
      email: userInfo.value.email || '',
      phone: userInfo.value.phone || '',
      aboutMe: userInfo.value.aboutMe || '',
      birthday: userInfo.value.birthday || '',
      avatar: userInfo.value.avatar || ''
    })
    editFormRef.value?.clearValidate()
  }
  
  // 保存个人资料
  const saveProfile = async () => {
    if (!editFormRef.value) return
    
    try {
      await editFormRef.value.validate()
      saving.value = true
      
      const updateData = {
        nickname: editForm.nickname,
        phone: editForm.phone,
        aboutMe: editForm.aboutMe,
        birthday: editForm.birthday
      }
      
      await userApi.updateProfile(updateData)
      
      // 更新本地数据
      Object.assign(userInfo.value, updateData)
      
      // 更新用户store
      if (userStore.user) {
        Object.assign(userStore.user, updateData)
        localStorage.setItem('user', JSON.stringify(userStore.user))
      }
      
      ElMessage.success('个人资料更新成功')
      showEditDialog.value = false
      
    } catch (error) {
      console.error('更新个人资料失败:', error)
      ElMessage.error('更新失败，请重试')
    } finally {
      saving.value = false
    }
  }
  
  // 获取用户信息
  const fetchUserInfo = async () => {
    try {
      userStore.initUserState()
      if (userStore.user) {
        userInfo.value = { ...userStore.user }
      } else {
        const response = await userApi.getProfile()
        userInfo.value = response.user
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
      ElMessage.error('获取用户信息失败')
    }
  }
  
  const fetchStats = async () => {
    try {
      userStore.initUserState()
      const uid = userStore.user?.id || userInfo.value.id || 0
      const s = await articleApi.getAuthorStats(uid)
      stats.value = {
        articleCount: s.publishedCount || 0,
        draftCount: s.draftCount || 0,
        followerCount: s.followerCount || 0,
        followingCount: s.followingCount || 0
      }
    } catch (error) {
      console.error('获取统计数据失败:', error)
    }
  }
  
  // 初始化
  onMounted(() => {
    fetchUserInfo()
    fetchStats()
  })
  
  return {
    userInfo,
    stats,
    showEditDialog,
    editForm,
    editRules,
    editFormRef,
    saving,
    navigateTo,
    navigateToMyProfile,
    navigateToMyCollections,
    handleAvatarChange,
    saveProfile,
    resetEditForm
  }
}
