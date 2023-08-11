// import reactLogo from "./assets/react.svg";
import "./App.css";
import axios from "axios";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
// TypeScript環境でReact Hook Formのフォーム作成の基礎を学ぶ  https://reffect.co.jp/react/react-hook-form-ts/
import { useForm } from "react-hook-form";
import Title from "./components/Title";
import Config from "./types/Config";
import InputWithLabel from "./components/InputWithLabel";

import BoolSwitch from "./components/BoolSwitch";

function App() {
  const baseURL = "/api";
  const { t } = useTranslation();
  // nullは型に含めず、useStateの初期値が決まらないという場合には、型アサーションで逃げる手もあります。 https://qiita.com/FumioNonaka/items/4361d1cdf34ffb5a5338
  // ただし、型アサーションはTypeScriptに値({})の型を偽っているだけだ、ということにご注意ください。状態変数(Config)の値を正しく扱うことは、コードの書き手に委ねられるのです。
  //誤ればランタイムエラーになってしまうかもしれません。
  const [config, setConfig] = useState<Config>({} as Config);



  // 使用react-hook-form的话，handleSubmit函数如下：
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    defaultValues: {
      Port: config.Port,
      Host: config.Host,
      StoresPath: config.StoresPath,
      MaxScanDepth: config.MaxScanDepth,
      OpenBrowser: config.OpenBrowser,
      DisableLAN: config.DisableLAN,
      UserName: config.UserName,
      Password: config.Password,
      Timeout: config.Timeout,
      CertFile: config.CertFile,
      KeyFile: config.KeyFile,
      EnableLocalCache: config.EnableLocalCache,
      CachePath: config.CachePath,
      ClearCacheExit: config.ClearCacheExit,
      EnableUpload: config.EnableUpload,
      UploadPath: config.UploadPath,
      EnableDatabase: config.EnableDatabase,
      ClearDatabase: config.ClearDatabase,
      ExcludeFileOrFolders: config.ExcludeFileOrFolders,
      SupportMediaType: config.SupportMediaType,
      SupportFileType: config.SupportFileType,
      MinImageNum: config.MinImageNum,
      TimeoutLimitForScan: config.TimeoutLimitForScan,
      PrintAllIP: config.PrintAllIP,
      Debug: config.Debug,
      LogToFile: config.LogToFile,
      LogFilePath: config.LogFilePath,
      LogFileName: config.LogFileName,
      ZipFileTextEncoding: config.ZipFileTextEncoding,
      EnableFrpcServer: config.EnableFrpcServer,
      FrpConfig: config.FrpConfig,
      GenerateMetaData: config.GenerateMetaData,
    },
  });

  const sendDataToBackend = async (data: Config) => {
    console.log("sendDataToBackend:", data);
    try {
      const response = await axios.post('/api/update_config', data);
      console.log('Data sent successfully:', response.data);
      // 可以在此处进行其他操作，例如更新状态或显示成功消息等
    } catch (error) {
      console.error('Error sending data:', error);
      // 可以在此处处理错误，例如显示错误消息等
    }
  };

  const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    console.log(e);
  }

  useEffect(() => {
    axios
      .get<Config>(`${baseURL}/config.json`)
      .then((response) => {
        setConfig(response.data);
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);

  //console.log(config);
  // React 通过  onChange 监听事件 实现数据的动态录入
  // html的属性props的类型
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    console.log(typeof value);
    setConfig({ ...config, [name]: value });
  };

  const setBoolConfig = (checked: boolean, boolValueName: string) => {
    // console.log("OnChange" + boolValueName, checked);

    setConfig({ ...config, [boolValueName]: checked });
  };


  //React 使用 value 或者 defaultValue 在 input 框中呈现内容
  return (
    <>
      <Title />

      <form
        onSubmit={onSubmit}
        className="card flex flex-col bg-slate-300 justify-center items-center"
      >
        <button
          type="submit"
          className="w-32 font-semibold bg-green-400 justify-center items-center m-1"
        >
          submit
        </button>

        <InputWithLabel label={t("Port")} name={"Port"} type={"number"} value={config.Port} onChange={onChange} placeholder={"Port"}   ></InputWithLabel>

        <InputWithLabel label={t("Host")} name={"Host"} type={"text"} value={config.Host} onChange={onChange} placeholder={"Host"}></InputWithLabel>

        <InputWithLabel label={t("StoresPath")} name={"StoresPath"} type={"text"} value={config.StoresPath} onChange={onChange} placeholder={"StoresPath"}></InputWithLabel>

        <InputWithLabel label={t("MaxScanDepth")} name={"MaxScanDepth"} type={"number"} value={config.MaxScanDepth} onChange={onChange} placeholder={"MaxScanDepth"}></InputWithLabel>

        <BoolSwitch label={t("OpenBrowser")} name={"OpenBrowser"} boolValue={config.OpenBrowser} setBoolConfig={setBoolConfig} ></BoolSwitch>

        <BoolSwitch name={"DisableLAN"} label={t("DisableLAN")} boolValue={config.DisableLAN} setBoolConfig={setBoolConfig} ></BoolSwitch>

        <InputWithLabel label={t("Username")} name={"UserName"} type={"text"} value={config.UserName} onChange={onChange} placeholder={"UserName"}></InputWithLabel>

        <InputWithLabel label={t("Password")} name={"Password"} type={"text"} value={config.Password} onChange={onChange} placeholder={"Password"} ></InputWithLabel>

        <InputWithLabel label={t("Timeout")} name={"Timeout"} type={"number"} value={config.Timeout} onChange={onChange} placeholder={"Timeout"} ></InputWithLabel>

        <InputWithLabel label={t("CertFile")} name={"CertFile"} type={"text"} value={config.CertFile} onChange={onChange} placeholder={"CertFile"} ></InputWithLabel>

        <InputWithLabel label={t("KeyFile")} name={"KeyFile"} type={"text"} value={config.KeyFile} onChange={onChange} placeholder={"KeyFile"} ></InputWithLabel>

        <BoolSwitch name={"ClearCacheExit"} label={t("ClearCacheExit")} boolValue={config.ClearCacheExit} setBoolConfig={setBoolConfig} ></BoolSwitch>

        <InputWithLabel label={t("CachePath")} name={"CachePath"} type={"text"} value={config.CachePath} onChange={onChange} placeholder={"CachePath"}></InputWithLabel>

        <BoolSwitch name={"EnableUpload"} label={t("EnableUpload")} boolValue={config.EnableUpload} setBoolConfig={setBoolConfig} ></BoolSwitch>

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
          name={"ExcludeFileOrFolders"}
          type={"text"}
          value={config.ExcludeFileOrFolders}
          onChange={onChange}
          placeholder={"ExcludeFileOrFolders"}
        />

        <InputWithLabel
          label={t("SupportedImageFileExtensions")}
          name={"SupportMediaType"}
          type={"text"}
          value={config.SupportMediaType}
          onChange={onChange}
          placeholder={"SupportMediaType"}
        />

        <InputWithLabel
          label={t("SupportedBookFileExtensions")}
          name={"SupportFileType"}
          type={"text"}
          value={config.SupportFileType}
          onChange={onChange}
          placeholder={"SupportFileType"}
        />



        <InputWithLabel
          label={t("MinImageCountInBook")}
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


        <BoolSwitch name={"PrintAllIP"} label={t("PrintAllPossibleQRCode")} boolValue={config.PrintAllIP} setBoolConfig={setBoolConfig} ></BoolSwitch>

        <BoolSwitch name={"EnableDatabase"} label={t("EnableDatabase")} boolValue={config.EnableDatabase} setBoolConfig={setBoolConfig} ></BoolSwitch>

        <BoolSwitch name={"ClearDatabase"} label={t("ClearDatabaseWhenExit")} boolValue={config.ClearDatabase} setBoolConfig={setBoolConfig} ></BoolSwitch>

        <BoolSwitch name={"Debug"} label={t("EnableDebugMode")} boolValue={config.Debug} setBoolConfig={setBoolConfig} ></BoolSwitch>

        <BoolSwitch name={"LogToFile"} label={t("LogToFile")} boolValue={config.LogToFile} setBoolConfig={setBoolConfig} ></BoolSwitch>

        <InputWithLabel
          label={t("LogFilePath")}
          name={"TimeoutLimitForScan"}
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

        <BoolSwitch name={"EnableFrpcServer"} label={t("EnableFrpc")} boolValue={config.EnableFrpcServer} setBoolConfig={setBoolConfig} ></BoolSwitch>

        <BoolSwitch name={"GenerateMetaData"} label={t("GenerateMetaData")} boolValue={config.GenerateMetaData} setBoolConfig={setBoolConfig} ></BoolSwitch>

      </form>
    </>
  );
}

export default App;
