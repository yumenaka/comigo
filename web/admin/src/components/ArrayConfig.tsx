import React from "react";
import { useTranslation } from "react-i18next";
interface Props {
    label: string;
    fieldDescription: string;
    name: string;
    value: string[];
    setStringArray: (name: string, value: string[]) => void;
    error?: string;
    InterfaceColor?: string;
    showDialog?: (title: string, content: string) => void;
}

const ArrayConfig: React.FC<Props> = ({
    label,
    fieldDescription,
    name,
    value,
    setStringArray,
    error,
    InterfaceColor,
    showDialog,
}) => {
    const { t } = useTranslation();
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
        } else {
            showDialog && showDialog(t("hint"), t("please_enter_content"))
        }
    }

    return (
        <div
            style={{
                backgroundColor: InterfaceColor, // 绑定样式
            }}
            className={`w-full m-1 p-2 flex flex-col shadow-md hover:shadow-2xl font-semibold rounded-md justify-start items-left`}>
            <label className="py-0 w-32" htmlFor={name}>
                {label}
            </label>
            <div className="py-1 w-3/4 flex flex-row flex-wrap">
                {/* {value.toString()} */}
                {value.map((item, index) => (
                    <div key={index} className="p-2  m-1 flex flex-row items-center rounded-2xl bg-blue-300 text-sm font-medium text-black">
                        {item}
                        {/* https://www.xicons.org/#/ */}
                        <svg onClick={() => remove(index)}
                            className="mx-1 h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><g fill="none"><path d="M11.5 4a1.5 1.5 0 0 0-3 0h-1a2.5 2.5 0 0 1 5 0H17a.5.5 0 0 1 0 1h-.554L15.15 16.23A2 2 0 0 1 13.163 18H6.837a2 2 0 0 1-1.987-1.77L3.553 5H3a.5.5 0 0 1-.492-.41L2.5 4.5A.5.5 0 0 1 3 4h8.5zm3.938 1H4.561l1.282 11.115a1 1 0 0 0 .994.885h6.326a1 1 0 0 0 .993-.885L15.438 5zM8.5 7.5c.245 0 .45.155.492.359L9 7.938v6.125c0 .241-.224.437-.5.437c-.245 0-.45-.155-.492-.359L8 14.062V7.939c0-.242.224-.438.5-.438zm3 0c.245 0 .45.155.492.359l.008.079v6.125c0 .241-.224.437-.5.437c-.245 0-.45-.155-.492-.359L11 14.062V7.939c0-.242.224-.438.5-.438z" fill="currentColor"></path></g></svg>
                    </div>
                ))}

                {/* https://www.hyperui.dev/components/application-ui/inputs */}
                <div className="relative">
                    <label htmlFor="Array" className="sr-only">
                        {t("type_or_paste_content")}
                    </label>
                    <input
                        type="text"
                        id="Array"
                        placeholder={t("type_or_paste_content")}
                        className="w-full rounded-md border-gray-400 py-2.5 pe-10 shadow-sm sm:text-sm"
                        onKeyDown={onEnter}
                    ></input>
                    <span className="absolute top-[0px] right-[-80px] place-content-center">
                        <button onClick={onClick} type="button" className="h-10 w-16 mx-2 my-1 bg-sky-300 border border-gray-500 text-center text-gray-700 transition hover:text-gray-900 rounded">
                            {t("submit")}
                        </button>
                    </span>
                </div>
            </div>
            <div className="ml-2 py-1 w-3/4 text-xs text-gray-500">{fieldDescription}</div>
            <div className="bg-red-600">{error && <div>{error}</div>}</div>
        </div>
    );
};

export default ArrayConfig;