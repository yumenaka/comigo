// 用Alpine Persist 注册全局变量
// https://alpinejs.dev/plugins/persist#using-alpine-persist-global

/**
 * 解析UserAgent获取浏览器信息
 * @returns {string} 浏览器名称
 */
function getBrowserInfo() {
    const ua = navigator.userAgent;
    let browser = 'Unknown';
    if (ua.indexOf('Firefox') > -1) {
        browser = 'Firefox';
    } else if (ua.indexOf('Edg') > -1) {
        browser = 'Edge';
    } else if (ua.indexOf('Chrome') > -1) {
        browser = 'Chrome';
    } else if (ua.indexOf('Safari') > -1) {
        browser = 'Safari';
    } else if (ua.indexOf('Opera') > -1 || ua.indexOf('OPR') > -1) {
        browser = 'Opera';
    } else if (ua.indexOf('Trident') > -1 || ua.indexOf('MSIE') > -1) {
        browser = 'IE';
    }
    return browser;
}

/**
 * 解析UserAgent获取系统信息
 * @returns {string} 系统名称
 */
function getSystemInfo() {
    const ua = navigator.userAgent;
    let os = 'Unknown';
    if (ua.indexOf('Win') > -1) {
        os = 'Windows';
    } else if (ua.indexOf('Mac') > -1) {
        os = 'MacOS';
    } else if (ua.indexOf('Linux') > -1) {
        os = 'Linux';
    } else if (ua.indexOf('Android') > -1) {
        os = 'Android';
    } else if (ua.indexOf('iOS') > -1 || ua.indexOf('iPhone') > -1 || ua.indexOf('iPad') > -1) {
        os = 'iOS';
    }
    return os;
}

/**
 * 生成随机字符串
 * @returns {string} 随机字符串
 */
function generateRandomString() {
    return (Date.now() % 10000000).toString(36) + Math.random().toString(36).substring(2, 5);
}

// 浏览器 系统信息 随机字符串
const browser = getBrowserInfo();
const system = getSystemInfo();
const randomString = generateRandomString();
// 生成userID: 使用UserAgent的哈希值 + 随机字符串，确保唯一性且长度适中
const initClientID = `Client_${randomString}_${system}_${browser}`;
Alpine.store('global', {
    nowPageNum: 1,
    allPageNum: 1,
    // 自动切边
    autoCrop: Alpine.$persist(false).as('global.autoCrop'),
    // 自动切边阈值,范围是0~100。多数情况下 1 就够了。
    autoCropNum: Alpine.$persist(1).as('global.autoCropNum'),
    // 是否压缩图片
    autoResize: Alpine.$persist(false).as('global.autoResize'),
    // 压缩图片限宽
    autoResizeWidth: Alpine.$persist(800).as('global.autoResizeWidth'),
    // bgPattern 背景花纹
    bgPattern: Alpine.$persist('grid-line').as('global.bgPattern'),
    // 是否禁止缓存（TODO：缓存功能优化与测试）
    noCache: Alpine.$persist(false).as('global.noCache'),
    // clientID 用于识别匿名用户与设备
    clientID: Alpine.$persist(initClientID).as('global.clientID'),
    // debugMode 是否开启调试模式
    debugMode: Alpine.$persist(true).as('global.debugMode'),
    //是否通过websocket同步翻页
    syncPageByWS: Alpine.$persist(true).as('global.syncPageByWS'),
    // bookSortBy 书籍排序方式 以按照文件名、修改时间、文件大小排序（或反向排序）
    bookSortBy: Alpine.$persist('name').as('global.bookSortBy'),
    // pageSortBy 书页排序顺序 以按照文件名、修改时间、文件大小排序（或反向排序）
    pageSortBy: Alpine.$persist('name').as('global.pageSortBy'),
    language: Alpine.$persist('en').as('global.language'),
    //是否保存阅读进度（页数）到本地存储
    saveReadingProgress: Alpine.$persist(true).as('global.saveReadingProgress'),
    // 从本地存储加载页码并跳转
    loadPageNumFromLocalStorage(book_id, callbackFunction) {
        if (!this.saveReadingProgress) {
            return;
        }
        try {
            const key = `pageNum_${book.id}`;
            const savedPageNum = localStorage.getItem(key);
            if (savedPageNum !== null && !isNaN(parseInt(savedPageNum))) {
                const pageNum = parseInt(savedPageNum);
                // 确保页码在有效范围内
                if (pageNum > 0 && pageNum <= Alpine.store('global').allPageNum) {
                    console.log(`加载到本地存储的页码: ${pageNum}`);
                    callbackFunction(); // 跳转函数,或发送书签更新信息
                }
            }
        } catch (e) {
            console.error("Error loading page number from localStorage:", e);
        }
    },
    // 保存当前页码到本地存储
    savePageNumToLocalStorage() {
        if (!this.saveReadingProgress) {
            return;
        }
        try {
            const key = `pageNum_${book.id}`;
            const nowPageNum = Alpine.store('global').nowPageNum;
            localStorage.setItem(key, nowPageNum);
        } catch (e) {
            console.error("Error saving page number to localStorage:", e);
        }
    },
    // readerMode 当前阅读模式: infinite_scroll  paged_scroll  flip_page
    readMode: Alpine.$persist('infinite_scroll').as('global.readMode'),
    // 切换阅读模式
    infiniteScrollLoadAllPage(mode) {
        this.readMode = "infinite_scroll";
        const url = new URL(window.location.href);
        url.searchParams.delete("start");
        window.location.href = url.href;
    },
    onChangeReadMode() {
        // 切换阅读模式时，如果在阅读，就修改URL路径 参考文献：https://developer.mozilla.org/zh-CN/docs/Web/API/URL
        const url = new URL(window.location.href);
        const pathname = url.pathname;
        // 使用 URLSearchParams 提取键值对
        const params = new URLSearchParams(url.search);
        // 分割路径为各层级关键词, filter(Boolean) 的作用是去除空字符串 如//aa/bb/ 会产生空字符串(虽然这里不会这么做)
        const pathSegments = url.pathname.split('/').filter(Boolean); // like ["scroll", "id3DcA1v9"]
        const book_id = pathSegments[1];
        console.log(`切换阅读模式到: ${this.readMode}, 当前路径: ${pathname},${pathSegments}, 查询参数: ${params.toString()}`);
        // 翻页模式
        if (this.readMode === 'page_flip') {
            // 如果已经是翻页模式
            if (pathSegments[0] === "flip") {
                console.log("已经是翻页模式，无需切换");
                console.log(`${pathSegments[0]} , ${params.get("start")}`);
                return;
            }
        }
        // 卷轴(分页)模式
        if (this.readMode === 'paged_scroll') {
            // 如果已经是分页卷轴模式
            if (pathSegments[0] === "scroll" && params.get("page") !== null) {
                console.log(`${pathSegments[0]} , ${params.get("page")}`);
                console.log("已经是分页卷轴模式，无需切换");
                return;
            }
        }
        // 卷轴(无限)模式
        if (this.readMode === 'infinite_scroll') {
            // 如果已经是无限卷轴模式
            if (pathSegments[0] === "scroll" && params.get("page") === null) {
                console.log(`${pathSegments[0]} , ${params.get("page")}`);
                console.log("已经是无限卷轴模式，无需切换");
                return;
            }
        }
        // 跳转到新的阅读模式URL
        window.location.href = this.getReadURL(book_id, this.nowPageNum);
    },
    getReadURL(book_id, start_index) {
        // TODO: 处理旧版本数据干扰的问题。若干个版本后大概就不需要了，到时候删除这段代码。
        if (this.readMode !== 'page_flip'&& this.readMode !== 'paged_scroll' && this.readMode !== 'infinite_scroll') {
            console.error(`未知的阅读模式: ${this.readMode}, 可能是旧版本数据干扰, 重置为 infinite_scroll`);
            this.readMode = 'infinite_scroll';
        }
        let PAGED_SIZE = 32;
        // console.log(`生成阅读模式URL: ${this.readMode}`);
        // console.log(`当前页码: ${start_index}`);
        const url = new URL(window.location.href);
        // 翻页(左右)
        if (this.readMode === 'page_flip') {
            let new_url = new URL(`/flip/${book_id}`, url.origin);
            if (start_index > 1) {
                new_url.searchParams.set("start", start_index.toString());
            }
            return new_url.href;
        }
        // 卷轴(分页)
        if (this.readMode === 'paged_scroll') {
            let new_url = new URL(`/scroll/${book_id}`, url.origin);
            let page = Math.floor(start_index / PAGED_SIZE) + 1;
            new_url.searchParams.set("page", page.toString());
            return new_url.href;
        }
        // 卷轴(无限)
        if (this.readMode === 'infinite_scroll') {
            let new_url = new URL(`/scroll/${book_id}`, url.origin);
            if (start_index > PAGED_SIZE) {
                new_url.searchParams.set("start", start_index.toString());
            }
            return new_url.href;
        }
        return "";
    },
    // 竖屏模式
    isPortrait: false,
    // 横屏模式
    isLandscape: true,
    // 获取cookie里面存储的值
    getCookieValue(bookID, valueName) {
        let pgCookie = "";
        const paramName = (bookID === "" ? `$${valueName}` : `${bookID}_${valueName}`);
        const cookies = document.cookie.split(";");
        for (let i = 0; i < cookies.length; i++) {
            const cookie = cookies[i].trim();
            if (cookie.startsWith(paramName)) {
                pgCookie = decodeURIComponent(cookie.substring(paramName.length + 1));
            }
        }
        return pgCookie;
    },
    setPaginationIndex(bookID, valueName, value) {
        const paramName = (bookID === "" ? `$${valueName}` : `${bookID}_${valueName}`);
        // 设置cookie，过期时间为365天
        const expirationDate = new Date();
        expirationDate.setDate(expirationDate.getDate() + 365);
        document.cookie = `${paramName}${encodeURIComponent(value)}; expires=${expirationDate.toUTCString()}; path=/; SameSite=Lax`;
        window.location.reload();
    },
    /**
     * 调用后端 /api/store_bookmark 接口，更新书签信息
     * @param {Object} params
     * @param {string} params.type - 书签类型，例如 'auto'
     * @param {string} params.bookId - 书籍ID
     * @param {number} params.pageIndex - 页码（1 起始）
     * @param {string} [params.label='自动书签'] - 书签名称，当前后端固定为自动书签，仅用于日志
     * @returns {Promise<Object|string>} 后端返回的响应体
     */
    async UpdateBookmark({ type = 'auto', bookId, pageIndex, label = '自动书签' } = {}) {
        if (!bookId) {
            const error = new Error('UpdateBookmark: bookId is required');
            if (this.debugMode) {
                console.error(error);
            }
            throw error;
        }
        if (!Number.isInteger(pageIndex) || pageIndex <= 0) {
            const error = new Error('UpdateBookmark: pageIndex must be a positive integer');
            if (this.debugMode) {
                console.error(error);
            }
            throw error;
        }

        const deviceDescription = `${browser} in ${system}`;
        const payload = {
            type,
            book_id: bookId,
            page_index: pageIndex,
            description: deviceDescription
        };
        const response = await fetch('/api/store_bookmark', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify(payload)
        });

        const contentType = response.headers.get('content-type') || '';
        const isJSON = contentType.includes('application/json');
        const responseBody = isJSON ? await response.json() : await response.text();
        if (!response.ok) {
            const error = new Error(`UpdateBookmark failed: ${response.status} ${response.statusText}`);
            if (this.debugMode) {
                console.error('[UpdateBookmark] error', error, responseBody);
            }
            throw error;
        }
        return responseBody;
    },
    // 检测并设置视口方向
    checkOrientation() {
        const isPortrait = window.innerHeight > window.innerWidth;
        this.isPortrait = isPortrait;
        this.isLandscape = !isPortrait;
        //console.log(`当前视口方向: ${isPortrait ? '竖屏' : '横屏'}`);
    },
    // 初始化方法
    init() {
        // 设置初始方向
        this.checkOrientation();
        // 添加视口变化监听
        window.addEventListener('resize', () => {
            this.checkOrientation();
        });
    }
})

// 初始化全局存储
document.addEventListener('alpine:initialized', () => {
    Alpine.store('global').init();
});
