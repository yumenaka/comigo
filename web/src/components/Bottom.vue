<template>
  <div class="bottom text-gray-900 h-12 py-4 space-x-2 text-base content-center">
    <span class="text-base">Power by</span>
    <a href="https://github.com/yumenaka/comi/releases" class="text-blue-700 hover:underline font-semibold">{{
      ServerName }}</a>
  </div>
</template>

<script lang ="ts">
import { defineComponent } from 'vue'
export default defineComponent({
  // name: "Bottom",
  props: ['ServerName'],
  data() {
    return {
      SomeFlag: "",
    };
  },
  methods: {
    //点击返回的时候，后退到上一页
    onClick() {
    },
  },
});
</script>

<style scoped>
</style>




