// VUEX TypeScript 支持
// https://vuex.vuejs.org/zh/guide/typescript-support.html

import { ComponentCustomProperties } from "@/vue";
import Store from "@/store"; //VueX

declare module "@vue/runtime-core" {
  // declare your own store states
  interface State {
    orbital;
  }

  // provide typings for `this.$store`
  interface ComponentCustomProperties {
    $store: Store<State>;
  }
}



