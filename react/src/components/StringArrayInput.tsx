import React from "react";

interface Props {
    label: string;
    name: string;
    value: string[];
    setStringArray: (name: string, value: string[]) => void;
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

    function onEnter(e: React.KeyboardEvent<HTMLInputElement>): void {
        if (e.key === 'Enter') {
            push(e.currentTarget.value)
            e.currentTarget.value = ''
        }
    }
    function onClick(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void {
        const input = e.currentTarget.parentElement?.parentElement?.firstElementChild?.nextElementSibling as HTMLInputElement
        if (input.value !== '') {
            push(input.value)
            input.value = ''
        }
    }

    return (
        <div
            className="m-1  w-2/3  flex flex-col font-semibold rounded-md shadow-md bg-yellow-100 justify-start items-left">
            <label className="ml-4 py-1 w-32" htmlFor={name}>
                {label}:
            </label>
            <div className="ml-4 py-1 w-3/4 flex flex-row flex-wrap">
                {/* {value.toString()} */}
                {value.map((item, index) => (
                    <div key={index} className="px-2 py-1 m-1 flex flex-row items-center rounded-2xl bg-gray-200  text-sm font-medium text-black">
                        {item}
                        {/* https://www.xicons.org/#/ */}
                        <svg onClick={() => remove(index)}
                            className="mx-1 h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><g fill="none"><path d="M11.5 4a1.5 1.5 0 0 0-3 0h-1a2.5 2.5 0 0 1 5 0H17a.5.5 0 0 1 0 1h-.554L15.15 16.23A2 2 0 0 1 13.163 18H6.837a2 2 0 0 1-1.987-1.77L3.553 5H3a.5.5 0 0 1-.492-.41L2.5 4.5A.5.5 0 0 1 3 4h8.5zm3.938 1H4.561l1.282 11.115a1 1 0 0 0 .994.885h6.326a1 1 0 0 0 .993-.885L15.438 5zM8.5 7.5c.245 0 .45.155.492.359L9 7.938v6.125c0 .241-.224.437-.5.437c-.245 0-.45-.155-.492-.359L8 14.062V7.939c0-.242.224-.438.5-.438zm3 0c.245 0 .45.155.492.359l.008.079v6.125c0 .241-.224.437-.5.437c-.245 0-.45-.155-.492-.359L11 14.062V7.939c0-.242.224-.438.5-.438z" fill="currentColor"></path></g></svg>
                    </div>
                ))}

                {/* https://www.hyperui.dev/components/application-ui/inputs */}
                <div className="relative">
                    <label htmlFor="Search" className="sr-only"> Search </label>
                    <input
                        type="text"
                        id="Search"
                        placeholder="Add new..."
                        className="w-full rounded-md border-gray-200 py-2.5 pe-10 shadow-sm sm:text-sm"
                        onKeyDown={onEnter}
                    ></input>

                    <span className="absolute inset-y-0 end-0 grid w-10 place-content-center">
                        <button onClick={onClick} type="button" className="text-gray-600 hover:text-gray-700">
                            <span className="sr-only">Search</span>
                            <svg className="h-6 w-6" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path d="M448 256c0-106-86-192-192-192S64 150 64 256s86 192 192 192s192-86 192-192z" fill="none" stroke="currentColor" stroke-miterlimit="10" stroke-width="32"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M256 176v160"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M336 256H176"></path></svg>
                        </button>
                    </span>
                </div>
            </div>

            <div className="bg-red-600">{error && <div>{error}</div>}</div>
        </div>
    );
};

export default StringArrayInput;