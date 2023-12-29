<template>
  <div class="my-1 mx-2 p-2 border rounded-lg bg-slate-300">
    <a class="font-medium text-lg flex flex-row justify-center items-center" :href="getBookURL()" :target="getTarget">
      <svg v-if='book_info.type == "dir"' class="w-6 h-6 pr-1" xmlns="http://www.w3.org/2000/svg" fill="none"
        viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round"
          d="M2.25 12.75V12A2.25 2.25 0 0 1 4.5 9.75h15A2.25 2.25 0 0 1 21.75 12v.75m-8.69-6.44-2.12-2.12a1.5 1.5 0 0 0-1.061-.44H4.5A2.25 2.25 0 0 0 2.25 6v12a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9a2.25 2.25 0 0 0-2.25-2.25h-5.379a1.5 1.5 0 0 1-1.06-.44Z" />
      </svg>

      <svg v-if='book_info.type == "video"' xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
        stroke-width="1.5" stroke="currentColor" class="w-6 h-6  pr-1">
        <path stroke-linecap="round" stroke-linejoin="round"
          d="m15.75 10.5 4.72-4.72a.75.75 0 0 1 1.28.53v11.38a.75.75 0 0 1-1.28.53l-4.72-4.72M4.5 18.75h9a2.25 2.25 0 0 0 2.25-2.25v-9a2.25 2.25 0 0 0-2.25-2.25h-9A2.25 2.25 0 0 0 2.25 7.5v9a2.25 2.25 0 0 0 2.25 2.25Z" />
      </svg>

      <svg v-if='book_info.type == "audio"' xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
        stroke-width="1.5" stroke="currentColor" class="w-6 h-6 pr-1">
        <path stroke-linecap="round" stroke-linejoin="round"
          d="m9 9 10.5-3m0 6.553v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 1 1-.99-3.467l2.31-.66a2.25 2.25 0 0 0 1.632-2.163Zm0 0V2.25L9 5.25v10.303m0 0v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 0 1-.99-3.467l2.31-.66A2.25 2.25 0 0 0 9 15.553Z" />
      </svg>
      <span>{{ shortTitle }}</span>
    </a>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "BookText",
  props: ["book_info", "bookCardMode", "readerMode", "simplifyTitle"],
  components: {},
  setup() {
    return {};
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
  },
});
</script>

<style scoped></style>
