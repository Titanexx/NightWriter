import { createApp } from 'vue'
import './index.scss'
import App from './App.vue'
import router from './routes'
// import Sortable from 'sortablejs'

createApp(App)
    .use(router)
    .mount('#app')