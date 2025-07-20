import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import enTranslation from './locales/en/translation.json';
import thTranslation from './locales/th/translation.json';
import { getSavedLanguage } from './lib/language-storage';

i18n.use(initReactI18next).init({
  resources: {
    en: { translation: enTranslation },
    th: { translation: thTranslation },
  },
  lng: getSavedLanguage() || 'th',
  fallbackLng: 'th',
  interpolation: {
    escapeValue: false,
  },
});

export default i18n;
