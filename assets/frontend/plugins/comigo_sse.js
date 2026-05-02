/**
 * 全局 SSE：接收 ui_suggest_reload（整页刷新建议）并转发 log 到设置页日志面板。
 */
function shouldEnableComigoSSE() {
    if (typeof window === 'undefined' || typeof EventSource === 'undefined') {
        return false
    }
    // 登录页没有 JWT，会导致 /api/sse 持续 401 重连
    const pathname = window.ComiGoRelativePath ? window.ComiGoRelativePath(window.location.pathname) : window.location.pathname
    return pathname !== '/login'
}

// 仅在书架与设置页弹出「建议刷新」；阅读页（flip/scroll 等）不打断
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
            window.location.reload()
        },
        onCancel: () => {
            window.__comigoReloadPromptOpen = false
        },
    })
}

function appendSharedLog(line) {
    if (typeof window.__comigoLogAppend === 'function') {
        window.__comigoLogAppend(line)
    }
}

function queueComigoSSEStart() {
    if (window.__comigoSSEStartQueued) {
        return
    }
    window.__comigoSSEStartQueued = true
    const start = () => {
        setTimeout(() => {
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
window.dispatchEvent(new Event('comigo:sse-ready'))
