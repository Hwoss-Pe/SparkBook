<template>
  <MainLayout>
    <div class="hot-container">
      <h1 class="page-title">热门榜单</h1>
      
      <div class="hot-content">
        <el-tabs v-model="activeTab" class="hot-tabs">
          <el-tab-pane label="综合榜" name="overall">
            <div class="ranking-list">
              <div 
                v-for="(item, index) in hotRankings" 
                :key="item.id" 
                class="ranking-item"
                @click="viewArticle(item.id)"
              >
                <div class="ranking-number" :class="{ 'top-three': index < 3 }">{{ index + 1 }}</div>
                <div class="ranking-content">
                  <h3 class="ranking-title">{{ item.title }}</h3>
                  <div class="ranking-abstract">{{ item.abstract }}</div>
                  <div class="ranking-meta">
                    <div class="author-info">
                      <el-avatar :size="24" :src="item.author.avatar">
                        {{ item.author.name.substring(0, 1) }}
                      </el-avatar>
                      <span class="author-name">{{ item.author.name }}</span>
                    </div>
                    <div class="interaction-info">
                      <span class="interaction-item">
                        <el-icon><View /></el-icon>
                        {{ formatNumber(item.readCount) }}
                      </span>
                      <span class="interaction-item">
                        <el-icon><ThumbsUp /></el-icon>
                        {{ formatNumber(item.likeCount) }}
                      </span>
                      <span class="interaction-item">
                        <el-icon><ChatDotRound /></el-icon>
                        {{ formatNumber(item.commentCount) }}
                      </span>
                    </div>
                  </div>
                </div>
                <div class="ranking-cover" v-if="item.coverImage">
                  <img :src="item.coverImage" :alt="item.title" />
                </div>
              </div>
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="日榜" name="daily">
            <div class="empty-state">
              <el-empty description="日榜数据正在更新中" />
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="周榜" name="weekly">
            <div class="empty-state">
              <el-empty description="周榜数据正在更新中" />
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="月榜" name="monthly">
            <div class="empty-state">
              <el-empty description="月榜数据正在更新中" />
            </div>
          </el-tab-pane>
        </el-tabs>
        
        <div class="update-info">
          <el-icon><Timer /></el-icon>
          <span>数据更新于: {{ formatDateTime(updateTime) }}</span>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import MainLayout from '@/components/layout/MainLayout.vue'
import { View, ChatDotRound, Timer } from '@element-plus/icons-vue'
import { Star as ThumbsUp } from '@element-plus/icons-vue'

const router = useRouter()
const activeTab = ref('overall')
const updateTime = ref(new Date())

// 模拟热榜数据
const hotRankings = ref([
  {
    id: 101,
    title: '年轻人为什么都在做副业？',
    abstract: '在当今社会，越来越多的年轻人开始尝试副业。这不仅是为了增加收入，更是为了寻找职业保障和自我实现...',
    coverImage: 'https://picsum.photos/id/1/200/200',
    author: {
      name: '经济观察家',
      avatar: 'https://picsum.photos/id/1005/100/100'
    },
    readCount: 125000,
    likeCount: 32000,
    commentCount: 2100
  },
  {
    id: 102,
    title: '这些小众景点比网红打卡地更值得去',
    abstract: '厌倦了人山人海的旅游景点？这篇文章为你推荐一些鲜为人知但风景绝美的小众旅行地...',
    coverImage: 'https://picsum.photos/id/1036/200/200',
    author: {
      name: '旅行达人',
      avatar: 'https://picsum.photos/id/1012/100/100'
    },
    readCount: 98000,
    likeCount: 25000,
    commentCount: 1800
  },
  {
    id: 103,
    title: '如何在30天内养成一个新习惯',
    abstract: '科学研究表明，养成一个新习惯需要21-30天的持续练习。本文将分享一套实用的方法...',
    coverImage: 'https://picsum.photos/id/1029/200/200',
    author: {
      name: '生活改造师',
      avatar: 'https://picsum.photos/id/1027/100/100'
    },
    readCount: 87000,
    likeCount: 21000,
    commentCount: 1500
  },
  {
    id: 104,
    title: '2025年最值得学习的5个技能',
    abstract: '随着AI和自动化技术的发展，未来的就业市场将发生巨大变化。这些技能将帮助你在未来职场中保持竞争力...',
    coverImage: 'https://picsum.photos/id/0/200/200',
    author: {
      name: '职场导师',
      avatar: 'https://picsum.photos/id/1000/100/100'
    },
    readCount: 76000,
    likeCount: 18000,
    commentCount: 1200
  },
  {
    id: 105,
    title: '这样做，让你的朋友圈不再单调',
    abstract: '朋友圈是展示生活和个性的窗口，如何让你的朋友圈更有趣、更有深度？这里有一些创意和技巧...',
    coverImage: 'https://picsum.photos/id/1062/200/200',
    author: {
      name: '社交达人',
      avatar: 'https://picsum.photos/id/1074/100/100'
    },
    readCount: 65000,
    likeCount: 16000,
    commentCount: 980
  },
  {
    id: 106,
    title: '低成本高品质的家居改造指南',
    abstract: '不需要大笔预算，也能让你的家焕然一新。这篇文章分享一些实用的家居改造技巧...',
    coverImage: 'https://picsum.photos/id/106/200/200',
    author: {
      name: '家居设计师',
      avatar: 'https://picsum.photos/id/1025/100/100'
    },
    readCount: 54000,
    likeCount: 13000,
    commentCount: 760
  },
  {
    id: 107,
    title: '如何科学健身：避开这些常见误区',
    abstract: '很多人健身不得法，不仅效果不佳，还可能带来伤害。本文将帮你避开这些常见的健身误区...',
    coverImage: 'https://picsum.photos/id/1070/200/200',
    author: {
      name: '健身教练',
      avatar: 'https://picsum.photos/id/1012/100/100'
    },
    readCount: 43000,
    likeCount: 11000,
    commentCount: 650
  },
  {
    id: 108,
    title: '读懂财报的10个关键指标',
    abstract: '无论是投资还是了解一家公司，财报都是最重要的信息来源。这10个关键指标将帮你快速读懂财报...',
    coverImage: 'https://picsum.photos/id/180/200/200',
    author: {
      name: '投资顾问',
      avatar: 'https://picsum.photos/id/1005/100/100'
    },
    readCount: 32000,
    likeCount: 8000,
    commentCount: 420
  }
])

// 格式化数字，例如1200显示为1.2k
const formatNumber = (num: number): string => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

// 格式化日期时间
const formatDateTime = (date: Date): string => {
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 查看文章详情
const viewArticle = (id: number) => {
  router.push(`/article/${id}`)
}

onMounted(() => {
  // 这里应该调用API获取热榜数据
  // 目前使用模拟数据
})
</script>

<style scoped>
.hot-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.page-title {
  font-size: 24px;
  margin-bottom: 20px;
  color: var(--text-primary);
}

.hot-content {
  background-color: var(--bg-primary);
  border-radius: var(--border-radius-md);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}

.hot-tabs {
  margin-bottom: 20px;
}

.ranking-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.ranking-item {
  display: flex;
  background-color: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  cursor: pointer;
  transition: transform 0.3s;
}

.ranking-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.ranking-number {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
  color: #666;
  border-radius: 50%;
  font-weight: bold;
  font-size: 18px;
  margin-right: 20px;
  flex-shrink: 0;
}

.ranking-number.top-three {
  background-color: var(--primary-color);
  color: white;
}

.ranking-content {
  flex: 1;
}

.ranking-title {
  margin: 0 0 10px;
  font-size: 18px;
  font-weight: 600;
  line-height: 1.4;
}

.ranking-abstract {
  margin: 0 0 15px;
  font-size: 14px;
  color: #666;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.ranking-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.author-name {
  font-size: 14px;
  color: #333;
}

.interaction-info {
  display: flex;
  gap: 12px;
}

.interaction-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #999;
}

.ranking-cover {
  width: 120px;
  height: 120px;
  overflow: hidden;
  border-radius: 8px;
  margin-left: 20px;
  flex-shrink: 0;
}

.ranking-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.update-info {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  color: #999;
  margin-top: 20px;
  justify-content: flex-end;
}

.empty-state {
  text-align: center;
  padding: 40px 0;
}
</style>
