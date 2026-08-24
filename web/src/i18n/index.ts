// i18n 语言包配置文件
import en from './en'
import zhCN from './zh-CN'

export const messages = {
  en,
  'zh-CN': zhCN
}

export const supportedLanguages = [
  { code: 'en', name: 'English' },
  { code: 'zh-CN', name: '简体中文' }
]

export default messages