<template>
  <div class="websocket">
    <input type="text" v-model="message.msg" />
    <button class="hint">{{hint}}</button>
    <button @click="send">发言</button>
    <div class="chat-title">聊天记录：</div>
    <div v-for="(item,index) in msgList" :key="index" class="chat-box">{{item.message.msg}}</div>
  </div>
</template>

<script>
export default {
  name: "Websocket",
  data() {
    return {
      message: {msg: "",user:"admin",page:0},
      page: 0,
      websock: "",
      hint:"",
      msgList: []
    };
  },
  created() {
      this.initWebSocket();
  },
  //离开后断开websocket连接
  destroyed() {
    this.websock.close();
  },
  methods: {
    send() {
      this.websock.send(JSON.stringify(this.message));
      this.message = {msg: "",user:"",page:0};
      this.hint="发送消息";
    },      
    initWebSocket(){ //初始化weosocket
        var host = document.location.host; //域名与端口
        this.websock = new WebSocket("ws://" + host + "/ws");
        this.websock.onmessage = this.websocketonmessage;
        this.websock.onopen = this.websocketonopen;
        this.websock.onerror = this.websocketonerror;
        this.websock.onclose = this.websocketclose;
      },
      websocketonopen(){ //连接建立
        console.log("Connection open ...");
        this.hint="连接成功";
        console.log("websocket连接成功");
        this.setButtonColor("green");
      },
      websocketonerror(){//失败重连
        console.log("Connection Error !!!");
        this.hint="连接出错";
        this.setButtonColor("red");
        this.initWebSocket();
      },
      websocketonmessage(evt){ //数据接收
        console.log(evt);
        this.msgList.push(JSON.parse(evt.data));
        this.hint="接收消息";
        this.setButtonColor("blue");
      },
      websocketsend(Data){//数据发送
        this.websock.send(Data);
      },
      websocketclose(e){  //关闭
      this.hint="连接断开";
        console.log('断开连接',e);
        this.setButtonColor("#888888");
      },
      setButtonColor(color) {
        var hintButton = document.getElementsByClassName("hint")[0];
        hintButton.style.background=color;
      },
  },
  mounted() {
  }
};
</script>

<style scoped>
.hint {
  background: red;
  color: #fff;
  border: none;
  padding: 5px 9px;
  border-radius: 50%;
  cursor: pointer;
  float: right;
}

.websocket {
  width: inherit;
  margin: 0 auto;
  background-color: aliceblue;
  height: 200px;
  text-align: center;
  padding-top: 20px;
}
.chat-title {
  text-align: left;
  margin-left: 100px;
  margin-top: 20px;
}
.chat-box {
  background-color: white;
  width: 100px;
  margin: 0 auto;
}
</style>
