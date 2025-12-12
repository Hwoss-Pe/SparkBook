<template>
  <MainLayout>
    <div class="my-profile-container">
      <div class="profile-header">
        <h1 class="page-title">我的</h1>
      </div>
      
      <div class="profile-content">
        <div class="profile-card">
          <div class="profile-info">
            <div class="avatar-section">
              <el-avatar :size="80" :src="userAvatarResolved">
                {{ userInfo.nickname ? userInfo.nickname.substring(0, 1) : 'U' }}
              </el-avatar>
              <el-button type="text" @click="showEditDialog = true" class="edit-btn">
                <el-icon><Edit /></el-icon>
                编辑资料
              </el-button>
            </div>
            
            <div class="info-section">
              <div class="info-item">
                <label>昵称</label>
                <span>{{ userInfo.nickname || '未设置' }}</span>
              </div>
              <div class="info-item">
                <label>邮箱</label>
                <span>{{ userInfo.email || '未设置' }}</span>
              </div>
              <div class="info-item">
                <label>手机号</label>
                <span>{{ userInfo.phone || '未绑定' }}</span>
              </div>
              <div class="info-item">
                <label>个人简介</label>
                <span>{{ userInfo.aboutMe || '这个人很懒，还没有填写个人简介' }}</span>
              </div>
              <div class="info-item">
                <label>生日</label>
                <span>{{ userInfo.birthday || '未设置' }}</span>
              </div>
            </div>
          </div>
          
          <div class="stats-section">
            <div class="stat-item">
              <div class="stat-value">{{ stats.articleCount }}</div>
              <div class="stat-label">发布文章</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ stats.draftCount }}</div>
              <div class="stat-label">草稿</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ stats.followerCount }}</div>
              <div class="stat-label">粉丝</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ stats.followingCount }}</div>
              <div class="stat-label">关注</div>
            </div>
          </div>
        </div>
        
        <div class="action-cards">
          <div class="action-card" @click="navigateTo('/create')">
            <el-icon class="action-icon"><Edit /></el-icon>
            <h3>创作中心</h3>
            <p>发布文章，管理草稿</p>
          </div>
          
          <div class="action-card" @click="navigateTo('/message')">
            <el-icon class="action-icon"><ChatDotRound /></el-icon>
            <h3>我的评论</h3>
            <p>查看我参与的评论与互动</p>
          </div>
          
          <div class="action-card" @click="navigateToMyCollections">
            <el-icon class="action-icon"><Star /></el-icon>
            <h3>我的收藏</h3>
            <p>查看收藏的文章</p>
          </div>
          
          <div class="action-card" @click="navigateToMyProfile">
            <el-icon class="action-icon"><User /></el-icon>
            <h3>个人主页</h3>
            <p>查看自己的个人主页</p>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 编辑资料弹窗 -->
    <el-dialog
      v-model="showEditDialog"
      title="编辑个人资料"
      width="500px"
      @close="resetEditForm"
    >
      <el-form :model="editForm" :rules="editRules" ref="editFormRef" label-width="80px">
        <el-form-item label="头像">
          <div class="avatar-upload">
            <el-upload
              class="avatar-uploader"
              action="#"
              :show-file-list="false"
              :auto-upload="false"
              :on-change="handleAvatarChange"
              accept="image/*"
            >
              <el-avatar :size="60" :src="editFormAvatarResolved">
                {{ editForm.nickname ? editForm.nickname.substring(0, 1) : 'U' }}
              </el-avatar>
              <div class="upload-overlay">
                <el-icon><Camera /></el-icon>
              </div>
            </el-upload>
          </div>
        </el-form-item>
        
        <el-form-item label="昵称" prop="nickname">
          <el-input
            v-model="editForm.nickname"
            placeholder="请输入昵称"
            maxlength="20"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="editForm.email"
            placeholder="请输入邮箱"
            disabled
          />
          <div class="form-tip">邮箱不可修改</div>
        </el-form-item>
        
        <el-form-item label="手机号" prop="phone">
          <el-input
            v-model="editForm.phone"
            placeholder="请输入手机号"
            maxlength="11"
          />
        </el-form-item>
        
        <el-form-item label="个人简介" prop="aboutMe">
          <el-input
            v-model="editForm.aboutMe"
            type="textarea"
            placeholder="介绍一下自己吧"
            maxlength="200"
            show-word-limit
            :rows="4"
          />
        </el-form-item>
        
        <el-form-item label="生日" prop="birthday">
          <el-date-picker
            v-model="editForm.birthday"
            type="date"
            placeholder="选择生日"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditDialog = false">取消</el-button>
          <el-button type="primary" @click="saveProfile" :loading="saving">
            保存
          </el-button>
        </span>
      </template>
    </el-dialog>
  </MainLayout>
</template>

<script setup lang="ts">
import MainLayout from '@/components/layout/MainLayout.vue'
import { Edit, ChatDotRound, Star, User, Camera } from '@element-plus/icons-vue'
import useMyProfileView from '@/scripts/views/MyProfileView'

const {
  userInfo,
  userAvatarResolved,
  stats,
  showEditDialog,
  editForm,
  editFormAvatarResolved,
  editRules,
  editFormRef,
  saving,
  navigateTo,
  navigateToMyProfile,
  navigateToMyCollections,
  handleAvatarChange,
  saveProfile,
  resetEditForm
} = useMyProfileView()
</script>

<style lang="scss">
@import '@/styles/views/MyProfileView.scss';
</style>
