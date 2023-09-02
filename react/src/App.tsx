import axios from "axios";
import React, { useEffect, useReducer } from "react";
import { useTranslation } from "react-i18next";
// TypeScript環境でReact Hook Formのフォーム作成の基礎を学ぶ  https://reffect.co.jp/react/react-hook-form-ts/
// import { useForm } from "react-hook-form";
import Title from "./components/Title";
import Config from "./types/Config";
import InputWithLabel from "./components/InputWithLabel";

import BoolSwitch from "./components/BoolSwitch";
import { configReducer, defaultConfig } from "./reducers/configReducer";

function App() {
  const baseURL = "/api";
  const { t } = useTranslation();
  // nullは型に含めず、useStateの初期値が決まらないという場合には、型アサーションで逃げる手もあります。 https://qiita.com/FumioNonaka/items/4361d1cdf34ffb5a5338
  const [config, dispatch] = useReducer(configReducer, defaultConfig);

  const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    axios.post("/api/post_config", config).then((response) => {
      console.log("Data sent successfully");
      //axios默认解析Json，所以 response.data 就是解析后的object
      console.info(response.data);
    })
      .catch((error) => {
        console.error("Error sending data:", error);
      });
  };
  // useEffect 用于在函数组件中执行副作用操作，例如数据获取、订阅、手动修改DOM等。
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
    // 通过传递第二个参数，你可以告诉 React 仅在某些值改变的时候才执行 effect。
    // 传递空数组([])作为第二个参数，effect 内部的 props 和 state 就会一直持有其初始值。也就是只在渲染的时候执行一次。
  }, []);

  // React 通过  onChange 监听事件 实现数据的动态录入
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    console.log(typeof value);
    //setConfig({ ...config, [name]: value });
    dispatch({
      type: 'boolConfig',
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

  
  return (
    <>
      <Title />
      <form
        onSubmit={onSubmit}
        className="card w-full flex flex-col bg-slate-300 justify-center items-center"
      >
        <button
          type="submit"
          className="w-32 font-semibold bg-green-400 justify-center items-center m-1"
        >
          submit
        </button>
        {/* React 使用 value 或者 defaultValue 在 input 框中呈现内容 */}

        <InputWithLabel
          label={t("Port")}
          name={"Port"}
          type={"number"}
          value={config.Port}
          onChange={onChange}
          placeholder={"Port"}
        ></InputWithLabel>

        <InputWithLabel
          label={t("Host")}
          name={"Host"}
          type={"text"}
          value={config.Host}
          onChange={onChange}
          placeholder={"Host"}
        ></InputWithLabel>

        <InputWithLabel
          label={t("StoresPath")}
          name={"StoresPath"}
          type={"text"}
          value={config.StoresPath}
          onChange={onChange}
          placeholder={"StoresPath"}
        ></InputWithLabel>

        <InputWithLabel
          label={t("MaxScanDepth")}
          name={"MaxScanDepth"}
          type={"number"}
          value={config.MaxScanDepth}
          onChange={onChange}
          placeholder={"MaxScanDepth"}
        ></InputWithLabel>

        <BoolSwitch
          label={t("OpenBrowser")}
          name={"OpenBrowser"}
          boolValue={config.OpenBrowser}
          setBoolValue={setBoolValue}
        ></BoolSwitch>

        <BoolSwitch
          name={"DisableLAN"}
          label={t("DisableLAN")}
          boolValue={config.DisableLAN}
          setBoolValue={setBoolValue}
        ></BoolSwitch>

        <InputWithLabel
          label={t("Username")}
          name={"Username"}
          type={"text"}
          value={config.Username}
          onChange={onChange}
          placeholder={"Username"}
        ></InputWithLabel>

        <InputWithLabel
          label={t("Password")}
          name={"Password"}
          type={"text"}
          value={config.Password}
          onChange={onChange}
          placeholder={"Password"}
        ></InputWithLabel>

        <InputWithLabel
          label={t("Timeout")}
          name={"Timeout"}
          type={"number"}
          value={config.Timeout}
          onChange={onChange}
          placeholder={"Timeout"}
        ></InputWithLabel>

        <InputWithLabel
          label={t("CertFile")}
          name={"CertFile"}
          type={"text"}
          value={config.CertFile}
          onChange={onChange}
          placeholder={"CertFile"}
        ></InputWithLabel>

        <InputWithLabel
          label={t("KeyFile")}
          name={"KeyFile"}
          type={"text"}
          value={config.KeyFile}
          onChange={onChange}
          placeholder={"KeyFile"}
        ></InputWithLabel>

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
          onChange={onChange}
          placeholder={"CachePath"}
        ></InputWithLabel>

        <BoolSwitch
          name={"EnableUpload"}
          label={t("EnableUpload")}
          boolValue={config.EnableUpload}
          setBoolValue={setBoolValue}
        ></BoolSwitch>

        <InputWithLabel
          label={t("UploadPath")}
          name={"UploadPath"}
          type={"text"}
          value={config.UploadPath}
          onChange={onChange}
          placeholder={"UploadPath"}
        />

        <InputWithLabel
          label={t("ExcludePath")}
          name={"ExcludePath"}
          type={"text"}
          value={config.ExcludePath}
          onChange={onChange}
          placeholder={"ExcludePath"}
        />

        <InputWithLabel
          label={t("SupportMediaType")}
          name={"SupportMediaType"}
          type={"text"}
          value={config.SupportMediaType}
          onChange={onChange}
          placeholder={"SupportMediaType"}
        />

        <InputWithLabel
          label={t("SupportFileType")}
          name={"SupportFileType"}
          type={"text"}
          value={config.SupportFileType}
          onChange={onChange}
          placeholder={"SupportFileType"}
        />

        <InputWithLabel
          label={t("MinImageNum")}
          name={"MinImageNum"}
          type={"number"}
          value={config.MinImageNum}
          onChange={onChange}
          placeholder={"MinImageNum"}
        ></InputWithLabel>

        <InputWithLabel
          label={t("TimeoutLimitForScan")}
          name={"TimeoutLimitForScan"}
          type={"numbers"}
          value={config.TimeoutLimitForScan}
          onChange={onChange}
          placeholder={"TimeoutLimitForScan"}
        />

        <BoolSwitch
          name={"PrintAllPossibleQRCode"}
          label={t("PrintAllPossibleQRCode")}
          boolValue={config.PrintAllPossibleQRCode}
          setBoolValue={setBoolValue}
        ></BoolSwitch>

        <BoolSwitch
          name={"EnableDatabase"}
          label={t("EnableDatabase")}
          boolValue={config.EnableDatabase}
          setBoolValue={setBoolValue}
        ></BoolSwitch>

        <BoolSwitch
          name={"ClearDatabaseWhenExit"}
          label={t("ClearDatabaseWhenExit")}
          boolValue={config.ClearDatabaseWhenExit}
          setBoolValue={setBoolValue}
        ></BoolSwitch>

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

        <InputWithLabel
          label={t("LogFilePath")}
          name={"LogFilePath"}
          type={"text"}
          value={config.LogFilePath}
          onChange={onChange}
          placeholder={"LogFilePath"}
        />

        <InputWithLabel
          label={t("LogFileName")}
          name={"LogFileName"}
          type={"text"}
          value={config.LogFileName}
          onChange={onChange}
          placeholder={"LogFileName"}
        />

        <InputWithLabel
          label={t("ZipFileTextEncoding")}
          name={"ZipFileTextEncoding"}
          type={"text"}
          value={config.ZipFileTextEncoding}
          onChange={onChange}
          placeholder={"ZipFileTextEncoding"}
        />

        <BoolSwitch
          name={"EnableFrpcServer"}
          label={t("EnableFrpcServer")}
          boolValue={config.EnableFrpcServer}
          setBoolValue={setBoolValue}
        ></BoolSwitch>

        <BoolSwitch
          name={"GenerateMetaData"}
          label={t("GenerateMetaData")}
          boolValue={config.GenerateMetaData}
          setBoolValue={setBoolValue}
        ></BoolSwitch>
      </form>
    </>
  );
}

export default App;