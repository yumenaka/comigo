// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package flip

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/yumenaka/comigo/htmx/embed_files"
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
		return templ_7745c5c3_Err
	})
}

func FlipMainArea(s *state.GlobalState, book *model.Book) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"mouseMoveArea\" class=\"grid-line flex flex-col items-center justify-center flex-1 w-full max-w-full bg-base-100 text-base-content\" :class=\"(theme.toString() ===&#39;light&#39;||theme.toString() ===&#39;dark&#39;||theme.toString() ===&#39;retro&#39;||theme.toString() ===&#39;lofi&#39;||theme.toString() ===&#39;nord&#39;) ? ($store.global.bgPattern !== &#39;none&#39;?$store.global.bgPattern+&#39; bg-base-300&#39;:&#39;bg-base-300&#39;):($store.global.bgPattern !== &#39;none&#39;?$store.global.bgPattern:&#39;&#39;)\"><div class=\"manga_area\" id=\"MangaMain\"><div class=\"manga_area_img_div\"><!-- 非自动拼合模式最简单,直接显示一张图 --><img id=\"NowImage\" class=\"max-w-full max-h-screen m-0 rounded shadow-lg\" src=\"/static/images/ball-triangle.svg\" alt=\"\"><!-- 简单拼合双页,不管单双页什么的 --><img id=\"NextImage\" x-show=\"$store.flip.doublePageMode &amp;&amp; $store.flipnowPageNum &lt; $store.flip.allPageNum\" src=\"/static/images/ball-triangle.svg\"></div></div></div><!-- 底部的阅读进度条 --><!-- https://flowbite.com/docs/forms/range/ --><!-- 宽度：w-5/6 https://www.tailwindcss.cn/docs/width 使用 w-{fraction} 或 w-full 将元素设置为基于百分比的宽度。 --><!-- 定位：https://www.tailwindcss.cn/docs/position  --><!-- 使用 fixed 来定位一个元素相对于浏览器窗视口的位置。偏移量是相对于视口计算的，且该元素将作为绝对定位的子元素的位置参考。 --><!-- 控制 flex 和 grid 项目如何沿着容器的主轴定位:https://www.tailwindcss.cn/docs/justify-content --><!-- Tailwind 的容器不会自动居中，也没有任何内置的水平方向的内边距。要使一个容器居中，使用 mx-auto 功能类： --><div id=\"steps-range_area\" onmouseover=\"showToolbar();\" onmouseout=\"hideToolbar()\" class=\"w-full px-2 overflow-hidden bg-gray-400 border border-blue-800 rounded toolbar h-14 opacity-80\" :class=\"Alpine.store(&#39;flip&#39;).autoHideToolbar? &#39;absolute fixed bottom-0&#39;:&#39;flex flex-col justify-center&#39;\"><label for=\"steps-range\" class=\"block m-0 text-sm font-medium text-center text-gray-900 dark:text-white\" x-text=\"$store.flip.nowPageNum+&#39;/&#39;+$store.flip.allPageNum\"></label> <input id=\"steps-range\" class=\"w-full h-2 mb-2 bg-yellow-800 rounded-lg appearance-none cursor-pointer dark:bg-gray-700\" type=\"range\" min=\"1\" :max=\"$store.flip.allPageNum\" x-model=\"$store.flip.nowPageNum\" onchange=\"setImageSrc()\" step=\"1\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.Raw("<style>"+embed_files.GetFileStr("static/flip.css")+"</style>").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
