import Alpine from 'alpinejs'
import persist from '@alpinejs/persist'
import morph from '@alpinejs/morph'

window.Alpine = Alpine // 将 Alpine 实例添加到窗口对象中。
Alpine.plugin(persist)
Alpine.plugin(morph) 