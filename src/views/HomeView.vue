<template>
  <el-container>
    <div class="home">
      <el-row class="banner" justify="center" align="middle">
        <el-col :span="24">
          <h1>欢迎来到博客</h1>
          <p>在这里分享你的想法和故事</p>
          <div class="action-buttons">
            <el-button v-if="isLoggedIn" type="success" size="large" @click="router.push('/add-post')">
              写博客
            </el-button>
            <el-button v-else type="success" size="large" @click="router.push('/login')">
              立即登录
            </el-button>
          </div>
        </el-col>
      </el-row>

      <el-row class="features" :gutter="20">
        <el-col :xs="24" :sm="8">
          <el-card>
            <template #header>
              <div class="card-header">
                <i class="el-icon-edit"></i>
                <span>写作自由</span>
              </div>
            </template>
            <div class="card-content">
              在这里，你可以自由地表达你的想法，分享你的知识和经验。
            </div>
          </el-card>
        </el-col>

        <el-col :xs="24" :sm="8">
          <el-card>
            <template #header>
              <div class="card-header">
                <i class="el-icon-reading"></i>
                <span>阅读发现</span>
              </div>
            </template>
            <div class="card-content">
              探索其他作者的精彩文章，获取新的知识和灵感。
            </div>
          </el-card>
        </el-col>

        <el-col :xs="24" :sm="8">
          <el-card>
            <template #header>
              <div class="card-header">
                <i class="el-icon-user"></i>
                <span>社区互动</span>
              </div>
            </template>
            <div class="card-content">
              与其他用户交流互动，分享观点，建立连接。
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </el-container>
</template>

<script>
import { computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'

export default {
  name: 'HomeView',
  setup() {
    const store = useStore()
    const router = useRouter()
    const isLoggedIn = computed(() => store.state.user.isLoggedIn)

    // 在组件挂载后检查登录状态并跳转
    onMounted(() => {
      if (isLoggedIn.value) {
        router.push('/home')
      }
    })

    return { isLoggedIn, router }
  }
}
</script>

<style scoped>
/* 保持原有的样式不变 */
.home {
  height: calc(100vh - 60px);
  width: 100%;
  background: #3498db;
  padding-top: 40px;
}

.banner {
  text-align: center;
  font-size: 20px;
  background-color: #f5f7fa;
  border-radius: 4px;
  margin: 0 auto;
  max-width: 1200px;
  height: 500px;
}

.features {
  margin-top: 40px;
}

.el-card {
  margin-bottom: 20px;
  height: 200px;
  flex-direction: column;
  text-align: center;
}
</style>

<style>
body, html {
  margin: 0;
  padding: 0;
}
</style>