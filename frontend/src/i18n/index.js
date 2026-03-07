import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import zhCN from './locales/zh-CN.json';
import en from './locales/en.json';
import zhTW from './locales/zh-TW.json';

const getCustomTranslations = () => {
  try {
    const customLang = localStorage.getItem('wenxi-custom-language');
    if (customLang) {
      return JSON.parse(customLang);
    }
  } catch (e) {
    console.error('Failed to load custom translations:', e);
  }
  return null;
};

const customTranslations = getCustomTranslations();

const resources = {
  'zh-CN': {
    translation: customTranslations?.['zh-CN'] || zhCN
  },
  'en': {
    translation: customTranslations?.['en'] || en
  },
  'zh-TW': {
    translation: customTranslations?.['zh-TW'] || zhTW
  }
};

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'zh-CN',
    debug: false,
    
    detection: {
      order: ['localStorage', 'navigator', 'htmlTag'],
      caches: ['localStorage'],
      lookupLocalStorage: 'wenxi-language',
    },

    interpolation: {
      escapeValue: false,
    },

    react: {
      useSuspense: false,
    },
  });

export const changeLanguage = (lng) => {
  i18n.changeLanguage(lng);
  localStorage.setItem('wenxi-language', lng);
};

export const getCurrentLanguage = () => {
  return i18n.language || 'zh-CN';
};

export const supportedLanguages = [
  { code: 'zh-CN', name: 'Simplified Chinese', nativeName: '简体中文', flag: 'CN' },
  { code: 'en', name: 'English', nativeName: 'English', flag: 'US' },
  { code: 'zh-TW', name: 'Traditional Chinese', nativeName: '繁體中文', flag: 'TW' },
];

export const isCustomLanguage = (lng) => {
  return lng?.startsWith('custom-');
};

export const getCustomLanguageName = () => {
  const customName = localStorage.getItem('wenxi-custom-language-name');
  return customName || 'Custom Language';
};

export default i18n;
