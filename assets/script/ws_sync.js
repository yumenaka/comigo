/**
 * ComiGo 共享 WebSocket 同步模块
 * 供 flip / scroll / shelf 等页面复用，连接 /api/ws，发送时附带 page_type 与可选的 book_id。
 * 参考: https://www.ruanyifeng.com/blog/2017/05/websocket.html
 */
(function () {
    'use strict'

    const token = 'your_token' // TODO: 使用真正的令牌
    const tabID = (Date.now() % 10000000).toString(36) + Math.random().toString(36).substring(2, 5)

    let socket = null
    let reconnectAttempts = 0
    let reconnectTimer = null
    let isIntentionallyClosed = false
    let options = null // { pageType, getBookId, getWsConfig, isDebug, onMessage }

    function getConfig() {
        if (!options || !options.getWsConfig) return { maxReconnectAttempts: 200, reconnectInterval: 3000 }
        return options.getWsConfig()
    }

    function isDebug() {
        return options && typeof options.isDebug === 'function' && options.isDebug()
    }

    function handleOpen() {
        console.log('WebSocket连接已建立')
        reconnectAttempts = 0
        isIntentionallyClosed = false
    }

    function handleMessage(event) {
        try {
            const message = JSON.parse(event.data)
            if (options && typeof options.onMessage === 'function') {
                options.onMessage(message)
            }
        } catch (error) {
            console.error('WebSocket 消息解析失败:', error)
        }
    }

    function handleClose() {
        console.log('WebSocket连接已关闭')
        if (!isIntentionallyClosed) {
            scheduleReconnect()
        }
    }

    function handleError(error) {
        console.error('WebSocket发生错误：', error)
        if (socket) {
            socket.close()
        }
    }

    function scheduleReconnect() {
        const config = getConfig()
        if (reconnectAttempts >= config.maxReconnectAttempts) {
            console.log('已达到最大重连次数，停止重连')
            return
        }
        reconnectAttempts++
        if (isDebug()) {
            console.log(`第 ${reconnectAttempts} 次重连...`)
        }
        reconnectTimer = setTimeout(() => {
            connect()
        }, config.reconnectInterval)
    }

    const ComiGoWS = {
        /**
         * 初始化（仅允许调用一次，或先 disconnect 后再 init）
         * @param {Object} opts - pageType: 'flip'|'scroll'|'shelf', getBookId?: ()=>string|undefined, getWsConfig: ()=>({maxReconnectAttempts, reconnectInterval}), isDebug?: ()=>boolean, onMessage: (message)=>void
         */
        init(opts) {
            if (!opts || !opts.pageType || !opts.getWsConfig || !opts.onMessage) {
                console.error('ComiGoWS.init 需要 pageType、getWsConfig、onMessage')
                return
            }
            options = opts
        },

        /**
         * 建立 WebSocket 连接
         */
        connect() {
            if (socket && (socket.readyState === WebSocket.CONNECTING || socket.readyState === WebSocket.OPEN)) {
                if (isDebug()) {
                    console.log('WebSocket 正在连接或已打开，跳过')
                }
                return
            }

            const config = getConfig()
            const wsProtocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://'
            const wsUrl = wsProtocol + window.location.host + '/api/ws'

            try {
                socket = new WebSocket(wsUrl)
                socket.onopen = () => handleOpen()
                socket.onmessage = (event) => handleMessage(event)
                socket.onclose = () => handleClose()
                socket.onerror = (error) => handleError(error)
            } catch (error) {
                console.error('WebSocket 连接失败:', error)
                scheduleReconnect()
            }
        },

        /**
         * 发送消息。自动附带 page_type、book_id（若 getBookId 有返回值）、tab_id、token
         * @param {string} type - 消息类型，如 'flip_mode_sync_page'、'heartbeat'
         * @param {Object} data - 会序列化为 data_string；同时 page_type 与可选的 book_id 放在顶层
         * @param {string} [detail] - 可选说明
         */
        send(type, data, detail) {
            if (!socket || socket.readyState !== WebSocket.OPEN) {
                if (isDebug()) {
                    console.log('WebSocket 未连接或未准备好，无法发送消息')
                }
                return
            }

            const bookId = options && typeof options.getBookId === 'function' ? options.getBookId() : undefined
            const msg = {
                type,
                status_code: 200,
                tab_id: tabID,
                Token: token,
                detail: detail || '',
                data_string: data != null ? JSON.stringify(data) : '',
                page_type: options ? options.pageType : undefined,
                book_id: bookId,
            }

            try {
                socket.send(JSON.stringify(msg))
            } catch (error) {
                console.error('WebSocket 发送消息失败:', error)
            }
        },

        /**
         * 主动断开连接
         */
        disconnect() {
            isIntentionallyClosed = true
            if (reconnectTimer) {
                clearTimeout(reconnectTimer)
                reconnectTimer = null
            }
            if (socket) {
                socket.close()
                socket = null
            }
        },

        /**
         * 获取当前标签页 ID（用于 onMessage 中判断是否为本页消息）
         */
        getTabId() {
            return tabID
        },
    }

    window.ComiGoWS = ComiGoWS
})()
