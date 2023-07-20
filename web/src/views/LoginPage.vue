<template>
  <n-modal :show="showModal">
    <n-card style="width: 600px" title="登录成功" size="huge" :bordered="false" role="dialog" aria-modal="true">
      返回
    </n-card>
  </n-modal>

  <div class="flex flex-col w-screen h-screen  bg-gray-50  items-center justify-center ">
    <div class="w-100 h-70 p-8 flex flex-col items-center bg-gray-300 rounded-lg border border-gray-800">
      <h2 class="text-2xl md:text-3xl font-bold text-indigo-900">登录</h2>
      <n-input class="h-10 my-3 w-9/4" type="text" v-model:value="username" placeholder="请输入用户名" required />
      <n-input class="h-10 my-3 w-9/12" type="password" v-model:value="password" placeholder="请输入密码" required
        @keydown.enter.native="login" />
      <n-button class="h-10 w-3/4 my-3  rounded-lg" type="primary" @click="login">Login</n-button>
    </div>
  </div>
</template>

<script lang="ts">
import Header from "@/components/Header.vue";
import Bottom from "@/components/Bottom.vue";
import { NCard, NForm, NFormItem, NInput, NModal, NButton, useMessage, useDialog } from "naive-ui";
import { ref, defineComponent, reactive, } from "vue";
import { useRoute, } from 'vue-router';
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
    NInput,// https://www.naiveui.com/zh-CN/os-theme/components/input
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
    const route = useRoute();
    const message = useMessage()
    const dialog = useDialog()
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
      backPage,
      dialog,
      handleConfirm() {
        dialog.warning({
          title: '警告',
          content: '格式错误，请检查输入',
          positiveText: '确定',
          // negativeText: '不确定',
          onPositiveClick: () => {
            message.success('OK')
          },
          // onNegativeClick: () => {
          //   message.error('不确定')
          // }
        })
      },
      handleSuccess() {
        dialog.success({
          title: '成功登录',
          content: '成功登录，返回原页面',
          positiveText: 'OK',
          onPositiveClick: () => {
            message.success('成功登录')//message
          }
        })
      },
      handleError() {
        console.log(dialog)
        dialog.error({
          title: '错误',
          content: '用户名或密码错误',
          positiveText: '确认',
          onPositiveClick: () => {
            // message.success('用户名或密码错误')
          }
        })
      },
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
    login(e: any) {
      e.preventDefault() //阻止浏览器默认行为
      axios.post("/login", {
        username: this.username,
        password: this.password,
      })
        .then((response) => {
          // 登录成功
          const respData: any = response.data;
          console.log(respData.token);
          sessionStorage.setItem('JWT_TOKEN', respData.token);
          this.backPage();
        })
        .catch((error) => {
          // 登录请求失败
          this.handleError();
        });
    },
  },
});
</script>

