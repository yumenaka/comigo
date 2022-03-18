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
        <n-radio-group v-model:value="nowTemplateLocal">
          <!-- 卷軸模式 -->
          <n-radio-button
            :checked="nowTemplateLocal === 'scroll'"
            @change="onChangeTemplate"
            value="scroll"
            name="basic-demo"
          >{{ $t('scroll_mode') }}</n-radio-button>
          <!-- 翻頁模式 -->
          <n-radio-button
            :checked="nowTemplateLocal === 'flip'"
            @change="onChangeTemplate"
            value="flip"
            name="basic-demo"
          >{{ $t('flip_mode') }}</n-radio-button>
        </n-radio-group>
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
          v-if="nowTemplate == 'flip' || nowTemplate == 'scroll'"
          @click="startSketchMode"
        >{{ $t('startSketchMode') }}</n-button>
        <n-button v-if="nowTemplate == 'sketch'" @click="stopSketchMode">{{ $t('stopSketchMode') }}</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<script>

import { useCookies } from "vue3-cookies";
import { NDrawer, NDivider, NDrawerContent, NSpace, NRadioGroup, NRadioButton, NButton, NSelect, NPopconfirm, } from 'naive-ui'
import { defineComponent, } from 'vue'
// import { useI18n } from 'vue-i18n'


export default defineComponent({
  name: "Drawer",
  props: ['book', 'initDrawerActive', 'initDrawerPlacement', 'nowTemplate'],
  emits: ["setT", "saveConfig", "startSketch", "stopSketch", "closeDrawer"],
  components: {
    NDrawer,//抽屉，可以从上下左右4个方向冒出. https://www.naiveui.com/zh-CN/os-theme/components/drawer
    NDrawerContent,//抽屉内容
    NSpace,//间距 https://www.naiveui.com/zh-CN/os-theme/components/space
    NDivider,//间隔
    NRadioGroup,//单选  https://www.naiveui.com/zh-CN/os-theme/components/radio
    NRadioButton,//单选 用按钮显得更优雅一点
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
      nowTemplateLocal: "",
    };
  },
  //挂载前
  beforeMount() {
    var lang = this.cookies.get("userLanguageSetting");
    if (lang) {
      this.$i18n.locale = lang;
    }
    this.nowTemplateLocal = this.nowTemplate;
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

    onClearLocalSetting() {
      // //清除localStorage保存的设定
      // localStorage.clear();
      // //刷新当前页面
      // location.reload();
      //载入新文档替换当前页面
      //window.location.replace("http://www.example.com")
    },

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




