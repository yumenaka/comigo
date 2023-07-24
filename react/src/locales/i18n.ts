import i18n from "i18next";
import { initReactI18next } from "react-i18next";

import Backend from "i18next-http-backend";
import LanguageDetector from "i18next-browser-languagedetector";

import enUS from "./en.json";
import zhCN from "./zh_CN.json";
import jaJP from "./ja.json";

const resources = {
  en: {
    translation: enUS,
  },
  zh: {
    translation: zhCN,
  },
  ja: {
    translation: jaJP,
  },
};

void i18n
  // load translation using http -> see /public/locales (i.e. https://github.com/i18next/react-i18next/tree/master/example/react/public/locales)
  // learn more: https://github.com/i18next/i18next-http-backend
  .use(Backend)
  // detect user language
  // learn more: https://github.com/i18next/i18next-browser-languageDetector
  .use(LanguageDetector)
  // pass the i18n instance to react-i18next.
  .use(initReactI18next)
  // init i18next
  // for all options read: https://www.i18next.com/overview/configuration-options
  .init({
    fallbackLng: "en",
    lng: "zh",
    debug: true,
    resources: resources,
    interpolation: {
      escapeValue: false, // 由于React默认情况下会转义，因此不需要
    },
  });

export default i18n;
