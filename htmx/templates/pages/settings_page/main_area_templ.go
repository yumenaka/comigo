// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.819
package settings_page

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/yumenaka/comigo/htmx/state"

func MainArea(s *state.GlobalState) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<header id=\"header\" class=\"toolbar flex justify-between w-full h-12 p-1 border-b bg-base-100 text-base-content border-slate-400\"><a href=\"/\" class=\"text-3xl w-10 pt-1 mx-2 rounded hover:ring\">🔙</a><!-- 老设置界面，开发时参考--><a href=\"/admin\" class=\"text-3xl w-10 pt-1 mx-2 rounded hover:ring\">⚙️</a><!-- 标题--><div class=\"flex flex-1  w-92 items-center justify-center p-0 m-0 font-semibold text-center truncate\"><button class=\"flex items-center justify-center min-w-24 mx-1 my-2 p-2 h-10 border border-solid rounded border-gray-400 bg-white/90\">📖&nbsp;Book</button> <button class=\"flex items-center justify-center min-w-24 mx-1 my-2 p-2 h-10 border border-solid rounded border-gray-400 bg-white/90\">🛜&nbsp;Net</button> <button class=\"flex items-center justify-center min-w-24 mx-1 my-2 p-2 h-10 border border-solid rounded border-gray-400 bg-white/90\">🧪&nbsp;Lab</button></div><!-- 溢出 overflow-x-auto :https://www.tailwindcss.cn/docs/overflow --><div class=\"flex justify-between p-0 m-0 max-w-64\"><div data-modal-target=\"qrcode-modal\" data-modal-toggle=\"qrcode-modal\" class=\"w-8 pt-1 mx-1 my-0 rounded hover:ring \"><svg class=\"static w-8\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 512 512\"><rect x=\"336\" y=\"336\" width=\"80\" height=\"80\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"272\" y=\"272\" width=\"64\" height=\"64\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"416\" y=\"416\" width=\"64\" height=\"64\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"432\" y=\"272\" width=\"48\" height=\"48\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"272\" y=\"432\" width=\"48\" height=\"48\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"336\" y=\"96\" width=\"80\" height=\"80\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"288\" y=\"48\" width=\"176\" height=\"176\" rx=\"16\" ry=\"16\" fill=\"none\" stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"32\"></rect><rect x=\"96\" y=\"96\" width=\"80\" height=\"80\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"48\" y=\"48\" width=\"176\" height=\"176\" rx=\"16\" ry=\"16\" fill=\"none\" stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"32\"></rect><rect x=\"96\" y=\"336\" width=\"80\" height=\"80\" rx=\"8\" ry=\"8\" fill=\"currentColor\"></rect><rect x=\"48\" y=\"288\" width=\"176\" height=\"176\" rx=\"16\" ry=\"16\" fill=\"none\" stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"32\"></rect></svg></div><div id=\"FullScreenIcon\" class=\"w-8 pt-1 mx-1 my-0 rounded hover:ring\"><svg class=\"static w-8\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><g fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><path d=\"M16 4h4v4\"></path> <path d=\"M14 10l6-6\"></path> <path d=\"M8 20H4v-4\"></path> <path d=\"M4 20l6-6\"></path> <path d=\"M16 20h4v-4\"></path> <path d=\"M14 14l6 6\"></path> <path d=\"M8 4H4v4\"></path> <path d=\"M4 4l6 6\"></path></g></svg></div></div></header><div class=\"w-full h-full flex justify-center font-extrabold text-lg items-center bg-yellow-300 flex-1\">MainArea</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
