<template>
  <div id="app" class="app_div">
    <!-- 下拉阅读 -->
    <MultiPage
      :book="book"
      :bookshelf="bookshelf"
      :defaultSetiing="defaultSetiing"
      v-if="defaultSetiing.template === 'multi'"
    >
    </MultiPage>

    <!-- 随机,或倒计时（绘图用） -->
    <RandomPage
      :book="book"
      :bookshelf="bookshelf"
      :defaultSetiing="defaultSetiing"
      v-if="defaultSetiing.template === 'random'"
    >
    </RandomPage>

    <!-- 单页阅读 -->
    <SinglePage
      :book="book"
      :bookshelf="bookshelf"
      :defaultSetiing="defaultSetiing"
      v-if="defaultSetiing.template === 'single'"
    >
    </SinglePage>
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
  components: {
    MultiPage,
    SinglePage,
    RandomPage,
  },
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
      bookshelf: {},
      defaultSetiing: {
        default_page_template:"???",
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
      //this.book = this.$store.book
      // this.$store.commit('syncBookDate');
      axios.get("/book.json").then((response) => (this.book = response.data));
      axios
        .get("/setting.json")
        .then((response) => (this.defaultSetiing = response.data));
      axios
        .get("/bookshelf.json")
        .then((response) => (this.bookshelf = response.data))
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

.app_div {
  /* margin: auto; */
  align-items: center;
}
</style>
