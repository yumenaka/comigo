<template>
  <div class="w-full h-screen flex flex-col">
    <Header class="header flex-none h-12" :bookIsFolder="false" :headerTitle="getUploadTitile()" :showReturnIcon="true"
      :showSettingsIcon="false" :bookID="null" :setDownLoadLink="false">
    </Header>
    <!-- 可悬浮  hoverable-->
    <!-- <n-card title="注册" hoverable> 卡片内容 </n-card> -->
    <!-- 原生Form的文档： https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/form -->
    <!-- 原生输入表单的文档：https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/input -->

    <div class="mian_area flex-grow">
      <!-- @submit.prevent   .prevent 表示提交以后不刷新页面，prevent是preventDefault,阻止标签默认行为，有些标签有默认行为，例如a标签的跳转链接属性href等。 -->

      <input type="text" v-model="username" placeholder="Username" required />
      <input type="password" v-model="password" placeholder="Password" required />
      <button @click="login">Login</button>
    </div>

    <Bottom class="bottom flex-none h-12" :softVersion="
      $store.state.server_status.ServerName
        ? $store.state.server_status.ServerName
        : 'Comigo'
    "></Bottom>
  </div>
</template>

<script lang="ts">
import Header from "@/components/Header.vue";
import Bottom from "@/components/Bottom.vue";
import { NCard, NForm } from "naive-ui";
import { defineComponent, reactive } from "vue";
import axios from "axios";

export default defineComponent({
  name: "LoginPage",
  props: ["readMode"],
  emits: ["setSome"],
  components: {
    Header, // 自定义页头
    Bottom, // 自定义页尾
    NCard, // https://www.naiveui.com/zh-CN/os-theme/components/card
    NForm, // https://www.naiveui.com/zh-CN/os-theme/components/form
  },
  setup() {
    // 背景颜色,颜色选择器用  // 此处不能使用this
    const model = reactive({
      interfaceColor: "#F5F5E4",
      backgroundColor: "#E0D9CD",
    });
    return {
      model,
    };
  },

  data() {
    return {
      book_num: 0,
      drawerActive: false,
      drawerPlacement: "right",
      PageTitle: "",
      username: "",
      password: "",
    };
  },
  created() {
    // 当前颜色
    const tempBackgroundColor = localStorage.getItem("BackgroundColor");
    if (typeof tempBackgroundColor === "string") {
      this.model.backgroundColor = tempBackgroundColor;
    }
    const tempInterfaceColor = localStorage.getItem("InterfaceColor");
    if (typeof tempInterfaceColor === "string") {
      this.model.interfaceColor = tempInterfaceColor;
    }
  },
  methods: {
    logout() {
      console.log("logout");
      axios.post("/logout", {
        username: "admin",
        password: this.password,
      })
        .then(resp => {
          //登录成功后的操作,例如跳转页面
          if (resp.data.code === 0) {
            console.log(resp.data.user)
            this.$router.replace('/')
          } else {
            console.log(resp.data.msg)
          }
        })
        //登录失败时的操作
        .catch(failResponse => { })
    },
    login() {
      console.log("Login Username:" + this.username)
      console.log("Login Password:" + this.password)
      axios.post("/login", {
        username: "admin",
        password: this.password,
      })
        .then(resp => {
          //登录成功后的操作,例如跳转页面
          if (resp.data.code === 0) {
            console.log(resp.data.user)
            this.$router.replace('/')
          } else {
            console.log(resp.data.msg)
          }
        })
        //登录失败时的操作
        .catch(failResponse => { })
    },
    getUploadTitile() {
      //如果没有一本书
      if (this.$store.state.server_status.SupportUploadFile === false) {
        return this.$t("no_support_upload_file");
      }
      //如果没有一本书
      if (this.$store.state.server_status.NumberOfBooks === 0) {
        return this.$t("no_book_found_hint");
      }
      return (
        this.$t("number_of_online_books") +
        this.$store.state.server_status.NumberOfBooks
      );
    },
    remoteIsWindows() {
      if (!this.$store.state.server_status) {
        return false;
      }
      console.dir(this.$store.state.server_status.Description);
      return (
        this.$store.state.server_status.Description.indexOf("windows") !== -1
      );
    },
  },
});
</script>

<style scoped>
.header {
  background: v-bind("model.interfaceColor");
}

.bottom {
  background: v-bind("model.interfaceColor");
}

.mian_area {
  background: v-bind("model.backgroundColor");
}
</style>
