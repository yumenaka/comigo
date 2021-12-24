<template>
  <div class="home">
    <ScrollMode
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
      selectTemplate: "",
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
    OnSetTemplate(value) {
      this.cookies.set("nowTemplate", value);
      this.selectTemplate = value; 
    },
    getNumber: function (number) {
      this.page = number;
      console.log(number);
    },
  },
  //计算属性
  computed: {
    book() {
      return this.$store.state.book;
    },
    setting() {
      return this.$store.state.setting;
    },
    getDefaultTemplate: function () {
      var localValue = this.cookies.get("nowTemplate");
      if (localValue !== null) {
        return localValue;
      }
      if (this.setting.template) {
        this.cookies.set("nowTemplate", this.setting.template)
        return this.setting.template;
      }
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
