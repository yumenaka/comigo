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
      <!-- "image1",page_mark取得。排列在左。page_mark初始值为1.-->
      <!-- 可以是第二张图片，不可以是最后一张图片。 -->
      <!-- page_mark 与 page_mark -1都为单页时才显示。 -->
      <img
        id="image1"
        v-on:click="nextPageSinglePage"
        v-if="
          page_mark < this.$store.state.book.all_page_num &&
          this.$store.state.book.pages[page_mark - 1].image_type ==
            'SinglePage' &&
          this.$store.state.book.pages[page_mark].image_type == 'SinglePage'
        "
        lazy-src="/resources/favicon.ico"
        v-bind:src="this.$store.state.book.pages[page_mark].url"
      /><img />
      <!-- "image2",page_mark-1。排列在右，可以是第一张图片，也可以是最后一张图片。 -->
      <!-- page_mark-1为单页，且page_mark为单页时，这一张共同显示。排列在右，点击后返回上一页。 -->

      <img
        id="image2"
        v-on:click="previousPageSinglePage"
        v-if="
          page_mark - 1 >= 0 &&
          this.$store.state.book.pages[page_mark - 1].image_type ==
            'SinglePage' &&
          this.$store.state.book.pages[page_mark].image_type == 'SinglePage'
        "
        lazy-src="/resources/favicon.ico"
        v-bind:src="this.$store.state.book.pages[page_mark - 1].url"
      /><img />

      <!-- page_mark-1为单页，且page_mark为双页时，这一张单独显示。与上面的图片的唯一不同，点击后翻下一页 -->
      <img
        id="image3"
        v-on:click="nextPageSinglePage"
        v-if="
          page_mark - 1 >= 0 &&
          this.$store.state.book.pages[page_mark - 1].image_type ==
            'SinglePage' &&
          this.$store.state.book.pages[page_mark].image_type == 'DoublePage'
        "
        lazy-src="/resources/favicon.ico"
        v-bind:src="this.$store.state.book.pages[page_mark - 1].url"
      /><img />

      <!-- page_mark-1， page_mark的前一张。-->
      <!-- 如为双叶，只显示这一张，可以是第一张图片  -->
      <img
        id="image4"
        v-on:click="nextPageDoublePage"
        v-if="
          page_mark - 1 >= 0 &&
          this.$store.state.book.pages[page_mark - 1].image_type == 'DoublePage'
        "
        lazy-src="/resources/favicon.ico"
        v-bind:src="this.$store.state.book.pages[page_mark - 1].url"
      /><img />
    </div>
    <v-pagination
      v-if="showPagination"
      v-model="page_mark"
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
        console.log("page_mark error", p);
      }
      this.page_mark = p;
      console.log(p);
    },
    nextPageDoublePage: function () {
      this.page_mark = this.page_mark + 1;
      console.log(this.page_mark);
    },
    //感觉与其弄得这么复杂，不如干脆弄个字典，分别表示页数与图片数组……
    nextPageSinglePage: function () {
      var _this = this;
      //mark指向最后一页+[已经翻完了]，或者小于0、或大于AllPageNum：打印错误值，什么都不做、返回。
      if (_this.page_mark >= _this.AllPageNum || _this.page_mark < 0) {
        console.log(_this.page_mark);
        return;
      }
      //mark指向倒数第1页，页数加1
      if (_this.page_mark == _this.AllPageNum - 1) {
        _this.page_mark = _this.page_mark + 1;
        return;
      }
      //mark指向倒数第2页，页数加1
      if (_this.page_mark == this.AllPageNum - 2) {
        _this.page_mark = _this.page_mark + 1;
        console.log(_this.page_mark);
      }
      //经过上面3轮判断，_this.page_mark < this.AllPageNum - 2
      //继续判断+1或+2
      //如果后两张都是单页，加2页
      if (
        this.$store.state.book.pages[_this.page_mark + 1].image_type ==
          "SinglePage" &&
        this.$store.state.book.pages[_this.page_mark + 2].image_type ==
          "SinglePage"
      ) {
        _this.page_mark = _this.page_mark + 2;
        console.log(_this.page_mark);
        //return;
      } else {
        //剩下的其他情况，比如有一张是双页，都只加1
        _this.page_mark = _this.page_mark + 1;
        console.log(_this.page_mark);
        //return;
      }
    },
    previousPageSinglePage: function () {
      //page_mark不应该小于或等于0，打印并返回。
      if (this.page_mark <= 0) {
        console.log(this.page_mark);
        return;
      }
      //page_mark == 1，因为page_mark 起始值为1，不应该继续减
      if (this.page_mark == 1) {
        console.log(this.page_mark);
        return;
      }
      //page_mark == 2，只能往前翻一页
      if (this.page_mark == 2) {
        this.page_mark = this.page_mark - 1;
        console.log(this.page_mark);
        return;
      }
      //到此为止，this.page_mark >=3
      //-2后，两页都是单页，页数-2
      if (
        this.$store.state.book.pages[this.page_mark - 2].image_type ==
          "SinglePage" &&
        this.$store.state.book.pages[this.page_mark - 2 - 1].image_type ==
          "SinglePage"
      ) {
        this.page_mark = this.page_mark - 2;
        console.log(this.page_mark);
        return;
      }
      this.page_mark = this.page_mark - 1;
      console.log(this.page_mark);
      return;
    },
    //键盘快捷键用的下一页、上一页函数
    nextPage: function () {
      if (this.page_mark > this.$store.state.book.all_page_num) {
        console.log(this.page_mark);
        return;
      }
      if (this.page_mark == this.$store.state.book.all_page_num) {
        console.log(this.page_mark);
        return;
      }
      if (
        this.$store.state.book.pages[this.page_mark].image_type == "SinglePage"
      ) {
        this.nextPageSinglePage();
        return;
      }
      if (
        this.$store.state.book.pages[this.page_mark].image_type == "DoublePage"
      ) {
        this.nextPageDoublePage();
        return;
      }
    },
    previousPage: function () {
      if (this.page_mark > this.$store.state.book.all_page_num) {
        console.log(this.page_mark);
        return;
      }
      if (this.page_mark == this.$store.state.book.all_page_num) {
        this.page_mark = this.page_mark - 1;
        return;
      }
      if (
        this.$store.state.book.pages[this.page_mark].image_type == "SinglePage"
      ) {
        this.previousPageSinglePage();
      } else if (
        this.$store.state.book.pages[this.page_mark].image_type ==
          "DoublePage" &&
        this.page_mark - 1 >= 0
      ) {
        this.page_mark = this.page_mark - 1;
        console.log(this.page_mark);
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
      page_mark: 1,
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
    // _this.page_mark=1;
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
