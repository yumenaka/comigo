<template>
  <div id="DoublePageTemplate">
    <Header v-if="showHeader">
      <h2>
        <a
          v-if="!this.$store.state.book.IsFolder"
          v-bind:href="'raw/' + this.$store.state.book.name"
          >{{ this.$store.state.book.name }}【Download】现在时刻：{{
            currentTime
          }}</a
        >
        <a
          v-if="this.$store.state.book.IsFolder"
          v-bind:href="'raw/' + this.$store.state.book.name"
          >{{ this.$store.state.book.name }}现在时刻：{{ currentTime }}</a
        >
      </h2>
    </Header>
    <div class="double_page_main">
      <img
        v-on:click="nextPageSinglePage"
        v-if="
          now_page < this.$store.state.book.all_page_num &&
          this.$store.state.book.pages[now_page - 1].image_type ==
            'SinglePage' &&
          this.$store.state.book.pages[now_page].image_type == 'SinglePage'
        "
        lazy-src="/resources/favicon.ico"
        v-bind:src="this.$store.state.book.pages[now_page].url"
      /><img />
      <img
        v-on:click="previousPageSinglePage"
        v-if="
          now_page - 1 >= 0 &&
          this.$store.state.book.pages[now_page - 1].image_type == 'SinglePage'
        "
        lazy-src="/resources/favicon.ico"
        v-bind:src="this.$store.state.book.pages[now_page - 1].url"
      /><img />

      <img
        v-on:click="nextPageDoublePage"
        v-if="
          now_page - 1 >= 0 &&
          this.$store.state.book.pages[now_page - 1].image_type == 'DoublePage'
        "
        lazy-src="/resources/favicon.ico"
        v-bind:src="this.$store.state.book.pages[now_page - 1].url"
      /><img />
    </div>
    <v-pagination
      v-if="showPagination"
      v-model="now_page"
      :length="this.$store.state.book.all_page_num"
      :total-visible="15"
      @input="toPage"
    >
    </v-pagination>
    <slot></slot>
  </div>
</template>

<script>
import Header from "./Header.vue";

export default {
  methods: {
    initPage() {},
    toPage: function (p) {
      if (p > this.$store.state.book.all_page_num && p < 0) {
        console.log("now_page error", p);
      }
      this.now_page = p;
      console.log(p);
    },

    nextPageDoublePage: function () {
      this.now_page = this.now_page + 1;
      console.log(this.now_page);
    },

    nextPageSinglePage: function () {
      var _this = this;
      //已经是最后一页，什么都不做，返回
      if (_this.now_page + 1 > _this.AllPageNum) {
        console.log(_this.now_page);
        return;
      }
      //还差一页是最后一页，页数只能加1
      if (_this.now_page + 1 == _this.AllPageNum) {
        _this.now_page = _this.now_page + 1;
        return;
      }
      //经过上面两轮判断，当前不是倒数前两页，+1或+2
      if (_this.now_page + 1 < _this.AllPageNum) {
        //如果目前是两张单页，加2页
        if (
          this.$store.state.book.pages[_this.now_page].image_type ==
            "SinglePage" &&
          this.$store.state.book.pages[_this.now_page - 1].image_type ==
            "SinglePage"
        ) {
          _this.now_page = _this.now_page + 2;
          console.log(_this.now_page);
          return;
        } else {
          _this.now_page = _this.now_page + 1;
          console.log(_this.now_page);
          return;
        }
      } else {
        _this.now_page = _this.now_page + 1;
        console.log(_this.now_page);
        return;
      }
    },
    previousPageSinglePage: function () {
      if (this.now_page == 0) {
        console.log(this.now_page);
        return;
      }
      if (this.now_page == 1) {
        this.now_page = 1;
        console.log(this.now_page);
        return;
      }
      if (this.now_page == 2) {
        this.now_page = this.now_page - 1;
        console.log(this.now_page);
        return;
      }
      //当前页与下一页都是单页，页数-2
      if (
        this.$store.state.book.pages[this.now_page - 2].image_type ==
          "SinglePage" &&
        this.$store.state.book.pages[this.now_page - 1].image_type ==
          "SinglePage"
      ) {
        this.now_page = this.now_page - 2;
        console.log(this.now_page);
        return;
      } else {
        this.now_page = this.now_page - 1;
        console.log(this.now_page);
        return;
      }
    },
    //键盘快捷键用的下一页、上一页函数
    nextPage: function () {
      if (this.now_page > this.$store.state.book.all_page_num) {
        console.log(this.now_page);
        return;
      }
      if (this.now_page == this.$store.state.book.all_page_num) {
        console.log(this.now_page);
        return;
      }
      if (
        this.$store.state.book.pages[this.now_page].image_type == "SinglePage"
      ) {
        this.nextPageSinglePage();
        return;
      }
      if (
        this.$store.state.book.pages[this.now_page].image_type == "DoublePage"
      ) {
        this.nextPageDoublePage();
        return;
      }
    },
    previousPage: function () {
      if (this.now_page > this.$store.state.book.all_page_num) {
        console.log(this.now_page);
        return;
      }
      if (this.now_page == this.$store.state.book.all_page_num) {
        this.now_page = this.now_page - 1;
        return;
      }
      if (
        this.$store.state.book.pages[this.now_page].image_type == "SinglePage"
      ) {
        this.previousPageSinglePage();
      } else if (
        this.$store.state.book.pages[this.now_page].image_type ==
          "DoublePage" &&
        this.now_page - 1 >= 0
      ) {
        this.now_page = this.now_page - 1;
        console.log(this.now_page);
      }
    },

    // 键盘事件
    handleKeyup(event) {
      const e = event || window.event || arguments.callee.caller.arguments[0];
      if (!e) return;
      //https://developer.mozilla.org/zh-CN/docs/Web/API/KeyboardEvent/keyCode
      switch (e.key) {
        // case "KeyH":
        case "PageUp":
        case "ArrowUp":
        case "ArrowLeft":
          this.previousPage(); //前一页
          break;
        // case "KeyL":
        case "Space":
        case "ArrowDown":
        case "PageDown":
        case "ArrowRight":
          this.nextPage(); //后一页
          break;
        // case "KeyJ":  
        case "Home":  
          this.toPage(1); //跳转到第一页
          break;
        // case "KeyK":
        case "End":  
          this.toPage(this.$store.state.book.all_page_num);
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

  components: {
    Header,
  },

  data() {
    return {
      book: null,
      bookshelf: null,
      defaultSetting: null,
      now_page: 1,
      showHeader: false,
      showPagination: true,
      AllPageNum: this.$store.state.book.all_page_num,
      time_cont: 0,
      alert: false,
      easing: "easeInOutCubic",
      timer: "", //定义一个定时器的变量
      currentTime: new Date(), // 获取当前时间
    };
  },
  created() {
    var _this = this; //声明一个变量指向Vue实例this，保证作用域一致
    // _this.now_page=1;
    this.timer = setInterval(function () {
      var date = new Date();
      var year = date.getFullYear();
      var month = date.getMonth() + 1;
      var strDate = date.getDate();
      if (month >= 1 && month <= 9) {
        month = "0" + month;
      }
      if (strDate >= 0 && strDate <= 9) {
        strDate = "0" + strDate;
      }
      var currentdate = year + " 年 " + month + " 月 " + strDate + " 日 ";
      var Hours = date.getHours();
      if (Hours >= 0 && Hours <= 9) {
        Hours = "0" + Hours;
      }
      var Minutes = date.getMinutes();
      if (Minutes >= 0 && Minutes <= 9) {
        Minutes = "0" + Minutes;
      }
      var Seconds = date.getSeconds();
      if (Seconds >= 0 && Seconds <= 9) {
        Seconds = "0" + Seconds;
      }
      _this.currentTime =
        currentdate + " " + Hours + ":" + Minutes + ":" + Seconds;
    }, 1000);
  },
  mounted() {
    this.time_cont = 0;
    this.$cookies.keys();
    // 增加监听
    window.addEventListener("keyup", this.handleKeyup);
    // window.addEventListener("scroll", this.handleScroll);
  },

  destroyed() {
    // 销毁监听
    window.removeEventListener("keyup", this.handleKeyup);
    // window.removeEventListener("scroll", this.handleScroll);
    if (this.timer) {
      clearInterval(this.timer); // 在Vue实例销毁前，清除我们的定时器
    }
  },
};
</script>


<style>
#DoublePageTemplate {
  align-items: center;
  width: 100vw;
  height: 100vh;
  align-self: center;
}

.double_page_main {
  width: 100%;
  height: 95vh;
  display: flex;
  justify-content: center;
  align-items: center;
}

.double_page_main img {
  max-width: 100%;
  max-height: 100%;
  height: 95vh;
  display: block;
  margin: center;
}
</style>
