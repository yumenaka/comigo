import 'htmx.org'
import Alpine from 'alpinejs'
import persist from '@alpinejs/persist'
import i18next from 'i18next'
import 'flowbite'
import LanguageDetector from 'i18next-browser-languagedetector'
import 'tw-colors'
// 将 Alpine 实例添加到窗口对象中。
window.Alpine = Alpine


// Alpine Persist 插件，用于持久化存储。默认存储到 localStorage。
// 详细用法参见： https://alpinejs.dev/plugins/persist
Alpine.plugin(persist)

// i18next 国际化插件，用于国际化。详细用法参见：
// https://www.i18next.com/overview/getting-started
import enLocale from './locales/en_US.json'
import zhLocale from './locales/zh_CN.json'
import jaLocale from './locales/ja_JP.json'
i18next
  .use(LanguageDetector)
  .init({
    debug: false,
    // // 在 setTimeout（默认异步行为）内的 init（） 中触发资源加载。如果您的后端同步加载资源，请将其设置为 false - 这样，
    // // 可以在 init（） 之后调用 i18next.t（） 而无需依赖初始化回调。此选项仅适用于同步（阻塞）加载后端，例如 i18next-fs-backend 和 i18next-sync-fs-backend！
    initImmediate: false,
    //lng: 'en', // if you're using a language detector, do not define the lng option
    // supportedLngs: ['en', 'cn', 'ja'],
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
  .then(function (t) {
    //console.log(t('test'))
    // i18next.changeLanguage('en', (err, t) => {
    //     if (err) return console.log('something went wrong loading', err);
    //     console.log(t('test'));
    // });
  })

window.i18next = i18next; // 使i18next在全局作用域中可用

import screenfull from 'screenfull';
document.getElementById('FullScreenIcon').addEventListener('click', () => {
    if (screenfull.isEnabled) {
        screenfull.toggle();
    } else {
        // Ignore or do something else
        i18next.t('not_support_fullscreen');
    }
});


// Start Alpine.
Alpine.start()
