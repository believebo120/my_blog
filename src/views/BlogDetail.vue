<template>
  <div class="blog-detail" v-if="blog">
    <el-card>
      <template #header>
        <div class="blog-header">
          <h1>{{ blog.title }}</h1>
          <div class="blog-meta">
            <span>作者: {{ blog.author }}</span>
            <span>发布时间: {{ formatDate(blog.created_at) }}</span>
          </div>
        </div>
      </template>
      
      <div class="blog-content" v-html="blog.content"></div>

      <div class="blog-actions" v-if="isAuthor">
        <el-button type="primary" @click="handleEdit">编辑</el-button>
        <el-button type="danger" @click="handleDelete">删除</el-button>
      </div>
    </el-card>
  </div>
  <div v-else class="loading">
    <el-empty description="博客不存在或已被删除"></el-empty>
  </div>
</template>

<script>
import { computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { formatDate } from '../utils/date'

export default {
  name: 'BlogDetail',
  setup() {
    const store = useStore()
    const route = useRoute()
    const router = useRouter()

    const blog = computed(() => store.state.currentBlog)
    const userInfo = computed(() => store.state.user.userInfo)
    const isAuthor = computed(() => blog.value?.author_id === userInfo.value?.id)

    onMounted(async () => {
      try {
        await store.dispatch('fetchBlogById', route.params.id)
      } catch (error) {
        ElMessage.error('获取博客详情失败')
      }
    })

    const handleEdit = () => {
      router.push(`/edit-post/${blog.value.id}`)
    }

    const handleDelete = () => {
      ElMessageBox.confirm('确定要删除这篇博客吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          await store.dispatch('deleteBlog', blog.value.id)
          ElMessage.success('删除成功')
          router.push('/blog')
        } catch (error) {
          ElMessage.error('删除失败')
        }
      }).catch(() => {})
    }

    return {
      blog,
      isAuthor,
      formatDate,
      handleEdit,
      handleDelete
    }
  }
}
</script>

<style scoped>
.blog-detail {
  max-width: 800px;
  margin: 20px auto;
  padding: 0 20px;
}

.blog-header {
  text-align: center;
}

.blog-meta {
  color: #909399;
  margin: 10px 0;
}

.blog-meta span {
  margin: 0 10px;
}

.blog-content {
  line-height: 1.8;
  margin: 20px 0;
}

.blog-actions {
  margin-top: 20px;
  text-align: right;
}

.blog-actions .el-button {
  margin-left: 10px;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 400px;
}
</style> 