import React from "react";

interface Props {
    label: string;
    name: string;
    type: string;
    value: string | number | [];
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    placeholder: string;
    error?: string;
    // register?: any;
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
        <div className="m-2 py-2 px-4 flex flex-col w-2/3  font-semibold rounded-md shadow-md bg-yellow-100 justify-start items-left">
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
            />
            <div className="bg-red-600">{error && <div>{error}</div>}</div>
        </div>
    );
};

export default InputWithLabel;