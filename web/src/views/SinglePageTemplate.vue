<template>
  <div id="SinglePageTemplate">
    <div class="single_page_main">
      <img
        v-on:click="addPage(1)"
        v-if="now_page <= this.$store.state.book.all_page_num && now_page >= 1"
        lazy-src="/resources/favicon.ico"
        v-bind:src="this.$store.state.book.pages[now_page - 1].url"
      /><img />
    </div>
    <v-pagination
      v-model="now_page"
      :length="this.$store.state.book.all_page_num"
      :total-visible="10"
      @input="toPage"
    >
    </v-pagination>
    <slot></slot>
  </div>
</template>

<style>
/* #SinglePageTemplate {
  align-items: center;
  width: 100vw;
  height: 100vh;
  align-self: center;
} */

.single_page_main {
  width: 100%;
  height: 95vh;
  display: flex;
  justify-content: center;
  align-items: center;
}

.single_page_main img {
  max-width: 100%;
  max-height: 100%;
  display: block;
  margin: center;
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
      now_page: 1,
      showPagination: true,
      alert: false,
      easing: "easeInOutCubic",
      book: null,
      bookshelf: null,
      setting: null,
    };
  },

  mounted() {
    // 注册监听
    window.addEventListener("keyup", this.handleKeyup);
    // window.addEventListener("scroll", this.handleScroll);
  },
  destroyed() {
    // 销毁监听
    window.removeEventListener("keyup", this.handleKeyup);
    // window.removeEventListener("scroll", this.handleScroll);
  },

  methods: {
    initPage() {
      //this.$cookies.keys();
    },
    addPage: function (num) {
      if (
        this.now_page + num <= this.$store.state.book.all_page_num &&
        this.now_page + num >= 1
      ) {
        this.now_page = this.now_page + num;
      }
      // console.log(this.now_page);
    },
    toPage: function (num) {
      if (num <= this.$store.state.book.all_page_num && num >= 1) {
        this.now_page = num;
      }
      // console.log(num);
    },
    // 键盘事件
    handleKeyup(event) {
      const e = event || window.event || arguments.callee.caller.arguments[0];
      if (!e) return;
      //https://developer.mozilla.org/zh-CN/docs/Web/API/KeyboardEvent/keyCode
      switch (e.key) {
        case "ArrowUp":
        case "PageUp":
        case "ArrowLeft":
          this.addPage(-1); //上一页
          break;
        case "Space":
        case "ArrowDown":
        case "PageDown":
        case "ArrowRight":
          this.addPage(1); //下一页
          break;
        case "Home":  
          this.toPage(1); //跳转到第一页
          break;
        case "End":  
          this.toPage(this.$store.state.book.all_page_num); //跳转到最后一页
          break;
        case "Ctrl":
          // Ctrl key pressed //组合键？
          //openOverlay();
          break;
      }
      // console.log(e.keyCode);
      // console.log(e.key);
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
