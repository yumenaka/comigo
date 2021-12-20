import { createApp } from "vue";
import VueAxios from "vue-axios";
import axios from "axios";
import App from "./App.vue";
import router from "./router"; //vue-router
import store from "./store"; //VueX
// import VueCookies from "vue3-cookies";
import VueLazyLoad from "vue3-lazyload";

// createApp(App).use(store).use(router).use(VueAxios,axios).mount('#app')
const app = createApp(App);
// 以后后端改成 /api/book/:id的形式
axios.defaults.baseURL = "/api"
app.use(VueAxios, axios);
app.use(store);
app.use(router);
// // 通用字体
// import 'vfonts/Lato.css'
// // 等宽字体
// import 'vfonts/FiraCode.css'
// vue3-lazyload
// https://github.com/murongg/vue3-lazyload
app.use(VueLazyLoad, {
    loading: "loading.jpg",
    error: "error.jpg",
    //懒加载相关设置
    observerOptions: { rootMargin: '500px', threshold: 0.1 },
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
            // 如果有缓存，读缓存。
            //还要避免默认占位图片的情况，目前远程网速较慢时似乎会出错
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

//自定义指令 https://v3.cn.vuejs.org/guide/custom-directive.html
// app.directive('demo', {
//     beforeUpdate(el, binding) {
//         // console.log(el.getAttribute("class"))
//         // console.log(binding.value.singleWidth) 
//         // console.log(binding.value.doubleWidth)
//         if (el.getAttribute("class") === "SinglePageImage") {
//             el.style.width = binding.value.singleWidth
//         } else if (el.getAttribute("class") === "DoublePageImage") {
//             el.style.width = binding.value.doubleWidth
//         }
//         console.log("directive isLandscape:"+binding.value.isLandscape)
//     }
// })

app.mount("#app"); // look index.html:  <div id="app"></div>