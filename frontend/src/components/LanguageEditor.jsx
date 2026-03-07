import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { X, Download, Upload, RotateCcw, Globe, Edit3, Save, Search } from 'lucide-react';
import { changeLanguage, supportedLanguages, getCurrentLanguage } from '../i18n';
import zhCN from '../i18n/locales/zh-CN.json';
import en from '../i18n/locales/en.json';
import zhTW from '../i18n/locales/zh-TW.json';

const defaultTranslations = {
  'zh-CN': zhCN,
  'en': en,
  'zh-TW': zhTW
};

export default function LanguageEditor({ isOpen, onClose }) {
  const { t } = useTranslation();
  const [selectedLang, setSelectedLang] = useState(getCurrentLanguage());
  const [customTranslations, setCustomTranslations] = useState({});
  const [searchTerm, setSearchTerm] = useState('');
  const [activeTab, setActiveTab] = useState('app');
  const [hasChanges, setHasChanges] = useState(false);
  const [saveSuccess, setSaveSuccess] = useState(false);

  useEffect(() => {
    if (isOpen) {
      loadCustomTranslations();
    }
  }, [isOpen]);

  const loadCustomTranslations = () => {
    try {
      const saved = localStorage.getItem('wenxi-custom-language');
      if (saved) {
        setCustomTranslations(JSON.parse(saved));
      }
    } catch (e) {
      console.error('Failed to load custom translations:', e);
    }
  };

  const getFlattenedKeys = (obj, prefix = '') => {
    const keys = [];
    for (const key in obj) {
      if (typeof obj[key] === 'object' && obj[key] !== null) {
        keys.push(...getFlattenedKeys(obj[key], `${prefix}${key}.`));
      } else {
        keys.push(`${prefix}${key}`);
      }
    }
    return keys;
  };

  const getValueByPath = (obj, path) => {
    return path.split('.').reduce((acc, part) => acc?.[part], obj);
  };

  const setValueByPath = (obj, path, value) => {
    const parts = path.split('.');
    const last = parts.pop();
    const target = parts.reduce((acc, part) => {
      if (!acc[part]) acc[part] = {};
      return acc[part];
    }, obj);
    target[last] = value;
    return { ...obj };
  };

  const handleTranslationChange = (key, value) => {
    const newCustom = { ...customTranslations };
    if (!newCustom[selectedLang]) {
      newCustom[selectedLang] = {};
    }
    newCustom[selectedLang] = setValueByPath(newCustom[selectedLang], key, value);
    setCustomTranslations(newCustom);
    setHasChanges(true);
    setSaveSuccess(false);
  };

  const handleSave = () => {
    try {
      localStorage.setItem('wenxi-custom-language', JSON.stringify(customTranslations));
      setHasChanges(false);
      setSaveSuccess(true);
      setTimeout(() => setSaveSuccess(false), 2000);
      window.location.reload();
    } catch (e) {
      console.error('Failed to save custom translations:', e);
      alert(t('common.error'));
    }
  };

  const handleExport = () => {
    const dataStr = JSON.stringify(customTranslations, null, 2);
    const dataBlob = new Blob([dataStr], { type: 'application/json' });
    const url = URL.createObjectURL(dataBlob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `wenxi-custom-translations-${new Date().toISOString().split('T')[0]}.json`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  };

  const handleImport = (event) => {
    const file = event.target.files[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const imported = JSON.parse(e.target.result);
        setCustomTranslations(imported);
        setHasChanges(true);
        alert(t('language.importSuccess'));
      } catch (err) {
        alert(t('language.importError'));
      }
    };
    reader.readAsText(file);
    event.target.value = '';
  };

  const handleReset = () => {
    if (confirm(t('language.reset') + '?')) {
      localStorage.removeItem('wenxi-custom-language');
      setCustomTranslations({});
      setHasChanges(false);
      window.location.reload();
    }
  };

  const getCurrentValue = (key) => {
    const customValue = getValueByPath(customTranslations[selectedLang], key);
    if (customValue !== undefined) return customValue;
    return getValueByPath(defaultTranslations[selectedLang], key) || '';
  };

  const allKeys = getFlattenedKeys(defaultTranslations[selectedLang]);
  const filteredKeys = allKeys.filter(key => 
    key.toLowerCase().includes(searchTerm.toLowerCase()) ||
    getCurrentValue(key).toLowerCase().includes(searchTerm.toLowerCase())
  );

  const groupedKeys = filteredKeys.reduce((acc, key) => {
    const group = key.split('.')[0];
    if (!acc[group]) acc[group] = [];
    acc[group].push(key);
    return acc;
  }, {});

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg w-full max-w-4xl mx-4 max-h-[90vh] flex flex-col">
        <div className="flex justify-between items-center p-4 border-b">
          <div className="flex items-center space-x-4">
            <h2 className="text-xl font-bold text-gray-900 flex items-center">
              <Globe className="h-5 w-5 mr-2" />
              {t('language.title')}
            </h2>
          </div>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600">
            <X className="h-5 w-5" />
          </button>
        </div>

        <div className="flex items-center justify-between p-4 border-b bg-gray-50">
          <div className="flex items-center space-x-4">
            <label className="text-sm font-medium text-gray-700">{t('language.select')}:</label>
            <select
              value={selectedLang}
              onChange={(e) => setSelectedLang(e.target.value)}
              className="px-3 py-1 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500"
            >
              {supportedLanguages.map(lang => (
                <option key={lang.code} value={lang.code}>
                  {lang.nativeName}
                </option>
              ))}
            </select>
          </div>

          <div className="flex items-center space-x-2">
            <button
              onClick={handleExport}
              className="flex items-center px-3 py-1.5 text-sm text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
            >
              <Download className="h-4 w-4 mr-1" />
              {t('language.export')}
            </button>
            <label className="flex items-center px-3 py-1.5 text-sm text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 cursor-pointer">
              <Upload className="h-4 w-4 mr-1" />
              {t('language.import')}
              <input type="file" accept=".json" onChange={handleImport} className="hidden" />
            </label>
            <button
              onClick={handleReset}
              className="flex items-center px-3 py-1.5 text-sm text-red-600 bg-white border border-red-300 rounded-md hover:bg-red-50"
            >
              <RotateCcw className="h-4 w-4 mr-1" />
              {t('language.reset')}
            </button>
          </div>
        </div>

        <div className="p-4 border-b">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
            <input
              type="text"
              placeholder="Search translation keys or values..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        <div className="flex-1 overflow-hidden flex">
          <div className="w-48 border-r overflow-y-auto bg-gray-50">
            {Object.keys(groupedKeys).map(group => (
              <button
                key={group}
                onClick={() => setActiveTab(group)}
                className={`w-full text-left px-4 py-2 text-sm capitalize ${
                  activeTab === group 
                    ? 'bg-blue-100 text-blue-700 font-medium' 
                    : 'text-gray-700 hover:bg-gray-100'
                }`}
              >
                {group}
                <span className="ml-2 text-xs text-gray-500">
                  ({groupedKeys[group].length})
                </span>
              </button>
            ))}
          </div>

          <div className="flex-1 overflow-y-auto p-4">
            {groupedKeys[activeTab]?.map(key => (
              <div key={key} className="mb-4 p-3 bg-gray-50 rounded-lg">
                <label className="block text-xs font-medium text-gray-600 mb-1 font-mono">
                  {key}
                </label>
                <input
                  type="text"
                  value={getCurrentValue(key)}
                  onChange={(e) => handleTranslationChange(key, e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500"
                />
              </div>
            ))}
          </div>
        </div>

        <div className="p-4 border-t bg-gray-50 flex justify-between items-center">
          <div>
            {saveSuccess && (
              <span className="text-sm text-green-600">{t('common.success')}!</span>
            )}
            {hasChanges && !saveSuccess && (
              <span className="text-sm text-orange-600">{t('common.edit')}...</span>
            )}
          </div>
          <div className="flex space-x-2">
            <button
              onClick={onClose}
              className="px-4 py-2 text-sm text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
            >
              {t('common.cancel')}
            </button>
            <button
              onClick={handleSave}
              disabled={!hasChanges}
              className="flex items-center px-4 py-2 text-sm text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              <Save className="h-4 w-4 mr-1" />
              {t('common.save')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
