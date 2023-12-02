import ConfigStatus from '../types/ConfigStatus';
import axios from "axios";
type ActionConfigStatus = {
    type:
    | "init"
    | "boolean"
    | "string";
    name: string;
    value: boolean | string 
    config: ConfigStatus;
};

export function configStatusReducer(c: ConfigStatus, action: ActionConfigStatus) {
    switch (action.type) {
        case "init":
            //console.log("config_status", action.config );
            return { ...action.config };
        case "boolean":
            axios
                .put("/api/config/status", { [action.name]: action.value })
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
        default:
            console.log(action);
            throw new Error();
    }
}

export const defaultConfigStatus: ConfigStatus = {
    CurrentConfig: "RAM",
    Home: false,
    Execution: false,
    Program: false,
};
