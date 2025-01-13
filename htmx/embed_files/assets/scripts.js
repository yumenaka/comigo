import 'htmx.org'
import Alpine from 'alpinejs'
import persist from '@alpinejs/persist'
import i18next from 'i18next'
import 'flowbite'
import LanguageDetector from 'i18next-browser-languagedetector'
import 'tw-colors'
import morph from '@alpinejs/morph'
// i18next 国际化插件，用于国际化。详细用法参见：
// https://www.i18next.com/overview/getting-started
import enLocale from './locales/en_US.json'
import zhLocale from './locales/zh_CN.json'
import jaLocale from './locales/ja_JP.json'
import screenfull from 'screenfull'
// 将 Alpine 实例添加到窗口对象中。
window.Alpine = Alpine

// Alpine Persist 插件，用于持久化存储。默认存储到 localStorage。
// 详细用法参见： https://alpinejs.dev/plugins/persist
Alpine.plugin(persist)
Alpine.plugin(morph)

i18next
    .use(LanguageDetector)
    .init({
        debug: false,
        // // 在 setTimeout（默认异步行为）内的 init（） 中触发资源加载。如果您的后端同步加载资源，请将其设置为 false - 这样，
        // // 可以在 init（） 之后调用 i18next.t（） 而无需依赖初始化回调。此选项仅适用于同步（阻塞）加载后端，例如 i18next-fs-backend 和 i18next-sync-fs-backend！
        initImmediate: true,
        //lng: 'en', // if you're using a language detector, do not define the lng option
        // supportedLngs: ['en', 'cn', 'ja'],
        supportedLngs: ['en-US', 'ja-JP', 'zh-CN', 'en', 'zh', 'ja'],
        fallbackLng: ['en', 'zh', 'ja'],
        resources: {
            'en-US': {
                translation: enLocale,
            },
            en: {
                translation: enLocale,
            },
            'zh-CN': {
                translation: zhLocale,
            },
            zh: {
                translation: zhLocale,
            },
            'ja-JP': {
                translation: jaLocale,
            },
            ja: {
                translation: jaLocale,
            },
        },
    })
    .then(function (t) {
        //console.log(t('test'))
        // i18next.changeLanguage('en', (err, t) => {
        //     if (err) return console.log('something went wrong loading', err);
        //     console.log(t('test'));
        // });
    })

window.i18next = i18next // 使i18next在全局作用域中可用
if(document.getElementById('FullScreenIcon')){
    document.getElementById('FullScreenIcon').addEventListener('click', () => {
        if (screenfull.isEnabled) {
            screenfull.toggle()
        } else {
            // Ignore or do something else
            i18next.t('not_support_fullscreen')
        }
    })
}


// 用Alpine Persist 注册全局变量
// https://alpinejs.dev/plugins/persist#using-alpine-persist-global

// global 全局设置
Alpine.store('global', {
    // bgPattern 背景花纹
    bgPattern: Alpine.$persist('grid-line').as('global.bgPattern'),
    autoCrop: Alpine.$persist(false).as('global.autoCrop'),
    autoCropNum: Alpine.$persist(1).as('global.autoCropNum'), // 自动切白边阈值,范围是0~100。大多数情况下 1 就够了。
    // userID 当前用户ID  用于同步阅读进度 随机生成
    userID: Alpine.$persist(Math.random().toString(36).substring(2)).as('global.userID'),
    // debugMode 是否开启调试模式
    debugMode: Alpine.$persist(true).as('global.debugMode'),
    // readerMode 当前阅读模式
    readMode: Alpine.$persist('scroll').as('global.readMode'),
    //是否通过websocket同步翻页
    syncPageByWS: Alpine.$persist(true).as('global.syncPageByWS'),
    // bookSortBy 书籍排序方式 以按照文件名、修改时间、文件大小排序（或反向排序）
    bookSortBy: Alpine.$persist('name').as('global.bookSortBy'),
    // pageSortBy 书页排序顺序 以按照文件名、修改时间、文件大小排序（或反向排序）
    pageSortBy: Alpine.$persist('name').as('global.pageSortBy'),
    language: Alpine.$persist('en').as('global.language'),
    toggleReadMode() {
        this.readMode = this.readMode === 'flip' ? 'scroll' : 'flip'
    },
})

// BookShelf 书架设置
Alpine.store('shelf', {
    bookCardMode: Alpine.$persist('gird').as('shelf.bookCardMode'), //gird,list,text
    showTitle: Alpine.$persist(true).as('shelf.showTitle'), //是否显示标题
    showFileIcon: Alpine.$persist(true).as('shelf.showFileIcon'), //是否显示文件图标
    simplifyTitle: Alpine.$persist(true).as('shelf.simplifyTitle'), //是否简化标题
    InfiniteDropdown: Alpine.$persist(false).as('shelf.InfiniteDropdown'), //卷轴模式下，是否无限下拉
    bookCardShowTitleFlag: Alpine.$persist(true).as('shelf.bookCardShowTitleFlag'), // 书库中的书籍是否显示文字版标题
    syncScrollFlag: false, // 同步滚动,目前还没做
    // 屏幕宽横比,inLandscapeMode的判断依据
    aspectRatio: 1.2,
    // 可见范围宽高的具体值
    clientWidth: 0,
    clientHeight: 0,
})

// Scroll 卷轴模式
Alpine.store('scroll', {
    nowPageNum: 1,
    simplifyTitle: Alpine.$persist(true).as('scroll.simplifyTitle'), //是否简化标题
    //下拉模式下，漫画页面的底部间距。单位px。
    marginBottomOnScrollMode: Alpine.$persist(10).as(
        'scroll.marginBottomOnScrollMode'
    ),
    //卷轴模式下，是否无限下拉
    InfiniteDropdown: Alpine.$persist(true).as('scroll.InfiniteDropdown'),
    syncScrollFlag: Alpine.$persist(false).as('scroll.syncScrollFlag'), // 同步滚动,目前还没做
    imageMaxWidth: 400,
    // 屏幕宽横比,inLandscapeMode的判断依据
    aspectRatio: 1.2,
    // 可见范围宽高的具体值
    clientWidth: 0,
    clientHeight: 0,
    //漫画页的单位,是否使用固定值
    widthUseFixedValue: Alpine.$persist(true).as('scroll.widthUseFixedValue'),
    //横屏(Landscape)状态的漫画页宽度,百分比
    singlePageWidth_Percent: Alpine.$persist(60).as('scroll.singlePageWidth_Percent'),
    doublePageWidth_Percent: Alpine.$persist(95).as('scroll.doublePageWidth_Percent'),
    //横屏(Landscape)状态的漫画页宽度。px。
    singlePageWidth_PX: Alpine.$persist(720).as('scroll.singlePageWidth_PX'),
    doublePageWidth_PX: Alpine.$persist(1200).as('scroll.doublePageWidth_PX'),
    //书籍数据,需要从远程拉取
    //是否显示顶部页头
    showHeaderFlag: true,
    //是否显示页数
    show_page_num: Alpine.$persist(false).as('scroll.show_page_num'),
    //ws翻页相关
    syncPageByWS: Alpine.$persist(false).as('scroll.syncPageByWS'), //是否通过websocket同步翻页
    // //此处修改不会实时生效，不要这么做
    // toggleSimplifyTitle() {
    //     this.simplifyTitle = ! this.simplifyTitle
    // }
})

// Flip 翻页模式
Alpine.store('flip', {
    nowPageNum: 1,
    allPageNum: 100,
    imageMaxWidth: 400,
    isLandscapeMode: true,
    isPortraitMode: false,
    //自动隐藏工具条
    autoHideToolbar: Alpine.$persist(true).as('flip.autoHideToolbar'),
    //是否显示页头
    show_header: Alpine.$persist(true).as('flip.show_header'),
    //是否显示页脚
    showFooter: Alpine.$persist(true).as('flip.showFooter'),
    //是否显示页数
    show_page_num: Alpine.$persist(false).as('flip.show_page_num'),
    //是否是右半屏翻页（从右到左）?日本漫画从左到右(false)
    rightToLeft: Alpine.$persist(false).as('flip.rightToLeft'),
    //双页模式
    doublePageMode: Alpine.$persist(false).as('flip.doublePageMode'),
    //自动拼合双页(TODO)
    autoDoublePageMode: Alpine.$persist(false).as(
        'flip.autoDoublePageModeFlag'
    ),
    //是否保存阅读进度（页数）
    saveReadingProgress: Alpine.$persist(true).as('flip.saveReadingProgress'),
    //素描模式标记
    sketchModeFlag: false,
    //是否显示素描提示
    showPageHint: Alpine.$persist(false).as(
        'flip.showPageHint'
    ),
    //翻页间隔时间
    sketchFlipSecond: 30,
    //计时用,从0开始
    sketchSecondCount: 0,
})

// 自定义主题
Alpine.store('theme', {
    theme: Alpine.$persist('light').as('theme'),
    interfaceColor: '#F5F5E4',
    backgroundColor: '#E0D9CD',
    textColor: '#000000',
    toggleTheme() {
        this.theme = this.theme === 'light' ? 'dark' : 'light'
    },
})


// 由于 Cookie “cookie.someCookieKey”缺少正确的“sameSite”属性值，缺少“SameSite”或含有无效值的 Cookie
// 即将被视作指定为“Lax”，该 Cookie 将无法发送至第三方上下文中。若您的应用程序依赖这组 Cookie 以在不同上下文中工作，
// 请添加“SameSite=None”属性。若要了解“SameSite”属性的更多信息，请参阅：https://developer.mozilla.org/docs/Web/HTTP/Headers/Set-Cookie/SameSite

// https://alpinejs.dev/plugins/persist#custom-storage
// 定义自定义存储对象，公开 getItem 函数和 setItem 函数
// 使用会话 cookie 作为存储
window.cookieStorage = {
    getItem(key) {
        let cookies = document.cookie.split(";");
        for (let i = 0; i < cookies.length; i++) {
            let cookie = cookies[i].split("=");
            if (key === cookie[0].trim()) {
                return decodeURIComponent(cookie[1]);
            }
        }
        return null;
    },
    setItem(key, value) {
        document.cookie = `${key}=${encodeURIComponent(value)}; SameSite=Lax`;//SameSite设置默认值（Lax），防止控制台报错。加载图像或框架（frame）的请求将不会包含用户的 Cookie。
    }
}
// 使用 cookieStorage 作为存储
Alpine.store('cookie', {
    someCookieKey: Alpine.$persist(false).using(cookieStorage).as('cookie.someCookieKey'),
})


//请求图片文件时，可添加的额外参数
const imageParameters = {
    resize_width: -1, // 缩放图片,指定宽度
    resize_height: -1, // 指定高度,缩放图片
    do_auto_resize: false,
    resize_max_width: 800, //图片宽度大于这个上限时缩小
    resize_max_height: -1, //图片高度大于这个上限时缩小
    do_auto_crop: false,
    auto_crop_num: 1, // 自动切白边阈值,范围是0~100,其实为1就够了
    gray: false, //黑白化
};

//添加各种字符串参数,不需要的话为空
const resize_width_str =
    imageParameters.resize_width > 0
        ? "&resize_width=" + imageParameters.resize_width
        : "";
const resize_height_str =
    imageParameters.resize_height > 0
        ? "&resize_height=" + imageParameters.resize_height
        : "";
const gray_str = imageParameters.gray ? "&gray=true" : "";
const do_auto_resize_str = imageParameters.do_auto_resize
    ? "&resize_max_width=" + imageParameters.resize_max_width
    : "";
const resize_max_height_str =
    imageParameters.resize_max_height > 0
        ? "&resize_max_height=" + imageParameters.resize_max_height
        : "";
const auto_crop_str = imageParameters.do_auto_crop
    ? "&auto_crop=" + imageParameters.auto_crop_num
    : "";
//所有附加的转换参数
let addStr =
    resize_width_str +
    resize_height_str +
    do_auto_resize_str +
    resize_max_height_str +
    auto_crop_str +
    gray_str;
if (addStr!=="") {
    addStr = "?" + addStr.substring(1);
    console.log("addStr:", addStr);
}

// Start Alpine.
Alpine.start()

// Document ready function to ensure the DOM is fully loaded.
document.addEventListener('DOMContentLoaded', function () {
    initFlowbite() // initialize Flowbite
})

// Add event listeners for all HTMX events.
document.body.addEventListener(
    'htmx:afterSwap htmx:afterRequest htmx:afterSettle',
    function () {
        initFlowbite() // initialize Flowbite
    }
)