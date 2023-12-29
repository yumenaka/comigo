<template>
  <!-- class='w-28 md:w-33 lg:w-48'   Width of 28 by default, 32 on medium screens, and 48 on large screens -->
  <!-- 响应式设计：https://www.tailwindcss.cn/docs/responsive-design -->
  <!-- sm<640px  md<768px lg<1024px  lg<1280px 2xl<1536px-->
  <a :href="getBookURL()" :target="getTarget"
    class="relative w-32 h-44 mx-4 my-4 bg-gray-200 rounded shadow-xl hover:shadow-2xl ring-1 ring-gray-400 hover:ring hover:ring-blue-500 bg-top bg-cover"
    :style="setBackgroundImage()">
    <!-- 书籍组：显示书籍数量 -->
    <div v-if="childBookNum != ''" class="flex flex-row justify-end pl-2">
      <svg class="opacity-70 w-7 h-7" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="rgb(234 179 8)" stroke-width="1.0"
        stroke="currentColor">
        <path
          d="M11.25 4.533A9.707 9.707 0 0 0 6 3a9.735 9.735 0 0 0-3.25.555.75.75 0 0 0-.5.707v14.25a.75.75 0 0 0 1 .707A8.237 8.237 0 0 1 6 18.75c1.995 0 3.823.707 5.25 1.886V4.533ZM12.75 20.636A8.214 8.214 0 0 1 18 18.75c.966 0 1.89.166 2.75.47a.75.75 0 0 0 1-.708V4.262a.75.75 0 0 0-.5-.707A9.735 9.735 0 0 0 18 3a9.707 9.707 0 0 0-5.25 1.533v16.103Z" />
      </svg>

      <span class="text-2xl text-yellow-500 font-black text-shadow">{{ childBookNum }}</span>
    </div>

    <!-- 文件夹 -->
    <svg v-if='book_info.type == "dir"' class="opacity-70 w-7 h-7 m-0 absolute top-1 right-1 shadow-2xl hover:shadow-2xl"
      xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="rgb(234 179 8)" stroke-width="1.0"
      stroke="currentColor">
      <path
        d="M19.5 21a3 3 0 0 0 3-3v-4.5a3 3 0 0 0-3-3h-15a3 3 0 0 0-3 3V18a3 3 0 0 0 3 3h15ZM1.5 10.146V6a3 3 0 0 1 3-3h5.379a2.25 2.25 0 0 1 1.59.659l2.122 2.121c.14.141.331.22.53.22H19.5a3 3 0 0 1 3 3v1.146A4.483 4.483 0 0 0 19.5 9h-15a4.483 4.483 0 0 0-3 1.146Z" />
    </svg>

    <!-- 视频  -->
    <svg v-if='book_info.type == "video"'
      class="opacity-70 w-7 h-7 m-0 absolute top-1 right-1 shadow-2xl hover:shadow-2xl" xmlns="http://www.w3.org/2000/svg"
      fill="rgb(234 179 8)" viewBox="0 0 24 24" stroke-width="1.0" stroke="currentColor">
      <path stroke-linecap="round" stroke-linejoin="round"
        d="m15.75 10.5 4.72-4.72a.75.75 0 0 1 1.28.53v11.38a.75.75 0 0 1-1.28.53l-4.72-4.72M4.5 18.75h9a2.25 2.25 0 0 0 2.25-2.25v-9a2.25 2.25 0 0 0-2.25-2.25h-9A2.25 2.25 0 0 0 2.25 7.5v9a2.25 2.25 0 0 0 2.25 2.25Z" />
    </svg>

    <!-- 音乐  -->
    <svg v-if='book_info.type == "audio"'
      class="opacity-70 w-7 h-7 m-0 absolute top-1 right-1 shadow-2xl hover:shadow-2xl" xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24" fill="rgb(234 179 8)" stroke-width="1.0" stroke="currentColor">
      <path fill-rule="evenodd"
        d="M19.952 1.651a.75.75 0 0 1 .298.599V16.303a3 3 0 0 1-2.176 2.884l-1.32.377a2.553 2.553 0 1 1-1.403-4.909l2.311-.66a1.5 1.5 0 0 0 1.088-1.442V6.994l-9 2.572v9.737a3 3 0 0 1-2.176 2.884l-1.32.377a2.553 2.553 0 1 1-1.402-4.909l2.31-.66a1.5 1.5 0 0 0 1.088-1.442V5.25a.75.75 0 0 1 .544-.721l10.5-3a.75.75 0 0 1 .658.122Z"
        clip-rule="evenodd" />
    </svg>

    <!-- PDF  -->
    <svg v-if='book_info.type == ".pdf"' class="opacity-70 w-7 h-7 m-0 absolute top-1 right-1 shadow-2xl hover:shadow-2xl"
      xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="rgb(234 179 8)" stroke-width="1.0"
      stroke="currentColor">
      <path
        d="M5.625 1.5c-1.036 0-1.875.84-1.875 1.875v17.25c0 1.035.84 1.875 1.875 1.875h12.75c1.035 0 1.875-.84 1.875-1.875V12.75A3.75 3.75 0 0 0 16.5 9h-1.875a1.875 1.875 0 0 1-1.875-1.875V5.25A3.75 3.75 0 0 0 9 1.5H5.625Z" />
      <path
        d="M12.971 1.816A5.23 5.23 0 0 1 14.25 5.25v1.875c0 .207.168.375.375.375H16.5a5.23 5.23 0 0 1 3.434 1.279 9.768 9.768 0 0 0-6.963-6.963Z" />
    </svg>

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

export default defineComponent({
  name: "BookCard",
  props: ["book_info", "readerMode", "showTitle", "simplifyTitle"],
  components: {
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
