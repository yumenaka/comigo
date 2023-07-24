// import reactLogo from "./assets/react.svg";
import "./App.css";
import axios from "axios";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

type Config = {
  Port: number;
  Host: string;
  StoresPath: [];
  MaxScanDepth: number;
  OpenBrowser: boolean;
  DisableLAN: boolean;
  DefaultMode: string;
  UserName: string;
  Password: string;
  Timeout: number;
  CertFile: string;
  KeyFile: string;
  EnableLocalCache: boolean;
  CachePath: string;
  ClearCacheExit: boolean;
  EnableUpload: boolean;
  UploadPath: string;
  EnableDatabase: boolean;
  ClearDatabase: boolean;
  ExcludeFileOrFolders: [];
  SupportMediaType: [];
  SupportFileType: [];
  MinImageNum: number;
  TimeoutLimitForScan: number;
  PrintAllIP: boolean;
  Debug: boolean;
  LogToFile: boolean;
  LogFilePath: string;
  LogFileName: string;
  ZipFileTextEncoding: string;
  EnableFrpcServer: boolean;
  FrpConfig: {
    FrpcCommand: string;
    ServerAddr: string;
    ServerPort: number;
    Token: string;
    FrpType: string;
    RemotePort: number;
    RandomRemotePort: boolean;
  };
  GenerateMetaData: boolean;
};

function App() {
  const baseURL = "/api";
  const { t, i18n } = useTranslation();
  const [config, setConfig] = useState<Config | null>(null);

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

  return (
    <>
      <h1>admin</h1>
      <div className="card">
        {t("Port")} : {config?.Port} <br />
        {t("Host")} : {config?.Host} <br />
        {t("StoresPath")} : {config?.StoresPath} <br />
        {t("MaxScanDepth")} : {config?.MaxScanDepth} <br />
        {t("OpenBrowser")} : {config?.OpenBrowser ? "true" : "false"} <br />
        {t("DisableLAN")} : {config?.DisableLAN ? "true" : "false"} <br />
        {t("Username")} : {config?.UserName} <br />
        {t("Password")} : {config?.Password} <br />
        {t("Timeout")} : {config?.Timeout} <br />
        {t("CertFile")} : {config?.CertFile} <br />
        {t("KeyFile")} : {config?.KeyFile} <br />
        {t("EnableLocalCache")} : {config?.EnableLocalCache ? "true" : "false"}{" "} <br />
        {t("CachePath")} : {config?.CachePath} <br />
        {t("ClearCacheExit")} : {config?.ClearCacheExit} <br />
        {t("EnableUpload")} : {config?.EnableUpload ? "true" : "false"} <br />
        {t("UploadPath")} : {config?.UploadPath} <br />
        {t("EnableDatabase")} : {config?.EnableDatabase} <br />
        {t("ClearDatabaseWhenExit")} :{" "}
        {config?.ClearDatabase ? "true" : "false"} <br />
        {t("ExcludePath")} : {config?.ExcludeFileOrFolders} <br />
        {t("SupportedImageFileExtensions")} : {config?.SupportMediaType} <br />
        {t("SupportedBookFileExtensions")} : {config?.SupportFileType} <br />
        {t("MinImageCountInBook")} : {config?.MinImageNum} <br />
        {t("TimeoutLimitForScan")} : {config?.TimeoutLimitForScan} <br />
        {t("PrintAllPossibleQRCode")} : {config?.PrintAllIP ? "true" : "false"}{" "}
        <br />
        {t("EnableDebugMode")} : {config?.Debug ? "true" : "false"} <br />
        {t("LogToFile")} : {config?.LogToFile} <br />
        {t("LogFilePath")} : {config?.LogFilePath} <br />
        {t("LogFileName")} : {config?.LogFileName} <br />
        {t("ZipFileTextEncoding")} : {config?.ZipFileTextEncoding} <br />
        {t("StartFrpClientInBackground")} :{" "}
        {config?.EnableFrpcServer ? "true" : "false"} <br />
        {t("GenerateBookMetadata")} :{" "}
        {config?.GenerateMetaData ? "true" : "false"} <br />
        {/* {t("FrpClientConfig")} : {config?.FrpConfig} */}
      </div>
    </>
  );
}

export default App;
