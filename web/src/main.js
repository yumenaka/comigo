import Vue from 'vue'
import App from './App.vue'
import router from './router'
import VueLazyload from 'vue-lazyload'
import websocket from 'vue-native-websocket'
import vuetify from './plugins/vuetify';
import VueCookies from 'vue-cookies'

Vue.use(websocket, "ws://" + document.location.host + "/ws", {//服务器的地址
  reconnection: true, // (Boolean)是否自动重连，默认false
  reconnectionAttempts: 500, // 重连次数
  reconnectionDelay: 1000, // 再次重连等待时间间隔(1000)
})

Vue.config.productionTip = false

// Vue.use(MintUI)
Vue.use(VueLazyload, {
  preLoad: 4.5,
  attempt: 10,
})

// https://github.com/cmp-cc/vue-cookies
Vue.use(VueCookies)
Vue.$cookies.config('30d')

new Vue({
  router,
  vuetify,
  render: h => h(App)
}).$mount('#app')
