// 扩展全局属性
// 某些插件会通过 app.config.globalProperties 为所有组件都安装全局可用的属性。
// 举例来说，我们可能为了请求数据而安装了 this.$http，或者为了国际化而安装了 this.$translate。
// 为了使 TypeScript 更好地支持这个行为，Vue 暴露了一个被设计为可以通过 TypeScript 模块扩展来扩展的 ComponentCustomProperties 接口.
// https://cn.vuejs.org/guide/typescript/options-api.html#augmenting-global-properties

import { ComponentCustomProperties } from 'vue'
import { WebSocket } from 'vue-native-websocket-vue3'

declare module '@vue/runtime-core' {
    interface ComponentCustomProperties {
        $socket : typeof WebSocket
    }
}