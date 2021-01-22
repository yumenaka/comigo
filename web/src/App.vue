<template>
  <div id="app">
    <!-- 下拉阅读 -->
    <MultiPage v-if="defaultSetiing.template === 'multi--'"> </MultiPage>

    <!-- 随机,或倒计时（绘图用） -->
    <RandomPage v-if="defaultSetiing.template === 'multi'"> </RandomPage>

    <!-- 单页阅读 -->
    <SinglePage v-if="defaultSetiing.template === 'single'"> </SinglePage>
  </div>
</template>

<script>
//代码参考：https://github.com/bradtraversy/vue_crash_todolist
import axios from "axios";
import MultiPage from "./views/MultiPage.vue";
import SinglePage from "./views/SinglePage.vue";
import RandomPage from "./views/RandomPage.vue";

export default {
  name: "app",
  //为了能在模板中使用，组件必须先注册以便 Vue 能够识别
  components: {
    MultiPage,
    SinglePage,
    RandomPage,
  },
  //组件的 data 选项必须是一个函数
  //每个实例可以维护一份被返回对象的独立的拷贝
  data() {
    return {
      book: {
        name: "loading",
        page_num: 1,
        pages: [
          {
            height: 2000,
            width: 1419,
            url: "/resources/favicon.ico",
            class: "Vertical",
          },
        ],
      },
      //如果你知道你会在晚些时候需要一个 property，但是一开始它为空或不存在，那么你仅需要设置一些初始值。
      bookshelf: {},
      defaultSetiing: {
        default_page_template: "???",
      },
      page: 1,
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
  },
  destroyed() {
    this.$socket.close();
  },
  methods: {
    initPage() {
      this.$cookies.keys();
      this.book = this.$store.book;
      this.defaultSetiing = this.$store.defaultSetiing;
      this.bookshelf = this.$store.bookshelf;
      // this.$store.commit('syncBookDate');
      axios
        .get("/book.json")
        .then((response) => (this.$store.state.book = response.data));
      axios
        .get("/setting.json")
        .then((response) => (this.defaultSetiing = response.data));
      axios
        .get("/bookshelf.json")
        .then((response) => (this.$store.state.bookshelf = response.data))
        .finally();
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
