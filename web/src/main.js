import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import VueLazyload from "vue-lazyload";
import websocket from "vue-native-websocket";
import vuetify from "./plugins/vuetify";
import VueCookies from "vue-cookies";
import Vuex from "vuex";
import axios from "axios";

Vue.use(websocket, "ws://" + document.location.host + "/ws", {
  //服务器的地址
  reconnection: true, // (Boolean)是否自动重连，默认false
  reconnectionAttempts: 500, // 重连次数
  reconnectionDelay: 1000, // 再次重连等待时间间隔(1000)
});

Vue.config.productionTip = false;

// Vue.use(MintUI)
Vue.use(VueLazyload, {
  preLoad: 4.5,
  attempt: 10,
});

// https://github.com/cmp-cc/vue-cookies
Vue.use(VueCookies);
Vue.$cookies.config("30d");

//https://vuex.vuejs.org/zh/guide/
Vue.use(Vuex);

const store = new Vuex.Store({
  state: {
    count: 0,
    todos: [
      { id: 1, text: '...', done: true },
      { id: 2, text: '...', done: false }
    ],
    now_page: 1,
    book: {
      name: "loading",
      page_num: 1,
      pages: [
        {
          height: 2000,
          width: 1419,
          url: "/resources/favicon.ico",
          class: "Vertical",
        },
      ],
    },
    bookshelf: {},
    defaultSetiing: {
      default_page_template: "???",
    },
    message: {
      user_uuid: "",
      server_status: "",
      now_book_uuid: "",
      read_percent: 0.0,
      msg: "",
    },
  },
  getters: {
    doneTodos: state => {
      return state.todos.filter(todo => todo.done)
    },
    now_page: state => {
      return state.now_page;
    },
    book: state => {
      return state.book;
    },
    bookshelf: state => {
      return state.bookshelf;
    },
    defaultSetiing: state => {
      return state.defaultSetiing;
    },
    message: state => {
      return state.message;
    },
  },
  mutations: {
    increment(state) {
      state.count++;
    },
    syncRemoteSetting() {
      axios
        .get("/bookshelf.json")
        .then((response) => (bookshelf = response.data))
        .finally();
    },
    syncBookDate() {
      axios.get("/book.json").then((response) => (ook = response.data));
    },
    syncBookShelfDate() {
      axios
        .get("/bookshelf.json")
        .then((response) => (bookshelf = response.data))
        .finally();
    },
  },
});

new Vue({
  router,
  vuetify,
  render: (h) => h(App),
  store: store,
}).$mount("#app");
