// Flip 翻页模式
Alpine.store('flip', {
    imageMaxWidth: 400,
    //自动隐藏工具条
    autoHideToolbar: Alpine.$persist(false).as('flip.autoHideToolbar'),
    // autoHideToolbar: Alpine.$persist((() => {
    //     const ua = navigator.userAgent || '';
    //     const isMobileUA = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini|Mobi/i.test(ua);
    //     const isTouch = (('ontouchstart' in window) || (navigator.maxTouchPoints > 1));
    //     //console.log('isMobileUA', isMobileUA);
    //     //console.log('isTouch', isTouch);
    //     // return isMobileUA || isTouch;
    //     return isMobileUA;
    // })()).as('flip.autoHideToolbar'),
    //自动对齐
    autoAlign: Alpine.$persist(true).as('flip.autoAlignTop'),
    //是否显示页头
    show_header: Alpine.$persist(true).as('flip.show_header'),
    //是否显示页脚
    showFooter: Alpine.$persist(true).as('flip.showFooter'),
    //是否显示页数
    showPageNum: Alpine.$persist(true).as('flip.showPageNum'),
    //是否是日本漫画【右半屏翻页,从左到右(true)】【右半屏翻页,从右到左(false)】
    mangaMode: Alpine.$persist(true).as('flip.mangaMode'),
    //触摸滑动翻页
    swipeTurn: Alpine.$persist(true).as('flip.swipeTurn'),
    //鼠标滚轮翻页
    wheelFlip: Alpine.$persist(false).as('flip.wheelFlip'),
    //双页模式
    doublePageMode: Alpine.$persist(false).as('flip.doublePageMode'),
    //自动拼合双页(TODO)
    autoDoublePageMode: Alpine.$persist(false).as(
        'flip.autoDoublePageModeFlag'
    ),
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
    
    // ============ 滑动动画配置（新增）============
    // 滑动动画持续时间（毫秒）
    swipeAnimationDuration: Alpine.$persist(300).as('flip.swipeAnimationDuration'),
    // 回弹动画持续时间（毫秒）
    resetAnimationDuration: Alpine.$persist(400).as('flip.resetAnimationDuration'),
    // 滑动阈值（像素）- 超过这个值才会触发翻页
    swipeThreshold: Alpine.$persist(100).as('flip.swipeThreshold'),
    // 快速滑动超时时间（毫秒）
    swipeTimeout: Alpine.$persist(300).as('flip.swipeTimeout'),
    
    // ============ 其他可配置参数（新增）============
    // 预加载图片范围
    preloadRange: Alpine.$persist(10).as('flip.preloadRange'),
    // 滚轮节流延迟（毫秒）
    wheelThrottleDelay: Alpine.$persist(250).as('flip.wheelThrottleDelay'),
    // WebSocket 最大重连次数
    websocketMaxReconnect: Alpine.$persist(200).as('flip.websocketMaxReconnect'),
    // WebSocket 重连间隔（毫秒）
    websocketReconnectInterval: Alpine.$persist(3000).as('flip.websocketReconnectInterval'),
}) 