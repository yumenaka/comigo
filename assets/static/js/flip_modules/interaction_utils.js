// 阅读交互几何判断工具 IIFE。
// 只集中“区域命中”这类纯判断；点击响应和鼠标图标仍由各阅读模式自己决定。
(function (global) {
    'use strict'

    function getSetAreaBounds(innerWidth, innerHeight, setAreaRatio = 0.15) {
        let minY = innerHeight * (0.5 - setAreaRatio)
        let maxY = innerHeight * (0.5 + setAreaRatio)
        let minX = innerWidth * 0.5 - (maxY - minY) * 0.5
        let maxX = innerWidth * 0.5 + (maxY - minY) * 0.5

        if (innerWidth < innerHeight) {
            minX = innerWidth * (0.5 - setAreaRatio)
            maxX = innerWidth * (0.5 + setAreaRatio)
            minY = innerHeight * 0.5 - (maxX - minX) * 0.5
            maxY = innerHeight * 0.5 + (maxX - minX) * 0.5
        }
        return { minX, maxX, minY, maxY }
    }

    function isInSetArea(clickX, clickY, innerWidth, innerHeight, setAreaRatio = 0.15) {
        const bounds = getSetAreaBounds(innerWidth, innerHeight, setAreaRatio)
        return clickX > bounds.minX && clickX < bounds.maxX && clickY > bounds.minY && clickY < bounds.maxY
    }

    function isPointInRect(rect, x, y) {
        if (!rect) return false
        return x >= rect.left && x <= rect.right && y >= rect.top && y <= rect.bottom
    }

    const api = {
        getSetAreaBounds,
        isInSetArea,
        isPointInRect,
    }

    global.ComiGoInteraction = api
    global.ComiGoFlip = global.ComiGoFlip || {}
    global.ComiGoFlip.interaction = api
    // 如果是 CommonJS 模块系统，则导出 api
    // 让 Node/Bun 测试环境也能 require() 这个模块，方便测试。
    if (typeof module !== 'undefined' && module.exports) {
        module.exports = api
    }
})(typeof window !== 'undefined' ? window : globalThis)
