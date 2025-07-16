import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import enTranslation from './locales/en/translation.json';
import thTranslation from './locales/th/translation.json';
import { useLanguageStorage } from './hooks/languageStorage';

const { getSavedLanguage } = useLanguageStorage();

i18n.use(initReactI18next).init({
  resources: {
    en: {
      translation: enTranslation,
    },
    th: {
      translation: thTranslation,
    },
  },
  lng: getSavedLanguage() || 'th',
  fallbackLng: 'en',
  interpolation: {
    escapeValue: false,
  },
});

export default i18n;
