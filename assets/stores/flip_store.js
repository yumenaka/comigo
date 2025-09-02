// Flip 翻页模式
Alpine.store('flip', {
    nowPageNum: 1,
    allPageNum: 100,
    imageMaxWidth: 400,
    //自动隐藏工具条
    autoHideToolbar: Alpine.$persist(true).as('flip.autoHideToolbar'),
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
    //swipeTurn or clickTurn
    swipeTurn: Alpine.$persist(true).as('flip.swipeTurn'),
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