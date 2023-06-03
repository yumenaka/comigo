/* eslint-disable */
declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

// 声明文件必需以 .d.ts 为后缀。
// 一般来说，ts 会解析项目中所有的 *.ts 文件，当然也包含以 .d.ts 结尾的文件。

// // 声明全局属性类型
// declare module "@vue/runtime-core" {
//   interface ComponentCustomProperties<T> {
//     $connect: (url?: string) => void;
//   }
// }