// reader 页面专用 PWA 支持：注册 Service Worker，并处理安装按钮。
'use strict'

let readerInstallPromptEvent = null

function getReaderInstallButton() {
    return document.getElementById('ReaderInstallPWAButton')
}

function readerPWAText(key, fallback) {
    if (typeof i18next !== 'undefined') {
        const text = i18next.t(key)
        if (text && text !== key) return text
    }
    return fallback
}

function notifyReaderPWA(message, type = 'info') {
    if (typeof showToast === 'function') {
        showToast(message, type)
        return
    }
    alert(message)
}

function isReaderPWAStandalone() {
    return window.matchMedia('(display-mode: standalone)').matches || window.navigator.standalone === true
}

function isReaderPWAInstallProtocolAllowed() {
    // beforeinstallprompt 只会在安全上下文中可靠触发；这里保持按钮策略和 SW 注册一致。
    return window.location.protocol === 'https:'
}

function setReaderInstallButtonState() {
    const button = getReaderInstallButton()
    if (!button) return

    if (!isReaderPWAInstallProtocolAllowed()) {
        button.classList.add('hidden')
        return
    }

    button.classList.remove('hidden')
    // beforeinstallprompt 并不是稳定的“可用状态”信号；按钮保持可点击，点击时再决定弹原生安装框或提示手动安装。
    button.disabled = false
    button.title = readerInstallPromptEvent
        ? readerPWAText('reader_pwa_install_ready', 'Ready to add')
        : readerPWAText('reader_pwa_install_unavailable', 'Use the browser menu to add this page as an app, or refresh and try again.')
}

function initReaderPWAInstallButton() {
    const button = getReaderInstallButton()
    if (!button) return

    setReaderInstallButtonState()

    button.addEventListener('click', async () => {
        if (isReaderPWAStandalone()) {
            notifyReaderPWA(readerPWAText('reader_pwa_already_installed', 'This page is already running as an app.'), 'success')
            return
        }

        if (!readerInstallPromptEvent) {
            notifyReaderPWA(
                readerPWAText('reader_pwa_install_unavailable', 'Use the browser menu to add this page as an app, or refresh and try again.'),
                'warning',
            )
            return
        }

        readerInstallPromptEvent.prompt()
        await readerInstallPromptEvent.userChoice
        readerInstallPromptEvent = null
        setReaderInstallButtonState()
    })
}

function registerReaderServiceWorker() {
    if (!('serviceWorker' in navigator)) return
    if (!window.isSecureContext) return

    const register = () => {
        // scope 固定在 reader 页面下，避免离线缓存影响普通书架、在线阅读等页面。
        const swPath = window.ComiGoPath ? window.ComiGoPath('/reader-sw.js') : '/reader-sw.js'
        const swScope = window.ComiGoPath ? window.ComiGoPath('/reader') : '/reader'
        navigator.serviceWorker.register(swPath, { scope: swScope }).catch((error) => {
            console.error('[reader-pwa] register service worker failed:', error)
        })
    }

    if (document.readyState === 'complete') {
        register()
    } else {
        window.addEventListener('load', register, { once: true })
    }
}

window.addEventListener('beforeinstallprompt', (event) => {
    event.preventDefault()
    readerInstallPromptEvent = event
    setReaderInstallButtonState()
})

window.addEventListener('appinstalled', () => {
    readerInstallPromptEvent = null
    setReaderInstallButtonState()
    notifyReaderPWA(readerPWAText('reader_pwa_install_completed', 'App added.'), 'success')
})

document.addEventListener('DOMContentLoaded', () => {
    initReaderPWAInstallButton()
    registerReaderServiceWorker()
})
