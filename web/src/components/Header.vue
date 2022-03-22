<template>
  <header class="header">
    <!-- 以后放返回箭头？ -->
    <!-- SVG资源来自 https://www.xicons.org/#/ -->
    <n-icon size="40" @click="onBackTop()">
      <book-outline v-if="!showReturnIcon" />
      <return-up-back v-if="showReturnIcon" />
    </n-icon>
    <!-- 标题，可下载压缩包 -->
    <n-space>
      <n-ellipsis style="max-width: 60vw;">
        <h2 v-if="bookIsFolder" :href="'raw/' + bookName">{{ bookName }}</h2>
        <h2>
          <a v-if="!this.bookIsFolder" :href="'raw/' + bookName">{{ bookName }}</a>
        </h2>
      </n-ellipsis>
    </n-space>
    <!-- slot，用来插入右边的设置图标 -->
    <slot></slot>
  </header>
</template>

<script>
import { useCookies } from "vue3-cookies";
import { NSpace, NIcon, NEllipsis, } from 'naive-ui'
import { BookOutline, ReturnUpBack } from '@vicons/ionicons5'
import { defineComponent } from 'vue'
export default defineComponent({
  name: "Header",
  props: ['bookIsFolder', 'bookName', 'showReturnIcon'],
  components: {
    NSpace,
    NIcon,
    NEllipsis,
    BookOutline,//图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）与导入
    ReturnUpBack,
  },
  setup() {
    const { cookies } = useCookies();
    return { cookies };
  },
  data() {
    return {
      someflag: "",
    };
  },
  methods: {
    //回首页
    onBackTop() {
      // 字符串路径
      this.$router.push('/')

    },
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.header {
  /* padding: 12px; */
  /* background: #f2f3df; */
  color: #111;
  text-align: center;

  font-size: 15px;
  display: flex;
  /* https://www.w3school.com.cn/tiy/t.asp?f=css3_flexbox_justify-content_space-between */
  justify-content: space-between;
  align-items: center;
  line-height: 75px;
}
</style>




