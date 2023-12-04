
type PropsType = {
    name: string
    label: string
    fieldDescription: string
    boolValue: boolean
    InterfaceColor: string
    setBoolValue: (valueName: string, checked: boolean) => void
    showDialogModal?: (message: string) => void
}

const BoolConfig = (props: PropsType) => {
    const { name: valueName, label, fieldDescription, boolValue, InterfaceColor, setBoolValue,showDialogModal} = props
    const handleCheckboxChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        console.log(valueName, event.target.checked)
        setBoolValue(valueName, event.target.checked)
        if (showDialogModal) {
            showDialogModal("Checkbox value changed");
        }
    }

    return (
        <div
            style={{
                backgroundColor: InterfaceColor, // 绑定样式
            }}
            className={`w-full m-1 p-2 flex flex-col shadow-md hover:shadow-2xl font-semibold rounded-md  justify-left items-left`}>
            <div className="w-32">{label}</div>
            <label htmlFor={valueName} className="relative h-8 w-14 cursor-pointer">
                <input type="checkbox" checked={boolValue} id={valueName} className="peer sr-only" onChange={handleCheckboxChange} />
                <span
                    className="absolute inset-0 rounded-full bg-gray-300 transition peer-checked:bg-green-500"
                ></span>
                <span
                    className="absolute inset-y-0 start-0 m-1 h-6 w-6 rounded-full bg-white transition-all peer-checked:start-6"
                ></span>
            </label>
            <div className="py-1 w-3/4 text-xs text-gray-500">{fieldDescription}</div>
        </div>
    )
}

export default BoolConfig
