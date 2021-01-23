<template>
  <div id="SketchPage" class="SketchPage">
    <Header>
      <h2>
        <a v-if="!this.$store.state.book.IsFolder" v-bind:href="'raw/' + this.$store.state.book.name"
          >{{ this.$store.state.book.name }}【Download】</a
        >
        <a v-if="this.$store.state.book.IsFolder" v-bind:href="'raw/' + this.$store.state.book.name">{{
          book.name
        }}</a>
      </h2>
    </Header>
    <div class="sketch_div" v-on:click="nextPage(2)">
      <img
        lazy-src="/resources/favicon.ico"
        v-bind:src="this.$store.state.book.pages[page - 1].url"
      /><img />
    </div>
    <v-pagination
      circle
      v-model="page"
      :length="this.$store.state.book.page_num"
      :total-visible="10"
      @input="toPage"
    >
    </v-pagination>
    <slot></slot>
  </div>
</template>

<style>
#RandomPage {
  align-items: center;
  width: 100vw;
  height: 100vh;
  align-self: center;
}

.sketch_div {
  width: 100%;
  height: 80vh;
  align-items: center;
}

.sketch_div img {
  max-width: 100%;
  max-height: 100%;
  height: 80vh;
  display: block;
  margin: center;
}
</style>

<script>
import Header from "./Header.vue";

export default {
  components: {
    Header,
  },

  data() {
    return {
      book: null,
      bookshelf: null,
      defaultSetting: null,
      page: 1,
      time_cont: 0,
      alert: false,
      easing: "easeInOutCubic",
    };
  },

  mounted() {
    this.time_cont = 0;
    this.$cookies.keys();
    this.book = this.$store.state.book;
    this.bookshelf = this.$store.state.bookshelf;
    this.defaultSetting = this.$store.state.defaultSetting;
    // this.initPage();
  },

  destroyed() {},

  methods: {
    initPage() {},
    nextPage: function (num) {
      if (this.page < this.book.page_num) {
        this.page = this.page + num;
      } else {
        this.alert = true;
      }
      console.log(num);
    },
    toPage: function (p) {
      this.page = p;
      console.log(p);
    },
    // 键盘事件
    handleKeyup(event) {
      const e = event || window.event || arguments.callee.caller.arguments[0];
      if (!e) return;
      const { key, keyCode } = e;
      console.log(keyCode);
      console.log(key);
    },
    //  滑轮事件
    handleScroll() {
      var e = document.body.scrollTop || document.documentElement.scrollTop;
      if (!e) return;
      // console.log(e);
    },
  },
};
</script>

