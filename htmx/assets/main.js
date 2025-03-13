import 'htmx.org'
import 'flowbite'
// 这种 import './config/i18n' 通常用于单纯执行该文件内的脚本或配置逻辑，确保它在程序启动或相关流程中被“触发”过。
// 没有显式使用 import { ... } from ... 的原因就是：这个模块没提供需要拿来用的东西，而是仅仅为了执行模块内部的副作用代码。
import './config/i18n'
import './config/alpine'
import './stores/cookie'
import './stores/global'
import './stores/shelf'
import './stores/scroll'
import './stores/flip'
import './stores/theme'
import './plugins/screenfull'
import './utils/imageParameters'

// Start Alpine.
Alpine.start()

// Document ready function to ensure the DOM is fully loaded.
document.addEventListener('DOMContentLoaded', function () {
    initFlowbite() // initialize Flowbite
})

// Add event listeners for all HTMX events.
document.body.addEventListener(
    'htmx:afterSwap htmx:afterRequest htmx:afterSettle',
    function () {
        initFlowbite() // initialize Flowbite
    }
) 