<template>
  <!-- 外边距: m-2 https://www.tailwindcss.cn/docs/margin -->
  <!-- 内边距： p-4 https://www.tailwindcss.cn/docs/padding -->
  <header class="header p-1 m-0 h-12 flex justify-between content-center">
    <!-- 返回箭头,点击返回上一页 -->
    <n-icon class="p-0 m-0" v-if="showReturnIcon" size="40" @click="onClickReturnIcon()">
      <return-up-back />
    </n-icon>

    <!-- 一本书，点击返回主页，但是目前应该没有任何反应 -->
    <n-icon class="p-0 m-0" v-if="!showReturnIcon" @click="onClickToTop()" size="40">
      <book-outline />
    </n-icon>

    <!-- 标题-->
    <!-- 文本颜色： https://www.tailwindcss.cn/docs/text-color -->
    <!-- 文本装饰（下划线）：https://www.tailwindcss.cn/docs/text-decoration -->
    <!-- 文本溢出：https://www.tailwindcss.cn/docs/text-overflow -->
    <!-- 字体粗细:https://www.tailwindcss.cn/docs/font-weight -->
    <div class="p-0 m-0 py-2 font-semibold content-center truncate">
      <!-- 标题，只显示 -->
      <span class="text-lg" v-if="!setDownLoadLink">{{ headerTitle }}</span>
      <!-- 标题，可下载压缩包 -->
      <span class="text-lg text-blue-700 text-opacity-100  hover:underline">
        <a v-if="this.setDownLoadLink" :href="'api/raw/' + bookID + '/' + headerTitle">{{ headerTitle }}</a>
      </span>
    </div>
    <!-- slot，用来插入右边的设置图标 -->
    <slot></slot>
  </header>
</template>

<script>
import { useCookies } from "vue3-cookies";
import { NIcon, } from 'naive-ui'
import { BookOutline, ReturnUpBack } from '@vicons/ionicons5'
import { defineComponent } from 'vue'
export default defineComponent({
  name: "ComigoHeader",
  props: ['setDownLoadLink', 'headerTitle', 'bookID', 'showReturnIcon',],
  components: {
    NIcon,
    BookOutline,//图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）与导入
    ReturnUpBack,
  },
  setup() {
    const { cookies } = useCookies();
    // console.log(window.history)
    return { cookies };
  },
  data() {
    return {
      someflag: "",
    };
  },
  methods: {
    //点击返回图标的时候，后退到上一页或主页
    onClickReturnIcon() {
      // console.log(window.history)
      //如果直接进入本页面，没有上一页，那么回到主页。不过这时候浏览器back按钮本来应该也不能按。
      if (window.history.length === 1) {
        this.$router.push('/')
        return
      }
      //如果是新建标签页,并不是从书架点进去的时候，window.history.state.back为null，直接回到主页。
      if (window.history.state.back === null) {
        this.$router.push('/')
        return
      }
      //其他情况下，后退一页。与单击浏览器中的“后退”按钮相同。
      this.$router.back();
      // location.reload();
    },
    //点击主页图标的时候，回到主页
    onClickToTop() {
      this.$router.push('/')
    },
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>




