// 本地压缩包阅读页面。文件只在浏览器内存中解压，不上传到服务器。
'use strict'

const readerState = {
    // WASM 初始化只执行一次，后续打开文件复用同一个 Promise，避免重复加载 wasm_exec。
    wasmReady: null,
    // book 是 WASM 返回的轻量元数据；真实图片按页读取后转成 object URL。
    book: null,
    fileName: '',
    objectURLs: new Map(),
    // 卷轴模式的懒加载与当前页计算状态。
    observer: null,
    centerUpdateRaf: null,
    lastPageNum: 0,
    scrollTopSave: 0,
    scrollDownFlag: false,
    showBackTopFlag: false,
    backTopButton: null,
    // 翻页模式的交互状态。这里不放进 Alpine store，避免高频 touchmove 触发响应式更新。
    flip: {
        initialized: false,
        touchStartX: 0,
        touchStartY: 0,
        touchEndX: 0,
        isSwiping: false,
        isAnimating: false,
        hasMoved: false,
        horizontalLocked: false,
        currentTranslate: 0,
        startTime: 0,
        animationID: 0,
        suppressClickUntil: 0,
        wheelThrottleTimer: null,
        toolbarHideTimer: null,
    },
}

function readerText(key, fallback) {
    if (window.i18next && i18next.exists && i18next.exists(key)) {
        return i18next.t(key)
    }
    return fallback
}

function getReaderCursorStyle(name) {
    const src = window.ComiGoReaderCursorImages?.[name]
    return src ? `url("${src}") 12 12, pointer` : 'pointer'
}

function setReaderStatus(message, type = 'info') {
    const status = document.getElementById('ReaderStatus')
    if (!status) return
    status.textContent = message
    status.className = `w-full mt-4 text-sm text-center ${type === 'error' ? 'text-red-600' : 'opacity-80'}`
}

function getReaderHeaderTitle() {
    return document.getElementById('headerTitle')
}

function chooseReaderArchiveAgain() {
    const input = document.getElementById('ReaderArchiveInput')
    if (!input) return
    input.value = ''
    input.click()
}

function setReaderDefaultHeaderTitle() {
    const title = getReaderHeaderTitle()
    if (!title) return
    const defaultTitle = title.dataset.readerDefaultTitle || readerText('reader_title', 'Local Reader')
    const version = title.dataset.readerVersion || ' - Comigo'
    title.title = defaultTitle
    title.className = 'text-lg font-semibold truncate'
    title.textContent = defaultTitle
    title.onclick = null
    document.title = `${defaultTitle}${version}`
}

function setReaderHeaderTitle(fileName, book) {
    const title = getReaderHeaderTitle()
    if (!title) return
    const version = title.dataset.readerVersion || ' - Comigo'
    const pageCount = book?.PageInfos?.length || book?.page_count || 0
    const displayName = fileName || book?.title || readerText('reader_title', 'Local Reader')
    const titleText = pageCount > 0 ? `${displayName} (${pageCount})` : displayName
    title.title = `${titleText} ${readerText('reader_choose_another_file', '点击重选')}`
    title.className = 'inline-flex items-center justify-center gap-1 max-w-full min-w-0 text-lg font-semibold text-center cursor-pointer'
    title.innerHTML = ''

    const name = document.createElement('span')
    name.className = 'truncate'
    name.textContent = titleText
    title.appendChild(name)

    const button = document.createElement('button')
    button.type = 'button'
    button.className = 'shrink-0 px-1.5 py-0.5 text-xs font-normal text-blue-700/90 rounded hover:underline hover:bg-base-200'
    button.textContent = readerText('reader_choose_another_file', '点击重选')
    title.appendChild(button)

    title.onclick = (event) => {
        event.preventDefault()
        chooseReaderArchiveAgain()
    }
    document.title = `${displayName}${version}`
}

async function loadArchiveWasm() {
    if (readerState.wasmReady) {
        return readerState.wasmReady
    }
    readerState.wasmReady = (async () => {
        if (!window.Go) {
            await loadScript(comigoPath('/assets/static/wasm/wasm_exec.js'))
        }
        const go = new Go()
        const result = await instantiateArchiveWasm(go.importObject)
        go.run(result.instance)
        return window.ComiGoArchive
    })()
    return readerState.wasmReady
}

async function instantiateArchiveWasm(importObject) {
    if (window.ComiGoReaderStaticWasmBase64) {
        const wasmBytes = base64ToUint8Array(window.ComiGoReaderStaticWasmBase64)
        return WebAssembly.instantiate(wasmBytes, importObject)
    }

    const wasmResponse = await fetch(comigoPath('/assets/static/wasm/archive.wasm'))
    try {
        return await WebAssembly.instantiateStreaming(wasmResponse.clone(), importObject)
    } catch (_) {
        // 某些嵌入式环境不会给 .wasm 返回 application/wasm，退回 ArrayBuffer 加载。
        return WebAssembly.instantiate(await wasmResponse.arrayBuffer(), importObject)
    }
}

function comigoPath(path) {
    return window.ComiGoPath ? window.ComiGoPath(path) : path
}

function base64ToUint8Array(base64) {
    const binary = atob(base64)
    const bytes = new Uint8Array(binary.length)
    for (let i = 0; i < binary.length; i += 1) {
        bytes[i] = binary.charCodeAt(i)
    }
    return bytes
}

function loadScript(src) {
    return new Promise((resolve, reject) => {
        const script = document.createElement('script')
        script.src = src
        script.onload = resolve
        script.onerror = () => reject(new Error(`load script failed: ${src}`))
        document.head.appendChild(script)
    })
}

async function openReaderArchive(file) {
    if (!file) return
    // 每次打开新压缩包前都释放旧 object URL，防止大图阅读时持续占用内存。
    cleanupReaderBook()
    readerState.fileName = file.name || ''
    setReaderStatus(readerText('reader_loading_wasm', 'Loading reader core...'))
    try {
        const archive = await loadArchiveWasm()
        setReaderStatus(readerText('reader_reading_archive', 'Reading archive...'))
        const book = await archive.open(await file.arrayBuffer(), file.name, {
            sortBy: 'default',
        })
        if (!book || !Array.isArray(book.PageInfos) || book.PageInfos.length === 0) {
            throw new Error(readerText('reader_no_images_found', 'No readable images found in this archive'))
        }
        readerState.book = book
        readerState.fileName = file.name || book.title || ''
        renderReaderBook(book)
        setReaderStatus('')
        if (typeof showToast === 'function') {
            showToast(readerText('reader_archive_ready', 'Archive ready'), 'success')
        }
    } catch (error) {
        console.error('[reader] open archive failed:', error)
        setReaderStatus(String(error?.message || error), 'error')
        if (typeof showToast === 'function') {
            showToast(readerText('reader_archive_failed', 'Failed to open archive'), 'error')
        }
    }
}

function cleanupReaderBook() {
    // object URL 不会被 GC 自动回收，关闭书籍或重新打开时必须显式 revoke。
    for (const url of readerState.objectURLs.values()) {
        URL.revokeObjectURL(url)
    }
    readerState.objectURLs.clear()
    cleanupReaderView()
    readerState.book = null
    readerState.fileName = ''
    if (window.ComiGoArchive && typeof window.ComiGoArchive.close === 'function') {
        window.ComiGoArchive.close()
    }
}

function cleanupReaderView() {
    if (readerState.observer) {
        readerState.observer.disconnect()
        readerState.observer = null
    }
    const mainArea = document.getElementById('ScrollMainArea')
    if (mainArea) {
        mainArea.innerHTML = ''
    }
    readerState.lastPageNum = 0
    readerState.scrollTopSave = 0
    readerState.scrollDownFlag = false
    readerState.showBackTopFlag = false
    const flipArea = document.getElementById('ReaderFlipArea')
    if (flipArea) {
        resetReaderFlipSlider()
    }
    restoreReaderPageLayout()
    restoreReaderHeaderToolbar()
}

function renderReaderBook(book) {
    Alpine.store('global').onlineBook = false
    normalizeReaderReadMode()
    Alpine.store('global').nowPageNum = 1
    Alpine.store('global').allPageNum = book.page_count || book.PageInfos.length
    Alpine.store('scroll').allPageNum = Alpine.store('global').allPageNum

    const picker = document.getElementById('ReaderFilePicker')
    const shell = document.getElementById('ReaderShell')
    const mainArea = document.getElementById('ScrollMainArea')
    if (!mainArea || !shell) return

    if (picker) picker.classList.add('hidden')
    shell.classList.remove('hidden')
    shell.classList.add('flex')
    setReaderHeaderTitle(readerState.fileName, book)

    renderReaderCurrentMode(book)
}

function renderReaderCurrentMode(book = readerState.book) {
    if (!book) return
    // 模式切换时先清空当前视图，避免卷轴/翻页 DOM 与事件状态互相影响。
    cleanupReaderView()
    normalizeReaderReadMode()
    if (Alpine.store('global').readMode === 'page_flip') {
        renderReaderFlipBook(book)
        return
    }
    renderReaderScrollBook(book)
}

function renderReaderScrollBook(book) {
    const mainArea = document.getElementById('ScrollMainArea')
    const flipArea = document.getElementById('ReaderFlipArea')
    const flipSteps = document.getElementById('ReaderFlipStepsRangeArea')
    const backTop = getReaderBackTopButton()
    if (!mainArea) return
    restoreReaderPageLayout()
    mainArea.classList.remove('hidden')
    mainArea.classList.add('flex')
    if (flipArea) {
        flipArea.classList.add('hidden')
        flipArea.classList.remove('flex')
    }
    if (flipSteps) {
        flipSteps.classList.add('hidden')
    }
    if (backTop) {
        backTop.style.display = 'none'
    }
    restoreReaderHeaderToolbar()

    mainArea.innerHTML = ''
    const fragment = document.createDocumentFragment()
    book.PageInfos.forEach((page, index) => {
        fragment.appendChild(createReaderPageNode(page, index, book.PageInfos.length))
    })
    mainArea.appendChild(fragment)
    Alpine.initTree(mainArea)
    // 卷轴模式先生成空 img，再由 IntersectionObserver 近屏加载图片，减少大压缩包首屏等待。
    initReaderLazyLoading()
    initReaderGestures()
    restoreReaderProgress(book)
    scheduleReaderCenterUpdate()
}

function createReaderPageNode(page, index, total) {
    const wrapper = document.createElement('div')
    wrapper.className = 'flex flex-col justify-start w-full max-w-full m-0 rounded item-center'
    wrapper.setAttribute(':style', "{ marginBottom: $store.scroll.marginBottomOnScrollMode + 'px' }")

    const img = document.createElement('img')
    img.dataset.scrollPageNum = String(index + 1)
    img.dataset.readerPageIndex = String(index)
    img.alt = page.name || `Page ${index + 1}`
    img.draggable = false
    img.className = 'w-full manga_image min-h-16 text-center select-none'
    img.setAttribute('x-data', '{ isDoublePage: false, loaded: false }')
    img.setAttribute('@load', 'isDoublePage=$event.target.naturalWidth > $event.target.naturalHeight; loaded = true; window.ComiGoReader?.scheduleCenterUpdate();')
    img.setAttribute('@error', 'loaded = false')
    img.setAttribute(
        ':style',
        "{ width: $store.global.isLandscape?($store.scroll.widthUseFixedValue? (isDoublePage ? $store.scroll.doublePageWidth_PX +'px': $store.scroll.singlePageWidth_PX +'px'): (isDoublePage ? $store.scroll.doublePageWidth_Percent + '%':$store.scroll.singlePageWidth_Percent + '%')): $store.scroll.portraitWidthPercent+'%', maxWidth: '100%'}",
    )
    wrapper.appendChild(img)

    const pageHint = document.createElement('div')
    pageHint.className = 'w-full mt-0 mb-1 text-sm font-semibold text-center page_hint'
    pageHint.textContent = `${index + 1} / ${total}`
    pageHint.setAttribute('x-show', '$store.scroll.showPageNum')
    wrapper.appendChild(pageHint)
    return wrapper
}

function initReaderLazyLoading() {
    const images = document.querySelectorAll('#ScrollMainArea img.manga_image[data-reader-page-index]')
    readerState.observer = new IntersectionObserver((entries) => {
        for (const entry of entries) {
            if (!entry.isIntersecting) continue
            const img = entry.target
            loadReaderImage(img)
            readerState.observer.unobserve(img)
        }
    }, {
        root: null,
        rootMargin: '1200px 0px',
        threshold: 0.01,
    })
    images.forEach((img) => readerState.observer.observe(img))
}

async function loadReaderImage(img) {
    const index = parseInt(img.dataset.readerPageIndex, 10)
    if (!Number.isInteger(index) || img.src) return
    img.src = await getReaderPageObjectURL(index)
}

async function getReaderPageObjectURL(index) {
    if (!Number.isInteger(index) || index < 0 || index >= (readerState.book?.PageInfos?.length || 0)) {
        return ''
    }
    if (readerState.objectURLs.has(index)) {
        return readerState.objectURLs.get(index)
    }
    try {
        // 图片数据只从 WASM 按需读取一次，之后通过 object URL 复用。
        const bytes = await window.ComiGoArchive.readPage(index)
        const page = readerState.book.PageInfos[index]
        const blob = new Blob([bytes], { type: guessMimeType(page?.name || '') })
        const objectURL = URL.createObjectURL(blob)
        readerState.objectURLs.set(index, objectURL)
        return objectURL
    } catch (error) {
        console.error('[reader] load image failed:', error)
        return ''
    }
}

function guessMimeType(filename) {
    const lower = filename.toLowerCase()
    if (lower.endsWith('.jpg') || lower.endsWith('.jpeg')) return 'image/jpeg'
    if (lower.endsWith('.png')) return 'image/png'
    if (lower.endsWith('.gif')) return 'image/gif'
    if (lower.endsWith('.webp')) return 'image/webp'
    if (lower.endsWith('.svg')) return 'image/svg+xml'
    if (lower.endsWith('.bmp')) return 'image/bmp'
    if (lower.endsWith('.avif')) return 'image/avif'
    if (lower.endsWith('.html') || lower.endsWith('.htm')) return 'text/html'
    return 'application/octet-stream'
}

function getReaderPageImages() {
    return Array.from(document.querySelectorAll('#ScrollMainArea img.manga_image[data-scroll-page-num]'))
}

function resolveReaderCenterPage() {
    const centerY = window.innerHeight / 2
    let active = null
    let closestDistance = Infinity
    for (const image of getReaderPageImages()) {
        const rect = image.getBoundingClientRect()
        if (rect.height <= 0) continue
        if (rect.top <= centerY && rect.bottom >= centerY) {
            active = image
            break
        }
        const distance = Math.abs((rect.top + rect.bottom) / 2 - centerY)
        if (distance < closestDistance) {
            closestDistance = distance
            active = image
        }
    }
    if (!active) return 1
    return parseInt(active.dataset.scrollPageNum, 10) || 1
}

function scheduleReaderCenterUpdate() {
    if (readerState.centerUpdateRaf !== null) return
    readerState.centerUpdateRaf = requestAnimationFrame(() => {
        readerState.centerUpdateRaf = null
        const pageNum = resolveReaderCenterPage()
        if (pageNum === readerState.lastPageNum) return
        readerState.lastPageNum = pageNum
        Alpine.store('global').nowPageNum = pageNum
        saveReaderProgress(pageNum)
    })
}

// 暴露给抽屉插件和内联 Alpine 表达式使用，保持与 scroll.js 的全局入口类似。
window.ComiGoReader = {
    scheduleCenterUpdate: scheduleReaderCenterUpdate,
    setReadMode: setReaderReadMode,
    refreshFlip: updateReaderFlipImages,
    inputFlipPage: inputReaderFlipPage,
    showFlipToolbar: syncReaderFlipToolbar,
}

function saveReaderProgress(pageNum) {
    if (!readerState.book?.id || !Alpine.store('global').saveReadingProgress) return
    localStorage.setItem(`reader.pageNum.${readerState.book.id}`, String(pageNum))
}

function restoreReaderProgress(book) {
    if (!book?.id || !Alpine.store('global').saveReadingProgress) return
    const raw = localStorage.getItem(`reader.pageNum.${book.id}`)
    const pageNum = parseInt(raw, 10)
    if (!Number.isInteger(pageNum) || pageNum <= 1 || pageNum > book.PageInfos.length) return
    const target = document.querySelector(`#ScrollMainArea img[data-scroll-page-num="${pageNum}"]`)
    if (target) {
        requestAnimationFrame(() => {
            target.scrollIntoView({ block: 'center' })
            scheduleReaderCenterUpdate()
        })
    }
}

function normalizeReaderReadMode() {
    const mode = Alpine.store('global').readMode
    if (mode === 'flip_page') {
        Alpine.store('global').readMode = 'page_flip'
        return
    }
    if (mode !== 'page_flip' && mode !== 'infinite_scroll') {
        Alpine.store('global').readMode = 'infinite_scroll'
    }
}

function setReaderReadMode(mode) {
    // reader 只提供本地无限卷轴和本地翻页，旧的 flip_page 值统一兼容为 page_flip。
    Alpine.store('global').readMode = mode === 'page_flip' || mode === 'flip_page' ? 'page_flip' : 'infinite_scroll'
    if (readerState.book) {
        renderReaderCurrentMode(readerState.book)
    } else if (Alpine.store('global').readMode !== 'page_flip') {
        restoreReaderPageLayout()
        restoreReaderHeaderToolbar()
    }
}

function getReaderStoredPage(book) {
    if (!book?.id || !Alpine.store('global').saveReadingProgress) return 1
    const raw = localStorage.getItem(`reader.pageNum.${book.id}`)
    const pageNum = parseInt(raw, 10)
    if (!Number.isInteger(pageNum) || pageNum < 1 || pageNum > book.PageInfos.length) return 1
    return pageNum
}

function renderReaderFlipBook(book) {
    const mainArea = document.getElementById('ScrollMainArea')
    const flipArea = document.getElementById('ReaderFlipArea')
    const flipSteps = document.getElementById('ReaderFlipStepsRangeArea')
    const backTop = getReaderBackTopButton()
    if (!flipArea) return
    applyReaderFlipLayout()

    if (mainArea) {
        mainArea.classList.add('hidden')
        mainArea.classList.remove('flex')
    }
    flipArea.classList.remove('hidden')
    flipArea.classList.add('flex')
    if (flipSteps) {
        flipSteps.classList.remove('hidden')
        flipSteps.classList.add('flex', 'flex-col', 'justify-center')
    }
    if (backTop) {
        backTop.style.display = 'none'
    }

    // 翻页模式不保留卷轴滚动位置，恢复的是最后阅读页码。
    Alpine.store('global').nowPageNum = getReaderStoredPage(book)
    initReaderFlipListeners()
    syncReaderFlipToolbar()
    updateReaderFlipImages()
    updateReaderFlipProgress()
}

function getReaderFlipElements() {
    return {
        area: document.getElementById('ReaderFlipArea'),
        steps: document.getElementById('ReaderFlipStepsRangeArea'),
        range: document.getElementById('ReaderFlipStepsRange'),
        sliderContainer: document.getElementById('reader-flip-slider-container'),
        slider: document.getElementById('reader-flip-slider'),
        leftSlide: document.getElementById('reader-flip-left-slide'),
        rightSlide: document.getElementById('reader-flip-right-slide'),
        singleImage: document.getElementById('ReaderFlipSingleImage'),
        doubleLeft: document.getElementById('ReaderFlipDoubleImageLeft'),
        doubleRight: document.getElementById('ReaderFlipDoubleImageRight'),
    }
}

function getReaderFlipPaginationUtils() {
    return window.ComiGoFlip?.pagination
}

function getReaderFlipInteractionUtils() {
    return window.ComiGoFlip?.interaction
}

function createReaderFlipImage(className) {
    const img = document.createElement('img')
    img.className = className
    img.draggable = false
    return img
}

async function setImageElementSrc(img, index) {
    if (!img) return
    if (index < 0 || index >= readerState.book.PageInfos.length) {
        img.removeAttribute('src')
        return
    }
    const url = await getReaderPageObjectURL(index)
    if (url) {
        img.src = url
        img.alt = readerState.book.PageInfos[index]?.name || `Page ${index + 1}`
    }
}

function getReaderFlipStepNext() {
    const util = getReaderFlipPaginationUtils()
    const nowPageNum = parseInt(Alpine.store('global').nowPageNum, 10)
    const allPageNum = Alpine.store('global').allPageNum
    const doublePageMode = Alpine.store('flip').doublePageMode === true
    return util?.getNextPageStep ? util.getNextPageStep(doublePageMode, nowPageNum, allPageNum) : 1
}

function getReaderFlipStepPrevious() {
    const util = getReaderFlipPaginationUtils()
    const nowPageNum = parseInt(Alpine.store('global').nowPageNum, 10)
    const doublePageMode = Alpine.store('flip').doublePageMode === true
    return util?.getPreviousPageStep ? util.getPreviousPageStep(doublePageMode, nowPageNum) : -1
}

async function updateReaderFlipImages() {
    if (!readerState.book || Alpine.store('global').readMode !== 'page_flip') return
    const elements = getReaderFlipElements()
    const nowPageNum = parseInt(Alpine.store('global').nowPageNum, 10)
    const allPageNum = Alpine.store('global').allPageNum
    const mangaMode = Alpine.store('flip').mangaMode
    const doublePageMode = Alpine.store('flip').doublePageMode
    const isPortrait = Alpine.store('global').isPortrait

    // 漫画模式是从右往左阅读，因此相邻页容器的左右方向需要反过来。
    if (mangaMode) {
        elements.leftSlide.style.transform = 'translateX(100%)'
        elements.rightSlide.style.transform = 'translateX(-100%)'
    } else {
        elements.leftSlide.style.transform = 'translateX(-100%)'
        elements.rightSlide.style.transform = 'translateX(100%)'
    }

    const singlePageClass = isPortrait
        ? 'object-contain w-auto max-w-full h-screen'
        : 'h-screen w-auto max-w-full object-contain'
    const doublePageClass = 'object-contain w-auto max-h-screen m-0 select-none max-w-1/2 grow-0'
    const singleImgClass = 'object-contain h-screen max-w-full max-h-screen m-0'
    const adjacentLoads = []

    // 中间页是当前可见页；左右页只用于拖动/滑出动画，必须提前准备好 src。
    if (!doublePageMode) {
        await setImageElementSrc(elements.singleImage, nowPageNum - 1)
        elements.leftSlide.innerHTML = ''
        if (nowPageNum > 1) {
            const prev = createReaderFlipImage(singlePageClass)
            elements.leftSlide.appendChild(prev)
            adjacentLoads.push(setImageElementSrc(prev, nowPageNum - 2))
        }
        elements.rightSlide.innerHTML = ''
        if (nowPageNum < allPageNum) {
            const next = createReaderFlipImage(singlePageClass)
            elements.rightSlide.appendChild(next)
            adjacentLoads.push(setImageElementSrc(next, nowPageNum))
        }
    } else {
        const leftIndex = mangaMode ? nowPageNum : nowPageNum - 1
        const rightIndex = mangaMode ? nowPageNum - 1 : nowPageNum
        await setImageElementSrc(elements.doubleLeft, leftIndex)
        await setImageElementSrc(elements.doubleRight, rightIndex)
        elements.leftSlide.innerHTML = ''
        elements.rightSlide.innerHTML = ''
        const prevStart = Math.max(0, nowPageNum + getReaderFlipStepPrevious() - 1)
        if (nowPageNum > 1) {
            const useSinglePrev = nowPageNum === 2
            const prevA = createReaderFlipImage(useSinglePrev ? singleImgClass : doublePageClass)
            elements.leftSlide.appendChild(prevA)
            adjacentLoads.push(setImageElementSrc(prevA, useSinglePrev ? nowPageNum - 2 : (mangaMode ? prevStart + 1 : prevStart)))
            if (!useSinglePrev) {
                const prevB = createReaderFlipImage(doublePageClass)
                elements.leftSlide.appendChild(prevB)
                adjacentLoads.push(setImageElementSrc(prevB, mangaMode ? prevStart : prevStart + 1))
            }
        }
        const nextStart = nowPageNum + getReaderFlipStepNext() - 1
        if (nextStart < allPageNum) {
            const useSingleNext = nextStart === allPageNum - 1
            const nextA = createReaderFlipImage(useSingleNext ? singleImgClass : doublePageClass)
            elements.rightSlide.appendChild(nextA)
            adjacentLoads.push(setImageElementSrc(nextA, useSingleNext ? nextStart : (mangaMode ? nextStart + 1 : nextStart)))
            if (!useSingleNext) {
                const nextB = createReaderFlipImage(doublePageClass)
                elements.rightSlide.appendChild(nextB)
                adjacentLoads.push(setImageElementSrc(nextB, mangaMode ? nextStart : nextStart + 1))
            }
        }
    }

    // 等待相邻页加载完成再复位 slider，避免快速滑动时动画末端闪白。
    await Promise.all(adjacentLoads)
    resetReaderFlipSlider()
    updateReaderFlipProgress()
}

function updateReaderFlipProgress() {
    const range = document.getElementById('ReaderFlipStepsRange')
    if (!range) return
    const allPageNum = Alpine.store('global').allPageNum || 1
    const nowPageNum = Alpine.store('global').nowPageNum || 1
    const value = Alpine.store('flip').mangaMode ? allPageNum - nowPageNum + 1 : nowPageNum
    range.value = value
    const percent = allPageNum <= 1 ? 0 : ((value - 1) / (allPageNum - 1)) * 100
    range.style.setProperty('--value-percent', `${percent}%`)
}

function inputReaderFlipPage(event) {
    const inputValue = parseInt(event.target.value, 10)
    const target = Alpine.store('flip').mangaMode
        ? Alpine.store('global').allPageNum - inputValue + 1
        : inputValue
    jumpReaderFlipPage(target)
}

function jumpReaderFlipPage(pageNum) {
    const allPageNum = Alpine.store('global').allPageNum
    const target = Math.min(Math.max(parseInt(pageNum, 10) || 1, 1), allPageNum)
    Alpine.store('global').nowPageNum = target
    saveReaderProgress(target)
    updateReaderFlipImages()
}

function addReaderFlipPage(step) {
    const nowPageNum = parseInt(Alpine.store('global').nowPageNum, 10)
    const allPageNum = Alpine.store('global').allPageNum
    const target = nowPageNum + step
    if (target > allPageNum) {
        if (typeof showToast === 'function') showToast(i18next.t('hint_last_page'), 'warning')
        return
    }
    if (target < 1) {
        if (typeof showToast === 'function') showToast(i18next.t('hint_first_page'), 'warning')
        return
    }
    jumpReaderFlipPage(target)
}

function toNextReaderFlipPage() {
    const step = getReaderFlipStepNext()
    if (step !== 0) addReaderFlipPage(step)
}

function toPreviousReaderFlipPage() {
    const step = getReaderFlipStepPrevious()
    if (step !== 0) addReaderFlipPage(step)
}

function resetReaderFlipSlider() {
    const elements = getReaderFlipElements()
    cancelAnimationFrame(readerState.flip.animationID)
    readerState.flip.isAnimating = false
    readerState.flip.isSwiping = false
    readerState.flip.hasMoved = false
    readerState.flip.horizontalLocked = false
    readerState.flip.currentTranslate = 0
    if (elements.slider) {
        // 复位必须跳过 Tailwind transition，否则下一次手势开始会出现残留动画。
        elements.slider.style.transition = 'none'
        elements.slider.style.transform = 'translateX(0)'
        elements.slider.offsetHeight
        elements.slider.style.transition = ''
    }
}

function getReaderPointerPoint(event) {
    const touch = event.touches?.[0] || event.changedTouches?.[0]
    const pointer = touch || event
    return {
        x: pointer?.clientX || 0,
        y: pointer?.clientY || 0,
    }
}

function setReaderFlipTranslate(value) {
    const elements = getReaderFlipElements()
    if (!elements.slider) return
    readerState.flip.currentTranslate = value
    elements.slider.style.transition = 'none'
    elements.slider.style.transform = `translate3d(${value}px, 0, 0)`
}

// 边界拖动使用阻尼，既给出反馈，也避免第一页/末页被硬推到空白区域。
function getReaderRubberBandTranslate(diffX) {
    const limit = Math.min(window.innerWidth * 0.18, 72)
    const distance = Math.abs(diffX)
    const eased = limit * (1 - Math.exp(-distance / limit))
    return Math.sign(diffX) * eased
}

function easeOutReaderFlip(x) {
    return 1 - Math.pow(1 - x, 3)
}

function animateReaderFlipTo(targetPosition, duration, onComplete) {
    const elements = getReaderFlipElements()
    if (!elements.slider) return
    cancelAnimationFrame(readerState.flip.animationID)
    readerState.flip.isAnimating = true
    const startPosition = readerState.flip.currentTranslate
    let startedAt = null

    function animate(timestamp) {
        if (startedAt === null) startedAt = timestamp
        const progress = Math.min((timestamp - startedAt) / duration, 1)
        const eased = easeOutReaderFlip(progress)
        setReaderFlipTranslate(startPosition + (targetPosition - startPosition) * eased)
        if (progress < 1) {
            readerState.flip.animationID = requestAnimationFrame(animate)
            return
        }
        readerState.flip.isAnimating = false
        readerState.flip.currentTranslate = targetPosition
        if (typeof onComplete === 'function') onComplete()
    }

    readerState.flip.animationID = requestAnimationFrame(animate)
}

function animateReaderFlipReset() {
    // 滑动距离不足或触碰边界时回弹到当前页。
    const duration = Alpine.store('flip').resetAnimationDuration || 320
    animateReaderFlipTo(0, duration, () => {
        readerState.flip.currentTranslate = 0
        readerState.flip.hasMoved = false
    })
}

function animateReaderFlipSlide(direction) {
    // 松手后先把当前 slider 补完到整屏宽度，再提交页码切换。
    const width = getReaderFlipElements().sliderContainer?.clientWidth || window.innerWidth
    const targetPosition = direction === 'left' ? -width : width
    const duration = Alpine.store('flip').swipeAnimationDuration || 260
    animateReaderFlipTo(targetPosition, duration, () => {
        const mangaMode = Alpine.store('flip').mangaMode
        const turnForward = (direction === 'left' && !mangaMode) || (direction === 'right' && mangaMode)
        if (turnForward) {
            toNextReaderFlipPage()
        } else {
            toPreviousReaderFlipPage()
        }
        readerState.flip.suppressClickUntil = Date.now() + 350
    })
}

function readerFlipShouldBlock(diffX) {
    const util = getReaderFlipPaginationUtils()
    if (util?.shouldBlockScrollBoundary) {
        return util.shouldBlockScrollBoundary(
            diffX,
            Alpine.store('flip').mangaMode,
            Alpine.store('global').nowPageNum,
            Alpine.store('global').allPageNum,
        )
    }
    return false
}

function readerFlipTouchStart(event) {
    if (Alpine.store('global').readMode !== 'page_flip' || !Alpine.store('flip').swipeTurn || readerState.flip.isAnimating) return
    const point = getReaderPointerPoint(event)
    readerState.flip.startTime = Date.now()
    readerState.flip.isSwiping = true
    readerState.flip.hasMoved = false
    readerState.flip.horizontalLocked = false
    readerState.flip.touchStartX = point.x
    readerState.flip.touchStartY = point.y
    readerState.flip.currentTranslate = 0
    cancelAnimationFrame(readerState.flip.animationID)
    const elements = getReaderFlipElements()
    if (elements.slider) {
        // 拖动期间 transform 必须完全跟手，不能使用 CSS transition。
        elements.slider.style.transition = 'none'
    }
}

function readerFlipTouchMove(event) {
    const elements = getReaderFlipElements()
    if (!readerState.flip.isSwiping || !Alpine.store('flip').swipeTurn || !elements.slider || readerState.flip.isAnimating) return
    const point = getReaderPointerPoint(event)
    const diffX = point.x - readerState.flip.touchStartX
    const diffY = point.y - readerState.flip.touchStartY

    if (!readerState.flip.horizontalLocked) {
        if (Math.abs(diffX) < 8 && Math.abs(diffY) < 8) return
        if (Math.abs(diffY) > Math.abs(diffX) * 1.15) {
            // 明显是纵向手势时交还给页面滚动，避免误触发翻页。
            readerState.flip.isSwiping = false
            return
        }
        readerState.flip.horizontalLocked = true
    }

    readerState.flip.hasMoved = true
    const translate = readerFlipShouldBlock(diffX) ? getReaderRubberBandTranslate(diffX) : diffX
    setReaderFlipTranslate(translate)
    if (event.cancelable) event.preventDefault()
}

function readerFlipTouchEnd(event) {
    if (!readerState.flip.isSwiping || !Alpine.store('flip').swipeTurn || readerState.flip.isAnimating) return
    readerState.flip.isSwiping = false
    const point = getReaderPointerPoint(event)
    const diffX = point.x - readerState.flip.touchStartX
    if (!readerState.flip.hasMoved || !readerState.flip.horizontalLocked) {
        readerState.flip.currentTranslate = 0
        return
    }
    const quick = Date.now() - readerState.flip.startTime < (Alpine.store('flip').swipeTimeout || 300) && Math.abs(diffX) > 44
    const threshold = Alpine.store('flip').swipeThreshold || 100
    let direction = null
    if (diffX < -threshold || (quick && diffX < 0)) direction = 'left'
    if (diffX > threshold || (quick && diffX > 0)) direction = 'right'
    if (readerFlipShouldBlock(diffX) || !direction) {
        animateReaderFlipReset()
        if (readerState.flip.hasMoved) readerState.flip.suppressClickUntil = Date.now() + 300
        return
    }
    // 触摸滑动结束后浏览器可能补发 click，短时间抑制以免二次翻页。
    readerState.flip.suppressClickUntil = Date.now() + 600
    animateReaderFlipSlide(direction)
}

function onReaderFlipClick(event) {
    if (Date.now() < readerState.flip.suppressClickUntil || readerState.flip.isSwiping || readerState.flip.isAnimating || Math.abs(readerState.flip.currentTranslate) > 10) return
    if (getInReaderSettingArea(event)) {
        if (Alpine.store('flip').autoAlign) {
            document.getElementById('reader-flip-slider-container')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
        }
        showReaderFlipToolbar()
        openReaderSettings()
        return
    }
    const leftSide = event.clientX < window.innerWidth * 0.5
    if ((leftSide && !Alpine.store('flip').mangaMode) || (!leftSide && Alpine.store('flip').mangaMode)) {
        toPreviousReaderFlipPage()
    } else {
        toNextReaderFlipPage()
    }
}

function onReaderFlipWheel(event) {
    if (!Alpine.store('flip').wheelFlip || Alpine.store('global').readMode !== 'page_flip') return
    if (event.deltaY === 0 || readerState.flip.wheelThrottleTimer || readerState.flip.isAnimating) return
    event.preventDefault()
    readerState.flip.wheelThrottleTimer = setTimeout(() => {
        readerState.flip.wheelThrottleTimer = null
    }, Alpine.store('flip').wheelThrottleDelay || 250)
    if (event.deltaY > 0) {
        toNextReaderFlipPage()
    } else {
        toPreviousReaderFlipPage()
    }
}

function onReaderFlipMouseMove(event) {
    if (Alpine.store('global').readMode !== 'page_flip') return
    const elements = getReaderFlipElements()
    const x = event.clientX
    const y = event.clientY
    const inSetArea = getInReaderSettingArea(event)
    if (inSetArea) {
        event.currentTarget.style.cursor = getReaderCursorStyle('settings')
        showReaderFlipToolbar()
        return
    }

    const stepsRect = elements.steps?.getBoundingClientRect()
    if (readerFlipPointInRect(stepsRect, x, y)) {
        event.currentTarget.style.cursor = 'default'
        return
    }

    const leftSide = x < window.innerWidth * 0.5
    const mangaMode = Alpine.store('flip').mangaMode
    const nowPageNum = Alpine.store('global').nowPageNum
    const allPageNum = Alpine.store('global').allPageNum
    if (leftSide) {
        if ((!mangaMode && nowPageNum === 1) || (mangaMode && nowPageNum === allPageNum)) {
            event.currentTarget.style.cursor = getReaderCursorStyle('prohibited')
            return
        }
        event.currentTarget.style.cursor = getReaderCursorStyle('left')
        return
    }
    if ((!mangaMode && nowPageNum === allPageNum) || (mangaMode && nowPageNum === 1)) {
        event.currentTarget.style.cursor = getReaderCursorStyle('prohibited')
        return
    }
    event.currentTarget.style.cursor = getReaderCursorStyle('right')
}

function onReaderFlipMouseLeave(event) {
    event.currentTarget.style.cursor = ''
}

function initReaderFlipListeners() {
    const elements = getReaderFlipElements()
    if (!elements.sliderContainer || readerState.flip.initialized) return
    readerState.flip.initialized = true
    // 监听器只注册一次；模式切换时复用同一套 DOM，避免重复绑定导致一次手势翻多页。
    elements.sliderContainer.addEventListener('touchstart', readerFlipTouchStart)
    elements.sliderContainer.addEventListener('touchmove', readerFlipTouchMove, { passive: false })
    elements.sliderContainer.addEventListener('touchend', readerFlipTouchEnd)
    elements.sliderContainer.addEventListener('touchcancel', readerFlipTouchEnd)
    elements.sliderContainer.addEventListener('mousedown', readerFlipTouchStart)
    elements.sliderContainer.addEventListener('mousemove', readerFlipTouchMove)
    elements.sliderContainer.addEventListener('mouseup', readerFlipTouchEnd)
    elements.sliderContainer.addEventListener('mouseleave', readerFlipTouchEnd)
    elements.area.addEventListener('click', onReaderFlipClick)
    elements.area.addEventListener('wheel', onReaderFlipWheel, { passive: false })
    elements.area.addEventListener('mousemove', onReaderFlipMouseMove)
    elements.area.addEventListener('mouseleave', onReaderFlipMouseLeave)
    document.addEventListener('mousemove', onReaderFlipDocumentMouseMove)
}

function getReaderFlipToolbarElements() {
    return {
        header: document.getElementById('header'),
        steps: document.getElementById('ReaderFlipStepsRangeArea'),
    }
}

function getReaderLayoutElements() {
    return {
        header: document.getElementById('header'),
        footer: document.getElementById('ReaderFooter') || document.querySelector('footer.footer'),
        root: document.getElementById('ReaderRoot'),
        shell: document.getElementById('ReaderShell'),
        flipArea: document.getElementById('ReaderFlipArea'),
    }
}

function isReaderFlipModeActive() {
    return Alpine.store('global').readMode === 'page_flip'
}

function clearReaderFlipToolbarTimer() {
    if (readerState.flip.toolbarHideTimer) {
        clearTimeout(readerState.flip.toolbarHideTimer)
        readerState.flip.toolbarHideTimer = null
    }
}

function restoreReaderHeaderToolbar() {
    clearReaderFlipToolbarTimer()
    const { header, steps } = getReaderFlipToolbarElements()
    if (header) {
        header.style.opacity = '1'
        header.style.transform = 'translateY(0)'
    }
    if (steps) {
        steps.style.opacity = '1'
        steps.style.transform = 'translateY(0)'
    }
}

function applyReaderFlipLayout() {
    const { header, footer, root, shell, flipArea } = getReaderLayoutElements()
    // 翻页模式需要占满视口；底部 footer 会挤压图片，因此进入翻页时隐藏。
    if (footer) {
        footer.style.setProperty('display', 'none', 'important')
    }
    if (root) {
        root.style.minHeight = '100vh'
    }
    if (shell) {
        shell.style.minHeight = '100vh'
    }
    if (flipArea) {
        flipArea.style.height = Alpine.store('flip').autoHideToolbar ? '100vh' : 'calc(100vh - 3rem)'
        flipArea.style.minHeight = Alpine.store('flip').autoHideToolbar ? '100vh' : 'calc(100vh - 3rem)'
    }
    if (!header) return
    if (Alpine.store('flip').autoHideToolbar) {
        header.style.position = 'fixed'
        header.style.top = '0'
        header.style.left = '0'
        header.style.right = '0'
        header.style.width = '100%'
        header.style.zIndex = '30'
        header.style.backdropFilter = 'blur(16px)'
        return
    }
    header.style.position = ''
    header.style.top = ''
    header.style.left = ''
    header.style.right = ''
    header.style.width = ''
    header.style.zIndex = ''
    header.style.backdropFilter = ''
}

function restoreReaderPageLayout() {
    const { header, footer, root, shell, flipArea } = getReaderLayoutElements()
    // 退出翻页模式时撤销所有内联样式，避免影响卷轴模式和普通页面布局。
    if (footer) {
        footer.style.removeProperty('display')
    }
    if (root) {
        root.style.minHeight = ''
    }
    if (shell) {
        shell.style.minHeight = ''
    }
    if (flipArea) {
        flipArea.style.height = ''
        flipArea.style.minHeight = ''
    }
    if (header) {
        header.style.position = ''
        header.style.top = ''
        header.style.left = ''
        header.style.right = ''
        header.style.width = ''
        header.style.zIndex = ''
        header.style.backdropFilter = ''
    }
}

function showReaderFlipToolbar() {
    const { header, steps } = getReaderFlipToolbarElements()
    if (!header || !steps) return
    clearReaderFlipToolbarTimer()
    if (!isReaderFlipModeActive()) {
        restoreReaderHeaderToolbar()
        return
    }
    header.style.opacity = Alpine.store('flip').autoHideToolbar ? '0.9' : '1'
    steps.style.opacity = Alpine.store('flip').autoHideToolbar ? '0.9' : '1'
    header.style.transform = 'translateY(0)'
    steps.style.transform = 'translateY(0)'
}

function hideReaderFlipToolbar() {
    const { header, steps } = getReaderFlipToolbarElements()
    if (!header || !steps || !isReaderFlipModeActive() || !Alpine.store('flip').autoHideToolbar) return
    header.style.opacity = '0'
    steps.style.opacity = '0'
    header.style.transform = 'translateY(-100%)'
    steps.style.transform = 'translateY(100%)'
}

function syncReaderFlipToolbar() {
    clearReaderFlipToolbarTimer()
    if (!isReaderFlipModeActive()) {
        restoreReaderPageLayout()
        restoreReaderHeaderToolbar()
        return
    }
    applyReaderFlipLayout()
    showReaderFlipToolbar()
    if (Alpine.store('flip').autoHideToolbar) {
        // 自动隐藏只作用于翻页模式，卷轴模式继续使用普通 header/footer。
        readerState.flip.toolbarHideTimer = setTimeout(hideReaderFlipToolbar, 1000)
    }
}

function readerFlipPointInRect(rect, x, y) {
    const util = getReaderFlipInteractionUtils()
    return util?.isPointInRect ? util.isPointInRect(rect, x, y) : Boolean(rect && x >= rect.left && x <= rect.right && y >= rect.top && y <= rect.bottom)
}

function onReaderFlipDocumentMouseMove(event) {
    if (!isReaderFlipModeActive()) return
    const { header, steps } = getReaderFlipToolbarElements()
    const x = event.clientX
    const y = event.clientY
    const nearHeader = Alpine.store('flip').autoHideToolbar ? y <= 80 : readerFlipPointInRect(header?.getBoundingClientRect(), x, y)
    const nearSteps = Alpine.store('flip').autoHideToolbar ? y >= window.innerHeight - 80 : readerFlipPointInRect(steps?.getBoundingClientRect(), x, y)
    if (getInReaderSettingArea(event) || nearHeader || nearSteps) {
        showReaderFlipToolbar()
    } else {
        hideReaderFlipToolbar()
    }
}

// 判断点击位置是否位于屏幕中央设置区域，逻辑与 scroll 阅读页保持一致。
function getInReaderSettingArea(event) {
    const pointer = event.touches ? event.touches[0] : event
    if (!pointer) return false

    const clickX = pointer.clientX
    const clickY = pointer.clientY
    const innerWidth = window.innerWidth
    const innerHeight = window.innerHeight
    const setArea = 0.15

    let minY = innerHeight * (0.5 - setArea)
    let maxY = innerHeight * (0.5 + setArea)
    let minX = innerWidth * 0.5 - (maxY - minY) * 0.5
    let maxX = innerWidth * 0.5 + (maxY - minY) * 0.5
    if (innerWidth < innerHeight) {
        minX = innerWidth * (0.5 - setArea)
        maxX = innerWidth * (0.5 + setArea)
        minY = innerHeight * 0.5 - (maxX - minX) * 0.5
        maxY = innerHeight * 0.5 + (maxX - minX) * 0.5
    }
    return clickX > minX && clickX < maxX && clickY > minY && clickY < maxY
}

function openReaderSettings() {
    const button = document.getElementById('OpenSettingButton')
    if (button) {
        button.click()
    }
}

function onReaderClick(event) {
    if (getInReaderSettingArea(event)) {
        openReaderSettings()
    }
}

function onReaderMouseMove(event) {
    event.currentTarget.style.cursor = getInReaderSettingArea(event) ? getReaderCursorStyle('settings') : ''
}

function initReaderGestures() {
    const mainArea = document.getElementById('ScrollMainArea')
    if (!mainArea || mainArea.dataset.readerGesturesReady === 'true') return

    mainArea.dataset.readerGesturesReady = 'true'
    // 卷轴模式保留中间区域打开设置的交互，和 scroll 阅读页保持一致。
    mainArea.addEventListener('mousemove', onReaderMouseMove)
    mainArea.addEventListener('click', onReaderClick)
    mainArea.addEventListener('touchstart', onReaderClick, { passive: true })
}

// 平滑滚动到页面顶部，移植自 scroll.js 的返回顶部行为。
function scrollReaderToTop(scrollDuration) {
    const scrollStep = -window.scrollY / (scrollDuration / 15)
    const scrollInterval = setInterval(() => {
        if (window.scrollY !== 0) {
            window.scrollBy(0, scrollStep)
        } else {
            clearInterval(scrollInterval)
        }
    }, 15)
}

function getReaderBackTopButton() {
    if (!readerState.backTopButton) {
        readerState.backTopButton = document.getElementById('BackTopButton')
    }
    return readerState.backTopButton
}

function initReaderInput() {
    const input = document.getElementById('ReaderArchiveInput')
    const dropArea = document.getElementById('ReaderDropArea')
    if (!input || !dropArea) return

    input.addEventListener('change', () => openReaderArchive(input.files?.[0]))
    const headerTitle = getReaderHeaderTitle()
    if (headerTitle) {
        headerTitle.addEventListener('click', () => {
            if (readerState.book) {
                chooseReaderArchiveAgain()
            }
        })
    }
    for (const eventName of ['dragenter', 'dragover']) {
        dropArea.addEventListener(eventName, (event) => {
            event.preventDefault()
            dropArea.classList.add('border-blue-500')
        })
    }
    for (const eventName of ['dragleave', 'drop']) {
        dropArea.addEventListener(eventName, (event) => {
            event.preventDefault()
            dropArea.classList.remove('border-blue-500')
        })
    }
    dropArea.addEventListener('drop', (event) => {
        openReaderArchive(event.dataTransfer?.files?.[0])
    })
}

function initReaderBackTop() {
    const btn = getReaderBackTopButton()
    if (!btn) return
    btn.addEventListener('click', () => scrollReaderToTop(500))
    window.addEventListener('scroll', () => {
        const scrollTop = document.documentElement.scrollTop || document.body.scrollTop
        readerState.scrollDownFlag = scrollTop > readerState.scrollTopSave
        const step = readerState.scrollTopSave - scrollTop
        readerState.scrollTopSave = scrollTop
        if (step < -10 || step > 10) {
            readerState.showBackTopFlag = scrollTop > 400 && !readerState.scrollDownFlag
            btn.style.display = readerState.showBackTopFlag ? 'block' : 'none'
        }
        scheduleReaderCenterUpdate()
    }, { passive: true })
}

function initReaderResize() {
    const onResize = () => {
        Alpine.store('scroll').imageMaxWidth = window.innerWidth
        Alpine.store('global').checkOrientation()
        scheduleReaderCenterUpdate()
        // 横竖屏或窗口变化会影响单双页尺寸，需要重建翻页相邻页。
        updateReaderFlipImages()
    }
    onResize()
    window.addEventListener('resize', onResize)
}

document.addEventListener('DOMContentLoaded', () => {
    Alpine.store('global').onlineBook = false
    setReaderDefaultHeaderTitle()
    normalizeReaderReadMode()
    loadArchiveWasm().catch((error) => {
        console.error('[reader] preload wasm failed:', error)
        setReaderStatus(String(error?.message || error), 'error')
    })
    initReaderInput()
    initReaderBackTop()
    initReaderResize()
    initReaderGestures()
})

window.addEventListener('keydown', (event) => {
    if (Alpine.store('global').readMode !== 'page_flip' || !readerState.book) return
    const key = event.key.toLowerCase()
    if (key === 'arrowleft' || key === 'h' || key === ',' || key === '<' || key === 'pageup') {
        event.preventDefault()
        if (Alpine.store('flip').mangaMode) {
            toNextReaderFlipPage()
        } else {
            toPreviousReaderFlipPage()
        }
    }
    if (key === 'arrowright' || key === 'l' || key === '.' || key === '>' || key === 'pagedown' || key === ' ') {
        event.preventDefault()
        if (Alpine.store('flip').mangaMode) {
            toPreviousReaderFlipPage()
        } else {
            toNextReaderFlipPage()
        }
    }
    if (key === 'home') {
        jumpReaderFlipPage(1)
    }
    if (key === 'end') {
        jumpReaderFlipPage(Alpine.store('global').allPageNum)
    }
})
