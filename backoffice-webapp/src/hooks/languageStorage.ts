export function useLanguageStorage() {
  const LANGUAGE_STORAGE_KEY = 'i18nextLng';

  const getSavedLanguage = () => {
    return localStorage.getItem(LANGUAGE_STORAGE_KEY);
  };

  const saveLanguage = (lng: string) => {
    localStorage.setItem(LANGUAGE_STORAGE_KEY, lng);
  };

  const clearLanguage = () => {
    localStorage.removeItem(LANGUAGE_STORAGE_KEY);
  };

  return {
    getSavedLanguage,
    saveLanguage,
    clearLanguage
  };
}
