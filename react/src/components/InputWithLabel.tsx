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
        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-start items-center">
            <label htmlFor={name} className="w-32 border border-black rounded-md">
                {label}:
            </label>
            <input
                className="rounded ml-2 px-1"
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