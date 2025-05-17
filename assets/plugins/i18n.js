import i18next from 'i18next'
import LanguageDetector from 'i18next-browser-languagedetector'
import enLocale from '../locale/en_US.json'
import zhLocale from '../locale/zh_CN.json'
import jaLocale from '../locale/ja_JP.json'

i18next
    .use(LanguageDetector)
    .init({
        debug: false,
        initImmediate: true,
        supportedLngs: ['en-US', 'ja-JP', 'zh-CN', 'en', 'zh', 'ja'],
        fallbackLng: ['en', 'zh', 'ja'],
        resources: {
            'en-US': {
                translation: enLocale,
            },
            en: {
                translation: enLocale,
            },
            'zh-CN': {
                translation: zhLocale,
            },
            zh: {
                translation: zhLocale,
            },
            'ja-JP': {
                translation: jaLocale,
            },
            ja: {
                translation: jaLocale,
            },
        },
    })

window.i18next = i18next // 使i18next在全局作用域中可用 