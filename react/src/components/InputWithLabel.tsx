import React from "react";

interface InputWithLabelProps {
    label: string;
    fieldDescription: string;
    name: string;
    type: string;
    value: string | number | [];
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    placeholder: string;
    error?: string;
    // register?: any;
}



const InputWithLabel = ({
    label,
    fieldDescription,
    name,
    type,
    value,
    onChange,
    placeholder,
    error,
}: InputWithLabelProps) => {

    // const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    //     event.preventDefault();
    //     console.log("Form submitted!");
    // };

    return (
        <div className="m-2 py-2 px-4 flex flex-col w-2/3  font-semibold rounded-md shadow-md bg-blue-100 justify-start items-left">
            <label htmlFor={name} className="w-64">
                {label}:
            </label>
            <input
                className="h-8  px-1 w-24 border border-black rounded-md shadow-sm sm:text-sm"
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

export default InputWithLabel;