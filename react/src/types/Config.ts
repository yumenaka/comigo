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
  EnableLogin: boolean;
  Username: string;
  Password: string;
  Timeout: number;
  EnableTLS: boolean;
  CertFile: string;
  KeyFile: string;
  UseCache: boolean;
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
  GenerateMetaData: boolean;
}

export default Config;
