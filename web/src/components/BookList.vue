<template>
  <a class="mx-2 my-2 flex flex-row justify-between  w-[32rem] h-[15rem] rounded bg-gray-100 shadow-xl hover:shadow-2xl ring-1 ring-gray-400 hover:ring hover:ring-blue-500"
    :href="getBookURL()" :target="getTarget">
    <div
      class="relative w-32 mx-4 my-4 bg-gray-200 bg-top bg-cover rounded shadow-xl h-44 hover:shadow-2xl ring-1 ring-gray-400 hover:ring hover:ring-blue-500 "
      :style="setBackgroundImage()">
    </div>
    <div class="top-0 flex flex-col w-2/3 p-4 my-2 border-blue-800 rounded-b">
      <div class="w-full my-1 text-xl font-bold text-left ">{{ book_info.title }}</div>
      <div class="w-full my-1 text-xl text-left" v-if="book_info.author !== ''">{{ $t('author', [book_info.author]) }}
      </div>
      <div class="w-full my-1 text-xl text-left" v-if="!isBookGroup && !isDirBook">{{ $t('filesize', [fileSizeString])
        }}
      </div>
      <div class="w-full my-1 text-xl text-left" v-if="!isBookGroup">{{ $t('allpagenum', [book_info.page_count]) }}
      </div>
      <div class="w-full my-1 text-xl text-left" v-if="book_info.child_book_num > 0">{{ bookNumHint }}</div>
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
  props: ["book_info", "bookCardMode", "readerMode", "showTitle", "simplifyTitle", "InfiniteDropdown"],
  components: {
  },
  setup() {
    const { cookies } = useCookies();
    return { cookies };
  },
  computed: {
    isBookGroup() {
      return this.book_info.type === 'book_group'
    },
    isDirBook() {
      return this.book_info.type === 'dir'
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
    bookNumHint(): string {
      if (this.book_info.child_book_num > 0) {
        return this.$t('childbookhint', [this.book_info.child_book_num]);
      }
      return "";
    },
    fileSizeString(): string {
      let bytes = this.book_info.file_size;
      if (bytes === 0) {
        return '0 Bytes';
      }
      const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
      const i = Math.floor(Math.log(bytes) / Math.log(1024));
      const formattedValue = parseFloat((bytes / Math.pow(1024, i)).toFixed(2));
      return `${formattedValue} ${sizes[i]}`;
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
        // short_title = short_title.replace(/.(zip|rar|cbr|cbz|tar|pdf|mp3|mp4|flv|gz|webm|gif|png|jpg|jpeg|webp|svg|psd|bmp|tif)/g, "");
        //域名（误伤多?）   参考了正则大全（https://any86.github.io/any-rule/）的网址(URL)
        // const domain_reg = /^(((ht|f)tps?):\/\/)?([^!@#$%^&*?.\s-]([^!@#$%^&*?.\s]{0,63}[^!@#$%^&*?.\s])?\.)+[a-zA-Z]{2,6}\/?/g;
        // short_title = short_title.replace(domain_reg, "");
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
    getBookURL() {
      let bookID = this.book_info.id;
      let bookType = this.book_info.type;
      let bookName = this.book_info.title;
      if (bookName === "Upload Book") {
        return "/#/upload";
      }
      // console.log("getBookCardOpenURL  bookID：" + bookID + " bookType：" + bookType)
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
      if (arrUrl[1]!==undefined&&arrUrl[1].includes("api")) {
        return `${this.book_info.cover.url}&resize_width=256&resize_height=360&thumbnail_mode=true`;
      }
      return this.book_info.cover.url;
    },
  },
});
</script>