// 翻页模式分页与边界判断工具
// （立即执行函数）：(function(){...})(...) 
// 创建一个“一次性、私有”的作用域。作用：避免污染全局作用域，即只暴露api，而不暴露实现细节。
(function (global) {
    'use strict'
    // 判断在第一页/最后一页时，当前滑动方向是否应该被“阻断”。
    function shouldBlockScrollBoundary(diffX, mangaMode, nowPageNum, allPageNum) {
        if (nowPageNum === 1) {
            return (diffX < 0 && mangaMode) || (diffX > 0 && !mangaMode)
        }
        if (nowPageNum === allPageNum) {
            return (diffX > 0 && mangaMode) || (diffX < 0 && !mangaMode)
        }
        return false
    }

    // 计算“下一次翻页应该加几页”：1 / 2 / 0。
    // 0 表示不该再翻（越界保护）。
    function getNextPageStep(doublePageMode, nowPageNum, allPageNum) {
        if (!doublePageMode) {
            return nowPageNum <= allPageNum ? 1 : 0
        }
        if (nowPageNum < allPageNum - 1) {
            return 2
        }
        return 1
    }

    // 计算“上一次翻页应该减几页”：-1 / -2 / 0。
    // 0 表示不该再翻（越界保护）。
    function getPreviousPageStep(doublePageMode, nowPageNum) {
        if (nowPageNum <= 1) {
            return 0
        }
        if (!doublePageMode) {
            return -1
        }
        return nowPageNum-2 > 0 ? -2 : -1
    }

    const api = {
        shouldBlockScrollBoundary,
        getNextPageStep,
        getPreviousPageStep,
    }
    global.ComiGoFlip = global.ComiGoFlip || {}
    global.ComiGoFlip.pagination = api

    if (typeof module !== 'undefined' && module.exports) {
        module.exports = api
    }
})(typeof window !== 'undefined' ? window : globalThis)
