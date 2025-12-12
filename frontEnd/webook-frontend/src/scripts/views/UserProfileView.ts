import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { userApi } from '@/api/user'
import type { User } from '@/api/user'
import { followApi } from '@/api/follow'
import { articleApi } from '@/api/article'
import type { Article as ApiArticle } from '@/api/article'
import { useUserStore } from '@/stores/user'
import { post, resolveStaticUrl } from '@/api/http'

// 定义类型接口
interface UserProfile {
  id: number;
  nickname: string;
  avatar: string;
  aboutMe: string;
  email?: string;
  birthday?: string;
  articleCount: number;
  followerCount: number;
  followingCount: number;
  isFollowing: boolean;
}

interface Article {
  id: number;
  title: string;
  abstract: string;
  coverImage: string;
  createTime: string;
  readCount: number;
  likeCount: number;
  collectCount: number;
  liked: boolean;
  collected: boolean;
}

interface Collection {
  id: number;
  title: string;
  abstract: string;
  coverImage: string;
  author: {
    name: string;
    avatar: string;
  };
  readCount: number;
  likeCount: number;
  collectCount: number;
  liked: boolean;
  collected: boolean;
}

interface FollowUser {
  id: number;
  name: string;
  avatar: string;
  description: string;
  isFollowing: boolean;
}

interface EditForm {
  nickname: string;
  aboutMe: string;
  avatarUrl: string;
}

interface UploadFile {
  raw: File;
}

export default function useUserProfileView() {
  const router = useRouter()
  const route = useRoute()

  // 当前用户ID（从用户状态中获取）
  const userStore = useUserStore()
  const currentUserId = ref<number>(userStore.user?.id || 0)

  // 用户资料
  const userProfile = ref<UserProfile>({
    id: 0,
    nickname: '',
    avatar: '',
    aboutMe: '',
    articleCount: 0,
    followerCount: 0,
    followingCount: 0,
    isFollowing: false
  })

  // 是否为当前用户的个人主页
  const isCurrentUser = computed(() => {
    return userProfile.value.id === currentUserId.value
  })

  const activeTab = ref<string>((useRoute().query.tab as string) || 'articles')

  // 用户文章列表
  const userArticles = ref<Article[]>([])
  const hasMoreArticles = ref(true)

  // 用户收藏列表
  const userCollections = ref<Collection[]>([])
  const hasMoreCollections = ref(false)

  // 关注/粉丝弹窗
  const showFollowDialog = ref(false)
  const followDialogTitle = ref('')
  const followDialogUsers = ref<FollowUser[]>([])
  const followDialogEmptyText = ref('')

  const onFollowDialogClosed = async () => {
    try {
      const staticResponse = await followApi.getFollowStatics({ followee: userProfile.value.id })
      const followStatic = staticResponse.followStatic
      userProfile.value.followerCount = followStatic.followers
      userProfile.value.followingCount = followStatic.followees
    } catch (error) {
      console.error('刷新关注/粉丝统计失败:', error)
    }
  }

  watch(showFollowDialog, (visible) => {
    if (!visible) {
      onFollowDialogClosed()
    }
  })

  // 编辑资料弹窗
  const showEditDialog = ref(false)
  const editForm = ref<EditForm>({
    nickname: '',
    aboutMe: '',
    avatarUrl: ''
  })

  // 格式化数字，例如1200显示为1.2k
  const formatNumber = (num: number): string => {
    if (num >= 10000) {
      return (num / 10000).toFixed(1) + 'w'
    } else if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'k'
    }
    return num.toString()
  }

  // 格式化日期
  const formatDate = (dateString: string): string => {
    const date = new Date(dateString)
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    })
  }

  // 查看文章详情
  const viewArticle = (id: number) => {
    router.push(`/article/${id}`)
  }

  // 查看用户主页
  const viewUser = (id: number) => {
    if (showFollowDialog.value) {
      showFollowDialog.value = false
    }
    router.push(`/user/${id}`)
  }

  // 关注/取消关注
  const toggleFollow = async () => {
    try {
      if (userProfile.value.isFollowing) {
        // 取消关注
        await followApi.cancelFollow(userProfile.value.id, currentUserId.value)
        userProfile.value.isFollowing = false
        userProfile.value.followerCount--
        ElMessage.success('已取消关注')
      } else {
        // 关注
        await followApi.follow(userProfile.value.id, currentUserId.value)
        userProfile.value.isFollowing = true
        userProfile.value.followerCount++
        ElMessage.success('关注成功')
      }
    } catch (error) {
      console.error('关注操作失败:', error)
      ElMessage.error('操作失败，请稍后重试')
    }
  }

  // 关注/取消关注用户（弹窗中）
  const toggleFollowUser = async (user: FollowUser) => {
    try {
      if (user.isFollowing) {
        // 取消关注
        await followApi.cancelFollow(user.id, currentUserId.value)
        user.isFollowing = false
      } else {
        // 关注
        await followApi.follow(user.id, currentUserId.value)
        user.isFollowing = true
      }
    } catch (error) {
      console.error('关注操作失败:', error)
      ElMessage.error('操作失败，请稍后重试')
    }
  }

  // 显示粉丝列表
  const showFollowers = async () => {
    followDialogTitle.value = '粉丝'
    followDialogEmptyText.value = '暂无粉丝'
    
    try {
      const response = await followApi.getFollower({
        followee: userProfile.value.id,
        offset: 0,
        limit: 20
      })
      followDialogUsers.value = response.follow_relations.map(relation => ({
        id: relation.follower,
        name: relation.name || `用户${relation.follower}`,
        avatar: relation.avatar ? resolveStaticUrl(relation.avatar) : `https://picsum.photos/id/${1000 + relation.follower}/100/100`,
        description: relation.about_me || '这个人很懒，还没有填写个人简介',
        isFollowing: false
      }))
    } catch (error) {
      console.error('获取粉丝列表失败:', error)
      // 使用模拟数据
      followDialogUsers.value = [
        {
          id: 102,
          name: '旅行笔记',
          avatar: 'https://picsum.photos/id/1012/100/100',
          description: '记录世界各地的美景与文化',
          isFollowing: true
        },
        {
          id: 103,
          name: '生活家',
          avatar: 'https://picsum.photos/id/1005/100/100',
          description: '分享高品质生活方式',
          isFollowing: false
        },
        {
          id: 104,
          name: '摄影师小王',
          avatar: 'https://picsum.photos/id/1062/100/100',
          description: '用镜头记录生活的美好瞬间',
          isFollowing: false
        }
      ]
    }
    
    showFollowDialog.value = true
  }

  // 显示关注列表
  const showFollowing = async () => {
    followDialogTitle.value = '关注'
    followDialogEmptyText.value = '暂无关注'
    
    try {
      const response = await followApi.getFollowee({
        follower: userProfile.value.id,
        offset: 0,
        limit: 20
      })
      followDialogUsers.value = response.follow_relations.map(relation => ({
        id: relation.followee,
        name: relation.name || `用户${relation.followee}`,
        avatar: relation.avatar ? resolveStaticUrl(relation.avatar) : `https://picsum.photos/id/${1000 + relation.followee}/100/100`,
        description: relation.about_me || '这个人很懒，还没有填写个人简介',
        isFollowing: true
      }))
    } catch (error) {
      console.error('获取关注列表失败:', error)
      // 使用模拟数据
      followDialogUsers.value = [
        {
          id: 105,
          name: '园艺爱好者',
          avatar: 'https://picsum.photos/id/1074/100/100',
          description: '分享家庭园艺技巧和心得',
          isFollowing: true
        },
        {
          id: 106,
          name: '健身教练',
          avatar: 'https://picsum.photos/id/1027/100/100',
          description: '专业健身指导，科学减脂增肌',
          isFollowing: true
        }
      ]
    }
    
    showFollowDialog.value = true
  }

  // 编辑个人资料
  const editProfile = () => {
    editForm.value.nickname = userProfile.value.nickname
    editForm.value.aboutMe = userProfile.value.aboutMe
    editForm.value.avatarUrl = userProfile.value.avatar
    
    showEditDialog.value = true
  }

  // 处理头像变更
  const handleAvatarChange = async (file: UploadFile) => {
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
    try {
      const form = new FormData()
      form.append('file', file.raw!)
      const url = await post<string>('/upload/avatar', form)
      editForm.value.avatarUrl = url
      ElMessage.success('头像上传成功')
    } catch (e) {
      console.error('头像上传失败:', e)
      ElMessage.error('头像上传失败，请重试')
    }
  }

  // 保存个人资料
  const saveProfile = async () => {
    try {
      // 调用API保存个人资料
      await userApi.updateProfile({
        id: userProfile.value.id,
        nickname: editForm.value.nickname,
        aboutMe: editForm.value.aboutMe,
        avatar: editForm.value.avatarUrl
      })
      
      userProfile.value.nickname = editForm.value.nickname
      userProfile.value.aboutMe = editForm.value.aboutMe
      userProfile.value.avatar = resolveStaticUrl(editForm.value.avatarUrl)
      
      showEditDialog.value = false
      ElMessage.success('个人资料已更新')
    } catch (error) {
      console.error('更新个人资料失败:', error)
      ElMessage.error('更新失败，请稍后重试')
    }
  }

  // 加载更多文章
  const loadMoreArticles = async () => {
    try {
      const pubs = await articleApi.getAuthorPublishedList(userProfile.value.id, { offset: userArticles.value.length, limit: 10 })
      const moreArticles: Article[] = pubs.map(article => ({
        id: article.id,
        title: article.title,
        abstract: article.abstract,
        coverImage: article.coverImage ? resolveStaticUrl(article.coverImage) : `https://picsum.photos/id/${400 + article.id}/400/300`,
        createTime: article.ctime,
        readCount: article.readCnt || 0,
        likeCount: article.likeCnt || 0,
        collectCount: article.collectCnt || 0,
        liked: !!article.liked,
        collected: !!article.collected
      }))

      userArticles.value = [...userArticles.value, ...moreArticles]

      // 如果返回的文章数量小于请求的数量，说明没有更多文章了
      if (pubs.length < 10) {
        hasMoreArticles.value = false
      }
    } catch (error) {
      console.error('加载更多文章失败:', error)
      // 使用模拟数据
      const moreArticles: Article[] = [
        {
          id: 4,
          title: '数字游民：如何边旅行边工作',
          abstract: '远程工作正在改变我们的生活和工作方式，本文分享如何成为一名成功的数字游民...',
          coverImage: 'https://picsum.photos/id/1081/400/300',
          createTime: '2024-11-10T10:30:00',
          readCount: 15300,
          likeCount: 4200,
          collectCount: 900
        },
        {
          id: 5,
          title: '家庭花园：从零开始的种植指南',
          abstract: '无论你是有着一片后院，还是只有一个小阳台，都可以打造自己的绿色天地...',
          coverImage: 'https://picsum.photos/id/145/400/300',
          createTime: '2024-11-05T15:45:00',
          readCount: 7600,
          likeCount: 1800,
          collectCount: 350
        }
      ]
      
      userArticles.value = [...userArticles.value, ...moreArticles]
      hasMoreArticles.value = false
    }
  }

  // 加载更多收藏
  const loadMoreCollections = async () => {
    try {
      const params = {
        offset: userCollections.value.length,
        limit: 10
      }
      
      const response = await articleApi.getCollectedList(params)
      
      if (response && response.length > 0) {
        // 转换 ArticlePub 格式到 Collection 格式
        const collections = response.map(article => ({
          id: article.id,
          title: article.title,
          abstract: article.abstract,
          coverImage: resolveStaticUrl(article.coverImage),
          author: {
            name: article.author.name,
            avatar: resolveStaticUrl(article.author.avatar)
          },
          readCount: article.readCnt,
          likeCount: article.likeCnt,
          collectCount: article.collectCnt,
          liked: !!article.liked,
          collected: article.collected ?? true
        }))
        
        userCollections.value = [...userCollections.value, ...collections]
        hasMoreCollections.value = response.length === 10 // 如果返回数量等于请求数量，说明可能还有更多
      } else {
        hasMoreCollections.value = false
      }
    } catch (error) {
      console.error('加载收藏列表失败:', error)
      hasMoreCollections.value = false
    }
  }

  // 获取用户资料
  const fetchUserProfile = async (userId: number) => {
    try {
      // 公开个人资料
      const profile = await userApi.getPublicProfile(userId)
      
      // 获取关注统计
      const staticResponse = await followApi.getFollowStatics({ followee: userId })
      const followStatic = staticResponse.followStatic
      
      // 获取关注状态
      let isFollowing = false
      if (userId !== currentUserId.value) {
        try {
          const followInfo = await followApi.getFollowInfo({
            follower: currentUserId.value,
            followee: userId
          })
          isFollowing = !!followInfo.follow_relation
        } catch (error) {
          console.error('获取关注状态失败:', error)
        }
      }
      
      // 作者文章统计（已发布数量）
      let publishedCount = 0
      try {
        const stats = await articleApi.getAuthorStats(userId)
        publishedCount = stats.publishedCount || 0
      } catch {}

      const av = (profile.Avatar || profile.avatar || '').replace(/[`]/g, '').trim()
      userProfile.value = {
        id: userId,
        nickname: profile.Nickname || profile.nickname || '',
        avatar: resolveStaticUrl(av),
        aboutMe: profile.AboutMe || profile.aboutMe || '',
        email: profile.Email || profile.email || '',
        birthday: profile.Birthday || profile.birthday || '',
        articleCount: publishedCount,
        followerCount: followStatic.followers,
        followingCount: followStatic.followees,
        isFollowing
      }
    } catch (error) {
      console.error('获取用户资料失败:', error)
      // 使用模拟数据
      if (userId === 101) {
        userProfile.value = {
          id: 101,
          nickname: '美食达人',
          avatar: 'https://picsum.photos/id/1027/100/100',
          aboutMe: '探索美食世界的专业吃货，分享美食制作技巧和餐厅推荐',
          articleCount: 45,
          followerCount: 12800,
          followingCount: 56,
          isFollowing: false
        }
      } else {
        userProfile.value = {
          id: userId,
          nickname: '用户' + userId,
          avatar: `https://picsum.photos/id/${1000 + userId}/100/100`,
          aboutMe: '这个人很懒，还没有填写个人简介',
          articleCount: Math.floor(Math.random() * 50),
          followerCount: Math.floor(Math.random() * 10000),
          followingCount: Math.floor(Math.random() * 100),
          isFollowing: Math.random() > 0.5
        }
      }
    }
  }

  // 获取用户文章
  const fetchUserArticles = async (userId: number) => {
    try {
      const pubs = await articleApi.getAuthorPublishedList(userId, { offset: 0, limit: 10 })
      userArticles.value = pubs.map(article => ({
        id: article.id,
        title: article.title,
        abstract: article.abstract,
        coverImage: article.coverImage ? resolveStaticUrl(article.coverImage) : `https://picsum.photos/id/${400 + article.id}/400/300`,
        createTime: article.ctime,
        readCount: article.readCnt || 0,
        likeCount: article.likeCnt || 0,
        collectCount: article.collectCnt || 0,
        liked: !!article.liked,
        collected: !!article.collected
      }))
      if (pubs.length < 10) {
        hasMoreArticles.value = false
      }
    } catch (error) {
      console.error('获取用户文章失败:', error)
      // 使用模拟数据
      userArticles.value = [
        {
          id: 1,
          title: '如何在家制作完美的提拉米苏',
          abstract: '提拉米苏是一道经典的意大利甜点，本文将分享专业大厨的独家秘方...',
          coverImage: 'https://picsum.photos/id/431/400/300',
          createTime: '2024-11-28T14:20:00',
          readCount: 12500,
          likeCount: 3200,
          collectCount: 640
        },
        {
          id: 2,
          title: '2025年最值得去的10个小众旅行地',
          abstract: '厌倦了人山人海的热门景点？这些鲜为人知的目的地将带给你全新的旅行体验...',
          coverImage: 'https://picsum.photos/id/1036/400/300',
          createTime: '2024-11-25T09:15:00',
          readCount: 18700,
          likeCount: 5400,
          collectCount: 980
        },
        {
          id: 3,
          title: '极简主义：如何通过断舍离改变你的生活',
          abstract: '极简主义不仅是一种生活方式，更是一种思维模式。本文将分享如何开始你的极简之旅...',
          coverImage: 'https://picsum.photos/id/106/400/300',
          createTime: '2024-11-20T16:40:00',
          readCount: 9800,
          likeCount: 2100,
          collectCount: 310
        }
      ]
    }
  }

  // 获取用户收藏
  const fetchUserCollections = async (userId: number) => {
    try {
      const params = {
        offset: 0,
        limit: 10
      }
      
      const response = await articleApi.getCollectedList(params)
      
      if (response && response.length > 0) {
        // 转换 ArticlePub 格式到 Collection 格式
        const collections = response.map(article => ({
          id: article.id,
          title: article.title,
          abstract: article.abstract,
          coverImage: article.coverImage,
          author: {
            name: article.author.name,
            avatar: article.author.avatar
          },
          readCount: article.readCnt,
          likeCount: article.likeCnt,
          collectCount: article.collectCnt,
          liked: !!article.liked,
          collected: article.collected ?? true
        }))
        
        userCollections.value = collections
        hasMoreCollections.value = response.length === 10 // 如果返回数量等于请求数量，说明可能还有更多
      } else {
        userCollections.value = []
        hasMoreCollections.value = false
      }
    } catch (error) {
      console.error('获取收藏列表失败:', error)
      userCollections.value = []
      hasMoreCollections.value = false
    }
  }

  // 封面加载失败占位
  const onCoverError = (art: { coverImage: string }) => {
    art.coverImage = ''
  }

  // 取消收藏并从列表移除
  const cancelCollection = async (id: number) => {
    try {
      await ElMessageBox.confirm('是否取消收藏？', '提示', {
        confirmButtonText: '取消收藏',
        cancelButtonText: '再想想',
        type: 'warning'
      })

      const defaultCid = 1
      await articleApi.cancelCollect(id, defaultCid)
      userCollections.value = userCollections.value.filter(a => a.id !== id)
      ElMessage.success('已取消收藏')
    } catch (error) {
      // 用户点击取消会抛异常，直接忽略
      if (error) {
        const msg = (error as any)?.message || ''
        if (!/cancel/i.test(msg)) {
          console.error('取消收藏失败:', error)
          ElMessage.error('操作失败，请稍后重试')
        }
      }
    }
  }

  onMounted(() => {
    const userId = parseInt(route.params.id as string)
    // 初始化标签（支持通过 query.tab 直接展示收藏等）
    if (route.query.tab && typeof route.query.tab === 'string') {
      activeTab.value = route.query.tab
    }
    fetchUserProfile(userId)
    fetchUserArticles(userId)
    
    if (userId === currentUserId.value) {
      fetchUserCollections(userId)
    }
  })

  watch(() => route.params.id, (newId) => {
    const uid = parseInt(newId as string)
    if (showFollowDialog.value) {
      showFollowDialog.value = false
    }
    fetchUserProfile(uid)
    fetchUserArticles(uid)
    if (uid === currentUserId.value) {
      fetchUserCollections(uid)
    } else {
      userCollections.value = []
      hasMoreCollections.value = false
    }
  })

  watch(() => route.query.tab, (newTab) => {
    if (typeof newTab === 'string') {
      activeTab.value = newTab
    }
  })

  return {
    userProfile,
    isCurrentUser,
    activeTab,
    userArticles,
    hasMoreArticles,
    userCollections,
    hasMoreCollections,
    showFollowDialog,
    followDialogTitle,
    followDialogUsers,
    followDialogEmptyText,
    onFollowDialogClosed,
    showEditDialog,
    editForm,
    currentUserId,
    formatNumber,
    formatDate,
    viewArticle,
    viewUser,
    toggleFollow,
    toggleFollowUser,
    showFollowers,
    showFollowing,
    editProfile,
    handleAvatarChange,
    saveProfile,
    loadMoreArticles,
    loadMoreCollections,
    onCoverError,
    cancelCollection
  }
}
