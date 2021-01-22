<template>
  <div id="SinglePage">
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
    <div class="singe_page" v-on:click="nextPage">
      <img
        lazy-src="/resources/favicon.ico"
        v-bind:src="book.pages[page - 1].url"
      /><img />
    </div>

    <slot></slot>
  </div>
</template>

<style>
.singe_page {
  width: 100%;
  height: calc(100vh - 100px);
  border: 2px solid rgb(84, 106, 233);
  display: table-cell;
  vertical-align: middle;
}

.singe_page img {
  width: 100%;
  height: calc(100vh - 120px);
  /* display: block; */
  margin: center;
  vertical-align: middle;
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
      defaultSetiing: null,
    };
  },

  mounted() {
      this.book = this.$store.state.book;
      this.bookshelf = this.$store.state.bookshelf;
      this.defaultSetiing = this.$store.state.defaultSetiing;
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
      this.$cookies.keys();
    },
    nextPage: function (p) {
      if (this.page < this.book.page_num) {
        this.page = this.page + 1;
      } else {
        // this.alert = true;
        alert("Last Page!");
      }
      console.log(p);
    },
    toPage: function (p) {
      this.page = p;
      console.log(p);
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
