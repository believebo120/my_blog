<template>
  <div class="add-post">
    <el-card>
      <template #header>
        <h2>写博客</h2>
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
            发布
          </el-button>
          <el-button @click="$router.push('/blog')">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

export default {
  name: 'AddPost',
  setup() {
    const store = useStore()
    const router = useRouter()
    const postFormRef = ref(null)
    const loading = ref(false)

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

    const handleSubmit = () => {
      postFormRef.value.validate(async (valid) => {
        if (valid) {
          loading.value = true
          try {
            await store.dispatch('createBlog', postForm)
            ElMessage.success('发布成功')
            router.push('/blog')
          } catch (error) {
            ElMessage.error(error.response?.data?.message || '发布失败')
          } finally {
            loading.value = false
          }
        }
      })
    }

    return {
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
.add-post {
  max-width: 800px;
  margin: 20px auto;
  padding: 0 20px;
}

.el-button {
  margin-right: 10px;
}
</style> 