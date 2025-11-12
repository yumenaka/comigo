import Alpine from 'alpinejs'
import persist from '@alpinejs/persist'
import morph from '@alpinejs/morph'
import intersect from '@alpinejs/intersect'

window.Alpine = Alpine // 将 Alpine 实例添加到窗口对象中。
Alpine.plugin(persist) // 用于在本地存储中持久化数据的插件
Alpine.plugin(morph) // 不丢失 Alpine 页面状态的情况下，根据服务器请求更新 HTML
Alpine.plugin(intersect)//  Intersection Observer 的一个便捷封装，在元素进入视口时做出反应。