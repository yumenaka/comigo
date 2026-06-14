/**
 * 全局 SSE：接收 ui_suggest_reload（整页刷新通知）并转发 log 到设置页日志面板。
 */
function shouldEnableComigoSSE() {
    if (typeof window === 'undefined' || typeof EventSource === 'undefined') {
        return false
    }
    // 登录页没有 JWT，会导致 /api/sse 持续 401 重连
    const pathname = window.ComiGoRelativePath ? window.ComiGoRelativePath(window.location.pathname) : window.location.pathname
    return pathname !== '/login'
}

const libraryRescanReloadReasons = new Set([
    'library_rescan_done',
    'auto_library_rescan_done',
    'single_store_rescan_done',
])

// 仅在书架与设置页处理整页刷新；阅读页（flip/scroll 等）不打断
function shouldShowUISuggestReloadPrompt() {
    const p = window.ComiGoRelativePath ? window.ComiGoRelativePath(window.location.pathname) : window.location.pathname
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

function autoReloadAfterLibraryRescan() {
    if (!shouldShowUISuggestReloadPrompt() || window.__comigoAutoReloadQueued) {
        return
    }
    // 若本页正在发起书库变更请求，先别 reload：
    // 否则 location.reload() 会把这个还没返回的 fetch 直接 abort 掉，
    // Firefox 报 “TypeError: NetworkError when attempting to fetch resource”，
    // 让明明扫描成功的操作弹出“网络错误，请重试”。
    // 改为挂起，等该 fetch 在自己的 finally 里结束后再刷新。
    if (window.__comigoRescanInFlight) {
        window.__comigoReloadPending = true
        return
    }
    window.__comigoAutoReloadQueued = true
    reloadComigoPage()
}

// 由发起重扫的页面在 fetch 结束(finally)后调用：执行被挂起的整页刷新。
function comigoRunPendingReload() {
    if (window.__comigoReloadPending && !window.__comigoAutoReloadQueued) {
        window.__comigoReloadPending = false
        window.__comigoAutoReloadQueued = true
        reloadComigoPage()
    }
}

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

function reloadComigoPage() {
    closeComigoSSE()
    window.location.reload()
}

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
            '<span style="color:oklch(62.7% 0.194 149.214)">[message]</span>' + e.data
        )
    }

    es.onopen = () => {
        const text =
            typeof i18next !== 'undefined' && i18next.t
                ? i18next.t('settings_log_sse_connected')
                : 'SSE connected'
        appendSharedLog(
            '<span style="color:oklch(62.7% 0.194 149.214)">[open]</span> ' + text
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
            '<span style="color:oklch(57.7% 0.245 27.325)">[error]</span> ' + text
        )
    }
}

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
    const sseURL = window.ComiGoPath ? window.ComiGoPath('/api/sse') : '/api/sse'
    const es = new EventSource(sseURL, { withCredentials: true })
    window.__comigoSSEInstance = es
    comigoAttachSSEListeners(es)
    return es
}

window.__comigoSSEInit = comigoSSEInit
window.__comigoRunPendingReload = comigoRunPendingReload

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

window.dispatchEvent(new Event('comigo:sse-ready'))
