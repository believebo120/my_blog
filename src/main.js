import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import axios from 'axios'

const app = createApp(App)

// 配置axios
axios.defaults.baseURL = '/api'
const token = localStorage.getItem('token')
if (token) {
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
}

app.config.globalProperties.$http = axios

app.use(store)
   .use(router)
   .use(ElementPlus)
   .mount('#app') 