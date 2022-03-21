// 参考：https://www.i4k.xyz/article/weixin_42174938/119764352
// import {createI18n} from 'vue-i18n'
import { createI18n } from 'vue-i18n/index'
import enLocale from './en.json'
import cnLocale from './zh_CN.json'
import jaLocale from './ja.json'
const messages = {
  en: {
    ...enLocale
  },
  zh: {
    ...cnLocale
  },
  ja: {
    ...jaLocale
  },
}
const i18n = new createI18n({
  locale: localStorage.getItem('lang') || 'zh',
  globalInjection: true,
  messages
})
export default i18n