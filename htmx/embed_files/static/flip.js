// 使用标准 <script> 标记插入的 JavaScript 代码
//https://templ.guide/syntax-and-usage/script-templates/

//https://www.runoob.com/js/js-strict.html
"use strict";
const book = JSON.parse(document.getElementById('NowBook').textContent);
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

//翻页函数，任意页
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
// 翻页函数，下一页
function toNextPage() {
  let doublePageMode = Alpine.store('flip').doublePageMode === true;
  let nowPageNum = parseInt(Alpine.store('flip').nowPageNum);
  let allPageNum = parseInt(Alpine.store('flip').allPageNum);
  // 单页模式
  if (!doublePageMode) {
    if (nowPageNum < allPageNum) {
      flipPage(1);
    }
  }
  //双页模式
  if (doublePageMode) {
    if (nowPageNum < allPageNum - 1) {
      flipPage(2);
    } else {
      flipPage(1);
    }
  }
}
// 翻页函数，前一页
function toPerviousPage() {
  let doublePageMode = Alpine.store('flip').doublePageMode === true;
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

//获取鼠标位置,决定是否打开设置面板
function onMouseClick(e) {
  // https://developer.mozilla.org/zh-CN/docs/Web/API/Event/stopPropagation
  // 相当于 Alpinejs的  @click.stop="onMouseMove()" ? https://alpinejs.dev/directives/on#prevent
  e.stopPropagation();
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
    inSetArea = true;
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
  if (inSetArea) {
    showToolbar();
  }
}


//获取鼠标位置,决定如何显示鼠标指针
function onMouseMove(e) {
  // https://developer.mozilla.org/zh-CN/docs/Web/API/Event/stopPropagation
  // 相当于 Alpinejs的  @click.stop="onMouseMove()" ? https://alpinejs.dev/directives/on#prevent
  let clickX = e.x; //获取鼠标的X坐标（鼠标与屏幕左侧的距离,单位为px）
  let clickY = e.y; //获取鼠标的Y坐标（鼠标与屏幕顶部的距离,单位为px）
  // 浏览器的视口,不包括工具栏和滚动条:
  let innerWidth = window.innerWidth;
  let innerHeight = window.innerHeight;

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
    rect3: document.getElementById("dropdownHover").getBoundingClientRect()
  };
}

document.addEventListener('mousemove', function (event) {
  const { rect1, rect2,rect3 } = getElementsRect();
  const x = event.clientX;
  const y = event.clientY;
  // 判断鼠标是否在元素 1 范围内
  const isInElement1 = (
    x >= rect1.left &&
    x <= rect1.right &&
    y >= rect1.top &&
    y <= rect1.bottom
  );
  // 判断鼠标是否在元素 2 范围内
  const isInElement2 = (
    x >= rect2.left &&
    x <= rect2.right &&
    y >= rect2.top &&
    y <= rect2.bottom
  );
  // 判断鼠标是否在元素 3 范围内
  const isInElement3 = (
    x >= rect3.left &&
    x <= rect3.right &&
    y >= rect3.top &&
    y <= rect3.bottom
    );

  let clickX = event.x; //获取鼠标的X坐标（鼠标与屏幕左侧的距离,单位为px）
  let clickY = event.y; //获取鼠标的Y坐标（鼠标与屏幕顶部的距离,单位为px）
  // 浏览器的视口,不包括工具栏和滚动条:
  let innerWidth = window.innerWidth;
  let innerHeight = window.innerHeight;

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
    showToolbar();
  }
  // 判断鼠标是否不在任一元素范围内
  if (!isInElement1 && !isInElement2 && !isInElement3) {
    //console.log('鼠标不在任何一个元素范围内');
    if (inSetArea === false) {
      hideToolbar();
    }
  }
  if (isInElement1 || isInElement2 || isInElement3 ||inSetArea) {
    //console.log('鼠标在元素范围内');
    showToolbar();
  }
});

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
// 鼠标移动时触发移动事件
mouseMoveArea.addEventListener('mousemove', onMouseMove)
//点击的时候触发点击事件
mouseMoveArea.addEventListener('click', onMouseClick)
