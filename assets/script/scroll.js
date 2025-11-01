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
// 最大页码
const MaxPageNum = Math.floor(parseInt(book.page_count) / 32) + 1

//滚动到顶部
function scrollToTop(scrollDuration) {
    let scrollStep = -window.scrollY / (scrollDuration / 15),
        scrollInterval = setInterval(function () {
            if (window.scrollY !== 0) {
                window.scrollBy(0, scrollStep)
            } else clearInterval(scrollInterval)
        }, 15)
}
// Button ID为BackTopButton的元素，点击后滚动到顶部
document.getElementById('BackTopButton').addEventListener('click', function () {
    scrollToTop(500)
})

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
    // console.log("scrollDownFlag:",scrollDownFlag,"scrollTop:",scrollTop,"step:", step);
    scrollTopSave = scrollTop
    if (step < -10 || step > 10) {
        showBackTopFlag = scrollTop > 400 && !scrollDownFlag
        if (showBackTopFlag) {
            document.getElementById('BackTopButton').style.display = 'block'
        } else {
            document.getElementById('BackTopButton').style.display = 'none'
        }
    }
}
window.addEventListener('scroll', onScroll)

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
//获取ID为 mouseMoveArea 的元素
let mouseMoveArea = document.getElementById('mouseMoveArea')
// 鼠标移动的时候触发移动事件
mouseMoveArea.addEventListener('mousemove', onMouseMove)
// 点击的时候触发点击事件
mouseMoveArea.addEventListener('click', onMouseClick)
// 触摸的时候也触发点击事件
mouseMoveArea.addEventListener('touchstart', onMouseClick)

// 键盘快捷键
/* 记录方向键当前的按压状态 */
// 1) 方向/动作当前状态
const state = { up: false, down: false, left: false, right: false, fire: false };

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
};

// 3) 通用按键处理器：down=true 表示按下，false 表示松开
function handle(e, down) {
    const k = e.key.toLowerCase();      // 统一小写
    const act = keyMap[k];              // 查映射表
    if (!act) return;                   // 映射表里没有，忽略
    state[act] = down;                  // 更新状态
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
        toPreviousPage()
    }
    // 触按下相当于右方向键的按键的时候
    if (act === "right" && down) {
        toNextPage()
    }
    // 直接跳转到第一页,同时长按的时候不执行多次
    if (act === "first_page" && down && !e.repeat) {
        jumpPageNum(1)
    }
    // 直接跳转到最后一页,同时长按的时候不执行多次
    if (act === "last_page" && down && !e.repeat) {
        jumpPageNum(MaxPageNum)
    }
}

// 4) 事件监听
// 监听键盘事件
// keydown 第一次按下和按住时会触发 //统一禁止长按： addEventListener("keydown", e => !e.repeat && handle(e, true));
addEventListener("keydown", e => handle(e, true));
// keyup 松开时触发
addEventListener("keyup", e => handle(e, false));

// 根据url获取当前页码 当前 url 类似 http://localhost:1234/scroll/somebookid?page=1
function getNowPageNum() {
    const urlParams = new URLSearchParams(window.location.search);
    const page = parseInt(urlParams.get('page'));
    return isNaN(page) ? 1 : page;
}

// 根据当前页码设置url并刷新，如果小于最小页码（1），打印错误并返回
function toPreviousPage() {
    // 如果是无限下拉模式，不响应翻页，直接返回
    if (!Alpine.store('scroll').fixedPagination){
        return;
    }
    const currentPage = getNowPageNum();
    if (currentPage <= 1) {
        console.warn(`已经是第一页了。MaxPageNum：${MaxPageNum}`);
        showToast(i18next.t('hint_first_page'), 'warning');
        return;
    }
    const newPage = currentPage - 1;
    const url = new URL(window.location.href);
    url.searchParams.set('page', newPage);
    window.location.href = url.toString();
}

// 根据当前页码设置url并刷新，如果大于最大页码（MaxPageNum），打印错误并返回
function toNextPage() {
    // 如果是无限下拉模式，不响应翻页，直接返回
    if (!Alpine.store('scroll').fixedPagination){
        return;
    }
    const currentPage = getNowPageNum();
    if (currentPage >= MaxPageNum) {
        console.warn(`已经是最后一页了。MaxPageNum：${MaxPageNum}`);
        showToast(i18next.t('hint_last_page'), 'warning')
        return;
    }
    const newPage = currentPage + 1;
    const url = new URL(window.location.href);
    url.searchParams.set('page', newPage);
    window.location.href = url.toString();
}

// 根据当前页码设置url并刷新，如果小于最小页码（1）或大于最大页码（MaxPageNum），打印错误并返回
function jumpPageNum(pageNum) {
    // 如果是无限下拉模式，不响应翻页，直接返回
    if (!Alpine.store('scroll').fixedPagination){
        return;
    }
    if (pageNum < 1 || pageNum > MaxPageNum) {
        console.warn(`页码超出范围，有效范围为1-${MaxPageNum}`);
        return;
    }
    const url = new URL(window.location.href);
    url.searchParams.set('page', pageNum);
    window.location.href = url.toString();
} 