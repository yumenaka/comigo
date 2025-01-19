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
            .get<ConfigStatus>(`/api/config/status`)
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
        { name: "WorkingDirectory", icon: "/admin/icon/working_directory.png", description: t("WorkingDirectory"), path: config_status.Path.WorkingDirectory },
        //https://icon-icons.com/icon/web-page-home/85808
        { name: "HomeDirectory", icon: "/admin/icon/home_directory.png", description: t("HomeDirectory"), path: config_status.Path.HomeDirectory },
        //https://icon-icons.com/icon/folder-sync-outline/139517
        { name: "ProgramDirectory", icon: "/admin/icon/program_directory.png", description: t("ProgramDirectory"), path: config_status.Path.ProgramDirectory },
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
            .post<ConfigStatus | { error: string }>(`/api/config/${selected}`)
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
            .delete<ConfigStatus | { error: string }>(`/api/config/${selected}`)
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
            className="flex flex-col items-center justify-center w-full p-2 m-1 font-semibold rounded-md shadow-md hover:shadow-2xl"
            style={{
                backgroundColor: InterfaceColor,
            }}>
            <label className="w-full py-0">
                {t("ConfigManager")}
            </label>
            <div className="flex flex-row w-full mx-0 my-1">
                {Object.entries(selectOption).map(([key, value]) => (
                    <div key={key} data-save_to={value.name} onClick={onSelect} className={"flex flex-col justify-center items-center text-xs font-normal pt-2 mx-1 w-1/3 min-h-20 border border-gray-500 rounded" + (selected === value.name ? " bg-cyan-200" : "")}>
                        <img className="h-7 w-7" src={value.icon} key={key} />
                        <div className="mt-1">{value.description}</div>
                        {/* 超过两行，显示省略号。https://zenn.dev/ilove/articles/8a93705d396e05 */}
                        {selectOption.map(s =>
                            value.name === s.name &&
                            <div key={s.name} className="mx-1 my-1 text-xs text-gray-500 line-clamp-2 hover:line-clamp-none active:line-clamp-none">{s.path}</div>
                        )}
                    </div>
                ))}
            </div>
            <div className="flex flex-row mx-4">
                <button onClick={onSaveConfig} className="w-24 h-10 mx-2 my-1 text-center text-gray-700 transition border border-gray-500 rounded bg-sky-300 hover:text-gray-900">SAVE</button>
                {selectOption.map(s =>
                    selected === s.name && s.path !== "" &&
                    <button key={s.name} onClick={onDeleteConfig} className="w-24 h-10 mx-2 my-1 text-center text-gray-700 transition bg-red-300 border border-gray-500 rounded hover:text-gray-900">DELETE</button>
                )}
            </div>
            <div className="w-full py-1 text-xs text-gray-500">{t("ConfigManagerDescription")}</div>
        </div>
        
    )
}

export default ConfigManager













  <div class="style-scope ytd-topbar-logo-renderer">
    <ytd-logo class="style-scope ytd-topbar-logo-renderer" enable-refresh-ringo2-web=""><!--css-build:shady--><!--css-build:shady--><yt-icon id="logo-icon" class="style-scope ytd-logo"><!--css-build:shady--><!--css-build:shady--><span class="yt-icon-shape style-scope yt-icon yt-spec-icon-shape"><div style="width: 100%; height: 100%; display: block; fill: currentcolor;"><svg xmlns="http://www.w3.org/2000/svg" id="yt-ringo2-svg_yt24" width="93" height="20" viewBox="0 0 93 20" focusable="false" aria-hidden="true" style="pointer-events: none; display: inherit; width: 100%; height: 100%;">
  <g>
    <path d="M14.4848 20C14.4848 20 23.5695 20 25.8229 19.4C27.0917 19.06 28.0459 18.08 28.3808 16.87C29 14.65 29 9.98 29 9.98C29 9.98 29 5.34 28.3808 3.14C28.0459 1.9 27.0917 0.94 25.8229 0.61C23.5695 0 14.4848 0 14.4848 0C14.4848 0 5.42037 0 3.17711 0.61C1.9286 0.94 0.954148 1.9 0.59888 3.14C0 5.34 0 9.98 0 9.98C0 9.98 0 14.65 0.59888 16.87C0.954148 18.08 1.9286 19.06 3.17711 19.4C5.42037 20 14.4848 20 14.4848 20Z" fill="#FF0033"></path>
    <path d="M19 10L11.5 5.75V14.25L19 10Z" fill="white"></path>
  </g>
  <g id="youtube-paths_yt24">
    <path d="M37.1384 18.8999V13.4399L40.6084 2.09994H38.0184L36.6984 7.24994C36.3984 8.42994 36.1284 9.65994 35.9284 10.7999H35.7684C35.6584 9.79994 35.3384 8.48994 35.0184 7.22994L33.7384 2.09994H31.1484L34.5684 13.4399V18.8999H37.1384Z"></path>
    <path d="M44.1003 6.29994C41.0703 6.29994 40.0303 8.04994 40.0303 11.8199V13.6099C40.0303 16.9899 40.6803 19.1099 44.0403 19.1099C47.3503 19.1099 48.0603 17.0899 48.0603 13.6099V11.8199C48.0603 8.44994 47.3803 6.29994 44.1003 6.29994ZM45.3903 14.7199C45.3903 16.3599 45.1003 17.3899 44.0503 17.3899C43.0203 17.3899 42.7303 16.3499 42.7303 14.7199V10.6799C42.7303 9.27994 42.9303 8.02994 44.0503 8.02994C45.2303 8.02994 45.3903 9.34994 45.3903 10.6799V14.7199Z"></path>
    <path d="M52.2713 19.0899C53.7313 19.0899 54.6413 18.4799 55.3913 17.3799H55.5013L55.6113 18.8999H57.6012V6.53994H54.9613V16.4699C54.6812 16.9599 54.0312 17.3199 53.4212 17.3199C52.6512 17.3199 52.4113 16.7099 52.4113 15.6899V6.53994H49.7812V15.8099C49.7812 17.8199 50.3613 19.0899 52.2713 19.0899Z"></path>
    <path d="M62.8261 18.8999V4.14994H65.8661V2.09994H57.1761V4.14994H60.2161V18.8999H62.8261Z"></path>
    <path d="M67.8728 19.0899C69.3328 19.0899 70.2428 18.4799 70.9928 17.3799H71.1028L71.2128 18.8999H73.2028V6.53994H70.5628V16.4699C70.2828 16.9599 69.6328 17.3199 69.0228 17.3199C68.2528 17.3199 68.0128 16.7099 68.0128 15.6899V6.53994H65.3828V15.8099C65.3828 17.8199 65.9628 19.0899 67.8728 19.0899Z"></path>
    <path d="M80.6744 6.26994C79.3944 6.26994 78.4744 6.82994 77.8644 7.73994H77.7344C77.8144 6.53994 77.8744 5.51994 77.8744 4.70994V1.43994H75.3244L75.3144 12.1799L75.3244 18.8999H77.5444L77.7344 17.6999H77.8044C78.3944 18.5099 79.3044 19.0199 80.5144 19.0199C82.5244 19.0199 83.3844 17.2899 83.3844 13.6099V11.6999C83.3844 8.25994 82.9944 6.26994 80.6744 6.26994ZM80.7644 13.6099C80.7644 15.9099 80.4244 17.2799 79.3544 17.2799C78.8544 17.2799 78.1644 17.0399 77.8544 16.5899V9.23994C78.1244 8.53994 78.7244 8.02994 79.3944 8.02994C80.4744 8.02994 80.7644 9.33994 80.7644 11.7299V13.6099Z"></path>
    <path d="M92.6517 11.4999C92.6517 8.51994 92.3517 6.30994 88.9217 6.30994C85.6917 6.30994 84.9717 8.45994 84.9717 11.6199V13.7899C84.9717 16.8699 85.6317 19.1099 88.8417 19.1099C91.3817 19.1099 92.6917 17.8399 92.5417 15.3799L90.2917 15.2599C90.2617 16.7799 89.9117 17.3999 88.9017 17.3999C87.6317 17.3999 87.5717 16.1899 87.5717 14.3899V13.5499H92.6517V11.4999ZM88.8617 7.96994C90.0817 7.96994 90.1717 9.11994 90.1717 11.0699V12.0799H87.5717V11.0699C87.5717 9.13994 87.6517 7.96994 88.8617 7.96994Z"></path>
  </g>
</svg></div></span></yt-icon></ytd-logo>
  </div>
  <ytd-yoodle-renderer class="style-scope ytd-topbar-logo-renderer" hide-lottie="" hidden=""><!--css-build:shady--><!--css-build:shady--><picture class="style-scope ytd-yoodle-renderer">
  <source type="image/webp" class="style-scope ytd-yoodle-renderer" srcset="">
  <img class="style-scope ytd-yoodle-renderer" src="">
</picture>
<ytd-logo class="style-scope ytd-yoodle-renderer" enable-refresh-ringo2-web="" hidden=""><!--css-build:shady--><!--css-build:shady--><yt-icon id="logo-icon" class="style-scope ytd-logo"><!--css-build:shady--><!--css-build:shady--><span class="yt-icon-shape style-scope yt-icon yt-spec-icon-shape"><div style="width: 100%; height: 100%; display: block; fill: currentcolor;"><svg xmlns="http://www.w3.org/2000/svg" id="yt-ringo2-svg_yt9" width="93" height="20" viewBox="0 0 93 20" focusable="false" aria-hidden="true" style="pointer-events: none; display: inherit; width: 100%; height: 100%;">
  <g>
    <path d="M14.4848 20C14.4848 20 23.5695 20 25.8229 19.4C27.0917 19.06 28.0459 18.08 28.3808 16.87C29 14.65 29 9.98 29 9.98C29 9.98 29 5.34 28.3808 3.14C28.0459 1.9 27.0917 0.94 25.8229 0.61C23.5695 0 14.4848 0 14.4848 0C14.4848 0 5.42037 0 3.17711 0.61C1.9286 0.94 0.954148 1.9 0.59888 3.14C0 5.34 0 9.98 0 9.98C0 9.98 0 14.65 0.59888 16.87C0.954148 18.08 1.9286 19.06 3.17711 19.4C5.42037 20 14.4848 20 14.4848 20Z" fill="#FF0033"></path>
    <path d="M19 10L11.5 5.75V14.25L19 10Z" fill="white"></path>
  </g>
  <g id="youtube-paths_yt9">
    <path d="M37.1384 18.8999V13.4399L40.6084 2.09994H38.0184L36.6984 7.24994C36.3984 8.42994 36.1284 9.65994 35.9284 10.7999H35.7684C35.6584 9.79994 35.3384 8.48994 35.0184 7.22994L33.7384 2.09994H31.1484L34.5684 13.4399V18.8999H37.1384Z"></path>
    <path d="M44.1003 6.29994C41.0703 6.29994 40.0303 8.04994 40.0303 11.8199V13.6099C40.0303 16.9899 40.6803 19.1099 44.0403 19.1099C47.3503 19.1099 48.0603 17.0899 48.0603 13.6099V11.8199C48.0603 8.44994 47.3803 6.29994 44.1003 6.29994ZM45.3903 14.7199C45.3903 16.3599 45.1003 17.3899 44.0503 17.3899C43.0203 17.3899 42.7303 16.3499 42.7303 14.7199V10.6799C42.7303 9.27994 42.9303 8.02994 44.0503 8.02994C45.2303 8.02994 45.3903 9.34994 45.3903 10.6799V14.7199Z"></path>
    <path d="M52.2713 19.0899C53.7313 19.0899 54.6413 18.4799 55.3913 17.3799H55.5013L55.6113 18.8999H57.6012V6.53994H54.9613V16.4699C54.6812 16.9599 54.0312 17.3199 53.4212 17.3199C52.6512 17.3199 52.4113 16.7099 52.4113 15.6899V6.53994H49.7812V15.8099C49.7812 17.8199 50.3613 19.0899 52.2713 19.0899Z"></path>
    <path d="M62.8261 18.8999V4.14994H65.8661V2.09994H57.1761V4.14994H60.2161V18.8999H62.8261Z"></path>
    <path d="M67.8728 19.0899C69.3328 19.0899 70.2428 18.4799 70.9928 17.3799H71.1028L71.2128 18.8999H73.2028V6.53994H70.5628V16.4699C70.2828 16.9599 69.6328 17.3199 69.0228 17.3199C68.2528 17.3199 68.0128 16.7099 68.0128 15.6899V6.53994H65.3828V15.8099C65.3828 17.8199 65.9628 19.0899 67.8728 19.0899Z"></path>
    <path d="M80.6744 6.26994C79.3944 6.26994 78.4744 6.82994 77.8644 7.73994H77.7344C77.8144 6.53994 77.8744 5.51994 77.8744 4.70994V1.43994H75.3244L75.3144 12.1799L75.3244 18.8999H77.5444L77.7344 17.6999H77.8044C78.3944 18.5099 79.3044 19.0199 80.5144 19.0199C82.5244 19.0199 83.3844 17.2899 83.3844 13.6099V11.6999C83.3844 8.25994 82.9944 6.26994 80.6744 6.26994ZM80.7644 13.6099C80.7644 15.9099 80.4244 17.2799 79.3544 17.2799C78.8544 17.2799 78.1644 17.0399 77.8544 16.5899V9.23994C78.1244 8.53994 78.7244 8.02994 79.3944 8.02994C80.4744 8.02994 80.7644 9.33994 80.7644 11.7299V13.6099Z"></path>
    <path d="M92.6517 11.4999C92.6517 8.51994 92.3517 6.30994 88.9217 6.30994C85.6917 6.30994 84.9717 8.45994 84.9717 11.6199V13.7899C84.9717 16.8699 85.6317 19.1099 88.8417 19.1099C91.3817 19.1099 92.6917 17.8399 92.5417 15.3799L90.2917 15.2599C90.2617 16.7799 89.9117 17.3999 88.9017 17.3999C87.6317 17.3999 87.5717 16.1899 87.5717 14.3899V13.5499H92.6517V11.4999ZM88.8617 7.96994C90.0817 7.96994 90.1717 9.11994 90.1717 11.0699V12.0799H87.5717V11.0699C87.5717 9.13994 87.6517 7.96994 88.8617 7.96994Z"></path>
  </g>
</svg></div></span></yt-icon></ytd-logo>
<ytd-lottie-player></ytd-lottie-player></ytd-yoodle-renderer>
