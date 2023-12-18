//导入各种库
import axios from "axios";
import React, { useReducer, useState } from "react";
import { useTranslation } from "react-i18next";
// react-use中文文档： https://github.com/zenghongtu/react-use-chinese/blob/master/README.md  英文文档：https://streamich.github.io/react-use/?path=%2Fstory%2Flifecycle-useeffectonce--docs
import { useEffectOnce } from 'react-use';
import Cookies from 'js-cookie';

//自定义组件
import Header from "./components/Header";
import NormalConfig from "./components/NormalInput";
import ArrayConfig from "./components/ArrayConfig";
import BoolConfig from "./components/BoolConfig";
import DialogModal from "./components/DialogModal";

//使用useReducer管理的远程数据
import Config from "./types/Config";
import { configReducer, defaultConfig } from "./reducers/configReducer";

import ConfigManager from "./components/ConfigManager";
// import { useForm } from "react-hook-form"; //sample：https://reffect.co.jp/react/react-hook-form-ts/  （TypeScript環境でReact Hook Formのフォーム作成の基礎を学ぶ）
import DialogStatus from './types/DialogStatus';

function App() {
  const baseURL = "/api";
  const { t, i18n } = useTranslation();
  const [headerGroup, setHeaderGroup] = useState("bookstore")
  const [BackgroundColor, setBackgroundColor] = useState("#e0d9cd")
  const [InterfaceColor, setInterfaceColor] = useState("#F5F5E4")
  //useReducer 和 useState 非常相似，但是它可以让你把状态更新逻辑从事件处理函数中移动到组件外部:https://zh-hans.react.dev/reference/react/useReducer
  //在用法上，它接收一个reducer函数作为第一个参数，第二个参数是初始化的state。
  //useReducer最终返回一个存储有当前状态值的数组和一个dispatch函数，该dispatch函数执行触发action，带来状态的变化。
  const [config, config_dispatch] = useReducer(configReducer, defaultConfig)



  // 只执行一次的useEffect，来自'react-use'库。
  useEffectOnce(() => {
    // 当前语言 jp zh en
    // document.cookie="userLanguageSetting=jp"  //手动设置cookie
    const lang = Cookies.get("userLanguageSetting");
    if (lang) {
      i18n.changeLanguage(lang).then(() => {
        console.log("i18n.changeLanguage", lang);
      }).catch((err) => {
        console.log("i18n.changeLanguage", err)
      });
    }
    // 主题色1
    const tempBackgroundColor = localStorage.getItem("BackgroundColor");
    if (tempBackgroundColor !== null) {
      console.log("tempBackgroundColor", tempBackgroundColor)
      setBackgroundColor(tempBackgroundColor)
    }
    // 主题色2
    const tempInterfaceColor = localStorage.getItem("InterfaceColor");
    if (tempInterfaceColor !== null) {
      console.log("tempInterfaceColor", tempInterfaceColor)
      setInterfaceColor(tempInterfaceColor)
    }

    // 获取comigo配置
    axios
      .get<Config>(`${baseURL}/config`)
      .then((response) => {
        config_dispatch({
          type: 'init',
          name: "",
          value: "",
          config: response.data
        });
      })
      .catch((error) => {
        console.error(error);
      });

  });

  //配置文件修改后，保存到后端的各种函数
  const setStringValueFunc = (name: string, value: string) => {
    console.log("setStringValue " + name, value);
    config_dispatch({
      type: 'string',
      name: name,
      value: value,
      config: config
    });
  };
  //字符串类型的配置文件修改后，保存到后端的函数
  const setStringValue = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setStringValueFunc(name, value);
  };
  //数字类型的配置文件修改后，保存到后端的函数
  const setNumberValue = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    config_dispatch({
      type: 'number',
      name: name,
      value: value,
      config: config
    });
  };
  //布尔类型的配置文件修改后，保存到后端的函数
  const setBoolValue = (name: string, value: boolean) => {
    console.log("setBoolValue " + name, value);
    config_dispatch({
      type: 'boolean',
      name: name,
      value: value,
      config: config
    });
  };
  //字符串数组类型的配置文件修改后，保存到后端的函数
  const setStringArray = (valueName: string, value: string[]) => {
    config_dispatch({
      type: 'boolean',
      name: valueName,
      value: value,
      config: config
    });
  };

  //弹窗提示
  const [dialogStatus, setDialogStatus] = useState({
    isOpen: false,
    title: "Test",
    OK: "OK",
    description: "bla bla bla"
  });
  //开启弹窗的函数，给需要弹窗提示组件用
  function openDialog(title: string, OK: string, description: string) {
    setDialogStatus({
      isOpen: true,
      title: title,
      OK: OK,
      description: description
    }as DialogStatus);
  }
  //关闭弹窗的函数，弹窗组件本身用
  function closeDialog() {
    setDialogStatus({
      isOpen: false,
      title: dialogStatus.title,
      OK: dialogStatus.OK,
      description: dialogStatus.description
    }as DialogStatus);
  }


  return (
    <div
      style={{
        backgroundColor: BackgroundColor, // 绑定样式
      }}
      className={`w-full h-full min-h-screen flex flex-col justify-start items-center`} >

      {/* 顶部标题 */}
      <Header group={headerGroup} setGroup={setHeaderGroup} InterfaceColor={InterfaceColor} />
      <DialogModal dialogStatus={dialogStatus} closeDialog={closeDialog} InterfaceColor={InterfaceColor} />

      <div
        className={`main-area w-3/5 min-w-[24rem] flex flex-col justify-center items-center`}
      >
        <ConfigManager
              label={"配置保存："}
              InterfaceColor={InterfaceColor}
            ></ConfigManager>
        {headerGroup === "bookstore" &&
          <>
   
            {/* <button className="h-15 w-full" onClick={() => setLang(lang === 'en' ? 'ja' : 'en')}>切换语言</button> */}
            <BoolConfig
              label={t("OpenBrowser")}
              fieldDescription={t("OpenBrowser_Description")}
              name={"OpenBrowser"}
              boolValue={config.OpenBrowser}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolConfig>

            <ArrayConfig
              label={t("StoresPath")}
              fieldDescription={t("StoresPath_Description")}
              name={"StoresPath"}
              value={config.StoresPath}
              InterfaceColor={InterfaceColor}
              setStringArray={setStringArray}
            ></ArrayConfig>

            <NormalConfig
              label={t("MaxScanDepth")}
              fieldDescription={t("MaxScanDepth_Description")}
              name={"MaxScanDepth"}
              type={"number"}
              value={config.MaxScanDepth}
              InterfaceColor={InterfaceColor}
              onChange={setNumberValue}
              placeholder={"MaxScanDepth"}
            ></NormalConfig>

            <NormalConfig
              label={t("MinImageNum")}
              fieldDescription={t("MinImageNum_Description")}
              name={"MinImageNum"}
              type={"number"}
              value={config.MinImageNum}
              InterfaceColor={InterfaceColor}
              onChange={setNumberValue}
              placeholder={"MinImageNum"}
            ></NormalConfig>

            <NormalConfig
              label={t("TimeoutLimitForScan")}
              fieldDescription={t("TimeoutLimitForScan_Description")}
              name={"TimeoutLimitForScan"}
              type={"number"}
              value={config.TimeoutLimitForScan}
              InterfaceColor={InterfaceColor}
              onChange={setNumberValue}
              placeholder={"TimeoutLimitForScan"}
            />

            <BoolConfig
              label={t("EnableUpload")}
              name={"EnableUpload"}
              fieldDescription={t("EnableUpload_Description")}
              boolValue={config.EnableUpload}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
              showDialog={() => {
                openDialog(t("EnableUpload"), "OK", t("EnableUpload_Description"))
              }}
            ></BoolConfig>

            {config.EnableUpload && <NormalConfig
              label={t("UploadPath")}
              fieldDescription={t("UploadPath_Description")}
              name={"UploadPath"}
              type={"text"}
              value={config.UploadPath}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"UploadPath"}
            />}

            <NormalConfig
              label={t("ZipFileTextEncoding")}
              fieldDescription={t("ZipFileTextEncoding_Description")}
              name={"ZipFileTextEncoding"}
              type={"text"}
              value={config.ZipFileTextEncoding}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"ZipFileTextEncoding"}
            />

            <ArrayConfig
              label={t("ExcludePath")}
              fieldDescription={t("ExcludePath_Description")}
              name={"ExcludePath"}
              value={config.ExcludePath}
              InterfaceColor={InterfaceColor}
              setStringArray={setStringArray}
            ></ArrayConfig>

            <ArrayConfig
              label={t("SupportMediaType")}
              fieldDescription={t("SupportMediaType_Description")}
              name={"SupportMediaType"}
              value={config.SupportMediaType}
              InterfaceColor={InterfaceColor}
              setStringArray={setStringArray}
            ></ArrayConfig>

            <ArrayConfig
              label={t("SupportFileType")}
              fieldDescription={t("SupportFileType_Description")}
              name={"SupportFileType"}
              value={config.SupportFileType}
              InterfaceColor={InterfaceColor}
              setStringArray={setStringArray}
            ></ArrayConfig>
          </>
        }

        {headerGroup === "internet" &&
          <>
            <NormalConfig
              label={t("Port")}
              fieldDescription={t("Port_Description")}
              name={"Port"}
              type={"number"}
              value={config.Port}
              InterfaceColor={InterfaceColor}
              onChange={setNumberValue}
              placeholder={"Port"}
            ></NormalConfig>

            <NormalConfig
              label={t("Host")}
              fieldDescription={t("Host_Description")}
              name={"Host"}
              type={"text"}
              value={config.Host}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"Host"}
            ></NormalConfig>

            <BoolConfig
              name={"DisableLAN"}
              label={t("DisableLAN")}
              fieldDescription={t("DisableLAN_Description")}
              boolValue={config.DisableLAN}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolConfig>

            <BoolConfig
              name={"EnableLogin"}
              label={t("EnableLogin")}
              fieldDescription={t("EnableLogin_Description")}
              boolValue={config.EnableLogin}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolConfig>

            {config.EnableLogin && <NormalConfig
              label={t("Username")}
              fieldDescription={t("Username_Description")}
              name={"Username"}
              type={"text"}
              value={config.Username}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"Username"}
            ></NormalConfig>}

            {config.EnableLogin && <NormalConfig
              label={t("Password")}
              fieldDescription={t("Password_Description")}
              name={"Password"}
              type={"text"}
              value={config.Password}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"Password"}
            ></NormalConfig>}

            <NormalConfig
              label={t("Timeout")}
              fieldDescription={t("Timeout_Description")}
              name={"Timeout"}
              type={"number"}
              value={config.Timeout}
              InterfaceColor={InterfaceColor}
              onChange={setNumberValue}
              placeholder={"Timeout"}
            ></NormalConfig>

            <BoolConfig
              name={"EnableTLS"}
              label={t("EnableTLS")}
              fieldDescription={t("EnableTLS_Description")}
              boolValue={config.EnableTLS}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolConfig>

            {config.EnableTLS && <NormalConfig
              label={t("CertFile")}
              fieldDescription={t("CertFile_Description")}
              name={"CertFile"}
              type={"text"}
              value={config.CertFile}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"CertFile"}
            ></NormalConfig>}

            {config.EnableTLS && <NormalConfig
              label={t("KeyFile")}
              fieldDescription={t("KeyFile_Description")}
              name={"KeyFile"}
              type={"text"}
              value={config.KeyFile}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"KeyFile"}
            ></NormalConfig>}

          </>
        }

        {headerGroup === "other" &&
          <>
            <div style={{ backgroundColor: InterfaceColor, }}// 绑定样式 
              className={`w-full m-1 p-2 flex flex-col shadow-md hover:shadow-2xl font-semibold rounded-md  justify-left items-left`}>
              还未完成的功能，开发与调整中。
            </div>

            <BoolConfig
              name={"EnableDatabase"}
              label={t("EnableDatabase")}
              fieldDescription={t("EnableDatabase_Description")}
              boolValue={config.EnableDatabase}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolConfig>

            {config.EnableDatabase && <BoolConfig
              name={"ClearDatabaseWhenExit"}
              label={t("ClearDatabaseWhenExit")}
              fieldDescription={t("ClearDatabaseWhenExit_Description")}
              boolValue={config.ClearDatabaseWhenExit}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolConfig>}

            <BoolConfig
              name={"Debug"}
              label={t("Debug")}
              fieldDescription={t("Debug_Description")}
              boolValue={config.Debug}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolConfig>

            <BoolConfig
              name={"LogToFile"}
              label={t("LogToFile")}
              fieldDescription={t("LogToFile_Description")}
              boolValue={config.LogToFile}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolConfig>

            {config.LogToFile &&
              <NormalConfig
                label={t("LogFilePath")}
                fieldDescription={t("LogFilePath_Description")}
                name={"LogFilePath"}
                type={"text"}
                value={config.LogFilePath}
                InterfaceColor={InterfaceColor}
                onChange={setStringValue}
                placeholder={"LogFilePath"}
              />}

            {config.LogToFile &&
              <NormalConfig
                label={t("LogFileName")}
                fieldDescription={t("LogFileName_Description")}
                name={"LogFileName"}
                type={"text"}
                value={config.LogFileName}
                InterfaceColor={InterfaceColor}
                onChange={setStringValue}
                placeholder={"LogFileName"}
              />
            }

            <BoolConfig
              name={"GenerateMetaData"}
              label={t("GenerateMetaData")}
              fieldDescription={t("GenerateMetaData_Description")}
              boolValue={config.GenerateMetaData}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolConfig>

            <BoolConfig
              name={"ClearCacheExit"}
              label={t("ClearCacheExit")}
              fieldDescription={t("ClearCacheExit_Description")}
              boolValue={config.ClearCacheExit}
              InterfaceColor={InterfaceColor}
              setBoolValue={setBoolValue}
            ></BoolConfig>

            <NormalConfig
              label={t("CachePath")}
              fieldDescription={t("CachePath_Description")}
              name={"CachePath"}
              type={"text"}
              value={config.CachePath}
              InterfaceColor={InterfaceColor}
              onChange={setStringValue}
              placeholder={"CachePath"}
            ></NormalConfig>
          </>
        }
      </div>
      {/* 返回主页的按钮 */}
      <a
        className="fixed top-2 left-2 inline-block rounded-full border  shadow-md hover:shadow-2xl bg-white border-indigo-600 p-3 text-indigo-600 hover:bg-indigo-600 hover:text-white focus:outline-none focus:ring active:bg-indigo-500"
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
        <a href="https://github.com/yumenaka/comi/releases" className="text-blue-700 hover:underline font-bold"> Power by Comigo</a>
      </div>
    </div>
  );
}

export default App;