<template>
  <div class="app">
    <n-message-provider>
      <!-- <ScrollMode
        v-if="selectTemplate === 'scroll'"
        :book="this.book"
        :nowTemplate="this.selectTemplate"
        @setTemplate="OnSetTemplate"
      ></ScrollMode>
      <FlipMode
        v-if="selectTemplate === 'flip' || selectTemplate === 'sketch'"
        :book="this.book"
        :nowTemplate="this.selectTemplate"
        @setTemplate="OnSetTemplate"
      ></FlipMode>
      <BookShelf
        v-if="selectTemplate === 'bookshelf'"
        :book="this.book"
        :nowTemplate="this.selectTemplate"
        @setTemplate="OnSetTemplate"  
      ></BookShelf>-->

      <!--使用 router-link 来导航 -->
      <h1>Comigo 路由测试</h1>
      <p>
        <router-link to="/">Router-link: home</router-link>
        <router-link to="/about">Router-link: about</router-link>
        <router-link to="/book/3AzY2">Router-link: book</router-link>
        <!-- <router-link :to="{ name: 'ScrollMode', params: { book_id: '3AzY2' } }">ScrollMode</router-link> -->
      </p>
      <!-- 路由出口 路由匹配到的组件将渲染在这里 -->
      <router-view></router-view>
    </n-message-provider>
  </div>
</template>

<script>
// @ is an alias to /src
// import ScrollMode from "@/views/ScrollMode.vue";
// import FlipMode from "@/views/FlipMode.vue";
// import BookShelf from "@/views/BookShelf.vue";
import { useCookies } from "vue3-cookies";
import { defineComponent } from 'vue'
import { NMessageProvider } from 'naive-ui'

export default defineComponent({
  name: "Home", //默认为 default。如果 <router-view>设置了名称，则会渲染对应的路由配置中 components 下的相应组件。
  components: {
    // ScrollMode,
    // FlipMode,
    // BookShelf,
    NMessageProvider,
  },
  setup() {
    const { cookies } = useCookies();
    return { cookies };
  },
  data() {
    return {
      selectTemplate: "",
      isAuthenticated: false,
    };
  },
  created() {
    this.$store.dispatch("syncBookDataAction");
    this.$store.dispatch("syncSettingDataAction");
    this.$store.dispatch("syncBookShelfDataAction");
    this.selectTemplate = this.getDefaultTemplate;
  },
  beforeMount() {
  },

  methods: {
    goToDashboard() {
      if (this.isAuthenticated) {
        this.$router.push('/dashboard')
      } else {
        this.$router.push('/login')
      }
    },
    OnSetTemplate(value) {
      localStorage.setItem("nowTemplate", value);
      this.selectTemplate = value;
    },
    getNumber: function (number) {
      this.page = number;
      console.log(number);
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
      //不管服务器设置，完全按照本地值来
      // if (this.setting.template) {
      //   localStorage.setItem("nowTemplate", this.setting.template)
      //   return this.setting.template;
      // }
      return "scroll";
    },
  },
});
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  /* 整体颜色，做成用户设定？ */
  background-color: #f6f7eb;
  align-items: center;
}
/* 覆盖8px的浏览器默认值 */
* {
  /* 外边距，不指定的话，浏览器默认设置成8px */
  margin: 0px;
  /* 内边框 */
  padding: 0px;
}
</style>
