// 使用标准 <script> 标记插入的 JavaScript 代码
//https://templ.guide/syntax-and-usage/script-templates/

//https://www.runoob.com/js/js-strict.html
"use strict";

const book = JSON.parse(document.getElementById('NowBook').textContent);
//console.log(book);
// const globalState = JSON.parse(document.getElementById('GlobalState').textContent);

// 设置图片资源，预加载等
// 需要 HTTP 响应头中允许缓存（没有使用 Cache-Control: no-cache），也就是 gin 不能用htmx/router/server.go 里面的noCache 中间件。
// 在预加载用到的图片资源URL
let preloadedImages = new Set();

function setImageSrc() {
    const nowPageNum = Alpine.store('flip').nowPageNum;
    const images = book.pages.images;
    const totalImages = images.length;

    if (nowPageNum !== 0 && nowPageNum <= totalImages) {
        // 加载当前图片
        document.getElementById('NowImage').src = images[nowPageNum - 1].url;
        preloadedImages.add(images[nowPageNum - 1].url);
        // 双页模式，加载下一张图片
        if (Alpine.store('flip').doublePageMode && nowPageNum < totalImages) {
            document.getElementById('NextImage').src = images[nowPageNum].url;
            preloadedImages.add(images[nowPageNum].url);
        }

        // 预加载前一张和后十张图片
        const preloadRange = 10; // 预加载范围，可以根据需要调整
        for (let i = nowPageNum - 2; i <= nowPageNum + preloadRange; i++) {
            if (i >= 0 && i < totalImages) {
                const imgUrl = images[i].url;
                if (!preloadedImages.has(imgUrl)) {
                    let img = new Image();
                    img.src = imgUrl;
                    preloadedImages.add(imgUrl);
                }
            }
        }
    } else {
        console.log("setImageSrc: nowPageNum is 0 or out of range", nowPageNum);
    }
}


//首次加载时
setImageSrc();
//设置初始值
Alpine.store('flip').allPageNum = book.page_count;

//翻页函数，加页或减页
function addPageNum(n = 1) {
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
    // 通过ws通道发送翻页数据
    if (Alpine.store('global').syncPageByWS === true) {
        sendFlipData(); // 发送翻页数据
    }
}

//翻页函数，跳转到指定页
function jumpPageNum(jumpNum) {
    let num = parseInt(jumpNum);
    let allPageNum = parseInt(Alpine.store('flip').allPageNum);
    if (num <= 0 || num > allPageNum) {
        alert(i18next.t("hintPageNumOutOfRange"));
        return;
    }
    Alpine.store('flip').nowPageNum = num;
    setImageSrc();
}


// 翻页函数，下一页
function toNextPage() {
    let doublePageMode = Alpine.store('flip').doublePageMode === true;
    let nowPageNum = parseInt(Alpine.store('flip').nowPageNum);
    let allPageNum = parseInt(Alpine.store('flip').allPageNum);
    // 单页模式
    if (!doublePageMode) {
        if (nowPageNum < allPageNum) {
            addPageNum(1);
        }
    }
    //双页模式
    if (doublePageMode) {
        if (nowPageNum < allPageNum - 1) {
            addPageNum(2);
        } else {
            addPageNum(1);
        }
    }
}

// 翻页函数，前一页
function toPreviousPage() {
    let doublePageMode = Alpine.store('flip').doublePageMode === true;
    let nowPageNum = parseInt(Alpine.store('flip').nowPageNum);
    //错误值,第0或第1页。
    if (nowPageNum <= 1) {
        alert(i18next.t("hintFirstPage"));
        return;
    }
    //简单合并模式
    if (doublePageMode) {
        if (nowPageNum - 2 > 0) {
            addPageNum(-2);
        } else {
            addPageNum(-1);
        }
    } else {
        addPageNum(-1);
    }
}

//隐藏工具栏的工具函数
// https://www.runoob.com/js/js-htmldom-events.html
let hideTimeout;
let header = document.getElementById("header");
let range = document.getElementById("steps-range_area");

// 显示工具栏
function showToolbar() {
    if (Alpine.store('flip').autoHideToolbar) {
        header.style.opacity = '0.9';
        range.style.opacity = '0.9';
    } else {
        header.style.opacity = '1';
        range.style.opacity = '1';
    }
}

// 隐藏工具栏
function hideToolbar() {
    if (Alpine.store('flip').autoHideToolbar) {
        header.style.opacity = '0';
        range.style.opacity = '0';
    }
}

function startHideTimer() {
    hideTimeout = setTimeout(hideToolbar, 1000); // 3 seconds
}

if (Alpine.store('flip').autoHideToolbar) {
    startHideTimer();
}

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
    // // https://developer.mozilla.org/zh-CN/docs/Web/API/Event/stopPropagation
    // // 相当于 Alpinejs的  @click.stop="onMouseMove()" ? https://alpinejs.dev/directives/on#prevent
    // e.stopPropagation();
    let clickX = e.x; //获取鼠标的X坐标（鼠标与屏幕左侧的距离,单位为px）
    //浏览器的视口宽,不包括工具栏和滚动条:
    let innerWidth = window.innerWidth
    let inSetArea = getInSetArea(e);
    if (inSetArea) {
        //获取ID为 OpenSettingButton的元素，然后模拟点击
        document.getElementById("OpenSettingButton").click();
    }
    if (!inSetArea) {
        //决定如何翻页
        if (clickX < innerWidth * 0.5) {
            //左边的翻页
            if (Alpine.store('flip').rightToLeft) {
                toPreviousPage();
            } else {
                toNextPage();
            }
        } else {
            //右边的翻页
            if (Alpine.store('flip').rightToLeft) {
                toNextPage();
            } else {
                toPreviousPage();
            }
        }
    }
    if (inSetArea) {
        showToolbar();
    }
}


//获取鼠标位置,决定如何显示鼠标指针
function onMouseMove(e) {
    // https://developer.mozilla.org/zh-CN/docs/Web/API/Event/stopPropagation
    // 相当于 Alpinejs的  @click.stop="onMouseMove()" ? https://alpinejs.dev/directives/on#prevent
    let clickX = e.x; //获取鼠标的X坐标（鼠标与屏幕左侧的距离,单位为px）
    // 浏览器的视口宽,不包括工具栏和滚动条:
    let innerWidth = window.innerWidth;

    //在设置区域
    let inSetArea = getInSetArea(e);
    if (inSetArea) {
        e.currentTarget.style.cursor =
            "url(/static/images/SettingsOutline.png), pointer";
    }
    if (!inSetArea) {
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
}

// 获取两个元素的边界信息
function getElementsRect() {
    return {
        rect1: header.getBoundingClientRect(),
        rect2: range.getBoundingClientRect(),
        rect3: document.getElementById("ReSortDropdown").getBoundingClientRect(),
        rect4: document.getElementById("QuickJumpDropdown").getBoundingClientRect(),
    };
}

document.addEventListener('mousemove', function (event) {
    const {rect1, rect2, rect3,rect4} = getElementsRect();
    const x = event.clientX;
    const y = event.clientY;
    // 判断鼠标是否在元素 1 范围内(Header)
    const isInElement1 = (
        x >= rect1.left &&
        x <= rect1.right &&
        y >= rect1.top &&
        y <= rect1.bottom
    );
    // 判断鼠标是否在元素 2 范围内(导航条)
    const isInElement2 = (
        x >= rect2.left &&
        x <= rect2.right &&
        y >= rect2.top &&
        y <= rect2.bottom
    );
    // 判断鼠标是否在元素 3 范围内(页面重新排序的下拉菜单。在菜单上面的时候，导航条需要保持显示状态。)
    const isInElement3 = (
        x >= rect3.left &&
        x <= rect3.right &&
        y >= rect3.top &&
        y <= rect3.bottom
    );
    // 判断鼠标是否在元素 4 范围内(快速跳转的下拉菜单。在菜单上面的时候，导航条需要保持显示状态。)
    const isInElement4 = (
        x >= rect4.left &&
        x <= rect4.right &&
        y >= rect4.top &&
        y <= rect4.bottom
    );

    //鼠标在设置区域
    let inSetArea = getInSetArea(event);
    if (inSetArea) {
        showToolbar();
    }
    // '鼠标不在任何一个元素范围内'
    if (!isInElement1 && !isInElement2 && !isInElement3) {
        if (inSetArea === false) {
            hideToolbar();
        }
    }
    // '鼠标在元素范围内'
    if (isInElement1 || isInElement2 || isInElement3 || isInElement4 || inSetArea) {
        //console.log('鼠标在元素范围内');
        showToolbar();
    }
});

//可见区域变化时，改变页面状态
function onResize() {
    Alpine.store("flip").imageMaxWidth = window.innerWidth
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
//初始化时,执行一次onResize()
onResize();
//文档视图调整大小时触发 resize 事件。 https://developer.mozilla.org/zh-CN/docs/Web/API/Window/resize_event
window.addEventListener("resize", onResize);

//离开区域的时候,清空鼠标样式
function onMouseLeave(e) {
    e.currentTarget.style.cursor = "";
}


//获取ID为 mouseMoveArea 的元素
let mouseMoveArea = document.getElementById("mouseMoveArea")
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
let socket;
let reconnectAttempts = 0;
const maxReconnectAttempts = 200;
const reconnectInterval = 3000; // 每次重连间隔3秒

// 用户ID和令牌，假设已在其他地方定义
const userID = Alpine.store('global').userID;
const token = "your_token";

// 翻页数据，假设已在其他地方定义
const flip_data = {
    book_id: book.id,
    now_page_num: Alpine.store('flip').nowPageNum,
    need_double_page_mode: false,
};

// 建立WebSocket连接的函数
function connectWebSocket() {
    // 根据当前协议选择ws或wss
    const wsProtocol = window.location.protocol === "https:" ? "wss://" : "ws://";
    const wsUrl = wsProtocol + window.location.host + "/api/ws";
    socket = new WebSocket(wsUrl);

    // 连接打开时的回调
    socket.onopen = function () {
        console.log("WebSocket连接已建立");
        reconnectAttempts = 0; // 重置重连次数
    };

    // 收到消息时的回调
    socket.onmessage = function (event) {
        const message = JSON.parse(event.data);
        handleMessage(message); // 调用处理函数
    };

    // 连接关闭时的回调
    socket.onclose = function () {
        console.log("WebSocket连接已关闭");
        attemptReconnect(); // 尝试重连
    };

    // 发生错误时的回调
    socket.onerror = function (error) {
        console.log("WebSocket发生错误：", error);
        socket.close(); // 关闭连接以触发重连
    };
}

// 处理收到的翻页消息
function handleMessage(message) {
    // console.log("收到消息：", message);
    // console.log("My userID：" + userID);
    // console.log("Remote userID：" + message.user_id);
    // 根据消息类型进行处理
    if (message.type === "flip_mode_sync_page" && message.user_id !== userID) {
        // 解析翻页数据
        const data = JSON.parse(message.data_string);
        if (Alpine.store('global').syncPageByWS && data.book_id === book.id) {
            //console.log("同步页数：", data);
            // 更新翻页数据
            flip_data.now_page_num = data.now_page_num;
            // 更新页面
            jumpPageNum(data.now_page_num);
        }
    } else if (message.type === "heartbeat") {
        console.log("收到心跳消息");
    } else {
        //console.log("不处理此消息"+message);
    }
}

// 发送翻页数据到服务器
function sendFlipData() {
    flip_data.now_page_num = Alpine.store('flip').nowPageNum;
    const flipMsg = {
        type: "flip_mode_sync_page", // 或 "heartbeat"
        status_code: 200,
        user_id: userID,
        token: token,
        detail: "翻页模式，发送数据",
        data_string: JSON.stringify(flip_data),
    };
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(flipMsg));
        //console.log("已发送翻页数据"+JSON.stringify(flipMsg));
        //console.log("已发送翻页数据");
    } else {
        console.log("WebSocket未连接，无法发送消息");
    }
}

// 尝试重连函数
function attemptReconnect() {
    if (reconnectAttempts < maxReconnectAttempts) {
        reconnectAttempts++;
        console.log(`第 ${reconnectAttempts} 次重连...`);
        setTimeout(() => {
            connectWebSocket();
        }, reconnectInterval);
    } else {
        console.log("已达到最大重连次数，停止重连");
    }
}

// 页面加载完成后建立WebSocket连接
window.onload = function () {
    connectWebSocket();
};
