import Alpine from 'alpinejs'
import persist from '@alpinejs/persist'

// 全站状态和浏览器偏好只需 Alpine 核心与持久化插件。
window.Alpine = Alpine
Alpine.plugin(persist)
