import { createApp } from 'vue'
import VueAxios from 'vue-axios'
import axios from 'axios'
import App from './App.vue'
import router from './router'
import store from './store'

// createApp(App).use(store).use(router).use(VueAxios,axios).mount('#app')

const app = createApp(App);
app.use(store);
app.use(router);
app.use(VueAxios,axios);
app.mount('#app');// look index.html:  <div id="app"></div>
