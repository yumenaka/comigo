<template>
  <div id="app">
    <!-- 初始化后才显示，避免 setting错误 -->
    <!-- <div v-if="this.$store.state.setting"> -->
    <div v-if="this.$store.state.setting">
      <!-- 下拉阅读 -->
      <ScrollTemplate v-if="nowTemplate === 'scroll'"> </ScrollTemplate>
      <!-- 绘图参考（倒计时速写什么的） -->
      <SketchTemplate v-if="nowTemplate === 'sketch'"> </SketchTemplate>
      <!-- 单页阅读 -->
      <SinglePageTemplate v-if="nowTemplate === 'single'"> </SinglePageTemplate>
      <!-- 双页阅读 -->
      <DoublePageTemplate v-if="nowTemplate === 'double'"> </DoublePageTemplate>
    </div>
    <!-- 加载中 -->
    <!-- <p v-else>loading.....</p> -->
  </div>
</template>

<script>
//代码参考：https://github.com/bradtraversy/vue_crash_todolist
import axios from "axios";
import ScrollTemplate from "./views/ScrollTemplate.vue";
import SketchTemplate from "./views/SketchTemplate.vue";
import SinglePageTemplate from "./views/SinglePageTemplate.vue";
import DoublePageTemplate from "./views/DoublePageTemplate.vue";

export default {
  name: "app",
  //为了能在模板中使用，组件必须先注册以便 Vue 能够识别
  components: {
    ScrollTemplate,
    SketchTemplate,
    SinglePageTemplate,
    DoublePageTemplate,
  },
  //组件的 data 选项必须是一个函数
  //每个实例可以维护一份被返回对象的独立的拷贝
  data() {
    return {
      book: null,
      //如果你知道你会在晚些时候需要一个 property，但是一开始它为空或不存在，那么你仅需要设置一些初始值。
      bookshelf: {},
      setting: {
        template: "scroll",
        sketch_count_seconds: 90,
      },
      now_page: 1,
      duration: 300,
      offset: 0,
      easing: "easeInOutCubic",
      message: {
        user_uuid: "",
        server_status: "",
        now_book_uuid: "",
        read_percent: 0.0,
        msg: "",
      },
    };
  },

  mounted() {
    this.initPage();
    this.$cookies.keys();
  },
  destroyed() {
    //this.$socket.close();
  },

  computed: {
    // 计算属性的 getter
    nowTemplate: function () {
      console.log("computed:" + this.$store.state.setting.template);
      // `this` 指向 vm 实例
      return this.$store.state.setting.template;
    },
  },
  methods: {
    initPage() {
      axios
        .get("/book.json")
        .then((response) => (this.$store.state.book = response.data))
        .finally((this.book = this.$store.book));
      axios
        .get("/setting.json")
        .then((response) => (this.$store.state.setting = response.data))
        .finally((this.setting = this.$store.setting));
      axios
        .get("/bookshelf.json")
        .then((response) => (this.$store.state.bookshelf = response.data))
        .finally((this.bookshelf = this.$store.bookshelf));
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
  text-align: center;
  background-color: #f6f7eb;
  align-items: center;
}
</style>
