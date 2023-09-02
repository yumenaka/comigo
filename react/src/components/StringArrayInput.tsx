import React from "react";
import useArray from "../tools/useArray"
interface Props {
    label: string;
    name: string;
    type: string;
    value: string;
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

    const arr = value.split(",")
    const { array, set, push, remove, filter, update, clear } = useArray(arr)

    return (
        <div className="flex flex-row w-2/3 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-start items-center">

            <label htmlFor={name} className="w-32 border border-black rounded-md">
                {label}:
            </label>
            <div>{array.join(", ")}</div>
            <button onClick={() => push("7")}>Add 7</button>
            <button onClick={() => update(1, "9")}>Change Second Element To 9</button>
            <button onClick={() => remove(1)}>Remove Second Element</button>
            <button onClick={() => filter(n => n === "3")}>
                Keep Numbers Less Than 4
            </button>
            <button onClick={() => set(["1", "2"])}>Set To 1, 2</button>
            <button onClick={clear}>Clear</button>

            <input
                className="rounded ml-4 px-1 w-11/12"
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

export default StringArrayInput;