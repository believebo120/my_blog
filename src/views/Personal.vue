<template>
  <el-container class="personal-center-container">
    <el-main>
      <el-card class="user-info-card">
        <!-- 背景图区域 -->
        <div class="user-bg-image" :style="{ backgroundImage: `url(${userInfo.backgroundImage})` }" @click="openBgDialog" :title="userInfo.backgroundImage || '点击修改背景图'"></div>

        <!-- 头像区域 -->
        <div class="user-avatar" @click="openAvatarDialog">
          <el-avatar :size="100" :src="userInfo.avatar" />
        </div>

        <!-- 用户信息显示/编辑区域 -->
        <div class="user-info">
          <!-- 用户名显示/编辑 -->
          <div class="info-row">
            <span class="info-label">用户名:</span>
            <span v-if="!isEditing.username" class="info-value" @click="toggleEdit('username')">
              {{ userInfo.username }}
              <span class="edit-icon" @click.stop>✏</span>
            </span>
            <el-input v-else v-model="userInfo.username" ref="usernameInput" @blur="handleBlur('username')" />
            <div v-else class="edit-actions">
              <el-button type="primary" @click="saveField('username')">保存</el-button>
              <el-button @click="cancelEdit('username')">取消</el-button>
            </div>
          </div>

          <!-- 邮箱显示/编辑 -->
          <div class="info-row">
            <span class="info-label">邮箱:</span>
            <span v-if="!isEditing.email" class="info-value" @click="toggleEdit('email')">
              {{ userInfo.email || '未设置' }}
              <span class="edit-icon" @click.stop>✏</span>
            </span>
            <el-input v-else v-model="userInfo.email" type="text" ref="emailInput" @blur="handleBlur('email')" />
            <div v-else class="edit-actions">
              <el-button type="primary" @click="saveField('email')">保存</el-button>
              <el-button @click="cancelEdit('email')">取消</el-button>
            </div>
          </div>

          <!-- 修改密码按钮 -->
          <div class="info-row">
            <span class="info-label">密码:</span>
            <el-button type="text" @click="openPasswordDialog">修改密码</el-button>
          </div>
        </div>

        <!-- 头像修改对话框 -->
        <el-dialog v-model="avatarDialogVisible" title="修改头像" width="400px">
          <el-upload
              class="avatar-uploader"
              action="/api/upload/avatar"
              :show-file-list="false"
              :before-upload="handleBeforeAvatarUpload"
              :on-success="handleUploadSuccess('avatar')"
          >
            <el-avatar :size="100" :src="userInfo.avatar" />
          </el-upload>
        </el-dialog>

        <!-- 背景图修改对话框 -->
        <el-dialog v-model="bgDialogVisible" title="修改背景图" width="600px">
          <el-upload
              class="bg-uploader"
              action="/api/upload/background"
              :show-file-list="false"
              :before-upload="handleBeforeBgUpload"
              :on-success="handleUploadSuccess('background')"
          >
            <div class="preview-bg" :style="{ backgroundImage: `url(${userInfo.backgroundImage})` }"></div>
          </el-upload>
        </el-dialog>

        <!-- 修改密码对话框 -->
        <el-dialog v-model="passwordDialogVisible" title="修改密码" width="400px">
          <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="100px">
            <el-form-item label="当前密码" prop="oldPassword">
              <el-input v-model="passwordForm.oldPassword" type="password" />
            </el-form-item>
            <el-form-item label="新密码" prop="newPassword">
              <el-input v-model="passwordForm.newPassword" type="password" />
            </el-form-item>
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input v-model="passwordForm.confirmPassword" type="password" />
            </el-form-item>
          </el-form>
          <template #footer>
            <el-button @click="passwordDialogVisible = false">取消</el-button>
            <el-button type="primary" @click="handleChangePassword">确定</el-button>
          </template>
        </el-dialog>
      </el-card>
    </el-main>
  </el-container>
</template>

<script setup>
import { ref, onMounted, nextTick, reactive } from 'vue';
import { ElMessage } from 'element-plus';
import { getUserInfo, updateUserInfo, changePassword } from '@/api/user';
const usernameInput = ref(null); // 定义用户名输入框的 ref
const emailInput = ref(null); // 定义邮箱输入框的 ref
// 用户信息
const userInfo = ref({
  id: null,
  username: '',
  email: '',
  avatar: '/uploads/default_avatar.jpg',
  backgroundImage: '/uploads/default_bg.jpg',
});

// 编辑状态管理
const isEditing = ref({
  username: false,
  email: false
});

// 原始值存储
const originalValues = ref({
  username: '',
  email: ''
});

// 密码表单
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
});

const passwordFormRef = ref(null);
const passwordDialogVisible = ref(false);

// 密码验证规则
const passwordRules = reactive({
  oldPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value) => {
        if (value !== passwordForm.newPassword) {
          return Promise.reject('两次输入的密码不一致');
        }
        return Promise.resolve();
      },
      trigger: 'blur'
    }
  ]
});

// 对话框状态
const avatarDialogVisible = ref(false);
const bgDialogVisible = ref(false);

// 获取用户信息
onMounted(() => {
  fetchUserInfo();
});

const fetchUserInfo = async () => {
  try {
    const res = await getUserInfo();
    userInfo.value = res.data;
    originalValues.value.username = res.data.username;
    originalValues.value.email = res.data.email;
  } catch (error) {
    ElMessage.error('获取用户信息失败');
    console.error(error);
  }
};

// 切换编辑状态
const toggleEdit = (field) => {
  isEditing.value[field] = true;
  nextTick(() => {
    if (field === 'username') {
      usernameInput.value.$refs.input.focus();
    } else if (field === 'email') {
      emailInput.value.$refs.input.focus();
    }
  });
};

// 保存字段修改
const saveField = async (field) => {
  if (!userInfo.value.id) return ElMessage.error('用户ID不存在');

  const newValue = userInfo.value[field];
  const originalValue = originalValues.value[field];

  if (newValue === originalValue) {
    return cancelEdit(field);
  }

  try {
    const updateData = { [field]: newValue };
    await updateUserInfo(userInfo.value.id, updateData);

    ElMessage.success(`${field === 'username' ? '用户名' : '邮箱'}修改成功`);
    originalValues.value[field] = newValue;
    isEditing.value[field] = false;
  } catch (error) {
    ElMessage.error(`修改失败: ${error.message || '请稍后重试'}`);
    console.error(error);
    // 恢复原始值
    userInfo.value[field] = originalValues.value[field];
    isEditing.value[field] = false;
  }
};

// 取消修改
const cancelEdit = (field) => {
  userInfo.value[field] = originalValues.value[field];
  isEditing.value[field] = false;
};

// 失去焦点时自动保存
const handleBlur = (field) => {
  if (isEditing.value[field] && userInfo.value[field] !== originalValues.value[field]) {
    saveField(field);
  }
};

// 打开密码修改对话框
const openPasswordDialog = () => {
  passwordDialogVisible.value = true;
  // 重置表单
  passwordForm.oldPassword = '';
  passwordForm.newPassword = '';
  passwordForm.confirmPassword = '';
  passwordFormRef.value.resetFields();
};

// 处理密码修改
const handleChangePassword = async () => {
  passwordFormRef.value.validate(async (valid) => {
    if (!valid) return;

    try {
      await changePassword(userInfo.value.id, {
        oldPassword: passwordForm.oldPassword,
        newPassword: passwordForm.newPassword
      });

      ElMessage.success('密码修改成功，请使用新密码登录');
      passwordDialogVisible.value = false;
    } catch (error) {
      ElMessage.error(`修改失败: ${error.message || '原密码不正确'}`);
      console.error(error);
    }
  });
};

// 头像上传相关方法（保持不变）
const handleBeforeAvatarUpload = (file) => {
  const isImage = /^image\/(jpg|jpeg|png|gif)$/.test(file.type);
  const isLt2M = file.size / 1024 / 1024 < 2;
  if (!isImage) {
    ElMessage.error('请上传JPG/JPEG/PNG/GIF格式的图片');
    return false;
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过2MB');
    return false;
  }
  return true;
};

// 背景图上传相关方法（保持不变）
const handleBeforeBgUpload = (file) => handleBeforeAvatarUpload(file);

// 通用上传成功处理函数
const handleUploadSuccess = (type) => (response, file) => {
  if (type === 'avatar') {
    userInfo.value.avatar = response.url;
  } else {
    userInfo.value.backgroundImage = response.url;
  }

  type === 'avatar' ? avatarDialogVisible.value = false : bgDialogVisible.value = false;
  ElMessage.success(`${type === 'avatar' ? '头像' : '背景图'}更新成功`);
};

// 打开头像对话框
const openAvatarDialog = () => {
  avatarDialogVisible.value = true;
};

// 打开背景图对话框
const openBgDialog = () => {
  bgDialogVisible.value = true;
};
</script>

<style scoped>
.personal-center-container {
  margin: 20px auto;
  max-width: 1200px;
}

.user-bg-image {
  width: 100%;
  height: 300px;
  background-size: cover;
  background-position: center;
  cursor: pointer;
  margin-bottom: 30px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.user-avatar {
  margin: -80px auto 30px;
  text-align: center;
  position: relative;
  z-index: 1;
}

.user-info {
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.info-row {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.info-label {
  width: 80px;
  color: #606266;
}

.info-value {
  flex: 1;
  color: #303133;
  cursor: pointer;
  display: flex;
  align-items: center;
}

.edit-icon {
  margin-left: 8px;
  font-size: 0.9em;
  color: #409eff;
  opacity: 0.7;
  transition: opacity 0.3s;
}

.info-value:hover .edit-icon {
  opacity: 1;
}

.edit-actions {
  margin-left: 10px;
  display: flex;
  gap: 10px;
}

.avatar-uploader,
.bg-uploader {
  width: 100%;
  padding: 30px;
  text-align: center;
}

.preview-bg {
  width: 100%;
  height: 400px;
  background-size: cover;
  background-position: center;
  border: 1px dashed #e0e0e0;
  margin-bottom: 20px;
}
</style>