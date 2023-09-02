import Config  from "../types/Config";

type Action = {
    type:
    | "downloadConfig"
    | "boolConfig"
    | "stringConfig"
    | "numberConfig"
    | "arrayConfig";
    name: string;
    value: boolean | string | number | string[];
    config: Config;
};

export function configReducer(c: Config, action: Action) {
    switch (action.type) {
        case "downloadConfig":
            return { ...action.config };
        case "boolConfig":
            return { ...c, [action.name]: action.value };
        case "stringConfig":
            return { ...c, [action.name]: action.value };
        case "numberConfig":
            return { ...c, [action.name]: action.value };
        case "arrayConfig":
            return { ...c, [action.name]: action.value };
        default:
            console.log(action);
            throw new Error();
    }
}

export const defaultConfig:Config={
    Port: 1234,
    Host: "",
    StoresPath: [],
    MaxScanDepth: 3,
    OpenBrowser: false,
    DisableLAN: false,
    DefaultMode: "all",
    Username: "",
    Password: "",
    Timeout: 30,
    CertFile: "",
    KeyFile: "",
    EnableLocalCache: true,
    CachePath: "",
    ClearCacheExit: true,
    EnableUpload: true,
    UploadPath: "",
    EnableDatabase: false,
    ClearDatabaseWhenExit: false,
    ExcludePath: [],
    SupportMediaType: [],
    SupportFileType: [],
    MinImageNum: 2,
    TimeoutLimitForScan: 10,
    PrintAllPossibleQRCode: false,
    Debug: false,
    LogToFile: false,
    LogFilePath: "",
    LogFileName: "",
    ZipFileTextEncoding: "utf-8",
    EnableFrpcServer: false,
    FrpConfig: {
        FrpcCommand: "",
        ServerAddr: "",
        ServerPort: 0,
        Token: "",
        FrpType: "",
        RemotePort: 0,
        RandomRemotePort: false
    },
    GenerateMetaData: false,
}
