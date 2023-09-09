import { useEffect, useState } from "react";

type PropsType = {
    name: string
    nameText: string
    BoolValue: boolean
    onChange: (checked: boolean) => void;
}

const NumberInput = (props: PropsType) => {
    const [InterfaceColor, setInterfaceColor] = useState("bg-[#F5F5E4]")
    useEffect(() => {
        // 当前颜色
        const tempInterfaceColor = localStorage.getItem("InterfaceColor");
        if (tempInterfaceColor !== null) {
            setInterfaceColor(tempInterfaceColor)
        }
    }, []);
    return (
        <div className={`w-full m-1 p-2 flex flex-row font-semibold rounded-md shadow-md justify-start items-center ${InterfaceColor}`}>
            <label htmlFor="Port" className="w-32 border border-black rounded-md">
                {props.nameText}:
            </label>
            <input
                className="px-1 w-32 rounded-md border-gray-400 py-2.5 pe-10 shadow-sm sm:text-sm"
                id="Port"
                type="number"
                placeholder="Port"
            />
            <div className="bg-red-600">
                {/* {errors.Port && <div>入力が必須の項目です(0~65535)</div>} */}
            </div>
        </div>
    );
};

export default NumberInput;