<template>
  <!-- 外边距: m-2 https://www.tailwindcss.cn/docs/margin -->
  <!-- 内边距： p-4 https://www.tailwindcss.cn/docs/padding  p-0 m-0  -->
  <header class="header p-1 h-12 w-full flex justify-between content-center">
    <div> <!-- 返回箭头,点击返回上一页 -->
      <n-icon class="p-0 mx-1 my-0" v-if="showReturnIcon" size="40" @click="onClickReturnIcon()">
        <return-up-back />
      </n-icon>

      <!-- 服务器设置 -->
      <n-icon class="p-0 mx-1 my-0" v-if="!showReturnIcon" @click="ToAdminPage()" size="40">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"
          class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M5.25 14.25h13.5m-13.5 0a3 3 0 01-3-3m3 3a3 3 0 100 6h13.5a3 3 0 100-6m-16.5-3a3 3 0 013-3h13.5a3 3 0 013 3m-19.5 0a4.5 4.5 0 01.9-2.7L5.737 5.1a3.375 3.375 0 012.7-1.35h7.126c1.062 0 2.062.5 2.7 1.35l2.587 3.45a4.5 4.5 0 01.9 2.7m0 0a3 3 0 01-3 3m0 3h.008v.008h-.008v-.008zm0-6h.008v.008h-.008v-.008zm-3 6h.008v.008h-.008v-.008zm0-6h.008v.008h-.008v-.008z" />
        </svg>
      </n-icon>
      <!-- 上传按钮，点击进入上传页面 -->
      <n-icon class="p-0 mx-1 my-0" v-if="!showReturnIcon" @click="gotoUploadPage()" size="40">
        <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 512 512">
          <path
            d="M320 367.79h76c55 0 100-29.21 100-83.6s-53-81.47-96-83.6c-8.89-85.06-71-136.8-144-136.8c-69 0-113.44 45.79-128 91.2c-60 5.7-112 43.88-112 106.4s54 106.4 120 106.4h56"
            fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32"></path>
          <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32"
            d="M320 255.79l-64-64l-64 64"></path>
          <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32"
            d="M256 448.21V207.79"></path>
        </svg>
      </n-icon>

      <!-- 列表图标 -->
      <n-dropdown v-if="showReSortIcon" trigger="hover" :options="options" @select="onSelect">
        <n-icon class="w-10" size="40">
          <Filter />
        </n-icon>
      </n-dropdown>
    </div>

    <!-- 标题-->
    <!-- 文本颜色： https://www.tailwindcss.cn/docs/text-color -->
    <!-- 文本装饰（下划线）：https://www.tailwindcss.cn/docs/text-decoration -->
    <!-- 文本溢出：https://www.tailwindcss.cn/docs/text-overflow -->
    <!-- 字体粗细:https://www.tailwindcss.cn/docs/font-weight -->
    <div class="p-0 m-0 py-0 flex flex-col justify-center font-semibold content-center truncate">
      <!-- 标题，只显示 -->
      <span class="text-lg" v-if="!setDownLoadLink && (inShelf)">{{ headerTitle }}</span>
      <!-- 标题，可下载压缩包 -->
      <span class="text-lg text-blue-700 text-opacity-100  hover:underline">
        <a v-if="setDownLoadLink && (inShelf)" :href="'api/raw/' + bookID + '/' + encodeURIComponent(headerTitle)">{{
          headerTitle
        }}</a>
      </span>
      <!-- 快速跳转 -->
      <QuickJumpBar v-if="!inShelf" class="self-center" :nowBookID="bookID" :readMode="readMode"></QuickJumpBar>
    </div>
    <!-- slot，用来插入自定义组件。但是目前没需求 -->
    <!-- <slot></slot> -->

    <!-- 溢出 overflow-x-auton :https://www.tailwindcss.cn/docs/overflow -->
    <div class="p-0 h-10 w-33 flex justify-between content-center overflow-x-auton">
      <!-- QRCode图，点击可以在屏幕正中显示二维码 -->
      <Qrcode class="w-10 p-0"></Qrcode>
      <!-- 全屏图标 -->
      <svg class="w-10 static" @click="onFullScreen" xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 24 24">
        <g fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M16 4h4v4"></path>
          <path d="M14 10l6-6"></path>
          <path d="M8 20H4v-4"></path>
          <path d="M4 20l6-6"></path>
          <path d="M16 20h4v-4"></path>
          <path d="M14 14l6 6"></path>
          <path d="M8 4H4v4"></path>
          <path d="M4 4l6 6"></path>
        </g>
      </svg>
      <!-- 右边的设置图标,点击屏幕中央也可以打开  可自定义方向 -->
      <n-icon v-if="showSettingsIcon" class="w-10" size="40" @click="onClickSettingIcon('right')">
        <settings-outline />
      </n-icon>
    </div>
  </header>
</template>

<script lang="ts">
import { useCookies } from "vue3-cookies";
import { NIcon, NDropdown, useMessage, } from 'naive-ui'
import { ReturnUpBack, SettingsOutline, Grid, List, Filter, Text } from '@vicons/ionicons5'
import { h, defineComponent } from 'vue'
import Qrcode from "@/components/Qrcode.vue";
import screenfull from 'screenfull'
import QuickJumpBar from "@/components/QuickJumpBar.vue";
import axios from "axios";

export default defineComponent({
  name: "ComigoHeader",
  props: ['setDownLoadLink', 'headerTitle', 'bookID', 'showReturnIcon', 'showSettingsIcon', 'showReSortIcon', 'readMode', 'inShelf', 'depth'],
  emits: ['drawerActivate', 'onResort'],
  components: {
    NDropdown,//下拉菜单 https://www.naiveui.com/zh-CN/os-theme/components/dropdown
    NIcon,
    //图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）与导入
    ReturnUpBack,
    Grid,
    List,
    Filter,
    Text,
    SettingsOutline, //图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）
    Qrcode,//https://github.com/scopewu/qrcode.vue
    QuickJumpBar,
  },
  setup() {
    const { cookies } = useCookies();
    // console.log(window.history)
    //警告信息
    const message = useMessage();
    return { cookies, message, };
  },
  data() {
    return {
      resort_hint_key: "filename", //书籍的排序方式。可以按照文件名、修改时间、文件大小排序（或反向排序）
      options: [
        {
          label: '卡片模式',
          icon() {
            return h(NIcon, null, {
              default: () => h(Grid)
            })
          },
          key: 'gird'
        },
        {
          label: '列表模式',
          icon() {
            return h(NIcon, null, {
              default: () => h(List)
            })
          },
          key: "list"
        },
        {
          label: '文字模式',
          icon() {
            return h(NIcon, null, {
              default: () => h(Text)
            })
          },
          key: "text"
        },
        {
          type: 'divider',
          key: 'd1'
        },
        {
          label: this.$t("sort_by_filename"),
          key: 'filename'
        },
        {
          label: this.$t("sort_by_modify_time"),
          key: "modify_time"
        },
        {
          label: this.$t("sort_by_filesize"),
          key: 'filesize'
        },
        {
          label: this.$t("sort_by_filename_reverse") + this.$t("sort_reverse"),
          key: 'filename_reverse'
        },
        {
          label: this.$t("sort_by_modify_time_reverse") + this.$t("sort_reverse"),
          key: "modify_time_reverse"
        },
        {
          label: this.$t("sort_by_filesize_reverse") + this.$t("sort_reverse"),
          key: 'filesize_reverse'
        },
      ],
      handleSelect(key: string | number) {
        console.info(String(key))
      }
    };
  },
  methods: {
    //点击返回图标的时候，后退到上一页或主页
    onClickReturnIcon() {
      axios
        .get("/parent_book_info?id=" + this.bookID)
        .then((response) => {
          //请求接口成功，跳转到书架页面
          location.href = `/\#/child_shelf/${response.data.id}`;
        }).catch((error) => {
          //console.log("请求接口失败" + error);
          this.$router.push('/')
        });
    },
    //根据文件名、修改时间、文件大小等参数重新排序
    onSelect(key: string) {
      // console.info(key);
      this.$emit('onResort', key);
    },
    //进入全屏，由screenfull实现 https://github.com/sindresorhus/screenfull
    //全屏 API： https://developer.mozilla.org/zh-CN/docs/Web/API/Fullscreen_API
    onFullScreen() {
      //如果不允许进入全屏，发提示
      if (!screenfull.isEnabled) {
        this.message.warning(this.$t('not_support_fullscreen'))
        return false
      }
      //切换提示
      if (!screenfull.isFullscreen) {
        this.message.success(this.$t('success_fullScreen'));
      } else {
        this.message.success(this.$t('exit_fullScreen'));
      }
      //切换全屏状态
      screenfull.toggle()
    },

    //点击上传的时候，去上传页
    gotoUploadPage() {
      this.$router.push({
        name: "UploadPage"
      });
    },
    ToAdminPage() {
      location.href = "/admin";
    },

    //点击主页图标的时候，回到主页
    onClickSettingIcon(place: string) {
      this.$emit("drawerActivate", place);
    },

  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>




