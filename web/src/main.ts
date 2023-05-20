import { createApp } from "vue";
import VueAxios from "vue-axios";
import axios from "axios";
import App from "@/App.vue";
import router from "@/router"; //vue-router
import store from "@/store"; //VueX
// import VueCookies from "vue3-cookies";
import VueLazyLoad from "vue3-lazyload";
import i18n from "@/locales";
import VueNativeSock from "vue-native-websocket-vue3";

// 后端改成 /api/book/:id的形式
axios.defaults.baseURL = "/api";

// axios 全局配置拦截器  每次向后端请求携带 头信息 Authorization。不知为什么不起作用，目前是cookie验证。
axios.interceptors.request.use(
    (config) => {
        if (sessionStorage.JWT_TOKEN) {
            // 判断是否存在token，如果存在的话，则每个http header都加上token
            config.headers.Authorization = "Bearer "+`${sessionStorage.JWT_TOKEN}`;
            // document.cookie= "id:"+`${sessionStorage.JWT_TOKEN}`;
        }
        return config;
    },
    (err) => {
        return Promise.reject(err);
    }
);

// createApp(App).use(store).use(router).use(VueAxios,axios).mount('#app')
const app = createApp(App);
app.use(i18n);
app.use(VueAxios, axios);
// 传入 injection key
app.use(store);
app.use(router);

// Tailwind CSS
import "./index.css";

// vue3-lazyload
// https://github.com/murongg/vue3-lazyload
app.use(VueLazyLoad, {
    //懒加载相关设置
    //https://www.cnblogs.com/niuzijie/p/13703710.html
    preLoad: 1.5, //预加载高度比例,默认1.3
    loading: "/images/loading.gif",
    error: "/images/error.jpg",
    attempt: 3, //尝试加载图片数量，默认3
    observerOptions: { rootMargin: "200px", threshold: 0.1 },
    lifecycle: {
        loading: (el: any) => {
            el.setAttribute("class", "LoadingImage");
            // console.log("loading", el);
        },
        error: (el: any) => {
            el.setAttribute("class", "ErrorImage");
            // console.log("error", el);
        },
        //可以在这里插入判断分辨率的函数
        loaded: (el: any) => {
            let image = new Image();
            image.src = el.src;
            // 图片是否完全加载完成。
            //https://developer.mozilla.org/zh-CN/docs/Web/API/HTMLImageElement/complete
            if (image.complete) {
                el.setAttribute("w", image.width);
                el.setAttribute("h", image.height);
                if (image.width < image.height) {
                    el.setAttribute("class", "SinglePageImage");
                } else {
                    el.setAttribute("class", "DoublePageImage");
                }
            } else {
                el.setAttribute("class", "SinglePageImage");
            }
            // console.log("loaded", el);
        },
    },
});

// 使用VueNativeSock插件，并进行相关配置
// 参考https://github.com/likaia/vue-native-websocket-vue3
var protocol = "ws://";
// console.log(window.location.protocol)
if (window.location.protocol === "https:") {
    protocol = "wss://";
}
var ws_url = protocol + window.location.host + "/api/ws";
app.use(VueNativeSock, ws_url, {
    // 启用Vuex集成
    store: store,
    // 数据发送/接收使用使用json
    format: "json",
    // 开启手动调用 connect() 连接服务器
    connectManually: false,
    // 开启自动重连
    reconnection: true,
    // 尝试重连的次数
    reconnectionAttempts: 60,
    // 重连间隔时间
    reconnectionDelay: 3000, //掉线后每3秒重连一次
});

app.mount("#app"); // look index.html:  <div id="app"></div>

//store的websockets需要导入main，所以要有这一句，参考： https://github.com/likaia/chat-system/blob/master/src/main.ts
export default app;
