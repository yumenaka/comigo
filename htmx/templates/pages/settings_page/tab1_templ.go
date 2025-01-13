// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.819
package settings_page

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/yumenaka/comigo/htmx/state"

// htmx 是一个用于构建现代 web 应用程序的库，它使用无需刷新页面的 AJAX 技术。请求是通过 HTTP 发送的，但是 htmx 会处理响应并更新页面的部分内容，而不是整个页面。
// https://htmx.org/docs/#parameters
// 默认情况下，引起请求的元素将包含其值（如果有）。如果元素是一个表单，它将包含其中所有输入的值。
// 与 HTML 表单一样，输入的 name 属性用作 htmx 发送的请求中的参数名称。
// 此外，如果该元素导致非 GET 请求，则将包含最近的封闭表单的所有输入的值。
// 此外，还可以使用 hx-vals（like: hx-vals='{"myVal": "My Value"}'）在请求中包含额外的值
func tab1(s *state.GlobalState) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div class=\"flex flex-col w-full p-2 m-1 font-semibold rounded-md shadow-md hover:shadow-2xl justify-left items-left\" style=\"background-color: rgb(245, 245, 228);\"><div class=\"w-32\">打开浏览器</div><label for=\"OpenBrowser\" class=\"relative h-8 cursor-pointer w-14\"><input type=\"checkbox\" id=\"OpenBrowser\" name=\"OpenBrowser\" class=\"sr-only peer\"> <span class=\"absolute inset-0 transition bg-gray-300 rounded-full peer-checked:bg-green-500\"></span> <span class=\"absolute inset-y-0 w-6 h-6 m-1 transition-all bg-white rounded-full start-0 peer-checked:start-6\"></span></label><div class=\"w-3/4 py-1 text-xs text-gray-500\">扫描完成后，是否同时打开浏览器。windows默认true，其他平台默认false。</div></div><div class=\"flex flex-col justify-start w-full p-2 m-1 font-semibold rounded-md shadow-md hover:shadow-2xl items-left\" style=\"background-color: rgb(245, 245, 228);\"><label class=\"w-32 py-0\" for=\"LocalStores\">书库文件夹</label><div class=\"flex flex-row flex-wrap w-3/4 py-1\"><div class=\"flex flex-row items-center p-2 m-1 text-sm font-medium text-black bg-blue-300 rounded-2xl\">/Users/bai/cvgo/htmx/temp/comi <svg class=\"w-5 h-5 mx-1\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\"><g fill=\"none\"><path d=\"M11.5 4a1.5 1.5 0 0 0-3 0h-1a2.5 2.5 0 0 1 5 0H17a.5.5 0 0 1 0 1h-.554L15.15 16.23A2 2 0 0 1 13.163 18H6.837a2 2 0 0 1-1.987-1.77L3.553 5H3a.5.5 0 0 1-.492-.41L2.5 4.5A.5.5 0 0 1 3 4h8.5zm3.938 1H4.561l1.282 11.115a1 1 0 0 0 .994.885h6.326a1 1 0 0 0 .993-.885L15.438 5zM8.5 7.5c.245 0 .45.155.492.359L9 7.938v6.125c0 .241-.224.437-.5.437c-.245 0-.45-.155-.492-.359L8 14.062V7.939c0-.242.224-.438.5-.438zm3 0c.245 0 .45.155.492.359l.008.079v6.125c0 .241-.224.437-.5.437c-.245 0-.45-.155-.492-.359L11 14.062V7.939c0-.242.224-.438.5-.438z\" fill=\"currentColor\"></path></g></svg></div><div class=\"flex flex-row items-center p-2 m-1 text-sm font-medium text-black bg-blue-300 rounded-2xl\">../test <svg class=\"w-5 h-5 mx-1\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\"><g fill=\"none\"><path d=\"M11.5 4a1.5 1.5 0 0 0-3 0h-1a2.5 2.5 0 0 1 5 0H17a.5.5 0 0 1 0 1h-.554L15.15 16.23A2 2 0 0 1 13.163 18H6.837a2 2 0 0 1-1.987-1.77L3.553 5H3a.5.5 0 0 1-.492-.41L2.5 4.5A.5.5 0 0 1 3 4h8.5zm3.938 1H4.561l1.282 11.115a1 1 0 0 0 .994.885h6.326a1 1 0 0 0 .993-.885L15.438 5zM8.5 7.5c.245 0 .45.155.492.359L9 7.938v6.125c0 .241-.224.437-.5.437c-.245 0-.45-.155-.492-.359L8 14.062V7.939c0-.242.224-.438.5-.438zm3 0c.245 0 .45.155.492.359l.008.079v6.125c0 .241-.224.437-.5.437c-.245 0-.45-.155-.492-.359L11 14.062V7.939c0-.242.224-.438.5-.438z\" fill=\"currentColor\"></path></g></svg></div><div class=\"flex flex-row items-center p-2 m-1 text-sm font-medium text-black bg-blue-300 rounded-2xl\"><svg class=\"w-5 h-5 mx-1\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\"><g fill=\"none\"><path d=\"M11.5 4a1.5 1.5 0 0 0-3 0h-1a2.5 2.5 0 0 1 5 0H17a.5.5 0 0 1 0 1h-.554L15.15 16.23A2 2 0 0 1 13.163 18H6.837a2 2 0 0 1-1.987-1.77L3.553 5H3a.5.5 0 0 1-.492-.41L2.5 4.5A.5.5 0 0 1 3 4h8.5zm3.938 1H4.561l1.282 11.115a1 1 0 0 0 .994.885h6.326a1 1 0 0 0 .993-.885L15.438 5zM8.5 7.5c.245 0 .45.155.492.359L9 7.938v6.125c0 .241-.224.437-.5.437c-.245 0-.45-.155-.492-.359L8 14.062V7.939c0-.242.224-.438.5-.438zm3 0c.245 0 .45.155.492.359l.008.079v6.125c0 .241-.224.437-.5.437c-.245 0-.45-.155-.492-.359L11 14.062V7.939c0-.242.224-.438.5-.438z\" fill=\"currentColor\"></path></g></svg></div><div class=\"relative\"><label for=\"Array\" class=\"sr-only\">键入或粘贴内容</label><input type=\"text\" id=\"Array\" placeholder=\"键入或粘贴内容\" class=\"w-full rounded-md border-gray-400 py-2.5 pe-10 shadow-sm sm:text-sm\"> <span class=\"absolute top-[0px] right-[-80px] place-content-center\"><button type=\"button\" class=\"w-16 h-10 mx-2 my-1 text-center text-gray-700 transition border border-gray-500 rounded bg-sky-300 hover:text-gray-900\">提交</button></span></div></div><div class=\"w-3/4 py-1 ml-2 text-xs text-gray-500\">书库文件夹，支持绝对目录与相对目录。相对目录以当前执行目录为基准</div><div class=\"bg-red-600\"></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
