function isScrollReadPage() {
	const pathname = window.ComiGoRelativePath ? window.ComiGoRelativePath(window.location.pathname) : window.location.pathname;
	return pathname.split("/").filter(Boolean)[0] === "scroll";
}

function normalizeScrollPageLimit(value) {
	const pageLimit = parseInt(value, 10);
	return Number.isInteger(pageLimit) && pageLimit > 0 ? pageLimit : 32;
}

// Scroll 卷轴阅读
Alpine.store("scroll", {
	simplifyTitle: Alpine.$persist(true).as("scroll.simplifyTitle"), //是否简化标题
	//下拉模式下，漫画页面的底部间距。单位px。
	marginBottomOnScrollMode: Alpine.$persist(0).as(
		"scroll.marginBottomOnScrollMode",
	),
	// 卷轴阅读的加载策略：infinite=一次加载全部，lazy=延迟加载，paged=分页加载
	loadMode: Alpine.$persist("infinite").as("scroll.loadMode"),
	// 延迟加载与分页加载每批页面上限
	pageLimit: Alpine.$persist(32).as("scroll.pageLimit"),
	// 卷轴阅读的同步滚动,目前还没做
	syncScrollFlag: Alpine.$persist(false).as("scroll.syncScrollFlag"),
	imageMaxWidth: 400,
	// 屏幕宽横比,inLandscapeMode的判断依据
	aspectRatio: 1.2,
	// 可见范围宽高的具体值
	clientWidth: 0,
	clientHeight: 0,
	//漫画页的单位,是否使用固定值
	widthUseFixedValue: Alpine.$persist(true).as("scroll.widthUseFixedValue"),
    //竖屏(Portrait)状态的漫画页宽度,百分比
	portraitWidthPercent: Alpine.$persist(100).as("scroll.portraitWidthPercent"),
	//横屏(Landscape)状态的漫画页宽度,百分比
	singlePageWidth_Percent: Alpine.$persist(60).as(
		"scroll.singlePageWidth_Percent",
	),
	doublePageWidth_Percent: Alpine.$persist(60).as(
		"scroll.doublePageWidth_Percent",
	),
	//横屏(Landscape)状态的漫画页宽度。px。
	singlePageWidth_PX: Alpine.$persist(800).as("scroll.singlePageWidth_PX"),
	doublePageWidth_PX: Alpine.$persist(800).as("scroll.doublePageWidth_PX"),
	//书籍数据,需要从远程拉取
	//是否显示顶部页头
	showHeaderFlag: true,
	//是否显示页数
	showPageNum: Alpine.$persist(false).as("scroll.showPageNum"),
});

if (isScrollReadPage()) {
	const scrollURLParams = new URLSearchParams(window.location.search);
	const scrollStore = Alpine.store("scroll");
	if (scrollURLParams.has("page")) {
		scrollStore.loadMode = "paged";
		scrollStore.pageLimit = normalizeScrollPageLimit(scrollURLParams.get("limit"));
	} else if (scrollStore.loadMode === "paged" || !["infinite", "lazy"].includes(scrollStore.loadMode)) {
		scrollStore.loadMode = "infinite";
	}
	scrollStore.pageLimit = normalizeScrollPageLimit(scrollStore.pageLimit);
}
