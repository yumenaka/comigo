<template>
  <div class="">
    <!-- 顶部,标题页头 -->
    <!-- 定位：https://www.tailwindcss.cn/docs/position -->
    <!-- 使用 fixed 来定位一个元素相对于浏览器窗视口的位置。偏移量是相对于视口计算的，且该元素将作为绝对定位的子元素的位置参考。 -->
    <!-- 控制 flex 和 grid 项目如何沿着容器的主轴定位:https://www.tailwindcss.cn/docs/justify-content -->
    <!-- 不透明度：opacity-70 https://www.tailwindcss.cn/docs/opacity -->
    <!-- Vue组建过渡：https://v3.cn.vuejs.org/guide/transitions-enterleave.html#%E5%8D%95%E5%85%83%E7%B4%A0-%E7%BB%84%E4%BB%B6%E7%9A%84%E8%BF%87%E6%B8%A1= -->
    <transition name="header-bottom">
      <Header v-if="this.showHeaderFlag_FlipMode" class="mx-auto w-full opacity-80"
        v-bind:class="{ 'fixed': this.hideToolbar, absolute: this.hideToolbar, 'top-0': this.hideToolbar }"
        v-bind:style="{ background: model.interfaceColor }" :setDownLoadLink="this.needDownloadLink()"
        :headerTitle="book.name" :bookID="this.book.id" :showReturnIcon="true" @drawerActivate="this.drawerActivate">
      </Header>
    </transition>

    <div class="main" v-bind:style="{ overflow: 'hidden' }">
      <!-- 主题,漫画div -->
      <!-- 事件修饰符： https://v3.cn.vuejs.org/guide/events.html#%E4%BA%8B%E4%BB%B6%E4%BF%AE%E9%A5%B0%E7%AC%A6 -->
      <div class="manga_area" id="MangaMain" @click.stop="onMouseClick" @mousemove.stop="onMouseMove"
        @mouseleave.stop="onMouseLeave">
        <div class="manga_area_img_div">
          <!-- 非自动拼合模式最简单,直接显示一张图 -->
          <img class="w-auto h-auto" v-bind:src="
            this.imageParametersString(book.pages.images[nowPageNum - 1].url)
          " v-bind:alt="nowPageNum" />

          <!-- 简单拼合双页,不管单双页什么的 -->
          <img v-if="
            !this.autoDoublePageModeFlag &&
            this.simpleDoublePageModeFlag &&
            this.nowPageNum < this.book.all_page_num
          " v-bind:src="
  this.imageParametersString(book.pages.images[nowPageNum].url)
" v-bind:alt="nowPageNum + 1" />

          <!-- 自动拼合模式当前页,如果开启自动拼合,右边可能显示拼合页 -->
          <img v-if="
            this.autoDoublePageModeFlag &&
            this.nowPageNum < this.book.all_page_num &&
            this.nowAndNextPageIsSingle()
          " v-bind:src="
  this.imageParametersString(book.pages.images[nowPageNum].url)
" v-bind:alt="nowPageNum + 1" />
        </div>
      </div>
    </div>
    <!-- 页数、倒计时文字提示 -->
    <!-- Top / Right / Bottom / Left :用于控制定位元素的位置的功能类。https://www.tailwindcss.cn/docs/top-right-bottom-left -->
    <div class="
        font-sans
        text-blue-700 text-opacity-90
        shadow-lg shadow-blue-500/50
        h-auto
        w-full
        opacity-70
        absolute
        bottom-0
        fixed
      " v-bind:class="{ 'text-2xl': this.sketchModeFlag, 'text-lg': (!this.sketchModeFlag) }"
      v-if="this.showPageHintFlag_FlipMode">
      {{ pageNumOrSketchHint }}
    </div>

    <!-- 底部的阅读进度条 -->
    <transition name="header-bottom">
      <!-- 进度条 -->
      <!-- 宽度：w-5/6 https://www.tailwindcss.cn/docs/width 使用 w-{fraction} 或 w-full 将元素设置为基于百分比的宽度。 -->
      <!-- 定位：https://www.tailwindcss.cn/docs/position  -->
      <!-- 使用 fixed 来定位一个元素相对于浏览器窗视口的位置。偏移量是相对于视口计算的，且该元素将作为绝对定位的子元素的位置参考。 -->
      <!-- 控制 flex 和 grid 项目如何沿着容器的主轴定位:https://www.tailwindcss.cn/docs/justify-content -->
      <!-- Tailwind 的容器不会自动居中，也没有任何内置的水平方向的内边距。要使一个容器居中，使用 mx-auto 功能类： -->
      <!-- absolute bottom-0  -->
      <!-- v-bind:class="{ absolute: this.hideToolbar, 'bottom-0': this.hideToolbar, 'fixed': this.hideToolbar, 'flex': !this.hideToolbar}" -->

      <div class="m-0 w-full h-10 opacity-80 overflow-hidden"
        v-bind:class="{ absolute: this.hideToolbar, 'bottom-0': this.hideToolbar, 'flex': (!this.hideToolbar) }"
        v-if="this.showFooterFlag_FlipMode">
        <div class="
          bg-yellow-400
          flex flex-row
          justify-center
          items-end
          mx-auto
          w-full
          h-10
        ">
          <!-- 底部滑动条,日漫模式，reverse翻转计数方向 -->
          <!-- 背景颜色：bg-blue-300  https://www.tailwindcss.cn/docs/background-color  -->
          <div class="
            bg-blue-300
            flex flex-row
            justify-center
            items-center
            w-5/6
            px-4
            h-full
          " v-if="!this.rightToLeftFlag">
            <span class="right">{{ this.book.all_page_num }}</span>
            <n-slider class="w-10/11" reverse v-model:value="nowPageNum" :max="this.book.all_page_num" :min="1"
              :step="1" :format-tooltip="(value) => `${value}`" @update:value="this.saveNowPageNumOnUpdate" />
            <span class="left">{{ this.nowPageNum }}</span>
          </div>

          <!-- 底部滑动条，美漫模式 -->
          <!-- h-full: 将一个元素的高度设置为其父元素的 100%，只要父元素有一个定义的高度。 https://www.tailwindcss.cn/docs/height -->
          <!-- 使用 items-center 沿着容器的交叉轴中心对齐项目：https://www.tailwindcss.cn/docs/align-items#center=  -->
          <div class="
            bg-blue-300
            flex flex-row
            justify-center
            items-center
            w-5/6
            px-4
            h-full
          " v-if="this.rightToLeftFlag">
            <span class="right">{{ this.nowPageNum }}</span>
            <n-slider class="bg-yellow-300" v-model:value="nowPageNum" :max="this.book.all_page_num" :min="1" :step="1"
              :format-tooltip="(value) => `${value}`" @update:value="this.saveNowPageNumOnUpdate" />
            <span class="left">{{ this.book.all_page_num }}</span>
          </div>
        </div>
      </div>
    </transition>

    <!-- 设置抽屉,一开始隐藏 -->
    <Drawer :initDrawerActive="this.drawerActive" :initDrawerPlacement="this.drawerPlacement"
      @saveConfig="this.saveConfigToLocal" @startSketch="this.startSketchMode" @stopSketch="this.stopSketchMode"
      @closeDrawer="this.drawerDeactivate" :readerMode="this.readerMode" :inBookShelf="false"
      :sketching="this.sketchModeFlag">
      <!-- 选择：切换页面模式 -->
      <n-space>
        <n-button @click="changeReaderModeToScrollMode">{{
            $t("switch_to_scrolling_mode")
        }}
        </n-button>
      </n-space>
      <!-- 分割线 -->
      <p> &nbsp;</p>

      <!-- 显示当前页数 -->
      <n-space>
        <n-switch size="large" v-model:value="this.showPageHintFlag_FlipMode" @update:value="setShowPageNumChange">
          <template #checked>{{ $t("showPageNum") }}</template>
          <template #unchecked>{{ $t("showPageNum") }}</template>
        </n-switch>
      </n-space>

      <!-- websocket同步 -->
      <n-space>
        <n-switch size="large" v-model:value="this.syncPageFlag" @update:value="this.setSyncPageFlag">
          <template #checked>{{ $t("sync_page") }}</template>
          <template #unchecked>{{ $t("sync_page") }}</template>
        </n-switch>
      </n-space>
      <!-- 保存阅读的页数 -->
      <n-space>
        <n-switch size="large" v-model:value="this.saveNowPageNumFlag" @update:value="this.setSavePageNumFlag">
          <template #checked>{{ $t("savePageNum") }}</template>
          <template #unchecked>{{ $t("savePageNum") }}</template>
        </n-switch>
      </n-space>

      <!-- 合并双页 -->
      <n-space>
        <n-switch size="large" v-model:value="this.simpleDoublePageModeFlag"
          @update:value="this.setSimpleDoublePage_FlipMode">
          <template #checked>{{ $t("simpleDoublePage") }}</template>
          <template #unchecked>{{ $t("simpleDoublePage") }}</template>
        </n-switch>
      </n-space>
      <!-- 自动切边 -->
      <n-space>
        <n-switch size="large" v-model:value="this.imageParameters.do_auto_crop"
          @update:value="setImageParameters_DoAutoCrop">
          <template #checked>{{ $t("auto_crop") }}</template>
          <template #unchecked>{{ $t("auto_crop") }}</template>
        </n-switch>
        <!-- 切白边阈值 -->
        <n-input-number :show-button="false" v-if="this.imageParameters.do_auto_crop"
          v-model:value="this.imageParameters.auto_crop_num" @update:value="setImageParameters_AutoCropNum" :max="10"
          :min="0">
          <template #prefix>{{ $t("energy_threshold") }}</template>
        </n-input-number>
      </n-space>

      <!-- 翻页模式,默认右开本（日漫）-->
      <n-space>
        <n-switch size="large" v-model:value="this.rightToLeftFlag" :rail-style="railStyle"
          @update:value="this.setFlipScreenFlag">
          <template #checked>{{ $t("leftScreenToNext") }}</template>
          <template #unchecked>{{ $t("rightScreenToNext") }}</template>
        </n-switch>
      </n-space>

      <!-- 自动隐藏工具栏 -->
      <n-space>
        <n-switch size="large" v-model:value="this.hideToolbar" @update:value="this.setHideToolbar">
          <template #checked>{{ $t('auto_hide_toolbar') }}</template>
          <template #unchecked>{{ $t('auto_hide_toolbar') }}</template>
        </n-switch>
      </n-space>

      <!-- 分割线 -->
      <!-- <n-divider /> -->

      <!-- Debug,开启一些不稳定功能 -->
      <!-- <n-space>
      <n-switch size="large" v-model:value="this.debugModeFlag" @update:value="this.setDebugModeFlag">
        <template #checked>{{ $t("debugMode") }}</template>
        <template #unchecked>{{ $t("debugMode") }}</template>
      </n-switch>
    </n-space>

    <n-space v-if="this.debugModeFlag">
      <n-switch
        size="large"
        v-model:value="this.autoDoublePageModeFlag"
        @update:value="this.setAutoDoublePage_FlipMode"
      >
        <template #checked>{{ $t('autoDoublePage') }}</template>
        <template #unchecked>{{ $t('autoDoublePage') }}</template>
      </n-switch>
    </n-space>-->

      <!-- 分割线 -->
      <n-divider v-if="this.readerMode === 'sketch'" />
      <!-- 自动翻页秒数 -->
      <!-- 数字输入% -->
      <n-input-number v-if="this.readerMode === 'sketch'" size="small" :show-button="false"
        v-model:value="this.sketchFlipSecond" :max="65535" :min="1" :update-value-on-input="false"
        @update:value="this.resetSketchSecondCount">
        <template #prefix>{{ $t("pageTurningSeconds") }}</template>
        <template #suffix>{{ $t("second") }}</template>
      </n-input-number>
      <!-- 滑动选择% -->
      <n-slider v-if="this.readerMode === 'sketch'" v-model:value="this.sketchFlipSecond" :step="1" :max="120" :min="1"
        :marks="marks" :format-tooltip="(value) => `${value}s`" @update:value="this.resetSketchSecondCount" />

      <!-- 分割线 -->
      <n-divider />
      <!-- 页面排序方式 -->
      <n-select placeholder=this.getSortHintText() @update:value="this.onResort" v-model:value="this.resort_hint_key"
        :options="options" />
    </Drawer>
  </div>
</template>

<script>
import { useCookies } from "vue3-cookies";
// 自定义组件
import Header from "@/components/Header.vue";
import Drawer from "@/components/Drawer.vue";
// import Bottom from "@/components/Bottom.vue";
import { defineComponent, reactive } from "vue";
// 直接导入组件并使用它。这种情况下,只有导入的组件才会被打包。
import {
  NDivider,
  NInputNumber,
  NSlider,
  NSpace,
  NSwitch,
  useMessage,
  NButton,
  NSelect,
} from "naive-ui";
import axios from "axios";
import md5 from "js-md5";

export default defineComponent({
  name: "FlipMode",
  props: [],
  emits: ["childMark"], // 向父组件传递参数的函数,用法： 子组件调用：this.$emit("childMark", value); 父组件：@childMark="this.fatherMethod"
  components: {
    Header,
    Drawer,
    // Bottom,
    NSpace, //间距 https://www.naiveui.com/zh-CN/os-theme/components/space
    NSlider, //滑动选择  Slider https://www.naiveui.com/zh-CN/os-theme/components/slider
    NSwitch, //开关   https://www.naiveui.com/zh-CN/os-theme/components/switch
    // NLayout,//布局 https://www.naiveui.com/zh-CN/os-theme/components/layout
    // NLayoutSider,
    // NLayoutContent,
    // NColorPicker, //颜色选择器 Color Picker https://www.naiveui.com/zh-CN/os-theme/components/color-picker
    NDivider, //分割线  https://www.naiveui.com/zh-CN/os-theme/components/divider
    NInputNumber, ///  https://www.naiveui.com/zh-CN/os-theme/components/input-number
    // useNotification, // https://www.naiveui.com/zh-CN/os-theme/components/notification
    NButton, //按钮，来自:https://www.naiveui.com/zh-CN/os-theme/components/button
    // NMessageProvider,
    // NDropdown, //下拉菜单 https://www.naiveui.com/zh-CN/os-theme/components/dropdown
    NSelect, //选择器 https://www.naiveui.com/zh-CN/os-theme/components/select
  },
  setup() {
    const { cookies } = useCookies();
    //背景颜色,颜色选择器用
    const model = reactive({
      backgroundColor: "#E0D9CD",
      interfaceColor: "#f5f5e4",
    });
    //请求图片文件时，可添加的额外参数
    const imageParameters = reactive({
      resize_width: -1, // 缩放图片,指定宽度
      resize_height: -1, // 指定高度,缩放图片
      do_auto_resize: false,
      resize_max_width: 800, //图片宽度大于这个上限时缩小
      resize_max_height: -1, //图片高度大于这个上限时缩小
      do_auto_crop: false,
      auto_crop_num: 1, // 自动切白边阈值,范围是0~100,其实为1就够了
      gray: false, //黑白化
    });

    //警告信息
    const message = useMessage();
    // const notification = useNotification()
    return {
      message,
      //背景色
      model,
      cookies,
      imageParameters, //获取图片所用的参数
      imageParametersString: (source_url) => {
        // var temp =
        if (source_url.substr(0, 12) === "api/getfile?") {
          //当前URL
          const url = document.location.toString();
          //按照“/”分割字符串
          const arrUrl = url.split("/");
          //拼一个完整的图片URL（因为路由路径会变化,所以不能用相对路径？）
          const base_str = arrUrl[0] + "//" + arrUrl[2] + "/" + source_url;
          //添加各种字符串参数,不需要的话为空
          const resize_width_str =
            imageParameters.resize_width > 0
              ? "&resize_width=" + imageParameters.resize_width
              : "";
          const resize_height_str =
            imageParameters.resize_height > 0
              ? "&resize_height=" + imageParameters.resize_height
              : "";
          const gray_str = imageParameters.gray ? "&gray=true" : "";
          const do_auto_resize_str = imageParameters.do_auto_resize
            ? "&resize_max_width=" + imageParameters.resize_max_width
            : "";
          const resize_max_height_str =
            imageParameters.resize_max_height > 0
              ? "&resize_max_height=" + imageParameters.resize_max_height
              : "";
          const auto_crop_str = imageParameters.do_auto_crop
            ? "&auto_crop=" + imageParameters.auto_crop_num
            : "";
          //所有附加的转换参数
          const addStr =
            resize_width_str +
            resize_height_str +
            do_auto_resize_str +
            resize_max_height_str +
            auto_crop_str +
            gray_str;
          //如果有附加转换参数，则设置成不缓存
          const nocache_str = addStr === "" ? "" : "&no-cache=true";
          return base_str + addStr + nocache_str;
        } else {
          return source_url;
        }
      },
      //开关的颜色
      railStyle: ({ focused, checked }) => {
        const style = {};
        if (checked) {
          style.background = "#18a058";
          if (focused) {
            style.boxShadow = "0 0 0 2px #d0305040";
          }
        } else {
          style.background = "#2080f0";
          if (focused) {
            style.boxShadow = "0 0 0 2px #2080f040";
          }
        }
        return style;
      },
      //滑动秒数建议值
      marks: {
        30: "30",
        60: "60",
        90: "90",
        120: "120",
      },
    };
  },
  data() {
    return {
      //自动隐藏工具条
      hideToolbar: true,
      resort_hint_key: "resort",
      options: [
        {
          label: this.$t("sort_by_filename"),
          value: "filename",
        },
        {
          label: this.$t("sort_by_modify_time"),
          value: "modify_time",
        },
        {
          label: this.$t("sort_by_filesize"),
          value: "filesize",
        },
        {
          label: this.$t("sort_by_filename") + this.$t("sort_reverse"),
          value: "filename_reverse",
        },
        {
          label: this.$t("sort_by_modify_time") + this.$t("sort_reverse"),
          value: "modify_time_reverse",
        },
        {
          label: this.$t("sort_by_filesize") + this.$t("sort_reverse"),
          value: "filesize_reverse",
        },
      ],
      book: {
        name: "loading",
        id: "loading",
        all_page_num: 2,
        book_type: "dir",
        pages: {
          sort_by: "",
          images: [
            {
              height: 500,
              width: 449,
              url: "/images/loading.gif",
            },
            {
              height: 500,
              width: 449,
              url: "/images/loading.gif",
            },
          ],
        },
      },
      readerMode: "flip",
      drawerActive: false,
      drawerPlacement: "right",
      //开发模式 未完成的功能与设置,开启Debug以后才能见到
      debugModeFlag: true,
      //是否通过websocket同步翻页
      syncPageFlag: true,
      //是否显示页头
      showHeaderFlag_FlipMode: true,
      //是否显示页脚
      showFooterFlag_FlipMode: true,
      //是否是右半屏翻页（从右到左）?日本漫画从左到右(false)
      rightToLeftFlag: false,
      //简单拼合双叶
      simpleDoublePageModeFlag: false,
      //自动拼合双叶,效果不太好
      autoDoublePageModeFlag: false,
      //是否保存当前页数
      saveNowPageNumFlag: true,
      //当前页数,注意语义,直接就是1开始的页数,不是数组下标,在pages数组当中用的时候需要-1
      nowPageNum: 1,
      //素描模式标记
      sketchModeFlag: false,
      //是否显示素描提示
      showPageHintFlag_FlipMode: false,
      //翻页间隔时间
      sketchFlipSecond: 30,
      //计时用,从0开始
      sketchSecondCount: 0,
    };
  },
  //在选项API中使用 Vue 生命周期钩子：
  created() {
    // 消息监听，即接收websocket服务端推送的消息. optionsAPI用法
    this.$options.sockets.onmessage = (data) => this.handlePacket(data);

    //根据文件名、修改时间、文件大小等要素排序的参数
    let sort_image_by_str = "";
    if (this.$route.query.sort_by) {
      sort_image_by_str = "&sort_by=" + this.$route.query.sort_by;
    }

    //根据路由参数获取特定书籍
    axios
      .get("/getbook?id=" + this.$route.params.id + sort_image_by_str)
      .then((response) => (this.book = response.data))
      .finally(() => {
        document.title = this.book.name;
        console.log("成功获取书籍数据,书籍ID:" + this.$route.params.id);
      });
    //监听路由参数的变化,刷新本地的Book数据
    this.$watch(
      () => this.$route.params.id,
      (id) => {
        // console.log(id)
        axios
          .get("/getbook?id=" + this.$route.params.id + sort_image_by_str)
          .then((response) => (this.book = response.data))
          .finally(console.log("路由参数改变,书籍ID:" + id));
      }
    );
    //初始化默认值
    // https://www.developers.pub/wiki/1006381/1013545
    // https://developer.mozilla.org/zh-CN/docs/Web/API/Storage/setItem
    //一个域名下存放的cookie的个数有限制,不同的浏览器存放的个数不一样,一般为20个。因为不需要上传,使用localStorage（本地存储）存储在浏览器,永不过期。
    if (localStorage.getItem("debugModeFlag") === "true") {
      this.debugModeFlag = true;
    } else if (localStorage.getItem("debugModeFlag") === "false") {
      this.debugModeFlag = false;
    }
    //是否显示标题
    if (localStorage.getItem("showHeaderFlag_FlipMode") === "true") {
      this.showHeaderFlag_FlipMode = true;
    } else if (localStorage.getItem("showHeaderFlag_FlipMode") === "false") {
      this.showHeaderFlag_FlipMode = false;
    }
    //是否显示页脚
    if (localStorage.getItem("showFooterFlag_FlipMode_FlipMode") === "true") {
      this.showFooterFlag_FlipMode = true;
    } else if (
      localStorage.getItem("showFooterFlag_FlipMode_FlipMode") === "false"
    ) {
      this.showFooterFlag_FlipMode = false;
    }
    //是否显示页数
    if (localStorage.getItem("showPageHintFlag_FlipMode") === "true") {
      this.showPageHintFlag_FlipMode = true;
    } else if (localStorage.getItem("showPageHintFlag_FlipMode") === "false") {
      this.showPageHintFlag_FlipMode = false;
    }
    //翻页方向、是否用右半屏翻页
    if (localStorage.getItem("rightToLeftFlag") === "true") {
      this.rightToLeftFlag = true;
    } else if (localStorage.getItem("rightToLeftFlag") === "false") {
      this.rightToLeftFlag = false;
    }
    //简单合并单页
    if (localStorage.getItem("simpleDoublePageModeFlag") === "true") {
      this.simpleDoublePageModeFlag = true;
    } else if (localStorage.getItem("simpleDoublePageModeFlag") === "false") {
      this.simpleDoublePageModeFlag = false;
    }
    //自动合并单页
    if (localStorage.getItem("autoDoublePageModeFlag") === "true") {
      this.autoDoublePageModeFlag = true;
    } else if (localStorage.getItem("autoDoublePageModeFlag") === "false") {
      this.autoDoublePageModeFlag = false;
    }
    //当前背景色
    if (localStorage.getItem("BackgroundColor") != null) {
      this.model.backgroundColor = localStorage.getItem("BackgroundColor");
    }
    if (localStorage.getItem("InterfaceColor") != null) {
      this.model.interfaceColor = localStorage.getItem("InterfaceColor");
    }
    //倒计时秒数
    if (localStorage.getItem("sketchFlipSecond") != null) {
      let saveNum = Number(localStorage.getItem("sketchFlipSecond"));
      if (!isNaN(saveNum)) {
        this.sketchFlipSecond = saveNum;
      }
    }

    // 图片处理相关
    //是否获取黑白图片
    if (localStorage.getItem("ImageParameters_Gray") === "true") {
      this.imageParameters.gray = true;
    } else if (localStorage.getItem("ImageParameters_Gray") === "false") {
      this.imageParameters.gray = false;
    }
    // console.log("读取设置并初始化: ImageParameters_Gray=" + this.imageParameters.gray);

    //是否压缩图片
    if (localStorage.getItem("ImageParameters_DoAutoResize") === "true") {
      this.imageParameters.do_auto_resize = true;
    } else if (
      localStorage.getItem("ImageParameters_DoAutoResize") === "false"
    ) {
      this.imageParameters.do_auto_resize = false;
    }

    //启用压缩的Width下限
    if (localStorage.getItem("ImageParametersResizeMaxWidth") != null) {
      let saveNum = Number(
        localStorage.getItem("ImageParametersResizeMaxWidth")
      );
      if (!isNaN(saveNum)) {
        this.imageParameters.resize_max_width = saveNum;
      }
    }

    //是否自动切白边
    if (localStorage.getItem("ImageParameters_DoAutoCrop") === "true") {
      this.imageParameters.do_auto_crop = true;
    } else if (localStorage.getItem("ImageParameters_DoAutoCrop") === "false") {
      this.imageParameters.do_auto_crop = false;
    }

    //切白边参数
    if (localStorage.getItem("ImageParameters_AutoCropNum") != null) {
      let saveNum = Number(localStorage.getItem("ImageParameters_AutoCropNum"));
      if (!isNaN(saveNum)) {
        this.imageParameters.auto_crop_num = saveNum;
      }
    }
    //是否保存页数
    if (localStorage.getItem("saveNowPageNumFlag") === "true") {
      this.saveNowPageNumFlag = true;
    } else if (localStorage.getItem("saveNowPageNumFlag") === "false") {
      this.saveNowPageNumFlag = false;
    }
    //是否通过websocket同步页数
    if (localStorage.getItem("SyncPageFlag") === "true") {
      this.syncPageFlag = true;
    } else if (localStorage.getItem("SyncPageFlag") === "false") {
      this.syncPageFlag = false;
    }

    //是否自动隐藏工具条
    if (localStorage.getItem("HideToolbar") === "true") {
      this.hideToolbar = true;
    } else if (localStorage.getItem("HideToolbar") === "false") {
      this.hideToolbar = false;
    }
  },
  // beforeMount : 指令第一次绑定到元素并且在挂载父组件之前调用。
  beforeMount() {
    // 自动开始Sketch模式
    if (localStorage.getItem("ReaderMode") === "sketch") {
      this.startSketchMode();
    }
    // 注册监听
    window.addEventListener("keyup", this.handleKeyup);
  },
  //卸载前
  beforeUnmount() {
    // 销毁监听
    window.removeEventListener("keyup", this.handleKeyup);
    //移除websockets消息监听
    delete this.$options.sockets.onmessage;
  },
  // mounted : 在绑定元素的父组件被挂载后调用。
  mounted() {
    //需要得书籍远程数据,避免初始化失败,所以延迟1秒执行
    setTimeout(this.setNowPageNumByLocalStorage, 1000);
  },
  updated() {
    //界面有更新就会调用,随便乱放会引起难以调试的BUG
  },
  methods: {
    //接收服务器发来的websocket消息，做各种反应（翻页、提示信息）
    handlePacket(data) {
      if (!this.syncPageFlag) {
        return;
      }
      // console.log("FlipMode 接收到Message:");
      //data.data也是个字符串，需要解析成对象
      const msg = JSON.parse(data.data);
      //心跳信息,直接返回
      if (msg.type === "heartbeat") {
        return;
      }
      //服务器发来翻页信息，来自于另一个用户，
      if (msg.type === "sync_page" && msg.user_id !== this.$store.userID) {
        const syncData = JSON.parse(msg.data_string);
        //正在读的是同一本书、就翻页。
        if (syncData.book_id === this.book.id && syncData.now_page_num !== this.nowPageNum) {
          // console.log(syncData);
          this.toPage(syncData.now_page_num, false);
        }
      }
    },

    //Websocket 发送消息
    sendNowPage() {
      const readPercent =
        parseFloat(this.nowPageNum) / parseFloat(this.book.all_page_num);
      // console.debug("ReadPercent: " + readPercent)
      const data = {
        book_id: this.book.id,
        now_page_num: this.nowPageNum,
        now_page_num_percent: 1.0,
        read_percent: readPercent,
      };
      // console.log("this.$store.userID: " + this.$store.state.userID)
      const newMsg = {
        type: "sync_page",
        status_code: 200,
        user_id: this.$store.state.userID,
        token: this.$store.state.token,
        detail: "翻页模式，发送数据", // 消息描述
        data_string: JSON.stringify(data),
      };
      // 配置为了json，可调用sendObj方法来发送数据
      this.$socket.sendObj(newMsg);
      // console.log("send:", newMsg);
    },

    // 辅助函数，用于从 Gravatar 获取头像。URL 的最后一段需要用户的 email 地址的 MD5 编码。
    gravatarURL: function (email) {
      return "http://www.gravatar.com/avatar/" + md5(email);
    },

    //页面排序相关
    onResort(key) {
      axios
        .get("/getbook?id=" + this.$route.params.id + "&sort_by=" + key)
        .then((response) => (this.book = response.data))
        .finally(() => {
          document.title = this.book.name;
          this.resort_hint_key = key;
          // 带查询参数，结果是 /#/flip/abc123?sort_by="filesize"
          this.$router.push({
            name: "FlipMode",
            replace: true,
            query: { sort_by: key },
          });
          console.log(
            "成功刷新书籍数据,书籍ID:" +
            this.$route.params.id +
            "  sort_by=" +
            key
          );
        });
    },

    //返回“重新排序”选择菜单的文字提示
    getSortHintText() {
      switch (this.resort_hint_key) {
        case "filename":
          return this.$t("sort_by_filename");
        case "modify_time":
          return this.$t("sort_by_modify_time");
        case "filesize":
          return this.$t("sort_by_filesize");
        case "filename_reverse":
          return this.$t("sort_by_filename") + this.$t("sort_reverse");
        case "modify_time_reverse":
          return this.$t("sort_by_modify_time") + this.$t("sort_reverse");
        case "filesize_reverse":
          return this.$t("sort_by_filesize") + this.$t("sort_reverse");
        default:
          return this.$t("re_sort");
      }
    },

    //图片处理相关
    //黑白化参数
    setImageParameters_Gray(value) {
      // console.log("value:" + value);
      this.imageParameters.gray = value;
      localStorage.setItem("ImageParameters_Gray", value);
      // console.log("成功保存设置: ImageParameters_Gray=" + localStorage.getItem("ImageParameters_Gray"));
    },
    //缩放图片大小的参数
    setImageParameters_DoAutoResize(value) {
      this.imageParameters.do_auto_resize = value;
      localStorage.setItem("ImageParameters_DoAutoResize", value);
      // console.log("成功保存设置: ImageParameters_DoAutoResize=" + localStorage.getItem("ImageParameters_DoAutoResize"));
    },
    //设置是否切白边
    setImageParameters_DoAutoCrop(value) {
      this.imageParameters.do_auto_crop = value;
      localStorage.setItem(
        "ImageParameters_DoAutoCrop",
        this.imageParameters.do_auto_crop ? "true" : "false"
      );
      // console.log("成功保存设置: ImageParameters_DoAutoCrop=" + localStorage.getItem("ImageParameters_DoAutoCrop"));
    },
    //切白边参数
    setImageParameters_AutoCropNum(value) {
      this.imageParameters.auto_crop_num = value;
      localStorage.setItem(
        "ImageParameters_AutoCropNum",
        this.imageParameters.auto_crop_num ? "true" : "false"
      );
    },
    //切换到卷轴模式
    changeReaderModeToScrollMode() {
      localStorage.setItem("ReaderMode", "scroll");
      //replace的作用类似于 router.push，唯一不同的是，它在导航时不会向 history 添加新记录，正如它的名字所暗示的那样——它取代了当前的条目。
      this.$router.replace({
        name: "ScrollMode",
        replace: true,
        params: { id: this.$route.params.id },
      });
    },
    //判断这本书是否需要提供原始压缩包下载
    needDownloadLink() {
      return this.book.book_type !== "dir";
    },
    // 分析单双页用
    nowAndNextPageIsSingle() {
      this.nowPageIsDouble = this.checkImageIsDoublePage_byPageNum(
        this.nowPageNum
      );
      this.nextPageIsDouble = this.checkImageIsDoublePage_byPageNum(
        this.nowPageNum + 1
      );
      if (this.nowPageIsDouble || this.nextPageIsDouble) {
        this.nowAndNextPageIsSingleFlag = false;
        return false;
      } else {
        this.nowAndNextPageIsSingleFlag = false;
        return true;
      }
    },
    //根据书籍UUID,设定当前页数,因为需要取得远程书籍数据（this.book）,所以延迟执行
    setNowPageNumByLocalStorage() {
      if (this.saveNowPageNumFlag) {
        let cookieValue = localStorage.getItem("nowPageNum" + this.book.id);
        if (cookieValue != null) {
          let saveNum = Number(cookieValue);
          if (!isNaN(saveNum)) {
            this.nowPageNum = saveNum;
            console.log("成功读取页数" + saveNum);
          } else {
            console.log("读取页数失败,this.nowPageNum = " + this.nowPageNum);
          }
        } else {
          console.log("本队存储里没找到:" + "nowPageNum = " + this.nowPageNum);
        }
      } else {
        //console.log("不读取页数,this.saveNowPageNumFlag=" + this.saveNowPageNumFlag);
      }
    },
    // // 设置当前模板-接收Drawer的参数,继续往父组件传
    // 改变阅读模式
    OnSetTemplate(value) {
      if (value === "scroll") {
        console.log("跳转到卷轴阅读模式");
      } else if (this.readerMode === "scroll" || this.readerMode === "sketch") {
        // 命名路由,并加上参数,让路由建立 url
        this.$router.push({ name: "ScrollMode", params: { id: this.book.id } });
      }
    },
    //打开抽屉
    drawerActivate(place) {
      this.drawerActive = true;
      this.drawerPlacement = place;
      //打开抽屉时，触发工具条显隐
      if (this.hideToolbar) {
        this.showHeaderFlag_FlipMode = true
        this.showFooterFlag_FlipMode = true
      }
    },
    //关闭抽屉
    drawerDeactivate() {
      this.drawerActive = false;
      //关闭抽屉时，触发工具条显隐
      if (this.hideToolbar) {
        this.showHeaderFlag_FlipMode = false
        this.showFooterFlag_FlipMode = false
      }
    },
    //开始速写倒计时
    startSketchMode() {
      this.readerMode = "sketch";
      this.message.success(this.$t("startSketchMessage"));
      this.drawerActive = false; //关闭设置抽屉
      this.sketchModeFlag = true;
      //是否倒计时提示文字
      this.showPageHintFlag_FlipMode = true;
      //是否显示页头
      this.showHeaderFlag_FlipMode = false;
      //是否显示页脚
      this.showFooterFlag_FlipMode = false;
      // this.$emit("setTemplate", "sketch");
      //setTimeout和setInterval函数,都返回一个表示计数器编号的整数值,将该整数传入clearTimeout和clearInterval函数,就可以取消对应的定时器。setInterval指定某个任务每隔一段时间就执行一次。setTimeout()用于在指定的毫秒数后调用函数或计算表达式  setTimeout('console.log(2)',1000);
      this.interval = setInterval(this.sketchCount, 1000);
    },
    //修改间隔的时候重新计秒数
    resetSketchSecondCount() {
      this.sketchSecondCount = 0;
    },
    //停止速写倒计时
    stopSketchMode() {
      this.message.success(this.$t("goodjob_and_byebye"));
      this.sketchModeFlag = false;
      this.showPageHintFlag_FlipMode = false;
      this.sketchSecondCount = 0;
      //是否显示页头
      this.showHeaderFlag_FlipMode = true;
      //是否显示页脚
      this.showFooterFlag_FlipMode = true;
      this.readerMode = "flip";
      // this.$emit("setTemplate", "flip");
      clearInterval(this.interval); // 清除定时器
    },
    //开始速写（quick sketch）,每秒执行一次
    sketchCount() {
      if (this.sketchModeFlag === false) {
        this.stopSketchMode();
      }
      this.sketchSecondCount = this.sketchSecondCount + 1;
      let nowSecond = this.sketchSecondCount % this.sketchFlipSecond;
      // console.log("sketchSecondCount=" + this.sketchSecondCount + " nowSecond:" + nowSecond)
      if (nowSecond === 0) {
        if (this.nowPageNum < this.book.all_page_num) {
          this.flipPage(1);
        } else {
          this.toPage(1);
        }
      }
    },
    // 关闭抽屉时,保存设置到cookies
    saveConfigToLocal() {
      localStorage.setItem("SyncPageFlag", this.syncPageFlag);
      localStorage.setItem("debugModeFlag", this.debugModeFlag);
      localStorage.setItem(
        "showHeaderFlag_FlipMode",
        this.showHeaderFlag_FlipMode
      );
      localStorage.setItem(
        "showFooterFlag_FlipMode",
        this.showFooterFlag_FlipMode
      );
      localStorage.setItem(
        "showPageHintFlag_FlipMode",
        this.showPageHintFlag_FlipMode
      );
      localStorage.setItem("rightToLeftFlag", this.rightToLeftFlag);
      localStorage.setItem(
        "simpleDoublePageModeFlag",
        this.simpleDoublePageModeFlag
      );
      localStorage.setItem(
        "autoDoublePageModeFlag",
        this.autoDoublePageModeFlag
      );
      localStorage.setItem("saveNowPageNumFlag", this.saveNowPageNumFlag);
      localStorage.setItem("nowPageNum" + this.book.id, this.nowPageNum);
      localStorage.setItem("BackgroundColor", this.model.backgroundColor);
      localStorage.setItem("sketchFlipSecond", this.sketchFlipSecond);
      //set对有setXXXChange函数的来说有些多余,但没有set函数的话就有必要了
      localStorage.setItem(
        "ImageParameters_DoAutoCrop",
        this.imageParameters.do_auto_crop ? "true" : "false"
      );
      localStorage.setItem(
        "ImageParametersResizeMaxWidth",
        this.imageParameters.resize_max_width ? "true" : "false"
      );
    },
    // 随机换一下背景色
    randomBackgroundColor() {
      let R = Math.ceil(Math.random() * 155) + 100;
      let G = Math.ceil(Math.random() * 155) + 100;
      let B = Math.ceil(Math.random() * 100) + 100;
      //rgb(185,175,145)
      let RGB = "rgb(" + R + "," + G + "," + B + ")";
      // console.log(RGB);
      this.model.backgroundColor = RGB;
    },
    //HTML DOM 事件 https://www.runoob.com/jsref/dom-obj-event.html
    // 进入绑定该事件的元素和其子元素均会触发该事件,所以有一个重复触发,冒泡过程。其对应的离开事件 mouseout
    onMouseOver() {
    },
    // 只有进入绑定该事件的元素才会触发事件,也就是不会冒泡。其对应的离开事件mouseleave
    onMouseEnter() {
      // this.randomColor = 'background-color: rgb(235,235,235)';
    },
    onMouseLeave(e) {
      //离开区域的时候,清空鼠标样式
      e.currentTarget.style.cursor = "";
    },
    //事件修饰符: https://v3.cn.vuejs.org/guide/events.html#%E4%BA%8B%E4%BB%B6%E4%BF%AE%E9%A5%B0%E7%AC%A6
    onMouseMove(e) {
      // // offsetX/Y获取到是触发点相对被触发dom的左上角距离
      // let offsetX = e.offsetX;
      // let offsetY = e.offsetY;
      //根据ID获取元素
      // let mangaDiv =document.getElementById("MangaMain")
      //不用自己获取元素
      // let offsetWidth = e.currentTarget.offsetWidth;
      // let offsetHeight = e.currentTarget.offsetHeight;
      let clickX = e.x; //获取鼠标的X坐标（鼠标与屏幕左侧的距离,单位为px）
      let clickY = e.y; //获取鼠标的Y坐标（鼠标与屏幕顶部的距离,单位为px）
      // console.log("clickX: " + clickX);
      // console.log("clickY: " + clickY);
      // 浏览器的视口,不包括工具栏和滚动条:
      let innerWidth = window.innerWidth;
      let innerHeight = window.innerHeight;

      //是否进入上下工具条附近区域时，触发工具条显隐
      const toolBarArea = 0.10;
      let inToolBarArea = false
      if (clickY <= (innerHeight * toolBarArea) || clickY >= (innerHeight * (1.0 - toolBarArea))) {
        inToolBarArea = true
      }

      //进入设置区域的时候，设置鼠标的形状
      //设置区域为正方形，边长按照宽或高里面，比较小的值决定
      const setArea = 0.15;
      // innerWidth >= innerHeight 的情况下
      let MinY = innerHeight * (0.5 - setArea);
      let MaxY = innerHeight * (0.5 + setArea);
      let MinX = innerWidth * 0.5 - (MaxY - MinY) * 0.5;
      let MaxX = innerWidth * 0.5 + (MaxY - MinY) * 0.5;
      if (innerWidth < innerHeight) {
        MinX = innerWidth * (0.5 - setArea);
        MaxX = innerWidth * (0.5 + setArea);
        MinY = innerHeight * 0.5 - (MaxX - MinX) * 0.5;
        MaxY = innerHeight * 0.5 + (MaxX - MinX) * 0.5;
      }
      //在设置区域
      let inSetArea = false
      if (clickX > MinX && clickX < MaxX && clickY > MinY && clickY < MaxY) {
        inSetArea = true
      }
      if (inSetArea) {
        e.currentTarget.style.cursor =
          "url(/images/SettingsOutline.png), pointer";
      } else {
        if (clickX < innerWidth * 0.5) {
          //设置左边的鼠标指针
          if (this.rightToLeftFlag && this.nowPageNum === 1) {
            //右边翻下一页,且目前是第一页的时候,左边的鼠标指针,设置为禁止翻页
            e.currentTarget.style.cursor =
              "url(/images/Prohibited28Filled.png), pointer";
          } else if (
            !this.rightToLeftFlag &&
            this.nowPageNum === this.book.all_page_num
          ) {
            //左边翻下一页,且目前是最后一页的时候,左边的鼠标指针,设置为禁止翻页
            e.currentTarget.style.cursor =
              "url(/images/Prohibited28Filled.png), pointer";
          } else {
            //正常情况下,左边是向左的箭头
            e.currentTarget.style.cursor =
              "url(/images/ArrowLeft.png), pointer";
          }
        } else {
          //设置右边的鼠标指针
          if (
            this.rightToLeftFlag &&
            this.nowPageNum === this.book.all_page_num
          ) {
            //右边翻下一页,且目前是最后页的时候,右边的鼠标指针,设置为禁止翻页
            e.currentTarget.style.cursor =
              "url(/images/Prohibited28Filled.png), pointer";
          } else if (!this.rightToLeftFlag && this.nowPageNum === 1) {
            //左边翻下一页,且目前是第一页的时候,右边的鼠标指针,设置为禁止翻页
            e.currentTarget.style.cursor =
              "url(/images/Prohibited28Filled.png), pointer";
          } else {
            //正常情况下,右边是向右的箭头
            e.currentTarget.style.cursor =
              "url(/images/ArrowRight.png), pointer";
          }
        }
      }
      //进入上下工具条附近区域时，触发工具条显隐
      if ((inToolBarArea || inSetArea) && this.hideToolbar) {
        this.showHeaderFlag_FlipMode = true
        this.showFooterFlag_FlipMode = true
      } else if (this.hideToolbar) {
        this.showHeaderFlag_FlipMode = false
        this.showFooterFlag_FlipMode = false
      }
    },

    //根据鼠标点击事件的位置,决定是左右翻页还是打开设置
    onMouseClick(e) {
      let clickX = e.x; //获取鼠标的X坐标（鼠标与屏幕左侧的距离,单位为px）
      let clickY = e.y; //获取鼠标的Y坐标（鼠标与屏幕顶部的距离,单位为px）
      //浏览器的可视区域宽高,不包括工具栏和滚动条:
      let innerHeight = window.innerHeight;
      let innerWidth = window.innerWidth;
      // 设置区域为正方形
      const setArea = 0.15;
      // innerWidth >= innerHeight 的情况下
      let MinY = innerHeight * (0.5 - setArea);
      let MaxY = innerHeight * (0.5 + setArea);
      let MinX = innerWidth * 0.5 - (MaxY - MinY) * 0.5;
      let MaxX = innerWidth * 0.5 + (MaxY - MinY) * 0.5;
      if (innerWidth < innerHeight) {
        MinX = innerWidth * (0.5 - setArea);
        MaxX = innerWidth * (0.5 + setArea);
        MinY = innerHeight * 0.5 - (MaxX - MinX) * 0.5;
        MaxY = innerHeight * 0.5 + (MaxX - MinX) * 0.5;
      }
      // console.log("鼠标点击：e.offsetX=" + offsetX, "e.offsetY=" + offsetY);
      if (clickX > MinX && clickX < MaxX && clickY > MinY && clickY < MaxY) {
        //点中了设置区域
        this.drawerActivate("right");
      } else {
        //决定如何翻页
        if (clickX < innerWidth * 0.5) {
          //左边的翻页
          if (this.rightToLeftFlag === true) {
            this.toPerviousPage();
          } else {
            this.toNextPage();
          }
        } else {
          //右边的翻页
          if (this.rightToLeftFlag === true) {
            this.toNextPage();
          } else {
            this.toPerviousPage();
          }
        }
      }
    },
    toNextPage() {
      //简单合并模式
      if (this.simpleDoublePageModeFlag) {
        if (this.nowPageNum < this.book.all_page_num - 1) {
          this.flipPage(2);
          return;
        } else {
          this.flipPage(1);
          return;
        }
      }

      //如果开启了自动合并模式,并且当前页应该被合并
      if (
        this.autoDoublePageModeFlag &&
        this.checkMergedStatus_ByPageNum(this.nowPageNum)
      ) {
        if (this.nowPageNum < this.book.all_page_num - 1) {
          this.flipPage(2);
        } else {
          this.flipPage(1);
        }
      } else {
        this.flipPage(1);
      }
    },
    toPerviousPage() {
      //错误值,第0或第1页。
      if (this.nowPageNum <= 1) {
        // console.log("Error toPerviousPage");
        this.message.info(this.$t("hintFirstPage"));
        return;
      }

      //简单合并模式
      if (this.simpleDoublePageModeFlag) {
        if (this.nowPageNum - 2 > 0) {
          this.flipPage(-2);
          return;
        } else {
          this.flipPage(-1);
          return;
        }
      }

      //自动合并模式
      //如果没有开启自动合并模式,或现在是第2页
      if (this.nowPageNum === 2 || !this.autoDoublePageModeFlag) {
        this.flipPage(-1);
        return;
      }
      //如果前一页是双开叶
      let pervIsDouble = this.checkImageIsDoublePage_byPageNum(
        this.nowPageNum - 1
      );
      if (pervIsDouble) {
        this.flipPage(-1);
        return;
      }
      //如果前前页是双开叶
      let PervPervIsDouble = this.checkImageIsDoublePage_byPageNum(
        this.nowPageNum - 2
      );
      if (PervPervIsDouble) {
        this.flipPage(-1);
        return;
      }
      //都没return掉,那么前两张可以合并,减两页
      this.flipPage(-2);
    },

    //给一个页数,然后判断自动双页模式下,是否应该预读并合并显示下一页
    checkMergedStatus_ByPageNum(pageNum) {
      //如果没有开启自动双页模式,当然不需要
      if (!this.autoDoublePageModeFlag) {
        return false;
      }
      //可能传入的错误值,打印到控制台
      if (pageNum <= 0 || pageNum >= this.book.all_page_num) {
        console.log("Error pageNum :" + pageNum);
        return false;
      }

      //已经是最后一页了,显然不需要自动合并下一页
      if (pageNum === this.book.all_page_num) {
        return false;
      }
      //分析现在这张图片的宽高比,看是不是双开页
      let now_page_is_double_page =
        this.checkImageIsDoublePage_byPageNum(pageNum);
      //分析下一张漫画的宽高比,看是不是双开页
      let next_page_is_double_page = this.checkImageIsDoublePage_byPageNum(
        pageNum + 1
      );
      //如果现在这张就是开页漫画,直接不用比
      //如果下一张漫画是开页,显然也没法合并
      return !(now_page_is_double_page || next_page_is_double_page);
    },
    checkImageIsDoublePage_byPageNum(pageNum) {
      //如果传进了错误值
      if (pageNum <= 0 || pageNum > this.book.all_page_num) {
        console.log("Error checkImageIsDoublePage_byPageNum:" + pageNum);
        return;
      }
      if (this.book.pages.images[pageNum - 1].image_type === "SinglePage") {
        return false;
      }
      if (this.book.pages.images[pageNum - 1].image_type === "DoublePage") {
        return true;
      }
      let image = new Image();
      let temp_flag = false; //返回结果用
      image.src = this.book.pages.images[pageNum - 1].url;
      // image.complete 图片是否完全加载完成。
      //https://developer.mozilla.org/zh-CN/docs/Web/API/HTMLImageElement/complete

      // https://corbusier.github.io/2017/06/29/%E5%9B%BE%E7%89%87%E7%9A%84%E5%BC%82%E6%AD%A5%E5%8A%A0%E8%BD%BD%E4%B8%8Eonload%E5%87%BD%E6%95%B0/
      // onload函数要写在改变src前,这样确保了onload函数一定会被调用
      // complete只是变向的在判断img是否已经触发了load事件,而且是不精准的判断
      // complete在不同浏览器下,表现不一致,不建议使用
      // 无论浏览器是否存在图片缓存,重新请求图片地址,都会触发onload事件

      // 牵扯到加载资源,这段代码需要改进……
      if (image.complete) {
        return image.width >= image.height;
      } else {
        //否则加载图片
        image.onload = function () {
          //是单页漫画
          if (image.width < image.height) {
            image.onload = null; // 避免重复加载
            temp_flag = false; //宽小于高,是是竖着的,单页漫画
          } else {
            //是双页漫画
            image.onload = null; // 避免重复加载
            temp_flag = true; //宽大于高,很可能是开页
          }
        };
        return temp_flag;
      }
    },

    //翻页,其实不限页数
    flipPage: function (num) {
      if (
        this.nowPageNum + num <= this.book.all_page_num &&
        this.nowPageNum + num >= 1
      ) {
        this.nowPageNum = this.nowPageNum + num;
      } else {
        // console.log("无法继续翻,Num:" + num)
        if (num > 0) {
          this.message.info(this.$t("hintLastPage"));
        } else {
          this.message.info(this.$t("hintFirstPage"));
        }
      }
      //保存页数
      this.saveNowPageNumOnUpdate(this.nowPageNum);
    },
    //拖动进度条,或翻页的时候保存页数
    saveNowPageNumOnUpdate(value, sendWSMessage = true) {
      if (this.saveNowPageNumFlag) {
        localStorage.setItem("nowPageNum" + this.book.id, value);
      }
      //发送翻页消息到服务器
      if (sendWSMessage && this.syncPageFlag) {
        this.sendNowPage();
      }
    },
    //跳转到指定页数
    toPage: function (num, sendWSMessage = true) {
      if (num <= this.book.all_page_num && num >= 1) {
        this.nowPageNum = num;
      }
      //保存页数
      this.saveNowPageNumOnUpdate(this.nowPageNum, sendWSMessage);
    },

    // 键盘事件
    handleKeyup(event) {
      //错误:(815, 49) 不允许从实参中引用 'caller' 和 'callee'
      const e = event || window.event;
      if (!e) return;
      //https://developer.mozilla.org/zh-CN/docs/Web/API/KeyboardEvent/keyCode
      switch (e.key) {
        case "ArrowUp":
        case "PageUp":
          this.flipPage(-1); //上一页
          break;
        case "ArrowLeft":
          this.rightToLeftFlag === true
            ? this.toPerviousPage()
            : this.toNextPage();
          break;
        case "ArrowRight":
          this.rightToLeftFlag === true
            ? this.toNextPage()
            : this.toPerviousPage();
          break;
        case "Space":
        case "ArrowDown":
        case "PageDown":
          this.flipPage(1); //下一页
          break;
        case "Home":
          this.toPage(1); //跳转到第一页
          break;
        case "End":
          this.toPage(this.book.all_page_num); //跳转到最后一页
          break;
        case "Ctrl":
          // Ctrl key pressed //组合键？
          //openOverlay();
          break;
      }
      // console.log(e.keyCode);
      // console.log(e.key);
    },

    setShowHeaderChange(value) {
      console.log("value:" + value);
      this.showHeaderFlag_FlipMode = value;
      localStorage.setItem("showHeaderFlag_FlipMode", value);
      console.log(
        "cookie设置完毕: showHeaderFlag_FlipMode=" +
        localStorage.getItem("showHeaderFlag_FlipMode")
      );
    },
    setShowFooterFlagChange(value) {
      console.log("value:" + value);
      this.showFooterFlag_FlipMode = value;
      localStorage.setItem("showFooterFlag_FlipMode", value);
      console.log(
        "cookie设置完毕: showFooterFlag_FlipMode=" +
        localStorage.getItem("showFooterFlag_FlipMode")
      );
    },

    setShowPageNumChange(value) {
      console.log("value:" + value);
      this.showPageHintFlag_FlipMode = value;
      localStorage.setItem("showPageHintFlag_FlipMode", value);
      console.log(
        "cookie设置完毕: showPageHintFlag_FlipMode=" +
        localStorage.getItem("showPageHintFlag_FlipMode")
      );
    },

    setFlipScreenFlag(value) {
      console.log("value:" + value);
      this.rightToLeftFlag = value;
      localStorage.setItem("rightToLeftFlag", value);
      console.log(
        "cookie设置完毕: rightToLeftFlag=" +
        localStorage.getItem("rightToLeftFlag")
      );
    },

    setHideToolbar(value) {
      console.log("value:" + value);
      this.hideToolbar = value;
      localStorage.setItem("HideToolbar", value);
      console.log(
        "cookie设置完毕: HideToolbar=" +
        localStorage.getItem("HideToolbar")
      );
    },

    setSyncPageFlag(value) {
      console.log("value:" + value);
      this.syncPageFlag = value;
      localStorage.setItem("SyncPageFlag", value);
      console.log(
        "cookie设置完毕: SyncPageFlag=" +
        localStorage.getItem("SyncPageFlag")
      );
    },

    setSavePageNumFlag(value) {
      console.log("value:" + value);
      this.saveNowPageNumFlag = value;
      localStorage.setItem("saveNowPageNumFlag", value);
      console.log(
        "cookie设置完毕: saveNowPageNumFlag=" +
        localStorage.getItem("saveNowPageNumFlag")
      );
    },

    setDebugModeFlag(value) {
      console.log("value:" + value);
      this.debugModeFlag = value;
      //关闭Debug模式的时候顺便也关上“自动合并单双页”的功能（因为还有BUG）
      if (value === false) {
        this.autoDoublePageModeFlag = false;
      }
      localStorage.setItem("debugModeFlag", value);
      console.log(
        "cookie设置完毕: debugModeFlag=" + localStorage.getItem("debugModeFlag")
      );
    },

    setAutoDoublePage_FlipMode(value) {
      console.log("value:" + value);
      this.autoDoublePageModeFlag = value;
      if (value === true) {
        this.simpleDoublePageModeFlag = false;
      }
      localStorage.setItem("autoDoublePageModeFlag", value);
      console.log(
        "cookie设置完毕: autoDoublePageModeFlag=" +
        localStorage.getItem("autoDoublePageModeFlag")
      );
    },

    setSimpleDoublePage_FlipMode(value) {
      console.log("value:" + value);
      this.simpleDoublePageModeFlag = value;
      if (value === true) {
        this.autoDoublePageModeFlag = false;
      }
      localStorage.setItem("simpleDoublePageModeFlag", value);
      console.log(
        "cookie设置完毕: simpleDoublePageModeFlag=" +
        localStorage.getItem("simpleDoublePageModeFlag")
      );
    },
  },

  computed: {
    //页数或素描模式的提示
    pageNumOrSketchHint() {
      if (this.sketchModeFlag) {
        let nowSecond = (this.sketchSecondCount % this.sketchFlipSecond) + 1;
        let donePage = parseInt(this.sketchSecondCount / this.sketchFlipSecond);
        let totalMinutes = (this.sketchSecondCount + 1) / 60;
        //计算几小时几分
        let MinutesAndHourString = "";
        //如果不满意1小时,就不显示小时
        if (totalMinutes / 60 === 0) {
          MinutesAndHourString = totalMinutes + this.$t("minute");
        } else {
          MinutesAndHourString =
            parseInt(totalMinutes / 60) +
            this.$t("hour") +
            parseInt(totalMinutes % 60) +
            this.$t("minute");
        }
        let AllTimeString =
          MinutesAndHourString +
          ((this.sketchSecondCount + 1) % 60) +
          this.$t("second");
        return (
          this.$t("now_is") +
          nowSecond +
          this.$t("second") +
          "  " +
          this.$t("totalTime") +
          AllTimeString +
          "   " +
          this.$t("interval") +
          this.sketchFlipSecond +
          this.$t("second") +
          "  " +
          this.$t("total_is") +
          donePage +
          this.$t("page")
        );
      } else {
        return this.nowPageNum + "/" + this.book.all_page_num;
      }
    },
    mangaAreaHeight() {
      let Height = 95;
      //页头和底部拖动条都显示,或有一个显示的时候,95%
      if (this.showFooterFlag_FlipMode && this.showHeaderFlag_FlipMode) {
        Height = 95;
      }
      if (this.showFooterFlag_FlipMode && !this.showHeaderFlag_FlipMode) {
        Height = 95;
      }
      if (!this.showFooterFlag_FlipMode && this.showHeaderFlag_FlipMode) {
        Height = 95;
      }
      //页头和底部拖动条都不显示的时候,漫画占满屏幕
      if (!this.showFooterFlag_FlipMode && !this.showHeaderFlag_FlipMode) {
        Height = 100;
      }
      return Height + "vh";
    },
    mangaImageHeight() {
      let Height = 95;
      //页头和底部拖动条都显示,或有一个显示的时候,95%
      if (this.showFooterFlag_FlipMode && this.showHeaderFlag_FlipMode) {
        Height = 95;
      }
      if (this.showFooterFlag_FlipMode && !this.showHeaderFlag_FlipMode) {
        Height = 95;
      }
      if (!this.showFooterFlag_FlipMode && this.showHeaderFlag_FlipMode) {
        Height = 95;
      }
      //页头和拖动条都不显示的时候,漫画占满屏幕
      if (!this.showFooterFlag_FlipMode && !this.showHeaderFlag_FlipMode) {
        Height = 100;
      }
      //与上面唯一的不同,减去素描提示的空间
      if (this.showPageHintFlag_FlipMode) {
        if (this.readerMode === "sketch") {
          Height = Height - 6;
        } else {
          Height = Height - 3;
        }
      }
      return Height + "vh";
    },
    //进入素描模式的时候,把高度放大一倍
    sketchHintHeight() {
      if (this.readerMode === "sketch") {
        return "6vh";
      } else {
        return "3vh";
      }
    },
    //进入素描模式的时候,把字体放大
    sketchHintFontSize() {
      if (this.readerMode === "sketch") {
        return "24px";
      } else {
        return "16px";
      }
    },
    //从左到右还是从右到左
    get_flex_direction() {
      if (this.rightToLeftFlag === true) {
        return "row";
      } else {
        return "row-reverse";
      }
    },
  },
});
</script>

<style scoped>
/* 题头的显隐效果 */
/* https://v3.cn.vuejs.org/guide/transitions-enterleave.html#%E5%8D%95%E5%85%83%E7%B4%A0-%E7%BB%84%E4%BB%B6%E7%9A%84%E8%BF%87%E6%B8%A1= */
/* 可以为进入和离开动画设置不同的持续时间和动画函数 */
/* v-enter-active：定义进入过渡生效时的状态。在整个进入过渡的阶段中应用，在元素被插入之前生效，在过渡/动画完成之后移除。这个类可以被用来定义进入过渡的过程时间，延迟和曲线函数。 */
.header-bottom-enter-active {
  transition: all 0.3s ease;
}

/* v-leave-active：定义离开过渡生效时的状态。在整个离开过渡的阶段中应用，在离开过渡被触发时立刻生效，在过渡/动画完成之后移除。这个类可以被用来定义离开过渡的过程时间，延迟和曲线函数。 */
/* cubic-bezier() 函数定义了一个贝塞尔曲线(Cubic Bezier) */
.header-bottom-leave-active {
  transition: all 0.8s cubic-bezier(1, 0.5, 0.8, 1);
}

/* v-enter-from：定义进入过渡的开始状态。在元素被插入之前生效，在元素被插入之后的下一帧移除。 */
/* v-leave-to：离开过渡的结束状态。在离开过渡被触发之后下一帧生效 (与此同时 v-leave-from 被移除)，在过渡/动画完成之后移除。 */
.header-bottom-enter-from,
.header-bottom-leave-to {
  /* transform: translateY(-20px); */
  opacity: 0;
}



/* 参考CSS盒子模型改 */
/* https://www.runoob.com/css/css-boxmodel.html */
/* CSS 高度和宽度 */
/* https://www.w3school.com.cn/css/css_dimension.asp */
/* CSS Flexbox 布局 */
/* https://www.w3school.com.cn/css/css3_flexbox.asp */

* {
  box-sizing: border-box;
}

.main {
  font-family: Arial;
  margin: 0;
  padding: 0;
  /* height: v-bind(mangaAreaHeight); */
  display: flex;
  max-width: 100%;
  /* flex-direction: column垂直堆叠 flex 项目（从上到下）,column-reverse从下到上 row从左到右 row-reverse从右到左 */
  flex-direction: column;
  /* justify-content 属性用于对齐 flex 项目。 将 justify-content 和 align-items 属性设置为居中,flex 项目将完美居中： */
  justify-content: center;
  /* center 值将 flex 项目在容器中间对齐： */
  align-items: center;
  background: v-bind("model.backgroundColor");
}

/* 漫画div */
.manga_area {
  width: 100vw;
  max-height: 100vh;
  height: 100vh;
  max-width: 100vw;
  padding: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: baseline;
  /* 不可以被选中 */
  user-select: none;
  /* 火狐 */
  -moz-user-select: none;
  /* 谷歌 */
  -webkit-user-select: none;
  /* border-radius: 3px; */
  /* box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19); */
}

/* 漫画div中的图片div*/
.manga_area_img_div {
  width: 100vw;
  /* height: v-bind(mangaAreaHeight); */
  height: 100vh;
  display: flex;
  flex-direction: v-bind(get_flex_direction);
  justify-content: center;
  align-items: center;
  margin: 0px;
}

/* 最后的一或两张图片*/
.manga_area_img_div img {
  max-height: 100vh;
  max-width: 100vw;
  min-width: 0;
  width: auto;
  display: block;
  user-select: none;
  /* 下面两句，是设定高度为100%，同时保持比例缩放的关键（类似stetch small pages） */
  min-height: 100vh; 
  object-fit: contain;
  /* 两张图片之间不要留空间*/
  padding: 0px;
  /* background-color: #aaa; */
}
</style>
