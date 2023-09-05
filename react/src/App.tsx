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
  const [show, setShow] = useState("internet")
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
        <button
          type="submit"
          className="m-2 inline-block rounded bg-indigo-600 px-8 py-3 text-sm font-medium text-white transition hover:scale-110 hover:shadow-xl focus:outline-none focus:ring active:bg-indigo-500"
        >
          保存设置
        </button>

        {show === "bookstore" &&
          <>
            <StringArrayInput
              label={t("StoresPath")}
              name={"StoresPath"}
              value={config.StoresPath}
              setStringArray={setStringArray}
            ></StringArrayInput>

            <InputWithLabel
              label={t("MaxScanDepth")}
              name={"MaxScanDepth"}
              type={"number"}
              value={config.MaxScanDepth}
              onChange={setNumberValue}
              placeholder={"MaxScanDepth"}
            ></InputWithLabel>

            <BoolSwitch
              label={t("OpenBrowser")}
              name={"OpenBrowser"}
              boolValue={config.OpenBrowser}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <BoolSwitch
              name={"ClearCacheExit"}
              label={t("ClearCacheExit")}
              boolValue={config.ClearCacheExit}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <InputWithLabel
              label={t("CachePath")}
              name={"CachePath"}
              type={"text"}
              value={config.CachePath}
              onChange={setStringValue}
              placeholder={"CachePath"}
            ></InputWithLabel>

            <BoolSwitch
              name={"EnableUpload"}
              label={t("EnableUpload")}
              boolValue={config.EnableUpload}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            {config.EnableUpload && <InputWithLabel
              label={t("UploadPath")}
              name={"UploadPath"}
              type={"text"}
              value={config.UploadPath}
              onChange={setStringValue}
              placeholder={"UploadPath"}
            />}

            <InputWithLabel
              label={t("MinImageNum")}
              name={"MinImageNum"}
              type={"number"}
              value={config.MinImageNum}
              onChange={setNumberValue}
              placeholder={"MinImageNum"}
            ></InputWithLabel>

            <InputWithLabel
              label={t("TimeoutLimitForScan")}
              name={"TimeoutLimitForScan"}
              type={"numbers"}
              value={config.TimeoutLimitForScan}
              onChange={setNumberValue}
              placeholder={"TimeoutLimitForScan"}
            />

            <InputWithLabel
              label={t("ZipFileTextEncoding")}
              name={"ZipFileTextEncoding"}
              type={"text"}
              value={config.ZipFileTextEncoding}
              onChange={setStringValue}
              placeholder={"ZipFileTextEncoding"}
            />

            <BoolSwitch
              name={"GenerateMetaData"}
              label={t("GenerateMetaData")}
              boolValue={config.GenerateMetaData}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <StringArrayInput
              label={t("ExcludePath")}
              name={"ExcludePath"}
              value={config.ExcludePath}
              setStringArray={setStringArray}
            ></StringArrayInput>

            <StringArrayInput
              label={t("SupportMediaType")}
              name={"SupportMediaType"}
              value={config.SupportMediaType}
              setStringArray={setStringArray}
            ></StringArrayInput>

            <StringArrayInput
              label={t("SupportFileType")}
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
              name={"Port"}
              type={"number"}
              value={config.Port}
              onChange={setNumberValue}
              placeholder={"Port"}
            ></InputWithLabel>

            <InputWithLabel
              label={t("Host")}
              name={"Host"}
              type={"text"}
              value={config.Host}
              onChange={setStringValue}
              placeholder={"Host"}
            ></InputWithLabel>

            <BoolSwitch
              name={"DisableLAN"}
              label={t("DisableLAN")}
              boolValue={config.DisableLAN}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <BoolSwitch
              name={"EnableLogin"}
              label={t("EnableLogin")}
              boolValue={config.EnableLogin}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            {config.EnableLogin && <InputWithLabel
              label={t("Username")}
              name={"Username"}
              type={"text"}
              value={config.Username}
              onChange={setStringValue}
              placeholder={"Username"}
            ></InputWithLabel>}

            {config.EnableLogin && <InputWithLabel
              label={t("Password")}
              name={"Password"}
              type={"text"}
              value={config.Password}
              onChange={setStringValue}
              placeholder={"Password"}
            ></InputWithLabel>}

            <InputWithLabel
              label={t("Timeout")}
              name={"Timeout"}
              type={"number"}
              value={config.Timeout}
              onChange={setNumberValue}
              placeholder={"Timeout"}
            ></InputWithLabel>

            <BoolSwitch
              name={"EnableTLS"}
              label={t("EnableTLS")}
              boolValue={config.EnableTLS}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            {config.EnableTLS && <InputWithLabel
              label={t("CertFile")}
              name={"CertFile"}
              type={"text"}
              value={config.CertFile}
              onChange={setStringValue}
              placeholder={"CertFile"}
            ></InputWithLabel>}

            {config.EnableTLS && <InputWithLabel
              label={t("KeyFile")}
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
              boolValue={config.EnableDatabase}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            {config.EnableDatabase && <BoolSwitch
              name={"ClearDatabaseWhenExit"}
              label={t("ClearDatabaseWhenExit")}
              boolValue={config.ClearDatabaseWhenExit}
              setBoolValue={setBoolValue}
            ></BoolSwitch>}

            <BoolSwitch
              name={"Debug"}
              label={t("Debug")}
              boolValue={config.Debug}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            <BoolSwitch
              name={"LogToFile"}
              label={t("LogToFile")}
              boolValue={config.LogToFile}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

            {config.LogToFile &&
              <InputWithLabel
                label={t("LogFilePath")}
                name={"LogFilePath"}
                type={"text"}
                value={config.LogFilePath}
                onChange={setStringValue}
                placeholder={"LogFilePath"}
              />}

            {config.LogToFile &&
              <InputWithLabel
                label={t("LogFileName")}
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
              boolValue={config.EnableFrpcServer}
              setBoolValue={setBoolValue}
            ></BoolSwitch>

          </>
        }

      </form>
    </>
  );
}

export default App;