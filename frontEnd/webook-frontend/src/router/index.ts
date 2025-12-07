import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/follow',
      name: 'follow',
      component: () => import('../views/FollowView.vue'),
      meta: { title: '关注 - 小微书', requiresAuth: true }
    },
    {
      path: '/hot',
      name: 'hot',
      component: () => import('../views/HotView.vue'),
      meta: { title: '热榜 - 小微书' }
    },
    {
      path: '/message',
      name: 'message',
      component: () => import('../views/MessageView.vue'),
      meta: { title: '消息 - 小微书', requiresAuth: true }
    },
    {
      path: '/article/:id',
      name: 'article',
      component: () => import('../views/ArticleDetailView.vue'),
      meta: { title: '文章详情 - 小微书' }
    },
    {
      path: '/search',
      name: 'search',
      component: () => import('../views/SearchView.vue'),
      meta: { title: '搜索 - 小微书' }
    },
    {
      path: '/user/:id',
      name: 'user',
      component: () => import('../views/UserProfileView.vue'),
      meta: { title: '用户主页 - 小微书' }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: { title: '登录/注册 - 小微书' }
    },
    {
      path: '/create',
      name: 'create',
      component: () => import('../views/CreateArticleView.vue'),
      meta: { title: '创作中心 - 小微书', requiresAuth: true }
    },
    // 关于页面已删除
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('../views/HomeView.vue'),
    }
  ],
})

// 全局前置守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = to.meta.title as string
  }
  
  // 检查是否需要登录
  if (to.meta.requiresAuth) {
    // 从 localStorage 中检查是否有 token
    const token = localStorage.getItem('token')
    const isLoggedIn = !!token
    
    if (!isLoggedIn) {
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
      return
    }
  }
  
  next()
})

export default router
