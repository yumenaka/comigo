<template>
  <!-- class='w-28 md:w-33 lg:w-48'   Width of 28 by default, 32 on medium screens, and 48 on large screens -->
  <!-- 响应式设计：https://www.tailwindcss.cn/docs/responsive-design -->
  <!-- sm<640px  md<768px lg<1024px  lg<1280px 2xl<1536px-->
  <!-- 宽，只有一些典型值：https://www.tailwindcss.cn/docs/width  -->
  <!-- 高，只有一些典型值：https://www.tailwindcss.cn/docs/height -->
  <!-- 什么是REM：https://www.runoob.com/w3cnote/px-em-rem-different.html -->
  <!-- 边框圆角 rounded-xl：https://www.tailwindcss.cn/docs/border-radius -->
  <!-- 盒阴影 shadow: https://www.tailwindcss.cn/docs/box-shadow -->
  <!-- 外边距 m-x m-y  https://www.tailwindcss.cn/docs/margin -->
  <!-- 字体粗细： https://www.tailwindcss.cn/docs/font-weight -->
  <a :href="openURL" :target="a_target">
    <div
      class="relative w-32 h-44 mx-4 my-4 bg-gray-200 rounded shadow-xl hover:shadow-2xl ring-1 ring-gray-400 hover:ring hover:ring-blue-500 .bg-top bg-cover"
      :style="getBackgroundImageStyle()">
      <div v-if="childBookNum != ''" class="absolute inset-x-0 top-0 text-right">
        <span class="text-2xl text-yellow-500 font-black text-shadow">{{ childBookNum }}</span>
      </div>
      <!-- 文本对齐:       https://www.tailwindcss.cn/docs/text-align -->
      <!-- 定位:          https://www.tailwindcss.cn/docs/top-right-bottom-left -->
      <!-- 背景色:        https://www.tailwindcss.cn/docs/background-color -->
      <!-- 背景色不透明度: https://www.tailwindcss.cn/docs/background-opacity -->
      <!-- 文本溢出：      https://www.tailwindcss.cn/docs/text-overflow -->
      <!-- 字体粗细：     https://www.tailwindcss.cn/docs/font-weight -->
      <div v-if="showTitle"
        class="absolute inset-x-0 bottom-0 h-1/4 bg-gray-100 bg-opacity-80 font-semibold border-blue-800 rounded-b">
        <!-- 如果把链接的 target 属性设置为 "_blank"，该链接会在新窗口中打开。 -->
        <span class="absolute inset-x-0  font-bold top-0 p-1 align-middle">{{
            shortTitle
        }}</span>
      </div>
    </div>
  </a>

</template>

<script lang="ts">
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
// import { NCard, } from 'naive-ui'
import { useCookies } from "vue3-cookies";
import { defineComponent } from "vue";

export default defineComponent({
  name: "BookCover",
  props: ["title", "image_src", "id", "readerMode", "showTitle", "childBookNum", "openURL", "a_target", "simplifyTitle"],
  components: {
    // NCard,
    // NEllipsis,
  },
  setup() {
    const { cookies } = useCookies();
    return { cookies };
  },
  computed: {
    shortTitle(): string {
      let short_title = this.title
      //使用 JavaScript replace() 方法替换掉一些字符串
      if (this.simplifyTitle) {
        //中：/[\u4e00-\u9fa5]/  日：/[\u0800-\u4e00]/  韩：/[\uac00-\ud7ff]/  空格：[\s]
        //左半部分：任意中日韩字符或空格，g表示多次匹配、不限次数
        short_title = short_title.replace(/[\\[\\(（【][A-Za-z0-9_\-×\s+\u4e00-\u9fa5\u0800-\u4e00\uac00-\ud7ff]+/g, "");
        //右半部分
        short_title = short_title.replace(/[\]）】\\)]/g, "");
        //.zip .rar 等文件名
        short_title = short_title.replace(/.(zip|rar|cbr|cbz|tar|pdf|mp3|mp4|flv|gz|webm|gif|png|jpg|jpeg|webp|svg|psd|bmp|tif)/g, "");
        //域名（误伤多?）   参考了正则大全（https://any86.github.io/any-rule/）的网址(URL)
        const domain_reg = /^(((ht|f)tps?):\/\/)?([^!@#$%^&*?.\s-]([^!@#$%^&*?.\s]{0,63}[^!@#$%^&*?.\s])?\.)+[a-zA-Z]{2,6}\/?/g;
        short_title = short_title.replace(domain_reg, "");
        //开头的空格
        short_title = short_title.replace(/^[\s]/g, "");
        //开头的特殊字符
        short_title = short_title.replace(/^[\\\-`~!@#$^&*()=|{}':;'@#￥……&*（）——|{}‘；：”“'。，、？]/, "");
      }
      //删减到空字符串，就不删减了
      if (short_title.length == 0) {
        return this.title;
      }

      if (short_title.length <= 15) {
        return short_title;
      }
      return `${short_title.substr(0, 15)}…`;
    },
  },
  data() {
    return {
      test: "",
    };
  },
  methods: {
    // 回首页
    onBackTop() {
      // 字符串路径
      this.$router.push("/");
    },
    getBackgroundImageStyle() {
      return `background-image: url(${this.getThumbnailsImageUrl()});`;
    },
    getThumbnailsImageUrl() {
      // 按照“/”分割字符串
      const arrUrl = this.image_src.split("/");
      // console.log(arrUrl)
      if (arrUrl[0] === "api") {
        return `${this.image_src}&resize_width=256&resize_height=360&thumbnail_mode=true`;
      }
      return this.image_src;
    },
    // 自己构建一个<a>链接，后来发现不如可以直接用router-link与命名路由
    getHref() {
      // 当前URL
      const url = document.location.toString();
      // 按照“/”分割字符串
      const arrUrl = url.split("/");
      // 拼一个完整的图片URL
      if (this.readerMode === "flip") {
        const new_url = `${arrUrl[0]}//${arrUrl[2]}/#` + `f/${this.id}`;
        console.log(new_url);
        return new_url;
      }
      if (this.readerMode === "scroll") {
        const new_url = `${arrUrl[0]}//${arrUrl[2]}/#` + `s/${this.id}`;
        console.log(new_url);
        return new_url;
      }
    },
  },
});
</script>
//自定义样式
<style scoped>

</style>
