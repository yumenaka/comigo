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
    // readerMode 当前阅读模式
    readMode: Alpine.$persist('scroll').as('global.readMode'),
    // readerMode 当前阅读模式
    readModeIsScroll: Alpine.$persist(true).as('global.readModeIsScroll'),
    //是否通过websocket同步翻页
    syncPageByWS: Alpine.$persist(true).as('global.syncPageByWS'),
    // bookSortBy 书籍排序方式 以按照文件名、修改时间、文件大小排序（或反向排序）
    bookSortBy: Alpine.$persist('name').as('global.bookSortBy'),
    // pageSortBy 书页排序顺序 以按照文件名、修改时间、文件大小排序（或反向排序）
    pageSortBy: Alpine.$persist('name').as('global.pageSortBy'),
    language: Alpine.$persist('en').as('global.language'),
    toggleReadMode() {
      if (this.readMode === 'flip') {
        this.readMode = 'scroll'
        this.readModeIsScroll = true
      } else {
        this.readMode = 'flip'
        this.readModeIsScroll = false
      }
    },
    // 竖屏模式
    isPortrait: false,
    // 横屏模式
    isLandscape: true,
    // 获取cookie里面存储的值
    getCookieValue(bookID,valueName) {
      let pgCookie = "";
      const paramName = (bookID === ""?`$${valueName}`:`${bookID}_${valueName}`);
      const cookies = document.cookie.split(";");
      for (let i = 0; i < cookies.length; i++) {
        const cookie = cookies[i].trim();
        if (cookie.startsWith(paramName)) {
          pgCookie = decodeURIComponent(cookie.substring(paramName.length + 1));
        }
      }
      return pgCookie;
    },
    setPaginationIndex(bookID, valueName,value) {
      const paramName = (bookID === ""?`$${valueName}`:`${bookID}_${valueName}`);
      // 设置cookie，过期时间为365天
      const expirationDate = new Date();
      expirationDate.setDate(expirationDate.getDate() + 365);
      document.cookie = `${paramName}${encodeURIComponent(value)}; expires=${expirationDate.toUTCString()}; path=/; SameSite=Lax`;
      window.location.reload();
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