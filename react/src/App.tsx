// import reactLogo from "./assets/react.svg";
import "./App.css";
import axios from "axios";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useForm } from "react-hook-form";
import Title from "./components/Title";
import Config from "./types/Config";

function App() {
  const baseURL = "/api";
  const { t } = useTranslation();
  const [config, setConfig] = useState<Config | null>(null);

  const { register, handleSubmit } = useForm();

  const onSubmit = (data: any) => console.log(data);

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

  console.log(config);
  // React 通过  onChage 监听事件 实现数据的动态录入
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (config !== null) {
      setConfig({
        ...config,
        Port: Number(e.target.value),
      });
    }
  };
  //React 使用 value 或者 defaultValue 在 input 框中呈现内容
  return (
    <>
      <Title />
      <h2 className="text-lg font-semibold">admin</h2>

      <form
        onSubmit={handleSubmit(onSubmit)}
        className="card flex flex-col bg-slate-300 justify-center"
      >
        <button type="submit">ログイン</button>
        <div className="w-full flex flex-row bg-blue-300 justify-center items-center">
          <label htmlFor="Port">{t("Port")}:</label>
          <input
            id="Port"
            {...register("config.Port")}
            value={config?.Port}
            onChange={onChange}
          />
        </div>
        <p>
          {t("Host")}: {config?.Host}
        </p>
        <p>
          {t("StoresPath")}: {config?.StoresPath}
        </p>
        <p>
          {t("MaxScanDepth")}: {config?.MaxScanDepth}
        </p>
        <p>
          {t("OpenBrowser")}: {config?.OpenBrowser ? "true" : "false"}
        </p>
        <p>
          {t("DisableLAN")}: {config?.DisableLAN ? "true" : "false"}
        </p>
        <p>
          {t("Username")}: {config?.UserName}
        </p>
        <p>
          {t("Password")}: {config?.Password}
        </p>
        <p>
          {t("Timeout")}: {config?.Timeout}
        </p>
        <p>
          {t("CertFile")}: {config?.CertFile}
        </p>
        <p>
          {t("KeyFile")}: {config?.KeyFile}
        </p>
        <p>
          {t("EnableLocalCache")}: {config?.EnableLocalCache ? "true" : "false"}
        </p>
        <p>
          {t("CachePath")}: {config?.CachePath}
        </p>
        <p>
          {t("ClearCacheExit")}: {config?.ClearCacheExit}
        </p>
        <p>
          {t("EnableUpload")}: {config?.EnableUpload ? "true" : "false"}
        </p>
        <p>
          {t("UploadPath")}: {config?.UploadPath}
        </p>
        <p>
          {t("EnableDatabase")}: {config?.EnableDatabase}
        </p>
        <p>
          {t("ClearDatabaseWhenExit")}:{" "}
          {config?.ClearDatabase ? "true" : "false"}
        </p>
        <p>
          {t("ExcludePath")}: {config?.ExcludeFileOrFolders}
        </p>
        <p>
          {t("SupportedImageFileExtensions")}: {config?.SupportMediaType}
        </p>
        <p>
          {t("SupportedBookFileExtensions")}: {config?.SupportFileType}
        </p>
        <p>
          {t("MinImageCountInBook")}: {config?.MinImageNum}
        </p>
        <p>
          {t("TimeoutLimitForScan")}: {config?.TimeoutLimitForScan}
        </p>
        <p>
          {t("PrintAllPossibleQRCode")}: {config?.PrintAllIP ? "true" : "false"}
        </p>
        <p>
          {t("EnableDebugMode")}: {config?.Debug ? "true" : "false"}
        </p>
        <p>
          {t("LogToFile")}: {config?.LogToFile}
        </p>
        <p>
          {t("LogFilePath")}: {config?.LogFilePath}
        </p>
        <p>
          {t("LogFileName")}: {config?.LogFileName}
        </p>
        <p>
          {t("ZipFileTextEncoding")}: {config?.ZipFileTextEncoding}
        </p>
        <p>
          {t("StartFrpClientInBackground")}:{" "}
          {config?.EnableFrpcServer ? "true" : "false"}
        </p>
        <p>
          {t("GenerateBookMetadata")}:{" "}
          {config?.GenerateMetaData ? "true" : "false"}
        </p>
        {/* {t("FrpClientConfig")}: {config?.FrpConfig} */}
      </form>
    </>
  );
}

export default App;
