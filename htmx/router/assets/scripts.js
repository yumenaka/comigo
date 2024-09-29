import 'htmx.org'
import Alpine from 'alpinejs'
import persist from '@alpinejs/persist'
import i18next from 'i18next'
import 'flowbite'
import LanguageDetector from 'i18next-browser-languagedetector'
import 'tw-colors'
// 将 Alpine 实例添加到窗口对象中。
window.Alpine = Alpine

// Alpine Persist 插件，用于持久化存储。默认存储到 localStorage。
// 详细用法参见： https://alpinejs.dev/plugins/persist
Alpine.plugin(persist)

// i18next 国际化插件，用于国际化。详细用法参见：
// https://www.i18next.com/overview/getting-started
import enLocale from './locales/en_US.json'
import zhLocale from './locales/zh_CN.json'
import jaLocale from './locales/ja_JP.json'
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

import screenfull from 'screenfull'
document.getElementById('FullScreenIcon').addEventListener('click', () => {
  if (screenfull.isEnabled) {
    screenfull.toggle()
  } else {
    // Ignore or do something else
    i18next.t('not_support_fullscreen')
  }
})

// 用Alpine Persist 注册全局变量
// https://alpinejs.dev/plugins/persist#using-alpine-persist-global

// global 全局设置
Alpine.store('global', {
  // debugMode 是否开启调试模式
  debugMode: Alpine.$persist(false).as('global.debugMode'),
  // readerMode 当前阅读模式
  readMode: Alpine.$persist('flip').as('global.readMode'),
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
  simplifyTitle: Alpine.$persist(true).as('shelf.simplifyTitle'), //是否简化标题
  InfiniteDropdown: true, //卷轴模式下，是否无限下拉
  bookCardShowTitleFlag: true, // 书库中的书籍是否显示文字版标题
  syncScrollFlag: false, // 同步滚动,目前还没做
  scrollTopSave: 0, //存储现在滚动的位置
  // 可见范围是否是横向
  isLandscapeMode: true,
  isPortraitMode: false,
  imageMaxWidth: 800,
  // 屏幕宽横比,inLandscapeMode的判断依据
  aspectRatio: 1.2,
  // 可见范围宽高的具体值
  clientWidth: 0,
  clientHeight: 0,
})

// Scroll 卷轴模式
Alpine.store('scroll', {
  nowPageNum: 0,
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
  widthUseFixedValue: true,
  //横屏(Landscape)状态的漫画页宽度,百分比
  singlePageWidth_Percent: 50,
  doublePageWidth_Percent: 95,
  //横屏(Landscape)状态的漫画页宽度。px。
  singlePageWidth_PX: 720,
  doublePageWidth_PX: 720,
  //可见范围是否是横向
  isLandscapeMode: true,
  isPortraitMode: false,
  //书籍数据,需要从远程拉取
  //是否显示顶部页头
  showHeaderFlag: true,
  //是否显示页数
  showPageNum: false,
  //ws翻页相关
  syncPageByWS: true, //是否通过websocket同步翻页
  // //此处修改不会实时生效，不要这么做
  // toggleSimplifyTitle() {
  //     this.simplifyTitle = ! this.simplifyTitle
  // }
})

// Flip 翻页模式
Alpine.store('Flip', {
  //自动隐藏工具条
  interval: 0,
  hideToolbar: true,
  //是否显示页头
  showHeaderFlag_FlipMode: true,
  //是否显示页脚
  showFooterFlag_FlipMode: true,
  //是否是右半屏翻页（从右到左）?日本漫画从左到右(false)
  rightToLeftFlag: Alpine.$persist(false).as('Flip.rightToLeftFlag'),
  //简单拼合双页
  doublePageModeFlag: Alpine.$persist(false).as('Flip.doublePageModeFlag'),
  //自动拼合双页,效果不太好
  autoDoublePageModeFlag: Alpine.$persist(false).as(
    'Flip.autoDoublePageModeFlag'
  ),
  //是否保存当前页数
  saveNowPageNumFlag: Alpine.$persist(true).as('Flip.saveNowPageNumFlag'),
  //素描模式标记
  sketchModeFlag: false,
  //是否显示素描提示
  showPageHintFlag_FlipMode: Alpine.$persist(false).as(
    'Flip.showPageHintFlag_FlipMode'
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

// Start Alpine.
Alpine.start()
