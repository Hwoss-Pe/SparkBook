<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="app">
    <header class="app-header">
      <h1>小微书</h1>
      <nav>
        <router-link to="/">首页</router-link> |
        <router-link to="/hot">热门</router-link> |
        <router-link to="/follow">关注</router-link> |
        <router-link to="/search">搜索</router-link> |
        <router-link to="/login">登录</router-link>
      </nav>
    </header>
    <main>
      <RouterView />
    </main>
  </div>
</template>

<style>
/* 全局样式 */
body {
  margin: 0;
  padding: 0;
  font-family: 'PingFang SC', 'Helvetica Neue', Helvetica, 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background-color: #f5f5f5;
  color: #333;
}

* {
  box-sizing: border-box;
}

a {
  text-decoration: none;
  color: inherit;
}

img {
  max-width: 100%;
}

.app {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.app-header {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
  text-align: center;
}

.app-header h1 {
  color: #ff2442;
  margin: 0 0 15px 0;
}

nav {
  padding: 10px 0;
}

nav a {
  margin: 0 10px;
  color: #333;
  font-weight: 500;
}

nav a:hover {
  color: #ff2442;
}

nav a.router-link-active {
  color: #ff2442;
  font-weight: bold;
}
</style>
