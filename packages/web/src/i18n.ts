import { createI18n } from 'vue-i18n'
import en from '@/i18n/locales/en.json'
import zh from '@/i18n/locales/zh.json'
import { computed } from 'vue'

type enMessageSchema = typeof en
type zhMessageSchema = typeof zh;


const i18n = createI18n<[enMessageSchema, zhMessageSchema], 'en' | 'zh'>({
  locale: 'zh',
  legacy: false,
  fallbackLocale: 'en',
  messages: {
    en,
    zh
  }
})


export default i18n

const t = i18n.global.t

export const i18nRef = (arg:string) => {
  return computed(() => {
    return t(arg)
  })
}
