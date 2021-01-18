<template>
  <div id="app" class="app_div">
    <!-- 下拉阅读 -->
    <MultiPage
      :book="book"
      :bookshelf="bookshelf"
      :defaultSetiing="defaultSetiing"
      v-if="defaultSetiing.default_template === 'multi'"
    >
    </MultiPage>

    <!-- 随机,或倒计时（绘图用） -->
    <RandomPage
      :book="book"
      :bookshelf="bookshelf"
      :defaultSetiing="defaultSetiing"
      v-if="defaultSetiing.default_template === 'random'"
    >
    </RandomPage>

    <!-- 单页阅读 -->
    <SinglePage
      :book="book"
      :bookshelf="bookshelf"
      :defaultSetiing="defaultSetiing"
      v-if="defaultSetiing.default_template === 'single'"
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
      book: this.$store.getters.book,
      bookshelf: this.$store.getters.bookshelf,
      defaultSetiing: this.$store.getters.bookshelf,
      page: this.$store.getters.now_page,
      duration: 300,
      offset: 0,
      easing: "easeInOutCubic",
      message: this.$store.getters.message,
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
      this.$store.commit('syncRemoteSetting');
      this.$store.commit('syncBookDate');
      this.$store.commit('syncBookShelfDate');
      
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
