/**
 * Wenxi网盘 - 主题上下文
 * 作者：Wenxi
 * 功能：提供全局主题状态管理，支持暗黑模式切换
 */

import React, { createContext, useContext, useState, useEffect } from 'react';

const ThemeContext = createContext();

export function ThemeProvider({ children }) {
  const [theme, setTheme] = useState('light');
  const [mounted, setMounted] = useState(false);

  // 初始化主题：优先从localStorage获取，其次跟随系统
  useEffect(() => {
    const savedTheme = localStorage.getItem('wenxi-theme');
    
    if (savedTheme) {
      setTheme(savedTheme);
    } else {
      // 跟随系统主题
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      setTheme(prefersDark ? 'dark' : 'light');
    }
    
    setMounted(true);
  }, []);

  // 监听系统主题变化
  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    
    const handleChange = (e) => {
      const savedTheme = localStorage.getItem('wenxi-theme');
      // 只有在用户没有手动设置主题时才跟随系统
      if (!savedTheme) {
        setTheme(e.matches ? 'dark' : 'light');
      }
    };

    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, []);

  // 应用主题到document
  useEffect(() => {
    if (!mounted) return;

    const root = window.document.documentElement;
    
    if (theme === 'dark') {
      root.classList.add('dark');
    } else {
      root.classList.remove('dark');
    }

    // 保存到localStorage
    localStorage.setItem('wenxi-theme', theme);
  }, [theme, mounted]);

  const toggleTheme = () => {
    setTheme((prev) => (prev === 'light' ? 'dark' : 'light'));
  };

  const setLightTheme = () => setTheme('light');
  const setDarkTheme = () => setTheme('dark');
  const followSystemTheme = () => {
    localStorage.removeItem('wenxi-theme');
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    setTheme(prefersDark ? 'dark' : 'light');
  };

  const value = {
    theme,
    isDark: theme === 'dark',
    toggleTheme,
    setLightTheme,
    setDarkTheme,
    followSystemTheme,
    mounted,
  };

  return (
    <ThemeContext.Provider value={value}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (context === undefined) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
}

export default ThemeContext;
