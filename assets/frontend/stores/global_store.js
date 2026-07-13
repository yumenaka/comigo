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
const comigoPath = (path) => (window.ComiGoPath ? window.ComiGoPath(path) : path);
const comigoRelativePath = (pathname) => (window.ComiGoRelativePath ? window.ComiGoRelativePath(pathname) : (pathname || window.location.pathname || '/'));
// random 是前端选择态，不直接作为 daisyUI 的 data-theme。
const randomThemeName = 'random';
// 持久化值可能为空或历史异常值，统一转字符串避免 toString 抛错。
const themeToString = (theme) => (theme === undefined || theme === null ? '' : theme.toString());
const url = new URL(window.location.href);
const currentRelativePath = comigoRelativePath(url.pathname);
const currentRemoteStore = url.searchParams.get('remote_store') || '';
// 运行环境状态集中在这里计算，模板只读取 store，避免各处重复解析 URL。
const wailsBook = window.ComiGoIsWails ? window.ComiGoIsWails() : url.protocol === 'wails:';
const serverReachable = wailsBook || url.protocol === 'http:' || url.protocol === 'https:';
const localBook = !wailsBook && (url.protocol === 'file:' || url.protocol === 'content:');
const staticHtmlBook = window.location.toString().endsWith('.html');
const readerPage = window.ComiGoReaderMode || currentRelativePath.includes('/reader');
const onlineBook = !readerPage && serverReachable && !staticHtmlBook;

if (window.ComiGoForceRandomTheme) {
    try {
        // Alpine Persist 当前使用无前缀 key；保留旧前缀 key 兼容历史数据。
        localStorage.setItem('global.theme', JSON.stringify(randomThemeName));
        localStorage.setItem('_x_global.theme', JSON.stringify(randomThemeName));
    } catch (error) {
        console.warn('无法强制保存随机模板设置:', error);
    }
}

// setURLQueryParam 给站内资源 URL 设置查询参数，返回仍可交给 ComiGoPath 处理的相对 URL。
function setURLQueryParam(rawURL, key, value) {
    if (!rawURL || value === undefined || value === null || value === '') {
        return rawURL;
    }
    try {
        const url = new URL(rawURL, window.location.origin);
        url.searchParams.set(key, String(value));
        if (url.origin !== window.location.origin) {
            return url.href;
        }
        return `${comigoRelativePath(url.pathname)}${url.search}${url.hash}`;
    } catch (error) {
        const separator = rawURL.includes('?') ? '&' : '?';
        return `${rawURL}${separator}${encodeURIComponent(key)}=${encodeURIComponent(value)}`;
    }
}

// 生成userID: 使用UserAgent的哈希值 + 随机字符串，确保唯一性且长度适中
const initClientID = `Client_${randomString}_${system}_${browser}`;
Alpine.store('global', {
    nowPageNum: 1,
    allPageNum: 1,
    // 在线书籍模式：可访问后端，且不是本地 reader 或静态 HTML。
    onlineBook: onlineBook,
    // 本地便携模式：file:// 或 Android content:// 打开。
    localBook: localBook,
    // 静态 HTML 导出模式：用于控制便携 HTML 标题等显示。
    staticHtmlBook: staticHtmlBook,
    // 当前页面是否可访问 HTTP 后端能力，例如二维码和阅读历史。
    serverReachable: serverReachable,
    // Wails 桌面壳使用自定义协议，但资源仍由内嵌服务处理。
    wailsBook: wailsBook,
    // 播放器：音量（0~100）
    playerVolume: Alpine.$persist(100).as('global.playerVolume'),
    // 播放器：是否静音
    playerMuted: Alpine.$persist(false).as('global.playerMuted'),
    // 播放器：是否自动播放下一曲
    autoPlayNext: Alpine.$persist(true).as('global.autoPlayNext'),
    // 播放器：是否循环播放播放列表
    loopPlaylist: Alpine.$persist(true).as('global.loopPlaylist'),
    // 自动切边
    autoCrop: Alpine.$persist(false).as('global.autoCrop'),
    // 自动切边阈值,范围是0~100。多数情况下 1 就够了。
    autoCropNum: Alpine.$persist(1).as('global.autoCropNum'),
    // 是否压缩图片
    autoResize: Alpine.$persist(false).as('global.autoResize'),
    // 压缩图片限宽
    autoResizeWidth: Alpine.$persist(800).as('global.autoResizeWidth'),
    // 初始主题
    theme: Alpine.$persist('retro').as('global.theme'),
    // 随机模板实际解析出的主题，选择态依然保留为 random。
    randomResolvedTheme: Alpine.$persist('').as('global.randomResolvedTheme'),
    // 本次页面加载是否已经解析过随机模板；不持久化，刷新或跳转后会重新随机。
    randomThemeResolvedThisPage: false,
    // custom 主题：组件颜色
    customBase100: Alpine.$persist('#dce6ff').as('global.customBase100'),
    // custom 主题：背景颜色
    customBase300: Alpine.$persist('#076c0a').as('global.customBase300'),
    // custom 主题：文字颜色
    customBaseContent: Alpine.$persist('#282425').as('global.customBaseContent'),
    // bgPattern 背景花纹
    bgPattern: Alpine.$persist('grid-line').as('global.bgPattern'),
    // 随机模板池：只包含内置非 custom 主题，不包含 random 本身。
    randomThemeList: ['light', 'dark', 'retro', 'cupcake', 'cyberpunk', 'red-white-game', 'dracula', 'valentine', 'cmyk', 'halloween', 'coffee', 'winter', 'nord'],
    // 需要保留 bg-base-300 的主题名单（例如 custom 主题也要使用该背景层级）
    bgBase300ThemeList: ['light', 'dark', 'retro', 'custom', 'cupcake', 'cyberpunk', 'red-white-game', 'nord'],
    // 自带完整背景的主题会覆盖纯色/网格线花纹选择，相关控件需要隐藏。
    ownBackgroundThemeList: ['cupcake', 'cyberpunk', 'red-white-game', 'dracula', 'valentine', 'cmyk', 'halloween', 'coffee', 'winter', 'nord'],
    // 主题下拉框统一走这里，选择 random 时立即解析一次和当前模板不同的实际主题。
    setTheme(theme) {
        const currentTheme = this.theme === randomThemeName ? themeToString(this.randomResolvedTheme) : themeToString(this.theme);
        this.theme = (theme || '').toString();
        if (this.theme === randomThemeName) {
            this.refreshRandomTheme(currentTheme);
        }
    },
    // 从随机池抽取主题，并排除当前实际主题，避免刷新后仍然是同一个模板。
    pickRandomTheme(currentTheme = '') {
        const candidates = this.randomThemeList.filter((theme) => theme !== currentTheme);
        const pool = candidates.length > 0 ? candidates : this.randomThemeList;
        return pool[Math.floor(Math.random() * pool.length)] || 'cmyk';
    },
    // 重新解析 random 对应的实际主题；排除当前实际主题，避免刷新后仍是同一个模板。
    refreshRandomTheme(excludedTheme = '') {
        const resolvedTheme = themeToString(excludedTheme || this.randomResolvedTheme);
        const currentTheme = this.randomThemeList.includes(resolvedTheme) ? resolvedTheme : '';
        const nextTheme = this.pickRandomTheme(currentTheme);
        this.randomResolvedTheme = nextTheme;
        this.randomThemeResolvedThisPage = true;
        return nextTheme;
    },
    // 页面加载时只解析一次 random；刷新或全页面跳转后会重新抽取。
    ensureRandomTheme() {
        const resolvedTheme = themeToString(this.randomResolvedTheme);
        const isResolvedThemeValid = this.randomThemeList.includes(resolvedTheme);
        if (this.randomThemeResolvedThisPage && isResolvedThemeValid) {
            return resolvedTheme;
        }
        return this.refreshRandomTheme();
    },
    // 返回真正写入 body[data-theme] 的主题；random 只保留为设置项选择值。
    getEffectiveTheme() {
        const selectedTheme = themeToString(this.theme);
        if (selectedTheme !== randomThemeName) {
            return selectedTheme;
        }
        return this.ensureRandomTheme();
    },
    canSelectBgPattern() {
        return !this.ownBackgroundThemeList.includes(this.getEffectiveTheme());
    },
    /**
     * 返回主区域背景类名：统一处理背景花纹和 bg-base-300 的组合逻辑
     * @returns {string} 例如 "grid-line bg-base-300" / "bg-base-300" / "grid-line" / ""
     */
    getMainAreaBgClass() {
        const classes = [];
        if (this.canSelectBgPattern() && this.bgPattern !== 'none') {
            classes.push(this.bgPattern);
        }
        if (this.bgBase300ThemeList.includes(this.getEffectiveTheme())) {
            classes.push('bg-base-300');
        }
        return classes.join(' ');
    },
    // 是否禁止图片接口缓存；阅读页会把该状态转换为 no-cache 查询参数。
    noCache: Alpine.$persist(false).as('global.noCache'),
    // clientID 用于识别匿名用户与设备
    clientID: Alpine.$persist(initClientID).as('global.clientID'),
    // debugMode 是否开启调试模式
    debugMode: Alpine.$persist(true).as('global.debugMode'),
    // 是否通过 WebSocket 同步阅读页码
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
            const key = `pageNum_${book_id}`;
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
    savePageNumToLocalStorage(book_id) {
        if (!this.saveReadingProgress) {
            return;
        }
        if (!book_id) {
            console.warn("savePageNumToLocalStorage: book_id is required");
            return;
        }
        try {
            const key = `pageNum_${book_id}`;
            const nowPageNum = Alpine.store('global').nowPageNum;
            localStorage.setItem(key, nowPageNum);
        } catch (e) {
            console.error("Error saving page number to localStorage:", e);
        }
    },
    // 当前阅读模式：scroll=卷轴阅读，flip=翻页阅读
    readMode: Alpine.$persist('scroll').as('global.readMode'),
    // 切换为卷轴阅读，并使用无限卷轴加载策略。
    infiniteScrollLoadAllPage() {
        Alpine.store('scroll').loadMode = 'infinite';
        this.readMode = "scroll";
        this.onChangeReadMode();
    },
    onChangeReadMode() {
        // 切换阅读模式时，如果在阅读，就修改URL路径 参考文献：https://developer.mozilla.org/zh-CN/docs/Web/API/URL
        const url = new URL(window.location.href);
        const pathname = comigoRelativePath(url.pathname);
        // 分割路径为各层级关键词, filter(Boolean) 的作用是去除空字符串 如//aa/bb/ 会产生空字符串(虽然这里不会这么做)
        const pathSegments = pathname.split('/').filter(Boolean); // like ["scroll", "id3DcA1v9"]
        const book_id = pathSegments[pathSegments.length - 1];
        console.log(`切换阅读模式到: ${this.readMode}, 当前路径: ${pathname},${pathSegments}`);
        // 跳转到新的阅读模式URL
        if (pathSegments.includes("scroll")||pathSegments.includes("flip")) {
            window.location.href = this.getReadURL(book_id, this.nowPageNum);
        }
    },
    getReadURL(book_id, start_index, remote_store = '') {
        const url = new URL(window.location.href);
        const pageNum = Math.max(1, parseInt(start_index, 10) || 1);
        const remoteStore = remote_store || currentRemoteStore;
        // 翻页阅读
        if (this.readMode === 'flip') {
            let new_url = new URL(comigoPath(`/flip/${book_id}`), url.origin);
            new_url.searchParams.set("page", pageNum.toString());
            if (remoteStore) {
                new_url.searchParams.set("remote_store", remoteStore);
            }
            return new_url.href;
        }
        // 卷轴阅读
        if (this.readMode === 'scroll') {
            let new_url = new URL(comigoPath(`/scroll/${book_id}`), url.origin);
            const scrollStore = Alpine.store('scroll');
            const loadMode = ['infinite', 'lazy', 'paged'].includes(scrollStore.loadMode) ? scrollStore.loadMode : 'infinite';
            const pageLimit = Math.max(1, parseInt(scrollStore.pageLimit, 10) || 32);
            // page 始终表示精确书页；limit 只用于标识并计算分页加载块。
            new_url.searchParams.set("page", pageNum.toString());
            if (loadMode === 'paged') {
                new_url.searchParams.set("limit", pageLimit.toString());
            }
            if (remoteStore) {
                new_url.searchParams.set("remote_store", remoteStore);
            }
            return new_url.href;
        }
        return "";
    },
    // getCoverURL 统一生成封面 URL，所有调用方都显式传入展示尺寸，避免不同尺寸共用同一个后端缓存。
    getCoverURL(bookInfo, resizeHeight = 352) {
        const rawCoverURL = bookInfo?.cover?.url || (bookInfo?.id ? `/api/get-cover?id=${encodeURIComponent(bookInfo.id)}` : "");
        if (!rawCoverURL) {
            return "";
        }
        const isResizableCover = rawCoverURL.includes("/api/get-file") || rawCoverURL.includes("/api/get-cover");
        if (!isResizableCover) {
            return comigoPath(rawCoverURL);
        }

        let coverURL = setURLQueryParam(rawCoverURL, "resize_height", resizeHeight);
        coverURL = setURLQueryParam(coverURL, "remote_store", bookInfo?.remote_store || currentRemoteStore);
        return comigoPath(coverURL);
    },
    // 竖屏模式
    isPortrait: false,
    // 横屏模式
    isLandscape: true,
    // 获取cookie里面存储的值
    getCookieValue(bookID, valueName) {
        let pgCookie = "";
        const paramName = (bookID === "" ? `$${valueName}` : `${bookID}_${valueName}`);
        const cookiePrefix = `${paramName}=`;
        const cookies = document.cookie.split(";");
        for (let i = 0; i < cookies.length; i++) {
            const cookie = cookies[i].trim();
            if (cookie.startsWith(cookiePrefix)) {
                pgCookie = decodeURIComponent(cookie.substring(cookiePrefix.length));
            }
        }
        return pgCookie;
    },
    setPaginationIndex(bookID, valueName, value) {
        const paramName = (bookID === "" ? `$${valueName}` : `${bookID}_${valueName}`);
        // 设置cookie，过期时间为365天
        const expirationDate = new Date();
        expirationDate.setDate(expirationDate.getDate() + 365);
        document.cookie = `${paramName}=${encodeURIComponent(value)}; expires=${expirationDate.toUTCString()}; path=/; SameSite=Lax`;
        window.location.reload();
    },
    /**
     * 调用后端 /api/store-bookmark 接口，更新书签信息
     * @param {Object} params
     * @param {string} params.type - 书签类型，例如 'auto'
     * @param {string} params.bookId - 书籍ID
     * @param {number} params.pageIndex - 页码（1 起始）
     * @param {string} [params.label='自动书签'] - 书签名称，当前后端固定为自动书签，仅用于日志
     * @returns {Promise<Object|string>} 后端返回的响应体
     */
    async UpdateBookmark({ type = 'auto', bookId, pageIndex, description = '' } = {}) {
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
        if (description === '') {
            description = `${browser} in ${system}`;
        }
        const payload = {
            type,
            book_id: bookId,
            page_index: pageIndex,
            description: description
        };
        let bookmarkURL = '/api/store-bookmark';
        const remoteStore = currentRemoteStore;
        if (remoteStore) {
            bookmarkURL = setURLQueryParam(bookmarkURL, 'remote_store', remoteStore);
        }
        const response = await fetch(comigoPath(bookmarkURL), {
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

// 旧版本可能留下非 scroll/flip 的 global.readMode，统一回落到卷轴阅读。
if (!['scroll', 'flip'].includes(Alpine.store('global').readMode)) Alpine.store('global').readMode = 'scroll';

// 初始化全局存储
document.addEventListener('alpine:initialized', () => {
    Alpine.store('global').init();
});

if (currentRelativePath.includes('/flip/')) {
    Alpine.store('global').readMode = 'flip';
} else if (currentRelativePath.includes('/scroll/')) {
    Alpine.store('global').readMode = 'scroll';
}
