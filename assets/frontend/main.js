//此文件需要编译，编译指令请参考 package.json
// 基础插件
import './plugins/i18n' // 这种 import 通常用于单纯执行该文件内的脚本或配置逻辑，确保它在程序启动或相关流程中被"触发"过。
import './plugins/alpine'

import './plugins/screenfull'
import './plugins/ui_controls'
import './stores/global_store'
import './stores/shelf_store'
import './stores/scroll_store'
import './stores/flip_store'
import './plugins/sse'

// Start Alpine.
Alpine.start()
