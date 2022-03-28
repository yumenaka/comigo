<template>
  <header class="header">
    <!-- 返回箭头,点击返回上一页 -->
    <n-icon v-if="showReturnIcon" size="40" @click="onClickReturnIcon()">
      <return-up-back />
    </n-icon>

    <!-- 一本书，点击返回主页，但是目前应该没有任何反应 -->
    <n-icon v-if="!showReturnIcon" @click="onClickToTop()" size="40">
      <book-outline />
    </n-icon>

    <!-- 标题，可下载压缩包 -->
    <n-space>
      <n-ellipsis style="max-width: 60vw;">
        <span class="text-lg" v-if="!setDownLoadLink">{{ bookName }}</span>
        <span class="text-lg">
          <a v-if="this.setDownLoadLink" :href="'raw/' + bookName">{{ bookName }}</a>
        </span>
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
  props: ['setDownLoadLink', 'bookName', 'showReturnIcon',],
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
    //点击返回的时候，后退到上一页
    onClickReturnIcon() {
      this.$router.back()
    },
    //点击返回的时候，后退到上一页
    onClickToTop() {
      this.$router.push('/')
    },
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.header {
  color: #111;
  text-align: center;
  display: flex;
  /* https://www.w3school.com.cn/tiy/t.asp?f=css3_flexbox_justify-content_space-between */
  justify-content: space-between;
  align-items: center;
  line-height: 75px;
}
</style>




