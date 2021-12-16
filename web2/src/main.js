import { createApp } from "vue";
import VueAxios from "vue-axios";
import axios from "axios";
import App from "./App.vue";
// import router from "./router"; //vue-router
import store from "./store"; //VueX
// import VueCookies from "vue3-cookies";
import VueLazyLoad from "vue3-lazyload";

// createApp(App).use(store).use(router).use(VueAxios,axios).mount('#app')
const app = createApp(App);
app.use(VueAxios, axios);
app.use(store);
// app.use(router);
app.mount("#app"); // look index.html:  <div id="app"></div>

// 以后后端改成 /api/book/:id的形式
axios.defaults.baseURL ="/api"

// Set default vue3-cookies config:
// https://github.com/KanHarI/vue3-cookies
// app.use(VueCookies, {
//     expireTimes: "30d",
//     path: "/",
//     domain: "",
//     secure: true,
//     sameSite: "None",
// });

// vue3-lazyload
// https://github.com/murongg/vue3-lazyload
app.use(VueLazyLoad, {
    loading: "",
    error: "",
    lifecycle: {
        loading: (el) => {
            console.log("loading", el);
        },
        error: (el) => {
            console.log("error", el);
        },
        loaded: (el) => {
            console.log("loaded", el);
        },
    },
});
