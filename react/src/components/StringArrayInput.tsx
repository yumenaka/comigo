import React from "react";

interface Props {
    label: string;
    name: string;
    value: string[];
    setStringArray: (name: string, value:string[]) => void;
    error?: string;
    // register?: any;
}

const StringArrayInput: React.FC<Props> = ({
    label,
    name,
    value,
    setStringArray,
    error,
}) => {


    function push(element: string): void {
        setStringArray(name, [...value, element])
    }

    function remove(index: number): void {
        setStringArray(name, [...value.slice(0, index), ...value.slice(index + 1)])
    }

    return (
        <div
            className="flex flex-row w-2/3 font-semibold rounded-md shadow-md bg-yellow-300 justify-start items-center">
            <label htmlFor={name} className="w-32 border border-black rounded-md">
                {label}:
            </label>
            <div className="flex flex-row flex-wrap mx-4">
                {/* {value.toString()} */}
                {value.map((item, index) => (
                    <div key={index} className="flex flex-row p-2 m-1 items-center rounded-2xl bg-gray-200  text-sm font-medium text-black">
                        {item}
                        <svg onClick={() => remove(index)}
                            className="mx-1 h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><g fill="none"><path d="M11.5 4a1.5 1.5 0 0 0-3 0h-1a2.5 2.5 0 0 1 5 0H17a.5.5 0 0 1 0 1h-.554L15.15 16.23A2 2 0 0 1 13.163 18H6.837a2 2 0 0 1-1.987-1.77L3.553 5H3a.5.5 0 0 1-.492-.41L2.5 4.5A.5.5 0 0 1 3 4h8.5zm3.938 1H4.561l1.282 11.115a1 1 0 0 0 .994.885h6.326a1 1 0 0 0 .993-.885L15.438 5zM8.5 7.5c.245 0 .45.155.492.359L9 7.938v6.125c0 .241-.224.437-.5.437c-.245 0-.45-.155-.492-.359L8 14.062V7.939c0-.242.224-.438.5-.438zm3 0c.245 0 .45.155.492.359l.008.079v6.125c0 .241-.224.437-.5.437c-.245 0-.45-.155-.492-.359L11 14.062V7.939c0-.242.224-.438.5-.438z" fill="currentColor"></path></g></svg>
                    </div>
                ))}
            </div>
            <input className="border border-black rounded-md" type="text" name={name} id={name} placeholder="Add new item" onKeyDown={(e) => {
                if (e.key === 'Enter') {
                    push(e.currentTarget.value)
                    e.currentTarget.value = ''
                }
            }}
            ></input>
            <div className="bg-red-600">{error && <div>{error}</div>}</div>
        </div>
    );
};

export default StringArrayInput;