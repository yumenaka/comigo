
type PropsType = {
    name: string
    label: string
    fieldDescription: string
    value: string
    optionalValue: string[]
    InterfaceColor: string
    setSelectedOption: (valueName: string, value: string) => void
}

const SelectConfig = (props: PropsType) => {
    const { name: valueName, label, fieldDescription, value, optionalValue, InterfaceColor } = props

    const handleSelectChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        props.setSelectedOption(props.name, event.target.value); // 更新 selectedOption 状态
    };

    return (
        <div
        className="w-full m-1 p-2 flex flex-col shadow-md hover:shadow-2xl font-semibold rounded-md  justify-left items-left"
        style={{
            backgroundColor: InterfaceColor, // 绑定样式
        }}>
            <label htmlFor={valueName} className="py-0 w-32">
                {label}
            </label>
            <select
                value={value} // 绑定 value 值
                onChange={handleSelectChange}
                name={valueName}
                id={valueName}
                className="mt-1.5 w-full rounded-lg border-gray-300 text-gray-700 sm:text-sm"
            >
                {/* 使用 option 元素来定义选项 */}
                {optionalValue.map((option, index) => (
                    <option key={index} value={option}>{option}</option>
                ))}
            </select>
            <div className="py-1 w-3/4 text-xs text-gray-500">{fieldDescription}</div>
        </div>

    )
}

export default SelectConfig