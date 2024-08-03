import 'htmx.org'
import Alpine from 'alpinejs'
import persist from '@alpinejs/persist'

// Add Alpine instance to window object.
window.Alpine = Alpine

// Alpine Persist 插件，用于持久化存储。默认存储到 localStorage。
// 详细用法参见： https://alpinejs.dev/plugins/persist
Alpine.plugin(persist)


// Start Alpine.
Alpine.start()
