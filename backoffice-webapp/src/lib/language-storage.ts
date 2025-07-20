const LANGUAGE_STORAGE_KEY = 'i18nextLng';

export function getSavedLanguage() {
  return localStorage.getItem(LANGUAGE_STORAGE_KEY);
}

export function saveLanguage(lng: string) {
  localStorage.setItem(LANGUAGE_STORAGE_KEY, lng);
}

export function clearLanguage() {
  localStorage.removeItem(LANGUAGE_STORAGE_KEY);
}
