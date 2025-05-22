<template>
  <div class="edit-post">
    <el-card v-if="blog">
      <template #header>
        <h2>编辑博客</h2>
      </template>

      <el-form :model="postForm" :rules="rules" ref="postFormRef" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="postForm.title" placeholder="请输入博客标题"></el-input>
        </el-form-item>

        <el-form-item label="摘要" prop="summary">
          <el-input
            v-model="postForm.summary"
            type="textarea"
            :rows="3"
            placeholder="请输入博客摘要"
          ></el-input>
        </el-form-item>

        <el-form-item label="内容" prop="content">
          <el-input
            v-model="postForm.content"
            type="textarea"
            :rows="10"
            placeholder="请输入博客内容"
          ></el-input>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            保存
          </el-button>
          <el-button @click="$router.push(`/blog/${blog.id}`)">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    <div v-else class="loading">
      <el-empty description="博客不存在或已被删除"></el-empty>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

export default {
  name: 'EditPost',
  setup() {
    const store = useStore()
    const route = useRoute()
    const router = useRouter()
    const postFormRef = ref(null)
    const loading = ref(false)

    const blog = ref(null)
    const postForm = reactive({
      title: '',
      summary: '',
      content: ''
    })

    const rules = {
      title: [
        { required: true, message: '请输入标题', trigger: 'blur' },
        { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
      ],
      summary: [
        { required: true, message: '请输入摘要', trigger: 'blur' },
        { max: 200, message: '长度不能超过 200 个字符', trigger: 'blur' }
      ],
      content: [
        { required: true, message: '请输入内容', trigger: 'blur' }
      ]
    }

    onMounted(async () => {
      try {
        const response = await store.dispatch('fetchBlogById', route.params.id)
        blog.value = response
        Object.assign(postForm, {
          title: response.title,
          summary: response.summary,
          content: response.content
        })
      } catch (error) {
        ElMessage.error('获取博客详情失败')
      }
    })

    const handleSubmit = () => {
      postFormRef.value.validate(async (valid) => {
        if (valid) {
          loading.value = true
          try {
            await store.dispatch('updateBlog', {
              id: blog.value.id,
              blogData: postForm
            })
            ElMessage.success('更新成功')
            router.push(`/blog/${blog.value.id}`)
          } catch (error) {
            ElMessage.error(error.response?.data?.message || '更新失败')
          } finally {
            loading.value = false
          }
        }
      })
    }

    return {
      blog,
      postForm,
      postFormRef,
      rules,
      loading,
      handleSubmit
    }
  }
}
</script>

<style scoped>
.edit-post {
  max-width: 800px;
  margin: 20px auto;
  padding: 0 20px;
}

.el-button {
  margin-right: 10px;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 400px;
}
</style> 