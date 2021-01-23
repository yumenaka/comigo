<template>
  <div id="SinglePage">
    <Header>
      <h2>
        <a
          v-if="!this.$store.state.book.IsFolder"
          v-bind:href="'raw/' + this.$store.state.book.name"
          >{{ this.$store.state.book.name }}【Download】</a
        >
        <a
          v-if="this.$store.state.book.IsFolder"
          v-bind:href="'raw/' + this.$store.state.book.name"
          >{{ this.$store.state.book.name }}</a
        >
      </h2>
    </Header>
    <div class="sindle_div" v-on:click="nextPage(2)">
      <p>现在时刻：{{}}</p>
      <img
        v-if="page < this.$store.state.book.page_num"
        lazy-src="/resources/favicon.ico"
        v-bind:src="this.$store.state.book.pages[page].url"
      /><img />
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
#SinglePage {
  /* width: 100%;
  height: calc(100vh - 100px);
  border: 2px solid rgb(84, 106, 233);
  display: table-cell;
  vertical-align: middle; */
  align-items: center;
  width: 100vw;
  height: 100vh;
  align-self: center;
}

.singe_giv {
  width: 100%;
  height: 80vh;
  display: flex;
  justify-content: center;
  align-items: center;
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
      page: 1,
      alert: false,
      easing: "easeInOutCubic",
      book: null,
      bookshelf: null,
      defaultSetting: null,
    };
  },

  mounted() {
      this.book = this.$store.state.book;
      this.bookshelf = this.$store.state.bookshelf;
      this.defaultSetting = this.$store.state.defaultSetting;
    // 增加监听
    // window.addEventListener("keyup", this.handleKeyup);
    // window.addEventListener("scroll", this.handleScroll);
  },
  destroyed() {
    // window.removeEventListener("keyup", this.handleKeyup);
    // window.removeEventListener("scroll", this.handleScroll);
  },

  methods: {
    initPage() {
      //this.$cookies.keys();
    },
    nextPage: function (num) {
      if (this.page < this.book.page_num) {
        this.page = this.page + num;
      } else {
        this.alert = true;
      }
      console.log(num);
    },
    toPage: function (num) {
      this.page = num;
      console.log(num);
    },
    moveSomething(e) {
      switch (e.keyCode) {
        case 37:
          // left key pressed
          //advancePage(-1);
          break;
        case 32:
          // spacebar pressed
          if (
            window.innerHeight + window.scrollY >=
            document.body.offsetHeight
          ) {
            //advancePage(1);
          }
          break;
        case 39:
          // right key pressed
          //advancePage(1);
          break;
        case 17:
          // Ctrl key pressed
          //openOverlay();
          break;
      }
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
