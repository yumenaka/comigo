// VUEX TypeScript 支持
// https://vuex.vuejs.org/zh/guide/typescript-support.html

import { ComponentCustomProperties } from "@/vue";
import Store from "@/store"; //VueX

declare module "@vue/runtime-core" {
  // 声明自己的 store state
  interface State {
    count: number
  }

  // 为 `this.$store` 提供类型声明
  interface ComponentCustomProperties {
    $store: Store<State>;
  }
}



