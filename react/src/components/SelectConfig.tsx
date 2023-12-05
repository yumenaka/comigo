import { useState } from "react";

type PropsType = {
    name: string
    label: string

    InterfaceColor: string
}

const SelectConfig = (props: PropsType) => {
    const { name: valueName, label,  InterfaceColor } = props

    const optionalValue = ["RAM","Home","NowDir","ProgramDir"]
    //RAM（内存，临时生效，程序关闭后消失）。
    //Home（家目录，这个目录中的配置，每次启动时候都会生效)。
    //NowDir（命令执行目录，在此文件夹下面执行，配置会被读取）
    //Program（程序所在目录，每次启动时生效。适合制作便携版。）
    const [value, setValue] = useState("RAM")
    const [description, setDescription] = useState("RAM"); // 选中的选项
   

    const handleSelectChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        let saveTo ="RAM"
        if(event.target.value.startsWith("RAM")){
            saveTo = "RAM"
            setDescription("RAM（启动过程中生效，程序关闭后消失）")
        }
        if(event.target.value.startsWith("Home")){
            saveTo = "Home"
            setDescription("Home（保存到Home目录，每次启动时候都被读取）")
        }
        if(event.target.value.startsWith("NowDir")){
            saveTo = "NowDir"
            setDescription("NowDir（保存到命令执行目录，在此文件夹下面执行，对应配置会被读取）")
        }
        if(event.target.value.startsWith("ProgramDir")){
            saveTo = "ProgramDir"
            setDescription("ProgramDir（保存到程序所在目录，每次启动时读取，适合制作便携版。）")
        }
        console.log("saveTo",saveTo);
        setValue(saveTo);

        //props.setSelectedOption(props.name, saveTo); // 更新 selectedOption 状态
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
            <div className="py-1 w-3/4 text-xs text-gray-500">{description}</div>
        </div>

    )
}

export default SelectConfig