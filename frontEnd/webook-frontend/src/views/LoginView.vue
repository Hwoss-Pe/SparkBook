<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h1 class="login-logo">小微书</h1>
        <p class="login-subtitle">探索、分享、连接</p>
      </div>
      
      <div v-if="!showSignup">
        <el-tabs v-model="activeTab" class="login-tabs">
          <el-tab-pane label="邮箱登录" name="email">
            <div class="login-form">
              <el-form 
                ref="emailFormRef"
                :model="emailForm"
                :rules="emailRules"
                label-position="top"
              >
                <el-form-item label="邮箱" prop="email">
                  <el-input v-model="emailForm.email" placeholder="请输入邮箱" />
                </el-form-item>
                <el-form-item label="密码" prop="password">
                  <el-input 
                    v-model="emailForm.password" 
                    type="password" 
                    placeholder="请输入密码" 
                    show-password
                  />
                </el-form-item>
                <div class="remember-me">
                  <el-checkbox v-model="rememberMe">记住我</el-checkbox>
                </div>
                <el-form-item>
                  <el-button 
                    type="primary" 
                    class="login-button" 
                    @click="handleEmailLogin" 
                    :loading="emailLoading"
                  >
                    登录
                  </el-button>
                </el-form-item>
              </el-form>
              
              <div class="form-footer">
                <span class="form-footer-text">
                  没有账号？
                  <span class="form-footer-link" @click="toggleSignup">立即注册</span>
                </span>
              </div>
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="手机号登录" name="phone">
            <div class="login-form">
              <el-form 
                ref="phoneFormRef"
                :model="phoneForm"
                :rules="phoneRules"
                label-position="top"
              >
                <el-form-item label="手机号" prop="phone">
                  <el-input v-model="phoneForm.phone" placeholder="请输入手机号" />
                </el-form-item>
                <el-form-item label="验证码" prop="code">
                  <el-input v-model="phoneForm.code" placeholder="请输入验证码">
                    <template #append>
                      <el-button 
                        class="verification-code-button"
                        @click="sendVerificationCode" 
                        :disabled="codeSent"
                      >
                        {{ codeSent ? `${countdown}秒后重新获取` : '获取验证码' }}
                      </el-button>
                    </template>
                  </el-input>
                </el-form-item>
                <el-form-item>
                  <el-button 
                    type="primary" 
                    class="login-button" 
                    @click="handlePhoneLogin" 
                    :loading="phoneLoading"
                  >
                    登录
                  </el-button>
                </el-form-item>
              </el-form>
              
              <div class="form-footer">
                <span class="form-footer-text">
                  没有账号？
                  <span class="form-footer-link" @click="toggleSignup">立即注册</span>
                </span>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>
        
        <div class="third-party-login">
          <div class="third-party-title">其他登录方式</div>
          <div class="third-party-buttons">
            <div class="third-party-button wechat-button" @click="handleWechatLogin">
              <el-icon><i-ep-chat-dot-round /></el-icon>
            </div>
            <div class="third-party-button qq-button" @click="handleQQLogin">
              Q
            </div>
            <div class="third-party-button weibo-button" @click="handleWeiboLogin">
              W
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="login-form">
        <el-form 
          ref="signupFormRef"
          :model="signupForm"
          :rules="signupRules"
          label-position="top"
        >
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="signupForm.email" placeholder="请输入邮箱" />
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input 
              v-model="signupForm.password" 
              type="password" 
              placeholder="请输入密码" 
              show-password
            />
          </el-form-item>
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input 
              v-model="signupForm.confirmPassword" 
              type="password" 
              placeholder="请确认密码" 
              show-password
            />
          </el-form-item>
          <el-form-item>
            <el-button 
              type="primary" 
              class="login-button" 
              @click="handleSignup" 
              :loading="signupLoading"
            >
              注册
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="form-footer">
          <span class="form-footer-text">
            已有账号？
            <span class="form-footer-link" @click="toggleSignup">立即登录</span>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ChatDotRound as IChatDotRound } from '@element-plus/icons-vue'
import useLoginView from '@/scripts/views/LoginView'

const {
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
} = useLoginView()
</script>

<style lang="scss">
@import '@/styles/views/LoginView.scss';
</style>