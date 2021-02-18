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
      <!-- [page_mark]单页+双页:排列在左（两张都是单页）或中间（这一张为单页，下一张双页||这一张为双页）。-->
      <!-- 上面三种情况，图片点击事件都是下一页 -->
      <!-- page_mark初始值为0或1（前两张为单页，初始化为1）,最大值为this.$store.state.book.all_page_num，等于最大值时本image不显示  -->
      <!-- 可以是第1张图片，不可以是最后的图片。 -->
      <!--  page_mark < this.$store.state.book.all_page_num 时显示。 -->
      <img
        id="image1"
        v-on:click="nextPageClick"
        v-if="page_mark < this.$store.state.book.all_page_num"
        lazy-src="/resources/favicon.ico"
        v-bind:src="this.$store.state.book.pages[page_mark].url"
      /><img />
      <!-- [page_mark - 1]单页的情况:排列在右，可以是第2张图片，也可以是最后一张图片。 -->
      <!-- page_mark为单页的前提下，page_mark-1为单页时，这一张作为右页共同显示。点击后返回上一页。 -->
      <!-- page_mark为双页，page_mark-1无需显示。 -->
      <img
        id="image2"
        v-on:click="previousPageClick"
        v-if="
          page_mark - 1 >= 0 &&
          page_mark < this.$store.state.book.all_page_num &&
          this.$store.state.book.pages[page_mark].image_type == 'SinglePage' &&
          this.$store.state.book.pages[page_mark - 1].image_type == 'SinglePage'
        "
        lazy-src="/resources/favicon.ico"
        v-bind:src="this.$store.state.book.pages[page_mark - 1].url"
      /><img />
    </div>
    <v-pagination
      v-if="showPagination"
      v-model="page_mark"
      :length="this.$store.state.book.all_page_num - 1"
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
  components: {
    Header,
  },

  data() {
    return {
      showHeader: true,
      localbook: {
        name: this.$store.state.book.name,
        all_image_num: this.$store.state.book.all_page_num,
        images: this.$store.state.book.pages,
        pages: null, //需要根据all_image_num、images计算
        all_page_num: 0, //需要根据all_image_num、images计算
      },
      bookshelf: null,
      setting: null,
      page_mark: 0, //初始值为0或1（根据单双页判断，initPageMark）,最大值为this.$store.state.book.all_page_num 最大值的时候，代码逻辑上需要一些特殊处理（page_mark数组越界，但page_mark-1依然有意义）。
      
      showPagination: true,
      AllPageNum: this.$store.state.book.all_page_num - 1,
      time_cont: 0,
      alert: false,
      easing: "easeInOutCubic",
      timer: "", //定义一个定时器的变量
      currentTime: new Date(), // 获取当前时间
    };
  },

  methods: {
    //与其弄得这么复杂，不如干脆声称本地book数据，重新设定页数与对应图片数组……
    initLocalBook() {},
    initPageMark() {
      if (this.$store.state.book.all_page_num < 2) {
        this.page_mark = 0;
        return;
      }
      if (
        this.$store.state.book.pages[0].image_type == "SinglePage" &&
        this.$store.state.book.pages[1].image_type == "SinglePage"
      ) {
        this.page_mark = 1; //前两张为单页，初始化为1
      } else {
        this.page_mark = 0; //前两张有一张为双页，初始化为0
      }
    },
    toPage: function (p) {
      if (p > this.$store.state.book.all_page_num || p < 0) {
        console.log("page_mark error", p);
      }
      this.page_mark = p;
      console.log(p);
    },

    nextPageClick: function () {
      //var _this = this;
      //错误处理：mark指向最后一页[已经翻完了]，或大于AllPageNum，或小于0：打印错误值，什么都不做、返回。
      if (this.page_mark >= this.AllPageNum || this.page_mark < 0) {
        console.log(this.page_mark);
        return;
      }
      //特殊处理：mark指向倒数第1页，页数加1
      if (this.page_mark == this.AllPageNum - 1) {
        this.page_mark = this.page_mark + 1;
        return;
      }
      //特殊处理：mark指向倒数第2页，如果最后两张全部为单页，页数加2，否则加1
      if (this.page_mark == this.AllPageNum - 2) {
        //this.page_mark +2 == this.AllPageNum
        if (
          this.$store.state.book.pages[this.AllPageNum - 1].image_type ==
            "SinglePage" &&
          this.$store.state.book.pages[this.AllPageNum - 2].image_type ==
            "SinglePage"
        ) {
          this.page_mark = this.page_mark + 2;
          console.log(this.page_mark);
        } else {
          this.page_mark = this.page_mark + 1;
          console.log(this.page_mark);
        }
        return;
      }
      //经过上面3轮判断，_this.page_mark < this.AllPageNum - 2
      //继续判断+1或+2：如果接下来两张都是单页，加2页。其他情况加1页
      if (
        this.$store.state.book.pages[this.page_mark + 1].image_type ==
          "SinglePage" &&
        this.$store.state.book.pages[this.page_mark + 2].image_type ==
          "SinglePage"
      ) {
        this.page_mark = this.page_mark + 2;
        console.log(this.page_mark);
        //return;
      } else {
        //剩下的其他情况，比如有一张是双页，都只加1
        this.page_mark = this.page_mark + 1;
        console.log(this.page_mark);
        //return;
      }
    },
    previousPageClick: function () {
      //page_mark不应该小于0，等于0也无法继续往后翻，打印并返回。
      if (this.page_mark <= 0) {
        console.log(this.page_mark);
        return;
      }
      //page_mark == 1，最多只能再减一页
      if (this.page_mark == 1) {
        this.page_mark = this.page_mark - 1;
        console.log(this.page_mark);
        return;
      }
      //page_mark == 2，只能往前翻一页
      if (this.page_mark == 2) {
        this.page_mark = this.page_mark - 1;
        console.log(this.page_mark);
        return;
      }
      //错误处理：mark指向最后一页[已经翻完了]，或大于AllPageNum，倒数第一页
      if (this.page_mark >= this.AllPageNum) {
        this.page_mark = this.AllPageNum - 1;
        console.log(this.page_mark);
        return;
      }
      // //特殊处理：mark指向倒数第1页，页数加1
      // if (_this.page_mark == _this.AllPageNum - 1) {
      //   _this.page_mark = _this.page_mark - 1;
      //   return;
      // }

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
      //其他情况，都只往前翻一页
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
      this.nextPageClick();
      // console.log(this.page_mark);
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
      this.previousPageClick();
      // console.log(this.page_mark);
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
          this.toPage(this.$store.state.book.all_page_num - 1);
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
    this.initPageMark();
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
#DoublePage_Pagination {
  color: #066eb4;
  background-color: #f6f7eb;
  align-items: center;
}

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
