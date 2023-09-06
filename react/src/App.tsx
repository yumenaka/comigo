import axios from "axios";
import React, { useEffect, useReducer } from "react";
import { useTranslation } from "react-i18next";
// TypeScript環境でReact Hook Formのフォーム作成の基礎を学ぶ  https://reffect.co.jp/react/react-hook-form-ts/
// import { useForm } from "react-hook-form";
import Contained from "./components/Contained";
import Config from "./types/Config";
import InputWithLabel from "./components/InputWithLabel";
import StringArrayInput from "./components/StringArrayInput";
import BoolSwitch from "./components/BoolSwitch";
import { useState } from "react";
import { configReducer, defaultConfig } from "./reducers/configReducer";

function App() {
  const baseURL = "/api";
  const { t } = useTranslation();
  const [show, setShow] = useState("bookstore")
  // nullは型に含めず、useStateの初期値が決まらないという場合には、型アサーションで逃げる手もあります。 https://qiita.com/FumioNonaka/items/4361d1cdf34ffb5a5338
  const [config, dispatch] = useReducer(configReducer, defaultConfig);

  const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    axios.post("/api/update_config", config).then((response) => {
      console.log("Data sent successfully");
      console.info(response.data);//axios默认解析Json，所以 response.data 就是解析后的object
    })
      .catch((error) => {
        console.error("Error sending data:", error);
      });
  };
  // useEffect 用于在函数组件中执行副作用操作，例如数据获取、订阅、手动修改DOM等。
  // 通过传递第二个参数，你可以告诉 React 仅在某些值改变的时候才执行 effect。
  // 传递空数组([])作为第二个参数，effect 内部的 props 和 state 就会一直持有其初始值。也就是只在渲染的时候执行一次。
  useEffect(() => {
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

  //  监听事件 实现数据的动态录入
  const setStringValue = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    console.log(typeof value);
    dispatch({
      type: 'stringConfig',
      name: name,
      value: value,
      config: config
    });
  };

  const setNumberValue = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    console.log(typeof value);
    dispatch({
      type: 'numberConfig',
      name: name,
      value: value,
      config: config
    });
  };

  const setBoolValue = (value: boolean, valueName: string) => {
    console.log("setBoolValue" + valueName, value);
    dispatch({
      type: 'boolConfig',
      name: valueName,
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
    <>
      <Contained show={show} setShow={setShow} />
      <form
        onSubmit={onSubmit}
        className="card w-full flex flex-col justify-center items-center"
      >
        {show === "bookstore" &&
          <>
            <StringArrayInput
              label={t("StoresPath")}
              fieldDescription="书库文件夹，支持绝对目录与相对目录。相对目录以当前执行目录为基准"
              name={"StoresPath"}
              value={config.StoresPath}
              setStringArray={setStringArray}
            ></StringArrayInput>

            <InputWithLabel
              label={t("MaxScanDepth")}
              fieldDescription="最大扫描深度。超过深度的文件不会被扫描。当前执行目录为基准。"
              name={"MaxScanDepth"}
              type={"number"}
              value={config.MaxScanDepth}
              onChange={setNumberValue}
              placeholder={"MaxScanDepth"}
            ></InputWithLabel>

            <BoolSwitch
              label={t("OpenBrowser")}
              fieldDescription="扫描完成后，是否同时打开浏览器。windows默认true，其他平台默认false。"
              name={"OpenBrowser"}
              boolValue={config.OpenBrowser}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <BoolSwitch
              name={"ClearCacheExit"}
              label={t("ClearCacheExit")}
              fieldDescription="退出程序的时候，清理web图片缓存。"
              boolValue={config.ClearCacheExit}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <InputWithLabel
              label={t("CachePath")}
              fieldDescription="本地图片缓存位置，默认系统临时文件夹。"
              name={"CachePath"}
              type={"text"}
              value={config.CachePath}
              onChange={setStringValue}
              placeholder={"CachePath"}
            ></InputWithLabel>

            <BoolSwitch
              name={"EnableUpload"}
              label={t("EnableUpload")}
              fieldDescription="启用上传功能。"
              boolValue={config.EnableUpload}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            {config.EnableUpload && <InputWithLabel
              label={t("UploadPath")}
              fieldDescription="上传文件存储位置，默认在当前执行目录下创建 upload 文件夹。"
              name={"UploadPath"}
              type={"text"}
              value={config.UploadPath}
              onChange={setStringValue}
              placeholder={"UploadPath"}
            />}

            <InputWithLabel
              label={t("MinImageNum")}
              fieldDescription="压缩包或文件夹内至少有几张图片，才算作书籍。"
              name={"MinImageNum"}
              type={"number"}
              value={config.MinImageNum}
              onChange={setNumberValue}
              placeholder={"MinImageNum"}
            ></InputWithLabel>

            <InputWithLabel
              label={t("TimeoutLimitForScan")}
              fieldDescription="扫描文件时，超过几秒钟，就放弃扫描这个文件，避免卡在过大文件上。"
              name={"TimeoutLimitForScan"}
              type={"numbers"}
              value={config.TimeoutLimitForScan}
              onChange={setNumberValue}
              placeholder={"TimeoutLimitForScan"}
            />

            <InputWithLabel
              label={t("ZipFileTextEncoding")}
              fieldDescription="非utf-8编码的ZIP文件，尝试用什么编码解析，默认GBK。"
              name={"ZipFileTextEncoding"}
              type={"text"}
              value={config.ZipFileTextEncoding}
              onChange={setStringValue}
              placeholder={"ZipFileTextEncoding"}
            />

            <BoolSwitch
              name={"GenerateMetaData"}
              label={t("GenerateMetaData")}
              fieldDescription="生成书籍元数据。"
              boolValue={config.GenerateMetaData}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <StringArrayInput
              label={t("ExcludePath")}
              fieldDescription=""
              name={"ExcludePath"}
              value={config.ExcludePath}
              setStringArray={setStringArray}
            ></StringArrayInput>

            <StringArrayInput
              label={t("SupportMediaType")}
              fieldDescription=""
              name={"SupportMediaType"}
              value={config.SupportMediaType}
              setStringArray={setStringArray}
            ></StringArrayInput>

            <StringArrayInput
              label={t("SupportFileType")}
              fieldDescription=""
              name={"SupportFileType"}
              value={config.SupportFileType}
              setStringArray={setStringArray}
            ></StringArrayInput>
          </>
        }

        {show === "internet" &&
          <>
            <InputWithLabel
              label={t("Port")}
              fieldDescription=""
              name={"Port"}
              type={"number"}
              value={config.Port}
              onChange={setNumberValue}
              placeholder={"Port"}
            ></InputWithLabel>

            <InputWithLabel
              label={t("Host")}
              fieldDescription=""
              name={"Host"}
              type={"text"}
              value={config.Host}
              onChange={setStringValue}
              placeholder={"Host"}
            ></InputWithLabel>

            <BoolSwitch
              name={"DisableLAN"}
              label={t("DisableLAN")}
              fieldDescription=""
              boolValue={config.DisableLAN}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <BoolSwitch
              name={"EnableLogin"}
              label={t("EnableLogin")}
              fieldDescription=""
              boolValue={config.EnableLogin}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            {config.EnableLogin && <InputWithLabel
              label={t("Username")}
              fieldDescription=""
              name={"Username"}
              type={"text"}
              value={config.Username}
              onChange={setStringValue}
              placeholder={"Username"}
            ></InputWithLabel>}

            {config.EnableLogin && <InputWithLabel
              label={t("Password")}
              fieldDescription=""
              name={"Password"}
              type={"text"}
              value={config.Password}
              onChange={setStringValue}
              placeholder={"Password"}
            ></InputWithLabel>}

            <InputWithLabel
              label={t("Timeout")}
              fieldDescription=""
              name={"Timeout"}
              type={"number"}
              value={config.Timeout}
              onChange={setNumberValue}
              placeholder={"Timeout"}
            ></InputWithLabel>

            <BoolSwitch
              name={"EnableTLS"}
              label={t("EnableTLS")}
              fieldDescription=""
              boolValue={config.EnableTLS}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            {config.EnableTLS && <InputWithLabel
              label={t("CertFile")}
              fieldDescription=""
              name={"CertFile"}
              type={"text"}
              value={config.CertFile}
              onChange={setStringValue}
              placeholder={"CertFile"}
            ></InputWithLabel>}

            {config.EnableTLS && <InputWithLabel
              label={t("KeyFile")}
              fieldDescription=""
              name={"KeyFile"}
              type={"text"}
              value={config.KeyFile}
              onChange={setStringValue}
              placeholder={"KeyFile"}
            ></InputWithLabel>}

          </>
        }

        {show === "other" &&
          <>
            <BoolSwitch
              name={"EnableDatabase"}
              label={t("EnableDatabase")}
              fieldDescription=""
              boolValue={config.EnableDatabase}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            {config.EnableDatabase && <BoolSwitch
              name={"ClearDatabaseWhenExit"}
              label={t("ClearDatabaseWhenExit")}
              fieldDescription=""
              boolValue={config.ClearDatabaseWhenExit}
              setBoolValue={setBoolValue}
            ></BoolSwitch>}

            <BoolSwitch
              name={"Debug"}
              label={t("Debug")}
              fieldDescription=""
              boolValue={config.Debug}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <BoolSwitch
              name={"LogToFile"}
              label={t("LogToFile")}
              fieldDescription=""
              boolValue={config.LogToFile}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            {config.LogToFile &&
              <InputWithLabel
                label={t("LogFilePath")}
                fieldDescription=""
                name={"LogFilePath"}
                type={"text"}
                value={config.LogFilePath}
                onChange={setStringValue}
                placeholder={"LogFilePath"}
              />}

            {config.LogToFile &&
              <InputWithLabel
                label={t("LogFileName")}
                fieldDescription=""
                name={"LogFileName"}
                type={"text"}
                value={config.LogFileName}
                onChange={setStringValue}
                placeholder={"LogFileName"}
              />
            }

            <BoolSwitch
              name={"EnableFrpcServer"}
              label={t("EnableFrpcServer")}
              fieldDescription=""
              boolValue={config.EnableFrpcServer}
              setBoolValue={setBoolValue}
            ></BoolSwitch>
          </>
        }

      </form>
      <button
          type="submit"
          className="m-2 inline-block rounded-2xl bg-indigo-600 px-8 py-3 text-sm font-medium text-white transition hover:scale-110 hover:shadow-xl focus:outline-none focus:ring active:bg-indigo-500"
        >
          保存设置
        </button>
    </>
  );
}

export default App;