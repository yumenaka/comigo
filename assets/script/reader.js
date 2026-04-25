// 本地压缩包阅读页面。文件只在浏览器内存中解压，不上传到服务器。
'use strict'

const readerState = {
    wasmReady: null,
    book: null,
    objectURLs: new Map(),
    observer: null,
    centerUpdateRaf: null,
    lastPageNum: 0,
    scrollTopSave: 0,
    scrollDownFlag: false,
    showBackTopFlag: false,
    backTopButton: null,
}

function readerText(key, fallback) {
    if (window.i18next && i18next.exists && i18next.exists(key)) {
        return i18next.t(key)
    }
    return fallback
}

function setReaderStatus(message, type = 'info') {
    const status = document.getElementById('ReaderStatus')
    if (!status) return
    status.textContent = message
    status.className = `w-full mt-4 text-sm text-center ${type === 'error' ? 'text-red-600' : 'opacity-80'}`
}

async function loadArchiveWasm() {
    if (readerState.wasmReady) {
        return readerState.wasmReady
    }
    readerState.wasmReady = (async () => {
        if (!window.Go) {
            await loadScript('/script/wasm/wasm_exec.js')
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

    const wasmResponse = await fetch('/script/wasm/archive.wasm')
    try {
        return await WebAssembly.instantiateStreaming(wasmResponse.clone(), importObject)
    } catch (_) {
        // 某些嵌入式环境不会给 .wasm 返回 application/wasm，退回 ArrayBuffer 加载。
        return WebAssembly.instantiate(await wasmResponse.arrayBuffer(), importObject)
    }
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
    cleanupReaderBook()
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
    for (const url of readerState.objectURLs.values()) {
        URL.revokeObjectURL(url)
    }
    readerState.objectURLs.clear()
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
    if (window.ComiGoArchive && typeof window.ComiGoArchive.close === 'function') {
        window.ComiGoArchive.close()
    }
}

function renderReaderBook(book) {
    Alpine.store('global').onlineBook = false
    Alpine.store('global').readMode = 'infinite_scroll'
    Alpine.store('global').nowPageNum = 1
    Alpine.store('global').allPageNum = book.page_count || book.PageInfos.length
    Alpine.store('scroll').allPageNum = Alpine.store('global').allPageNum

    const picker = document.getElementById('ReaderFilePicker')
    const shell = document.getElementById('ReaderShell')
    const title = document.getElementById('ReaderBookTitle')
    const mainArea = document.getElementById('ScrollMainArea')
    if (!mainArea || !shell) return

    if (picker) picker.classList.add('hidden')
    shell.classList.remove('hidden')
    shell.classList.add('flex')
    if (title) title.textContent = `${book.title} (${book.PageInfos.length})`
    document.title = `${book.title} - Comigo`

    mainArea.innerHTML = ''
    const fragment = document.createDocumentFragment()
    book.PageInfos.forEach((page, index) => {
        fragment.appendChild(createReaderPageNode(page, index, book.PageInfos.length))
    })
    mainArea.appendChild(fragment)
    Alpine.initTree(mainArea)
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
    if (readerState.objectURLs.has(index)) {
        img.src = readerState.objectURLs.get(index)
        return
    }
    try {
        const bytes = await window.ComiGoArchive.readPage(index)
        const page = readerState.book.PageInfos[index]
        const blob = new Blob([bytes], { type: guessMimeType(page?.name || '') })
        const objectURL = URL.createObjectURL(blob)
        readerState.objectURLs.set(index, objectURL)
        img.src = objectURL
    } catch (error) {
        console.error('[reader] load image failed:', error)
        img.alt = readerText('reader_page_load_failed', 'Failed to load page')
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
    event.currentTarget.style.cursor = getInReaderSettingArea(event) ? 'pointer' : ''
}

function initReaderGestures() {
    const mainArea = document.getElementById('ScrollMainArea')
    if (!mainArea || mainArea.dataset.readerGesturesReady === 'true') return

    mainArea.dataset.readerGesturesReady = 'true'
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
    const chooseAnotherButton = document.getElementById('ReaderChooseAnotherButton')
    if (!input || !dropArea) return

    input.addEventListener('change', () => openReaderArchive(input.files?.[0]))
    if (chooseAnotherButton) {
        chooseAnotherButton.addEventListener('click', () => {
            input.value = ''
            input.click()
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
    }
    onResize()
    window.addEventListener('resize', onResize)
}

document.addEventListener('DOMContentLoaded', () => {
    Alpine.store('global').onlineBook = false
    Alpine.store('global').readMode = 'infinite_scroll'
    loadArchiveWasm().catch((error) => {
        console.error('[reader] preload wasm failed:', error)
        setReaderStatus(String(error?.message || error), 'error')
    })
    initReaderInput()
    initReaderBackTop()
    initReaderResize()
    initReaderGestures()
})
