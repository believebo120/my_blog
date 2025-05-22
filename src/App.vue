<template>
  <el-container class="app-container">
    <!-- 导航栏（层级高于背景） -->
    <el-header class="nav-bar">
      <div class="nav-links">
        <router-link to="/" class="nav-item">首页</router-link>
        <router-link to="/blog" class="nav-item">关注</router-link>
        <template v-if="!isLoggedIn">
          <router-link to="/personal" class="nav-item">个人中心</router-link>
          <router-link to="/login" class="nav-item">登录</router-link>
        </template>
        <template v-if="isLoggedIn">
          <router-link to="/profile" class="nav-item">
            {{ userInfo?.username||'用户' }}
          </router-link>
        </template>
      </div>
    </el-header>

    <!-- 背景容器（层级低于导航栏，覆盖整个页面） -->
    <div class="global-background"></div>

    <!-- 内容区域 -->
    <el-main class="content-area">
      <router-view />
    </el-main>
  </el-container>
</template>

<script>
import { computed } from 'vue';
import { useStore } from 'vuex';
import { useRouter } from 'vue-router';

export default {
  name: 'App',
  setup() {
    const store = useStore();
    const router = useRouter();
    const isLoggedIn = computed(() => store.state.user.isLoggedIn);//获取登录状态
    const userInfo=computed(() => store.state.user.userInfo)

    const handleLogout = async () => {
      await store.dispatch('user/logout');
    };

    return {
      isLoggedIn,
      userInfo,
      handleLogout
    };
  }
};
</script>

<style>
/* 全局样式：背景覆盖整个页面，层级在导航栏下方 */
.global-background {
  position: absolute; /* 绝对定位，脱离文档流 */
  top: 0;
  left: 0;
  width: 100%;
  height: 100vh; /* 占满视口高度 */
  background: #3498db; /* 背景颜色或图片 */
  z-index: -1; /* 层级低于导航栏 */
}

/* 导航栏样式：层级高于背景 */
.nav-bar {
  position: relative; /* 保持正常布局，z-index 生效 */
  z-index: 1; /* 层级优先级高于背景 */
  line-height: 60px;
  height: 60px;
  background-color: #ffffff;
  border-bottom: 1px solid #e4e7ed;
  margin: 0 20px;
  border-radius: 10px;
  display: flex;
  align-items: center;
}

.nav-links {
  display: flex;
  justify-content: space-evenly;
  flex-grow: 1;
  list-style: none;
  margin: 0;
  padding: 0;
  width: 100%;
}

.nav-item {
  flex: 1;
  text-align: center;
  text-decoration: none;
  color: #2c3e50;
  transition: color 0.3s;
  font-size: 25px;
}

.nav-item:hover,
.nav-item.router-link-exact-active {
  color: #42b983;
}

/* 内容区域：避免被背景覆盖，添加内边距 */
.content-area {
  padding-top: 60px; /* 预留导航栏高度 */
  min-height: calc(100vh - 60px); /* 内容区域高度适配 */
}
</style>