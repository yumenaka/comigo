<template>
  <div class="bottom text-gray-900 h-12 py-4 space-x-2 text-base content-center">
    <span class="text-base">Power by</span>
    <a
      href="https://github.com/yumenaka/comi/releases"
      class="text-blue-700 hover:underline font-semibold"
    >{{ softVersion }}</a>
    <!-- <span>
      <a
        href="https://www.yumenaka.net"
        class="text-blue-700 hover:underline text-base"
      >www.yumenaka.net &nbsp;</a>
    </span>-->
    <!-- <a href="https://github.com/yumenaka/comi/issues/new" class="text-gray-900">File an Issue</a> -->
  </div>
</template>

<script lang ="ts">
import { defineComponent } from 'vue'
export default defineComponent({
  name: "ComigoBottom",
  props: ['softVersion',],
  components: {

  },
  setup() {

  },
  data() {
    return {
      someflag: "",
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




