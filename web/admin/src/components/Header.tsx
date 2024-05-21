// Header.tsx
import Contained from "./Contained";
import React from "react";
interface PropsType {
    group: string;
    setGroup: React.Dispatch<React.SetStateAction<string>>;
    InterfaceColor: string;
}

export default function Header(props: PropsType) {
    const { group: show, setGroup: setShow, InterfaceColor, } = props
    return (
        <div className="relative  w-full h-16 mb-1 rounded shadow flex flex-row justify-center items-center" style={{
            backgroundColor: InterfaceColor, // 绑定样式
        }}>
            {/* 返回主页的按钮 */}
            <a
                className="absolute left-2  rounded-full border  shadow-md hover:shadow-2xl bg-white border-indigo-600 p-3 text-indigo-600 hover:bg-indigo-600 hover:text-white focus:outline-none focus:ring active:bg-indigo-500"
                href="/"
            >
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18" />
                </svg>
            </a>
            <Contained group={show} setGroup={setShow} InterfaceColor={InterfaceColor} />
        </div>
    );
}

