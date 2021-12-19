<template>
  <div class="home">
    <ScrollMode v-if="nowTemplate === 'scroll'"></ScrollMode>
    <SingleMode v-if="nowTemplate === 'single'"></SingleMode>
  </div>
</template>

<script>
// @ is an alias to /src
import ScrollMode from "@/components/ScrollMode.vue";
import SingleMode from "@/components/SingleMode.vue";
import { useCookies } from "vue3-cookies";
export default {
  name: "Home", //默认为 default。如果 <router-view>设置了名称，则会渲染对应的路由配置中 components 下的相应组件。
  components: {
    ScrollMode,
    SingleMode,
  },
  setup() {
    const { cookies } = useCookies(); 
    return { cookies };
  },
  data() {
    return {
      setting: null,
    };
  },
  beforeMount() {
    this.axios
      .get("/setting.json")
      .then((response) => {
        if (response.status == 200) {
          this.setting = response.data;
          console.log("get setting : "+this.setting);
        }
      })
      .catch((error) => alert(error));
      
  },
  computed: {
    // 计算属性的 getter
    nowTemplate: function () {
      // var localValue ='scroll'
      //document.cookie="nowTemplate=scroll"
      // this.cookies.set("nowTemplate",'scroll');
      var localValue = this.cookies.get("nowTemplate");
      console.log("nowTemplate is "+localValue);
      if (localValue !== null) {
        return localValue;
      } else {
        if (this.setting.template !== null) {
          return this.setting.template;
        } else {
          return ""
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
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  background-color: #f6f7eb;
  align-items: center;
}
/* 覆盖8px的浏览器默认值 */
body {
  margin: 0px;
  padding: 0px;
}
</style>
