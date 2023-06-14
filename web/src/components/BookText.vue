<template>
  <a class="m-2 bg-slate-300 text-lg" :href="openURL" :target="a_target">{{ shortTitle }}</a>
</template>

<script lang="ts">

import { defineComponent } from "vue";

export default defineComponent({
  name: "BookText",
  props: ["bookCardMode", "title", "image_src", "id", "readerMode", "showTitle", "childBookNum", "openURL", "a_target", "simplifyTitle"],
  components: {
  },
  setup() {
    return {};
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
        //short_title = short_title.replace(/.(zip|rar|cbr|cbz|tar|pdf|mp3|mp4|flv|gz|webm|gif|png|jpg|jpeg|webp|svg|psd|bmp|tif)/g, "");
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
  },
});
</script>

<style scoped></style>
