import { createStore } from 'vuex'
import axios from 'axios'

const store = createStore({
  //开启严格模式，避免直接的修改
  // strict :true,
  //state 存放所有组件的共用数据
  state: {
    count: 0,
    todos: [
      { id: 1, text: "...", done: true },
      { id: 2, text: "...", done: false },
    ],
    now_page: 1,
    book: {
      name: "loading",
      all_page_num: 1,
      pages: [
        {
          height: 500,
          width: 449,
          url: "/images/loading.jpg",
        },
      ],
    },
    bookshelf: {},
    setting: {
      template: "scroll",
      sketch_count_seconds: 30,
    },
    message: {
      user_id: "",
      server_status: "",
      now_book_id: "",
      read_percent: 0.0,
      msg: "",
    },
  },
  //mutaitions改变 只能执行同步操作。不能直接调用。需要使用 store.commit('函数名') 方法
  mutations: {
    change_template_to_scroll(state) {
      state.setting.template = "scroll";
      console.log("change_template_to_scroll:" + state.setting.template);
    },
    change_template_to_flip(state) {
      state.setting.template = "flip";
      console.log("change_template_to_flip:" + state.setting.template);
    },
    change_template_to_sketch(state) {
      state.setting.template = "sketch";
      console.log("change_template_to_sketch:" + state.setting.template);
    },
    increment(state) {
      state.count++;
    },
    //使用mutation()函数（store.commit('')）的时候，还可以传入额外的参数，也就是载荷payload
    syncSettingData(state, payload) {
      state.setting = payload.message;
    },
    syncBookData(state, payload) {
      state.book = payload.message;
    },
    syncBookShelfData(state, payload) {
      state.bookshelf = payload.message;
    },
  },
  //actions 可以包含任意异步操作，通过 store.dispatch 方法触发
  //接收context，与store具有相同方法与属性，不是store本身，
  //context可以访问state与getter，还可以用context.dispatch调用其他Action
  actions: {
    incrementAction(context) {
      context.commit("increment");
    },
    //拉取远程设定数据
    async syncSettingDataAction(context) {
      const msg = await axios.get("/setting.json").then(
        (res) => res.data,
        () => ""
      );
      const payload = {
        message: msg,
      };
      context.commit("syncSettingData", payload);
      console.log("syncSettingData!");
    },
    //拉取当前阅读书籍数据
    async syncBookDataAction(context) {
      const msg = await axios.get("/book.json").then(
        (res) => res.data,
        () => ""
      );
      const payload = {
        message: msg,
      };
      context.commit("syncBookData", payload);
      console.log("syncBookData!");
    },
    //拉取书架数据
    async syncBookShelfDataAction(context) {
      const msg = await axios.get("/bookshelf.json").then(
        (res) => res.data,
        () => ""
      );
      const payload = {
        message: msg,
      };
      context.commit("syncBookShelfData", payload);
      console.log("syncBookShelfData!");
    },
  },
  //相当于store的计算属性，会被缓存，变化的时候才重新计算
  getters: {
    doneTodos: (state) => {
      return state.todos.filter((todo) => todo.done);
    },
    now_page: (state) => {
      return state.now_page;
    },
    book: (state) => {
      return state.book;
    },
    bookshelf: (state) => {
      return state.bookshelf;
    },
    setting: (state) => {
      return state.setting;
    },
    message: (state) => {
      return state.message;
    },
  },
  modules: {
  }
})


export default store