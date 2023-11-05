import { defineStore } from "pinia";
import axios from "axios";

//生成一个随机ID
var tempUserID = "Comigo_" + Math.floor(Math.random() * 100000); //可均衡获取 0 到 99999 的随机整数。
var temp = localStorage.getItem("ComigoTempUserID");
if (typeof temp === "string") {
    tempUserID = temp;
} else {
    localStorage.setItem("ComigoTempUserID", tempUserID);
}

// 从 Vuex ≤4 迁移  https://pinia.vuejs.org/zh/cookbook/migration-vuex.html

// defineStore() 第一个参数是 storeId ，第二个参数是一个选项对象
export const userStore = defineStore("user", {
    // 1. 状态值定义
    state: () => {
        return {
            now_page: 1,
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
            onlineDerviceNumber: 1,
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
        };
    },
    // 2. getters 。getter 可以直接访问 state ，也可以调用其他 getter 。
    //  需要预处理数据，比如返回排序好的数据时用
    getters: {
        doubleNowPage: (state) => state.now_page * 2,
    },
    // 3. actions，状态更改方法，同步或异步修改 state 的地方
    // 再也没有什么mutations了、State 的更新全部用 actions 做。
    actions: {
        // 获取服务器状态
        syncSeverStatusData() {
            let _this = this;
            axios
                .get("get_status")
                .then(function (response) {
                    _this.server_status = response.data;
                    console.log(response.data);
                })
                .catch(function (error) {
                    console.log(error);
                });
            console.log("syncSeverStatusData!");
        },
    },
});
