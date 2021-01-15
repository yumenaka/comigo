<template>
  <div id="app" class="app_div">
    <!-- 下拉阅读 -->
    <MultiPage v-if="defaultSetiing.default_page_mode === 'multi'"> </MultiPage>

    <!-- 随机,或倒计时（绘图用） -->
    <RandomPage v-if="defaultSetiing.default_page_mode === 'random'">
    </RandomPage>

    <!-- 单页阅读 -->
    <SinglePage v-if="defaultSetiing.default_page_mode === 'single'">
    </SinglePage>
  </div>
</template>

<script>
//代码参考：https://github.com/bradtraversy/vue_crash_todolist
import axios from "axios";
// import Header from "./views/Header.vue";
import MultiPage from "./views/MultiPage.vue";
import SinglePage from "./views/SinglePage.vue";

export default {
  name: "app",
  components: {
    // Header,
    MultiPage,
    SinglePage,
  },
  data() {
    return {
      book: {
        name: "null",
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
      bookshelf: {},
      defaultSetiing: {
        default_page_mode: "single",
      },
      page: 1,
      page_mode: "multi",
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
    axios.get("/book.json").then((response) => (this.book = response.data));
    axios
      .get("/setting.json")
      .then((response) => (this.defaultSetiing = response.data));
    axios
      .get("/bookshelf.json")
      .then((response) => (this.bookshelf = response.data))
      .finally();
  },
  destroyed() {
    this.$socket.close();
  },
  methods: {
    initPage() {
      this.$cookies.keys();
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
}

.app_div {
  margin: auto;
  align-items: center;
}
</style>
