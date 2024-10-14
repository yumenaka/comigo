// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package common

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/gin-gonic/gin"

type HeaderProps struct {
	Title           string
	ShowReturnIcon  bool
	ReturnUrl       string
	SetDownLoadLink bool
	InShelf         bool
	DownLoadLink    string
	SetTheme        bool
}

// 网页主题颜色的使用可以帮助你创建一个一致且视觉上吸引人的界面。以下是这些颜色的一般使用指南：

// 1. **primary (主要颜色)**：
//    - 主要用于品牌标识和主要的互动元素，例如按钮和链接。
//    - 在整个网站中保持一致，以帮助用户识别和导航。

// 2. **secondary (次要颜色)**：
//    - 辅助主要颜色，用于次要按钮、链接或背景。
//    - 可以在一些强调的地方使用，但不应喧宾夺主。

// 3. **accent (强调颜色)**：
//    - 用于引起注意的重要信息、通知或警示。
//    - 可以用于图标、标签或其他需要视觉突出的元素。

// 4. **neutral (中性色)**：
//    - 用于背景、文本或不需要吸引注意力的元素。
//    - 包括各种灰色调，可以帮助创建层次感和对比度。

// 5. **base-100**：
//    - 通常指最浅的背景颜色，通常是白色或接近白色的颜色。
//    - 用于主要的背景颜色，以确保文本和其他内容的可读性。

// ### 实际应用示例

// - **按钮**：
//   - 主要按钮：使用 primary 颜色
//   - 次要按钮：使用 secondary 颜色
//   - 危险操作按钮：使用 accent 颜色

// - **背景和文本**：
//   - 主背景：使用 base-100 颜色
//   - 次级背景：使用 neutral 颜色
//   - 主文本：通常使用黑色或深色
//   - 次级文本：使用浅灰色

// - **通知和警示**：
//   - 信息通知：使用 accent 颜色

//PrimaryColor：主题颜色。app的主要颜色，即整个屏幕和所有控件的主要颜色，首选颜色。

//SecondaryColor：提示性颜色。这颜色一般比PrimaryColor亮一些或暗一些，取决于白天模式还是黑暗模式。一般用于提示相关动作或信息，提示性颜色。

//AccentColor：交互性颜色。这颜色一般用于交互性的控件颜色，比如FloatingButton、TextField、Cursor、ProgressBar、Selection、Links等具体交互性的颜色

// https://github.com/L-Blondy/tw-colors
// https://github.com/RyanClementsHax/tailwindcss-themer

// 使用这些颜色，创建一个视觉和谐、易于导航的用户界面。如果有品牌指南或设计规范，可以依据这些规范进一步调整。
func Header(c *gin.Context, prop HeaderProps) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex flex-row justify-between w-full h-12 p-1 border-b bg-base-100 text-base-content border-slate-400\"><div class=\"flex flex-row\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if prop.ShowReturnIcon {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!-- 返回箭头,点击返回上一页 --> <a href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 templ.SafeURL = templ.SafeURL(prop.ReturnUrl)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var2)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><div class=\"w-8 pt-1 m-0\"><svg class=\"static w-8 rounded hover:ring\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 512 512\"><path fill=\"none\" stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"32\" d=\"M112 160l-64 64l64 64\"></path><path d=\"M64 224h294c58.76 0 106 49.33 106 108v20\" fill=\"none\" stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"32\"></path></svg></div></a>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!-- 服务器设置 --><div onclick=\"window.location.href=&#39;/admin/&#39;\" class=\"w-8 pt-1 m-0 rounded hover:ring\"><svg class=\"static w-8\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M5.25 14.25h13.5m-13.5 0a3 3 0 01-3-3m3 3a3 3 0 100 6h13.5a3 3 0 100-6m-16.5-3a3 3 0 013-3h13.5a3 3 0 013 3m-19.5 0a4.5 4.5 0 01.9-2.7L5.737 5.1a3.375 3.375 0 012.7-1.35h7.126c1.062 0 2.062.5 2.7 1.35l2.587 3.45a4.5 4.5 0 01.9 2.7m0 0a3 3 0 01-3 3m0 3h.008v.008h-.008v-.008zm0-6h.008v.008h-.008v-.008zm-3 6h.008v.008h-.008v-.008zm0-6h.008v.008h-.008v-.008z\"></path></svg></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if false {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!-- 上传按钮，点击进入上传页面 --> <div class=\"w-8 pt-1 mx-1 my-0 rounded hover:ring \"><svg class=\"static w-8\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 512 512\"><path d=\"M320 367.79h76c55 0 100-29.21 100-83.6s-53-81.47-96-83.6c-8.89-85.06-71-136.8-144-136.8c-69 0-113.44 45.79-128 91.2c-60 5.7-112 43.88-112 106.4s54 106.4 120 106.4h56\" fill=\"none\" stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"32\"></path> <path fill=\"none\" stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"32\" d=\"M320 255.79l-64-64l-64 64\"></path> <path fill=\"none\" stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"32\" d=\"M256 448.21V207.79\"></path></svg></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!-- 文件重排序 Dropdown Menu --><!-- 来自： https://flowbite.com/docs/components/dropdowns/ --><svg id=\"dropdownHoverButton\" data-dropdown-toggle=\"dropdownHover\" data-dropdown-trigger=\"hover\" class=\"w-10 pt-1 mx-0 rounded hover:ring\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 24 24\"><path d=\"M10 18h4v-2h-4v2zM3 6v2h18V6H3zm3 7h12v-2H6v2z\" fill=\"currentColor\"></path></svg><!-- Dropdown menu --><div id=\"dropdownHover\" class=\"z-10 hidden divide-y divide-gray-100 rounded-lg shadow w-44 bg-white/90 dark:bg-gray-700/90\"><ul class=\"py-2 text-sm text-gray-700 dark:text-gray-200\" aria-labelledby=\"dropdownHoverButton\"><li><a href=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 templ.SafeURL = templ.URL(AddQuery(c, "sortBy", "filename"))
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var3)))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" x-text=\"i18next.t(&#39;sort_by_filename&#39;)\" class=\"block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white\"></a></li><li><a href=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 templ.SafeURL = templ.URL(AddQuery(c, "sortBy", "modify_time"))
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var4)))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" x-text=\"i18next.t(&#39;sort_by_modify_time&#39;)\" class=\"block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white\"></a></li><li><a href=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 templ.SafeURL = templ.URL(AddQuery(c, "sortBy", "filesize"))
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var5)))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" x-text=\"i18next.t(&#39;sort_by_filesize&#39;)\" class=\"block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white\"></a></li><li><a href=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 templ.SafeURL = templ.URL(AddQuery(c, "sortBy", "filename_reverse"))
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var6)))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" x-text=\"i18next.t(&#39;sort_by_filename_reverse&#39;)\" class=\"block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white\">Sign out</a></li><li><a href=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var7 templ.SafeURL = templ.URL(AddQuery(c, "sortBy", "modify_time_reverse"))
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var7)))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" x-text=\"i18next.t(&#39;sort_by_modify_time_reverse&#39;)\" class=\"block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white\">Sign out</a></li><li><a href=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var8 templ.SafeURL = templ.URL(AddQuery(c, "sortBy", "filesize_reverse"))
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var8)))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" x-text=\"i18next.t(&#39;sort_by_filesize_reverse&#39;)\" class=\"block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white\">Sign out</a></li></ul></div></div><!-- 标题--><div class=\"flex flex-col justify-center flex-1 p-0 m-0 font-semibold text-center truncate\"><!-- 标题，快速跳转 or 可下载压缩包 or 只显示 -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if prop.InShelf {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"text-lg font-semibold\">QuickJumpBar</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			if prop.SetDownLoadLink {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"text-lg font-semibold text-blue-700 text-opacity-100 hover:underline\"><a href=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var9 templ.SafeURL = templ.URL(prop.DownLoadLink)
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var9)))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var10 string
				templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(prop.Title)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/common/header.templ`, Line: 159, Col: 59}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a></span>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"text-lg font-semibold\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var11 string
				templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(prop.Title)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/common/header.templ`, Line: 162, Col: 53}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><!-- 溢出 overflow-x-auto :https://www.tailwindcss.cn/docs/overflow --><div class=\"flex justify-between p-0 m-0 max-w-64\"><!-- 点击显示二维码 --><!-- 已知问题： 打开的时候body会加上一个 overflow-hidden 属性。此时滑动条被禁用，导致背景元素移动 --><!-- 临时解决方法：router/static/scripts.js，把这个句删掉 document.body.classList.add(\"overflow-hidden\"); --><!-- 不管用的解决方法：添加点击事件，把这个属性删掉  @click=\"document.body.classList.remove('overflow-hidden');\"  --><!-- 根本的解决方案（patch）(TODO:研究 patch 用法)：. https://bun.sh/docs/install/patch  --><!-- https://github.com/themesberg/flowbite/blob/c22565d406246749a09c5d556c540c102e0f98ae/src/components/modal/index.ts#L305 --><div data-modal-target=\"qrcode-modal\" data-modal-toggle=\"qrcode-modal\" class=\"w-8 pt-1 mx-1 my-0  rounded hover:ring \"><svg class=\"static w-8\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 512 512\"><rect x=\"336\" y=\"336\" width=\"80\" height=\"80\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"272\" y=\"272\" width=\"64\" height=\"64\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"416\" y=\"416\" width=\"64\" height=\"64\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"432\" y=\"272\" width=\"48\" height=\"48\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"272\" y=\"432\" width=\"48\" height=\"48\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"336\" y=\"96\" width=\"80\" height=\"80\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"288\" y=\"48\" width=\"176\" height=\"176\" rx=\"16\" ry=\"16\" fill=\"none\" stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"32\"></rect><rect x=\"96\" y=\"96\" width=\"80\" height=\"80\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"48\" y=\"48\" width=\"176\" height=\"176\" rx=\"16\" ry=\"16\" fill=\"none\" stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"32\"></rect><rect x=\"96\" y=\"336\" width=\"80\" height=\"80\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"48\" y=\"288\" width=\"176\" height=\"176\" rx=\"16\" ry=\"16\" fill=\"none\" stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"32\"></rect></svg></div><!-- 全屏按钮 --><div id=\"FullScreenIcon\" class=\"w-8 pt-1 mx-1 my-0 rounded hover:ring\"><svg class=\"static w-8\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 24 24\"><g fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><path d=\"M16 4h4v4\"></path> <path d=\"M14 10l6-6\"></path> <path d=\"M8 20H4v-4\"></path> <path d=\"M4 20l6-6\"></path> <path d=\"M16 20h4v-4\"></path> <path d=\"M14 14l6 6\"></path> <path d=\"M8 4H4v4\"></path> <path d=\"M4 4l6 6\"></path></g></svg></div><!-- 阅读器设定,点击屏幕中央也可以打开  可自定义方向 --><!-- data-drawer-body-scrolling=\"true\"  允许鼠标穿透，滚动下面的页面，设置此项有个好处，就是打开抽屉时背景不抖动 --><!-- https://flowbite.com/docs/components/drawer/#body-scrolling --><div id=\"OpenSettingButton\" data-drawer-target=\"drawer-right\" data-drawer-show=\"drawer-right\" aria-controls=\"drawer-right\" data-drawer-placement=\"right\" data-drawer-body-scrolling=\"true\" class=\"w-8 pt-1 mx-1 my-0 drawer-button  rounded hover:ring\"><svg class=\"static w-8\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 512 512\"><path d=\"M262.29 192.31a64 64 0 1 0 57.4 57.4a64.13 64.13 0 0 0-57.4-57.4zM416.39 256a154.34 154.34 0 0 1-1.53 20.79l45.21 35.46a10.81 10.81 0 0 1 2.45 13.75l-42.77 74a10.81 10.81 0 0 1-13.14 4.59l-44.9-18.08a16.11 16.11 0 0 0-15.17 1.75A164.48 164.48 0 0 1 325 400.8a15.94 15.94 0 0 0-8.82 12.14l-6.73 47.89a11.08 11.08 0 0 1-10.68 9.17h-85.54a11.11 11.11 0 0 1-10.69-8.87l-6.72-47.82a16.07 16.07 0 0 0-9-12.22a155.3 155.3 0 0 1-21.46-12.57a16 16 0 0 0-15.11-1.71l-44.89 18.07a10.81 10.81 0 0 1-13.14-4.58l-42.77-74a10.8 10.8 0 0 1 2.45-13.75l38.21-30a16.05 16.05 0 0 0 6-14.08c-.36-4.17-.58-8.33-.58-12.5s.21-8.27.58-12.35a16 16 0 0 0-6.07-13.94l-38.19-30A10.81 10.81 0 0 1 49.48 186l42.77-74a10.81 10.81 0 0 1 13.14-4.59l44.9 18.08a16.11 16.11 0 0 0 15.17-1.75A164.48 164.48 0 0 1 187 111.2a15.94 15.94 0 0 0 8.82-12.14l6.73-47.89A11.08 11.08 0 0 1 213.23 42h85.54a11.11 11.11 0 0 1 10.69 8.87l6.72 47.82a16.07 16.07 0 0 0 9 12.22a155.3 155.3 0 0 1 21.46 12.57a16 16 0 0 0 15.11 1.71l44.89-18.07a10.81 10.81 0 0 1 13.14 4.58l42.77 74a10.8 10.8 0 0 1-2.45 13.75l-38.21 30a16.05 16.05 0 0 0-6.05 14.08c.33 4.14.55 8.3.55 12.47z\" fill=\"none\" stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"32\"></path></svg></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
