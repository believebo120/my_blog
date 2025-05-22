<template>
  <el-container>
    <el-main>
      <div class="welcome">
        <h1>欢迎回来，{{ userInfo?.username || '用户' }}</h1>
        <p>继续你的创作之旅</p>
        <el-button type="primary" @click="router.push('/add-post')">
          写新文章
        </el-button>
      </div>

      <!-- 登录后可见的内容区域 -->
      <el-row class="features" :gutter="20">
        <!-- 卡片内容保持不变 -->
      </el-row>
    </el-main>
  </el-container>
</template>

<script>
import { computed } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'

export default {
  name: 'HomeLoggedIn',
  setup() {
    const store = useStore()
    const router = useRouter()

    const userInfo = computed(() => store.state.user.userInfo)
    const activeMenu = computed(() => '1') // 当前激活的菜单

    return { userInfo, router, activeMenu }
  }
}
</script>

<style scoped>
.header {
  background-color: #3498db;
  color: white;
}

.welcome {
  text-align: center;
  padding: 40px;
}
</style>