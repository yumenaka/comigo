// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.819
package settings_page

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func ConfigManager(initSaveTo string, WorkingDirectoryConfig string, HomeDirectoryConfig string, ProgramDirectoryConfig string) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div id=\"config-container\" class=\"flex flex-col justify-start w-full p-2 m-1 font-semibold border rounded-md shadow-md hover:shadow-2xl items-left bg-base-100 text-base-content border-slate-400\" x-data=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs("{ selectedDir: '" + initSaveTo + "',initinitSaveTo:'" + initSaveTo + "'}")
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/settings_page/config_manager.templ`, Line: 7, Col: 85}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "\"><!-- 标题 --><label class=\"w-full py-0\" x-text=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(getTranslations("ConfigManager"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/settings_page/config_manager.templ`, Line: 10, Col: 70}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "\"></label><!-- 配置选项区域：点击一个卡片，就设置隐藏字段的值。 --><div class=\"flex flex-row w-full mx-0 my-1\"><div data-save_to=\"WorkingDirectory\" class=\"flex flex-col items-center justify-center w-1/3 pt-2 mx-1 text-xs font-normal border border-gray-500 rounded min-h-20 cursor-pointer\" :class=\"selectedDir === &#39;WorkingDirectory&#39;?&#39;bg-cyan-200&#39;:&#39;&#39;\" @click=\"selectedDir = &#39;WorkingDirectory&#39;\"><img class=\"h-7 w-7\" src=\"/static/images/working_directory.png\"><div class=\"mt-1\" x-text=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(getTranslations("WorkingDirectory"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/settings_page/config_manager.templ`, Line: 20, Col: 66}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "\">当前工作目录</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if WorkingDirectoryConfig != "" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "<div class=\"mx-1 my-1 text-xs text-gray-500 line-clamp-2 hover:line-clamp-none active:line-clamp-none\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(WorkingDirectoryConfig)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/settings_page/config_manager.templ`, Line: 24, Col: 30}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "</div><div data-save_to=\"HomeDirectory\" class=\"flex flex-col items-center justify-center w-1/3 pt-2 mx-1 text-xs font-normal border border-gray-500 rounded min-h-20 cursor-pointer\" :class=\"selectedDir === &#39;HomeDirectory&#39;?&#39;bg-cyan-200&#39;:&#39;&#39;\" @click=\"selectedDir = &#39;HomeDirectory&#39;\"><img class=\"h-7 w-7\" src=\"/static/images/home_directory.png\"><div class=\"mt-1\" x-text=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 string
		templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(getTranslations("HomeDirectory"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/settings_page/config_manager.templ`, Line: 34, Col: 63}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "\">用户主目录</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if HomeDirectoryConfig != "" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "<div class=\"mx-1 my-1 text-xs text-gray-500 line-clamp-2 hover:line-clamp-none active:line-clamp-none\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var7 string
			templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(HomeDirectoryConfig)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/settings_page/config_manager.templ`, Line: 38, Col: 27}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, "</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, "</div><div data-save_to=\"ProgramDirectory\" class=\"flex flex-col items-center justify-center w-1/3 pt-2 mx-1 text-xs font-normal border border-gray-500 rounded min-h-20 cursor-pointer\" :class=\"selectedDir === &#39;ProgramDirectory&#39;?&#39;bg-cyan-200&#39;:&#39;&#39;\" @click=\"selectedDir = &#39;ProgramDirectory&#39;\"><img class=\"h-7 w-7\" src=\"/static/images/program_directory.png\"><div class=\"mt-1\" x-text=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var8 string
		templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(getTranslations("ProgramDirectory"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/settings_page/config_manager.templ`, Line: 48, Col: 66}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, "\">程序所在目录</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if ProgramDirectoryConfig != "" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 13, "<div class=\"mx-1 my-1 text-xs text-gray-500 line-clamp-2 hover:line-clamp-none active:line-clamp-none\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var9 string
			templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(ProgramDirectoryConfig)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/settings_page/config_manager.templ`, Line: 52, Col: 30}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 14, "</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 15, "</div></div><!-- 用来记录用户当前选择的隐藏字段（和Alpine双向绑定） --><input id=\"selectedDir\" name=\"selectedDir\" type=\"hidden\" x-model=\"selectedDir\"><!-- SAVE 按钮：发起 HTMX 请求到 /api/config-save --><div class=\"flex flex-row justify-center w-full\"><!-- hx-include 将 #selectedDir 表单值一并提交 --><!-- hx-target 将返回结果替换到本容器内 --><button class=\"w-24 h-10 mx-2 my-1 text-center text-gray-700 transition border border-gray-500 rounded bg-sky-300 hover:text-gray-900\" hx-post=\"/api/config-save\" hx-include=\"#selectedDir\" hx-target=\"#config-container\" hx-swap=\"outerHTML\">SAVE</button> <button class=\"h-10 w-24 mx-2 my-1 bg-red-300 border border-gray-500 text-center text-gray-700 transition hover:text-gray-900 rounded\" hx-post=\"/api/config-delete\" hx-include=\"#selectedDir\" hx-target=\"#config-container\" hx-swap=\"outerHTML\">DELETE</button></div><!-- 说明文字 --><div class=\"w-full py-1 text-xs text-gray-500\" x-text=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var10 string
		templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(getTranslations("ConfigManagerDescription"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/settings_page/config_manager.templ`, Line: 87, Col: 101}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 16, "\">点击Save，会将当前配置上传到服务器，并覆盖已经存在的设定文件。</div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
