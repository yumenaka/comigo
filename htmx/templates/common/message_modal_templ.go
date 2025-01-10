// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.819
package common

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func MessageModal() templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!-- 模态遮罩层 --><div id=\"modal-overlay\" class=\"z-50 fixed inset-0 flex items-center justify-center hidden bg-black bg-opacity-50\"><!-- 模态窗口 --><div id=\"modal\" class=\"max-w-sm p-6 mx-auto bg-white rounded-lg\"><!-- 提示信息 --><p id=\"modal-message\" class=\"text-lg text-gray-800\"></p><!-- 按钮容器 --><div id=\"modal-buttons\" class=\"flex justify-end mt-4 space-x-2\"><!-- 按钮将动态插入 --></div></div></div><script>\nfunction showMessage(options) {\n  const overlay = document.getElementById('modal-overlay');\n  const messageElem = document.getElementById('modal-message');\n  const buttonsContainer = document.getElementById('modal-buttons');\n\n  // 设置提示信息\n  messageElem.textContent = options.message || '默认提示信息';\n\n  // 清空之前的按钮\n  buttonsContainer.innerHTML = '';\n\n  // 根据选项创建按钮\n  if (options.buttons === 'confirm') {\n    const confirmButton = document.createElement('button');\n    confirmButton.textContent = '确定';\n    confirmButton.className = 'px-4 py-2 bg-blue-500 text-white rounded';\n    confirmButton.onclick = function() {\n      hideModal();\n      if (typeof options.onConfirm === 'function') {\n         options.onConfirm();\n      }\n    };\n    buttonsContainer.appendChild(confirmButton);\n  } else if (options.buttons === 'yesno') {\n    const yesButton = document.createElement('button');\n    yesButton.textContent = '是';\n    yesButton.className = 'px-4 py-2 bg-green-500 text-white rounded';\n    yesButton.onclick = function() {\n      hideModal();\n      if (typeof options.onYes === 'function') {\n        options.onYes();\n      }\n    };\n    const noButton = document.createElement('button');\n    noButton.textContent = '否';\n    noButton.className = 'px-4 py-2 bg-gray-500 text-white rounded';\n    noButton.onclick = function() {\n      hideModal();\n    };\n    buttonsContainer.appendChild(yesButton);\n    buttonsContainer.appendChild(noButton);\n  }\n\n  // 显示模态\n  overlay.classList.remove('hidden');\n}\n\nfunction hideModal() {\n  const overlay = document.getElementById('modal-overlay');\n  overlay.classList.add('hidden');\n}\n\n</script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
