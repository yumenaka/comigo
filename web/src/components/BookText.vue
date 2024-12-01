<template>
  <div class="p-2 mx-2 my-1 border rounded-lg bg-slate-300">
    <a class="flex flex-row items-center justify-center text-lg font-medium" :href="getBookURL()" :target="getTarget">
      <!-- 书籍类型图标 -->

      <span>{{ getIcon }}{{ shortTitle }}</span>
    </a>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "BookText",
  props: ["book_info", "bookCardMode", "readerMode", "simplifyTitle", "InfiniteDropdown"],
  components: {},
  setup() {
    return {};
  },
  computed: {
    getIcon() {
      const bookType = this.book_info.type;
      switch (bookType) {
        case "dir":
          return "📁";
        case "audio":
          return "🎵";
        case "book_group":
          return "📚";
        case ".pdf":
          return "📄";
        case "video":
          return "💽";
        case ".zip" || ".rar" || ".cbr" || ".cbz" || ".tar" || ".gz":
          return "💼";
        default:
          return "❔";
      }
    },
    getTarget() {
      let bookType = this.book_info.type;
      if (
        bookType === "video" ||
        bookType === "audio" ||
        bookType === "unknown"
      ) {
        return "_blank";
      }
      return "_self";
    },
    shortTitle(): string {
      let short_title = this.book_info.title;
      //使用 JavaScript replace() 方法替换掉一些字符串
      if (this.simplifyTitle) {
        //中：/[\u4e00-\u9fa5]/  日：/[\u0800-\u4e00]/  韩：/[\uac00-\ud7ff]/  空格：[\s]
        //左半部分：任意中日韩字符或空格，g表示多次匹配、不限次数
        short_title = short_title.replace(
          /[\\[\\(（【][A-Za-z0-9_\-×\s+\u4e00-\u9fa5\u0800-\u4e00\uac00-\ud7ff]+/g,
          ""
        );
        //右半部分
        short_title = short_title.replace(/[\]）】\\)]/g, "");
        //.zip .rar 等文件名
        //short_title = short_title.replace(/.(zip|rar|cbr|cbz|tar|pdf|mp3|mp4|flv|gz|webm|gif|png|jpg|jpeg|webp|svg|psd|bmp|tif)/g, "");
        //域名（误伤多?）   参考了正则大全（https://any86.github.io/any-rule/）的网址(URL)
        const domain_reg =
          /^(((ht|f)tps?):\/\/)?([^!@#$%^&*?.\s-]([^!@#$%^&*?.\s]{0,63}[^!@#$%^&*?.\s])?\.)+[a-zA-Z]{2,6}\/?/g;
        short_title = short_title.replace(domain_reg, "");
        //开头的空格
        short_title = short_title.replace(/^[\s]/g, "");
        //开头的特殊字符
        short_title = short_title.replace(
          /^[\\\-`~!@#$^&*()=|{}':;'@#￥……&*（）——|{}‘；：”“'。，、？]/,
          ""
        );
      }
      return short_title;
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
      if (bookName === "Upload Book") {
        return "/#/upload";
      }
      if (bookType === "book_group") {
        return "/#/child_shelf/" + bookID + "/";
      }
      if (
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
        let query_string = "";
        if (this.InfiniteDropdown === false) {
          query_string = "?page=1"
        }
        // 命名路由,并加上参数,让路由建立 url
        return "/#/scroll/" + bookID + query_string;
      }
    },
  },
});
</script>

<style scoped></style>
