// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package common

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/yumenaka/comigo/htmx/templates/common/svg"
	"github.com/yumenaka/comigo/model"
)

type HeaderProps struct {
	Title             string
	ShowReturnIcon    bool
	ReturnUrl         string
	SetDownLoadLink   bool
	InShelf           bool
	DownLoadLink      string
	SetTheme          bool
	FlipMode          bool
	ShowQuickJumpBar  bool
	QuickJumpBarBooks *model.BookInfoList
}

func Header(prop HeaderProps) templ.Component {
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
		var templ_7745c5c3_Var2 = []any{"toolbar flex flex-row justify-between w-full h-12 p-1 border-b bg-base-100 text-base-content border-slate-400", templ.KV("mx-auto opacity-90", prop.FlipMode)}
		templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var2...)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<header id=\"header\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if prop.FlipMode {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, " x-data=\"{ FlipMode: true }\" onmouseover=\"showToolbar();\" onmouseout=\"hideToolbar()\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if !prop.FlipMode {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, " x-data=\"{ FlipMode: false }\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, " class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(templ.CSSClasses(templ_7745c5c3_Var2).String())
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/common/header.templ`, Line: 1, Col: 0}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if prop.FlipMode {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, " :class=\"{ &#39;fixed absolute top-0&#39;: $store.flip.autoHideToolbar}\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "><div class=\"flex flex-row\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if prop.ShowReturnIcon {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "<!-- 返回箭头,点击返回上一页 --> <a href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 templ.SafeURL = templ.SafeURL(prop.ReturnUrl)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var4)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "\"><div class=\"flex justify-center items-center w-8 h-10 mx-1 my-0 rounded hover:ring\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = svg.Return().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, "</div></a>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, "<!-- 新服务器设置 --><a href=\"/settings\" class=\"flex justify-center items-center w-8 h-10 mx-1 my-0 rounded hover:ring\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = svg.ServerDisk().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, "</a> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if prop.InShelf {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 13, "<!-- 上传按钮，点击进入上传页面 --> <a href=\"/upload\" class=\"w-8 pt-1 mx-1 my-0 rounded hover:ring \">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = svg.Upload().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 14, "</a>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 15, "<!-- 文件重排序 Dropdown Menu： https://flowbite.com/docs/components/dropdowns/ --><svg class=\"flex justify-center items-center w-8 h-10 mx-1 my-0 rounded hover:ring\" id=\"dropdownHoverButton\" data-dropdown-toggle=\"ReSortDropdown\" data-dropdown-trigger=\"hover\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 24 24\"><path d=\"M10 18h4v-2h-4v2zM3 6v2h18V6H3zm3 7h12v-2H6v2z\" fill=\"currentColor\"></path></svg><div id=\"ReSortDropdown\" x-data=\"{\n                    addQuery(key, value) {\n                    const url = new URL(window.location.href);\n                    // 设置或更新查询参数\n                    url.searchParams.set(key, value);\n                    // 返回完整的 URL，包含新的查询参数\n                    return url.toString();\n                    }\n                }\" class=\"z-10 hidden divide-y divide-gray-100 rounded-lg shadow max-w-fit bg-white/90 dark:bg-gray-700\"><ul class=\"py-2 text-sm text-gray-700 dark:text-gray-200\" aria-labelledby=\"dropdownHoverButton\"><li><a :href=\"addQuery(&#39;sortBy&#39;,&#39;filename&#39;)\" x-text=\"i18next.t(&#39;sort_by_filename&#39;)\" class=\"block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white\"></a></li><li><a :href=\"addQuery(&#39;sortBy&#39;,&#39;modify_time&#39;)\" x-text=\"i18next.t(&#39;sort_by_modify_time&#39;)\" class=\"block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white\"></a></li><li><a :href=\"addQuery(&#39;sortBy&#39;,&#39;filesize&#39;)\" x-text=\"i18next.t(&#39;sort_by_filesize&#39;)\" class=\"block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white\"></a></li><li><a :href=\"addQuery(&#39;sortBy&#39;,&#39;filename_reverse&#39;)\" x-text=\"i18next.t(&#39;sort_by_filename_reverse&#39;)\" class=\"block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white\">Sign out</a></li><li><a :href=\"addQuery(&#39;sortBy&#39;,&#39;modify_time_reverse&#39;)\" x-text=\"i18next.t(&#39;sort_by_modify_time_reverse&#39;)\" class=\"block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white\">Sign out</a></li><li><a :href=\"addQuery(&#39;sortBy&#39;,&#39;filesize_reverse&#39;)\" x-text=\"i18next.t(&#39;sort_by_filesize_reverse&#39;)\" class=\"block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white\">Sign out</a></li></ul></div></div><!-- 标题--><div class=\"flex flex-col items-center justify-center flex-1 p-0 m-0 font-semibold text-center truncate\"><!-- 标题，快速跳转 or 可下载压缩包 or 只显示 -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if prop.InShelf {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 16, "<span class=\"text-lg font-semibold\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(prop.Title)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/common/header.templ`, Line: 99, Col: 52}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 17, "</span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			if prop.ShowQuickJumpBar && prop.QuickJumpBarBooks != nil {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 18, "<button id=\"quickJumpBarDropdownButton\" data-dropdown-toggle=\"QuickJumpDropdown\" data-dropdown-trigger=\"hover\" class=\"max-w-fit px-5 py-2.5 inline-flex  hover:bg-gray-400 rounded-lg font-semibold text-center items-center\" type=\"button\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var6 string
				templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(prop.Title)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/common/header.templ`, Line: 109, Col: 18}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = svg.ArrowDown().Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 19, "</button><!-- Dropdown menu --> <div id=\"QuickJumpDropdown\" class=\"max-w-full md:max-w-2xl md:truncate z-10 hidden divide-y divide-gray-100 rounded-lg shadow bg-white/90  dark:bg-gray-700\"><ul class=\"py-2 mt-0 text-sm text-gray-700 dark:text-gray-200\" aria-labelledby=\"quickJumpBarDropdownButton\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				for _, book := range prop.QuickJumpBarBooks.BookInfos {
					if prop.FlipMode {
						templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 20, "<li><a href=\"")
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
						var templ_7745c5c3_Var7 templ.SafeURL = templ.SafeURL("/flip/" + book.BookID)
						_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var7)))
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
						templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 21, "\" class=\"block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white\">")
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
						var templ_7745c5c3_Var8 string
						templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(book.Title)
						if templ_7745c5c3_Err != nil {
							return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/common/header.templ`, Line: 118, Col: 159}
						}
						_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
						templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 22, "</a></li>")
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
					} else {
						templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 23, "<li><a href=\"")
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
						var templ_7745c5c3_Var9 templ.SafeURL = templ.SafeURL("/scroll/" + book.BookID)
						_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var9)))
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
						templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 24, "\" class=\"block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white\">")
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
						var templ_7745c5c3_Var10 string
						templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(book.Title)
						if templ_7745c5c3_Err != nil {
							return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/common/header.templ`, Line: 122, Col: 161}
						}
						_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
						templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 25, "</a></li>")
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
					}
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 26, "</ul></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				if prop.SetDownLoadLink {
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 27, "<span class=\"text-lg font-semibold text-blue-700/90 hover:underline\"><a href=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var11 templ.SafeURL = templ.URL(prop.DownLoadLink)
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var11)))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 28, "\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var12 string
					templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(prop.Title)
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/common/header.templ`, Line: 131, Col: 60}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 29, "</a></span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else {
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 30, "<span class=\"text-lg font-semibold\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var13 string
					templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(prop.Title)
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/common/header.templ`, Line: 134, Col: 54}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 31, "</span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 32, "</div><!-- 溢出 overflow-x-auto :https://www.tailwindcss.cn/docs/overflow --><div class=\"flex justify-between p-0 m-0 max-w-64\"><!-- 图标：点击显示二维码 --><div data-modal-target=\"qrcode-modal\" data-modal-toggle=\"qrcode-modal\" class=\"flex justify-center items-center w-8 h-10 mx-1 my-0 rounded hover:ring\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = svg.QRCode().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 33, "</div><!-- 全屏按钮 --><div id=\"FullScreenIcon\" class=\"flex justify-center items-center w-8 h-10 mx-1 my-0 rounded hover:ring\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = svg.FullScreen().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 34, "</div><!-- 阅读器设定,点击屏幕中央也可以打开  可自定义方向 --><!-- data-drawer-body-scrolling=\"true\"  允许鼠标穿透，滚动下面的页面，设置此项有个好处，就是打开抽屉时背景不抖动 --><!-- https://flowbite.com/docs/components/drawer/#body-scrolling --><div class=\"flex justify-center items-center w-8 h-10 mx-1 my-0 rounded hover:ring\" id=\"OpenSettingButton\" data-drawer-target=\"drawer-right\" data-drawer-show=\"drawer-right\" aria-controls=\"drawer-right\" data-drawer-placement=\"right\" data-drawer-body-scrolling=\"true\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = svg.Setting().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 35, "</div></div></header>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
