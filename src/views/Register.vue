<template>
  <div class="register-container">
    <el-card class="register-card">
      <template #header>
        <h2>注册</h2>
      </template>
      <el-form :model="registerForm" :rules="rules" ref="registerFormRef" label-width="80px">
        <el-form-item prop="username" label="用户名">
          <el-input v-model="registerForm.username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item prop="email" label="邮箱">
          <el-input v-model="registerForm.email" placeholder="请输入邮箱"></el-input>
        </el-form-item>
        <el-form-item prop="password" label="密码">
          <el-input v-model="registerForm.password" type="password" placeholder="请输入密码"></el-input>
        </el-form-item>
        <el-form-item prop="confirmPassword" label="确认密码">
          <el-input v-model="registerForm.confirmPassword" type="password" placeholder="请再次输入密码"></el-input>
        </el-form-item>

        <!-- 头像上传组件 -->
        <el-form-item label="头像">
          <div class="avatar-upload-wrapper">
            <el-upload
                class="square-uploader"
                action="#"
                :show-file-list="false"
                :on-success="handleAvatarSuccess"
                :before-upload="beforeAvatarUpload"
                :http-request="customUpload">
              <div class="upload-area">
                <img v-if="imageUrl" :src="imageUrl" class="avatar-preview" alt="用户头像" />
                <div v-else class="plus-icon">+</div>
              </div>
            </el-upload>
            <div class="upload-tip">支持 JPG、PNG 格式，最大 2MB</div>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleRegister" :loading="loading">注册</el-button>
          <router-link to="/login">
            <el-button>返回登录</el-button>
          </router-link>
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
  name: 'Register',
  setup() {
    const store = useStore()
    const router = useRouter()
    const registerFormRef = ref(null)
    const loading = ref(false)
    const imageUrl = ref('')  // 头像预览 URL
    const file = ref(null)    // 存储选中的文件

    const registerForm = reactive({
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
      avatar: ''  // 存储头像数据
    })

    // 邮箱验证函数
    const validateEmail = (rule, value, callback) => {
      const emailRegex = /^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$/
      if (!value) {
        callback(new Error('请输入邮箱'))
      } else if (!emailRegex.test(value)) {
        callback(new Error('请输入有效的邮箱地址'))
      } else {
        callback()
      }
    }

    const validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'))
      } else if (value !== registerForm.password) {
        callback(new Error('两次输入密码不一致!'))
      } else {
        callback()
      }
    }

    const rules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
      ],
      email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { validator: validateEmail, trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, message: '请再次输入密码', trigger: 'blur' },
        { validator: validatePass2, trigger: 'blur' }
      ]
    }

    // 头像上传成功处理
    const handleAvatarSuccess = (response) => {
      if (response && response.url) {
        imageUrl.value = response.url
        registerForm.avatar = response.url  // 存储服务器返回的 URL
      }
    }

    // 自定义上传函数（将文件转为 Base64）
    const customUpload = (params) => {
      return new Promise((resolve, reject) => {
        const reader = new FileReader()
        reader.onload = (e) => {
          imageUrl.value = e.target.result
          registerForm.avatar = e.target.result  // 存储 Base64 数据
          file.value = params.file  // 保存文件对象用于后续提交
          resolve({ success: true, url: e.target.result })
        }
        reader.onerror = (err) => {
          ElMessage.error('头像处理失败')
          reject(err)
        }
        reader.readAsDataURL(params.file)
      })
    }

    // 上传前验证
    const beforeAvatarUpload = (file) => {
      const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'
      const isLt2M = file.size / 1024 / 1024 < 2

      if (!isJpgOrPng) {
        ElMessage.error('头像只能是 JPG/PNG 格式!')
      }
      if (!isLt2M) {
        ElMessage.error('头像大小不能超过 2MB!')
      }
      return isJpgOrPng && isLt2M
    }

    const handleRegister = () => {
      registerFormRef.value.validate(async (valid) => {
        if (valid) {
          loading.value = true;
          try {
            // 创建 FormData 对象
            const formData = new FormData();
            formData.append('username', registerForm.username);
            formData.append('email', registerForm.email);
            formData.append('password', registerForm.password);

            // 添加文件（如果有）
            if (file.value) {
              formData.append('avatar', file.value);
            }

            // 提交表单数据
            await store.dispatch('user/register', formData);
            ElMessage.success('注册成功');
            router.push('/login');
          } catch (error) {
            ElMessage.error(error.response?.data?.message || '注册失败');
          } finally {
            loading.value = false;
          }
        }
      });
    };

    return {
      registerForm,
      registerFormRef,
      rules,
      loading,
      handleRegister,
      imageUrl,
      handleAvatarSuccess,
      beforeAvatarUpload,
      customUpload
    }
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: calc(100vh - 60px);
  background-color: #f5f7fa;
}

.register-card {
  width: 400px;
}

.el-button {
  width: 100%;
  margin-bottom: 10px;
}

/* 头像上传区域样式 */
.avatar-upload-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.square-uploader {
  display: flex;
  justify-content: center;
}

.upload-area {
  width: 120px;
  height: 120px;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  transition: all 0.3s;
}

.upload-area:hover {
  border-color: #409EFF;
  background-color: #f5f7fa;
}

.avatar-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
}

.plus-icon {
  font-size: 48px;
  color: #8c939d;
  font-weight: lighter;
}

.upload-tip {
  color: #999;
  font-size: 12px;
  margin-top: 8px;
}
</style>