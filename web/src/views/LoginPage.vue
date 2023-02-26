<template>
  <div class="w-full h-screen flex flex-col">
    <Header
      class="header flex-none h-12"
      :bookIsFolder="false"
      :headerTitle="getUploadTitile()"
      :showReturnIcon="true"
      :showSettingsIcon="false"
      :bookID="null"
      :setDownLoadLink="false"
    >
    </Header>
    <!-- 可悬浮  hoverable-->
    <!-- <n-card title="注册" hoverable> 卡片内容 </n-card> -->
    <!-- 原生Form的文档： https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/form -->
    <!-- 原生输入表单的文档：https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/input -->

    <div class="mian_area flex-grow">
      <form action="/api/form" method="post" class="form-example">
        <div class="form-example">
          <label for="username">Enter your name: </label>
          <input
            type="text"
            name="username"
            id="username"
            value="admin"
            required
          />
        </div>
        <div class="form-example">
          <label for="password">Enter your password: </label>
          <input
            type="password"
            name="password"
            id="password"
            value="admin"
            required
          />
        </div>
        <div class="form-example">
          <input type="submit" value="Subscribe!" />
        </div>
      </form>
    </div>

    <Bottom
      class="bottom flex-none h-12"
      :softVersion="
        $store.state.server_status.ServerName
          ? $store.state.server_status.ServerName
          : 'Comigo'
      "
    ></Bottom>
  </div>
</template>

<script lang="ts">
import Header from "@/components/Header.vue";
import Bottom from "@/components/Bottom.vue";
import { NCard, NForm } from "naive-ui";

import { defineComponent, reactive } from "vue";
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
