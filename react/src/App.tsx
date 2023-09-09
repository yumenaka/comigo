import axios from "axios";
import React, { useEffect, useReducer } from "react";
import { useTranslation } from "react-i18next";
// TypeScript環境でReact Hook Formのフォーム作成の基礎を学ぶ  https://reffect.co.jp/react/react-hook-form-ts/
// import { useForm } from "react-hook-form";
import Contained from "./components/Contained";
import Config from "./types/Config";
import NormalInput from "./components/NormalInput";
import StringArrayInput from "./components/StringArrayInput";
import BoolSwitch from "./components/BoolInput";
import { useState } from "react";
import { configReducer, defaultConfig } from "./reducers/configReducer";

function App() {
  const baseURL = "/api";
  const { t } = useTranslation();
  const [show, setShow] = useState("bookstore")
  const [config, dispatch] = useReducer(configReducer, defaultConfig);
  const [BackgroundColor, setBackgroundColor] = useState("#e0d9cd")
  const [InterfaceColor, setInterfaceColor] = useState("#F5F5E4")

  // useEffect 用于在函数组件中执行副作用操作，例如数据获取、订阅、手动修改DOM等。
  // 通过传递第二个参数，你可以告诉 React 仅在某些值改变的时候才执行 effect。
  // 传递空数组([])作为第二个参数，effect 内部的 props 和 state 就会一直持有其初始值。也就是只在渲染的时候执行一次。
  useEffect(() => {
    // 当前颜色
    const tempBackgroundColor = localStorage.getItem("BackgroundColor");
    if (tempBackgroundColor !== null) {
      console.log("tempBackgroundColor", tempBackgroundColor)
      setBackgroundColor(tempBackgroundColor)
    }
    // 当前颜色
    const tempInterfaceColor = localStorage.getItem("InterfaceColor");
    if (tempInterfaceColor !== null) {
      console.log("tempInterfaceColor", tempInterfaceColor)
      setInterfaceColor(tempInterfaceColor)
    }
    // 从后端获取配置文件
    axios
      .get<Config>(`${baseURL}/config.json`)
      .then((response) => {
        dispatch({
          type: 'downloadConfig',
          name: "",
          value: "",
          config: response.data
        });
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);

  //配置文件修改后，保存到后端
  const setStringValue = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    dispatch({
      type: 'stringConfig',
      name: name,
      value: value,
      config: config
    });
  };

  const setNumberValue = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    dispatch({
      type: 'numberConfig',
      name: name,
      value: value,
      config: config
    });
  };

  const setBoolValue = (name: string, value: boolean) => {
    console.log("setBoolValue " + name, value);
    dispatch({
      type: 'boolConfig',
      name: name,
      value: value,
      config: config
    });
  };

  const setStringArray = (valueName: string, value: string[]) => {
    dispatch({
      type: 'boolConfig',
      name: valueName,
      value: value,
      config: config
    });
  };

  return (
    <div
      style={{
        backgroundColor: BackgroundColor, // 绑定样式
      }}
      className={`w-full h-full min-h-screen flex flex-col justify-start items-center`} >
      {/* 顶部标题 */}
      <div className="w-full h-16 mb-1 rounded shadow flex flex-row justify-center items-center" style={{
        backgroundColor: InterfaceColor, // 绑定样式
      }}>
        <Contained show={show} setShow={setShow} BackgroundColor={BackgroundColor} />
      </div>

      <div
        className={`main-area w-3/5 min-w-[24rem] flex flex-col justify-center items-center`}
      >
        {show === "bookstore" &&
          <>
            <BoolSwitch
              label={t("OpenBrowser")}
              fieldDescription="扫描完成后，是否同时打开浏览器。windows默认true，其他平台默认false。"
              name={"OpenBrowser"}
              boolValue={config.OpenBrowser}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <BoolSwitch
              label={t("EnableUpload")}
              name={"EnableUpload"}
              fieldDescription="启用上传功能。"
              boolValue={config.EnableUpload}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <StringArrayInput
              label={t("StoresPath")}
              fieldDescription="书库文件夹，支持绝对目录与相对目录。相对目录以当前执行目录为基准"
              name={"StoresPath"}
              value={config.StoresPath}
              InterfaceColor={InterfaceColor}
              setStringArray={setStringArray}
            ></StringArrayInput>

            <NormalInput
              label={t("MaxScanDepth")}
              fieldDescription="最大扫描深度。超过深度的文件不会被扫描。当前执行目录为基准。"
              name={"MaxScanDepth"}
              type={"number"}
              value={config.MaxScanDepth}
              InterfaceColor={InterfaceColor}
              onChange={setNumberValue}
              placeholder={"MaxScanDepth"}
            ></NormalInput>

            {config.EnableUpload && <NormalInput
              label={t("UploadPath")}
              fieldDescription="自定义上传文件存储位置，默认在当前执行目录下创建 upload 文件夹。"
              name={"UploadPath"}
              type={"text"}
              value={config.UploadPath}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"UploadPath"}
            />}

            <NormalInput
              label={t("MinImageNum")}
              fieldDescription="压缩包或文件夹内至少有几张图片，才算作书籍。"
              name={"MinImageNum"}
              type={"number"}
              value={config.MinImageNum}
              InterfaceColor={InterfaceColor}
              onChange={setNumberValue}
              placeholder={"MinImageNum"}
            ></NormalInput>

            <NormalInput
              label={t("TimeoutLimitForScan")}
              fieldDescription="扫描文件时，超过几秒钟，就放弃扫描这个文件，避免卡在过大文件上。"
              name={"TimeoutLimitForScan"}
              type={"number"}
              value={config.TimeoutLimitForScan}
              InterfaceColor={InterfaceColor}
              onChange={setNumberValue}
              placeholder={"TimeoutLimitForScan"}
            />

            <NormalInput
              label={t("ZipFileTextEncoding")}
              fieldDescription="非utf-8编码ZIP文件，尝试用什么编码解析。默认GBK。"
              name={"ZipFileTextEncoding"}
              type={"text"}
              value={config.ZipFileTextEncoding}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"ZipFileTextEncoding"}
            />

            <StringArrayInput
              label={t("ExcludePath")}
              fieldDescription="扫描书籍的时候，需要排除的文件或文件夹的名字"
              name={"ExcludePath"}
              value={config.ExcludePath}
              InterfaceColor={InterfaceColor}
              setStringArray={setStringArray}
            ></StringArrayInput>

            <StringArrayInput
              label={t("SupportMediaType")}
              fieldDescription="扫描压缩包时，用于统计图片数量的图片文件后缀"
              name={"SupportMediaType"}
              value={config.SupportMediaType}
              InterfaceColor={InterfaceColor}
              setStringArray={setStringArray}
            ></StringArrayInput>

            <StringArrayInput
              label={t("SupportFileType")}
              fieldDescription="扫描文件时，用于决定跳过，还是算作书籍处理的文件后缀"
              name={"SupportFileType"}
              value={config.SupportFileType}
              InterfaceColor={InterfaceColor}
              setStringArray={setStringArray}
            ></StringArrayInput>
          </>
        }

        {show === "internet" &&
          <>
            <NormalInput
              label={t("Port")}
              fieldDescription="网页服务端口，此项配置不支持热重载"
              name={"Port"}
              type={"number"}
              value={config.Port}
              InterfaceColor={InterfaceColor}
              onChange={setNumberValue}
              placeholder={"Port"}
            ></NormalInput>

            <NormalInput
              label={t("Host")}
              fieldDescription="自定义二维码显示的主机名。默认为网卡IP。"
              name={"Host"}
              type={"text"}
              value={config.Host}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"Host"}
            ></NormalInput>

            <BoolSwitch
              name={"DisableLAN"}
              label={t("DisableLAN")}
              fieldDescription="只在本机提供阅读服务，不对外共享，此项配置不支持热重载"
              boolValue={config.DisableLAN}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <BoolSwitch
              name={"EnableLogin"}
              label={t("EnableLogin")}
              fieldDescription="是否启用登录。默认不需要登陆。此项配置不支持热重载。"
              boolValue={config.EnableLogin}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            {config.EnableLogin && <NormalInput
              label={t("Username")}
              fieldDescription="启用登陆后，登录界面需要的用户名。"
              name={"Username"}
              type={"text"}
              value={config.Username}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"Username"}
            ></NormalInput>}

            {config.EnableLogin && <NormalInput
              label={t("Password")}
              fieldDescription="启用登陆后，登录界面需要的密码。"
              name={"Password"}
              type={"text"}
              value={config.Password}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"Password"}
            ></NormalInput>}

            <NormalInput
              label={t("Timeout")}
              fieldDescription="启用登陆后，cookie过期的时间。单位为分钟。默认180分钟后过期。"
              name={"Timeout"}
              type={"number"}
              value={config.Timeout}
              InterfaceColor={InterfaceColor}
              onChange={setNumberValue}
              placeholder={"Timeout"}
            ></NormalInput>

            <BoolSwitch
              name={"EnableTLS"}
              label={t("EnableTLS")}
              fieldDescription="是否启用HTTPS协议。需要设置证书于key文件。"
              boolValue={config.EnableTLS}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            {config.EnableTLS && <NormalInput
              label={t("CertFile")}
              fieldDescription='TLS/SSL 证书文件路径 (default: 、"~/.config/.comigo/cert.crt")'
              name={"CertFile"}
              type={"text"}
              value={config.CertFile}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"CertFile"}
            ></NormalInput>}

            {config.EnableTLS && <NormalInput
              label={t("KeyFile")}
              fieldDescription='TLS/SSL key文件路径 (default: "~/.config/.comigo/key.key")'
              name={"KeyFile"}
              type={"text"}
              value={config.KeyFile}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"KeyFile"}
            ></NormalInput>}

          </>
        }

        {show === "other" &&
          <>
            <BoolSwitch
              name={"EnableDatabase"}
              label={t("EnableDatabase")}
              fieldDescription="启用本地数据库，保存扫描到的书籍数据。此项配置不支持热重载。"
              boolValue={config.EnableDatabase}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            {config.EnableDatabase && <BoolSwitch
              name={"ClearDatabaseWhenExit"}
              label={t("ClearDatabaseWhenExit")}
              fieldDescription="启用本地数据库时，扫描完成后，清除不存在的书籍。"
              boolValue={config.ClearDatabaseWhenExit}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolSwitch>}

            <BoolSwitch
              name={"Debug"}
              label={t("Debug")}
              fieldDescription=""
              boolValue={config.Debug}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <BoolSwitch
              name={"LogToFile"}
              label={t("LogToFile")}
              fieldDescription="是否保存程序Log到本地文件。默认不保存。"
              boolValue={config.LogToFile}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            {config.LogToFile &&
              <NormalInput
                label={t("LogFilePath")}
                fieldDescription="Log文件的保存位置"
                name={"LogFilePath"}
                type={"text"}
                value={config.LogFilePath}
                InterfaceColor={InterfaceColor}
                onChange={setStringValue}
                placeholder={"LogFilePath"}
              />}

            {config.LogToFile &&
              <NormalInput
                label={t("LogFileName")}
                fieldDescription="Log文件名"
                name={"LogFileName"}
                type={"text"}
                value={config.LogFileName}
                InterfaceColor={InterfaceColor}
                onChange={setStringValue}
                placeholder={"LogFileName"}
              />
            }

            <BoolSwitch
              name={"GenerateMetaData"}
              label={t("GenerateMetaData")}
              fieldDescription="生成书籍元数据。当前未生效。"
              boolValue={config.GenerateMetaData}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <BoolSwitch
              name={"ClearCacheExit"}
              label={t("ClearCacheExit")}
              fieldDescription="退出程序的时候，清理web图片缓存。"
              boolValue={config.ClearCacheExit}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <NormalInput
              label={t("CachePath")}
              fieldDescription="本地图片缓存位置，默认系统临时文件夹。"
              name={"CachePath"}
              type={"text"}
              value={config.CachePath}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"CachePath"}
            ></NormalInput>
          </>
        }
      </div>
      {/* 返回主页的按钮 */}
      <a
        className="fixed top-2 left-2 inline-block rounded-full border bg-white border-indigo-600 p-3 text-indigo-600 hover:bg-indigo-600 hover:text-white focus:outline-none focus:ring active:bg-indigo-500"
        href="/"
      >
        <span className="sr-only"> Download </span>
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-6 h-6">
          <path strokeLinecap="round" strokeLinejoin="round" d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18" />
        </svg>
      </a>

      {/* 底部提示 */}
      <div className="w-full mt-auto flex flex-col justify-center items-center text-gray-900 h-12 py-4 space-x-2 text-base content-center" style={{
        backgroundColor: InterfaceColor, // 绑定样式
      }}>
        <a href="https://github.com/yumenaka/comi/releases" className="text-blue-700 hover:underline font-semibold"> Power by Comigo</a>
      </div>
    </div>
  );
}

export default App;