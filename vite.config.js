import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  define: {
    __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: true, // 启用 hydration 不匹配的详细信息
    __VUE_OPTIONS_API: true
  },
  server: {
    port: 8081, // 固定前端端口为 8080
    strictPort: true, // 端口被占用时直接报错
    host: '0.0.0.0', // 若需通过 IP 访问（如局域网测试）
},})