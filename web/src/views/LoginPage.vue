<template>
  <n-modal :show="showModal">
    <n-card style="width: 600px" title="登录成功" size="huge" :bordered="false" role="dialog" aria-modal="true">
      返回
    </n-card>
  </n-modal>

  <div class="flex flex-col w-screen h-screen  bg-gray-50 items-center justify-center ">

    <div class="bg-gray-50 my-10 w-80 h-48 flex-grow">
      <h2 class="text-2xl md:text-4xl font-bold text-indigo-900">用户登录</h2>
      <n-input class="h-10 my-3" type="text" v-model="username" placeholder="请输入用户名" required />
      <n-input class="h-10 my-3" type="password" v-model="password" placeholder="请输入密码" required
        @keydown.enter.native="inputEnter" />
      <n-button class="h-10 w-80ad my-3" @click="login">Login</n-button>
    </div>
  </div>
</template>

<script lang="ts">
import Header from "@/components/Header.vue";
import Bottom from "@/components/Bottom.vue";
import { NCard, NForm, NFormItem, NInput, NModal, NButton } from "naive-ui";
import { ref, defineComponent, reactive, } from "vue";
import { useRoute, useRouter } from 'vue-router';
import axios from "axios";
import { PersonOutline, LockClosedOutline } from '@vicons/ionicons5'


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
    NButton,// https://www.naiveui.com/zh-CN/os-theme/components/button
    NFormItem,
    NInput,
  },
  setup() {
    // 背景颜色,颜色选择器用  // 此处不能使用this
    const model = reactive({
      interfaceColor: "#F5F5E4",
      backgroundColor: "#E0D9CD",
    });
    const showModalRef = ref(false)
    const timeoutRef = ref(1000)
    // router オブジェクトを取得
    const router = useRouter();
    const route = useRoute();
    const countdown = () => {
      if (timeoutRef.value <= 0) {
        showModalRef.value = false
        // 登录成功后，获取query中的redirect属性，然后跳转到这个地址
        window.location.replace(route.query.redirect?.toString() || '/')
        // router.push({
        //   name: "BookShelf",
        // });
      } else {
        console.log(document.referrer);
        timeoutRef.value -= 1000
        setTimeout(countdown, 1000)
      }
    }

    const backPage = () => {
      showModalRef.value = true
      timeoutRef.value = 3000
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
          //退出登录成功后的操作
          if (resp.data.code === 200) {
            console.log(resp.data.user)
            this.$router.replace('/');
          } else {
            console.log(resp.data.msg)
          }
        })
        .catch(failResponse => { })//失败时的操作
    },
    inputEnter(e: any) {
      this.login(e);
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
