import { useState, useEffect, useReducer,  } from "react"
import axios from "axios"
import ConfigStatus from "../types/ConfigStatus"
import { configStatusReducer, defaultConfigStatus } from "../reducers/configStatusReducer"


type PropsType = {
    name: string
    label: string
    InterfaceColor: string
}

const ConfigManager = (props: PropsType) => {
    const { name: valueName, label, InterfaceColor } = props
    const [config_status, config_status_dispatch] = useReducer(configStatusReducer, defaultConfigStatus);

    const saveToPos = {
        "用户目录": "用户的主目录",
        "执行目录": "命令运行目录",
        "程序目录": "程序所在目录"
    }
    const [saveTo, setSaveTo] = useState("HOME");
    
    //useEffect
    useEffect(() => {
        // comigo配置状态
        axios
            .get<ConfigStatus>(`api/config/status`)
            .then((response) => {
                config_status_dispatch({
                    type: 'init',
                    name: "",
                    value: "",
                    config: response.data
                });
                //console.log("config_status", config_status);
            })
            .catch((error) => {
                console.error(error);
            });
    }, []);


    const handleSaveConfig = (event: React.ChangeEvent<HTMLSelectElement>) => {
        if (event.target.value.startsWith("HomeDir")) {
            setSaveTo("RAM");
        }
        if (event.target.value.startsWith("NowDir")) {
            setSaveTo("RAM");
        }
        if (event.target.value.startsWith("ProgramDir")) {
            setSaveTo("RAM");
        }
        console.log("saveTo", saveTo);
        setSaveTo(saveTo);
    };

    return (
        <div
            className="w-full m-1 p-2 flex flex-col shadow-md hover:shadow-2xl font-semibold rounded-md justify-left items-left"
            style={{
                backgroundColor: InterfaceColor, // 绑定样式
            }}>
            <label htmlFor={valueName} className="py-0 w-32">
                {label}
            </label>
            <div className="flex flex-row mx-0 my-1">
                {Object.entries(saveToPos).map(([key, value]) => (
                    <div className="text-xs font-normal flex flex-col justify-center items-center p-1 mx-1 w-64 h-16 border border-gray-500 rounded">
                        <div key={key}>
                            {key}
                        </div>
                        <div>{value}</div>
                    </div>
                ))}
            </div>
        </div>

    )
}

export default ConfigManager