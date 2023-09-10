import Config from "../types/Config";
import axios from "axios";
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
            axios
                .post("/api/config_update", { [action.name]: action.value })
                .then((response) => {
                    console.log("Data sent successfully");
                    console.info(response.data); //axios默认解析Json，所以 response.data 就是解析后的object
                })
                .catch((error) => {
                    console.error("Error sending data:", error);
                });
            return { ...c, [action.name]: action.value };
        case "stringConfig":
            axios
                .post("/api/config_update", { [action.name]: action.value })
                .then((response) => {
                    console.log("Data sent successfully");
                    console.info(response.data); //axios默认解析Json，所以 response.data 就是解析后的object
                })
                .catch((error) => {
                    console.error("Error sending data:", error);
                });
            return { ...c, [action.name]: action.value };
        case "numberConfig":
            axios
                .post("/api/config_update", { [action.name]: action.value })
                .then((response) => {
                    console.log("Data sent successfully");
                    console.info(response.data); //axios默认解析Json，所以 response.data 就是解析后的object
                })
                .catch((error) => {
                    console.error("Error sending data:", error);
                });
            return { ...c, [action.name]: action.value };
        case "arrayConfig":
            axios
                .post("/api/config_update", { [action.name]: action.value })
                .then((response) => {
                    console.log("Data sent successfully");
                    console.info(response.data); //axios默认解析Json，所以 response.data 就是解析后的object
                })
                .catch((error) => {
                    console.error("Error sending data:", error);
                });
            return { ...c, [action.name]: action.value };
        default:
            console.log(action);
            throw new Error();
    }
}

export const defaultConfig: Config = {
    Port: 1234,
    Host: "",
    StoresPath: [],
    MaxScanDepth: 3,
    OpenBrowser: false,
    DisableLAN: false,
    DefaultMode: "all",
    EnableLogin: false,
    Username: "",
    Password: "",
    Timeout: 30,
    EnableTLS: false,
    CertFile: "",
    KeyFile: "",
    UseCache: true,
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
    GenerateMetaData: false,
    ConfigLocation: "",
};
