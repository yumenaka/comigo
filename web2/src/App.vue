<template>
  <div class="home">
    <BookScroll></BookScroll>
  </div>
</template>

<script>
// @ is an alias to /src
import BookScroll from "@/components/BookScroll.vue";
export default {
  name: "Home", //默认为 default。如果 <router-view>设置了名称，则会渲染对应的路由配置中 components 下的相应组件。
  components: {
    BookScroll,
  },
  setup() {},
  data() {
    return {
    };
  },
  computed: {
    // 计算属性的 getter
    nowTemplate: function () {
      var localValue = this.$cookies.get("nowTemplate");
      console.log("computed 1:" + localValue);
      if (localValue !== null) {
        return localValue;
      } else {
        return this.$store.state.setting.template;
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
