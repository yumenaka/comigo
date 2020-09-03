<template>
  <div id="app" class="manga_div">
    <h2>
      <a v-bind:href="'raw/' + book.name">{{ book.name }}</a>
    </h2>
    <h4>总页数：{{ book.page_num }}</h4>
    <!-- <v-alert type="success" fixed>I'm a success alert.</v-alert> -->
    <!-- <div>bookshelf:{{bookshelf}}</div> -->
    <div v-for="b in bookshelf" :key="b.uuid" class="bookshelf">
      <v-btn @click="onChangeBook(b.uuid)" class="book_button">{{ b.name }} ({{ b.page_num }})</v-btn>
    </div>
    <div v-for="page in book.pages" :key="page.num" class="manga">
      <img
        v-lazy="page.url"
        v-bind:H="page.height"
        v-bind:W="page.width"
        v-bind:class="page.class | capitalize(page.url)"
      />
    </div>
    <p></p>
    <v-btn v-scroll="onScroll" v-show="btnFlag" fab color="#bbcbff" bottom right @click="toTop">▲</v-btn>
  </div>
</template>

<script>
//代码参考：https://github.com/bradtraversy/vue_crash_todolist
import axios from "axios";
// import * as easings from 'vuetify/es5/services/goto/easing-patterns';
export default {
  name: "app",
  components: {
    //WebSocketTest
    // Header
  },
  data() {
    return {
      book: null,
      bookshelf: null,
      btnFlag: false,
      duration: 300,
      offset: 0,
      easing: "easeInOutCubic",
      message: {
        user_uuid: "",
        server_status: "",
        now_book_uuid: "",
        read_percent: 0.0,
        msg: "",
      },
    };
  },
  filters: {
    capitalize: function (value, image_url) {
      //if (!value) return "Vertical";
      value = value.toString();
      //如果已经预先算好了
      if (value=="Vertical"||value=="Horizontal")
      {
        return value
      }
      function getUrlInfo(url) {
        let image = new Image();
        image.src = url;
        // 如果有缓存，读缓存
        if (image.complete) {
          if (image.width < image.height) {
            value = "Vertical";
          } else {
            value = "Horizontal";
          }
        } else {
          //否则加载图片
          image.onload = function () {
            image.onload = null; // 避免重复加载
            if (image.width < image.height) {
              value = "Vertical";
            } else {
              value = "Horizontal";
            }
          };
        }
      }
      getUrlInfo(image_url);
      return value;
    },
  },
  mounted() {
    this.getBook();
    this.getBookShelf();
    this.reFreshBook();
    this.hintMessage();
    this.initWebSocket();
    // window.addEventListener("scroll", this.scrollToTop);
  },
  destroyed() {
    this.$socket.close();
    // window.removeEventListener("scroll", this.scrollToTop); //销毁时解除绑定
  },
  methods: {
    getBook() {
      axios.get("/book.json").then((response) => (this.book = response.data));
      axios
        .get("/bookshelf.json")
        .then((response) => (this.bookshelf = response.data))
        .finally();
    },
    getBookShelf() {
      // axios
      //   .get("/bookshelf.json")
      //   .then(response => (this.bookshelf = response.data));
      //axios.get("/bookshelf.json").then(response => (this.bookshelf = response.data));
      // .finally(console.log(this.bookshelf));
    },
    reFreshBook() {
      // var root = this;
      // setTimeout(function() {
      //   if (root.book.extract_complete == true) {
      //     console.log("解压完成");
      //     //刷新当前文档
      //     location.reload();
      //     return;
      //   }else {
      //     console.log(root.book);
      //   }
      //   //定时执行，1秒后执行
      //   root.getBook();
      // }, 500);
    },
    //缩放提示
    hintMessage() {
      // this.$toast({
      //   message: "PC可使用“CTRL +”、“CTRL -”缩放",
      //   position: "bottom",
      //   duration: 3000
      // });
    },

    onScroll(e) {
      if (typeof window === "undefined") return;
      const top = window.pageYOffset || e.target.scrollTop || 0;
      this.btnFlag = top > 20;
    },
    toTop() {
      this.$vuetify.goTo(0);
    },

    initWebSocket() {
      //初始化weosocket
      //初始化weosocket,这些onopen是一个事件
      this.$socket.onopen = this.websocketonopen;
      this.$socket.onerror = this.websocketonerror;
      this.$socket.onmessage = this.websocketonmessage;
      this.$socket.onclose = this.websocketclose;
      this.hint = "连接建立";
      this.setButtonColor("green");
    },
    websocketonopen(e) {
      //连接建立
      //  链接ws服务器，e.target.readyState = 0/1/2/3
      //0 CONNECTING ,1 OPEN, 2 CLOSING, 3 CLOSED
      this.hint = "连接成功";
      this.setButtonColor("green");
      console.log("连接建立", e);
    },
    websocketonerror(e) {
      //连接失败
      this.hint = "连接出错";
      this.setButtonColor("#999999");
      this.initWebSocket();
      console.log("Connection Error !!!", e);
    },
    websocketonmessage(e) {
      //数据接收
      console.log(e);
      this.msgList.push(JSON.parse(e.data));
      this.hint = "接收消息";
      this.setButtonColor("blue", e);
    },
    onChangeBook: function (e, uuid) {
      // 当前元素
      this.message.now_book_uuid = uuid;
      this.message.msg = "ChangeBook";
      this.$socket.send(JSON.stringify(this.message));
      this.getBook();
    },
    websocketsend(e) {
      //数据发送
      this.$socket.send(JSON.stringify(this.message));
      //this.$socket.close(1000)
      console.log(this.$socket.readyState, e);
    },
    websocketclose(e) {
      //关闭
      this.hint = "连接断开";
      console.log("断开连接", e);
      this.setButtonColor("#888888");
      //关闭链接时触发
      var code = e.code; //  状态码表 https://developer.mozilla.org/zh-CN/docs/Web/API/CloseEvent
      var reason = e.reason;
      var wasClean = e.wasClean;
      console.log(code, reason, wasClean);
      //手动重连
      let timer = setInterval(() => {
        this.$socket.onopen();
        if (e.target.readyState == 0) {
          clearInterval(timer);
        }
      }, 3000);
    },
    setButtonColor(color) {
      var hintButton = document.getElementsByClassName("hint")[0];
      hintButton.style.background = color;
    },
  },
};
</script>

<style>
#app {
  text-align: center;
  background-color: #f6f7eb;
}

.manga_div {
  margin: auto;
  align-items: center;
}

.manga img {
  margin: auto;
  max-width: inherit;
  padding-top: 3px;
  padding-bottom: 3px;
  padding-right: 0px;
  padding-left: 0px;
  border-radius: 7px;
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}
/* 显示区域为横时的CSS样式，IE无效 */
@media screen and (min-aspect-ratio: 19/16) {
  .Vertical {
    width: 800px;
  }
  .Horizontal {
    width: 100%;
  }
}
/* 竖CSS样式，IE无效 */
@media screen and (max-aspect-ratio: 19/16) {
  .Vertical {
    width: 100%;
  }
  .Horizontal {
    width: 100%;
  }
}
</style>
