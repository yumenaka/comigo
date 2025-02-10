// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package shelf

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func ShelfDrawerSlot() templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!-- 阅读模式 --><!-- toggle组件来自： https://flowbite.com/docs/forms/toggle/ --><label class=\"inline-flex items-center w-full my-4 cursor-pointer outline outline-offset-8 outline-dotted hover:outline outline-2\"><input type=\"checkbox\" :value=\"$store.global.readMode === &#39;scroll&#39;\" x-on:click=\"$store.global.toggleReadMode()\" class=\"sr-only peer\"><div class=\"relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[&#39;&#39;] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600\"></div><span class=\"text-sm font-medium ms-3\" x-text=\"$store.global.readMode === &#39;scroll&#39;?i18next.t(&#39;scroll_mode&#39;):i18next.t(&#39;flip_mode&#39;)\">Toggle me</span></label><!-- 显示书名  --><label class=\"inline-flex items-center w-full my-4 cursor-pointer outline outline-offset-8 outline-dotted hover:outline outline-2\"><input type=\"checkbox\" :value=\"$store.shelf.showTitle\" x-on:click=\"$store.shelf.showTitle =!$store.shelf.showTitle\" class=\"sr-only peer\"><div class=\"relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[&#39;&#39;] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600\"></div><span class=\"text-sm font-medium ms-3\" x-text=\"i18next.t(&#39;show_book_titles&#39;)\"></span></label><!-- 简化书名 --><label x-show=\"$store.shelf.showTitle\" class=\"inline-flex items-center w-full my-4 cursor-pointer outline outline-offset-8 outline-dotted hover:outline outline-2\"><input type=\"checkbox\" :value=\"$store.shelf.simplifyTitle\" x-on:click=\"$store.shelf.simplifyTitle =!$store.shelf.simplifyTitle\" class=\"sr-only peer\"><div class=\"relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[&#39;&#39;] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600\"></div><span class=\"text-sm font-medium ms-3\" x-text=\"i18next.t(&#39;simplify_book_titles&#39;)\"></span></label><!-- 显示文件类型  --><label class=\"inline-flex items-center w-full my-4 cursor-pointer outline outline-offset-8 outline-dotted hover:outline outline-2\"><input type=\"checkbox\" :value=\"$store.shelf.showFileType\" x-on:click=\"$store.shelf.showFileType =!$store.shelf.showFileType\" class=\"sr-only peer\"><div class=\"relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[&#39;&#39;] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600\"></div><span class=\"text-sm font-medium ms-3\" x-text=\"i18next.t(&#39;show_file_type&#39;)\"></span></label><!-- debug mode --><label class=\"inline-flex items-center w-full my-4 cursor-pointer outline outline-offset-8 outline-dotted hover:outline outline-2\"><input type=\"checkbox\" :value=\"$store.global.debugMode\" x-on:click=\"$store.global.debugMode =!$store.global.debugMode\" class=\"sr-only peer\"><div class=\"relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[&#39;&#39;] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600\"></div><span class=\"text-sm font-medium ms-3\" x-text=\"i18next.t(&#39;debugMode&#39;)\"></span></label>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
