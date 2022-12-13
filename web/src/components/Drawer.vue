<template>
  <!-- 外边距: m-2 https://www.tailwindcss.cn/docs/margin -->
  <!-- 内边距： p-4 https://www.tailwindcss.cn/docs/padding  p-0 m-0  -->
  <n-drawer v-bind:show="drawerActive" @update:show="saveConfigToCookie" :height="275" :width="251"
    :placement="drawerPlacement">
    <n-drawer-content closable>
      <!-- 抽屉：自定义头部 -->
      <template #header>
        <span>{{ $t('ReaderSettings') }}</span>
      </template>

      <!-- 父组件在此处插入自定义内容 -->
      <n-space vertical>
        <slot></slot>
        <n-button v-if="sketching === false && inBookShelf === false" @click="startSketchMode">{{
        $t('startSketchMode')
        }}</n-button>
        <n-button v-if="sketching === true && inBookShelf === false" @click="stopSketchMode">{{
        $t('stopSketchMode')
        }}</n-button>
        <!-- <n-divider /> -->
        <span>{{ $t('scan_qrcode') }}</span>
        <Qrcode></Qrcode>
      </n-space>
      <n-popconfirm @positive-click="handlePositiveClick" @negative-click="handleNegativeClick">
        <template #trigger>
          <n-button>{{ $t('reset_all_settings') }}</n-button>
        </template>
        {{ $t('do_you_reset_all_settings') }}
      </n-popconfirm>

      <!-- 抽屉：自定义底部 -->
      <template #footer>
        <n-button @click="onFullSreen">{{ $t('fullscreen') }}</n-button>
        <n-select placeholder="{{ $t('select-language') }}" v-model:value="$i18n.locale" :options="languageOptions"
          @update:value="OnChangeLanguage" />
      </template>
    </n-drawer-content>
  </n-drawer>
</template>
<script lang="ts">
import screenfull from 'screenfull'
import { useCookies } from "vue3-cookies";
import { NDrawer, NDrawerContent, NButton, NSelect, NPopconfirm, useMessage, NSpace, } from 'naive-ui'
import { defineComponent, } from 'vue'
// import { useI18n } from 'vue-i18n'
import Qrcode from "@/components/Qrcode.vue";
export default defineComponent({
  name: "SettingsDrawer",
  props: ['book', 'initDrawerActive', 'initDrawerPlacement', 'readerMode', 'inBookShelf', "sketching"],
  emits: ["setRM", "saveConfig", "startSketch", "stopSketch", "closeDrawer"],//用于向父组件传递信息，父组件的语法为@setRM="OnSetReaderMode"
  components: {
    NDrawer,//抽屉，可以从上下左右4个方向冒出. https://www.naiveui.com/zh-CN/os-theme/components/drawer
    NDrawerContent,//抽屉内容
    NSpace,//间距 https://www.naiveui.com/zh-CN/os-theme/components/space
    // NDivider,//间隔
    // NRadioGroup,//单选  https://www.naiveui.com/zh-CN/os-theme/components/radio
    // NRadioButton,//单选 用按钮显得更优雅一点
    NButton,//按钮，来自:https://www.naiveui.com/zh-CN/os-theme/components/button
    // NAvatar,//头像 https://www.naiveui.com/zh-CN/os-theme/components/avatar
    NSelect, //选择器 https://www.naiveui.com/zh-CN/os-theme/components/select
    NPopconfirm, //弹出确认 https://www.naiveui.com/zh-CN/os-theme/components/popconfirm
    Qrcode,//https://github.com/scopewu/qrcode.vue
  },
  setup() {
    const { cookies } = useCookies();
    const message = useMessage(); //需要导入 'naive-ui'的useMessage
    return {
      message,
      handlePositiveClick() {
        // message.info("是的");
        //清除localStorage保存的设定
        localStorage.clear();
        //刷新当前页面
        location.reload();
      },
      handleNegativeClick() {
        // message.info("并不");
      },

      cookies,
      languageOptions: [
        {
          label: "English",
          value: 'en',
        },
        {
          label: '日本語',
          value: 'ja'
        },
        {
          label: '中文',
          value: 'zh'
        },
      ],
    };
  },
  data() {
    return {
      someflag: "",
      isFullscreen: false,
      readModeLocal: "",
    };
  },
  //挂载前
  beforeMount() {
    var lang = this.cookies.get("userLanguageSetting");
    if (lang) {
      this.$i18n.locale = lang;
    }
    this.readModeLocal = this.readerMode;
  },
  computed: {
    drawerActive() {
      return this.initDrawerActive;
    },
    drawerPlacement() {
      return this.initDrawerPlacement;
    },
  },
  methods: {
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
    OnChangeLanguage(value: string) {
      this.cookies.set("userLanguageSetting", value);
    },
    // 关闭抽屉时，保存设置到cookies
    saveConfigToCookie(show: boolean) {
      if (show == false) {
        this.$emit('closeDrawer');
        this.$emit('saveConfig');
      }
    },
    startSketchMode() {
      this.$emit('startSketch');
    },
    stopSketchMode() {
      this.$emit('stopSketch');
    },
    //切换模板的函数，需要配合vue-router
    onChangeTemplate() {
      if (this.readModeLocal === "scroll") {
        this.$emit("setRM", "scroll");
      }
      if (this.readModeLocal === "flip") {
        this.$emit("setRM", "flip");
      }
      if (this.readModeLocal === "sketch") {
        this.$emit("setRM", "sketch");
      }
      // location.reload(); //需要刷新？ 以后研究VueRouter并去掉
    },
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>




