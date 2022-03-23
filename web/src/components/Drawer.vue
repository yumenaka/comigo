<template>
  <n-drawer
    v-bind:show="drawerActive"
    @update:show="saveConfigToCookie"
    :height="275"
    :width="251"
    :placement="drawerPlacement"
  >
    <n-drawer-content closable>
      <!-- 抽屉：自定义头部 -->
      <template #header>
        <span>{{ $t('ReaderSettings') }}</span>
      </template>
      <!-- 选择：切换页面模式 -->
      <n-space>
        <n-button @click="toFlipMode">切换为翻页阅读</n-button>
        <n-button @click="toScrollMode">切换为滚动阅读</n-button>
      </n-space>
      <!-- 分割线 -->
      <n-divider />
      <!-- 父组件在此处插入自定义内容 -->
      <slot></slot>
      <n-divider />
      <n-popconfirm @positive-click="handlePositiveClick" @negative-click="handleNegativeClick">
        <template #trigger>
          <n-button>{{ $t('reset_all_settings') }}</n-button>
        </template>
        {{ $t('do_you_reset_all_settings') }}
      </n-popconfirm>
      <!-- 抽屉：自定义底部 -->
      <template #footer>
        <n-select
          placeholder="select language"
          v-model:value="this.$i18n.locale"
          :options="this.languageOptions"
          @update:value="OnChangeLanguage"
        />
        <n-button
          v-if="this.sketching == false"
          @click="startSketchMode"
        >{{ $t('startSketchMode') }}</n-button>
        <n-button v-if="this.sketching == true" @click="stopSketchMode">{{ $t('stopSketchMode') }}</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>
<script>

import { useCookies } from "vue3-cookies";
import { NDrawer, NDivider, NDrawerContent, NSpace, NButton, NSelect, NPopconfirm, } from 'naive-ui'
import { defineComponent, } from 'vue'
// import { useI18n } from 'vue-i18n'

export default defineComponent({
  name: "Drawer",
  props: ['book', 'initDrawerActive', 'initDrawerPlacement', 'ReaderMode', "sketching"],
  emits: ["setRM", "saveConfig", "startSketch", "stopSketch", "closeDrawer"],//用于向父组件传递信息，父组件的语法为@setRM="OnSetReaderMode"
  components: {
    NDrawer,//抽屉，可以从上下左右4个方向冒出. https://www.naiveui.com/zh-CN/os-theme/components/drawer
    NDrawerContent,//抽屉内容
    NSpace,//间距 https://www.naiveui.com/zh-CN/os-theme/components/space
    NDivider,//间隔
    // NRadioGroup,//单选  https://www.naiveui.com/zh-CN/os-theme/components/radio
    // NRadioButton,//单选 用按钮显得更优雅一点
    NButton,//按钮，来自:https://www.naiveui.com/zh-CN/os-theme/components/button
    // NAvatar,//头像 https://www.naiveui.com/zh-CN/os-theme/components/avatar
    NSelect, //选择器 https://www.naiveui.com/zh-CN/os-theme/components/select
    NPopconfirm, //弹出确认 https://www.naiveui.com/zh-CN/os-theme/components/popconfirm
  },
  setup() {
    const { cookies } = useCookies();
    // const message = useMessage(); 需要导入 'naive-ui'的useMessage

    return {
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
      readModeLocal: "",
    };
  },
  //挂载前
  beforeMount() {
    var lang = this.cookies.get("userLanguageSetting");
    if (lang) {
      this.$i18n.locale = lang;
    }
    this.readModeLocal = this.readMode;
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
    OnChangeLanguage(value) {
      this.cookies.set("userLanguageSetting", value);
    },
    // 关闭抽屉时，保存设置到cookies
    saveConfigToCookie(show) {
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




