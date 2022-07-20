

<template>
  <!-- 上传 Upload: https://www.naiveui.com/zh-CN/os-theme/components/upload -->
  <!-- 在 web 应用程序中使用文件: https://developer.mozilla.org/zh-CN/docs/Web/API/File_API/Using_files_from_web_applications -->
  <!-- 使用 type="file" 的 <input> 元素: https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/Input/file -->
  <!--HTML 上传文件允许选择文件夹功能: https://blog.jijian.link/2020-01-08/html-upload-folder/ -->
  <!-- <input type="file"></input> -->

  <!-- multiple:是否支持上传多个文件 -->
  <!-- directory-dnd:是否支持目录拖拽上传 -->
  <!-- max:限制上传文件数量 :max="10" -->
  <div id="UploadFile">
    <n-upload multiple directory-dnd :show-remove-button="false" :action="this.actionUrl()"
      @finish="this.onFinishUpload">
      <n-upload-dragger>
        <div style="margin-bottom: 12px">
          <n-icon size="48" :depth="3">
            <archive-icon />
          </n-icon>
        </div>
        <n-text style="font-size: 16px">
          {{ this.$t('drop_to_upload') }}
        </n-text>
        <n-p depth="3" style="margin: 8px 0 0 0">
          {{ this.$t('uploaded_folder_hint') }}
        </n-p>
      </n-upload-dragger>
    </n-upload>
    <!-- 上传文件扫描好了的提示 -->
    <n-p v-if="this.$store.state.server_status.NumberOfBooks > 0" depth="3" style="margin: 8px 0 0 0">
      {{ this.$t('scanned_hint').replace("XX", this.$store.state.server_status.NumberOfBooks) }}
    </n-p>
    <!-- 上传完毕刷新页面的按钮 -->
    <n-button class="w-22 h-12"  color="#ff69b4" v-if="this.$store.state.server_status.NumberOfBooks > 0" @click="this.onRefreshPage">{{ $t('refresh_page')
    }}</n-button>
  </div>
</template>

<script>
import { NUpload, NUploadDragger, NText, NP, NIcon, NButton } from "naive-ui";
import { defineComponent, reactive } from 'vue'
import { ArchiveOutline as ArchiveIcon } from "@vicons/ionicons5";
export default defineComponent({
  name: "AboutPage",
  props: ['readMode'],
  emits: ["setSome"],
  components: {
    ArchiveIcon,
    NUpload,//上传 https://www.naiveui.com/zh-CN/os-theme/components/upload
    NUploadDragger,
    NText,
    NIcon,
    NP,
    NButton,
  },
  setup() {
    // 背景颜色,颜色选择器用
    const model = reactive({
      interfaceColor: "#F5F5E4",
      backgroundColor: "#E0D9CD",
    });
    return {
      // 背景色
      model,
    }
  },
  data() {
    return {
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
  },
  //挂载前
  beforeMount() {
  },
  onMounted() {
    // https://v3.cn.vuejs.org/api/options-lifecycle-hooks.html#beforemount
    this.$nextTick(function () {
      //视图渲染之后运行的代码
    })
  },
  //卸载前
  beforeUnmount() {
  },
  methods: {
    onRefreshPage() {
      location.reload();
    },
    //上传结束的回调
    onFinishUpload({ file }) {
      console.log(file);
      //每次上传完成后，触发轮询的次数
      let minTryNum = 20;
      const pollTimer = setInterval(() => {
        //服务器拉取最新状态，看是否新加了书籍
        this.$store.dispatch("syncSeverStatusDataAction");
        //ES6语法，需要反单引号
        console.log(`this.$store.state.server_status.NumberOfBooks：${this.$store.state.server_status.NumberOfBooks}`);
        if (this.$store.state.server_status.NumberOfBooks > 0) {
          //重新拉取书架数据
          this.$store.dispatch("syncBookShelfDataAction");
          minTryNum = minTryNum - 1;
          if (minTryNum <= 0) {
            clearInterval(pollTimer);
          }
        }
      }, 1500);
    },
    // 拼接上传接口路径
    actionUrl() {
      var protocol = 'http://'
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
  padding-bottom: 20px;
  padding-left: 20px;
  padding-right: 20px;
  padding-top: 20px;
  background: v-bind("model.color");
}
</style>
