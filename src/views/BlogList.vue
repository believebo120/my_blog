<template>
  <div class="blog-list">
    <div class="blog-header">
      <h2>博客列表</h2>
      <el-button v-if="isLoggedIn" type="primary" @click="router.push('/add-post')">
        写博客
      </el-button>
    </div>

    <el-row :gutter="20">
      <el-col v-for="blog in blogs" :key="blog.id" :xs="24" :sm="12" :md="8" :lg="6">
        <el-card class="blog-card" shadow="hover">
          <template #header>
            <div class="blog-title">
              <router-link :to="'/blog/' + blog.id">{{ blog.title }}</router-link>
            </div>
          </template>
          <div class="blog-content">
            <p>{{ blog.summary }}</p>
          </div>
          <div class="blog-footer">
            <span>作者: {{ blog.author }}</span>
            <span>{{ formatDate(blog.created_at) }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <div v-if="blogs.length === 0" class="no-data">
      暂无博客内容
    </div>
  </div>
</template>

<script>
import { onMounted, computed } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { formatDate } from '../utils/date.js'

export default {
  name: 'BlogList',
  setup() {
    const store = useStore()

    const blogs = computed(() => store.state.blogs)
    const isLoggedIn = computed(() => store.state.user.isLoggedIn)

    onMounted(async () => {
      try {
        await store.dispatch('fetchBlogs')
      } catch (error) {
        console.error('获取博客列表失败:', error)
      }
    })

    return {
      blogs,
      isLoggedIn,
      formatDate
    }
  }
}
</script>

<style scoped>
.blog-list {
  padding: 20px;
}

.blog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.blog-card {
  margin-bottom: 20px;
  height: 100%;
}

.blog-title {
  font-size: 18px;
  font-weight: bold;
}

.blog-title a {
  color: #303133;
  text-decoration: none;
}

.blog-title a:hover {
  color: #409EFF;
}

.blog-content {
  margin: 10px 0;
  color: #606266;
}

.blog-footer {
  display: flex;
  justify-content: space-between;
  color: #909399;
  font-size: 14px;
}

.no-data {
  text-align: center;
  color: #909399;
  margin-top: 100px;
  font-size: 16px;
}
</style> 