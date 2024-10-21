// 使用标准 <script> 标记插入的 JavaScript 代码
//https://templ.guide/syntax-and-usage/script-templates/
// const globalState = JSON.parse(document.getElementById('GlobalState').textContent);
// console.log(globalState);
const book = JSON.parse(document.getElementById('NowBook').textContent);

function setImageSrc() {
    if(Alpine.store('flip').nowPageNum !== 0 && Alpine.store('flip').nowPageNum !== book.pages.images.length+1){
        document.getElementById('NowImage').src = book.pages.images[Alpine.store('flip').nowPageNum - 1].url;
    }else{
        console.log("setImageSrc: nowPageNum is 0 or out of range", Alpine.store('flip').nowPageNum);
    }
}

setImageSrc();
Alpine.store('flip').allPageNum = book.page_count;

//翻页,其实不限页数
function flipPage(n = 1) {
    //防止数字转换为字符串
    let num = parseInt(n);
    let nowPageNum = parseInt(Alpine.store('flip').nowPageNum);
    let allPageNum = parseInt(Alpine.store('flip').allPageNum);
    //  无法继续翻
    if (nowPageNum + num > allPageNum) {
        alert(i18next.t("hintLastPage"));
        return;
    }
    if (nowPageNum + num < 1) {
        alert(i18next.t("hintFirstPage"));
        return;
    }
    // 翻页
    Alpine.store('flip').nowPageNum = nowPageNum + num;
    setImageSrc();
}

function toNextPage() {
    let doublePageMode = Alpine.store('flip').doublePageMode===true;
    let nowPageNum = parseInt(Alpine.store('flip').nowPageNum);
    let allPageNum = parseInt(Alpine.store('flip').allPageNum);
    if (!doublePageMode) {
        if (nowPageNum < allPageNum) {
            flipPage(1);
        }
    }
    //简单合并模式
    if (doublePageMode) {
        if (nowPageNum < allPageNum - 1) {
            flipPage(2);
        } else {
            flipPage(1);
        }
    }
}

function toPerviousPage() {
    let doublePageMode = Alpine.store('flip').doublePageMode===true;
    let nowPageNum = parseInt(Alpine.store('flip').nowPageNum);
    let allPageNum = parseInt(Alpine.store('flip').allPageNum);
    //错误值,第0或第1页。
    if (nowPageNum <= 1) {
        alert(i18next.t("hintFirstPage"));
        return;
    }
    //简单合并模式
    if (doublePageMode) {
        if (nowPageNum - 2 > 0) {
            flipPage(-2);
        } else {
            flipPage(-1);
        }
    } else {
        flipPage(-1);
    }
}

//获取鼠标位置,决定是否打开设置面板
function onMouseClick(e) {
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
        console.log("点中了设置区域！");
        inSetArea = true
    }
    if (inSetArea) {
        //获取ID为 OpenSettingButton的元素，然后模拟点击
        document.getElementById("OpenSettingButton").click();
    }

    if (!inSetArea) {
        //决定如何翻页
        if (clickX < innerWidth * 0.5) {
            //左边的翻页
            if (Alpine.store('flip').rightToLeft) {
                toPerviousPage();
            } else {
                toNextPage();
            }
        } else {
            //右边的翻页
            if (Alpine.store('flip').rightToLeft) {
                toNextPage();
            } else {
                toPerviousPage();
            }
        }
    }
}

//获取鼠标位置,决定如何显示鼠标指针
function onMouseMove(e) {
    let clickX = e.x; //获取鼠标的X坐标（鼠标与屏幕左侧的距离,单位为px）
    let clickY = e.y; //获取鼠标的Y坐标（鼠标与屏幕顶部的距离,单位为px）
    // console.log("clickX: " + clickX);
    // console.log("clickY: " + clickY);
    // 浏览器的视口,不包括工具栏和滚动条:
    let innerWidth = window.innerWidth;
    let innerHeight = window.innerHeight;

    //是否进入上下工具条附近区域时，触发工具条显隐
    const toolBarArea = 0.10;
    let inToolBarArea = false
    if (clickY <= (innerHeight * toolBarArea) || clickY >= (innerHeight * (1.0 - toolBarArea))) {
        inToolBarArea = true
    }
    //进入设置区域的时候，设置鼠标的形状
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
    if (clickX > MinX && clickX < MaxX && clickY > MinY && clickY < MaxY) {
        inSetArea = true
    }
    if (inSetArea) {
        e.currentTarget.style.cursor =
            "url(/static/images/SettingsOutline.png), pointer";
    } else {
        if (clickX < innerWidth * 0.5) {
            //设置左边的鼠标指针
            if (Alpine.store('flip').rightToLeft && Alpine.store('flip').nowPageNum === 1) {
                //右边翻下一页,且目前是第一页的时候,左边的鼠标指针,设置为禁止翻页
                e.currentTarget.style.cursor =
                    "url(/static/images/Prohibited28Filled.png), pointer";
            } else if (
                !Alpine.store('flip').rightToLeft &&
                Alpine.store('flip').nowPageNum === Alpine.store('flip').allPageNum
            ) {
                //左边翻下一页,且目前是最后一页的时候,左边的鼠标指针,设置为禁止翻页
                e.currentTarget.style.cursor =
                    "url(/static/images/Prohibited28Filled.png), pointer";
            } else {
                //正常情况下,左边是向左的箭头
                e.currentTarget.style.cursor =
                    "url(/static/images/ArrowLeft.png), pointer";
            }
        } else {
            //设置右边的鼠标指针
            if (
                Alpine.store('flip').rightToLeft &&
                Alpine.store('flip').nowPageNum === Alpine.store('flip').allPageNum
            ) {
                //右边翻下一页,且目前是最后页的时候,右边的鼠标指针,设置为禁止翻页
                e.currentTarget.style.cursor =
                    "url(/static/images/Prohibited28Filled.png), pointer";
            } else if (!Alpine.store('flip').rightToLeft && Alpine.store('flip').nowPageNum === 1) {
                //左边翻下一页,且目前是第一页的时候,右边的鼠标指针,设置为禁止翻页
                e.currentTarget.style.cursor =
                    "url(/static/images/Prohibited28Filled.png), pointer";
            } else {
                //正常情况下,右边是向右的箭头
                e.currentTarget.style.cursor =
                    "url(/static/images/ArrowRight.png), pointer";
            }
        }
    }
    //进入上下工具条附近区域时，触发工具条显隐
    if ((inToolBarArea || inSetArea) && Alpine.store('flip').autoHideToolbar) {
        Alpine.store('flip').showHeader = true
        Alpine.store('flip').showFooter = true
    } else if (Alpine.store('flip').autoHideToolbar) {
        Alpine.store('flip').showHeader = false
        Alpine.store('flip').showFooter = false
    }
}

//可见区域变化时，改变页面状态
function onResize() {
    let imageMaxWidth = window.innerWidth
    let clientWidth = document.documentElement.clientWidth
    let clientHeight = document.documentElement.clientHeight
    // var aspectRatio = window.innerWidth / window.innerHeight
    let aspectRatio = clientWidth / clientHeight
    // 为了调试的时候方便,阈值是正方形
    if (aspectRatio > (19 / 19)) {
        Alpine.store('flip').isLandscapeMode = true
        Alpine.store('flip').isPortraitMode = false
    } else {
        Alpine.store('flip').isLandscapeMode = false
        Alpine.store('flip').isPortraitMode = true
    }
}

onResize();
//文档视图调整大小时触发 resize 事件。 https://developer.mozilla.org/zh-CN/docs/Web/API/Window/resize_event
window.addEventListener("resize", onResize);

//离开区域的时候,清空鼠标样式
function onMouseLeave(e) {
    e.currentTarget.style.cursor = "";
}

//获取ID为 mouseMoveArea 的元素
let mouseMoveArea = document.getElementById("mouseMoveArea")
// 鼠标移动的时候触发移动事件
mouseMoveArea.addEventListener('mousemove', onMouseMove)
// 点击的时候触发点击事件
mouseMoveArea.addEventListener('click', onMouseClick)
// 触摸的时候也触发点击事件
mouseMoveArea.addEventListener('touchstart', onMouseClick)

function hideComponent() {
    return {
        showDiv: true,
        hideTimeout: null,
        startHideTimer() {
            // Cancel any existing timer
            this.cancelHideTimer();
            // Start a new timer to hide the div after 3 seconds
            this.hideTimeout = setTimeout(() => {
                this.showDiv = false;
            }, 3000);
        },
        cancelHideTimer() {
            // Clear the hide timer if it exists
            if (this.hideTimeout) {
                clearTimeout(this.hideTimeout);
                this.hideTimeout = null;
            }
            // Ensure the div is shown
            this.showDiv = true;
        }
    }
}
