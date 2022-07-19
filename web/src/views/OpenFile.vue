

<template>
  <div id="openfile">
    <!-- https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/Input/file -->
    <!-- https://blog.jijian.link/2020-01-08/html-upload-folder/ -->
    <!-- <input type="file"></input> -->
  </div>
</template>

<script>
import { defineComponent, reactive } from 'vue'
import { useCookies } from "vue3-cookies";// https://github.com/KanHarI/vue3-cookies


export default defineComponent({
  name: "AboutPage",
  props: ['readMode'],
  emits: ["setSome"],
  components: {
  },
  setup() {
    //此处不能使用this
    const { cookies } = useCookies();
    //背景颜色，颜色选择器用
    const model = reactive({
      color: "#E0D9CD",
      colorHeader: "#d1c9c1",
    });
    //单选按钮绑定的数值
    // const checkedValueRef = ref(null)
    return {
      cookies,
      //背景色
      model,
    }
  },
  data() {
    return {
      readerMode: "",
    };
  },
  //Vue3生命周期:  https://v3.cn.vuejs.org/api/options-lifecycle-hooks.html#beforecreate
  // created : 在绑定元素的属性或事件监听器被应用之前调用。
  // beforeMount : 指令第一次绑定到元素并且在挂载父组件之前调用。
  // mounted : 在绑定元素的父组件被挂载后调用。
  // beforeUpdate: 在更新包含组件的 VNode 之前调用。。
  // updated: 在包含组件的 VNode 及其子组件的 VNode 更新后调用。
  // beforeUnmount: 当指令与在绑定元素父组件卸载之前时，只调用一次。
  // unmounted: 当指令与元素解除绑定且父组件已卸载时，只调用一次。
  created() {
  },
  //挂载前
  beforeMount() {
  },
  onMounted() {
    // https://v3.cn.vuejs.org/api/options-lifecycle-hooks.html#beforemount
    this.$nextTick(function () {
      //视图渲染之后运行的代码
    })
  },
  //卸载前
  beforeUnmount() {

  },
  methods: {
    onOpenBook() {
    },
    computed: {
    }
  },
});
</script>

<style scoped>
#about {
  padding-bottom: 20px;
  padding-left: 20px;
  padding-right: 20px;
  padding-top: 20px;
  max-width: 100%;
  min-height: 93vh;
  height: 93vh;
  background: v-bind("model.color");
}
</style>
