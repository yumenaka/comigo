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
        <div className="w-full h-16 mb-1 rounded shadow flex flex-row justify-center items-center" style={{
            backgroundColor: InterfaceColor, // 绑定样式
        }}>
            <Contained group={show} setGroup={setShow} InterfaceColor={InterfaceColor} />
        </div>
    );
}

