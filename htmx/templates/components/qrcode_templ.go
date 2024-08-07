// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/yumenaka/comi/htmx/state"

func serverHostBindStr(serverHost string) string {
	//"{ serverHost: 'abc.com' }"
	return "{ serverHost: '" + serverHost + "' }"
}

func QRCode(s *state.GlobalState) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!-- Main modal --><div x-data=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(serverHostBindStr(s.ServerStatus.ServerHost))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/qrcode.templ`, Line: 12, Col: 61}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" id=\"qrcode-modal\" tabindex=\"-1\" aria-hidden=\"true\" class=\"hidden overflow-y-auto overflow-x-hidden fixed top-0 right-0 left-0 z-50 justify-center items-center w-full md:inset-0 h-[calc(100%-1rem)] max-h-full\"><!-- Modal content --><div class=\"relative w-64 h-64 p-1 m-1\" x-data=\"{ qrcodeSrc: &#39;&#39; }\" x-init=\"qrcodeSrc = window.location.origin +&#39;/api/qrcode.png?qrcode_str=&#39;+ encodeURIComponent(window.location.toString().replace(window.location.hostname,serverHost))\"><img class=\"w-64 h-64\" :src=\"qrcodeSrc\"></img> <a x-ref=\"content\" :href=\"window.location.toString().replace(window.location.hostname,serverHost)\" target=\"_blank\"><div class=\"min-w-64 p-1 mb-4 text-center text-white text-xs font-semibold\" x-text=\"window.location.toString().replace(window.location.hostname,serverHost)\"></div></a> <button type=\"button\" class=\"absolute buttom-0 left-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-primary m-2 p-2 rounded text-sm self-center\" x-on:click=\"navigator.clipboard.writeText($refs.content.href);alert(&#39;You copy it&#39;);\">Copy URL</button></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
