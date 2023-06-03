// https://github.com/likaia/pinia-websocket-project
import { defineStore } from "pinia";
import { store } from "@/stores/store";
import main from "@/main";
import { SocketStore } from "@/stores/PiniaType";

export const useSocketStore = defineStore({
  id: "socket",
  state: (): SocketStore => ({
    // 连接状态
    isConnected: false,
    // 消息内容
    message: "",
    // 重新连接错误
    reconnectError: false,
    // 心跳消息发送时间
    heartBeatInterval: 50000,
    // 心跳定时器
    heartBeatTimer: 0
  }),
  actions: {
    // 连接打开
    SOCKET_ONOPEN(event: any) {
      console.log("successful websocket connection");
      main.config.globalProperties.$socket = event.currentTarget;
      this.isConnected = true;
      // 连接成功时启动定时发送心跳消息，避免被服务器断开连接
      this.heartBeatTimer = window.setInterval(() => {
        const message = "心跳消息";
        this.isConnected &&
          main.config.globalProperties.$socket.sendObj({
            code: 200,
            msg: message
          });
      }, this.heartBeatInterval);
    },
    // 连接关闭
    SOCKET_ONCLOSE(event: any) {
      this.isConnected = false;
      // 连接关闭时停掉心跳消息
      window.clearInterval(this.heartBeatTimer);
      this.heartBeatTimer = 0;
      console.log("连接已断开: " + new Date());
      console.log(event);
    },
    // 发生错误
    SOCKET_ONERROR(event: any) {
      console.error(event);
    },
    // 收到服务端发送的消息
    SOCKET_ONMESSAGE(message: any) {
      this.message = message;
    },
    // 自动重连
    SOCKET_RECONNECT(count: any) {
      console.info("消息系统重连中...", count);
    },
    // 重连错误
    SOCKET_RECONNECT_ERROR() {
      this.reconnectError = true;
    }
  }
});

// Need to be used outside the setup
export function useSocketStoreWithOut() {
  return useSocketStore(store);
}