/**
 * 复用一个全局 SSE 连接接收后端事件：处理界面刷新通知，并把日志转交给设置页日志面板。
 * 后端的 ui_suggest_reload 与 log 事件来自 tools/sse_hub；日志面板通过 __comigoLogAppend 接入同一连接。
 */
// 登录页没有 JWT，且旧浏览器可能没有 EventSource；这两种情况都不建立连接。
function shouldEnableComigoSSE() {
    if (typeof window === 'undefined' || typeof EventSource === 'undefined') {
        return false
    }
    const pathname = window.ComiGoRelativePath
        ? window.ComiGoRelativePath(window.location.pathname)
        : window.location.pathname
    return pathname !== '/login'
}

// 这三类通知表示书库数据已经完成重扫，书架和设置页可直接刷新，无需再次确认。
const libraryRescanReloadReasons = new Set([
    'library_rescan_done',
    'auto_library_rescan_done',
    'single_store_rescan_done',
])

// 仅在展示书库或设置数据的页面处理整页刷新；阅读页不应被后台通知打断。
function shouldShowUISuggestReloadPrompt() {
    const p = window.ComiGoRelativePath
        ? window.ComiGoRelativePath(window.location.pathname)
        : window.location.pathname
    if (p === '/settings') {
        return true
    }
    if (p === '/' || p === '/index.html' || p === '/search') {
        return true
    }
    if (p.startsWith('/shelf/')) {
        return true
    }
    return false
}

function isLibraryRescanReloadReason(reason) {
    return libraryRescanReloadReasons.has(reason)
}

// reason 与 locale 中 ui_suggest_reload_reason_* 的后缀一致，缺少翻译时回退到通用提示。
function getReloadPromptMessage(reason) {
    const key = 'ui_suggest_reload_reason_' + reason
    const translated =
        typeof i18next !== 'undefined' && i18next.t ? i18next.t(key) : key
    if (translated && translated !== key) {
        return translated
    }
    return typeof i18next !== 'undefined' && i18next.t
        ? i18next.t('ui_suggest_reload_default')
        : 'Data was updated on the server. Reload the page to see the latest UI?'
}

// 非书库重扫通知使用项目现有的 showMessage 确认框，且同一时间只显示一个刷新提示。
function showReloadPrompt(reason) {
    if (!shouldShowUISuggestReloadPrompt()) {
        return
    }
    if (typeof showMessage !== 'function' || window.__comigoReloadPromptOpen) {
        return
    }
    window.__comigoReloadPromptOpen = true
    showMessage({
        message: getReloadPromptMessage(reason),
        buttons: 'confirm_cancel',
        onConfirm: () => {
            window.__comigoReloadPromptOpen = false
            reloadComigoPage()
        },
        onCancel: () => {
            window.__comigoReloadPromptOpen = false
        },
    })
}

// 书库重扫完成后，书架和设置页直接刷新；全局标记避免同批通知重复触发 reload。
function autoReloadAfterLibraryRescan() {
    if (!shouldShowUISuggestReloadPrompt() || window.__comigoAutoReloadQueued) {
        return
    }
    window.__comigoAutoReloadQueued = true
    reloadComigoPage()
}

// 设置页日志面板尚未挂载时直接丢弃日志；连接本身仍继续接收后续事件。
function appendSharedLog(line) {
    if (typeof window.__comigoLogAppend === 'function') {
        window.__comigoLogAppend(line)
    }
}

// 取消尚未执行的延迟连接，页面即将卸载时不能再新建 EventSource。
function clearQueuedComigoSSEStart() {
    if (window.__comigoSSEStartTimer) {
        clearTimeout(window.__comigoSSEStartTimer)
        window.__comigoSSEStartTimer = null
    }
    window.__comigoSSEStartQueued = false
}

// 主动关闭当前 SSE；由 reload/pagehide 共用，避免卸载时遗留被浏览器标记为中断的长连接。
function closeComigoSSE() {
    clearQueuedComigoSSEStart()
    if (!window.__comigoSSEInstance) {
        return
    }
    try {
        window.__comigoSSEInstance.close()
    } catch (_) {}
    window.__comigoSSEInstance = null
}

// 刷新前先主动关闭长连接，避免浏览器把正常卸载记录为 SSE 请求异常。
function reloadComigoPage() {
    closeComigoSSE()
    window.location.reload()
}

// 等待页面 load 后再延迟一秒连接；pageshow 恢复时也复用这套去重逻辑。
function queueComigoSSEStart() {
    if (window.__comigoSSEStartQueued) {
        return
    }
    window.__comigoSSEStartQueued = true
    const start = () => {
        window.__comigoSSEStartTimer = setTimeout(() => {
            window.__comigoSSEStartTimer = null
            window.__comigoSSEStartQueued = false
            comigoSSEInit()
        }, 1000)
    }
    if (document.readyState === 'complete') {
        start()
    } else {
        window.addEventListener('load', start, { once: true })
    }
}

// 集中处理界面刷新、日志、显式 tick 与默认 message 事件；默认 message 额外显示全局提示。
function comigoAttachSSEListeners(es) {
    es.addEventListener('ui_suggest_reload', (e) => {
        let reason = 'default'
        try {
            const data = JSON.parse(e.data || '{}')
            if (data.reason) {
                reason = data.reason
            }
        } catch (_) {}
        if (isLibraryRescanReloadReason(reason)) {
            autoReloadAfterLibraryRescan()
            return
        }
        showReloadPrompt(reason)
    })

    es.addEventListener('log', (e) => {
        appendSharedLog(e.data)
    })

    es.addEventListener('tick', (e) => {
        appendSharedLog('[tick] ' + e.data)
    })

    es.onmessage = (e) => {
        if (typeof showToast === 'function') {
            showToast(e.data, 'info')
        }
        appendSharedLog(
            '<span style="color:oklch(62.7% 0.194 149.214)">[message]</span>' +
                e.data,
        )
    }

    es.onopen = () => {
        const text =
            typeof i18next !== 'undefined' && i18next.t
                ? i18next.t('settings_log_sse_connected')
                : 'SSE connected'
        appendSharedLog(
            '<span style="color:oklch(62.7% 0.194 149.214)">[open]</span> ' +
                text,
        )
    }

    es.onerror = () => {
        const closed =
            typeof EventSource !== 'undefined' &&
            es.readyState === EventSource.CLOSED
        const text = closed
            ? typeof i18next !== 'undefined' && i18next.t
                ? i18next.t('settings_log_sse_closed')
                : 'closed'
            : typeof i18next !== 'undefined' && i18next.t
              ? i18next.t('settings_log_sse_retrying')
              : 'retrying'
        appendSharedLog(
            '<span style="color:oklch(57.7% 0.245 27.325)">[error]</span> ' +
                text,
        )
    }
}

// 返回当前可用连接；只有旧连接已关闭时才创建新的 EventSource。
function comigoSSEInit() {
    if (!shouldEnableComigoSSE()) {
        return null
    }
    if (window.__comigoSSEInstance) {
        if (window.__comigoSSEInstance.readyState === EventSource.CLOSED) {
            window.__comigoSSEInstance = null
        } else {
            return window.__comigoSSEInstance
        }
    }
    if (window.__comigoSSEStartQueued) {
        return window.__comigoSSEInstance
    }
    // 页面初次加载时稍后再连，避免浏览器把 SSE 长连接误报为加载中断。
    if (document.readyState !== 'complete') {
        queueComigoSSEStart()
        return null
    }
    const sseURL = window.ComiGoPath
        ? window.ComiGoPath('/api/sse')
        : '/api/sse'
    const es = new EventSource(sseURL, { withCredentials: true })
    window.__comigoSSEInstance = es
    comigoAttachSSEListeners(es)
    return es
}

// 日志面板可能晚于主包脚本执行，因此暴露初始化入口供其确认共享连接已经建立。
window.__comigoSSEInit = comigoSSEInit

// 全局启动 SSE；具体事件处理仍由上面的路径判断决定，阅读页不会被重扫通知打断。
queueComigoSSEStart()

// 页面卸载时主动关闭 SSE，并取消尚未执行的延迟启动，避免卸载过程中创建/留下
// 被中断(aborted)的 /api/sse 请求。
if (typeof window.addEventListener === 'function') {
    window.addEventListener('pagehide', () => {
        closeComigoSSE()
    })
    window.addEventListener('pageshow', () => {
        queueComigoSSEStart()
    })
}

// 通知设置页日志面板：全局初始化入口已经可以调用。
window.dispatchEvent(new Event('comigo:sse-ready'))
