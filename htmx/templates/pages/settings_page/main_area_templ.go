// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.819
package settings_page

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/yumenaka/comigo/htmx/state"
	"github.com/yumenaka/comigo/htmx/templates/common/svg"
)

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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<header id=\"header\" hx-target=\"#tab-contents\" role=\"tablist\" hx-on:htmx-after-on-load=\"let currentTab = document.querySelector(&#39;[aria-selected=true]&#39;);\n                                           currentTab.setAttribute(&#39;aria-selected&#39;, &#39;false&#39;)\n                                           currentTab.classList.remove(&#39;selected&#39;)\n                                           let newTab = event.target\n                                           newTab.setAttribute(&#39;aria-selected&#39;, &#39;true&#39;)\n                                           newTab.classList.add(&#39;selected&#39;)\" class=\"flex justify-between w-full h-12 py-1 border-b bg-base-100 text-base-content border-slate-400\"><a href=\"/\" class=\"flex justify-center items-center w-10 h-10 mx-1 my-0 rounded hover:ring\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = svg.Return().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</a><!-- examples: https://htmx.org/examples/tabs-javascript/--><div id=\"tabs\" class=\"flex flex-1  justify-center w-80 p-0 m-0 items-center font-semibold text-sm drop-shadow text-center focus:relative truncate\"><button role=\"tab\" aria-controls=\"tab-contents\" aria-selected=\"true\" hx-get=\"/htmx/settings/tab1\" class=\"selected flex items-center justify-center min-w-20 mx-0.5 my-2 h-9 rounded\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = svg.Book().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "<span x-text=\"i18next.t(&#39;book_shelf&#39;)\">Book</span></button> <button role=\"tab\" aria-controls=\"tab-contents\" aria-selected=\"false\" hx-get=\"/htmx/settings/tab2\" class=\"flex items-center justify-center min-w-20 mx-0.5 my-2 h-9 rounded\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = svg.Network().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "<span x-text=\"i18next.t(&#39;network&#39;)\">Network</span></button> <button role=\"tab\" aria-controls=\"tab-contents\" aria-selected=\"false\" hx-get=\"/htmx/settings/tab3\" class=\"flex items-center justify-center min-w-20 mx-0.5  my-2 h-9 rounded\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = svg.Labs().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "<span x-text=\"i18next.t(&#39;labs&#39;)\">Lab</span></button></div><style>\n            button {\n                background-color: #b1b5bb;\n                --tw-text-opacity: 1;\n                color: #6b7280; /* text-gray-500 */;\n            }\n            button:hover {\n                 --tw-text-opacity: 1;\n                 color: #374151 /* text-gray-700 */;\n            }\n            button.selected {\n                background-color: #f9f9f9;\n                --tw-text-opacity: 1;\n                color: #3b82f6; /* text-blue-500 */;\n            }\n        </style><!-- qrcode icon--><div data-modal-target=\"qrcode-modal\" data-modal-toggle=\"qrcode-modal\" class=\"flex justify-center items-center w-10 h-10 mx-1 my-0 rounded hover:ring\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = svg.QRCode().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</div></header><div class=\"flex flex-col justify-start items-center flex-1  w-full  h-full w-3/5 min-w-[24rem] font-semibold text-lg text-base-content\" :class=\"(theme.toString() ===&#39;light&#39;||theme.toString() ===&#39;dark&#39;||theme.toString() ===&#39;retro&#39;||theme.toString() ===&#39;lofi&#39;||theme.toString() ===&#39;nord&#39;) ? ($store.global.bgPattern !== &#39;none&#39;?$store.global.bgPattern+&#39; bg-base-300&#39;:&#39;bg-base-300&#39;):($store.global.bgPattern !== &#39;none&#39;?$store.global.bgPattern:&#39;&#39;)\"><div id=\"tab-contents\" role=\"tabpanel\" hx-get=\"/htmx/settings/tab1\" hx-trigger=\"load\">Loading……</div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
