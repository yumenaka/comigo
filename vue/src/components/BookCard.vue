<template>
  <!-- class='w-28 md:w-33 lg:w-48'   Width of 28 by default, 32 on medium screens, and 48 on large screens -->
  <!-- 响应式设计：https://www.tailwindcss.cn/docs/responsive-design -->
  <!-- sm<640px  md<768px lg<1024px  lg<1280px 2xl<1536px-->
  <a :href="getBookURL()" :target="getTarget"
    class="relative w-32 h-44 mx-4 my-4 bg-gray-200 rounded shadow-xl hover:shadow-2xl ring-1 ring-gray-400 hover:ring hover:ring-blue-500 bg-top bg-cover"
    :style="setBackgroundImage()">

    <!-- 书籍类型图标 -->
    <SvgBookIcon :book_info="book_info" :childBookNum="childBookNum"></SvgBookIcon>
    <!-- 图书封面 -->
    <div v-if="showTitle"
      class="absolute inset-x-0 bottom-0 h-1/4 bg-gray-100 bg-opacity-80 font-semibold border-blue-800 rounded-b">
      <!-- 如果把链接的 target 属性设置为 "_blank"，该链接会在新窗口中打开。 -->
      <span class="absolute inset-x-0  font-bold top-0 p-1 align-middle">{{
        shortTitle
      }}</span>
    </div>
  </a>
</template>

<script lang="ts">
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
// import { NCard, } from 'naive-ui'
import { useCookies } from "vue3-cookies";
import { defineComponent } from "vue";
import SvgBookIcon from "@/components/SvgBookIcon.vue";

export default defineComponent({
  name: "BookCard",
  props: ["book_info", "readerMode", "showTitle", "simplifyTitle"],
  components: {
    SvgBookIcon,
  },
  setup() {
    const { cookies } = useCookies();
    return { cookies };
  },
  computed: {
    getTarget() {
      let bookType = this.book_info.type;
      if (
        bookType === ".pdf" ||
        bookType === "video" ||
        bookType === "audio" ||
        bookType === "unknown"
      ) {
        return "_blank";
      }
      return "_self";
    },
    childBookNum(): string {
      if (this.book_info.type === 'dir') {
        return ''
      }
      if (this.book_info.child_book_num > 0) {
        return "x" + this.book_info.child_book_num.toString();
      }
      return "";
    },
    shortTitle(): string {
      let short_title = this.book_info.title
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
    getBookURL() {
      let bookID = this.book_info.id;
      let bookType = this.book_info.type;
      let bookName = this.book_info.title;
      if (bookType === "book_group") {
        return "/#/child_shelf/" + bookID + "/";
      }
      if (
        bookType === ".pdf" ||
        bookType === "video" ||
        bookType === "audio" ||
        bookType === "unknown"
      ) {
        return "/api/raw/" + bookID + "/" + encodeURIComponent(bookName);
      }
      if (this.readerMode === "flip" || this.readerMode === "sketch") {
        return "/#/flip/" + bookID;
      }
      if (this.readerMode === "scroll") {
        // 命名路由,并加上参数,让路由建立 url
        return "/#/scroll/" + bookID;
      }
    },
    // 回首页
    onBackTop() {
      // 字符串路径
      this.$router.push("/");
    },
    setBackgroundImage() {
      return `background-image: url(${this.getThumbnailsImageUrl()});`;
    },
    getThumbnailsImageUrl() {
      // 按照“/”分割字符串
      const arrUrl = this.book_info.cover.url.split("/");
      // console.log(arrUrl)
      if (arrUrl[0] === "api") {
        return `${this.book_info.cover.url}&resize_width=256&resize_height=360&thumbnail_mode=true`;
      }
      return this.book_info.cover.url;
    },
  },
});
</script>
