import React from "react";

interface Props {
    label: string;
    name: string;
    type: string;
    value: string | number | [];
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    placeholder: string;
    error?: string;
    register?: any;
}

const InputWithLabel: React.FC<Props> = ({
    label,
    name,
    type,
    value,
    onChange,
    placeholder,
    error,
}) => {
    return (
        <div className="flex flex-row w-2/3 m-2 py-1 px-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-start items-center">
            <label htmlFor={name} className="w-32 m-2 border border-black rounded-md">
                {label}:
            </label>
            <input
                className="h-8 rounded ml-4 px-1 w-11/12 border-gray-200 shadow-sm sm:text-sm"
                id={name}
                name={name}
                type={type}
                value={value}
                onChange={onChange}
                placeholder={placeholder}
            />
            <div className="bg-red-600">{error && <div>{error}</div>}</div>
        </div>
    );
};

export default InputWithLabel;