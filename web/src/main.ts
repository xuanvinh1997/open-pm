import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import { useToast } from './composables/useToast'
import './assets/styles/main.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)

app.config.errorHandler = (err) => {
  console.error('Unhandled error:', err)
  const { error } = useToast()
  error(err instanceof Error ? err.message : 'An unexpected error occurred')
}

app.mount('#app')
