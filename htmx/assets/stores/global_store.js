// 用Alpine Persist 注册全局变量
// https://alpinejs.dev/plugins/persist#using-alpine-persist-global
Alpine.store('global', {
    // bgPattern 背景花纹
    bgPattern: Alpine.$persist('grid-line').as('global.bgPattern'),
    // 自动切边
    autoCrop: Alpine.$persist(false).as('global.autoCrop'),
    // 自动切边阈值,范围是0~100。多数情况下 1 就够了。
    autoCropNum: Alpine.$persist(1).as('global.autoCropNum'),
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
      if (this.readMode === 'flip') {
        this.readMode = 'scroll'
      } else {
        this.readMode = 'flip'
      }
    },
    // 竖屏模式
    isPortrait: false,
    // 横屏模式
    isLandscape: true,
    // 检测并设置视口方向
    checkOrientation() {
        const isPortrait = window.innerHeight > window.innerWidth;
        this.isPortrait = isPortrait;
        this.isLandscape = !isPortrait;
        console.log(`当前视口方向: ${isPortrait ? '竖屏' : '横屏'}`);
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