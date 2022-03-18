<template>
  <header class="header">
    <!-- 以后放返回箭头？ -->
    <!-- SVG资源来自 https://www.xicons.org/#/ -->
    <n-icon size="40">
      <book-outline />
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
import { BookOutline } from '@vicons/ionicons5'
import { defineComponent } from 'vue'
export default defineComponent({
  name: "Header",
  props: ['bookIsFolder', 'bookName'],
  components: {
    NSpace,
    NIcon,
    NEllipsis,
    BookOutline,//图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）与导入
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
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.header {
  /* background: #f2f3df; */
  color: #111;
  text-align: center;
  padding: 12px;
  font-size: 12px;
  display: flex;
  /* https://www.w3school.com.cn/tiy/t.asp?f=css3_flexbox_justify-content_space-between */
  justify-content: space-between;
  align-items: center;
}
</style>




