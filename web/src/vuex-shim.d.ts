// VUEX TypeScript 支持
// https://vuex.vuejs.org/zh/guide/typescript-support.html

import { ComponentCustomProperties } from "@/vue";
import Store from "@/store"; //VueX

declare module "@vue/runtime-core" {
  // declare your own store states
  interface State {
    orbital;
  }

  // interface Connect {
  //   orbital
  // }

  // provide typings for `this.$store`
  interface ComponentCustomProperties {
    $store: Store<State>;
    // $connect: Connect<Connect>
  }
}

// 扩展全局属性

import { AxiosInstance } from "axios";

declare module "@vue/runtime-core" {
  interface ComponentCustomProperties {
    $axios: AxiosInstance;
  }
}

// import Store from "@/store"; //VueX
// declare module '@vue/runtime-core' {
//   interface ComponentCustomProperties {
//     $store: Store
//   }
// }
