<template>
  <!-- 外边距: m-2 https://www.tailwindcss.cn/docs/margin -->
  <!-- 内边距： p-4 https://www.tailwindcss.cn/docs/padding  p-0 m-0  -->
  <header class="header p-1 h-12 w-full flex justify-between content-center">
    <div> <!-- 返回箭头,点击返回上一页 -->
      <n-icon class="p-0 m-0" v-if="showReturnIcon" size="40" @click="onClickReturnIcon()">
        <return-up-back />
      </n-icon>

      <!-- 上传按钮，点击进入上传页面 -->
      <n-icon class="p-0 m-0" v-if="!showReturnIcon" @click="gotoUploadPage()" size="40">
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
      <n-dropdown trigger="hover" :options="options" @select="onSelect">
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
    <div class="p-0 m-0 py-2 font-semibold content-center truncate">
      <!-- 标题，只显示 -->
      <span class="text-lg" v-if="!setDownLoadLink">{{ headerTitle }}</span>
      <!-- 标题，可下载压缩包 -->
      <span class="text-lg text-blue-700 text-opacity-100  hover:underline">
        <a v-if="setDownLoadLink" :href="'api/raw/' + bookID + '/' + headerTitle">{{ headerTitle }}</a>
      </span>
    </div>
    <!-- slot，用来插入自定义组件。但是目前没需求 -->
    <!-- <slot></slot> -->

    <!-- 溢出 overflow-x-auton :https://www.tailwindcss.cn/docs/overflow -->
    <div class="p-0 h-10 w-33 flex justify-between content-center overflow-x-auton">

      <!-- QRCode图，点击可以在屏幕正中显示二维码 -->
      <Qrcode class="w-10 p-0"></Qrcode>
      <!-- 全屏图标 -->
      <svg class="w-10 static" @click="onFullSreen" xmlns="http://www.w3.org/2000/svg"
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
import { ReturnUpBack, SettingsOutline, Grid, List, Filter } from '@vicons/ionicons5'
import { h, defineComponent } from 'vue'
import Qrcode from "@/components/Qrcode.vue";
import screenfull from 'screenfull'
export default defineComponent({
  name: "ComigoHeader",
  props: ['setDownLoadLink', 'headerTitle', 'bookID', 'showReturnIcon', 'showSettingsIcon',],
  emits: ['drawerActivate', 'onResort'],
  components: {
    NDropdown,//下拉菜单 https://www.naiveui.com/zh-CN/os-theme/components/dropdown
    NIcon,
    // BookOutline,//图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）与导入
    ReturnUpBack,
    Grid,
    List,
    Filter,
    SettingsOutline, //图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）
    Qrcode,//https://github.com/scopewu/qrcode.vue
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
          label: this.$t("sort_by_filename") + this.$t("sort_reverse"),
          key: 'filename_reverse'
        },
        {
          label: this.$t("sort_by_modify_time") + this.$t("sort_reverse"),
          key: "modify_time_reverse"
        },
        {
          label: this.$t("sort_by_filesize") + this.$t("sort_reverse"),
          key: 'filesize_reverse'
        },
      ],
      handleSelect(key: string | number) {
        console.info(String(key))
      }
    };
  },
  methods: {
    //根据文件名、修改时间、文件大小等参数重新排序
    onSelect(key: string) {
      // console.info(key);
      this.$emit('onResort', key);
    },
    //进入全屏，由screenfull实现 https://github.com/sindresorhus/screenfull
    //全屏 API： https://developer.mozilla.org/zh-CN/docs/Web/API/Fullscreen_API
    onFullSreen() {
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
    //点击返回图标的时候，后退到上一页或主页
    onClickReturnIcon() {
      // console.log(window.history)
      //如果直接进入本页面，没有上一页，那么回到主页。不过这时候浏览器back按钮本来应该也不能按。
      if (window.history.length === 1) {
        this.$router.push('/')
        return
      }
      //如果是新建标签页,并不是从书架点进去的时候，window.history.state.back为null，直接回到主页。
      if (window.history.state.back === null) {
        this.$router.push('/')
        return
      }
      //其他情况下，后退一页。与单击浏览器中的“后退”按钮相同。
      this.$router.back();
      // location.reload();
    },
    //点击上传的时候，去上传页
    gotoUploadPage() {
      this.$router.push({
        name: "UploadPage"
      });
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




