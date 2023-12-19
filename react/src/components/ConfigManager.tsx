import { useState, } from "react"
import axios from "axios"
import ConfigStatus from "../types/ConfigStatus"
import { useEffectOnce } from 'react-use';

type PropsType = {
    label: string
    InterfaceColor: string
    showDialogFunc: (title: string, content: string) => void
}

const ConfigManager = (props: PropsType) => {
    const { label, InterfaceColor,showDialogFunc } = props
    const [config_status, setConfigStatus] = useState({
        ConfigDirectory: "",
        Home: false,
        Execution: false,
        Program: false,
    } as ConfigStatus);

    // 获取comigo配置的状态
    // 可以用React query代替useEffectOnce，获得loading，error，retry,cache等功能。 https://reffect.co.jp/react/react-use-query 
    const updateConfigStatus = () => {
        axios
            .get<ConfigStatus>(`api/config/status`)
            .then((response) => {
                setConfigStatus(response.data);
            })
            .catch((error) => {
                console.error(error);
            });
    };
    useEffectOnce(() => {
        updateConfigStatus();
    });

    const [selected, setSelected] = useState("WorkingDirectory");
    const selectOption = {
        "WorkingDirectory": ["icon/working_directory.png", "当前运行目录"],//https://icon-icons.com/icon/coding-program/71231
        "HomeDirectory": ["icon/home_directory.png", "用户主目录"],//https://icon-icons.com/icon/web-page-home/85808
        "ProgramDirectory": ["icon/program_directory.png", "程序所在目录"]//https://icon-icons.com/icon/folder-sync-outline/139517
    }

    const onSelect = (event: React.MouseEvent) => {
        //console.log(event.currentTarget.getAttribute("data-save_to"));
        setSelected(event.currentTarget.getAttribute("data-save_to") ?? "RAM");
    };

    const onSaveConfig = () => {
        axios
            .post<ConfigStatus>(`api/config/${selected}`)
            .then((response) => {
                console.log(response.data);
                setConfigStatus(response.data);
                showDialogFunc("提示", "配置已保存");
            })
            .catch((error) => {
                console.error(error);
            });
    }

    const onDeleteConfig = (event: React.MouseEvent) => {
        // get element data
        console.log(event.currentTarget.getAttribute("data-save_to"));
        axios
            .delete<ConfigStatus>(`api/config/${selected}`)
            .then((response) => {
                console.log(response.data);
                setConfigStatus(response.data);
                showDialogFunc("提示", "配置已删除");
            })
            .catch((error) => {
                console.error(error);
            });
    }

    return (
        <div
            className="w-full m-1 p-2 flex flex-col shadow-md hover:shadow-2xl font-semibold rounded-md justify-center items-center"
            style={{
                backgroundColor: InterfaceColor, // 绑定样式
            }}>
            <label className="py-0 w-full">
                {label}
            </label>
            <div className="flex flex-row mx-0 my-1 w-full">
                {Object.entries(selectOption).map(([key, value]) => (
                    <div key={key} data-save_to={key} onClick={onSelect} className={"text-xs font-normal flex flex-col justify-center items-center p-1 mx-1 w-1/3 h-20 border border-gray-500 rounded" + (selected === key ? " bg-cyan-200" : "")}>
                        <img className="h-7 w-7" src={value[0]} alt={key} />
                        <div className="mt-1">{value[1]}</div>
                        {config_status.ConfigDirectory === "HomeDirectory" && <div className="my-1 text-xs text-gray-500">配置生效</div>}
                    </div>
                ))}
            </div>
            <div className="flex flex-row mx-4">
                <button onClick={onSaveConfig} className="h-10 w-24 mx-2 my-1 bg-sky-300 border border-gray-500 text-center text-gray-700 transition hover:text-gray-900 rounded">SAVE</button>
                <button onClick={onDeleteConfig} className="h-10 w-24 mx-2 my-1 bg-red-300 border border-gray-500 text-center text-gray-700 transition hover:text-gray-900 rounded">DELETE</button>
            </div>

        </div>
    )
}

export default ConfigManager