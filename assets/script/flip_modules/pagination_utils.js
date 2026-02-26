// 翻页模式分页与边界判断工具（可复用、可测试）
(function (global) {
    'use strict'

    function shouldBlockScrollBoundary(diffX, mangaMode, nowPageNum, allPageNum) {
        if (nowPageNum === 1) {
            return (diffX < 0 && mangaMode) || (diffX > 0 && !mangaMode)
        }
        if (nowPageNum === allPageNum) {
            return (diffX > 0 && mangaMode) || (diffX < 0 && !mangaMode)
        }
        return false
    }

    function getNextPageStep(doublePageMode, nowPageNum, allPageNum) {
        if (!doublePageMode) {
            return nowPageNum <= allPageNum ? 1 : 0
        }
        if (nowPageNum < allPageNum - 1) {
            return 2
        }
        return 1
    }

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
