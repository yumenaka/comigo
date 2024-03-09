import Config from "../types/Config";
import axios from "axios";
type Action = {
    type:
    | "init"
    | "boolean"
    | "string"
    | "number"
    | "array";
    name: string;
    value: boolean | string | number | string[];
    config: Config;
};

export function configReducer(c: Config, action: Action) {
    switch (action.type) {
        case "init":
            return { ...action.config };
        case "boolean":
            axios
                .put("/api/config", { [action.name]: action.value })
                .then((response) => {
                    console.log("Data sent successfully");
                    console.info(response.data); //axios默认解析Json，所以 response.data 就是解析后的object
                })
                .catch((error) => {
                    console.error("Error sending data:", error);
                });
            return { ...c, [action.name]: action.value };
        case "string":
            axios
                .put("/api/config", { [action.name]: action.value })
                .then((response) => {
                    console.log("Data sent successfully");
                    console.info(response.data); //axios默认解析Json，所以 response.data 就是解析后的object
                })
                .catch((error) => {
                    console.error("Error sending data:", error);
                });
            return { ...c, [action.name]: action.value };
        case "number":
            axios
                .put("/api/config", { [action.name]: action.value })
                .then((response) => {
                    console.log("Data sent successfully");
                    console.info(response.data); //axios默认解析Json，所以 response.data 就是解析后的object
                })
                .catch((error) => {
                    console.error("Error sending data:", error);
                });
            return { ...c, [action.name]: action.value };
        case "array":
            axios
                .put("/api/config", { [action.name]: action.value })
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
    LocalStores: [],
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
    UseCache: false,
    CachePath: "",
    ClearCacheExit: true,
    EnableUpload: true,
    UploadPath: "",
    EnableDatabase: false,
    ClearDatabaseWhenExit: false,
    ExcludePath: [],
    SupportMediaType: [],
    SupportFileType: [],
    MinImageNum: 3,
    TimeoutLimitForScan: 10,
    PrintAllPossibleQRCode: false,
    Debug: false,
    LogToFile: false,
    LogFilePath: "",
    LogFileName: "",
    ZipFileTextEncoding: "utf-8",
    GenerateMetaData: false,
    ConfigSaveTo: "RAM",
};
