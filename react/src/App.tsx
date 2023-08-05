// import reactLogo from "./assets/react.svg";
import "./App.css";
import axios from "axios";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
// TypeScript環境でReact Hook Formのフォーム作成の基礎を学ぶ  https://reffect.co.jp/react/react-hook-form-ts/
import { useForm } from "react-hook-form";
import Title from "./components/Title";
import Config from "./types/Config";


function App() {
  const baseURL = "/api";
  const { t } = useTranslation();
  // nullは型に含めず、useStateの初期値が決まらないという場合には、型アサーションで逃げる手もあります。 https://qiita.com/FumioNonaka/items/4361d1cdf34ffb5a5338
  // ただし、型アサーションはTypeScriptに値({})の型を偽っているだけだ、ということにご注意ください。状態変数(Config)の値を正しく扱うことは、コードの書き手に委ねられるのです。
  //誤ればランタイムエラーになってしまうかもしれません。
  const [config, setConfig] = useState<Config>({} as Config);

  //不使用react-hook-form的话，一般的handleSubmit函数如下：
  // const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
  //   e.preventDefault();
  //   console.log({
  //     email,
  //     password,
  //   });
  // };

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
    // sendDataToBackend(config);
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
    console.log("OnChange:", name, value);
    // 在HTML中，属性的类型是字符串。所以，我们需要将字符串转换为布尔值。
    if (value === "true" || value === "false") {
      setConfig({ ...config, [name]: value === "true" });
      return;
    }
    setConfig({ ...config, [name]: value });
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
        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-start items-center">
          <label htmlFor="Port" className="w-32 border border-black rounded-md">
            {t("Port")}:
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("Port", { min: 0, max: 65535 })}
            id="Port"
            type="number"
            value={config.Port}
            onChange={onChange}
            placeholder="Port"
          />
          <div className="bg-red-600">
            {errors.Port && <div>入力が必須の項目です(0~65535)</div>}
          </div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label htmlFor="Host" className="w-32 border border-black rounded-md">
            {t("Host")}:
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("Host")}
            id="Host"
            type="text"
            value={config.Host}
            onChange={onChange}
            placeholder="Host"
          />
          <div className="bg-red-600">{errors.Host && <div></div>}</div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="StoresPath"
            className="w-32 border border-black rounded-md"
          >
            {t("StoresPath")}:
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("StoresPath")}
            id="StoresPath"
            type="text"
            value={config.StoresPath}
            onChange={onChange}
            placeholder="StoresPath"
          />
          <div className="bg-red-600">{errors.StoresPath && <div></div>}</div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="MaxScanDepth"
            className="w-32 border border-black rounded-md"
          >
            {t("MaxScanDepth")}
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("MaxScanDepth")}
            id="MaxScanDepth"
            type="numbers"
            value={config.MaxScanDepth}
            onChange={onChange}
            placeholder="MaxScanDepth"
          />
          <div className="bg-red-600">{errors.MaxScanDepth && <div></div>}</div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <div className="w-32 border border-black rounded-md">
            {t("OpenBrowser")}
          </div>

          <input
            type="radio"
            id="OpenBrowserTrue"
            value={config.OpenBrowser ? "true" : "false"}
            {...register("OpenBrowser")}
            onChange={onChange}
          />
          <label htmlFor="OpenBrowserTrue">是</label>
          <input
            type="radio"
            id="OpenBrowserFalse"
            value={!config.OpenBrowser ? "true" : "false"}
            checked={config.OpenBrowser === false}

            {...register("OpenBrowser", { required: true })}
          />
          <label htmlFor="OpenBrowserFalse">否</label>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <div className="w-32 border border-black rounded-md">
            {t("DisableLAN")}
          </div>
          <input
            type="radio"
            id="DisableLANTrue"
            name="DisableLAN"
            value={config.DisableLAN ? "true" : "false"}
          />
          <label htmlFor="DisableLANTrue">是</label>
          <input
            type="radio"
            id="DisableLANFalse"
            name="DisableLAN"
            value={!config.DisableLAN ? "true" : "false"}
          />
          <label htmlFor="DisableLANFalse">否</label>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="UserName"
            className="w-32 border border-black rounded-md"
          >
            {t("Username")}:
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("UserName")}
            id="UserName"
            type="text"
            value={config.UserName}
            onChange={onChange}
            placeholder="UserName"
          />
          <div className="bg-red-600">{errors.UserName && <div></div>}</div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="Password"
            className="w-32 border border-black rounded-md"
          >
            {t("Password")}:
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("Password")}
            id="Password"
            type="text"
            value={config.Password}
            onChange={onChange}
            placeholder="Password"
          />
          <div className="bg-red-600">{errors.Password && <div></div>}</div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="Timeout"
            className="w-32 border border-black rounded-md"
          >
            {t("Timeout")}
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("Timeout")}
            id="Timeout"
            type="numbers"
            value={config.Timeout}
            onChange={onChange}
            placeholder="Timeout"
          />
          <div className="bg-red-600">{errors.Timeout && <div></div>}</div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="CertFile"
            className="w-32 border border-black rounded-md"
          >
            {t("CertFile")}:
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("CertFile")}
            id="CertFile"
            type="text"
            value={config.CertFile}
            onChange={onChange}
            placeholder="CertFile"
          />
          <div className="bg-red-600">{errors.CertFile && <div></div>}</div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="KeyFile"
            className="w-32 border border-black rounded-md"
          >
            {t("KeyFile")}:
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("KeyFile")}
            id="KeyFile"
            type="text"
            value={config.KeyFile}
            onChange={onChange}
            placeholder="KeyFile"
          />
          <div className="bg-red-600">{errors.KeyFile && <div></div>}</div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <div className="w-32 border border-black rounded-md">
            {t("ClearCacheExit")}
          </div>
          <input
            type="radio"
            id="ClearCacheExitTrue"
            name="ClearCacheExit"
            value={config.ClearCacheExit ? "true" : "false"}
          />
          <label htmlFor="ClearCacheExitTrue">是</label>
          <input
            type="radio"
            id="ClearCacheExitFalse"
            name="ClearCacheExit"
            value={!config.ClearCacheExit ? "true" : "false"}
          />
          <label htmlFor="ClearCacheExitFalse">否</label>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="CachePath"
            className="w-32 border border-black rounded-md"
          >
            {t("CachePath")}:
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("CachePath")}
            id="CachePath"
            type="text"
            value={config.CachePath}
            onChange={onChange}
            placeholder="CachePath"
          />
          <div className="bg-red-600">{errors.CachePath && <div></div>}</div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <div className="w-32 border border-black rounded-md">
            {t("ClearCacheExit")}
          </div>
          <input
            type="radio"
            id="ClearCacheExitTrue"
            name="ClearCacheExit"
            value={config.ClearCacheExit ? "true" : "false"}
          />
          <label htmlFor="ClearCacheExitTrue">是</label>
          <input
            type="radio"
            id="ClearCacheExitFalse"
            name="ClearCacheExit"
            value={!config.ClearCacheExit ? "true" : "false"}
          />
          <label htmlFor="ClearCacheExitFalse">否</label>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <div className="w-32 border border-black rounded-md">
            {t("EnableUpload")}
          </div>
          <input
            type="radio"
            id="EnableUploadTrue"
            name="EnableUpload"
            value={config.EnableUpload ? "true" : "false"}
          />
          <label htmlFor="EnableUploadTrue">是</label>
          <input
            type="radio"
            id="EnableUploadFalse"
            name="EnableUpload"
            value={!config.EnableUpload ? "true" : "false"}
          />
          <label htmlFor="EnableUploadFalse">否</label>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="UploadPath"
            className="w-32 border border-black rounded-md"
          >
            {t("UploadPath")}:
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("UploadPath")}
            id="UploadPath"
            type="text"
            value={config.UploadPath}
            onChange={onChange}
            placeholder="UploadPath"
          />
          <div className="bg-red-600">{errors.UploadPath && <div></div>}</div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <div className="w-32 border border-black rounded-md">
            {t("EnableDatabase")}
          </div>
          <input
            type="radio"
            id="EnableDatabaseTrue"
            name="EnableDatabase"
            value={config.EnableDatabase ? "true" : "false"}
          />
          <label htmlFor="EnableDatabaseTrue">是</label>
          <input
            type="radio"
            id="EnableDatabaseFalse"
            name="EnableDatabase"
            value={!config.EnableDatabase ? "true" : "false"}
          />
          <label htmlFor="EnableDatabaseFalse">否</label>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <div className="w-32 border border-black rounded-md">
            {t("ClearDatabaseWhenExit")}
          </div>
          <input
            type="radio"
            id="ClearDatabaseTrue"
            name="ClearDatabase"
            value={config.ClearDatabase ? "true" : "false"}
          />
          <label htmlFor="ClearDatabaseTrue">是</label>
          <input
            type="radio"
            id="ClearDatabaseFalse"
            name="ClearDatabase"
            value={!config.ClearDatabase ? "true" : "false"}
          />
          <label htmlFor="ClearDatabaseFalse">否</label>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="ExcludeFileOrFolders"
            className="w-32 border border-black rounded-md"
          >
            {t("ExcludePath")}:
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("ExcludeFileOrFolders")}
            id="ExcludeFileOrFolders"
            type="text"
            value={config.ExcludeFileOrFolders}
            onChange={onChange}
            placeholder="ExcludeFileOrFolders"
          />
          <div className="bg-red-600">
            {errors.ExcludeFileOrFolders && <div></div>}
          </div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="SupportMediaType"
            className="w-32 border border-black rounded-md"
          >
            {t("SupportedImageFileExtensions")}:
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("SupportMediaType")}
            id="SupportMediaType"
            type="text"
            value={config.SupportMediaType}
            onChange={onChange}
            placeholder="SupportMediaType"
          />
          <div className="bg-red-600">
            {errors.SupportMediaType && <div></div>}
          </div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="SupportFileType"
            className="w-32 border border-black rounded-md"
          >
            {t("SupportedBookFileExtensions")}:
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("SupportFileType")}
            id="SupportFileType"
            type="text"
            value={config.SupportFileType}
            onChange={onChange}
            placeholder="SupportFileType"
          />
          <div className="bg-red-600">
            {errors.SupportFileType && <div></div>}
          </div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="MinImageNum"
            className="w-32 border border-black rounded-md"
          >
            {t("MinImageCountInBook")}
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("MinImageNum")}
            id="MinImageNum"
            type="numbers"
            value={config.MinImageNum}
            onChange={onChange}
            placeholder="MinImageNum"
          />
          <div className="bg-red-600">{errors.MinImageNum && <div></div>}</div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="TimeoutLimitForScan"
            className="w-32 border border-black rounded-md"
          >
            {t("TimeoutLimitForScan")}
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("TimeoutLimitForScan")}
            id="TimeoutLimitForScan"
            type="numbers"
            value={config.TimeoutLimitForScan}
            onChange={onChange}
            placeholder="TimeoutLimitForScan"
          />
          <div className="bg-red-600">
            {errors.TimeoutLimitForScan && <div></div>}
          </div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <div className="w-32 border border-black rounded-md">
            {t("PrintAllPossibleQRCode")}
          </div>
          <input
            type="radio"
            id="PrintAllIPTrue"
            name="PrintAllIP"
            value={config.PrintAllIP ? "true" : "false"}
          />
          <label htmlFor="PrintAllIPTrue">是</label>
          <input
            type="radio"
            id="PrintAllIPFalse"
            name="PrintAllIP"
            value={!config.PrintAllIP ? "true" : "false"}
          />
          <label htmlFor="PrintAllIPFalse">否</label>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <div className="w-32 border border-black rounded-md">
            {t("EnableDebugMode")}
          </div>
          <input
            type="radio"
            id="DebugTrue"
            name="Debug"
            value={config.Debug ? "true" : "false"}
          />
          <label htmlFor="DebugTrue">是</label>
          <input
            type="radio"
            id="DebugFalse"
            name="Debug"
            value={!config.Debug ? "true" : "false"}
          />
          <label htmlFor="DebugFalse">否</label>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <div className="w-32 border border-black rounded-md">
            {t("LogToFile")}
          </div>
          <input
            type="radio"
            id="LogToFileTrue"
            name="LogToFile"
            value={config.LogToFile ? "true" : "false"}
          />
          <label htmlFor="LogToFileTrue">是</label>
          <input
            type="radio"
            id="LogToFileFalse"
            name="LogToFile"
            value={!config.LogToFile ? "true" : "false"}
          />
          <label htmlFor="LogToFileFalse">否</label>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="LogFilePath"
            className="w-32 border border-black rounded-md"
          >
            {t("LogFilePath")}:
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("LogFilePath")}
            id="LogFilePath"
            type="text"
            value={config.LogFilePath}
            onChange={onChange}
            placeholder="LogFilePath"
          />
          <div className="bg-red-600">{errors.LogFilePath && <div></div>}</div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="LogFileName"
            className="w-32 border border-black rounded-md"
          >
            {t("LogFileName")}:
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("LogFileName")}
            id="LogFileName"
            type="text"
            value={config.LogFileName}
            onChange={onChange}
            placeholder="LogFileName"
          />
          <div className="bg-red-600">{errors.LogFileName && <div></div>}</div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <label
            htmlFor="ZipFileTextEncoding"
            className="w-32 border border-black rounded-md"
          >
            {t("ZipFileTextEncoding")}:
          </label>
          <input
            className="rounded ml-2 px-1"
            {...register("ZipFileTextEncoding")}
            id="ZipFileTextEncoding"
            type="text"
            value={config.ZipFileTextEncoding}
            onChange={onChange}
            placeholder="ZipFileTextEncoding"
          />
          <div className="bg-red-600">
            {errors.ZipFileTextEncoding && <div></div>}
          </div>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <div className="w-32 border border-black rounded-md">
            {t("EnableFrpc")}
          </div>
          <input
            type="radio"
            id="EnableFrpcServerTrue"
            name="EnableFrpcServer"
            value={config.EnableFrpcServer ? "true" : "false"}
          />
          <label htmlFor="EnableFrpcServerTrue">是</label>
          <input
            type="radio"
            id="EnableFrpcServerFalse"
            name="EnableFrpcServer"
            value={!config.EnableFrpcServer ? "true" : "false"}
          />
          <label htmlFor="EnableFrpcServerFalse">否</label>
        </div>

        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
          <div className="w-32 border border-black rounded-md">
            {t("GenerateMetaData")}
          </div>
          <input
            type="radio"
            id="GenerateMetaDataTrue"
            name="GenerateMetaData"
            value={config.GenerateMetaData ? "true" : "false"}
          />
          <label htmlFor="GenerateMetaDataTrue">是</label>
          <input
            type="radio"
            id="GenerateMetaDataFalse"
            name="GenerateMetaData"
            value={!config.GenerateMetaData ? "true" : "false"}
          />
          <label htmlFor="GenerateMetaDataFalse">否</label>
        </div>

        {/* {t("FrpClientConfig")}: {config.FrpConfig} */}
      </form>
    </>
  );
}

export default App;
