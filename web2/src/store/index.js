import { createStore } from 'vuex'

export default createStore({
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
      page_num: 1,
      pages: [
        {
          height: 500,
          width: 449,
          url: "/resources/favicon.ico",
          class: "Vertical",
        },
      ],
    },
    bookshelf: {},
    setting: {
      template: "scroll",
      sketch_count_seconds: 90,
    },
    message: {
      user_uuid: "",
      server_status: "",
      now_book_uuid: "",
      read_percent: 0.0,
      msg: "",
    },
  },
  //mutaitions 只能执行同步操作
  mutations: {
    change_template_to_scroll(state) {
      state.setting.template = "scroll";
      console.log("change_template_to_scroll:" + state.setting.template);
    },
    change_template_to_double(state) {
      state.setting.template = "double";
      console.log("change_template_to_double:" + state.setting.template);
    },
    change_template_to_single(state) {
      state.setting.template = "single";
      console.log("change_template_to_single:" + state.setting.template);
    },
    change_template_to_sketch(state) {
      state.setting.template = "sketch";
      console.log("change_template_to_sketch:" + state.setting.template);
    },
    increment(state) {
      state.count++;
    },
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
  actions: {
    incrementAction(context) {
      context.commit("increment");
    },
    //拉取远程设定数据
    async syncSettingDataAction(context) {
      const msg = await this.axios.get("/setting.json").then(
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
      const msg = await this.axios.get("/book.json").then(
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
      const msg = await this.axios.get("/bookshelf.json").then(
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
