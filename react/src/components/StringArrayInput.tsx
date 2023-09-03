import React from "react";

interface Props {
    label: string;
    name: string;
    type: string;
    value: string[];
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    placeholder: string;
    error?: string;
    register?: any;
}

const StringArrayInput: React.FC<Props> = ({
    label,
    name,
    type,
    value,
    onChange,
    placeholder,
    error,
}) => {

    // const { array, push, remove} = useArray(value)
    // console.log(value)
    // console.log(array)

    return (
        <div
            className="flex flex-row w-2/3 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-start items-center">
            <label htmlFor={name} className="w-32 border border-black rounded-md">
                {label}:
            </label>
            <div className="flex flex-row mx-4">
                {/* {value.toString()} */}
                {value.map((item, index) => (
                    <div key={index} className="flex flex-row  p-1 m-2  items-center rounded-2xl bg-gray-200 px-8 py-3 text-sm font-medium text-black">
                        {item}
                        <div
                            className="mx-2 h-6 w-6 rounded-2xl bg-red-700 text-center py-0.5 text-white transition hover:scale-140 hover:shadow-xl focus:outline-none focus:ring active:bg-red-500">
                            X
                        </div>
                    </div>
                ))}
            </div>
            <div className="bg-red-600">{error && <div>{error}</div>}</div>
        </div>
    );
};

export default StringArrayInput;