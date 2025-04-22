//此文件静态导入，不需要编译
//https://www.runoob.com/js/js-strict.html
'use strict'

//设置初始值
const book = JSON.parse(document.getElementById('NowBook').textContent)
const images = book.pages.images
Alpine.store('flip').allPageNum = parseInt(book.page_count)
const ALL_PAGE_NUM = Alpine.store('flip').allPageNum
// // 打印调试信息
// if (Alpine.store('global').debugMode) {
// 	//console.log(book);
// 	// const globalState = JSON.parse(document.getElementById('GlobalState').textContent);
// 	console.log('book.page_count:', book.page_count)
// 	console.log('images.length:', images.length)
// }

// 滑动相关变量
let touchStartX = 0
let touchEndX = 0
let isSwiping = false
let currentTranslate = 0
let startTime = 0
let animationID = 0
const sliderContainer = document.getElementById('slider-container')
const slider = document.getElementById('slider')
const leftSlide = document.getElementById('left-slide')
const middleSlide = document.getElementById('middle-slide')
const rightSlide = document.getElementById('right-slide')
const threshold = 100 // 滑动阈值，超过这个值才会触发翻页
const swipeTimeout = 300 // 滑动超时时间（毫秒）

// 设置图片资源，预加载等
// 需要 HTTP 响应头中允许缓存（没有使用 Cache-Control: no-cache），也就是 gin 不能用htmx/router/server.go 里面的noCache 中间件。
// 在预加载用到的图片资源URL
let preloadedImages = new Set()

//首次加载时
setImageSrc()

// === 新增：页面加载时读取本地存储的页码 ===
loadPageNumFromLocalStorage();
// ====================================

// 加载图片资源
function setImageSrc() {
	const nowPageNum = Alpine.store('flip').nowPageNum
	if (nowPageNum === 0 && nowPageNum >= ALL_PAGE_NUM) {
		console.log('setImageSrc: nowPageNum is 0 or out of range', nowPageNum)
		return
	}
	// console.log("setImageSrc: nowPageNum", nowPageNum);
	// console.log("setImageSrc: NowImage", images[nowPageNum - 1].url);
	// console.log("setImageSrc: NowImage+1=", images[nowPageNum].url);
	// 加载当前图片
	document.getElementById('Single-NowImage').src =
		images[nowPageNum - 1].url
	if (!Alpine.store('flip').mangaMode) {
		document.getElementById('Double-NowImage-Left').src = images[nowPageNum - 1].url;
	} else {
		document.getElementById('Double-NowImage-Right').src = images[nowPageNum - 1].url;
	}

	preloadedImages.add(images[nowPageNum - 1].url)
	// 更新滑动容器图片
	updateSliderImages(nowPageNum, images)

	// 为双页模式，加载下一张图片。
	// 因为用户有可能随时切换到双页模式，所以单页模式也预加载图片（尽管看不到）
	if (nowPageNum < ALL_PAGE_NUM) {
		if (Alpine.store('flip').mangaMode) {
			document.getElementById('Double-NowImage-Left').src = images[nowPageNum].url;
		} else {
			document.getElementById('Double-NowImage-Right').src = images[nowPageNum].url;
		}
		preloadedImages.add(images[nowPageNum].url)
	}

	// 预加载前一张和后十张图片
	const preloadRange = 10 // 预加载范围，可以根据需要调整
	for (let i = nowPageNum - 2; i <= nowPageNum + preloadRange; i++) {
		if (i >= 0 && i < ALL_PAGE_NUM) {
			const imgUrl = images[i].url
			if (!preloadedImages.has(imgUrl)) {
				let img = new Image()
				img.src = imgUrl
				preloadedImages.add(imgUrl)
			}
		}
	}
}

// === 保存当前页码到本地存储 ===
function savePageNumToLocalStorage() {
	if (!Alpine.store('flip').saveReadingProgress) {
		return;
	}
    try {
        const key = `pageNum_${book.id}`;
        const nowPageNum = Alpine.store('flip').nowPageNum;
        localStorage.setItem(key, nowPageNum);
    } catch (e) {
        console.error("Error saving page number to localStorage:", e);
    }
}
// ====================================

// 更新滑动容器图片
function updateSliderImages(nowPageNum, images) {
	// 根据阅读方向设置滑动元素的位置
	const prevSlideElement = document.getElementById('left-slide')
	const nextSlideElement = document.getElementById('right-slide')
	if (Alpine.store('flip').mangaMode) {
		// 日漫模式：prev在右侧，next在左侧
		prevSlideElement.style.transform = 'translateX(100%)'
		nextSlideElement.style.transform = 'translateX(-100%)'
	} else {
		// 美漫模式：prev在左侧，next在右侧
		prevSlideElement.style.transform = 'translateX(-100%)'
		nextSlideElement.style.transform = 'translateX(100%)'
	}
	// ------------ 单页模式设置 ------------
	if (!Alpine.store('flip').doublePageMode) {
		// 添加前一张图片（如果存在）
		if (nowPageNum > 1) {
			const prevImg = document.createElement('img')
			prevImg.src = images[nowPageNum - 2].url
			prevImg.className = 'object-contain m-0 max-w-full max-h-full h-full'
			prevImg.draggable = false
			leftSlide.innerHTML = ''
			leftSlide.appendChild(prevImg)
		} else {
			leftSlide.innerHTML = ''
		}
		// // 更新当前图片 (确保当前图片也在这里更新，以防万一)
		const currentImgElement = document.getElementById('Single-NowImage')
		if (currentImgElement && nowPageNum >= 1 && nowPageNum <= ALL_PAGE_NUM) {
			currentImgElement.src = images[nowPageNum - 1].url
		}
		// 添加后一张图片（如果存在）
		if (nowPageNum < ALL_PAGE_NUM) {
			const nextImg = document.createElement('img')
			nextImg.src = images[nowPageNum].url
			nextImg.className = 'object-contain m-0 max-w-full max-h-full h-full'
			nextImg.draggable = false
			rightSlide.innerHTML = ''
			rightSlide.appendChild(nextImg)
		} else {
			rightSlide.innerHTML = ''
		}
	}
	// ------------ 双页模式设置 ------------
	if (Alpine.store('flip').doublePageMode) {
		// 添加双页模式前一屏图片（如果存在）
		if (nowPageNum  === 2) {
			const prevImg = document.createElement('img')
			prevImg.src = images[nowPageNum - 2].url
			prevImg.className = 'object-contain m-0 max-w-full max-h-full h-full'
			prevImg.draggable = false
			leftSlide.innerHTML = ''
			leftSlide.appendChild(prevImg)
		}
		if (nowPageNum >= 3) {
			const prevImg_1 = document.createElement('img')
			prevImg_1.src = images[nowPageNum - 2].url
			prevImg_1.className = 'select-none object-contain m-0 max-w-1/2 w-auto max-h-screen grow-0'
			prevImg_1.draggable = false
			const prevImg_2 = document.createElement('img')
			prevImg_2.src = images[nowPageNum - 3].url
			prevImg_2.className = 'select-none object-contain m-0 max-w-1/2 w-auto max-h-screen grow-0'
			prevImg_2.draggable = false
			leftSlide.innerHTML = ''
			if (Alpine.store('flip').mangaMode) {
				leftSlide.appendChild(prevImg_1)
				leftSlide.appendChild(prevImg_2)
			} else {
				leftSlide.appendChild(prevImg_2)
				leftSlide.appendChild(prevImg_1)
			}
		}
		if (nowPageNum <= 1) {
			leftSlide.innerHTML = ''
		}
		// 添加后一屏图片（如果存在）
		if (nowPageNum  === ALL_PAGE_NUM - 3) {
			const nextImg = document.createElement('img')
			nextImg.src = images[nowPageNum - 2].url
			nextImg.className = 'object-contain m-0 max-w-full max-h-full h-full'
			nextImg.draggable = false
			rightSlide.innerHTML = ''
			rightSlide.appendChild(nextImg)
		}
		if (nowPageNum < ALL_PAGE_NUM - 3) {
			const nextImg_1 = document.createElement('img')
			nextImg_1.src = images[nowPageNum + 1].url
			nextImg_1.className = 'select-none object-contain m-0 max-w-1/2 w-auto max-h-screen grow-0'
			nextImg_1.draggable = false
			const nextImg_2 = document.createElement('img')
			nextImg_2.src = images[nowPageNum + 2].url
			nextImg_2.className = 'select-none object-contain m-0 max-w-1/2 w-auto max-h-screen grow-0'
			nextImg_2.draggable = false
			rightSlide.innerHTML = ''
			if (Alpine.store('flip').mangaMode) {
				rightSlide.appendChild(nextImg_2)
				rightSlide.appendChild(nextImg_1)
			} else {
				rightSlide.appendChild(nextImg_1)
				rightSlide.appendChild(nextImg_2)
			}
		}
		if (nowPageNum === ALL_PAGE_NUM - 1) {
			rightSlide.innerHTML = ''
		}
	}

	// 确保滑动容器在更新图片后回到初始位置 (没有动画)
	slider.style.transition = 'none' // 暂时禁用过渡效果，防止闪烁
	slider.style.transform = 'translateX(0)'
	// 强制浏览器重新计算样式，确保 `transition = 'none'` 生效
	slider.offsetHeight // 读取offsetHeight可以触发重排
	slider.style.transition = '' // 恢复过渡效果
	resetSlider() // 清理状态 (currentTranslate = 0, cancel animation)
}

// 重置滑动状态
function resetSlider() {
	cancelAnimationFrame(animationID)
	// 不再立即设置 transform
	currentTranslate = 0
}

// 触摸开始事件处理
function touchStart(e) {
	// 根据swipeTurn的值决定是否启用滑动翻页
	if (!Alpine.store('flip').swipeTurn)
		return
	startTime = new Date().getTime()
	isSwiping = true
	touchStartX = e.type === 'touchstart' ? e.touches[0].clientX : e.clientX

	// 停止任何正在进行的动画
	cancelAnimationFrame(animationID)
}

// 如果在第一页或最后一页尝试向前翻或向后翻，阻止默认滚动
function shouldBlockScroll(diffX) {
	const mangaMode = Alpine.store('flip').mangaMode
	const nowPageNum = Alpine.store('flip').nowPageNum
	// 判断是否应该阻止默认滚动
	let blockScroll = false
	// 如果是第一页尝试向前翻
	if (nowPageNum === 1) {
		// 日漫模式
		if (diffX < 0 && mangaMode) {
			blockScroll = true
		}
		// 美漫模式
		if (diffX > 0 && !mangaMode) {
			blockScroll = true
		}
	}
	// 如果是最后一页尝试向后翻
	if (nowPageNum === ALL_PAGE_NUM) {
		// 日漫模式
		if (diffX > 0 && mangaMode) {
			blockScroll = true
		}
		// 美漫模式
		if (diffX < 0 && !mangaMode) {
			blockScroll = true
		}
	}
	return blockScroll;
}

// 触摸移动事件处理
function touchMove(e) {
	if (!isSwiping)
		return
	if (!Alpine.store('flip').swipeTurn)
		return
	const currentX = e.type === 'touchmove' ? e.touches[0].clientX : e.clientX
	const diffX = currentX - touchStartX

	// 在第一页时不能向前翻，在最后一页时不能向后翻
	const nowPageNum = Alpine.store('flip').nowPageNum
	const mangaMode = Alpine.store('flip').mangaMode
	// 设置当前滑动距离
	currentTranslate = diffX
	// 如果在第一页或最后一页尝试向前翻或向后翻，阻止默认滚动
	if (shouldBlockScroll(diffX)) {
		if (diffX < 0) {
			currentTranslate = -30
		} else {
			currentTranslate = 30
		}
	}
	// 应用变换
	slider.style.transform = `translateX(${currentTranslate}px)`
	// 防止页面滚动
	if (Math.abs(diffX) > 10) {
		e.preventDefault()
	}
}

// 触摸结束事件处理
function touchEnd(e) {
	// 根据swipeTurn的值决定是否滑动翻页
	if (!isSwiping || !Alpine.store('flip').swipeTurn)
		return
	const nowPageNum = Alpine.store('flip').nowPageNum
	const mangaMode = Alpine.store('flip').mangaMode
	// 取消滑动状态
	isSwiping = false
	const endTime = new Date().getTime()
	const timeElapsed = endTime - startTime
	touchEndX = e.type === 'touchend' ? e.changedTouches[0].clientX : e.clientX
	const diffX = touchEndX - touchStartX

	// 判断是否应该翻页（基于滑动距离和速度）
	const isQuickSwipe = timeElapsed < swipeTimeout && Math.abs(diffX) > 50
	// 用于记录滑动方向
	let direction = null
	if (diffX < -threshold || (isQuickSwipe && diffX < 0)) {
		// 向左滑动
		direction = 'left'
	} else if (diffX > threshold || (isQuickSwipe && diffX > 0)) {
		// 向右滑动
		direction = 'right'
	}
	// if (Alpine.store('global').debugMode) {
	// 	console.log('mangaMode:', mangaMode)
	// 	console.log('touchEnd: diffX', diffX)
	// 	console.log('touchEnd: direction', direction)
	// 	console.log('nowPageNum:', nowPageNum)
	// 	console.log('ALL_PAGE_NUM:', ALL_PAGE_NUM)
	// }
	// 如果在第一页或最后一页尝试向前翻或向后翻，阻止默认滚动
	if (shouldBlockScroll(diffX) || direction === null) {
		// 没有足够的滑动距离或在边界，回到原始位置
		animateReset()
		return
	}
	// 如果确定了滑动方向，执行滑动动画及后续翻页
	animateSlide(direction)
}

// 修改 animateSlideOut 为 animateSlide，处理滑动动画和翻页逻辑
function animateSlide(direction) {
	const width = window.innerWidth

	// 根据滑动方向确定目标位置
	let targetPosition = direction === 'left' ? -width : width
	// 左滑是下一页（移到左侧），右滑是上一页（移到右侧）
	const mangaMode = Alpine.store('flip').mangaMode
	let startTime = null
	const duration = 300 // 动画持续时间，单位毫秒
	const startPosition = currentTranslate // 记录动画开始时的位置

	function animate(timestamp) {
		if (!startTime) startTime = timestamp
		const elapsed = timestamp - startTime
		const progress = Math.min(elapsed / duration, 1)

		// 使用缓动函数使动画更平滑
		const easeProgress = easeOutCubic(progress)

		// 计算当前位置（从startPosition到targetPosition的过渡）
		const position =
			startPosition + (targetPosition - startPosition) * easeProgress

		// 应用变换
		slider.style.transform = `translateX(${position}px)`

		if (progress < 1) {
			animationID = requestAnimationFrame(animate)
		} else {
			// 动画完成后执行翻页逻辑
			// 1. 确定调用哪个翻页函数
			let flipFunction
			if (mangaMode) {
				flipFunction = direction === 'left' ? toPreviousPage : toNextPage
			} else {
				flipFunction = direction === 'left' ? toNextPage : toPreviousPage
			}
			// 2. 执行翻页 (这会触发页面号码更新和 setImageSrc -> updateSliderImages)
			if (flipFunction) {
				// updateSliderImages 会负责加载新内容并将 slider transform 重置为 translateX(0)
				flipFunction()
			} else {
				// 以防万一没有确定翻页函数，动画重置回去
				animateReset()
			}
		}
	}

	// 启动动画
	animationID = requestAnimationFrame(animate)
}

// 动画回到原始位置
function animateReset() {
	let startTime = null
	const duration = 400 // 动画持续时间，单位毫秒
	const startPosition = currentTranslate

	function animate(timestamp) {
		if (!startTime) startTime = timestamp
		const elapsed = timestamp - startTime
		const progress = Math.min(elapsed / duration, 1)

		// 使用缓动函数使动画更平滑
		const easeProgress = easeOutCubic(progress)

		// 计算当前位置（从startPosition到0的过渡）
		const position = startPosition * (1 - easeProgress)

		// 应用变换
		slider.style.transform = `translateX(${position}px)`

		if (progress < 1) {
			animationID = requestAnimationFrame(animate)
		} else {
			// 动画完成后，确保transform为0并清理状态
			if (slider) {
				slider.style.transform = 'translateX(0)'
			}
			resetSlider() // 清理状态 (currentTranslate = 0, cancel animation)
		}
	}

	// 启动动画
	animationID = requestAnimationFrame(animate)
}

// 缓动函数 - 使动画更自然
function easeOutCubic(x) {
	return 1 - Math.pow(1 - x, 3)
}

// 为滑动容器添加事件监听器
document.addEventListener('DOMContentLoaded', function () {
	// 触摸事件（移动设备）
	// 设置初始值
	sliderContainer.addEventListener('touchstart', touchStart)
	// 移动中
	sliderContainer.addEventListener('touchmove', touchMove, {passive: false})
	// 移动结束
	sliderContainer.addEventListener('touchend', touchEnd)

	// 鼠标事件（PC）
	// 设置初始值
	sliderContainer.addEventListener('mousedown', touchStart)
	// 移动中
	sliderContainer.addEventListener('mousemove', touchMove)
	// 移动结束
	sliderContainer.addEventListener('mouseup', touchEnd)
	sliderContainer.addEventListener('mouseleave', touchEnd)

	// 初始化滑动容器中的图片
	const nowPageNum = Alpine.store('flip').nowPageNum
	const images = book.pages.images
	updateSliderImages(nowPageNum, images)
})

// 监听mangaMode值的变化，动态更新滑动方向
document.addEventListener('alpine:initialized', () => {
	if (window.Alpine) {
		// 当mangaMode值变化时更新滑动方向
		Alpine.effect(() => {
			if (nowPageNum > 0) {
				updateSliderImages(Alpine.store('flip').nowPageNum, book.pages.images)
			}
		})
	}
})

//翻页函数，加页或减页
function addPageNum(n = 1) {
	//防止数字转换为字符串
	let num = parseInt(n)
	let nowPageNum = parseInt(Alpine.store('flip').nowPageNum)
	//    无法继续翻
	if (nowPageNum + num > ALL_PAGE_NUM) {
		showToast(i18next.t('hint_last_page'), 'error')
		return
	}
	if (nowPageNum + num < 1) {
		showToast(i18next.t('hint_first_page'), 'error')
		return
	}
	// 翻页
	Alpine.store('flip').nowPageNum = nowPageNum + num
	setImageSrc()
	// 设置标签页标题
	setTitle();
	// 通过ws通道发送翻页数据
	if (Alpine.store('global').syncPageByWS === true) {
		sendFlipData() // 发送翻页数据
	}
    // === 调用保存函数 ===
    savePageNumToLocalStorage();
    // =======================
}

// === 从本地存储加载页码并跳转 ===
function loadPageNumFromLocalStorage() {
	if (!Alpine.store('flip').saveReadingProgress) {
		return;
	}
    try {
        const key = `pageNum_${book.id}`;
        const savedPageNum = localStorage.getItem(key);
        if (savedPageNum !== null && !isNaN(parseInt(savedPageNum))) {
            const pageNum = parseInt(savedPageNum);
            // 确保页码在有效范围内
            if (pageNum > 0 && pageNum <= ALL_PAGE_NUM && pageNum !== Alpine.store('flip').nowPageNum) {
                console.log(`加载到本地存储的页码: ${pageNum}`);
                jumpPageNum(pageNum); // 使用跳转函数更新页面
            }
        }
    } catch (e) {
        console.error("Error loading page number from localStorage:", e);
    }
}
// ====================================

//翻页函数，跳转到指定页
function jumpPageNum(jumpNum) {
	let num = parseInt(jumpNum)
	if (num <= 0 || num > ALL_PAGE_NUM) {
		alert(i18next.t('hintPageNumOutOfRange'))
		return
	}
	Alpine.store('flip').nowPageNum = num
	setImageSrc()
    // === 调用保存页数函数 ===
    savePageNumToLocalStorage();
    // =======================
}

// 翻页函数，下一页
function toNextPage() {
	let doublePageMode = Alpine.store('flip').doublePageMode === true
	let nowPageNum = parseInt(Alpine.store('flip').nowPageNum)
	// 单页模式
	if (!doublePageMode) {
		if (nowPageNum <= ALL_PAGE_NUM) {
			addPageNum(1)
		}
	}
	//双页模式
	if (doublePageMode) {
		if (nowPageNum < ALL_PAGE_NUM - 1) {
			addPageNum(2)
		} else {
			addPageNum(1)
		}
	}
}

// 翻页函数，前一页
function toPreviousPage() {
	let doublePageMode = Alpine.store('flip').doublePageMode === true
	let nowPageNum = parseInt(Alpine.store('flip').nowPageNum)
	//错误值,第0或第1页。
	if (nowPageNum <= 1) {
		showToast(i18next.t('hint_first_page'), 'error')
		return
	}
	//简单合并模式
	if (doublePageMode) {
		if (nowPageNum - 2 > 0) {
			addPageNum(-2)
		} else {
			addPageNum(-1)
		}
	} else {
		addPageNum(-1)
	}
}

//隐藏工具栏的工具函数
// https://www.runoob.com/js/js-htmldom-events.html
let hideTimeout
let header = document.getElementById('header')
let range = document.getElementById('StepsRangeArea')

// 显示工具栏
function showToolbar() {
	if (Alpine.store('flip').autoHideToolbar) {
		header.style.opacity = '0.9'
		range.style.opacity = '0.9'
		header.style.transform = 'translateY(0)'
		range.style.transform = 'translateY(0)'
	} else {
		header.style.opacity = '1'
		range.style.opacity = '1'
		header.style.transform = 'translateY(0)'
		range.style.transform = 'translateY(0)'
	}
}

// 隐藏工具栏
function hideToolbar() {
	if (Alpine.store('flip').autoHideToolbar) {
		header.style.opacity = '0'
		range.style.opacity = '0'
		header.style.transform = 'translateY(-100%)'
		range.style.transform = 'translateY(100%)'
	}
}

function startHideTimer() {
	hideTimeout = setTimeout(hideToolbar, 1500) // 1.5 seconds
}

if (Alpine.store('flip').autoHideToolbar) {
	startHideTimer()
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

// 翻页模式功能优化：不隐藏工具栏的时候。点击设置区域，自动漫画区域居中。
function scrollToMangaMain() {
	if (!Alpine.store('flip').autoHideToolbar) {
		// 1. 获取 manga_area 元素
		const mangaMains = document.getElementsByClassName('manga_area')
		// 2. 将 manga_area 顶部对齐到浏览器可见区域顶部
		//        这样它的高度（100vh）就能正好占满整个可见区域
		for (let i = 0; i < mangaMains.length; i++) {
			const mangaMain = mangaMains[i]
			mangaMain.scrollIntoView({
				behavior: 'smooth', // 平滑滚动
				block: 'start', // 与可视区顶部对齐
			})
		}
	}
}

//获取鼠标位置,决定是否打开设置面板
function onMouseClick(e) {
	// 如果正在滑动，则不处理点击事件
	if (isSwiping || Math.abs(currentTranslate) > 10) {
		return
	}

	let clickX = e.x //获取鼠标的X坐标（鼠标与屏幕左侧的距离,单位为px）
	//浏览器的视口宽,不包括工具栏和滚动条:
	let innerWidth = window.innerWidth
	let inSetArea = getInSetArea(e)
	if (inSetArea) {
		//获取ID为 OpenSettingButton的元素，然后模拟点击
		document.getElementById('OpenSettingButton').click()
		scrollToMangaMain()
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
	if (inSetArea) {
		showToolbar()
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
	if (inSetArea) {
		e.currentTarget.style.cursor = 'url(/images/SettingsOutline.png), pointer'
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
				Alpine.store('flip').nowPageNum === 1
			) {
				//右边翻下一页,且目前是第一页的时候,左边的鼠标指针,设置为禁止翻页
				e.currentTarget.style.cursor =
					'url(/images/Prohibited28Filled.png), pointer'
			} else if (
				Alpine.store('flip').mangaMode &&
				Alpine.store('flip').nowPageNum === ALL_PAGE_NUM
			) {
				//左边翻下一页,且目前是最后一页的时候,左边的鼠标指针,设置为禁止翻页
				e.currentTarget.style.cursor =
					'url(/images/Prohibited28Filled.png), pointer'
			} else {
				//正常情况下,左边是向左的箭头
				e.currentTarget.style.cursor = 'url(/images/ArrowLeft.png), pointer'
			}
		} else {
			//设置右边的鼠标指针
			if (
				!Alpine.store('flip').mangaMode &&
				Alpine.store('flip').nowPageNum === ALL_PAGE_NUM
			) {
				//右边翻下一页,且目前是最后页的时候,右边的鼠标指针,设置为禁止翻页
				e.currentTarget.style.cursor =
					'url(/images/Prohibited28Filled.png), pointer'
			} else if (
				Alpine.store('flip').mangaMode &&
				Alpine.store('flip').nowPageNum === 1
			) {
				//左边翻下一页,且目前是第一页的时候,右边的鼠标指针,设置为禁止翻页
				e.currentTarget.style.cursor =
					'url(/images/Prohibited28Filled.png), pointer'
			} else {
				//正常情况下,右边是向右的箭头
				e.currentTarget.style.cursor = 'url(/images/ArrowRight.png), pointer'
			}
		}
	}
}

// 获取两个元素的边界信息
function getElementsRect() {
	return {
		rect1_header: header.getBoundingClientRect(),
		rect2_range: range.getBoundingClientRect(),
		rect3_sort_dropdown: document
			.getElementById('ReSortDropdownMenu')
			.getBoundingClientRect(),
		rect4_dropdown_quick_jump: document
			.getElementById('QuickJumpDropdown')
			.getBoundingClientRect(),
		rect5_steps_range_area: document
			.getElementById('StepsRangeArea')
			.getBoundingClientRect(),
	}
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
	// 因为header需要收起来，所以不能用left、right、top、bottom判断y是否在header的范围内
	if (Alpine.store('flip').autoHideToolbar) {
		// 判断鼠标是否在元素 1 范围内(Header)。。
		inInElement1 = y <= rect1_header.height - 20
		// 判断鼠标是否在元素 2 范围内(导航条)。因为header可能隐藏，所以不能直接用left、right、top、bottom判断y是否在header的范围内。
		inInElement2 = y >= window.innerHeight - rect2_range.height - 20
	}
	// 如果工具栏不自动隐藏，用left、right、top、bottom判断y是否在header的范围内
	if (!Alpine.store('flip').autoHideToolbar) {
		// 判断鼠标是否在元素 1 范围内(Header)
		inInElement1 =
			x >= rect1_header.left &&
			x <= rect1_header.right &&
			y >= rect1_header.top &&
			y <= rect1_header.bottom
		// 判断鼠标是否在元素 2 范围内(导航条)
		inInElement2 =
			x >= rect2_range.left &&
			x <= rect2_range.right &&
			y >= rect2_range.top &&
			y <= rect2_range.bottom
	}

	// 判断鼠标是否在元素 3 范围内(页面重新排序的下拉菜单。在菜单上面的时候，导航条需要保持显示状态。)
	const inInElement3 =
		x >= rect3_sort_dropdown.left &&
		x <= rect3_sort_dropdown.right &&
		y >= rect3_sort_dropdown.top &&
		y <= rect3_sort_dropdown.bottom
	// 判断鼠标是否在元素 4 范围内(快速跳转的下拉菜单。在菜单上面的时候，导航条需要保持显示状态。)
	const inInElement4 =
		x >= rect4_dropdown_quick_jump.left &&
		x <= rect4_dropdown_quick_jump.right &&
		y >= rect4_dropdown_quick_jump.top &&
		y <= rect4_dropdown_quick_jump.bottom

	// 判断鼠标是否在元素 5 范围内(翻页条)
	const inInElement5 =
		x >= rect5_steps_range_area.left &&
		x <= rect5_steps_range_area.right &&
		y >= rect5_steps_range_area.top &&
		y <= rect5_steps_range_area.bottom

	//鼠标在设置区域
	let inSetArea = getInSetArea(event)
	if (inSetArea) {
		showToolbar()
	}
	// '鼠标不在任何一个元素范围内'
	if (!inInElement1 && !inInElement2 && !inInElement3) {
		if (inSetArea === false) {
			hideToolbar()
		}
	}
	// '鼠标在元素范围内'
	if (
		inInElement1 ||
		inInElement2 ||
		inInElement3 ||
		inInElement4 ||
		inInElement5 ||
		inSetArea
	) {
		//console.log('鼠标在元素范围内');
		showToolbar()
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

//获取ID为 mouseMoveArea 的元素
let mouseMoveArea = document.getElementById('mouseMoveArea')
// 鼠标移动时触发移动事件
mouseMoveArea.addEventListener('mousemove', onMouseMove)
//点击的时候触发点击事件
mouseMoveArea.addEventListener('click', onMouseClick)
//离开的时候触发离开事件
mouseMoveArea.addEventListener('mouseleave', onMouseLeave)

// Websocket 连接和消息处理
// https://www.ruanyifeng.com/blog/2017/05/websocket.html
// https://developer.mozilla.org/zh-CN/docs/Web/API/WebSocket

// 定义WebSocket变量和重连参数
let socket
let reconnectAttempts = 0
const maxReconnectAttempts = 200
const reconnectInterval = 3000 // 每次重连间隔3秒

// 用户ID和令牌，假设已在其他地方定义
const userID = Alpine.store('global').userID
// 假设token是一个有效的令牌 TODO:使用真正的令牌
const token = 'your_token'

// 翻页数据，假设已在其他地方定义
const flip_data = {
	book_id: book.id,
	now_page_num: Alpine.store('flip').nowPageNum,
	need_double_page_mode: false,
}

// 建立WebSocket连接的函数
function connectWebSocket() {
	// 根据当前协议选择ws或wss
	const wsProtocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://'
	const wsUrl = wsProtocol + window.location.host + '/api/ws'
	socket = new WebSocket(wsUrl)

	// 连接打开时的回调
	socket.onopen = function () {
		console.log('WebSocket连接已建立')
		reconnectAttempts = 0 // 重置重连次数
	}

	// 收到消息时的回调
	socket.onmessage = function (event) {
		const message = JSON.parse(event.data)
		handleMessage(message) // 调用处理函数
	}

	// 连接关闭时的回调
	socket.onclose = function () {
		console.log('WebSocket连接已关闭')
		attemptReconnect() // 尝试重连
	}

	// 发生错误时的回调
	socket.onerror = function (error) {
		console.log('WebSocket发生错误：', error)
		socket.close() // 关闭连接以触发重连
	}
}

// 处理收到的翻页消息
function handleMessage(message) {
	// console.log("收到消息：", message);
	// console.log("My userID：" + userID);
	// console.log("Remote userID：" + message.user_id);
	// 根据消息类型进行处理
	if (message.type === 'flip_mode_sync_page' && message.user_id !== userID) {
		// 解析翻页数据
		const data = JSON.parse(message.data_string)
		if (Alpine.store('global').syncPageByWS && data.book_id === book.id) {
			//console.log("同步页数：", data);
			// 更新翻页数据
			flip_data.now_page_num = data.now_page_num
			// 更新页面
			jumpPageNum(data.now_page_num)
		}
	} else if (message.type === 'heartbeat') {
		console.log('收到心跳消息')
	} else {
		//console.log("不处理此消息"+message);
	}
}

// 发送翻页数据到服务器
function sendFlipData() {
	flip_data.now_page_num = Alpine.store('flip').nowPageNum
	const flipMsg = {
		type: 'flip_mode_sync_page', // 或 "heartbeat"
		status_code: 200,
		user_id: userID,
		token: token,
		detail: '翻页模式，发送数据',
		data_string: JSON.stringify(flip_data),
	}
	if (socket.readyState === WebSocket.OPEN) {
		socket.send(JSON.stringify(flipMsg))
		//console.log("已发送翻页数据"+JSON.stringify(flipMsg));
		//console.log("已发送翻页数据");
	} else {
		console.log('WebSocket未连接，无法发送消息')
	}
}

// 尝试重连函数
function attemptReconnect() {
	if (reconnectAttempts < maxReconnectAttempts) {
		reconnectAttempts++
		console.log(`第 ${reconnectAttempts} 次重连...`)
		setTimeout(() => {
			connectWebSocket()
		}, reconnectInterval)
	} else {
		console.log('已达到最大重连次数，停止重连')
	}
}

// 页面加载完成后建立WebSocket连接
window.onload = function () {
	connectWebSocket()
}

// 设置标签页标题
function setTitle(name) {
	let numStr = ''
	if (Alpine.store('flip').showPageNum){
		numStr = ` ${Alpine.store('flip').nowPageNum}/${Alpine.store('flip').allPageNum} `
	}
	//简化标题
	if(Alpine.store('shelf').simplifyTitle){
		document.title = `${numStr}${shortName(book.title)}`;
	}else {
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
	const runes          = Array.from(shortTitle);
	const originalRunes  = Array.from(title);

	/* —— 与 Go 的长度判断逻辑保持完全一致 —— */

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
