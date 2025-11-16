// BookShelf 书架设置
Alpine.store('shelf', {
    bookCardMode: Alpine.$persist('gird').as('shelf.bookCardMode'), //gird,list,text
    showFilename: Alpine.$persist(true).as('shelf.showFilename'), //是否显示文件名
    showFileIcon: Alpine.$persist(true).as('shelf.showFileIcon'), //是否显示文件图标
    simplifyTitle: Alpine.$persist(true).as('shelf.simplifyTitle'), //是否简化标题
    InfiniteDropdown: Alpine.$persist(false).as('shelf.InfiniteDropdown'), //卷轴模式下，是否无限下拉
    bookCardShowTitleFlag: Alpine.$persist(true).as('shelf.bookCardShowTitleFlag'), // 书库中的书籍是否显示文字版标题
    readingProgressPercent: Alpine.$persist(false).as('shelf.readingProgressPercent'),// 书库中的书籍是否显示阅读进度百分比
    syncScrollFlag: false, // 同步滚动,目前还没做
    // 屏幕宽横比,inLandscapeMode的判断依据
    aspectRatio: 1.2,
    // 可见范围宽高的具体值
    clientWidth: 0,
    clientHeight: 0,
}) 