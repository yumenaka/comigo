// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package scroll

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yumenaka/comigo/entity"
	"github.com/yumenaka/comigo/htmx/state"
)

// getImageUrlForAlpine 用于获取图片的URL
func getImageUrlForAlpine(Url string) string {
	return `{ imageUrl: '` + Url + `' + ($store.global.autoCrop  ? "&auto_crop=1" : '') }`
	//return `{ imageUrl: '` + Url + `' + ($store.global.autoCrop  ? "" : '') }`
}

// ScrollMainArea 定义 BodyHTML
// 需要更复杂的屏幕状态判断的时候，可以参考：https://developer.mozilla.org/zh-CN/docs/Web/API/Screen/orientation
// orientation: (screen.orientation || {}).type ||  screen.mozOrientation || screen.msOrientation
// tips：hx-get 用于获取图片的URL，hx-trigger 用于触发加载，hx-swap 用于替换元素，innerHTML默认值，将内容放在目标元素内 outerHTML用返回的内容替换整个目标元素  hx-target 用于指定目标元素
// https://htmx.org/docs/#triggers  https://htmx.org/docs/#swapping
// tips： Alpine.js 动态CSS，只支持内联写法
func ScrollMainArea(s *state.GlobalState, book *entity.Book) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"mouseMoveArea\" x-data=\"{\n            isDoublePage: false,\n            imageWidth: &#39;&#39;,\n            orientation: &#39;&#39;,\n            updateOrientation() {\n                this.orientation = (window.innerWidth / window.innerHeight &gt;= 1) ? &#39;landscape&#39; : &#39;portrait&#39;;\n            }\n        }\" x-init=\"updateOrientation();\" @resize.window=\"updateOrientation()\" class=\"flex flex-col items-center justify-center flex-1 w-full max-w-full pt-0 bg-base-100 text-base-content\" :class=\"(theme.toString() ===&#39;light&#39;||theme.toString() ===&#39;dark&#39;||theme.toString() ===&#39;retro&#39;||theme.toString() ===&#39;lofi&#39;||theme.toString() ===&#39;nord&#39;) ? ($store.global.bgPattern !== &#39;none&#39;?$store.global.bgPattern+&#39; bg-base-300&#39;:&#39;bg-base-300&#39;):($store.global.bgPattern !== &#39;none&#39;?$store.global.bgPattern:&#39;&#39;)\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for key, image := range book.Pages.Images {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex flex-col justify-start w-full max-w-full m-0 rounded item-center\" :style=\"{ marginBottom: $store.scroll.marginBottomOnScrollMode + &#39;px&#39; }\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if strings.Contains(image.Url, ".html") && !strings.Contains(image.Url, "hidden.") {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div hx-get=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var2 string
				templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(image.Url)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/scroll/scroll_main_area.templ`, Line: 47, Col: 28}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-trigger=\"load\" hx-swap=\"innerHTML\" class=\"w-full m-0\"></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if !strings.Contains(image.Url, "hidden.") && !strings.Contains(image.Url, ".html") {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<img x-data=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var3 string
				templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(getImageUrlForAlpine(image.Url))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/scroll/scroll_main_area.templ`, Line: 51, Col: 48}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"w-full manga_image\" :src=\"imageUrl\" @load=\"\n                            if ($event.target.naturalWidth &gt; $event.target.naturalHeight) {\n                                // 双页\n                                isDoublePage = true;\n                                $el.classList.add(&#39;double&#39;);\n                                $el.classList.remove(&#39;single&#39;);\n                            } else {\n                                // 单页\n                                isDoublePage = false;\n                                $el.classList.add(&#39;single&#39;);\n                                $el.classList.remove(&#39;double&#39;);\n                            }\" @resize.window=\"updateOrientation()\" :style=\"{ width: orientation.toString().includes(&#39;landscape&#39;) ?(Alpine.store(&#39;scroll&#39;).widthUseFixedValue? (isDoublePage ? Alpine.store(&#39;scroll&#39;).doublePageWidth_PX +&#39;px&#39;: Alpine.store(&#39;scroll&#39;).singlePageWidth_PX +&#39;px&#39;): (isDoublePage ? Alpine.store(&#39;scroll&#39;).doublePageWidth_Percent + &#39;%&#39;: Alpine.store(&#39;scroll&#39;).singlePageWidth_Percent + &#39;%&#39;)): &#39;100%&#39;, maxWidth: &#39;100%&#39;}\" alt=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var4 string
				templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(key))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/scroll/scroll_main_area.templ`, Line: 68, Col: 29}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<template x-if=\"$store.scroll.showPageNum\"><div class=\"w-full mt-0 mb-1 text-sm font-semibold text-center page_hint \">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d / %d", key+1, book.BookInfo.PageCount))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/scroll/scroll_main_area.templ`, Line: 72, Col: 135}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></template></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><button id=\"BackTopButton\" style=\"display: none\" class=\"fixed flex items-center justify-center w-10 h-10 text-white bg-blue-500 rounded-full shadow-lg bottom-4 right-4\"><svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 512 512\"><path d=\"M256 48C141.13 48 48 141.13 48 256s93.13 208 208 208s208-93.13 208-208S370.87 48 256 48zm96 270.63l-96-96l-96 96L137.37 296L256 177.37L374.63 296z\" fill=\"currentColor\"></path></svg></button><style>\n    /* https://developer.mozilla.org/zh-CN/docs/Web/CSS/object-fit */\n    .manga_image {\n        margin: auto;\n        box-shadow: 0px 6px 3px 0px rgba(0, 0, 0, 0.19);\n    }\n\n    .page_hint {\n        /* 文字颜色 */\n        color: #413d3d;\n        /* 文字阴影：https://www.w3school.com.cn/css/css3_shadows.asp*/\n        text-shadow: -1px 0 rgb(240, 229, 229), 0 1px rgb(253, 242, 242),\n            1px 0 rgb(206, 183, 183), 0 -1px rgb(196, 175, 175);\n    }\n\n    .LoadingImage {\n        width: 90vw;\n        max-width: 90vw;\n    }\n\n    .ErrorImage {\n        width: 90vw;\n        max-width: 90vw;\n    }\n\n    /* 横屏（显示区域）时的CSS样式,IE无效 */\n    @media screen and (min-aspect-ratio: 19/19) {\n        .SinglePageImage {\n            width: v-bind(sPWL);\n            max-width: 100%;\n        }\n\n        .DoublePageImage {\n            width: v-bind(dPWL);\n            max-width: 100%;\n        }\n    }\n\n    /* 竖屏(显示区域)CSS样式,IE无效 */\n    @media screen and (max-aspect-ratio: 19/19) {\n        .SinglePageImage {\n            width: v-bind(sPWP);\n            max-width: 100%;\n        }\n\n        .DoublePageImage {\n            /* width: 100%; */\n            width: v-bind(dPWP);\n            max-width: 100%;\n        }\n    }\n    </style>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate