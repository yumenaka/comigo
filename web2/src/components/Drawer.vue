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
      </template>
      <!-- 选择：切换页面模式 -->
      <n-space>
        <n-radio-group v-model:value="selectedTemplate">
          <n-radio-button
            :checked="selectedTemplate === 'scroll'"
            @change="onChangeTemplate"
            value="scroll"
            name="basic-demo"
          >{{ this.$t('message.scroll_mode') }}</n-radio-button>
          <n-radio-button
            :checked="selectedTemplate === 'flip'"
            @change="onChangeTemplate"
            value="flip"
            name="basic-demo"
          >{{ this.$t('message.flip_mode') }}</n-radio-button>
        </n-radio-group>
      </n-space>

      <slot></slot>

      <!-- 抽屉：自定义底部 -->
      <template #footer>
        <n-button @click="startSketchMode">{{ this.$t('message.startSketchMode') }}</n-button>
        <n-avatar size="small" src="/favicon.ico" />
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<script>
import { useCookies } from "vue3-cookies";
import { NDrawer, NDrawerContent, NSpace, NRadioGroup,NRadioButton, NAvatar, NButton, } from 'naive-ui'
import { defineComponent, } from 'vue'
export default defineComponent({
  name: "Drawer",
  props: ['book', 'initDrawerActive', 'initDrawerPlacement'],
  components: {
    NDrawer,//抽屉，可以从上下左右4个方向冒出. https://www.naiveui.com/zh-CN/os-theme/components/drawer
    NDrawerContent,//抽屉内容
    NSpace,
    NRadioGroup,//单选  https://www.naiveui.com/zh-CN/os-theme/components/radio
    NRadioButton,//单选 用按钮显得更优雅一点
    NButton,//按钮，来自:https://www.naiveui.com/zh-CN/os-theme/components/button
    NAvatar,//头像 https://www.naiveui.com/zh-CN/os-theme/components/avatar
    // NSwitch, 
  },
  setup() {
    const { cookies } = useCookies();
    return {
      cookies,
    };
  },
  data() {
    return {
      someflag: "",
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
      if (show==false){
        this.$emit('closeDrawer');
        this.$emit('saveConfig');
      }
    },
    startSketchMode() {
      this.$emit('startSketch');
    },
    //切换模板的函数，需要配合vue-router
    onChangeTemplate() {
      if (this.selectedTemplate === "scroll") {
        this.cookies.set("nowTemplate", "scroll");
      }
      if (this.selectedTemplate === "flip") {
        this.cookies.set("nowTemplate", "flip");
      }
      if (this.selectedTemplate === "sketch") {
        this.cookies.set("nowTemplate", "sketch");
      }
      location.reload(); //暂时无法动态刷新，研究vue-router去掉
    },
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>




