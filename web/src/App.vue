<template>
  <div class="app">
    <n-dialog-provider>
      <n-message-provider>
        <router-view v-slot="{Component}">
          <keep-alive>
          <component :is="Component" v-if="$route.meta.keepAlive" :key="$route.name" />
          </keep-alive>
          <component :is="Component" v-if="!$route.meta.keepAlive" :key="$route.name" />
        </router-view> <!-- 路由出口 路由匹配到的组件将渲染在这里 -->
      </n-message-provider>
    </n-dialog-provider>
  </div>
</template>

<script lang="ts">
// @ is an alias to /src
// import ScrollMode from "@/views/ScrollMode.vue";
// import FlipMode from "@/views/FlipMode.vue";
// import BookShelf from "@/views/BookShelf.vue";
import { useCookies } from "vue3-cookies";
import { defineComponent } from 'vue'
import { NMessageProvider, NDialogProvider, darkTheme, lightTheme } from 'naive-ui'
import { StringChain } from "lodash";
export default defineComponent({
  name: "ComigoHome", //默认为 default。如果 <router-view>设置了名称，则会渲染对应的路由配置中 components 下的相应组件。
  components: {
    NMessageProvider,
    NDialogProvider,
    // NConfigProvider,//调整主题：https://www.naiveui.com/zh-CN/light/docs/customize-theme
  },
  setup() {
    const { cookies } = useCookies();
    return { cookies, darkTheme, lightTheme };
  },
  data() {
    return {
      selectTemplate: "",
      isAuthenticated: false,
    };
  },
  created() {
    // this.$store.dispatch("syncBookDataAction");
    this.$store.dispatch("syncSeverStatusDataAction");
    // this.$store.dispatch("syncBookShelfDataAction");
    this.selectTemplate = this.getDefaultTemplate;

    // // 连接websocket服务器，参数为websocket服务地址
    // var protocol = 'ws://'
    // if (window.location.protocol === "https") {
    //   protocol = 'wss://'
    // }
    // var ws_url = protocol + window.location.host + '/api/ws';
    // this.$connect(ws_url);
    // console.log("ws_url:" + ws_url)
  },
  beforeMount() {
    if (this.$store.state.server_status.ServerName != null) {
      document.title = this.$store.state.server_status.ServerName
    }
  },
  methods: {
    goToDashboard() {
      if (this.isAuthenticated) {
        this.$router.push('/dashboard')
      } else {
        this.$router.push('/login')
      }
    },
    OnSetTemplate(value: string) {
      localStorage.setItem("nowTemplate", value);
      this.selectTemplate = value;
    },
  },
  //计算属性
  computed: {
    username() {
      // 我们很快就会看到 `params` 是什么
      return this.$route.params.username
    },
    book() {
      return this.$store.state.book;
    },
    setting() {
      return this.$store.state.setting;
    },
    getDefaultTemplate: function () {
      var localValue = localStorage.getItem("nowTemplate");
      if (localValue !== null) {
        return localValue;
      }
      return "scroll";
    },
  },
});
</script>

<style>
.app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  align-items: center;
}

/* 覆盖8px的浏览器默认值 */
* {
  /* 外边距，如果不指定，浏览器默认设置成8px 造成一些现实问题 */
  margin: 0px;
  /* 内边框 */
  padding: 0px;
}
</style>