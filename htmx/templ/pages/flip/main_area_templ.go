// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package flip

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/yumenaka/comigo/htmx/embed"
	"github.com/yumenaka/comigo/htmx/state"
	"github.com/yumenaka/comigo/model"
)

func InsertData(bookData any, stateData any) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templ.JSONScript("NowBook", bookData).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.JSONScript("GlobalState", stateData).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func MainArea(s *state.GlobalState, book *model.Book) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div id=\"mouseMoveArea\" class=\"flex flex-col items-center justify-center flex-1 w-full max-w-full h-full min-h-full bg-base-100\" :class=\"(theme.toString() ===&#39;light&#39;||theme.toString() ===&#39;dark&#39;||theme.toString() ===&#39;retro&#39;||theme.toString() ===&#39;lofi&#39;||theme.toString() ===&#39;nord&#39;) ? ($store.global.bgPattern !== &#39;none&#39;?$store.global.bgPattern+&#39; bg-base-300&#39;:&#39;bg-base-300&#39;):($store.global.bgPattern !== &#39;none&#39;?$store.global.bgPattern:&#39;&#39;)\"><!-- 单页模式 --><div x-show=\"!$store.flip.doublePageMode\" class=\"manga_area scroll-hidden flex flex-row w-full h-full min-h-[100vh] m-0 p-0 justify-center item-center\"><div id=\"slider-container\" class=\"relative w-full h-full flex justify-center items-center overflow-hidden\"><div id=\"slider\" class=\"flex transition-transform duration-300 ease-out w-full h-full justify-center items-center\"><div id=\"prev-slide\" draggable=\"false\" class=\"slide w-full h-full flex justify-center items-center absolute\"><!-- 前一张图片 将在JS中动态加载 --></div><div id=\"current-slide\" class=\"slide w-full h-full flex justify-center items-center absolute\"><img id=\"SinglePageModeNowImage\" draggable=\"false\" class=\"object-contain m-0 max-w-full max-h-full h-full\"></div><div id=\"next-slide\" draggable=\"false\" class=\"slide w-full h-full flex justify-center items-center absolute\"><!-- 后一张图片 将在JS中动态加载 --></div></div></div></div><!-- 双页模式+日漫模式--><div x-show=\"$store.flip.doublePageMode &amp;&amp; !$store.flip.rightToLeft\" class=\"manga_area flex flex-row w-full h-full min-h-[100vh] m-0 p-0 justify-center-safe item-center\"><!-- 双页模式第二页 --><img id=\"DoublePageModeNextImageLTR\" x-show=\"$store.flip.nowPageNum &lt; $store.flip.allPageNum\" class=\"select-none object-contain m-0 max-w-1/2 w-auto max-h-screen grow-0\"><!-- 双页模式第一页 --><img id=\"DoublePageModeNowImageLTR\" class=\"select-none object-contain m-0 max-w-1/2 w-auto max-h-screen grow-0\"></div><!-- 双页模式+美漫模式--><div x-show=\"$store.flip.doublePageMode &amp;&amp; $store.flip.rightToLeft\" class=\"manga_area flex flex-row w-full h-full min-h-[100vh] m-0 p-0 justify-center-safe item-center\"><!-- 双页模式第一页 --><img id=\"DoublePageModeNowImageRTL\" class=\"select-none object-contain m-0 max-w-1/2 w-auto max-h-screen grow-0 \"><!-- 双页模式第二页 --><img id=\"DoublePageModeNextImageRTL\" x-show=\"$store.flip.nowPageNum &lt; $store.flip.allPageNum\" class=\"select-none object-contain m-0 max-w-1/2 w-auto max-h-screen grow-0\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = StepsRangeArea().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.Raw("<style>"+embed.GetFileStr("static/flip.css")+"</style>").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "<style>\n\t\t/* 滑动切换相关样式 */\n\t\t#slider-container {\n\t\t\ttouch-action: pan-y; /* 允许垂直滚动，但水平滑动会被捕获 */\n\t\t\tuser-select: none; /* 防止文本选择 */\n\t\t\tmin-height: 100vh; /* 至少占据视口高度的100% */\n\t\t}\n\t\t\n\t\t#slider {\n\t\t\twidth: 100%;\n\t\t\theight: 100%;\n\t\t\tmin-height: 100vh; /* 至少占据视口高度的100% */\n\t\t\twill-change: transform; /* 优化动画性能 */\n\t\t}\n\t\t\n\t\t.slide {\n\t\t\twidth: 100%;\n\t\t\theight: 100%;\n\t\t\tmin-height: 100vh; /* 至少占据视口高度的100% */\n\t\t\tflex-shrink: 0;\n\t\t\ttransform: translateX(0);\n\t\t}\n\t\t\n\t\t/* 基本位置设置，会通过JavaScript动态调整 */\n\t\t#prev-slide {\n\t\t\ttransform: translateX(-100%);\n\t\t}\n\t\t\n\t\t#current-slide {\n\t\t\ttransform: translateX(0);\n\t\t}\n\t\t\n\t\t#next-slide {\n\t\t\ttransform: translateX(100%);\n\t\t}\n\t\t\n\t\t/* 滑动翻页模式下的样式 */\n\t\t.swipe-enabled #slider-container {\n\t\t\tcursor: grab;\n\t\t}\n\t\t\n\t\t.swipe-enabled #slider-container:active {\n\t\t\tcursor: grabbing;\n\t\t}\n\t\t\n\t\t/* 确保manga_area有足够的高度 */\n\t\t.manga_area {\n\t\t\tmin-height: 100vh !important; /* 强制设置最小高度 */\n\t\t}\n\t</style>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
