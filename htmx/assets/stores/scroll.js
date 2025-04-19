// Scroll 卷轴模式
Alpine.store('scroll', {
    nowPageNum: 1,
    simplifyFilename: Alpine.$persist(true).as('scroll.simplifyFilename'), //是否简化标题
    //下拉模式下，漫画页面的底部间距。单位px。
    marginBottomOnScrollMode: Alpine.$persist(0).as(
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
}) 