<template>
  <div class="ws">
    <input type="text" v-model="message.msg" />
    <button class="hint">{{hint}}</button>
    <button @click="send">发言</button>
    <div class="chat-title">聊天记录：</div>
    <div v-for="(item,index) in msgList" :key="index" class="chat-box">{{item.msg}}</div>
  </div>
</template>

<script>
export default {
  name: "ws",
  data() {
    return {
      message: {user_uuid: "",server_status:"",now_book_uuid:"",read_percent:0.0,msg:""},
      page: 0,
      websock: "",
      hint:"发送消息",
      msgList: []
    };
  },
  created() {
  },
  //离开后断开websocket连接
  destroyed() {
    this.$socket.close();
  },
  mounted() {
    //this.$connect();
    this.initWebSocket();
  },
  methods: {
    send() {
      this.$socket.send(JSON.stringify(this.message));
      this.message = {user_uuid: "",server_status:"",now_book_uuid:"",now_page:0.0,msg:""};
      this.hint="发送消息";
    },      
    initWebSocket(){ //初始化weosocket
      //初始化weosocket,这些onopen是一个事件，当
      this.$socket.onopen = this.websocketonopen;
      this.$socket.onerror = this.websocketonerror;
      this.$socket.onmessage = this.websocketonmessage;
      this.$socket.onclose = this.websocketclose;
      this.hint="连接建立";
      this.setButtonColor("green");
      },
      websocketonopen(e){ //连接建立
        //  链接ws服务器，e.target.readyState = 0/1/2/3
        //0 CONNECTING ,1 OPEN, 2 CLOSING, 3 CLOSED
        this.hint="连接成功";
        this.setButtonColor("green");
        console.log("连接建立",e);
      },
      websocketonerror(e){//连接失败
        this.hint="连接出错";
        this.setButtonColor("#999999");
        this.initWebSocket();
        console.log("Connection Error !!!",e);
      },
      websocketonmessage(e){ //数据接收
        console.log(e);
        this.msgList.push(JSON.parse(e.data));
        this.hint="接收消息";
        this.setButtonColor("blue",e);
      },
      websocketsend(e){//数据发送
        this.$socket.send(JSON.stringify(this.message));
        //this.$socket.close(1000)
        console.log(this.$socket.readyState,e);
      },
      websocketclose(e){  //关闭
      this.hint="连接断开";
        console.log('断开连接',e);
        this.setButtonColor("#888888");
        //关闭链接时触发
        var code = e.code;//  状态码表 https://developer.mozilla.org/zh-CN/docs/Web/API/CloseEvent
        var reason = e.reason;
        var wasClean = e.wasClean;
        console.log(code,reason,wasClean);
        //手动重连
        let timer = setInterval(() => {
          this.$socket.onopen();
            if(e.target.readyState==0){
              clearInterval(timer);
            }
        }, 3000);
      },
      setButtonColor(color) {
        var hintButton = document.getElementsByClassName("hint")[0];
        hintButton.style.background=color;
      },
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

.ws {
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
