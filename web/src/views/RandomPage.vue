<template>
  <div id="RandomPage">
    <Header>
      <h2>
        <a v-if="!book.IsFolder" v-bind:href="'raw/' + book.name"
          >{{ book.name }}【Download】</a
        >
        <a v-if="book.IsFolder" v-bind:href="'raw/' + book.name">{{
          book.name
        }}</a>
      </h2>
    </Header>
    <div class="singe_page_main" v-on:click="nextPage">
      <img
        lazy-src="/resources/favicon.ico"
        v-bind:src="book.pages[page - 1].url"
      /><img />
    </div>
    <v-alert v-model="alert" type="info" close-text="Close Alert" dismissible>
      已经翻到最后一页。
    </v-alert>
    <v-pagination
      circle
      v-model="page"
      :length="book.page_num"
      :total-visible="10"
      @input="toPage"
    >
    </v-pagination>

    <slot></slot>
  </div>
</template>

<style>
#RandomPage {
  margin: 1000px 50px;
  align-items: center;
}

.random_main {
  max-width: 80%;
  max-height: 100%;
  align-items: center;
  display: block;
  margin: auto;
}

img {
  max-width: 80%;
  max-height: 100%;
  display: block;
  margin: auto;
}
</style>

<script>
// import Header from "./Header.vue";

export default {
  components: {
    // Header,
  },

  data() {
    return {
      page: 1,
      time_cont: 0,
      alert: false,
      easing: "easeInOutCubic",
    };
  },

  mounted() {
    this.time_cont = 0;
    // // 增加监听
    // window.addEventListener("keyup", this.handleKeyup);
    // window.addEventListener("scroll", this.handleScroll);
  },
  destroyed() {
    // window.removeEventListener("keyup", this.handleKeyup);
    // window.removeEventListener("scroll", this.handleScroll);
  },
  methods: {
    initPage() {
      this.$cookies.keys();
    },
    nextPage: function (p) {
      if (this.page < this.book.page_num) {
        this.page = this.page + 1;
      } else {
        this.alert = true;
      }
      console.log(p);
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

