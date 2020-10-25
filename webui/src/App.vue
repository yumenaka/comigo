<template>
  <div id="app" class="app_div">
    <Header v-if="defaultSetiing.page_mode === 'multi_page_mode'">
      <h2>
        <a v-bind:href="'raw/' + book.name">{{ book.name }}</a>
      </h2>
      <h4>总页数：{{ book.page_num }}</h4>
    </Header>
<MultiView>
    <div  v-for="(page, key) in book.pages" :key="page.url" class="manga">
      <img
        v-lazy="page.url"
        v-bind:H="page.height"
        v-bind:W="page.width"
        v-bind:key="k"
        v-bind:class="page.class | check_image(page.url)"
      />
      <p>{{key+1}}/{{book.page_num}}</p>
    </div>
</MultiView>
<v-pagination v-if="defaultSetiing.page_mode === 'single_page_mode'" v-model="book" :length="book.page_num">
</v-pagination>

    <p></p>
    <v-btn v-scroll="onScroll" v-show="btnFlag" fab color="#bbcbff" bottom right @click="toTop">▲</v-btn>
  </div>
</template>

<script>
//代码参考：https://github.com/bradtraversy/vue_crash_todolist
import axios from "axios";
import Header from "./components/Header.vue";
import MultiView from "./views/MultiPage.vue";

export default {
  name: "app",
  components: {
    //WebSocketTest
    Header,
    MultiView,
  },
  data() {
    return {
      book: null,
      bookshelf: null,
      defaultSetiing: {
        page_mode:"multi_page_mode",//single_page_mode,multi_page_mode,select_mode
      },
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

  mounted() {
    this.initPage();
    this.getBook();
    this.getBookShelf();
    this.hintMessage();
    this.initWebSocket();
    // window.addEventListener("scroll", this.scrollToTop);
  },
  destroyed() {
    this.$socket.close();
    // window.removeEventListener("scroll", this.scrollToTop); //销毁时解除绑定
  },
  methods: {
    initPage() {
      //Header.book=this.book;
      this.$cookies.keys();
    },
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
  // 判断图片横宽比
  filters: {
    check_image: function (value, image_url) {
      //if (!value) return "Vertical";
      value = value.toString();
      //如果已经预先算好了
      if (value == "Vertical" || value == "Horizontal" || image_url == "") {
        return value;
      }
      //value=this.$options.methods.getImageInfo(image_url);
      function getImageInfo(url) {
        let image = new Image();
        image.src = url;
        // 如果有缓存，读缓存
        if (image.complete) {
          if (image.width < image.height) {
            return "Vertical";
          } else {
            return "Horizontal";
          }
        } else {
          //否则加载图片
          image.onload = function () {
            image.onload = null; // 避免重复加载
            if (image.width < image.height) {
              return "Vertical";
            } else {
              return "Horizontal";
            }
          };
        }
      }
      value = getImageInfo(image_url);
      return value;
    },
  },
};
</script>

<style>
#app {
  text-align: center;
  background-color: #f6f7eb;
}

.app_div {
  margin: auto;
  align-items: center;
}

.manga img{
  margin: auto;
  max-width: inherit;
  padding-top: 3px;
  padding-bottom: 3px;
  padding-right: 0px;
  padding-left: 0px;
  border-radius: 7px;
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}


/* 竖屏(显示区域)CSS样式，IE无效 */
@media screen and (max-aspect-ratio: 19/19) {
  .Vertical {
    width: 100%;
  }
  .Horizontal {
    width: 100%;
  }
}


/* 横屏（显示区域）时的CSS样式，IE无效 */ 
@media screen and (min-aspect-ratio: 19/19) {
  .Vertical {
    width: 800px;
  }
  .Horizontal {
    width: 90%;
  }
}


/* 高分横屏（显示区域）时的CSS样式，IE无效 */ 
/* min-width 输出设备中页面最小可视区域宽度 大于这个width时，其中的css起作用 超宽屏 */
@media screen and (min-aspect-ratio: 19/19) and (min-width: 1922px) {
  .Vertical {
    width: 800px;
  }
  .Horizontal {
    width: 1900px;
  }
}


/* max-width 输出设备中页面最大可视区域宽度 小于这个width时，其中的css起作用 1920x1080屏 */
/* 
@media screen and (min-aspect-ratio: 19/16) and (max-width: 1921px) {
  .Vertical {
    width: 800px;
  }
  .Horizontal {
    width: 1700px;
  }
}
 */

 
/* max-width 输出设备中页面最大可视区域宽度 小于这个width时，其中的css起作用 1366x768屏 */
/*
@media screen and (min-aspect-ratio: 19/16) and (max-width: 1367px) {
  .Vertical {
    width: 800px;
  }
  .Horizontal {
    width: 1200px;
  }
}
 */

</style>
