<template>
  <div class="home">
    <ScrollMode v-if="nowTemplate === 'scroll'" :book="this.book"></ScrollMode>
    <FlipMode v-if="nowTemplate === 'flip' || nowTemplate === 'sketch'" :book="this.book"></FlipMode>
  </div>
</template>

<script>
// @ is an alias to /src
import ScrollMode from "@/views/ScrollMode.vue";
import FlipMode from "@/views/FlipMode.vue";
import { useCookies } from "vue3-cookies";
import { defineComponent } from 'vue'

export default defineComponent({
  name: "Home", //默认为 default。如果 <router-view>设置了名称，则会渲染对应的路由配置中 components 下的相应组件。
  components: {
    ScrollMode,
    FlipMode,
  },
  setup() {
    const { cookies } = useCookies();
    return { cookies };
  },
  data() {
    return {
    };
  },
  created() {
    this.$store.dispatch("syncBookDataAction");
    this.$store.dispatch("syncSettingDataAction");
    this.$store.dispatch("syncBookShelfDataAction");
  },
  beforeMount() {
  },
  computed: {
    book(){
      return this.$store.state.book;
    },
    setting(){
      return this.$store.state.setting;
    },
    // 计算属性的 getter
    nowTemplate: function () {
      var localValue = this.cookies.get("nowTemplate");
      if (localValue !== null) {
        return localValue;
      } else {
        if (this.setting.template) {
          this.cookies.set("nowTemplate", this.setting.template)
          return this.setting.template;
        } else {
          return "scroll"
        }
      }
    },
  },
  methods: {
    toTop() {

    },
    getNumber: function (number) {
      this.page = number;
      console.log(number);
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
