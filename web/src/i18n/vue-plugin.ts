import { createI18n } from 'vue-i18n'
import { messages, supportedLanguages } from './index'

export { supportedLanguages }

const savedLocale = localStorage.getItem('sundash-locale') || 'zh-CN'

export const i18n = createI18n({
  legacy: false,
  locale: savedLocale,
  fallbackLocale: 'zh-CN',
  messages,
})