// 翻页模式交互几何判断工具（可复用、可测试）
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

    global.ComiGoFlip = global.ComiGoFlip || {}
    global.ComiGoFlip.interaction = api

    if (typeof module !== 'undefined' && module.exports) {
        module.exports = api
    }
})(typeof window !== 'undefined' ? window : globalThis)
