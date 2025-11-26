//此文件静态导入，不需要编译
//https://www.runoob.com/js/js-strict.html
'use strict'

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

let headerArea = document.getElementById('header');
// 显示工具栏
headerArea.addEventListener('mouseover', showToolbar);
// 隐藏工具栏
headerArea.addEventListener('mouseout', hideToolbar);

// 初始化：如果autohidetoolbar为真,则自动隐藏工具栏
if (Alpine.store('flip').autoHideToolbar) {
    setTimeout(hideToolbar, 1000)
}

// Base64编码静态资源图片（鼠标图标）：
// base64 -i SettingsOutline.svg ，然后// 把下面这行换成输出的
const SettingsOutlineBase64 = 'iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAACXBIWXMAAAsSAAALEgHS3X78AAAKZklEQVRYhZVXbUxUVxp+zrl3LsPMIG7Dh4pVR0atgBhwt9WuaMCqRFM1McTYH0sbm25T7IYsMlhp7WZXy0JbowKxAqXNxjRNFdpqu+vGD1jsRBQH1/GjOiDMgEO7EpRBBubj3vPuD2emgN1s+v659yTnvO9zzvu8X8xqtWLbtm0YGRmBJEkQQsBoNOKrr75Ce3s7jEYjiAiBQABZWVlsx44dC3U63RJVVYkxBgAgIsiyjLGxsesVFRVdAwMDZDKZIIQAYwyBQABz585FXV0dAEBVVciyjI8//hgypggRQVEUOBwOnDt3DrIsQ1VVlpiYSPv27dt06dKlXTdu3FAVReFEBAARI2Lp0qXSunXrKl9//fVvx8bGGGOMGGMQQsBisUw1BQBPAggLlyQJAEhRFKaqqqiqqlr64MGD37/yyiu/VVWVMcYwEQARkSzL7NSpUz9WVVX17dy583pMTAwnIgoEAkyWZQAQ/xcAEXEhhFBVFQAQDAY5AMydOzfY1taWrKoqy8zMVFetWsX9fj8AQK/Xo62tTXM4HLoLFy6kPPfccz4A8Pv9jHMuAFBYH58KIgogciNFUUhRlIzi4uLY559//v7777/vXrZsGdLT03/T2NioA4Bt27bxPXv28AhIWZZRUVEBh8MBl8sVU1hYmLV48WK3Xq/XCgoK5pnN5kQhxLiqqrfCLzEZQNhPTKfT0YwZM/IvX768x+VymdasWdOxceNGmyRJ90ZGRjJdLtcsAJg/fz4DACF+uozZbGYA4HK5ZgcCgayGhoaR1NRUc3d39yqbzbbUbDY/dLvdf5k9e/YZnU7HADz2X1lZGex2O29tbUV/f/+i48eP/91oNIYAkMFg0LZs2XLvyJEjl7Zv3+6SJIkA0JUrV4iISNM00jSNiIjsdjsBIM45FRQUuI8cOWLPz88fjI2N1cK61MbGxua+vr4FmqaBiHhDQ8NjAB0dHayzsxMDAwO/y8nJeQCAzGazZjAYBABijFEYMZWVlRERkRCCIhL5Ly8vj+6NfI1Go7BYLBoAysnJ+eHq1at/8Hg8RiJCfX09g9Vq5R0dHejq6so9cOCADQAlJCRo3d3ddOvWLSouLhYLFy4UL7zwgrh48eIkg0KI6H/kJTo6Oig/P1+kpaWJXbt2iTt37lBPTw8lJSVpAOi9994763Q6VxIRGhoaOHbt2gWPx2O4fPny4YULF/oAUEVFhSAiUlWViIhCoVDUgKZppKrqEy+gqmp0T2Q9UUdVVZUAQAsWLHh09uzZSlVVYxsbG8FlWUZsbOyGzz77LNfpdBoWL14sdu7cySLk0jQNsiyDMQZN0wAAkiSBMQav1wuv1wvGGMJ5A2H/QpKk6H4AKCoqYhkZGaKrq8vU1NS0ZnBwcO3o6Ci4yWSa3dfXt+LYsWMLAKC8vJybTCaEQiFIkgRJkkBEYIyBMQbOOZqbm7FmzRqkpqYiNTUVeXl5aG5uBuc8uicCQpIkhEIhGAwGvPXWWxwAvvjii3SXy/Xr+/fvp/BAIBACoIVDA59++qlwu93Q6XTRMIvkCM45rFYrtm7divPnz8Pv98Pv96OlpQVbt25FWVlZ1HikTgghoNPp0Nvbi4aGBhFJXOKxchWlpaVwuVwZH3300enExMQAAEpOTqa2trZJPiciOnHiBAGguLg4qq6uJo/HQx6Ph6qrq8lkMhEAampqivo+womWlhZKTEwkADRz5kz/wYMHv37w4EFaXV0dYLVaud1ux7Vr1xafPHnyy2efffYhAFqyZInw+/2TGL5+/XoCQNXV1VEjEXCHDh0iAJSfnz/pzPj4OKWnpwsAtGLFiqHm5ubjvb29i8JRwDhjTACQhoaGvk9LS9v39ttvn5FlWQwODjKv1wsA4JzD6/Wis7MTRqMRL774IqbK5s2bYTAYYLfb4fV6wTkHAHi9XgwNDTGdTidKSkq+zM7O/mtKSsodABIREQce12eTyQS9Xs8YY7qI/36p/Nw5IUS0L2CM6SRJioIDAK5pGp8+fbqWmpqaeePGjT/t378/LxQK8YSEBJo2bVpUSXx8PLKzs+Hz+XDy5MknDJ06dQo+nw/Z2dmIj4+fRGBJkigYDPLKysrNdrt97+DgYBoADQDHnj17MD4+vrS+vr41OTk5AIASExOptbX1F5EwLi6OANCJEyeeIOGZM2coOTmZAFBSUlLw8OHDp30+X0ZdXR2wb9++ZIfDcWDmzJkBAJSbmyt6enomEWli+i0tLY3WBaPRSEajMbq2Wq3/s0643W5au3atBoBmzZrlv3jx4p/feeedZB4MBnWhUEgNBoMCAF599VVmNpsRCoWivorEtRACVVVVaGpqQl5eHvR6PfR6PfLy8tDU1ITKysqov2lCtxQKhTBnzhwUFhZyAAiFQuLx/TQdSkpKcO/evS1vvvnmvwFQenq65vP5ngiziCsmvsrw8DANDw9H11PrxMTzjx49omeeeUYDQEVFRe39/f0bDh06BM4Yw+jo6D9feumlcxaLxXfz5k1eW1tLEXJFcvrP5fv4+HjEx8eDiCbVCSKCqqrR/QBQU1NDt2/f5osWLRrZtGnTmaSkpNZp06aBc875o0ePxk0m09evvfbaFQD44IMPqLe3F3fv3kVpaSktWbKENmzYQFeuXInm+8gzR9zDOQfnHO3t7Vi/fj2lp6dTcXGxcDqd6O7uxocffkgAUFhYaHv66af/oSjKmKZpHGVlZcxut6OlpcV469atP65cufJHAGSxWFSj0TipIWGMUXl5+f8kWllZWZSQkTMGg0GYzWYVAK1evfo/165d2zQ0NMSIiNXX14MTETHGuBDCFwqFThUVFV0yGAxad3e3pGma2Lhx4w+1tbXtBQUFbsYY9u/fj87OzigpI6Sz2+2orKyEJEnYvn373aNHj363efPmPgCit7dXMplMamFh4fm4uLjvp0+fTo/5yaJdseCcM6/X25WRkVFTW1s7ze12G1avXn05Njb20sjIyEBubu66gYGBHTabLcHpdFJ2djaLdMWKoqCrq4sAsOXLlw/s3r37bx6P57u9e/fOKC0tXWGz2ZaZzWZfWlpaQ0pKShcRMYTbc3lCmBEADA0NncvIyBhISEiIOX/+/FBNTU3fU089hba2tllms3nAZrMl9Pb2EgA2MaX29PQQADZnzpwBRVFuvPHGG/8aHh7WSkpKbPPmzUvw+Xx+i8XyvaIoUFU1SvJokx6JW3o85908evQovvnmG8iyLD18+FBzOp0d8+fPDwHA559/rnk8HgoEAgCAmJgYXLhwQQOgzJs3z3/79u2rLpdLAyC9++67biGE22Kx4OWXX/4J8VQAE0Rwznl4gCBFUUhVVTidzphVq1bdk2U52+Fw6K5fvz51NOOyLCMnJ6e/p6fHFAZGANgvGs0iICJxHQwGCQArLS11VFdXN3zyySe/unr16phOp4vehjGGYDAosrKyYvv7+49ZrdbrAFgwGBThChkZzZ6QJwCElSEzMxNjY2MwGo0AQKOjozh9+vS3u3fvvpOZmblICEETy2/43J2DBw/eXb58OUwmE00dz39O/gtuwODKgfux3wAAAABJRU5ErkJggg=='
// base64 -i Prohibited28Filled.png
const Prohibited28FilledBase64 = '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'
// base64 -i ArrowLeft.png
const ArrowLeftBase64 = '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'
// base64 -i ArrowRight.png
const ArrowRightBase64 = 'iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAAAIRlWElmTU0AKgAAAAgABQESAAMAAAABAAEAAAEaAAUAAAABAAAASgEbAAUAAAABAAAAUgEoAAMAAAABAAIAAIdpAAQAAAABAAAAWgAAAAAAAABIAAAAAQAAAEgAAAABAAOgAQADAAAAAQABAACgAgAEAAAAAQAAACCgAwAEAAAAAQAAACAAAAAAX7wP8AAAAAlwSFlzAAALEwAACxMBAJqcGAAAAVlpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IlhNUCBDb3JlIDYuMC4wIj4KICAgPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4KICAgICAgPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIKICAgICAgICAgICAgeG1sbnM6dGlmZj0iaHR0cDovL25zLmFkb2JlLmNvbS90aWZmLzEuMC8iPgogICAgICAgICA8dGlmZjpPcmllbnRhdGlvbj4xPC90aWZmOk9yaWVudGF0aW9uPgogICAgICA8L3JkZjpEZXNjcmlwdGlvbj4KICAgPC9yZGY6UkRGPgo8L3g6eG1wbWV0YT4KGV7hBwAABRVJREFUWAntVktMZEUUfa8/MwINSBMmEQyBQJC/QRKDxiHI/+MC+YcNxJULXejeBOJWElwbPgsMPwMsJEAiIMpEmWDCL7CABEj4LfhNMgJ29+vynOK9TnfTMIAzs7KS26/frap7Tt17quopyv/tv2VA9Zvu/+7XffXVdNV1O8/8/LxVCBEGS4ClwR7BrC0tLfeOeTtkjNLBY/f29j4fGxv7uaen54/Z2dkfnE5nHkiEDAwMmG8d7B4DVYCEnpycfNbe3n6QkJAggoKCRG1trQDwHPrKSALvL58EArPGJBA7PT3dl5OTI/DuhDlgrvz8fJKYN0i8inJIkQEgfXh4eDI2NlaYzWZnZ2fnc4CTiFZUVCT6+/v/0kkEw3cnTbxIxQaBzJGRkan09HRmwLG1tXW2sbFBAjR3cXExM+EhcatMIKUWsLbBomHvwFKvMSr+04mJid/i4+NlCVZXV8/hE0tLS/+AgIskSkpKvElcqwkLBisEz8vLizw/P3+8sLDwCQQWjdQGzAZw1JCQEAtW/bbVauV0FY1PJSMj48HKyooTmVFA0ISx78H9bVVV1TfIyO+tra0XyAYJ+jSKynZ2dlbT3d29VlpaKlJTU0VaWppISUnxMfqTk5NFXFycsNlsIioqihnQACozgG2IUEIsLy9LUaLPXVhY6C1MW6DdQQIx2MfdFBAmkSEDGDUN9JRpxhiRm5vrPDo6chDYjaZpGv8aJDhX89odpegKvqIJOFMgql8yMzOlqCYnJy8QVOzu7or9/f2AhkNI9h8fHxPRDZMNHMRNJJCljzHwAYjJJjWAfxaTyaSi7tIZExOj2e32yxEv/vXZZoYeAKJAC1aUwwltiKmpKTcOrOzQ0NAvUOY1hD1gaM9kKgnsJZzLdakT453BbjI5yevHIMH5JLG2tiaysrLco6OjCkT+0cXFRSKHI6ZqZMBrOoqmaVLWeEpgn84AL0bmAnRJFwlBtJ5ulOFypfAYBEBGCJRBDgoODqYWFH2bSd9NP5jKvegzhD7G47ZECWTgyspKBTr7FWVgCThHGAScSBcVLINsbm6aw8PDFaRKBvGJ7PVCEJwJWkREBAE8DJh6P3C1vLzcVF9fP4ed9h3GPtPHy4UyzY9mZmY6ysrK6CALioDPm4xpFAUFBa7Dw0O5Dal+Gpt+FshtyLh9fX1/wv0BjEr3aI8ZIPNn2dnZAzU1Ne+Defr6+rqn9piAbt+GQ0vBNpRObFn14OBAi4yMtHLlFovFSDvjmqB4U1NT01x1dfXXeH+ql8pHA2JwcNCFE+pJY2Pjl4mJidWnp6dxuh48aTUoYIUqamgFycS2trZ4koUIJUuAC9TcpddcJXhzc/McUv8Vwevq6hjmss5GQOPJ04npgT2EvQmzB7BI+GiPx8fHp5OSkgjsNC6jxcVF4zLSmPbe3l5+oDDtpkBHsIHt/eSKr6zaa4DsQ8DMoaGhad4X6HNsb2//7XUdaxAcwZ9iu32IsebbgkscTODdcJ1JtaM/DQQmo6OjZQY6OjqeU4wIoFVUVEhwjLk7uNdKA/4lMXSQ3FsQ34+4iEiAwB6162mX4PC//O9CBKVWHuKyqoMI93nOh4WFiYaGBm61J0i7sdVeCTjwFQUfFxSrfWdnpw636E9dXV0zuGy+dzgc78JvvnLdylmBf24SXOAZuhdAxunHQ/4N2DH2OL9I8FBZmtfSPCca0FiaOy/ozhP8l6WDMg6X/tpW7s/j3u//Atk+vmGsFsM6AAAAAElFTkSuQmCC'


//设置初始值
const book = JSON.parse(document.getElementById('NowBook').textContent)
// 打印调试信息
if (Alpine.store('global').debugMode) {
    // const globalState = JSON.parse(document.getElementById('GlobalState').textContent);
    console.log('book', book)
    console.log('book.page_count:', book.page_count)
}
const images = book.PageInfos
Alpine.store('global').allPageNum = parseInt(book.page_count)
// 临时用户标签ID
const tabID = (Date.now() % 10000000).toString(36) + Math.random().toString(36).substring(2, 5)
// 假设token是一个有效的令牌 TODO:使用真正的令牌
const token = 'your_token'

// 滑动相关变量
let touchStartX = 0
let touchEndX = 0
let isSwiping = false
let currentTranslate = 0
let startTime = 0
let animationID = 0
const sliderContainer = document.getElementById('manga_area')
const slider = document.getElementById('slider')
const leftSlide = document.getElementById('left-slide')
//const middleSlide = document.getElementById('middle-slide')
const rightSlide = document.getElementById('right-slide')
const threshold = 100 // 滑动阈值，超过这个值才会触发翻页
const swipeTimeout = 300 // 滑动超时时间（毫秒）

// 设置图片资源，预加载等
// 需要 HTTP 响应头中允许缓存（没有使用 Cache-Control: no-cache），也就是 gin 不能用htmx/router/server.go 里面的noCache 中间件。
// 在预加载用到的图片资源URL
let preloadedImages = new Set()

//首次加载时
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

function GetImageSrc(index) {
    if (index < 0 || index >= images.length) {
        console.log(`Error,cannot use this index: ${index}`);
        return
    }
    const Url = images[index].url
    if (Alpine.store('global').onlineBook){
        const autoCrop = Alpine.store('global').autoCrop ? "&auto_crop=" + Alpine.store('global').autoCropNum : ''
        const autoResize = Alpine.store('global').autoResize ? "&resize_max_width=" + Alpine.store('global').autoResizeWidth : ''
        const noCache = Alpine.store('global').noCache ? "&no-cache=true" : ''
        return `${Url}${autoCrop}${autoResize}${noCache}`
    }
    return `${Url}`
}

// 加载图片资源
function setImageSrc() {
    const nowPageNum = Alpine.store('global').nowPageNum
    if (nowPageNum === 0 && nowPageNum >= Alpine.store('global').allPageNum) {
        console.log('setImageSrc: nowPageNum is 0 or out of range', nowPageNum)
        return
    }
    // console.log("setImageSrc: nowPageNum", nowPageNum);
    // console.log("setImageSrc: NowImage", images[nowPageNum - 1].url);
    // console.log("setImageSrc: NowImage+1=", images[nowPageNum].url);
    // 加载当前图片
    document.getElementById('Single-NowImage').src =
        GetImageSrc(nowPageNum - 1)
    if (!Alpine.store('flip').mangaMode) {
        document.getElementById('Double-NowImage-Left').src = GetImageSrc(nowPageNum - 1);
    } else {
        document.getElementById('Double-NowImage-Right').src = GetImageSrc(nowPageNum - 1);
    }

    preloadedImages.add(GetImageSrc(nowPageNum - 1))
    // 更新滑动容器图片
    updateSliderImages(nowPageNum)

    // 为双页模式，加载下一张图片。
    // 因为用户有可能随时切换到双页模式，所以单页模式也预加载图片（尽管看不到）
    if (nowPageNum < Alpine.store('global').allPageNum) {
        if (Alpine.store('flip').mangaMode) {
            document.getElementById('Double-NowImage-Left').src = GetImageSrc(nowPageNum);
        } else {
            document.getElementById('Double-NowImage-Right').src = GetImageSrc(nowPageNum);
        }
        preloadedImages.add(GetImageSrc(nowPageNum))
    }

    // 预加载前一张和后十张图片
    const preloadRange = 10 // 预加载范围，可以根据需要调整
    for (let i = nowPageNum - 2; i <= nowPageNum + preloadRange; i++) {
        if (i >= 0 && i < Alpine.store('global').allPageNum) {
            const imgUrl = GetImageSrc(i)
            if (!preloadedImages.has(imgUrl)) {
                let img = new Image()
                img.src = imgUrl
                preloadedImages.add(imgUrl)
            }
        }
    }
}


// 更新滑动容器图片
function updateSliderImages(nowPageNum) {
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
            prevImg.src = GetImageSrc(nowPageNum - 2)
            prevImg.className = Alpine.store('global').isPortrait ? 'object-contain w-auto max-w-full h-screen' : 'h-screen w-auto max-w-full object-contain'
            prevImg.draggable = false
            leftSlide.innerHTML = ''
            leftSlide.appendChild(prevImg)
        } else {
            leftSlide.innerHTML = ''
        }
        // // 更新当前图片 (确保当前图片也在这里更新，以防万一)
        const currentImgElement = document.getElementById('Single-NowImage')
        if (currentImgElement && nowPageNum >= 1 && nowPageNum <= Alpine.store('global').allPageNum) {
            currentImgElement.src = GetImageSrc(nowPageNum - 1)
        }
        // 添加后一张图片（如果存在）
        if (nowPageNum < Alpine.store('global').allPageNum) {
            const nextImg = document.createElement('img')
            nextImg.src = GetImageSrc(nowPageNum)
            nextImg.className = Alpine.store('global').isPortrait ? 'object-contain w-auto max-w-full h-screen' : 'h-screen w-auto max-w-full object-contain'
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
        if (nowPageNum === 2) {
            const prevImg = document.createElement('img')
            prevImg.src = GetImageSrc(nowPageNum - 2)
            prevImg.className = 'object-contain h-screen max-w-full max-h-screen m-0'
            prevImg.draggable = false
            leftSlide.innerHTML = ''
            leftSlide.appendChild(prevImg)
        }
        if (nowPageNum >= 3) {
            const prevImg_1 = document.createElement('img')
            prevImg_1.src = GetImageSrc(nowPageNum - 2)
            prevImg_1.className = 'object-contain w-auto max-h-screen m-0 select-none max-w-1/2 grow-0'
            prevImg_1.draggable = false
            const prevImg_2 = document.createElement('img')
            prevImg_2.src = GetImageSrc(nowPageNum - 3)
            prevImg_2.className = 'object-contain w-auto max-h-screen m-0 select-none max-w-1/2 grow-0'
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
        if (nowPageNum === Alpine.store('global').allPageNum - 3) {
            const nextImg = document.createElement('img')
            nextImg.src = GetImageSrc(nowPageNum - 2)
            nextImg.className = 'object-contain h-screen max-w-full max-h-screen m-0'
            nextImg.draggable = false
            rightSlide.innerHTML = ''
            rightSlide.appendChild(nextImg)
        }
        if (nowPageNum < Alpine.store('global').allPageNum - 3) {
            const nextImg_1 = document.createElement('img')
            nextImg_1.src = GetImageSrc(nowPageNum + 1)
            nextImg_1.className = 'object-contain w-auto max-h-screen m-0 select-none max-w-1/2 grow-0'
            nextImg_1.draggable = false
            const nextImg_2 = document.createElement('img')
            nextImg_2.src = GetImageSrc(nowPageNum + 2)
            nextImg_2.className = 'object-contain w-auto max-h-screen m-0 select-none max-w-1/2 grow-0'
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
        if (nowPageNum === Alpine.store('global').allPageNum - 1) {
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
    //console.log('touchStart,swipeTurn:' + Alpine.store('flip').swipeTurn)
    startTime = new Date().getTime()
    isSwiping = true
    touchStartX = e.type === 'touchstart' ? e.touches[0].clientX : e.clientX

    // 停止任何正在进行的动画
    cancelAnimationFrame(animationID)
}

// 如果在第一页或最后一页尝试向前翻或向后翻，阻止默认滚动
function shouldBlockScroll(diffX) {
    const mangaMode = Alpine.store('flip').mangaMode
    const nowPageNum = Alpine.store('global').nowPageNum
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
    if (nowPageNum === Alpine.store('global').allPageNum) {
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
    // 定义动画函数
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

    // 定义动画函数
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
    const nowPageNum = Alpine.store('global').nowPageNum
    updateSliderImages(nowPageNum)
})

//翻页函数，加页或减页
function addPageNum(n = 1) {
    // 防止n为字符串，转换为数字
    let nowPageNum = parseInt(Alpine.store('global').nowPageNum)
    // 无法继续翻
    if (nowPageNum + n > Alpine.store('global').allPageNum) {
        showToast(i18next.t('hint_last_page'), 'warning')
        return
    }
    if (nowPageNum + n < 1) {
        showToast(i18next.t('hint_first_page'), 'warning')
        return
    }
    // 翻页
    Alpine.store('global').nowPageNum = nowPageNum + n
    // 更新书签
    if (!!book && !!book.id && Alpine.store('global').onlineBook) {
        Alpine.store('global').UpdateBookmark({
            type: 'auto',
            bookId: book.id,
            pageIndex: Alpine.store('global').nowPageNum,
        });
    }
    setImageSrc()
    // 设置标签页标题
    setTitle();
    // 通过ws通道发送翻页数据
    if (Alpine.store('global').syncPageByWS === true) {
        sendFlipData() // 发送翻页数据
    }
    // 调用保存页数函数
    Alpine.store('global').savePageNumToLocalStorage();
}

//翻页函数，跳转到指定页
function inputPageNum(event) {
    const i = parseInt(event.target.value)
    let num = Alpine.store('flip').mangaMode ? (Alpine.store('global').allPageNum - i + 1) : i
    //console.log(num)
    jumpPageNum(num)
}


//翻页函数，跳转到指定页
function jumpPageNum(jumpNum, sync = true) {
    let num = parseInt(jumpNum)

    if (num <= 0 || num > Alpine.store('global').allPageNum) {
        alert(i18next.t('hintPageNumOutOfRange'))
        return
    }
    Alpine.store('global').nowPageNum = num
    if (Alpine.store('global').onlineBook) {
        Alpine.store('global').UpdateBookmark({
            type: 'auto',
            bookId: book.id,
            pageIndex: Alpine.store('global').nowPageNum,
        });
        if (sync) {
            // 通过ws通道发送翻页数据
            if (Alpine.store('global').syncPageByWS === true) {
                sendFlipData() // 发送翻页数据
            }
        }
    }
    // 调用保存页数函数
    Alpine.store('global').savePageNumToLocalStorage();
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
        // 将 manga_area 顶部对齐到浏览器可见区域顶部
        const mangaMain = document.getElementById('manga_area')
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
    // 现在设定为固定的80px高度，这样会比较自然
    if (Alpine.store('flip').autoHideToolbar) {
        // 判断鼠标是否在元素 1 范围内(Header)。。
        inInElement1 = (y <= 80)
        // 判断鼠标是否在元素 2 范围内(导航条)。因为header可能隐藏，所以不能直接用left、right、top、bottom判断y是否在header的范围内。
        inInElement2 = (y >= window.innerHeight - 80)
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

    // 鼠标在设置区域
    let inSetArea = getInSetArea(event)
    // 鼠标不在设置区域 + 不在任何一个元素范围内
    if (inSetArea || inInElement1 || inInElement2 || inInElement3 || inInElement4 || inInElement5) {
        showToolbar()
    } else {
        // '鼠标不在设置区域 + 不在任何一个元素范围内'
        //console.log(`inSetArea:${inSetArea}`)
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

//获取ID为 mouseMoveArea 的元素
let mouseMoveArea = document.getElementById('mouseMoveArea')
// 鼠标移动时触发移动事件
mouseMoveArea.addEventListener('mousemove', onMouseMove)
//点击的时候触发点击事件
mouseMoveArea.addEventListener('click', onMouseClick)
// 触摸开始时触发点击事件
mouseMoveArea.addEventListener('touchstart', onMouseClick)
//离开的时候触发离开事件
mouseMoveArea.addEventListener('mouseleave', onMouseLeave)

// Websocket 连接和消息处理
// https://www.ruanyifeng.com/blog/2017/05/websocket.html
// https://developer.mozilla.org/zh-CN/docs/Web/API/WebSocket

// 定义WebSocket变量和重连参数
let socket = null // 初始化为 null
let reconnectAttempts = 0
const maxReconnectAttempts = 200
const reconnectInterval = 3000 // 每次重连间隔3秒


// 翻页数据，假设已在其他地方定义
const flip_data = {
    book_id: book.id,
    now_page_num: Alpine.store('global').nowPageNum,
    need_double_page_mode: false,
}

// 建立WebSocket连接的函数
function connectWebSocket() {
    // 根据当前协议选择ws或wss
    // 检查是否已存在连接或正在连接
    if (socket && (socket.readyState === WebSocket.CONNECTING || socket.readyState === WebSocket.OPEN)) {
        console.log("WebSocket 正在连接或已打开，跳过。");
        return;
    }

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
    // console.log("Local Tab：" + tabID);
    // console.log("message_sender_id：" + message.user_id);// 用message_sender_id或许比较好区分？
    // 根据消息类型进行处理
    if (message.type === 'flip_mode_sync_page' && message.tab_id !== tabID) {
        // 解析翻页数据
        const data = JSON.parse(message.data_string)
        if (Alpine.store('global').syncPageByWS && data.book_id === book.id) {
            //console.log("同步页数：", data);
            // 更新页面(跳转到指定页，但是不发送翻页消息，因为这样会引起是循环)
            jumpPageNum(data.now_page_num, false)
        }
    } else if (message.type === 'heartbeat') {
        console.log('收到心跳消息')
    } else {
        //console.log("不处理此消息"+message);
    }
}

// 发送翻页数据到服务器
function sendFlipData() {
    const flip_data = {
        book_id: book.id,
        now_page_num: Alpine.store('global').nowPageNum,
    }
    const flipMsg = {
        type: 'flip_mode_sync_page', // 或 "heartbeat"
        status_code: 200,
        tab_id: tabID,
        token: token,
        detail: '翻页模式，发送数据',
        data_string: JSON.stringify(flip_data),
    }
    // 确保 socket 已初始化并且处于 OPEN 状态
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(flipMsg))
    } else {
        console.log('WebSocket 未连接或未准备好，无法发送消息。 State:', socket ? socket.readyState : 'null');
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
document.addEventListener('DOMContentLoaded', () => {
    if (Alpine.store('global').onlineBook) {
        connectWebSocket()
    }
})

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
