// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.819
package common

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Drawer(serverHost string, slot templ.Component) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!-- drawer component --><!-- https://flowbite.com/docs/components/drawer/ --><div id=\"drawer-right\" x-data=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(ServerHostBindStr(serverHost))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/common/drawer.templ`, Line: 8, Col: 40}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "\" class=\"fixed flex flex-col top-0 right-0 z-40 w-64 h-dvh p-4 overflow-y-auto transition-transform translate-x-full bg-base-100 text-base-content\" tabindex=\"-1\" aria-labelledby=\"drawer-right-label\"><div class=\"w-full m-0 p-0 cursor-pointer border-slate-500 border-double border-b-2\"><h5 id=\"drawer-right-label\" class=\"shadow-lg inline-flex items-center mb-4 text-base font-semibold\"><div x-text=\"i18next.t(&#39;reader_settings&#39;)\"></div><img class=\"shadow-xl mx-2 w-6 h-6 rounded-full border border-stone-500 border-dotted outline-offset-4\" src=\"/static/favicon.png\" alt=\"comigo favicon\"></h5><button type=\"button\" data-drawer-hide=\"drawer-right\" aria-controls=\"drawer-right\" class=\"font-bold rounded text-sm w-8 h-8 absolute top-2.5 end-2.5 inline-flex items-center justify-center bg-transparent hover:ring-2 hover:ring-blue-300 dark:hover:ring-gray-300 dark:hover:text-white\"><svg class=\"w-8 h-8\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 1024 1024\"><path d=\"M685.4 354.8c0-4.4-3.6-8-8-8l-66 .3L512 465.6l-99.3-118.4l-66.1-.3c-4.4 0-8 3.5-8 8c0 1.9.7 3.7 1.9 5.2l130.1 155L340.5 670a8.32 8.32 0 0 0-1.9 5.2c0 4.4 3.6 8 8 8l66.1-.3L512 564.4l99.3 118.4l66 .3c4.4 0 8-3.5 8-8c0-1.9-.7-3.7-1.9-5.2L553.5 515l130.1-155c1.2-1.4 1.8-3.3 1.8-5.2z\" fill=\"currentColor\"></path><path d=\"M512 65C264.6 65 64 265.6 64 513s200.6 448 448 448s448-200.6 448-448S759.4 65 512 65zm0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372s372 166.6 372 372s-166.6 372-372 372z\" fill=\"currentColor\"></path></svg></button></div><div class=\"drawer_slot flex flex-col flex-grow items-center justify-end p-1 my-2 rounded text-accent-content dark:text-white\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if slot != nil {
			templ_7745c5c3_Err = slot.Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "<div class=\"place-holder w-full flex-1\"></div><!-- 二维码 --><div class=\"p-1 mt-4 mb-2 w-36 h-36\" x-data=\"{ qrcodeSrc: &#39;&#39; }\" x-init=\"qrcodeSrc = window.location.origin +&#39;/api/qrcode.png?qrcode_str=&#39;+ encodeURIComponent(window.location.toString().replace(window.location.hostname,serverHost))\"><img class=\"w-32 h-32\" :src=\"qrcodeSrc\"></div><!-- 选择主题的select --><select x-model=\"theme\" x-on:change=\"theme = $event.target.value;console.log(theme);\" class=\"w-full h-10 mt-auto mb-2 border rounded bg-base-100 text-accent-content focus:outline-none\"><option value=\"retro\">Retro</option> <option value=\"light\">Light</option> <option value=\"dark\">Dark</option> <option value=\"dracula\">Dracula</option> <option value=\"cupcake\">Cupcake</option> <option value=\"cyberpunk\">Cyberpunk</option> <option value=\"valentine\">Valentine</option> <option value=\"aqua\">Aqua</option> <option value=\"lofi\">Lofi</option> <option value=\"halloween\">Halloween</option> <option value=\"coffee\">Coffee</option> <option value=\"winter\">Winter</option> <option value=\"nord\">Nord</option></select><!-- 选择背景花纹的select --><select x-show=\"$store.global.debugMode\" x-model=\"$store.global.bgPattern\" x-on:change=\"$store.global.bgPattern = $event.target.value;console.log($store.global.bgPattern);\" class=\"w-full h-10 mt-auto mb-2 border rounded bg-base-100 text-accent-content focus:outline-none\"><option value=\"none\">None</option> <option value=\"grid-line\">Grid Line</option> <option value=\"grid-point\">Grid Point</option> <option value=\"grid-mosaic\">Grid Mosaic</option></select><!-- 选择语言的select 此处需要与自动探测到的结果一致，所以才是 \"en-US\" \"zh-CN\" \"ja\"这种不统一的形式\"--><select x-model=\"i18next.language\" x-on:change=\"i18next.changeLanguage($event.target.value).then(location.reload())\" class=\"w-full h-10 border rounded bg-base-100 text-accent-content focus:outline-none\"><option value=\"en-US\">English</option> <option value=\"zh-CN\">中文</option> <option value=\"ja\">日本語</option></select></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
