<template>
  <div id="multiPage">
    <Header>
      <h2>
        <a v-if="!this.$store.state.book.IsFolder" v-bind:href="'raw/' + this.$store.state.book.name"
          >{{ this.$store.state.book.name }}【Download】</a
        >
        <a v-if="this.$store.state.book.IsFolder" v-bind:href="'raw/' + this.$store.state.book.name">{{
          this.$store.state.book.name
        }}</a>
      </h2>
      <h4>总页数：{{ this.$store.state.book.page_num }}</h4>
    </Header>
    <div v-for="(page, key) in this.$store.state.book.pages" :key="page.url" class="manga">
      <img
        v-lazy="page.url"
        v-bind:H="page.height"
        v-bind:W="page.width"
        v-bind:key="key"
        v-bind:class="page.class | check_image(page.url)"
      />
      <p>{{ key + 1 }}/{{ AllPageNum }}</p>
    </div>
    <p></p>
    <v-btn
      v-scroll="onScroll"
      v-show="btnFlag"
      fab
      color="#bbcbff"
      bottom
      right
      @click="toTop"
      >▲</v-btn
    >
    <slot></slot>
  </div>
</template>

<script>
import Header from "./Header.vue";

export default {
  components: {
    Header,
  },
  // props: ['book'],
  //组件的 data 选项必须是一个函数
  //每个实例可以维护一份被返回对象的独立的拷贝
  data() {
    return {
      // book: this.$store.state.book,
      // bookshelf: null,
      // defaultSetting: null,
      page_mode: "multi",
      btnFlag: false,
      duration: 300,
      offset: 0,
      easing: "easeInOutCubic",
      AllPageNum: this.$store.state.book.page_num,
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
    //以后再研究WebSocks
    //this.initWebSocket();
  },
  destroyed() {
    //this.$socket.close();
  },
  methods: {
    initPage() {
      this.$cookies.keys();
    },
    getBook: function () {
       return this.$store.state.book;
    },
    getNumber: function (number) {
      this.page = number;
      console.log(number);
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
    },
    websocketonopen(e) {
      //连接建立
      //  链接ws服务器，e.target.readyState = 0/1/2/3
      //0 CONNECTING ,1 OPEN, 2 CLOSING, 3 CLOSED
      this.hint = "连接成功";
      console.log("连接建立", e);
    },
    websocketonerror(e) {
      //连接失败
      this.hint = "连接出错";
      this.initWebSocket();
      console.log("Connection Error !!!", e);
    },
    websocketonmessage(e) {
      //数据接收
      console.log(e);
      this.msgList.push(JSON.parse(e.data));
      this.hint = "接收消息";
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
  },
  // 判断图片横宽比
  filters: {
    check_image: function (value, image_url) {
      // console.log(value);
      // console.log(image_url);
      //如果已经算好了
      value = value.toString();
      if (value == "Vertical" || value == "Horizontal") {
        return value;
      }
      if (value == "") {
        console.log("图片信息为空，开始本地JS分析" + image_url);
      }
      function getImageInfo(url) {
        let image = new Image();
        image.src = url;
        // 如果有缓存，读缓存。
        //还要避免默认占位图片的情况，目前远程网速较慢时似乎会出错
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
    width: 900px;
  }
  .Horizontal {
    width: 95%;
  }
}

/* 高分横屏（显示区域）时的CSS样式，IE无效 */
/* min-width 输出设备中页面最小可视区域宽度 大于这个width时，其中的css起作用 超宽屏 */
@media screen and (min-aspect-ratio: 19/19) and (min-width: 1922px) {
  .Vertical {
    width: 1000px;
  }
  .Horizontal {
    width: 1900px;
  }
}
</style>