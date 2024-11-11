

<template>
  <!-- 上传 Upload: https://www.naiveui.com/zh-CN/os-theme/components/upload -->
  <!-- 在 web 应用程序中使用文件: https://developer.mozilla.org/zh-CN/docs/Web/API/File_API/Using_files_from_web_applications -->
  <!-- 使用 type="file" 的 <input> 元素: https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/Input/file -->
  <!--HTML 上传文件允许选择文件夹功能: https://blog.jijian.link/2020-01-08/html-upload-folder/ -->
  <!-- <input type="file"></input> -->

  <!-- multiple:是否支持上传多个文件 -->
  <!-- directory-dnd:是否支持目录拖拽上传 -->
  <!-- max:限制上传文件数量 :max="10" -->
  <div id="UploadFile" class="">
    <n-upload color="#ff69b4" multiple="true" directory-dnd :show-remove-button="false" :action="actionUrl()"
      @finish="onFinishUpload">
      <n-upload-dragger>
        <div style="margin-bottom: 12px">
          <n-icon size="48" :depth="3">
            <svg v-if="$store.state.server_status.SupportUploadFile" xmlns="http://www.w3.org/2000/svg"
              xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 512 512">
              <path
                d="M320 367.79h76c55 0 100-29.21 100-83.6s-53-81.47-96-83.6c-8.89-85.06-71-136.8-144-136.8c-69 0-113.44 45.79-128 91.2c-60 5.7-112 43.88-112 106.4s54 106.4 120 106.4h56"
                fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32">
              </path>
              <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32"
                d="M320 255.79l-64-64l-64 64"></path>
              <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32"
                d="M256 448.21V207.79"></path>
            </svg>

            <svg v-if="!$store.state.server_status.SupportUploadFile" xmlns="http://www.w3.org/2000/svg"
              xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 24 24">
              <path
                d="M19 19H5V5h14v14zM3 3v18h18V3H3zm14 12.59L15.59 17L12 13.41L8.41 17L7 15.59L10.59 12L7 8.41L8.41 7L12 10.59L15.59 7L17 8.41L13.41 12L17 15.59z"
                fill="currentColor"></path>
            </svg>

          </n-icon>
        </div>
        <n-text v-if="$store.state.server_status.SupportUploadFile" style="font-size: 16px">
          {{ $t('drop_to_upload') }}
        </n-text>

        <n-text v-if="!$store.state.server_status.SupportUploadFile" style="font-size: 16px">
          {{ $t('please_enable_upload') }}
        </n-text>

        <n-p depth="3" style="margin: 8px 0 0 0">
          {{ $t('uploaded_folder_hint') }}
        </n-p>
      </n-upload-dragger>
    </n-upload>
    <!-- 上传完毕提示 -->
    <!-- <n-p v-if="$store.state.server_status.NumberOfBooks > 0" depth="3" style="margin: 8px 0 0 0">
      {{ $t('scanned_hint').replace("XX", ($store.state.server_status.NumberOfBooks - beforeBookNum)) }}
    </n-p> -->
    <!-- 上传完毕按钮 -->
    <n-button class="h-12 w-22" color="#ff69b4" v-if="$store.state.server_status.NumberOfBooks > 0"
      @click="onBackToBookShelf">{{
      $t('back_to_bookshelf')
      }}</n-button>
  </div>
</template>

<script lang="ts">
import { NUpload, NUploadDragger, NText, NP, NIcon, NButton, useMessage } from "naive-ui";
import { defineComponent } from 'vue'
export default defineComponent({
  name: "AboutPage",
  props: ['readMode'],
  emits: ["setSome"],
  components: {
    NUpload,//上传 https://www.naiveui.com/zh-CN/os-theme/components/upload
    NUploadDragger,
    NText,
    NIcon,
    NP,
    NButton,
  },
  setup() {
    const message = useMessage()
    return {
      message,
    }
  },
  data() {
    return {
      beforeBookNum: 0,
      readerMode: "",
      upLoadHint: "",
    };
  },
  //Vue3生命周期:  https://v3.cn.vuejs.org/api/options-lifecycle-hooks.html#beforecreate
  // created : 在绑定元素的属性或事件监听器被应用之前调用。
  // beforeMount : 指令第一次绑定到元素并且在挂载父组件之前调用。
  // mounted : 在绑定元素的父组件被挂载后调用。
  // beforeUpdate: 在更新包含组件的 VNode 之前调用。。
  // updated: 在包含组件的 VNode 及其子组件的 VNode 更新后调用。
  // beforeUnmount: 当指令与在绑定元素父组件卸载之前时，只调用一次。
  // unmounted: 当指令与元素解除绑定且父组件已卸载时，只调用一次。
  created() {
    this.beforeBookNum = this.$store.state.server_status.NumberOfBooks;
    this.$store.dispatch("syncSeverStatusDataAction");
  },
  //挂载前
  beforeMount() {
  },
  onMounted() {

  },
  //卸载前
  beforeUnmount() {
  },
  methods: {
    onBackToBookShelf() {
      this.$router.push({
        name: "BookShelf"
      });
    },
    onRefreshPage() {
      location.reload();
    },
    //上传结束的回调
    onFinishUpload({ file }: any) {
      // console.log(file);
      this.message.success(file.name);
      //每次上传完成后，触发轮询的次数
      let minTryNum = 4;
      let StartNum = this.$store.state.server_status.NumberOfBooks;
      const pollTimer = setInterval(() => {
        //服务器拉取最新状态，看是否新加了书籍
        this.$store.dispatch("syncSeverStatusDataAction");
        //重新拉取书架数据,目前的写法并不需要执行
        // this.$store.dispatch("syncBookShelfDataAction");
        //ES6语法的格式化，字符串需要使用反单引号
        console.log(`this.$store.state.server_status.NumberOfBooks: ${this.$store.state.server_status.NumberOfBooks}`);
        minTryNum = minTryNum - 1;
        if (this.$store.state.server_status.NumberOfBooks > StartNum) {
          minTryNum = minTryNum - 1;
        }
        if (minTryNum <= 0) {
          clearInterval(pollTimer);
        }

      }, 2000);
    },
    // 拼接上传接口路径
    actionUrl() {
      let protocol = 'http://'
      if (window.location.protocol === "https") {
        protocol = 'https://'
      }
      return protocol + window.location.host + '/api/upload';
    },
  },
  computed: {
  },
});
</script>

<style scoped>
#UploadFile {
  padding: 20px;
}

.n-upload {
  --tw-bg-opacity: 0.5;
  background-color: rgba(249, 250, 251, var(--tw-bg-opacity));
}
</style>
