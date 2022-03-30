import { createApp } from "vue";
import VueAxios from "vue-axios";
import axios from "axios";
import App from "@/App.vue";
import router from "@/router"; //vue-router
import store from "@/store"; //VueX
// import VueCookies from "vue3-cookies";
import VueLazyLoad from "vue3-lazyload";
import i18n from '@/locales'
// 后端改成 /api/book/:id的形式
axios.defaults.baseURL = "/api"

// createApp(App).use(store).use(router).use(VueAxios,axios).mount('#app')
const app = createApp(App);
app.use(i18n)
app.use(VueAxios, axios);
app.use(store);
app.use(router);

// Tailwind CSS
import './index.css'

// vue3-lazyload
// https://github.com/murongg/vue3-lazyload
app.use(VueLazyLoad, {
    //懒加载相关设置
    //https://www.cnblogs.com/niuzijie/p/13703710.html
    preLoad: 1.9,//预加载高度比例,默认1.3
    loading: "/images/loading.jpg",
    error: "/images/error.jpg",
    attempt: 4,//尝试加载图片数量，默认3
    observerOptions: { rootMargin: '200px', threshold: 0.1 },
    lifecycle: {
        loading: (el) => {
            el.setAttribute("class", "LoadingImage");
            // console.log("loading", el);
        },
        error: (el) => {
            el.setAttribute("class", "ErrorImage");
            // console.log("error", el);
        },
        //可以在这里插入判断分辨率的函数
        loaded: (el) => {
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

app.mount("#app"); // look index.html:  <div id="app"></div>