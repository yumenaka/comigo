import reactLogo from "./assets/react.svg";
import "./App.css";
import axios from "axios";
import { useEffect, useState } from "react";

type Config = {
  Port: number;
  Host: string;
  StoresPath: [];
  MaxDepth: number;
  OpenBrowser: boolean;
  DisableLAN: boolean;
  DefaultMode: string;
  UserName: string;
  Password: string;
  Timeout: number;
  CertFile: string;
  KeyFile: string;
  CacheEnable: boolean;
  CachePath: string;
  CacheClean: boolean;
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
  const [count, setCount] = useState(0);

  const baseURL = "/api";
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
      <div>
        <a
          className="flex flex-row justify-center items-center w-full"
          href="https://react.dev"
          target="_blank"
        >
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>React</h1>

      <div className="card">Now port is {config?.Port}</div>
    </>
  );
}

export default App;
