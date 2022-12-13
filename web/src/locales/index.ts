// 参考：https://www.i4k.xyz/article/weixin_42174938/119764352
// import {createI18n} from 'vue-i18n'
import { createI18n } from "vue-i18n";
import enLocale from "./en.json";
import cnLocale from "./zh_CN.json";
import jaLocale from "./ja.json";
// const messages = {
//   en: {
//     ...enLocale
//   },
//   zh: {
//     ...cnLocale
//   },
//   ja: {
//     ...jaLocale
//   },
// }

// Type-define 'en-US' as the master schema for the resource
type MessageSchema = typeof cnLocale;

const i18n = createI18n<[MessageSchema], "en" | "ja" | "zh">({
  locale: localStorage.getItem("lang") || "zh",
  globalInjection: true,
  messages: {
    en: enLocale,
    ja: jaLocale,
    zh: cnLocale,
  },
});
export default i18n;
