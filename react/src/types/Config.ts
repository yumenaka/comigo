// 取得したいデータの型を定義する
//http://note.alvinhtml.com/web/typescript.html
interface Config {
  Port: number;
  Host: string;
  StoresPath: [];
  MaxScanDepth: number;
  OpenBrowser: boolean;
  DisableLAN: boolean;
  DefaultMode: string;
  Username: string;
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
  ClearDatabaseWhenExit: boolean;
  ExcludePath: [];
  SupportMediaType: [];
  SupportFileType: [];
  MinImageNum: number;
  TimeoutLimitForScan: number;
  PrintAllPossibleQRCode: boolean;
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
}

export default Config;