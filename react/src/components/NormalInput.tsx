import React from "react";
interface Props {
    label: string;
    fieldDescription: string;
    name: string;
    type: string;
    value: string | number | [];
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    placeholder: string;
    error?: string;
    InterfaceColor?: string;
    // register?: any;
}

const NormalConfig = ({
    label,
    fieldDescription,
    name,
    type,
    value,
    onChange,
    placeholder,
    error,
    InterfaceColor,
}: Props) => {

    return (
        <div
        style={{
            backgroundColor: InterfaceColor, // 绑定样式
        }} 
        className={`w-full m-1 p-2 flex flex-col shadow-md hover:shadow-2xl font-semibold rounded-md justify-start items-left`}>
            <label htmlFor={name} className="w-64">
                {label}:
            </label>
            <input
                className="px-1 w-64 rounded-md border-gray-400 py-2.5 pe-10 shadow-sm sm:text-sm"
                id={name}
                name={name}
                type={type}
                value={value}
                onChange={onChange}
                placeholder={placeholder}
                // onBlur={handleSubmit} // 使用 onBlur 事件来监听输入框的失去焦点事件
            />
            <div className="py-1 w-3/4 text-xs text-gray-500">{fieldDescription}</div>
            <div className="bg-red-600">{error && <div>{error}</div>}</div>
        </div>
    );
};

export default NormalConfig;