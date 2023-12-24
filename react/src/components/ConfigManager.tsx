import { useState, } from "react"
import axios from "axios"
import ConfigStatus from "../types/ConfigStatus"
import { useEffectOnce } from 'react-use';
import { useTranslation } from "react-i18next";


type PropsType = {
    InterfaceColor: string
    showDialogFunc: (title: string, content: string) => void
}

const ConfigManager = (props: PropsType) => {
    const { InterfaceColor, showDialogFunc } = props
    const [config_status, setConfigStatus] = useState({
        In: "",
        Path: {
            WorkingDirectory: "",
            HomeDirectory: "",
            ProgramDirectory: ""
        }
    } as ConfigStatus);
    const { t } = useTranslation();

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
    const selectOption = [
        //https://icon-icons.com/icon/coding-program/71231
        { name: "WorkingDirectory", icon: "icon/working_directory.png", description: t("WorkingDirectory"), path: config_status.Path.WorkingDirectory },
        //https://icon-icons.com/icon/web-page-home/85808
        { name: "HomeDirectory", icon: "icon/home_directory.png", description: t("HomeDirectory"), path: config_status.Path.HomeDirectory },
        //https://icon-icons.com/icon/folder-sync-outline/139517
        { name: "ProgramDirectory", icon: "icon/program_directory.png", description: t("ProgramDirectory"), path: config_status.Path.ProgramDirectory },
    ];

    const onSelect = (event: React.MouseEvent) => {
        setSelected(event.currentTarget.getAttribute("data-save_to") ?? "");
    };

    const onSaveConfig = () => {
        for (let i = 0; i < selectOption.length; i++) {
            if (selected !== selectOption[i].name && selectOption[i].path !== "") {
                showDialogFunc(t("hint"), `【${selectOption[i].description}】`+t("ConfigManagerSaveHint"));
                return;
            }
        }

        axios
            .post<ConfigStatus | { error: string }>(`api/config/${selected}`)
            .then((response) => {
                if (response.status === 200) {
                    console.log(response.data);
                    setConfigStatus(response.data as ConfigStatus);
                    showDialogFunc(t("hint"), t("ConfigManagerSaveSuccess"));
                }
                if (response.status === 400) {
                    const error = response.data as { error: string };
                    console.log(error.error);
                    showDialogFunc(t("hint"), error.error);
                }
            })
            .catch((error) => {
                console.error(error);
            });
    }

    const onDeleteConfig = (event: React.MouseEvent) => {
        // get element data
        console.log(event.currentTarget.getAttribute("data-save_to"));
        axios
            .delete<ConfigStatus | { error: string }>(`api/config/${selected}`)
            .then((response) => {
                if (response.status === 200) {
                    console.log(response.data);
                    setConfigStatus(response.data as ConfigStatus);
                    showDialogFunc(t("hint"), t("ConfigManagerDeleteSuccess"));
                }
                if (response.status === 400) {
                    const error = response.data as { error: string };
                    console.log(error.error);
                    showDialogFunc(t("hint"), error.error);
                }
            })
            .catch((error) => {
                console.error(error);
            });
    }

    return (
        <div
            className="w-full m-1 p-2 flex flex-col shadow-md hover:shadow-2xl font-semibold rounded-md justify-center items-center"
            style={{
                backgroundColor: InterfaceColor,
            }}>
            <label className="py-0 w-full">
                {t("ConfigManager")}
            </label>
            <div className="flex flex-row  mx-0 my-1 w-full">
                {Object.entries(selectOption).map(([key, value]) => (
                    <div key={value.name} data-save_to={value.name} onClick={onSelect} className={"flex flex-col justify-center items-center text-xs font-normal pt-2 mx-1 w-1/3 min-h-20 border border-gray-500 rounded" + (selected === value.name ? " bg-cyan-200" : "")}>
                        <img className="h-7 w-7" src={value.icon} alt={key} />
                        <div className="mt-1">{value.description}</div>
                        {/* 超过两行，显示省略号。https://zenn.dev/ilove/articles/8a93705d396e05 */}
                        {selectOption.map(s =>
                            value.name === s.name &&
                            <div className="mx-1 my-1 text-xs text-gray-500 line-clamp-2 hover:line-clamp-none active:line-clamp-none">{s.path}</div>
                        )}
                    </div>
                ))}
            </div>
            <div className="flex flex-row mx-4">
                <button onClick={onSaveConfig} className="h-10 w-24 mx-2 my-1 bg-sky-300 border border-gray-500 text-center text-gray-700 transition hover:text-gray-900 rounded">SAVE</button>
                {selectOption.map(s =>
                    selected === s.name && s.path !== "" &&
                    <button onClick={onDeleteConfig} className="h-10 w-24 mx-2 my-1 bg-red-300 border border-gray-500 text-center text-gray-700 transition hover:text-gray-900 rounded">DELETE</button>
                )}
            </div>
        </div>
    )
}

export default ConfigManager