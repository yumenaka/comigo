import { createStore, useStore as baseUseStore, Store  } from "vuex";
import axios from "axios";
import main from "../main";
import { InjectionKey } from 'vue'

// Vuex TypeScript 支持 


// 为 store state 声明类型
export interface State {
  count: number
}
// 定义 injection key
export const key: InjectionKey<Store<State>> = Symbol()
// 定义自己的 `useStore` 组合式函数
export function useStore () {
  return baseUseStore(key)
}


//生成一个随机ID
var tempUserID = "Comigo_" + Math.floor(Math.random() * 100000); //可均衡获取 0 到 99999 的随机整数。
var temp = localStorage.getItem("ComigoTempUserID");
if (typeof temp === "string") {
  tempUserID = temp;
} else {
  localStorage.setItem("ComigoTempUserID", tempUserID);
}

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
          url: "/images/loading.gif",
        },
      ],
    },
    bookshelf: [
      {
        name: "loading",
        all_page_num: 1,
        pages: [
          {
            height: 500,
            width: 449,
            url: "/images/loading.gif",
          },
        ],
      },
    ],
    server_status: {
      template: "scroll",
      sketch_count_seconds: 30,
    },
    message: {
      user_id: tempUserID,
      server_status: "",
      now_book_id: "",
      read_percent: 0.0,
      msg: "",
    },
    userID: tempUserID,
    token: "abc123",
    refreshToken: "xxx321",
    // 用户头像
    profilePicture: "default",
    username: "adminTest",
    // 在线人数
    onlineUsers: 0,
    currentComponentName: null,
    //websockets相关，参考https://github.com/likaia/vue-native-websocket-vue3
    socket: {
      // 连接状态
      isConnected: false,
      // 消息内容
      message: "",
      // 重新连接错误
      reconnectError: false,
      // 心跳消息发送时间
      heartBeatInterval: 50000, //50000（50秒）一次的信条消息。为了测试，有时会把间隔改小一点
      // 心跳定时器
      heartBeatTimer: 0,
    },
  },
  //mutaitions改变 只能执行同步操作。不能直接调用。需要使用 store.commit('函数名') 方法
  mutations: {
    // 连接打开
    SOCKET_ONOPEN(state, event) {
      main.config.globalProperties.$socket = event.currentTarget;
      state.socket.isConnected = true;
      // 连接成功时启动定时发送心跳消息，避免被服务器断开连接
      state.socket.heartBeatTimer = setInterval(() => {
        var date_json = new Date(new Date()).toJSON();
        var date_str = new Date(+new Date(date_json) + 8 * 3600 * 1000)
          .toISOString()
          .replace(/T/g, " ")
          .replace(/\.[\d]{3}Z/, "");
        const detail = "【Websockets】heart Beat。" + date_str;
        state.socket.isConnected &&
          main.config.globalProperties.$socket.sendObj({
            type: "heartbeat",
            status_code: 200,
            user_id: state.userID,
            token: state.token,
            detail: detail,
          });
      }, state.socket.heartBeatInterval);
      console.log("临时客户端ID：", state.userID);
      console.log(
        "【Websockets】连接建立。 " +
          new Date().toLocaleDateString().replace(/\//g, "-") +
          " " +
          new Date().toTimeString().substr(0, 8)
      );
    },
    // 连接关闭
    SOCKET_ONCLOSE(state, event) {
      state.socket.isConnected = false;
      // 连接关闭时停掉心跳消息
      clearInterval(state.socket.heartBeatTimer);
      state.socket.heartBeatTimer = 0;
      console.log("【Websockets】连接已断开: " + new Date());
      // console.log(event);
    },
    // 发生错误
    SOCKET_ONERROR(state, event) {
      console.error(state, event);
    },
    // 收到服务端发送的消息
    SOCKET_ONMESSAGE(state, message) {
      state.socket.message = message;
      // console.info("收到服务器消息：",message);
    },
    // 自动重连
    SOCKET_RECONNECT(state, count) {
      console.info("【Websockets】消息系统重连中...", state, count);
    },
    // 重连错误
    SOCKET_RECONNECT_ERROR(state) {
      state.socket.reconnectError = true;
    },
    increment(state) {
      state.count++;
    },
    //使用mutation()函数（store.commit('')）的时候，还可以传入额外的参数，也就是载荷payload
    syncSeverStatusData(state, payload) {
      state.server_status = payload.message;
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
    async syncSeverStatusDataAction(context) {
      const msg = await axios
        .get("get_status")
        .then(
          (res) => res.data,
          () => ""
        )
        .finally(() => {});
      const payload = {
        message: msg,
      };
      context.commit("syncSeverStatusData", payload);
      console.log("syncSeverStatusData!");
    },
    //拉取当前阅读书籍数据
    async syncBookDataAction(context) {
      const msg = await axios.get("book.json").then(
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
      const msg = await axios.get("getshelf").then(
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
    message: (state) => {
      return state.message;
    },
  },
  modules: {},
});
export default store;
