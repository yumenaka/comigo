<template>
  <a class="bg-gray-100 my-4 flex flex-row justify-between  w-11/12  max-w-2xl rounded shadow-xl hover:shadow-2xl ring-1 ring-gray-400 hover:ring hover:ring-blue-500"
    :href="getBookCardURL()">
    <div
      class="bg-top bg-cover w-32 h-44 mx-4 my-4 bg-gray-200 rounded shadow-xl hover:shadow-2xl ring-1 ring-gray-400 hover:ring hover:ring-blue-500 "
      :style="setBackgroundImage()">
    </div>
    <div class="w-2/3 flex flex-col   top-0 p-4 align-middle  border-blue-800 rounded-b">
      <div class="font-bold text-xl">{{ shortTitle }}</div>
      <div v-if="book_info.child_book_num > 0" class="text-2xl font-black">{{
        bookNumHint
      }}</div>
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
  props: ["book_info", "bookCardMode", "readerMode", "showTitle",  "simplifyTitle"],
  components: {
  },
  setup() {
    const { cookies } = useCookies();
    return { cookies };
  },
  computed: {
    bookNumHint(): string {
      if (this.book_info.child_book_num > 0) {
        return `共${this.book_info.child_book_num}本`;
      }
      return "";
    },
    shortTitle(): string {
      let short_title = this.book_info.name
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
      return short_title;
    },
  },
  data() {
    return {
      test: "",
    };
  },
  methods: {
    getBookCardURL() {
      let bookID = this.book_info.id;
      let bookType = this.book_info.book_type;
      let bookName = this.book_info.name;
      // console.log("getBookCardOpenURL  bookID：" + bookID + " bookType：" + bookType)
      if (bookType === "book_group") {
        return "/#/child_shelf/" + bookID + "/";
      }
      if (
        bookType === ".pdf" ||
        bookType === "video" ||
        bookType === "audio" ||
        bookType === "unknown"
      ) {
        return "/api/raw/" + bookID + "/" + bookName;
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