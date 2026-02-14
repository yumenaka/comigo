//此文件静态导入，不需要编译
//严格模式,不能使用未声明的变量
'use strict'

// ============ DOM 元素缓存（避免重复查询）============
const DOM = {
    header: null,
    range: null,
    sliderContainer: null,
    slider: null,
    leftSlide: null,
    rightSlide: null,
    flipMainArea: null,
    singleNowImage: null,
    doubleNowImageLeft: null,
    doubleNowImageRight: null,
    openSettingButton: null,
    reSortDropdownMenu: null,
    quickJumpDropdown: null,
}

// 初始化 DOM 缓存（在 DOMContentLoaded 后调用）
function initDOMCache() {
    DOM.header = document.getElementById('header')
    DOM.range = document.getElementById('StepsRangeArea')
    DOM.sliderContainer = document.getElementById('slider_container')
    DOM.slider = document.getElementById('slider')
    DOM.leftSlide = document.getElementById('left-slide')
    DOM.rightSlide = document.getElementById('right-slide')
    DOM.flipMainArea = document.getElementById('FlipMainArea')
    DOM.singleNowImage = document.getElementById('Single-NowImage')
    DOM.doubleNowImageLeft = document.getElementById('Double-NowImage-Left')
    DOM.doubleNowImageRight = document.getElementById('Double-NowImage-Right')
    DOM.openSettingButton = document.getElementById('OpenSettingButton')
    DOM.reSortDropdownMenu = document.getElementById('ReSortDropdownMenu')
    DOM.quickJumpDropdown = document.getElementById('QuickJumpDropdown')
}

// ============ 配置获取函数 ============
/**
 * 从 Alpine Store 获取配置
 * @returns {Object} 配置对象
 */
const getConfig = () => ({
    // 滑动动画配置
    swipeAnimationDuration: Alpine.store('flip').swipeAnimationDuration,
    resetAnimationDuration: Alpine.store('flip').resetAnimationDuration,
    swipeThreshold: Alpine.store('flip').swipeThreshold,
    swipeTimeout: Alpine.store('flip').swipeTimeout,
    // 其他配置
    preloadRange: Alpine.store('flip').preloadRange,
    wheelThrottleDelay: Alpine.store('flip').wheelThrottleDelay,
    // WebSocket 配置
    maxReconnectAttempts: Alpine.store('flip').websocketMaxReconnect,
    reconnectInterval: Alpine.store('flip').websocketReconnectInterval,
    // 静态配置
    setAreaRatio: 0.15,  // 设置区域比例
    slideBlockMargin: 30, // 边界回弹距离
})

// ============ 工具函数 ============
/**
 * 防抖函数
 * @param {Function} func - 要防抖的函数
 * @param {number} wait - 等待时间（毫秒）
 * @returns {Function} 防抖后的函数
 */
function debounce(func, wait) {
    let timeout
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout)
            func(...args)
        }
        clearTimeout(timeout)
        timeout = setTimeout(later, wait)
    }
}

/**
 * 节流函数
 * @param {Function} func - 要节流的函数
 * @param {number} limit - 时间限制（毫秒）
 * @returns {Function} 节流后的函数
 */
function throttle(func, limit) {
    let inThrottle
    return function executedFunction(...args) {
        if (!inThrottle) {
            func(...args)
            inThrottle = true
            setTimeout(() => inThrottle = false, limit)
        }
    }
}

// ============ 工具栏显示/隐藏 ============
let hideTimeout

// 显示工具栏
function showToolbar() {
    if (!DOM.header || !DOM.range) return
    
    if (Alpine.store('flip').autoHideToolbar) {
        DOM.header.style.opacity = '0.9'
        DOM.range.style.opacity = '0.9'
        DOM.header.style.transform = 'translateY(0)'
        DOM.range.style.transform = 'translateY(0)'
    } else {
        DOM.header.style.opacity = '1'
        DOM.range.style.opacity = '1'
        DOM.header.style.transform = 'translateY(0)'
        DOM.range.style.transform = 'translateY(0)'
    }
}

// 隐藏工具栏
function hideToolbar() {
    if (!DOM.header || !DOM.range) return
    
    if (Alpine.store('flip').autoHideToolbar) {
        DOM.header.style.opacity = '0'
        DOM.range.style.opacity = '0'
        DOM.header.style.transform = 'translateY(-100%)'
        DOM.range.style.transform = 'translateY(100%)'
    }
}

// 初始化工具栏事件监听器（在 DOM 加载后调用）
function initToolbarListeners() {
    if (!DOM.header) return
    
    // 显示工具栏
    DOM.header.addEventListener('mouseover', showToolbar)
    // 隐藏工具栏
    DOM.header.addEventListener('mouseout', hideToolbar)
    
    // 初始化：如果 autohidetoolbar 为真，则自动隐藏工具栏
    if (Alpine.store('flip').autoHideToolbar) {
        setTimeout(hideToolbar, 1000)
    }
}

// Base64编码静态资源图片（鼠标图标）：
// base64 -i SettingsOutline.svg ，然后// 把下面这行换成输出
const SettingsOutlineBase64 = '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'
// base64 -i Prohibited28Filled.png
const Prohibited28FilledBase64 = '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'
// base64 -i ArrowLeft.png
const ArrowLeftBase64 = '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'
// base64 -i ArrowRight.png
const ArrowRightBase64 = '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'


// ============ 全局数据初始化 ============
const book = JSON.parse(document.getElementById('NowBook').textContent)
const images = book.PageInfos
Alpine.store('global').allPageNum = parseInt(book.page_count)

// ============ 滑动相关变量 ============
let touchStartX = 0
let touchEndX = 0
let isSwiping = false
let currentTranslate = 0
let startTime = 0
let animationID = 0

// 设置图片资源，预加载等
// 需要 HTTP 响应头中允许缓存（没有使用 Cache-Control: no-cache），也就是 gin 不能用htmx/router/server.go 里面的noCache 中间件。
// 在预加载用到的图片资源 URL
let preloadedImages = new Set()

// 首次加载时设置图片
setImageSrc()

// 首次加载时，检查URL参数中是否有页码
const url = new URL(window.location.href);
const params = new URLSearchParams(url.search);
if (params.has('start')) {
    // 优先使用URL参数中的页码
    let pageNum = parseInt(params.get('start'))
    jumpPageNum(pageNum, false); // 使用跳转函数更新页面
    // 翻页模式自动跳转后删除start查询参数
    params.delete('start');
    if (params.toString() !== '') {
        const newUrl = `${url.origin}${url.pathname}?${params.toString()}`;
        window.history.replaceState({}, document.title, newUrl);
    }
    if (params.toString() === '') {
        const newUrl = `${url.origin}${url.pathname}`;
        window.history.replaceState({}, document.title, newUrl);
    }
} else {
    // 页面加载时读取本地存储的页码
    Alpine.store('global').loadPageNumFromLocalStorage(book.id, () => {
        let pageNum = parseInt(localStorage.getItem(`pageNum_${book.id}`))
        jumpPageNum(pageNum, false); // 使用跳转函数更新页面
    });
}


//判断当前浏览器是不是Safari，暂时没啥用
// const isSafari = navigator.userAgent.indexOf('Safari') !== -1 && navigator.userAgent.indexOf('Chrome') === -1

/**
 * 获取图片源 URL
 * @param {number} index - 图片索引
 * @returns {string|undefined} 图片 URL
 */
function GetImageSrc(index) {
    if (index < 0 || index >= images.length) {
        if (Alpine.store('global').debugMode) {
            console.error(`GetImageSrc: 索引越界 ${index}`)
        }
        return
    }
    
    const url = images[index].url
    if (!Alpine.store('global').onlineBook) {
        return url
    }
    
    // 在线书籍：添加图片处理参数
    const autoCrop = Alpine.store('global').autoCrop 
        ? `&auto_crop=${Alpine.store('global').autoCropNum}` 
        : ''
    const autoResize = Alpine.store('global').autoResize 
        ? `&resize_max_width=${Alpine.store('global').autoResizeWidth}` 
        : ''
    const noCache = Alpine.store('global').noCache ? '&no-cache=true' : ''
    
    return `${url}${autoCrop}${autoResize}${noCache}`
}

/**
 * 加载图片资源
 */
function setImageSrc() {
    const nowPageNum = Alpine.store('global').nowPageNum
    const allPageNum = Alpine.store('global').allPageNum
    const mangaMode = Alpine.store('flip').mangaMode
    const config = getConfig()
    
    if (nowPageNum === 0 || nowPageNum > allPageNum) {
        if (Alpine.store('global').debugMode) {
            console.log('setImageSrc: nowPageNum is out of range', nowPageNum)
        }
        return
    }
    
    // 加载当前图片
    if (DOM.singleNowImage) {
        DOM.singleNowImage.src = GetImageSrc(nowPageNum - 1)
    }
    
    // 根据漫画模式设置双页图片
    if (mangaMode) {
        if (DOM.doubleNowImageRight) {
            DOM.doubleNowImageRight.src = GetImageSrc(nowPageNum - 1)
        }
    } else {
        if (DOM.doubleNowImageLeft) {
            DOM.doubleNowImageLeft.src = GetImageSrc(nowPageNum - 1)
        }
    }
    
    preloadedImages.add(GetImageSrc(nowPageNum - 1))
    
    // 更新滑动容器图片
    updateSliderImages(nowPageNum)
    
    // 为双页模式预加载下一张图片
    if (nowPageNum < allPageNum) {
        const nextImgSrc = GetImageSrc(nowPageNum)
        if (mangaMode) {
            if (DOM.doubleNowImageLeft) {
                DOM.doubleNowImageLeft.src = nextImgSrc
            }
        } else {
            if (DOM.doubleNowImageRight) {
                DOM.doubleNowImageRight.src = nextImgSrc
            }
        }
        preloadedImages.add(nextImgSrc)
    }
    
    // 预加载前后图片
    const preloadRange = config.preloadRange
    for (let i = nowPageNum - 2; i <= nowPageNum + preloadRange; i++) {
        if (i >= 0 && i < allPageNum) {
            const imgUrl = GetImageSrc(i)
            if (!preloadedImages.has(imgUrl)) {
                const img = new Image()
                img.src = imgUrl
                preloadedImages.add(imgUrl)
            }
        }
    }
}


/**
 * 创建图片元素
 * @param {string} src - 图片源
 * @param {string} className - CSS 类名
 * @returns {HTMLImageElement} 图片元素
 */
function createImageElement(src, className) {
    const img = document.createElement('img')
    img.src = src
    img.className = className
    img.draggable = false
    return img
}

/**
 * 更新滑动容器图片
 * @param {number} nowPageNum - 当前页码
 */
function updateSliderImages(nowPageNum) {
    if (!DOM.leftSlide || !DOM.rightSlide || !DOM.slider) return
    
    const mangaMode = Alpine.store('flip').mangaMode
    const doublePageMode = Alpine.store('flip').doublePageMode
    const allPageNum = Alpine.store('global').allPageNum
    const isPortrait = Alpine.store('global').isPortrait
    
    // 根据阅读方向设置滑动元素的位置
    if (mangaMode) {
        DOM.leftSlide.style.transform = 'translateX(100%)'
        DOM.rightSlide.style.transform = 'translateX(-100%)'
    } else {
        DOM.leftSlide.style.transform = 'translateX(-100%)'
        DOM.rightSlide.style.transform = 'translateX(100%)'
    }
    
    // CSS 类名
    const singlePageClass = isPortrait 
        ? 'object-contain w-auto max-w-full h-screen' 
        : 'h-screen w-auto max-w-full object-contain'
    const doublePageClass = 'object-contain w-auto max-h-screen m-0 select-none max-w-1/2 grow-0'
    const singleImgClass = 'object-contain h-screen max-w-full max-h-screen m-0'
    
    // ========== 单页模式 ==========
    if (!doublePageMode) {
        // 前一张图片
        DOM.leftSlide.innerHTML = ''
        if (nowPageNum > 1) {
            const prevImg = createImageElement(GetImageSrc(nowPageNum - 2), singlePageClass)
            DOM.leftSlide.appendChild(prevImg)
        }
        
        // 更新当前图片
        if (DOM.singleNowImage && nowPageNum >= 1 && nowPageNum <= allPageNum) {
            DOM.singleNowImage.src = GetImageSrc(nowPageNum - 1)
        }
        
        // 后一张图片
        DOM.rightSlide.innerHTML = ''
        if (nowPageNum < allPageNum) {
            const nextImg = createImageElement(GetImageSrc(nowPageNum), singlePageClass)
            DOM.rightSlide.appendChild(nextImg)
        }
    }
    
    // ========== 双页模式 ==========
    if (doublePageMode) {
        // 前一屏图片
        DOM.leftSlide.innerHTML = ''
        if (nowPageNum === 2) {
            const prevImg = createImageElement(GetImageSrc(nowPageNum - 2), singleImgClass)
            DOM.leftSlide.appendChild(prevImg)
        } else if (nowPageNum >= 3) {
            const prevImg_1 = createImageElement(GetImageSrc(nowPageNum - 2), doublePageClass)
            const prevImg_2 = createImageElement(GetImageSrc(nowPageNum - 3), doublePageClass)
            if (mangaMode) {
                DOM.leftSlide.appendChild(prevImg_1)
                DOM.leftSlide.appendChild(prevImg_2)
            } else {
                DOM.leftSlide.appendChild(prevImg_2)
                DOM.leftSlide.appendChild(prevImg_1)
            }
        }
        
        // 后一屏图片
        DOM.rightSlide.innerHTML = ''
        if (nowPageNum === allPageNum - 3) {
            const nextImg = createImageElement(GetImageSrc(nowPageNum - 2), singleImgClass)
            DOM.rightSlide.appendChild(nextImg)
        } else if (nowPageNum < allPageNum - 3) {
            const nextImg_1 = createImageElement(GetImageSrc(nowPageNum + 1), doublePageClass)
            const nextImg_2 = createImageElement(GetImageSrc(nowPageNum + 2), doublePageClass)
            if (mangaMode) {
                DOM.rightSlide.appendChild(nextImg_2)
                DOM.rightSlide.appendChild(nextImg_1)
            } else {
                DOM.rightSlide.appendChild(nextImg_1)
                DOM.rightSlide.appendChild(nextImg_2)
            }
        }
    }

    // 确保滑动容器在更新图片后回到初始位置（没有动画）
    DOM.slider.style.transition = 'none'
    DOM.slider.style.transform = 'translateX(0)'
    DOM.slider.offsetHeight // 触发重排
    DOM.slider.style.transition = ''
    resetSlider()
}

/**
 * 重置滑动状态
 */
function resetSlider() {
    cancelAnimationFrame(animationID)
    currentTranslate = 0
}

/**
 * 触摸开始事件处理
 * @param {TouchEvent|MouseEvent} e - 事件对象
 */
function touchStart(e) {
    if (!Alpine.store('flip').swipeTurn) return
    
    startTime = new Date().getTime()
    isSwiping = true
    touchStartX = e.type === 'touchstart' ? e.touches[0].clientX : e.clientX
    
    // 停止任何正在进行的动画
    cancelAnimationFrame(animationID)
}

/**
 * 判断是否应该阻止滚动（在边界时）
 * @param {number} diffX - 滑动距离
 * @returns {boolean} 是否阻止滚动
 */
function shouldBlockScroll(diffX) {
    const mangaMode = Alpine.store('flip').mangaMode
    const nowPageNum = Alpine.store('global').nowPageNum
    const allPageNum = Alpine.store('global').allPageNum
    
    // 第一页尝试向前翻
    if (nowPageNum === 1) {
        return (diffX < 0 && mangaMode) || (diffX > 0 && !mangaMode)
    }
    
    // 最后一页尝试向后翻
    if (nowPageNum === allPageNum) {
        return (diffX > 0 && mangaMode) || (diffX < 0 && !mangaMode)
    }
    
    return false
}

/**
 * 触摸移动事件处理
 * @param {TouchEvent|MouseEvent} e - 事件对象
 */
function touchMove(e) {
    if (!isSwiping || !Alpine.store('flip').swipeTurn || !DOM.slider)
        return
    
    const currentX = e.type === 'touchmove' ? e.touches[0].clientX : e.clientX
    const diffX = currentX - touchStartX
    const config = getConfig()
    
    // 设置当前滑动距离
    currentTranslate = diffX
    
    // 如果在第一页或最后一页尝试向前翻或向后翻，阻止默认滚动
    if (shouldBlockScroll(diffX)) {
        currentTranslate = diffX < 0 ? -config.slideBlockMargin : config.slideBlockMargin
    }
    
    // 应用变换
    DOM.slider.style.transform = `translateX(${currentTranslate}px)`
    
    // 防止页面滚动
    if (Math.abs(diffX) > 10) {
        e.preventDefault()
    }
}

/**
 * 触摸结束事件处理
 * @param {TouchEvent|MouseEvent} e - 事件对象
 */
function touchEnd(e) {
    if (!isSwiping || !Alpine.store('flip').swipeTurn)
        return
    
    isSwiping = false
    const config = getConfig()
    const endTime = new Date().getTime()
    const timeElapsed = endTime - startTime
    touchEndX = e.type === 'touchend' ? e.changedTouches[0].clientX : e.clientX
    const diffX = touchEndX - touchStartX
    
    // 判断是否应该翻页（基于滑动距离和速度）
    const isQuickSwipe = timeElapsed < config.swipeTimeout && Math.abs(diffX) > 50
    
    // 用于记录滑动方向
    let direction = null
    if (diffX < -config.swipeThreshold || (isQuickSwipe && diffX < 0)) {
        direction = 'left'
    } else if (diffX > config.swipeThreshold || (isQuickSwipe && diffX > 0)) {
        direction = 'right'
    }
    
    // 如果在第一页或最后一页尝试向前翻或向后翻，或没有足够的滑动距离，回到原始位置
    if (shouldBlockScroll(diffX) || direction === null) {
        animateReset()
        return
    }
    
    // 执行滑动动画及后续翻页
    animateSlide(direction)
}

/**
 * 滑动动画 - 处理滑动动画和翻页逻辑
 * @param {string} direction - 滑动方向 ('left' 或 'right')
 */
function animateSlide(direction) {
    if (!DOM.slider) return
    
    const config = getConfig()
    const width = window.innerWidth
    const targetPosition = direction === 'left' ? -width : width
    const mangaMode = Alpine.store('flip').mangaMode
    const startPosition = currentTranslate
    const duration = config.swipeAnimationDuration
    
    let animStartTime = null
    
    function animate(timestamp) {
        if (!animStartTime) animStartTime = timestamp
        const elapsed = timestamp - animStartTime
        const progress = Math.min(elapsed / duration, 1)
        const easeProgress = easeOutCubic(progress)
        const position = startPosition + (targetPosition - startPosition) * easeProgress

        DOM.slider.style.transform = `translateX(${position}px)`

        if (progress < 1) {
            animationID = requestAnimationFrame(animate)
        } else {
            // 动画完成后执行翻页逻辑
            const flipFunction = mangaMode
                ? (direction === 'left' ? toPreviousPage : toNextPage)
                : (direction === 'left' ? toNextPage : toPreviousPage)
            
            if (flipFunction) {
                flipFunction()
            } else {
                animateReset()
            }
        }
    }

    animationID = requestAnimationFrame(animate)
}

/**
 * 回弹动画 - 将滑块回到原始位置
 */
function animateReset() {
    if (!DOM.slider) return
    
    const config = getConfig()
    const startPosition = currentTranslate
    const duration = config.resetAnimationDuration
    
    let animStartTime = null
    
    function animate(timestamp) {
        if (!animStartTime) animStartTime = timestamp
        const elapsed = timestamp - animStartTime
        const progress = Math.min(elapsed / duration, 1)
        const easeProgress = easeOutCubic(progress)
        const position = startPosition * (1 - easeProgress)
        
        DOM.slider.style.transform = `translateX(${position}px)`
        
        if (progress < 1) {
            animationID = requestAnimationFrame(animate)
        } else {
            DOM.slider.style.transform = 'translateX(0)'
            resetSlider()
        }
    }

    animationID = requestAnimationFrame(animate)
}

/**
 * 缓动函数 - 使动画更自然
 * @param {number} x - 进度值 (0-1)
 * @returns {number} 缓动后的值
 */
function easeOutCubic(x) {
    return 1 - Math.pow(1 - x, 3)
}

/**
 * 初始化滑动事件监听器
 */
function initSlideListeners() {
    if (!DOM.sliderContainer) return
    
    // 触摸事件（移动设备）
    DOM.sliderContainer.addEventListener('touchstart', touchStart)
    DOM.sliderContainer.addEventListener('touchmove', touchMove, {passive: false})
    DOM.sliderContainer.addEventListener('touchend', touchEnd)
    
    // 鼠标事件（PC）
    DOM.sliderContainer.addEventListener('mousedown', touchStart)
    DOM.sliderContainer.addEventListener('mousemove', touchMove)
    DOM.sliderContainer.addEventListener('mouseup', touchEnd)
    DOM.sliderContainer.addEventListener('mouseleave', touchEnd)
    
    // 初始化滑动容器中的图片
    const nowPageNum = Alpine.store('global').nowPageNum
    updateSliderImages(nowPageNum)
}

/**
 * 翻页函数 - 增加或减少页码
 * @param {number} n - 要增加的页数（负数表示减少）
 */
function addPageNum(n = 1) {
    const nowPageNum = parseInt(Alpine.store('global').nowPageNum)
    const allPageNum = Alpine.store('global').allPageNum
    const targetPage = nowPageNum + n
    
    // 检查边界
    if (targetPage > allPageNum) {
        showToast(i18next.t('hint_last_page'), 'warning')
        return
    }
    if (targetPage < 1) {
        showToast(i18next.t('hint_first_page'), 'warning')
        return
    }
    
    // 更新页码
    Alpine.store('global').nowPageNum = targetPage
    
    // 更新书签（仅在线书籍）
    if (book && book.id && Alpine.store('global').onlineBook) {
        Alpine.store('global').UpdateBookmark({
            type: 'auto',
            bookId: book.id,
            pageIndex: targetPage,
        })
    }
    
    // 更新图片
    setImageSrc()
    
    // 更新标签页标题
    setTitle()
    
    // 通过 WebSocket 同步翻页数据
    if (Alpine.store('global').syncPageByWS) {
        sendFlipData()
    }
    
    // 保存页数到本地存储
    Alpine.store('global').savePageNumToLocalStorage(book.id)
}

/**
 * 从输入框跳转到指定页
 * @param {Event} event - 输入事件
 */
function inputPageNum(event) {
    const i = parseInt(event.target.value)
    const num = Alpine.store('flip').mangaMode 
        ? (Alpine.store('global').allPageNum - i + 1) 
        : i
    jumpPageNum(num)
}

/**
 * 跳转到指定页
 * @param {number} jumpNum - 目标页码
 * @param {boolean} sync - 是否同步到其他设备（默认 true）
 */
function jumpPageNum(jumpNum, sync = true) {
    const num = parseInt(jumpNum)
    const allPageNum = Alpine.store('global').allPageNum

    if (num <= 0 || num > allPageNum) {
        alert(i18next.t('hintPageNumOutOfRange'))
        return
    }
    
    Alpine.store('global').nowPageNum = num
    
    if (Alpine.store('global').onlineBook) {
        Alpine.store('global').UpdateBookmark({
            type: 'auto',
            bookId: book.id,
            pageIndex: num,
        })
        
        if (sync && Alpine.store('global').syncPageByWS) {
            sendFlipData()
        }
    }
    
    Alpine.store('global').savePageNumToLocalStorage(book.id)
    setImageSrc()
}

// 翻页函数，下一页
function toNextPage() {
    let doublePageMode = Alpine.store('flip').doublePageMode === true
    let nowPageNum = parseInt(Alpine.store('global').nowPageNum)
    // 单页模式
    if (!doublePageMode) {
        if (nowPageNum <= Alpine.store('global').allPageNum) {
            addPageNum(1)
        }
    }
    //双页模式
    if (doublePageMode) {
        if (nowPageNum < Alpine.store('global').allPageNum - 1) {
            addPageNum(2)
        } else {
            addPageNum(1)
        }
    }
}

// 翻页函数，前一页
function toPreviousPage() {
    //错误值,第0或第1页。
    if (Alpine.store('global').nowPageNum <= 1) {
        showToast(i18next.t('hint_first_page'), 'warning')
        return
    }
    //双页模式
    if (Alpine.store('flip').doublePageMode) {
        if (Alpine.store('global').nowPageNum - 2 > 0) {
            addPageNum(-2)
        } else {
            addPageNum(-1)
        }
    } else {
        addPageNum(-1)
    }
}


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

// 翻页模式功能：显示工具栏时，点击设置区域，自动漫画区域居中。
function scrollToMangaMain() {
    if (!Alpine.store('flip').autoHideToolbar) {
        // 将 slider_container 顶部对齐到浏览器可见区域顶部
        const mangaMain = document.getElementById('slider_container')
        mangaMain.scrollIntoView({
            behavior: 'smooth', // 平滑滚动
            block: 'start', // 与可视区顶部对齐
        })
    }
}

//获取鼠标位置,决定是否打开设置面板
function onMouseClick(e) {
    // 如果正在滑动，则不处理点击事件
    if (isSwiping || Math.abs(currentTranslate) > 10) {
        return
    }
    //获取鼠标的X坐标（鼠标与屏幕左侧的距离,单位为px）
    let clickX = e.x
    //浏览器的视口宽,不包括工具栏和滚动条:
    let innerWidth = window.innerWidth
    let inSetArea = getInSetArea(e)
    if (inSetArea) {
        // 高度对齐
        if (Alpine.store('flip').autoAlign) {
            scrollToMangaMain()
        }
        // 如果工具栏是隐藏的，点击设置区域时，显示工具栏
        if (Alpine.store('flip').autoHideToolbar === true) {
            showToolbar()
        }
        //获取ID为 OpenSettingButton的元素，然后模拟点击
        document.getElementById('OpenSettingButton').click()
    }
    if (!inSetArea) {
        //决定如何翻页
        if (clickX < innerWidth * 0.5) {
            //左边的翻页
            if (!Alpine.store('flip').mangaMode) {
                toPreviousPage()
            } else {
                toNextPage()
            }
        } else {
            //右边的翻页
            if (!Alpine.store('flip').mangaMode) {
                toNextPage()
            } else {
                toPreviousPage()
            }
        }
    }
}


//获取鼠标位置,决定如何显示鼠标指针
function onMouseMove(e) {
    // https://developer.mozilla.org/zh-CN/docs/Web/API/Event/stopPropagation
    // 相当于 Alpinejs的    @click.stop="onMouseMove()" ? https://alpinejs.dev/directives/on#prevent
    let clickX = e.x //获取鼠标的X坐标（鼠标与屏幕左侧的距离,单位为px）
    // 浏览器的视口宽,不包括工具栏和滚动条:
    let innerWidth = window.innerWidth
    //在设置区域
    let inSetArea = getInSetArea(e)
    // https://developer.mozilla.org/zh-CN/docs/Web/CSS/cursor
    //e.currentTarget.style.cursor = 'default'
    if (inSetArea) {
        e.currentTarget.style.cursor = `url("data:image/png;base64,${SettingsOutlineBase64}") 12 12, pointer`
        showToolbar()
    }
    let stepsRangeArea = document
        .getElementById('StepsRangeArea')
        .getBoundingClientRect()
    //判断鼠标是否在翻页条上
    let inRangeArea =
        clickX >= stepsRangeArea.left &&
        clickX <= stepsRangeArea.right &&
        e.y >= stepsRangeArea.top &&
        e.y <= stepsRangeArea.bottom
    // 判断鼠标是否在翻页条上,如果在翻页条上,就设置为默认的鼠标指针
    if (inRangeArea) {
        e.currentTarget.style.cursor = 'default'
    }
    //设置鼠标指针
    if (!inSetArea && !inRangeArea) {
        if (clickX < innerWidth * 0.5) {
            //设置左边的鼠标指针
            if (
                !Alpine.store('flip').mangaMode &&
                Alpine.store('global').nowPageNum === 1
            ) {
                //右边翻下一页,且目前是第一页的时候,左边的鼠标指针,设置为禁止翻页
                e.currentTarget.style.cursor =
                    `url("data:image/png;base64,${Prohibited28FilledBase64}") 12 12, pointer`
            } else if (
                Alpine.store('flip').mangaMode &&
                Alpine.store('global').nowPageNum === Alpine.store('global').allPageNum
            ) {
                //左边翻下一页,且目前是最后一页的时候,左边的鼠标指针,设置为禁止翻页
                e.currentTarget.style.cursor =
                    `url("data:image/png;base64,${Prohibited28FilledBase64}") 12 12, pointer`
            } else {
                //正常情况下,左边是向左的箭头
                e.currentTarget.style.cursor = `url("data:image/png;base64,${ArrowLeftBase64}") 12 12, pointer`
            }
        } else {
            //设置右边的鼠标指针
            if (
                !Alpine.store('flip').mangaMode &&
                Alpine.store('global').nowPageNum === Alpine.store('global').allPageNum
            ) {
                //右边翻下一页,且目前是最后页的时候,右边的鼠标指针,设置为禁止翻页
                e.currentTarget.style.cursor =
                    `url("data:image/png;base64,${Prohibited28FilledBase64}") 12 12, pointer`
            } else if (
                Alpine.store('flip').mangaMode &&
                Alpine.store('global').nowPageNum === 1
            ) {
                //左边翻下一页,且目前是第一页的时候,右边的鼠标指针,设置为禁止翻页
                e.currentTarget.style.cursor =
                    `url("data:image/png;base64,${Prohibited28FilledBase64}") 12 12, pointer`
            } else {
                //正常情况下,右边是向右的箭头
                e.currentTarget.style.cursor = `url("data:image/png;base64,${ArrowRightBase64}") 12 12, pointer`
            }
        }
    }
}

/**
 * 获取元素的边界信息
 * @returns {Object} 各元素的边界矩形
 */
function getElementsRect() {
    return {
        rect1_header: DOM.header ? DOM.header.getBoundingClientRect() : null,
        rect2_range: DOM.range ? DOM.range.getBoundingClientRect() : null,
        rect3_sort_dropdown: DOM.reSortDropdownMenu ? DOM.reSortDropdownMenu.getBoundingClientRect() : null,
        rect4_dropdown_quick_jump: DOM.quickJumpDropdown ? DOM.quickJumpDropdown.getBoundingClientRect() : null,
        rect5_steps_range_area: DOM.range ? DOM.range.getBoundingClientRect() : null,
    }
}

/**
 * 检查矩形区域是否包含坐标点
 * @param {DOMRect|null} rect - 矩形区域
 * @param {number} x - X 坐标
 * @param {number} y - Y 坐标
 * @returns {boolean} 是否在区域内
 */
function isPointInRect(rect, x, y) {
    if (!rect) return false
    return x >= rect.left && x <= rect.right && y >= rect.top && y <= rect.bottom
}

document.addEventListener('mousemove', function (event) {
    const {
        rect1_header,
        rect2_range,
        rect3_sort_dropdown,
        rect4_dropdown_quick_jump,
        rect5_steps_range_area,
    } = getElementsRect()
    const x = event.clientX
    const y = event.clientY
    let inInElement1 = false
    let inInElement2 = false
    
    // 因为 header 需要收起来，所以不能用 left、right、top、bottom 判断 y 是否在 header 的范围内
    // 现在设定为固定的 80px 高度，这样会比较自然
    if (Alpine.store('flip').autoHideToolbar) {
        inInElement1 = (y <= 80)
        inInElement2 = (y >= window.innerHeight - 80)
    } else {
        // 如果工具栏不自动隐藏，用边界判断
        inInElement1 = isPointInRect(rect1_header, x, y)
        inInElement2 = isPointInRect(rect2_range, x, y)
    }

    // 判断鼠标是否在各个下拉菜单范围内
    const inInElement3 = isPointInRect(rect3_sort_dropdown, x, y)
    const inInElement4 = isPointInRect(rect4_dropdown_quick_jump, x, y)
    const inInElement5 = isPointInRect(rect5_steps_range_area, x, y)

    // 鼠标在设置区域
    const inSetArea = getInSetArea(event)
    
    // 如果鼠标在任何感兴趣的区域，显示工具栏
    if (inSetArea || inInElement1 || inInElement2 || inInElement3 || inInElement4 || inInElement5) {
        showToolbar()
    } else {
        hideToolbar()
    }
})

//可见区域变化时，改变页面状态
function onResize() {
    Alpine.store('flip').imageMaxWidth = window.innerWidth
    let clientWidth = document.documentElement.clientWidth
    let clientHeight = document.documentElement.clientHeight
    // var aspectRatio = window.innerWidth / window.innerHeight
    let aspectRatio = clientWidth / clientHeight
    // 为了调试的时候方便,阈值是正方形
    if (aspectRatio > 19 / 19) {
        Alpine.store('flip').isLandscapeMode = true
        Alpine.store('flip').isPortraitMode = false
    } else {
        Alpine.store('flip').isLandscapeMode = false
        Alpine.store('flip').isPortraitMode = true
    }
}

//初始化时,执行一次onResize()
onResize()
//文档视图调整大小时触发 resize 事件。 https://developer.mozilla.org/zh-CN/docs/Web/API/Window/resize_event
window.addEventListener('resize', onResize)

//离开区域的时候,清空鼠标样式
function onMouseLeave(e) {
    e.currentTarget.style.cursor = ''
}

/**
 * 初始化 FlipMainArea 事件监听器
 */
function initFlipMainAreaListeners() {
    if (!DOM.flipMainArea) return
    
    // 鼠标移动时触发移动事件
    DOM.flipMainArea.addEventListener('mousemove', onMouseMove)
    // 点击的时候触发点击事件
    DOM.flipMainArea.addEventListener('click', onMouseClick)
    // 触摸开始时触发点击事件
    DOM.flipMainArea.addEventListener('touchstart', onMouseClick)
    // 离开的时候触发离开事件
    DOM.flipMainArea.addEventListener('mouseleave', onMouseLeave)
    
    // 绑定滚轮事件（使用 passive: false 以允许阻止默认滚动行为）
    DOM.flipMainArea.addEventListener('wheel', onWheel, {passive: false})
}

// ============ WebSocket 同步（使用共享模块 ComiGoWS）============
/** 通过 ComiGoWS 发送翻页数据，供 jumpPageNum、setImageSrc 等处调用 */
function sendFlipData() {
    if (typeof window.ComiGoWS === 'undefined') return
    window.ComiGoWS.send('flip_mode_sync_page', {
        book_id: book.id,
        now_page_num: Alpine.store('global').nowPageNum,
    }, '翻页模式，发送数据')
}

// 设置标签页标题
function setTitle(name) {
    let numStr = ''
    if (Alpine.store('flip').showPageNum) {
        numStr = ` ${Alpine.store('global').nowPageNum}/${Alpine.store('global').allPageNum} `
    }
    //简化标题
    if (Alpine.store('shelf').simplifyTitle) {
        document.title = `${numStr}${shortName(book.title)}`;
    } else {
        document.title = `${numStr}${book.title}`;
    }
}

setTitle();

/**
 * 生成简短标题
 * @param {string} title – 原始标题
 * @returns {string} – 处理后的短标题
 */
function shortName(title) {
    let shortTitle = title;

    /* 1. 移除常见文件扩展名（忽略大小写） */
    shortTitle = shortTitle.replace(
        /\.(epub|pdf|mobi|azw3|cbz|cbr|zip|rar|7z|txt|docx?)$/i,
        ""
    );

    /* 2. 顺序移除各种括号及其内容 */
    shortTitle = shortTitle
        .replace(/\([^)]*\)/g, "")   // ()
        .replace(/\[[^\]]*\]/g, "")  // []
        .replace(/（[^）]*）/g, "")   // （）
        .replace(/【[^】]*】/g, "");  // 【】

    /* 3. 移除域名（含 http/https 前缀） */
    shortTitle = shortTitle.replace(/https?:\/\/[^\s/]+/gi, "");

    /* 4 & 5. 去掉前后空白 */
    shortTitle = shortTitle.trimStart();
    shortTitle = shortTitle.trimEnd();

    /* 6. 去除开头的标点符号（使用 Unicode 属性，需要 Node ≥ v10） */
    shortTitle = shortTitle.replace(/^[\p{P}\p{S}]+/u, "");

    /* 7. 最后再 trim 一次，防止前一步留下空格 */
    shortTitle = shortTitle.trim();

    /* 将字符串按 Unicode 码点拆分 */
    const runes = Array.from(shortTitle);
    const originalRunes = Array.from(title);
    // 简化后过短（<2）
    if (runes.length < 2) {
        if (originalRunes.length > 15) {
            return originalRunes.slice(0, 15).join("") + "…";
        }
        if (originalRunes.length > 0) {
            return originalRunes.length <= 15
                ? title
                : originalRunes.slice(0, 15).join("") + "…";
        }
        return "";
    }
    // 简化后 ≤15：直接返回
    if (runes.length <= 15) {
        return shortTitle;
    }
    // 超过 15：截断并加省略号
    return runes.slice(0, 15).join("") + "…";
}

// 设定键盘快捷键
/* 记录方向键当前的按压状态 */
// 1) 方向/动作当前状态
const state = {up: false, down: false, left: false, right: false, fire: false};

/* 2) 键 → 动作 的映射表
 *   - 左边写 `event.key`（大小写无关，统一用小写）。
 *   - 键盘键位表：https://developer.mozilla.org/zh-CN/docs/Web/API/UI_Events/Keyboard_event_key_values
 *   - 右边写动作名称（小写）
 *   - 这里的动作名称可以是任意字符串，建议用小写
 *   - 同一个动作可以对应多组键：方向键 + WASD + 自定义
 */
const keyMap = {
    // 方向键 ↑
    arrowup: "up",
    // 方向键 ↓
    arrowdown: "down",
    // 方向键 ←
    arrowleft: "left",
    // 方向键 →
    arrowright: "right",
    // 长得像方向键的键位当作方向键
    "<": "left",
    ">": "right",
    // 英语键盘上，与 < 键在一起
    ",": "left",
    // 英语键盘上，与 > 键在一起
    ".": "right",
    // vim键位 hjkl 当做方向键
    h: "left",
    j: "down",
    k: "up",
    l: "right",
    // 游戏当中常用的 WSAD 当做方向键
    w: "up",
    s: "down",
    a: "left",
    d: "right",
    // Home 键
    home: "first_page",
    // End 键
    end: "last_page",
    // PageUp 键
    pageup: "pre_page",
    // PageDown 键
    pagedown: "next_page",
    // 加减相关键位当作方翻页键
    // + 键
    "+": "next_page",
    // - 键
    "-": "pre_page",
    // = 键 英语键盘，与 + 键在一起
    "=": "next_page",
    // —— 键 英语键盘，与 - 键在一起
    "——": "pre_page",
    // 空格键
    " ": "next_page",
};

// 3) 通用按键处理器：down=true 表示按下，false 表示松开
function handle(e, down) {
    const k = e.key.toLowerCase();      // 统一小写
    const act = keyMap[k];              // 查映射表
    if (!act) return;                   // 映射表里没有，忽略
    state[act] = down;                  // 更新状态
    //e.preventDefault();               // 阻止滚动等默认行为（可选）
    // 上一页
    if (act === "pre_page" && down) {
        toPreviousPage()
    }
    // 下一页
    if (act === "next_page" && down) {
        toNextPage()
    }
    // 触按下相当于左方向键的按键的时候
    if (act === "left" && down) {
        if (Alpine.store('flip').mangaMode) {
            toNextPage()
        } else {
            toPreviousPage()
        }
    }
    // 触按下相当于右方向键的按键的时候
    if (act === "right" && down) {
        if (Alpine.store('flip').mangaMode) {
            toPreviousPage()
        } else {
            toNextPage()
        }
    }
    // 直接跳转到第一页,同时长按的时候不执行多次
    if (act === "first_page" && down && !e.repeat) {
        jumpPageNum(1)
    }
    // 直接跳转到最后一页,同时长按的时候不执行多次
    if (act === "last_page" && down && !e.repeat) {
        jumpPageNum(Alpine.store('global').allPageNum)
    }
}

// 4) 事件监听
// 监听键盘事件
// keydown 第一次按下和按住时会触发 //统一禁止长按： addEventListener("keydown", e => !e.repeat && handle(e, true));
addEventListener("keydown", e => handle(e, true));
// keyup 松开时触发
addEventListener("keyup", e => handle(e, false));

// // 2 手柄方向键 (Gamepad API) TODO
// // https://developer.mozilla.org/zh-CN/docs/Web/API/Gamepad_API/Using_the_Gamepad_API
// const gamepads = {};          // 按 index 存储 Gamepad 对象
//
// window.addEventListener("gamepadconnected",   e => {
// 	gamepads[e.gamepad.index] = e.gamepad;
// 	console.log("已连接:", e.gamepad.id);
// });
//
// window.addEventListener("gamepaddisconnected", e => {
// 	delete gamepads[e.gamepad.index];
// 	console.log("已断开:", e.gamepad.id);
// });

// ============ 鼠标滚轮翻页功能 ============
let wheelThrottleTimer = null

/**
 * 滚轮事件处理函数
 * @param {WheelEvent} e - 滚轮事件对象
 */
function onWheel(e) {
    if (!Alpine.store('flip').wheelFlip) return
    
    // 如果正在滑动翻页，则不处理滚轮事件
    if (isSwiping || Math.abs(currentTranslate) > 10) return
    
    const deltaY = e.deltaY
    const config = getConfig()
    
    // 如果滚轮有有效移动，阻止默认滚动行为
    if (deltaY !== 0) {
        e.preventDefault()
    }
    
    // 节流处理：如果上次触发还在节流时间内，则忽略本次事件
    if (wheelThrottleTimer !== null) return
    
    // 设置节流定时器
    wheelThrottleTimer = setTimeout(() => {
        wheelThrottleTimer = null
    }, config.wheelThrottleDelay)
    
    // 向下滚动 → 下一页，向上滚动 → 上一页
    if (deltaY > 0) {
        toNextPage()
    } else if (deltaY < 0) {
        toPreviousPage()
    }
}

// ============ 主初始化函数 ============
/**
 * 主初始化函数 - 在 DOM 加载完成后调用
 */
function initFlipMode() {
    try {
        // 1. 初始化 DOM 缓存
        initDOMCache()
        
        // 2. 初始化工具栏事件监听器
        initToolbarListeners()
        
        // 3. 初始化滑动事件监听器
        initSlideListeners()
        
        // 4. 初始化 FlipMainArea 事件监听器
        initFlipMainAreaListeners()
        
        // 5. 初始化共享 WebSocket 并连接（在线书籍时）
        if (typeof window.ComiGoWS !== 'undefined') {
            window.ComiGoWS.init({
                pageType: 'flip',
                getBookId: () => book?.id,
                getWsConfig: () => ({
                    maxReconnectAttempts: Alpine.store('flip').websocketMaxReconnect,
                    reconnectInterval: Alpine.store('flip').websocketReconnectInterval,
                }),
                isDebug: () => Alpine.store('global').debugMode,
                onMessage(msg) {
                    if (msg.type === 'flip_mode_sync_page' && msg.tab_id !== window.ComiGoWS.getTabId()) {
                        try {
                            const data = JSON.parse(msg.data_string || '{}')
                            if (Alpine.store('global').syncPageByWS && data.book_id === book.id) {
                                jumpPageNum(data.now_page_num, false)
                            }
                        } catch (e) {
                            console.error('WebSocket 翻页同步数据解析失败:', e)
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
        
        if (Alpine.store('global').debugMode) {
            console.log('翻页模式初始化完成')
        }
    } catch (error) {
        console.error('翻页模式初始化失败:', error)
    }
}

// DOM 加载完成后初始化
document.addEventListener('DOMContentLoaded', initFlipMode)
