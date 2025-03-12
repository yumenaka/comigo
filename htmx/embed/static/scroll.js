// 使用标准 <script> 标记插入的 JavaScript 代码
"use strict";

//https://templ.guide/syntax-and-usage/script-templates/
const book = JSON.parse(document.getElementById('NowBook').textContent);
//const globalState = JSON.parse(document.getElementById('GlobalState').textContent);
//console.log(book);
//console.log(globalState);

//滚动到顶部
function scrollToTop(scrollDuration) {
    let scrollStep = -window.scrollY / (scrollDuration / 15),
        scrollInterval = setInterval(function () {
            if (window.scrollY !== 0) {
                window.scrollBy(0, scrollStep);
            }
            else clearInterval(scrollInterval);
        }, 15);
}
// Button ID为BackTopButton的元素，点击后滚动到顶部
document.getElementById("BackTopButton").addEventListener("click", function () {
    scrollToTop(500);
});

//滚动到一定位置显示返回顶部按钮
let scrollTopSave = 0
let scrollDownFlag = false
let showBackTopFlag = false
let step = 0
function onScroll() {
    let scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
    scrollDownFlag = scrollTop > scrollTopSave;
    //防手抖,小于一定数值状态就不变 Math.abs()会导致报错
    step = scrollTopSave - scrollTop;
    // console.log("scrollDownFlag:",scrollDownFlag,"scrollTop:",scrollTop,"step:", step);
    scrollTopSave = scrollTop
    if (step < -10 || step > 10) {
        showBackTopFlag = ((scrollTop > 400) && !scrollDownFlag);
        if (showBackTopFlag) {
            document.getElementById("BackTopButton").style.display = "block";
        } else {
            document.getElementById("BackTopButton").style.display = "none";
        }
    }
}
window.addEventListener("scroll", onScroll);

let isLandscapeMode = true
let isPortraitMode = false
//可见区域变化的时候改变页面状态
function onResize() {
    Alpine.store("scroll").imageMaxWidth = window.innerWidth
    let clientWidth = document.documentElement.clientWidth
    let clientHeight = document.documentElement.clientHeight
    let aspectRatio =  clientWidth / clientHeight
    // 为了调试的时候方便,阈值是正方形
    if (aspectRatio > (19 / 19)) {
        isLandscapeMode = true
        isPortraitMode = false
    } else {
        isLandscapeMode = false
        isPortraitMode = true
    }
}
//初始化时,执行一次onResize()
onResize();
//文档视图调整大小时触发 resize 事件。 https://developer.mozilla.org/zh-CN/docs/Web/API/Window/resize_event
window.addEventListener("resize", onResize);

//鼠标是否在设置区域
function getInSetArea(e) {
    let clickX = e.x //获取鼠标的X坐标（鼠标与屏幕左侧的距离,单位为px）
    let clickY = e.y //获取鼠标的Y坐标（鼠标与屏幕顶部的距离,单位为px）
    //浏览器的视口,不包括工具栏和滚动条:
    let innerWidth = window.innerWidth
    let innerHeight = window.innerHeight
    //设置区域为正方形，边长按照宽或高里面，比较小的值决定
    const setArea = 0.15;
    // innerWidth >= innerHeight 的情况下
    let MinY = innerHeight * (0.5 - setArea);
    let MaxY = innerHeight * (0.5 + setArea);
    let MinX = innerWidth * 0.5 - (MaxY - MinY) * 0.5;
    let MaxX = innerWidth * 0.5 + (MaxY - MinY) * 0.5;
    if (innerWidth < innerHeight) {
        MinX = innerWidth * (0.5 - setArea);
        MaxX = innerWidth * (0.5 + setArea);
        MinY = innerHeight * 0.5 - (MaxX - MinX) * 0.5;
        MaxY = innerHeight * 0.5 + (MaxX - MinX) * 0.5;
    }
    //在设置区域
    let inSetArea = false
    if ((clickX > MinX && clickX < MaxX) && (clickY > MinY && clickY < MaxY)) {
        inSetArea = true
    }
    return inSetArea;
}

//获取鼠标位置,决定是否打开设置面板
function onMouseClick(e) {
    if (getInSetArea(e)) {
        //获取ID为 OpenSettingButton的元素，然后模拟点击
        document.getElementById("OpenSettingButton").click();
    }
}
//获取鼠标位置,决定是否打开设置面板
function onMouseMove(e) {
    if (getInSetArea(e)) {
        e.currentTarget.style.cursor = 'url(/images/SettingsOutline.png), pointer';
    } else {
        e.currentTarget.style.cursor = '';
    }
}
//获取ID为 mouseMoveArea 的元素
let mouseMoveArea = document.getElementById("mouseMoveArea")
// 鼠标移动的时候触发移动事件
mouseMoveArea.addEventListener('mousemove', onMouseMove)
// 点击的时候触发点击事件
mouseMoveArea.addEventListener('click', onMouseClick)
// 触摸的时候也触发点击事件
mouseMoveArea.addEventListener('touchstart', onMouseClick)