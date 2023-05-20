<template>
  <div class="w-full h-screen flex flex-col">
    <!-- <Header class="header flex-none h-12" :bookIsFolder="false" :headerTitle="getUploadTitile()" :showReturnIcon="true"
      :showSettingsIcon="false" :bookID="null" :setDownLoadLink="false">
    </Header> -->
    <!-- 可悬浮  hoverable-->
    <!-- <n-card title="注册" hoverable> 卡片内容 </n-card> -->
    <!-- 原生Form的文档： https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/form -->
    <!-- 原生输入表单的文档：https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/input -->

    <n-modal :show="showModal">
      <n-card style="width: 600px" title="登录成功" size="huge" :bordered="false" role="dialog" aria-modal="true">
        倒计时 {{ timeout / 1000 }} 秒
      </n-card>
    </n-modal>

    <div class="mian_area flex-grow">
      <!-- @submit.prevent   .prevent 表示提交以后不刷新页面，prevent是preventDefault,阻止标签默认行为，有些标签有默认行为，例如a标签的跳转链接属性href等。 -->

      <input type="text" v-model="username" placeholder="Username" required />
      <input type="password" v-model="password" placeholder="Password" required />
      <button @click="login">Login</button>
    </div>

    <Bottom class="bottom flex-none h-12" :softVersion="$store.state.server_status.ServerName
      ? $store.state.server_status.ServerName
      : 'Comigo'
      "></Bottom>
  </div>
</template>

<script lang="ts">
import Header from "@/components/Header.vue";
import Bottom from "@/components/Bottom.vue";
import { NCard, NForm, NModal } from "naive-ui";
import { ref, defineComponent, reactive,  } from "vue";
// import { useRouter } from 'vue-router';
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
    NModal,// https://www.naiveui.com/zh-CN/os-theme/components/modal
  },
  setup() {
    // 背景颜色,颜色选择器用  // 此处不能使用this
    const model = reactive({
      interfaceColor: "#F5F5E4",
      backgroundColor: "#E0D9CD",
    });
    const showModalRef = ref(false)
    const timeoutRef = ref(3000)
    // router オブジェクトを取得
    // const router = useRouter();
    const countdown = () => {
      if (timeoutRef.value <= 0) {
        showModalRef.value = false
        // document.referrer 属性会返回上一页的 URL，然后将其传递给 window.location.replace() 方法，页面会被替换成上一页并重新加载。
        // 如果需要保留历史记录，可以使用 window.history.back() 方法来返回上一页，但这种方法不会刷新页面。
        window.location.replace(document.referrer);
        // router.push({
        //   name: "BookShelf",
        // });
      } else {
        timeoutRef.value -= 1000
        setTimeout(countdown, 1000)
      }
    }

    const backPage = () => {
      showModalRef.value = true
      timeoutRef.value = 6000
      countdown()
    }

    return {
      model,
      showModal: showModalRef,
      timeout: timeoutRef,
      backPage
    };
  },

  data() {
    return {
      book_num: 0,
      drawerActive: false,
      drawerPlacement: "right",
      PageTitle: "",
      username: "admin",
      password: "admin",
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
        username: this.username,
        password: this.password,
      })
        .then(resp => {
          //退出登录成功后的操作,例如跳转页面
          if (resp.data.code === 200) {
            console.log(resp.data.user)
            this.$router.replace('/');
          } else {
            console.log(resp.data.msg)
          }
        })
        .catch(failResponse => { })//失败时的操作
    },
    login(e: any) {
      e.preventDefault() //阻止浏览器默认行为
      console.log("Login Username:" + this.username)
      console.log("Login Password:" + this.password)

      axios.post("/login", {
        username: this.username,
        password: this.password,
      }).then(resp => {
        //登录成功后的操作,例如跳转页面
        console.log(resp.data.code)
        if (resp.data.code === 200) {
          console.log(resp.data.token);
          sessionStorage.setItem('JWT_TOKEN', resp.data.token)
          this.backPage();

          // 头信息 Authorization。不知为什么不起作用，目前是cookie验证。
          const myHeaders = new Headers();
          myHeaders.append('Authorization', "Bearer " + resp.data.token);

          //手动测试接口：
          // fetch('/api/getlist?max_depth=1&sort_by=modify_time', {
          //   method: 'GET',
          //   headers: myHeaders
          // }).then(function (data) {
          //   console.log(data);
          // })

        } else {
          //登录失败时的操作
          console.log(resp.data)
        }

      })
        .catch(failResponse => {
        })

      // // JSON Obj
      // var userObj = {
      //   username: this.username,
      //   password: this.username,
      // };
      // // Obj to string
      // var body_data = JSON.stringify(userObj);
      // fetch('/api/login', {
      //   method: 'post',
      //   body: body_data,
      //   headers: {
      //     'Content-Type': 'application/json'
      //   }
      // }).then(function (data) {
      //   console.log(data);
      //   console.log(data.data.token);
      //   //将返回的结果保存在 sessionStorage 中
      //   sessionStorage.setItem('JWT_TOKEN', data.data.token);

      //   return data.text();
      // }).then(function (data) {
      //   console.log(data);
      // })
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
