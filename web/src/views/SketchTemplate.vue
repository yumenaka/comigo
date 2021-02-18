<template>
  <div id="SketchPage">
    <Header v-if="showHeader">
      <h2>
        <a v-bind:href="'raw/' + this.$store.state.book.name"
          >{{ this.$store.state.book.name }}现在时刻：{{ currentTime }}</a
        >
      </h2>
    </Header>

    <div class="sketch_main">
      <div id="SketchHint">
        <p>{{ getNowCount() }}/{{ getALLSeconds()}}⏳</p>
      </div>
      <img
        v-on:click="addPage(1)"
        lazy-src="/resources/favicon.ico"
        v-bind:src="this.$store.state.book.pages[now_page - 1].url"
      /><img />
      <div id="SketchHint">
        <p>🕒{{ currentTime }}</p>
      </div>
    </div>
    <v-pagination
      v-if="showPagination"
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
#SketchPage {
  align-items: center;
  width: 100vw;
  height: 100vh;
  align-self: center;
}

#SketchHint {
  font-family: "Josefin Sans", -apple-system, "PingFang SC", "Hiragino Sans GB",
    "Microsoft Yahei", "WenQuanYi Micro Hei", "ST Heiti", sans-serif;
  color: #066eb4;
  text-align: center;
  /* position: absolute;
  left: 50%;
  top: 50%;
  -webkit-transform: translate(-50%, -50%);
  transform: translate(-50%, -50%);
  color: #daf6ff; */
  text-shadow: 0 0 20px #0aafe6, 0 0 20px rgba(10, 175, 230, 0);
  letter-spacing: 0.05em;
  font-size: 20px;
  padding: 5px 0;
}

.sketch_main {
  width: 100%;
  height: 95vh;
  display: flex;
  /* space-between space-evenly space-around */
  justify-content: center;
  align-items: center;
}

.sketch_main div {
  display: flex;
  /* space-between space-evenly space-around */
  justify-content: center;
}

.sketch_main img {
  max-width: 100%;
  max-height: 100%;
  /* height: 95vh; */
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
      showHeader: true,
      time_cont: 1,
      WaitSeconds: this.$store.state.setting.sketch_count_seconds,
      book: null,
      bookshelf: null,
      setting: null,
      showPagination: true,
      now_page: 1,
      alert: false,
      easing: "easeInOutCubic",
      timer: "", //定义一个定时器的变量
      currentTime: null, // 获取当前时间
    };
  },
  created() {
    var _this = this; //声明一个变量指向Vue实例this，保证作用域一致
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
      var currentdate = year + "年" + month + "月" + strDate + "日";
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
      //_this.currentTime =currentdate + " " + Hours + ":" + Minutes + ":" + Seconds;
      _this.currentTime = Hours + ":" + Minutes + ":" + Seconds;
      //每 WaitSeconds 秒翻页
      console.log(currentdate + "time_cont：" + _this.time_cont);
      if (_this.time_cont < _this.WaitSeconds) {
        _this.time_cont++;
      } else {
        _this.time_cont = 0;
        console.log("时间到，翻页：" + _this.currentTime + "秒");
        if (_this.now_page < _this.$store.state.book.all_page_num) {
          _this.now_page += 1;
        } else {
          _this.now_page = 1;
        }
      }
    }, 1000);
  },
  mounted() {
    this.time_cont = 0;
    this.$cookies.keys();
    // 注册监听
    window.addEventListener("keyup", this.handleKeyup);
    // window.addEventListener("scroll", this.handleScroll);
  },

  destroyed() {
    // 销毁监听
    window.removeEventListener("keyup", this.handleKeyup);
    // window.removeEventListener("scroll", this.handleScroll);
    if (this.timer) {
      clearInterval(this.timer); // 在Vue实例销毁前，清除定时器
    }
  },

  methods: {
    initPage() {},
    getWaitSeconds() {
      //console.log(this.$store.state.setting)
      return this.$store.state.setting.sketch_count_seconds;
    },
    getNowCount() {
      var Seconds =
        this.$store.state.setting.sketch_count_seconds - this.time_cont;
      if (Seconds >= 0 && Seconds <= 9) {
        Seconds = "0" + Seconds;
      }

      return Seconds;
    },
    getALLSeconds() {
      var AllSeconds = this.$store.state.setting.sketch_count_seconds;
      if (AllSeconds >= 0 && AllSeconds <= 9) {
        AllSeconds = "0" + AllSeconds;
      }
      return AllSeconds;
    },
    addPage: function (num) {
      if (
        this.now_page + num < this.$store.state.book.all_page_num &&
        this.now_page + num >= 1
      ) {
        this.now_page = this.now_page + num;
      }
      // console.log(this.now_page);
    },
    toPage: function (p) {
      this.now_page = p;
      console.log(p);
    },
    // 键盘事件
    handleKeyup(event) {
      const e = event || window.event || arguments.callee.caller.arguments[0];
      if (!e) return;
      //https://developer.mozilla.org/zh-CN/docs/Web/API/KeyboardEvent/keyCode
      switch (e.key) {
        case "PageUp":
        case "ArrowUp":
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
          this.toPage(this.$store.state.book.all_page_num - 1); //跳转到最后一页
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

