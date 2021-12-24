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
        <span>{{ this.$t('message.ReaderSettings') }}</span>
        <n-avatar size="small" src="/favicon.ico" />
      </template>
      <!-- 选择：切换页面模式 -->
      <n-space>
        <n-radio-group v-model:value="nowTemplateLocal">
          <!-- 卷軸模式 -->
          <n-radio-button
            :checked="nowTemplateLocal === 'scroll'"
            @change="onChangeTemplate"
            value="scroll"
            name="basic-demo"
          >{{ this.$t('message.scroll_mode') }}</n-radio-button>
          <!-- 翻頁模式 -->
          <n-radio-button
            :checked="nowTemplateLocal === 'flip'"
            @change="onChangeTemplate"
            value="flip"
            name="basic-demo"
          >{{ this.$t('message.flip_mode') }}</n-radio-button>
        </n-radio-group>
      </n-space>
      <!-- 分割线 -->
      <n-divider />
      <!-- 父组件在此处插入自定义内容 -->
      <slot></slot>

      <!-- 抽屉：自定义底部 -->
      <template #footer>
        <!-- <n-space vertical>
          <n-select
            v-model:value="this.locale"
            :options="this.languageOptions"
            @on-update:value="onchangeLanguage"
          />
        </n-space>-->
        <n-button
          v-if="nowTemplateDrawer == 'flip' || nowTemplateDrawer == 'scroll'"
          @click="startSketchMode"
        >{{ this.$t('message.startSketchMode') }}</n-button>
        <n-button
          v-if="nowTemplateDrawer == 'sketch'"
          @click="stopSketchMode"
        >{{ this.$t('message.stopSketchMode') }}</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<script>

import { useCookies } from "vue3-cookies";
import { NDrawer,NDivider, NDrawerContent, NSpace, NRadioGroup, NRadioButton, NAvatar, NButton } from 'naive-ui'
import { defineComponent, } from 'vue'
// import { useI18n } from 'vue-i18n'

export default defineComponent({
  name: "Drawer",
  props: ['book', 'initDrawerActive', 'initDrawerPlacement', 'nowTemplateDrawer'],
  emits: ["setT", "saveConfig", "startSketch", "stopSketch", "closeDrawer"],
  components: {
    NDrawer,//抽屉，可以从上下左右4个方向冒出. https://www.naiveui.com/zh-CN/os-theme/components/drawer
    NDrawerContent,//抽屉内容
    NSpace,
    NDivider,//间隔
    NRadioGroup,//单选  https://www.naiveui.com/zh-CN/os-theme/components/radio
    NRadioButton,//单选 用按钮显得更优雅一点
    NButton,//按钮，来自:https://www.naiveui.com/zh-CN/os-theme/components/button
    NAvatar,//头像 https://www.naiveui.com/zh-CN/os-theme/components/avatar
    // NSelect, //选择器 https://www.naiveui.com/zh-CN/os-theme/components/select
  },
  setup() {
    const { cookies } = useCookies();
    // const { locale } = useI18n();
    // const onchangeLanguage = (value) => {
    //   //切换语言
    //   locale.value = value
    // }
    return {
      cookies,
      // onchangeLanguage,
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
          label: '简体中文',
          value: 'zh_CN'
        },
      ],
    };
  },
  data() {
    return {
      someflag: "",
      nowTemplateLocal: "",
    };
  },
  //挂载前
  beforeMount() {
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
      // this.$emit("greet", this.nowTemplateLocal);
      // console.log("onChangeTemplate:"+value);
      if (this.nowTemplateLocal === "scroll") {
        this.$emit("setT", "scroll");
      }
      if (this.nowTemplateLocal === "flip") {
        this.$emit("setT", "flip");
      }
      if (this.nowTemplateLocal === "sketch") {
        this.$emit("setT", "sketch");
      }
      // location.reload(); //需要刷新？ 以后研究VueRouter并去掉
    },
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>




