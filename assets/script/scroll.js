//此文件静态导入，不需要编译

// 使用标准 <script> 标记插入的 JavaScript 代码
'use strict'

//https://templ.guide/syntax-and-usage/script-templates/
//设置初始值
const book = JSON.parse(document.getElementById('NowBook').textContent)
const images = book.PageInfos
Alpine.store('scroll').allPageNum = parseInt(book.page_count)
// 用户ID和令牌，假设已在其他地方定义
const userID = Alpine.store('global').clientID
// 分页卷轴每页图片数（与后端 PAGED_SIZE 保持一致）
const PAGED_SIZE = 32
// 最大页码
const MaxPageNum = Math.floor(parseInt(book.page_count) / PAGED_SIZE) + 1

// ====== 卷轴模式 WebSocket 同步常量 ======
const SCROLL_SYNC_PENDING_KEY = `comigo_scroll_sync_${book.id}` // 跨页导航时暂存远端同步数据的 sessionStorage 键
const SCROLL_REMOTE_SUPPRESS_MS = 1800       // 接收远端同步后，抑制本端回传的时长（ms）
const SCROLL_SYNC_PERCENT_THRESHOLD = 0.08   // 同页内触发同步发送的最小 percent 变化量
const SCROLL_SYNC_THROTTLE_MS = 300          // 同步发送的最短间隔（ms）
const SCROLL_REMOTE_SYNC_SETTLE_MS = 60      // 远端动画结束后，延迟更新中心追踪的等待时长
const SCROLL_REMOTE_FOLLOW_IDLE_MS = 140     // 普通位置动画吸附前，需目标稳定的最短空闲时长
const SCROLL_REMOTE_FOLLOW_SNAP_PX = 1.5     // 普通位置动画吸附的距离阈值（px）
const SCROLL_REMOTE_EDGE_SNAP_PX = 3         // 顶底边缘动画吸附的距离阈值（px），比普通更宽容
const SCROLL_DOCUMENT_EDGE_TOLERANCE_PX = 2  // 判断文档是否处于顶部/底部的容差（px）
const SCROLL_WS_CONFIG = {
    maxReconnectAttempts: 200,
    reconnectInterval: 3000,
}

// 卷轴同步运行时状态（非持久化）
const scrollSyncState = {
    autoScrollEnabled: false,                  // 自动滚动插件是否启用
    rafID: null,                               // 中心追踪 requestAnimationFrame ID
    suppressBroadcastUntil: 0,                 // 抑制回传的截止时间戳
    hasInitializedCenterTracking: false,        // 中心追踪是否已完成首次初始化
    lastTrackedPageNum: 0,                     // 上次追踪到的页码
    lastBookmarkedPageNum: 0,                  // 上次写入书签的页码
    lastSentPercent: 0,                        // 上次发送的阅读百分比
    lastSyncSendTime: 0,                       // 上次发送同步消息的时间戳
    remoteScrollAnimationFrameID: null,        // 远端跟随动画的 RAF ID（非 null 时动画进行中）
    remoteScrollAnimationTargetTop: 0,         // 当前动画目标 scrollTop
    remoteScrollAnimationLastFrameTime: 0,     // 动画上一帧的时间戳
    remoteScrollAnimationLastTargetAt: 0,      // 最近一次目标更新的时间戳
    remoteScrollAnimationTriggerSource: 'manual', // 动画触发来源（'auto'=自动滚动 / 'manual'=手动）
    remoteScrollSettleTimer: null,             // 动画结束后延迟更新中心追踪的定时器
    lastSentAtEdge: false,                     // 上次发送时是否处于文档边缘
}

// 将数值限制在 [min, max] 范围内
function clamp(value, min, max) {
    return Math.min(Math.max(value, min), max)
}

// 获取文档可滚动的最大 scrollTop 值
function getMaxScrollTop() {
    return Math.max(0, document.documentElement.scrollHeight - window.innerHeight)
}

// 判断文档是否已滚动到顶部（含容差）
function isDocumentAtTop() {
    return window.scrollY <= SCROLL_DOCUMENT_EDGE_TOLERANCE_PX
}

// 判断文档是否已滚动到底部（含容差）
function isDocumentAtBottom() {
    return getMaxScrollTop() - window.scrollY <= SCROLL_DOCUMENT_EDGE_TOLERANCE_PX
}

// 根据目标图片和阅读百分比计算 scrollTop；边缘标志优先直接返回顶/底
function getTargetScrollTopForImage(image, percent, edgeFlags = {}) {
    if (edgeFlags.isAtTop) {
        return 0
    }
    if (edgeFlags.isAtBottom) {
        return getMaxScrollTop()
    }

    const rect = image.getBoundingClientRect()
    const absoluteTop = window.scrollY + rect.top
    return clamp(
        absoluteTop + rect.height * clamp(percent, 0, 1) - window.innerHeight / 2,
        0,
        getMaxScrollTop(),
    )
}

// 清除远端动画结束后的定时器
function clearRemoteScrollSettleTimer() {
    if (scrollSyncState.remoteScrollSettleTimer !== null) {
        clearTimeout(scrollSyncState.remoteScrollSettleTimer)
        scrollSyncState.remoteScrollSettleTimer = null
    }
}

// 远端跟随动画是否正在运行
function isRemoteScrollAnimationActive() {
    return scrollSyncState.remoteScrollAnimationFrameID !== null
}

// 取消远端跟随动画
function cancelRemoteScrollAnimation() {
    if (scrollSyncState.remoteScrollAnimationFrameID !== null) {
        cancelAnimationFrame(scrollSyncState.remoteScrollAnimationFrameID)
        scrollSyncState.remoteScrollAnimationFrameID = null
    }
    scrollSyncState.remoteScrollAnimationLastFrameTime = 0
}

// 远端动画停止后，延迟触发一次中心追踪更新
function scheduleRemoteScrollSettle() {
    clearRemoteScrollSettleTimer()
    scrollSyncState.remoteScrollSettleTimer = setTimeout(() => {
        scrollSyncState.remoteScrollSettleTimer = null
        scheduleCenterPageUpdate()
    }, SCROLL_REMOTE_SYNC_SETTLE_MS)
}

function getRemoteScrollFollowTimeConstant(distance, triggerSource) {
    // 时间常数越大，跟随越“慢”但更稳；连续同步时保持持续拖尾感
    if (distance < 32) {
        return triggerSource === 'auto' ? 90 : 110
    }
    if (distance < 180) {
        return triggerSource === 'auto' ? 120 : 145
    }
    if (distance < 720) {
        return triggerSource === 'auto' ? 150 : 185
    }
    if (distance < 2000) {
        return triggerSource === 'auto' ? 185 : 225
    }
    return triggerSource === 'auto' ? 220 : 260
}

// 使用指数平滑驱动远端同步滚动动画；连续目标更新只修改目标值不重启 RAF
function animateRemoteScrollTo(targetTop, triggerSource) {
    const clampedTargetTop = clamp(targetTop, 0, getMaxScrollTop())
    clearRemoteScrollSettleTimer()
    scrollSyncState.remoteScrollAnimationTargetTop = clampedTargetTop
    scrollSyncState.remoteScrollAnimationTriggerSource = triggerSource
    scrollSyncState.remoteScrollAnimationLastTargetAt = performance.now()

    if (scrollSyncState.remoteScrollAnimationFrameID !== null) {
        return
    }

    // 远端同步复用同一条动画链路；新目标到来时只更新目标，不重启动画
    const step = (timestamp) => {
        if (scrollSyncState.remoteScrollAnimationLastFrameTime === 0) {
            scrollSyncState.remoteScrollAnimationLastFrameTime = timestamp
        }

        const dt = Math.max(1, timestamp - scrollSyncState.remoteScrollAnimationLastFrameTime)
        scrollSyncState.remoteScrollAnimationLastFrameTime = timestamp

        const currentTop = window.scrollY
        const targetTopNow = clamp(
            scrollSyncState.remoteScrollAnimationTargetTop,
            0,
            getMaxScrollTop(),
        )
        const delta = targetTopNow - currentTop
        const distance = Math.abs(delta)
        const idleForMs = timestamp - scrollSyncState.remoteScrollAnimationLastTargetAt

        // 边缘目标（顶部/底部）：不等 idle，距离足够小就直接吸附
        const isEdgeTarget = targetTopNow === 0 || targetTopNow >= getMaxScrollTop()
        if (isEdgeTarget && distance <= SCROLL_REMOTE_EDGE_SNAP_PX) {
            window.scrollTo({
                top: targetTopNow,
                behavior: 'auto',
            })
            scrollSyncState.remoteScrollAnimationFrameID = null
            scrollSyncState.remoteScrollAnimationLastFrameTime = 0
            scheduleRemoteScrollSettle()
            return
        }

        if (distance <= SCROLL_REMOTE_FOLLOW_SNAP_PX && idleForMs >= SCROLL_REMOTE_FOLLOW_IDLE_MS) {
            window.scrollTo({
                top: targetTopNow,
                behavior: 'auto',
            })
            scrollSyncState.remoteScrollAnimationFrameID = null
            scrollSyncState.remoteScrollAnimationLastFrameTime = 0
            scheduleRemoteScrollSettle()
            return
        }

        const timeConstant = getRemoteScrollFollowTimeConstant(
            distance,
            scrollSyncState.remoteScrollAnimationTriggerSource,
        )
        const followAlpha = 1 - Math.exp(-dt / Math.max(1, timeConstant))
        const nextTop = currentTop + delta * followAlpha

        window.scrollTo({
            top: nextTop,
            behavior: 'auto',
        })

        scrollSyncState.remoteScrollAnimationFrameID = requestAnimationFrame(step)
    }

    scrollSyncState.remoteScrollAnimationFrameID = requestAnimationFrame(step)
}

// 判断当前是否为分页卷轴模式
function isPagedScrollMode() {
    return new URLSearchParams(window.location.search).has('page')
}

// 判断当前是否为无限滚动模式
function isInfiniteScrollMode() {
    return !isPagedScrollMode()
}

// 判断卷轴同步功能是否开启（需在线且全局开关启用）
function isScrollSyncEnabled() {
    return Alpine.store('global').onlineBook && Alpine.store('global').syncPageByWS
}

// 生成无限滚动模式下指定起始页码的 URL
function getInfiniteScrollURL(pageNum) {
    const targetURL = new URL(`/scroll/${book.id}`, window.location.origin)
    if (pageNum > 1) {
        targetURL.searchParams.set('start', pageNum.toString())
    }
    return targetURL.toString()
}

// 根据图片页码计算所属分页块号
function getPagedChunkForPageNum(pageNum) {
    return Math.floor((pageNum - 1) / PAGED_SIZE) + 1
}

// 生成指定图片页码所在分页块的 URL
function getPagedScrollURL(pageNum) {
    const chunkPage = getPagedChunkForPageNum(pageNum)
    const targetURL = new URL(`/scroll/${book.id}`, window.location.origin)
    targetURL.searchParams.set('page', chunkPage.toString())
    return targetURL.toString()
}

// 暴露给外部（如自动滚动插件）的卷轴同步控制接口
window.ComiGoScrollSync = {
    setAutoScrollEnabled(enabled) {
        scrollSyncState.autoScrollEnabled = !!enabled
    },
    scheduleCenterUpdate() {
        scheduleCenterPageUpdate()
    },
}

// 将远端同步数据暂存到 sessionStorage，供跨页导航后恢复
function savePendingRemoteSync(data) {
    try {
        sessionStorage.setItem(
            SCROLL_SYNC_PENDING_KEY,
            JSON.stringify({
                ...data,
                expire_at: Date.now() + 15000,
            }),
        )
    } catch (error) {
        console.error('保存卷轴同步定位失败:', error)
    }
}

// 从 sessionStorage 读取待恢复的远端同步数据
function loadPendingRemoteSync() {
    try {
        const raw = sessionStorage.getItem(SCROLL_SYNC_PENDING_KEY)
        if (!raw) {
            return null
        }
        const data = JSON.parse(raw)
        if (!data || data.expire_at < Date.now()) {
            sessionStorage.removeItem(SCROLL_SYNC_PENDING_KEY)
            return null
        }
        return data
    } catch (error) {
        sessionStorage.removeItem(SCROLL_SYNC_PENDING_KEY)
        console.error('读取卷轴同步定位失败:', error)
        return null
    }
}

// 清除 sessionStorage 中的待恢复同步数据
function clearPendingRemoteSync() {
    try {
        sessionStorage.removeItem(SCROLL_SYNC_PENDING_KEY)
    } catch (_) {}
}

// 获取当前 DOM 中所有带页码属性的漫画图片元素
function getScrollPageImages() {
    return Array.from(
        document.querySelectorAll('#ScrollMainArea img.manga_image[data-scroll-page-num]'),
    ).filter((image) => {
        const pageNum = parseInt(image.dataset.scrollPageNum, 10)
        return Number.isInteger(pageNum) && pageNum > 0
    })
}

// 根据页码获取对应的图片 DOM 元素
function getScrollImageByPageNum(pageNum) {
    return document.querySelector(
        `#ScrollMainArea img.manga_image[data-scroll-page-num="${pageNum}"]`,
    )
}

// 获取当前已加载图片的起止页码范围
function getLoadedPageRange() {
    const pageNums = getScrollPageImages()
        .map((image) => parseInt(image.dataset.scrollPageNum, 10))
        .filter((pageNum) => Number.isInteger(pageNum))
        .sort((a, b) => a - b)

    if (pageNums.length === 0) {
        return {
            startLoadPageNum: 1,
            endLoadPageNum: 1,
        }
    }

    return {
        startLoadPageNum: pageNums[0],
        endLoadPageNum: pageNums[pageNums.length - 1],
    }
}

// 计算屏幕中线在指定图片上的纵向百分比（0=顶部，1=底部）
function getTrackedPercent(image) {
    const rect = image.getBoundingClientRect()
    if (rect.height <= 0) {
        return 0
    }
    return clamp((window.innerHeight / 2 - rect.top) / rect.height, 0, 1)
}

// 解析当前屏幕中线所在的图片及其阅读百分比
function resolveCenterTrackedPage() {
    const centerY = window.innerHeight / 2
    const pageImages = getScrollPageImages()
    if (pageImages.length === 0) {
        if (Alpine.store('global').debugMode) {
            console.log('[scroll-sync] resolveCenterTrackedPage: 未找到带 data-scroll-page-num 的图片')
        }
        return null
    }

    let activeImage = null
    for (const image of pageImages) {
        const rect = image.getBoundingClientRect()
        if (rect.height <= 0) {
            continue
        }
        if (rect.top <= centerY && rect.bottom >= centerY) {
            activeImage = image
            break
        }
    }

    if (!activeImage) {
        let closestDistance = Infinity
        for (const image of pageImages) {
            const rect = image.getBoundingClientRect()
            if (rect.height <= 0) {
                continue
            }
            const imageCenter = (rect.top + rect.bottom) / 2
            const distance = Math.abs(imageCenter - centerY)
            if (distance < closestDistance) {
                closestDistance = distance
                activeImage = image
            }
        }
    }

    if (!activeImage) {
        return null
    }

    return {
        image: activeImage,
        pageNum: parseInt(activeImage.dataset.scrollPageNum, 10),
        percent: getTrackedPercent(activeImage),
    }
}

// 更新卷轴模式的自动书签
function updateScrollBookmark(pageNum) {
    if (!book || !book.id || !Alpine.store('global').onlineBook) {
        return
    }
    if (scrollSyncState.lastBookmarkedPageNum === pageNum) {
        return
    }

    scrollSyncState.lastBookmarkedPageNum = pageNum
    Alpine.store('global')
        .UpdateBookmark({
            type: 'auto',
            bookId: book.id,
            pageIndex: pageNum,
        })
        .catch((error) => {
            console.error('更新卷轴模式自动书签失败:', error)
        })
}

// 持久化当前追踪页码到 store 和书签
function persistTrackedPage(pageNum) {
    Alpine.store('global').nowPageNum = pageNum
    Alpine.store('global').savePageNumToLocalStorage(book.id)
    updateScrollBookmark(pageNum)
}

// 当前是否处于回传抑制期
function isSuppressingRemoteBroadcast() {
    return Date.now() < scrollSyncState.suppressBroadcastUntil
}

// 开始应用远端同步数据，设置回传抑制计时器
function beginRemoteApply() {
    scrollSyncState.suppressBroadcastUntil = Date.now() + SCROLL_REMOTE_SUPPRESS_MS
}

// 通过 WebSocket 发送当前滚动位置同步数据
function sendScrollSyncData(tracked) {
    if (typeof window.ComiGoWS === 'undefined') {
        return
    }

    const loadRange = getLoadedPageRange()
    const triggerSource = scrollSyncState.autoScrollEnabled ? 'auto' : 'manual'
    const syncData = {
        book_id: book.id,
        now_page_num: tracked.pageNum,
        now_page_num_percent: tracked.percent,
        is_at_top: isDocumentAtTop(),
        is_at_bottom: isDocumentAtBottom(),
        start_load_page_num: loadRange.startLoadPageNum,
        end_load_page_num: loadRange.endLoadPageNum,
        trigger_source: triggerSource,
    }

    if (Alpine.store('global').debugMode) {
        console.log('[scroll-sync] sendScrollSyncData:', syncData)
    }

    window.ComiGoWS.send(
        'scroll_mode_sync_page',
        syncData,
        `卷轴模式，${triggerSource === 'auto' ? '自动' : '手动'}发送同步页数`,
    )
}

// 处理中心追踪结果：更新页码、发送同步（含节流与边缘强制发送）
function applyTrackedPage(tracked) {
    if (!tracked || !Number.isInteger(tracked.pageNum) || tracked.pageNum <= 0) {
        return
    }

    const debugMode = Alpine.store('global').debugMode

    if (!scrollSyncState.hasInitializedCenterTracking) {
        scrollSyncState.hasInitializedCenterTracking = true
        scrollSyncState.lastTrackedPageNum = tracked.pageNum
        scrollSyncState.lastBookmarkedPageNum = tracked.pageNum
        scrollSyncState.lastSentPercent = tracked.percent
        Alpine.store('global').nowPageNum = tracked.pageNum
        if (debugMode) {
            console.log('[scroll-sync] applyTrackedPage: 初始化完成, page=%d, percent=%.2f', tracked.pageNum, tracked.percent)
        }
        return
    }

    const pageChanged = tracked.pageNum !== scrollSyncState.lastTrackedPageNum
    scrollSyncState.lastTrackedPageNum = tracked.pageNum
    Alpine.store('global').nowPageNum = tracked.pageNum

    if (pageChanged) {
        persistTrackedPage(tracked.pageNum)
    }

    if (isScrollSyncEnabled() && !isSuppressingRemoteBroadcast() && !isRemoteScrollAnimationActive()) {
        const now = Date.now()
        const percentDelta = Math.abs(tracked.percent - scrollSyncState.lastSentPercent)
        // 页码变化立即发送；同一页内 percent 变化超过阈值且通过节流也发送
        const shouldSendByPercent = !pageChanged
            && percentDelta >= SCROLL_SYNC_PERCENT_THRESHOLD
            && (now - scrollSyncState.lastSyncSendTime) >= SCROLL_SYNC_THROTTLE_MS

        // 边缘强制发送：到达顶部/底部时跳过 percent 阈值，仅受时间节流约束
        const atEdge = isDocumentAtTop() || isDocumentAtBottom()
        const edgeChanged = atEdge && !scrollSyncState.lastSentAtEdge
        const shouldSendByEdge = edgeChanged
            || (atEdge && (now - scrollSyncState.lastSyncSendTime) >= SCROLL_SYNC_THROTTLE_MS)

        if (pageChanged || shouldSendByPercent || shouldSendByEdge) {
            scrollSyncState.lastSentPercent = tracked.percent
            scrollSyncState.lastSyncSendTime = now
            scrollSyncState.lastSentAtEdge = atEdge
            sendScrollSyncData(tracked)
        } else if (debugMode) {
            console.log('[scroll-sync] applyTrackedPage: 未满足发送条件 (pageChanged=%s, percentDelta=%.3f, timeDelta=%dms)',
                pageChanged, percentDelta, now - scrollSyncState.lastSyncSendTime)
        }
    } else if (debugMode) {
        console.log('[scroll-sync] applyTrackedPage: 同步未启用或处于抑制期 (syncEnabled=%s, suppressing=%s)',
            isScrollSyncEnabled(), isSuppressingRemoteBroadcast())
    }
}

// 调度一次中心追踪更新（通过 RAF 去重）
function scheduleCenterPageUpdate() {
    if (scrollSyncState.rafID !== null) {
        return
    }

    scrollSyncState.rafID = requestAnimationFrame(() => {
        scrollSyncState.rafID = null
        const tracked = resolveCenterTrackedPage()
        if (!tracked) {
            return
        }
        applyTrackedPage(tracked)
    })
}

// 等待图片加载完成后执行回调（带重试）
function waitForImageReady(image, callback, retryCount = 20) {
    if (!image) {
        return
    }

    const rect = image.getBoundingClientRect()
    if ((image.complete || image.naturalHeight > 0) && rect.height > 0) {
        callback()
        return
    }

    if (retryCount <= 0) {
        callback()
        return
    }

    setTimeout(() => {
        waitForImageReady(image, callback, retryCount - 1)
    }, 150)
}

// 将页面滚动到指定图片的指定百分比位置
function scrollImageToTrackedPercent(image, percent, triggerSource, edgeFlags = {}) {
    const targetTop = getTargetScrollTopForImage(image, percent, edgeFlags)
    animateRemoteScrollTo(targetTop, triggerSource)
}

// 应用远端发来的同步数据：定位图片并滚动，或跳转页面
function applyRemoteScrollSync(data, { allowNavigation = true } = {}) {
    if (!data || data.book_id !== book.id) {
        return
    }

    const pageNum = parseInt(data.now_page_num, 10)
    if (!Number.isInteger(pageNum) || pageNum < 1) {
        return
    }

    const percent = clamp(parseFloat(data.now_page_num_percent) || 0, 0, 1)
    const triggerSource = data.trigger_source === 'auto' ? 'auto' : 'manual'
    const edgeFlags = {
        isAtTop: data.is_at_top === true,
        isAtBottom: data.is_at_bottom === true,
    }
    const targetImage = getScrollImageByPageNum(pageNum)

    beginRemoteApply()

    if (!targetImage) {
        if (allowNavigation) {
            savePendingRemoteSync({
                book_id: data.book_id,
                now_page_num: pageNum,
                now_page_num_percent: percent,
                is_at_top: edgeFlags.isAtTop,
                is_at_bottom: edgeFlags.isAtBottom,
                trigger_source: data.trigger_source === 'auto' ? 'auto' : 'manual',
            })
            // 根据当前模式选择正确的导航 URL
            window.location.href = isPagedScrollMode()
                ? getPagedScrollURL(pageNum)
                : getInfiniteScrollURL(pageNum)
        }
        return
    }

    waitForImageReady(targetImage, () => {
        scrollImageToTrackedPercent(targetImage, percent, triggerSource, edgeFlags)
    })
}

// 页面加载后恢复 sessionStorage 中暂存的远端同步数据
function restorePendingRemoteSync(retryCount = 20) {
    const pending = loadPendingRemoteSync()
    if (!pending) {
        return
    }

    if (!getScrollImageByPageNum(parseInt(pending.now_page_num, 10))) {
        if (retryCount <= 0) {
            clearPendingRemoteSync()
            return
        }
        setTimeout(() => {
            restorePendingRemoteSync(retryCount - 1)
        }, 200)
        return
    }

    clearPendingRemoteSync()
    applyRemoteScrollSync(pending, { allowNavigation: false })
}

// 初始化卷轴模式的 WebSocket 同步连接与消息处理
function initScrollModeSync() {
    if (typeof window.ComiGoWS === 'undefined') {
        return
    }

    window.ComiGoWS.init({
        pageType: 'scroll',
        getBookId: () => book?.id,
        getWsConfig: () => SCROLL_WS_CONFIG,
        isDebug: () => Alpine.store('global').debugMode,
        onConnect() {
            // WS 连接就绪后，延迟发送一次当前位置（绕过 applyTrackedPage 的阈值逻辑）
            setTimeout(() => {
                const tracked = resolveCenterTrackedPage()
                if (tracked && isScrollSyncEnabled() && !isSuppressingRemoteBroadcast()) {
                    sendScrollSyncData(tracked)
                    scrollSyncState.lastSentPercent = tracked.percent
                    scrollSyncState.lastSyncSendTime = Date.now()
                }
            }, 500)
        },
        onMessage(msg) {
            if (
                msg.type === 'scroll_mode_sync_page' &&
                msg.tab_id !== window.ComiGoWS.getTabId()
            ) {
                try {
                    const data = JSON.parse(msg.data_string || '{}')
                    if (Alpine.store('global').syncPageByWS && data.book_id === book.id) {
                        applyRemoteScrollSync(data)
                    }
                } catch (error) {
                    console.error('卷轴模式 WebSocket 同步数据解析失败:', error)
                }
            } else if (msg.type === 'heartbeat' && Alpine.store('global').debugMode) {
                console.log('收到心跳消息')
            }
        },
    })

    if (Alpine.store('global').onlineBook) {
        window.ComiGoWS.connect()
    }
}

//滚动到顶部
function scrollToTop(scrollDuration) {
    let scrollStep = -window.scrollY / (scrollDuration / 15),
        scrollInterval = setInterval(function () {
            if (window.scrollY !== 0) {
                window.scrollBy(0, scrollStep)
            } else clearInterval(scrollInterval)
        }, 15)
}

// 带缓存地获取"返回顶部"按钮 DOM 元素
let _backTopButton = null
function getBackTopButton() {
    if (!_backTopButton) {
        _backTopButton = document.getElementById('BackTopButton')
    }
    return _backTopButton
}

// Button ID为BackTopButton的元素，点击后滚动到顶部
const backTopBtn = getBackTopButton()
if (backTopBtn) {
    backTopBtn.addEventListener('click', function () {
        scrollToTop(500)
    })
}

//滚动到一定位置显示返回顶部按钮
let scrollTopSave = 0
let scrollDownFlag = false
let showBackTopFlag = false
let step = 0
function onScroll() {
    let scrollTop = document.documentElement.scrollTop || document.body.scrollTop
    scrollDownFlag = scrollTop > scrollTopSave
    //防手抖,小于一定数值状态就不变 Math.abs()会导致报错
    step = scrollTopSave - scrollTop
    scrollTopSave = scrollTop
    if (step < -10 || step > 10) {
        showBackTopFlag = scrollTop > 400 && !scrollDownFlag
        const btn = getBackTopButton()
        if (btn) {
            btn.style.display = showBackTopFlag ? 'block' : 'none'
        }
    }
    scheduleCenterPageUpdate()
}
window.addEventListener('scroll', onScroll, { passive: true })

let isLandscapeMode = true
let isPortraitMode = false
//可见区域变化的时候改变页面状态
function onResize() {
    Alpine.store('scroll').imageMaxWidth = window.innerWidth
    let clientWidth = document.documentElement.clientWidth
    let clientHeight = document.documentElement.clientHeight
    let aspectRatio = clientWidth / clientHeight
    // 为了调试的时候方便,阈值是正方形
    if (aspectRatio > 19 / 19) {
        isLandscapeMode = true
        isPortraitMode = false
    } else {
        isLandscapeMode = false
        isPortraitMode = true
    }
    scheduleCenterPageUpdate()
}
//初始化时,执行一次onResize()
onResize()
//文档视图调整大小时触发 resize 事件。 https://developer.mozilla.org/zh-CN/docs/Web/API/Window/resize_event
window.addEventListener('resize', onResize)

//鼠标是否在设置区域
function getInSetArea(e) {
    let clickX = e.x //获取鼠标的X坐标（鼠标与屏幕左侧的距离,单位为px）
    let clickY = e.y //获取鼠标的Y坐标（鼠标与屏幕顶部的距离,单位为px）
    //浏览器的视口,不包括工具栏和滚动条:
    let innerWidth = window.innerWidth
    let innerHeight = window.innerHeight
    //设置区域为正方形，边长按照宽或高里面，比较小的值决定
    const setArea = 0.15
    // innerWidth >= innerHeight 的情况下
    let MinY = innerHeight * (0.5 - setArea)
    let MaxY = innerHeight * (0.5 + setArea)
    let MinX = innerWidth * 0.5 - (MaxY - MinY) * 0.5
    let MaxX = innerWidth * 0.5 + (MaxY - MinY) * 0.5
    if (innerWidth < innerHeight) {
        MinX = innerWidth * (0.5 - setArea)
        MaxX = innerWidth * (0.5 + setArea)
        MinY = innerHeight * 0.5 - (MaxX - MinX) * 0.5
        MaxY = innerHeight * 0.5 + (MaxX - MinX) * 0.5
    }
    //在设置区域
    let inSetArea = false
    if (clickX > MinX && clickX < MaxX && clickY > MinY && clickY < MaxY) {
        inSetArea = true
    }
    return inSetArea
}

//获取鼠标位置,决定是否打开设置面板
function onMouseClick(e) {
    if (getInSetArea(e)) {
        //获取ID为 OpenSettingButton的元素，然后模拟点击
        document.getElementById('OpenSettingButton').click()
    }
}

// base64 -i SettingsOutline.svg ，然后// 把下面这行换成输出的Base64编码
const SettingsOutlineBase64 = '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'

//获取鼠标位置,决定是否打开设置面板
function onMouseMove(e) {
    if (getInSetArea(e)) {
        e.currentTarget.style.cursor = `url("data:image/png;base64,${SettingsOutlineBase64}") 12 12, pointer`
    } else {
        e.currentTarget.style.cursor = ''
    }
}
//获取ID为 ScrollMainArea 的元素
let ScrollMainArea = document.getElementById('ScrollMainArea')
if (ScrollMainArea) {
    // 鼠标移动的时候触发移动事件
    ScrollMainArea.addEventListener('mousemove', onMouseMove)
    // 点击的时候触发点击事件
    ScrollMainArea.addEventListener('click', onMouseClick)
    // 触摸的时候也触发点击事件
    ScrollMainArea.addEventListener('touchstart', onMouseClick)
}

// 键盘快捷键
/* 记录方向键当前的按压状态 */
// 1) 方向/动作当前状态
const state = { up: false, down: false, left: false, right: false, fire: false }

/* 2) 键 → 动作 的映射表
 *   - 左边写 `event.key`（大小写无关，统一用小写）。
 *   - 键盘键位表：https://developer.mozilla.org/zh-CN/docs/Web/API/UI_Events/Keyboard_event_key_values
 *   - 右边写动作名称（小写）
 *   - 这里的动作名称可以是任意字符串，建议用小写
 *   - 同一个动作可以对应多组键：方向键 + WASD + 自定义
 */
const keyMap = {
    // 方向键 ↑
    arrowup: 'up',
    // 方向键 ↓
    arrowdown: 'down',
    // 方向键 ←
    arrowleft: 'left',
    // 方向键 →
    arrowright: 'right',
    // 长得像方向键的键位当作方向键
    '<': 'left',
    '>': 'right',
    // 英语键盘上，与 < 键在一起
    ',': 'left',
    // 英语键盘上，与 > 键在一起
    '.': 'right',
    // vim键位 hjkl 当做方向键
    h: 'left',
    j: 'down',
    k: 'up',
    l: 'right',
    // 游戏当中常用的 WSAD 当做方向键
    w: 'up',
    s: 'down',
    a: 'left',
    d: 'right',
    // Home 键
    home: 'first_page',
    // End 键
    end: 'last_page',
    // PageUp 键
    pageup: 'pre_page',
    // PageDown 键
    pagedown: 'next_page',
    // 加减相关键位当作方翻页键
    '+': 'next_page',
    '-': 'pre_page',
    '=': 'next_page',
    '——': 'pre_page',
}

// 3) 通用按键处理器：down=true 表示按下，false 表示松开
function handle(e, down) {
    const k = e.key.toLowerCase()
    const act = keyMap[k]
    if (!act) return
    state[act] = down
    if (act === 'pre_page' && down) {
        toPreviousPage()
    }
    if (act === 'next_page' && down) {
        toNextPage()
    }
    if (act === 'left' && down) {
        toPreviousPage()
    }
    if (act === 'right' && down) {
        toNextPage()
    }
    if (act === 'first_page' && down && !e.repeat) {
        jumpPageNum(1)
    }
    if (act === 'last_page' && down && !e.repeat) {
        jumpPageNum(MaxPageNum)
    }
}

// 4) 事件监听
addEventListener('keydown', (e) => handle(e, true))
addEventListener('keyup', (e) => handle(e, false))

// 根据url获取当前页码 当前 url 类似 http://localhost:1234/scroll/somebookid?page=1
function getNowPageNum() {
    const urlParams = new URLSearchParams(window.location.search)
    const page = parseInt(urlParams.get('page'))
    return isNaN(page) ? 1 : page
}

// 根据当前页码设置url并刷新，如果小于最小页码（1），打印错误并返回
function toPreviousPage() {
    if (!isPagedScrollMode()) {
        return
    }
    const currentPage = getNowPageNum()
    if (currentPage <= 1) {
        console.warn(`已经是第一页了。MaxPageNum：${MaxPageNum}`)
        showToast(i18next.t('hint_first_page'), 'warning')
        return
    }
    const newPage = currentPage - 1
    const url = new URL(window.location.href)
    url.searchParams.set('page', newPage)
    window.location.href = url.toString()
}

// 根据当前页码设置url并刷新，如果大于最大页码（MaxPageNum），打印错误并返回
function toNextPage() {
    if (!isPagedScrollMode()) {
        return
    }
    const currentPage = getNowPageNum()
    if (currentPage >= MaxPageNum) {
        console.warn(`已经是最后一页了。MaxPageNum：${MaxPageNum}`)
        showToast(i18next.t('hint_last_page'), 'warning')
        return
    }
    const newPage = currentPage + 1
    const url = new URL(window.location.href)
    url.searchParams.set('page', newPage)
    window.location.href = url.toString()
}

// 根据当前页码设置url并刷新，如果小于最小页码（1）或大于最大页码（MaxPageNum），打印错误并返回
function jumpPageNum(pageNum) {
    if (!isPagedScrollMode()) {
        return
    }
    if (pageNum < 1 || pageNum > MaxPageNum) {
        console.warn(`页码超出范围，有效范围为1-${MaxPageNum}`)
        return
    }
    const url = new URL(window.location.href)
    url.searchParams.set('page', pageNum)
    window.location.href = url.toString()
}

// 卷轴模式总初始化入口
function initScrollMode() {
    initScrollModeSync()
    // 无限滚动和分页模式均启用中心追踪与 pending sync 恢复
    scheduleCenterPageUpdate()
    setTimeout(scheduleCenterPageUpdate, 120)
    window.addEventListener(
        'load',
        () => {
            scheduleCenterPageUpdate()
            restorePendingRemoteSync()
        },
        { once: true },
    )
    restorePendingRemoteSync()
}

// ============ 卷轴模式 WebSocket 同步 ============
document.addEventListener('DOMContentLoaded', initScrollMode)
