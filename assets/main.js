//此文件需要编译，编译指令请参考 package.json
import 'htmx.org'
import 'flowbite'
// 基础插件
import './plugins/i18n' // 这种 import 通常用于单纯执行该文件内的脚本或配置逻辑，确保它在程序启动或相关流程中被"触发"过。
import './plugins/alpine'
import './plugins/screenfull'
// 声明各种变量
import './stores/cookie_store'
import './stores/global_store'
import './stores/shelf_store'
import './stores/scroll_store'
import './stores/flip_store'
import './stores/theme_store'
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