/**
 * Wenxi网盘 - 主题切换按钮组件
 * 作者：Wenxi
 * 功能：提供美观的主题切换按钮，支持太阳/月亮动画
 */

import React from 'react';
import { Sun, Moon, Monitor } from 'lucide-react';
import { useTheme } from '../contexts/ThemeContext';

export default function ThemeToggle({ variant = 'button', size = 'md', className = '' }) {
  const { theme, isDark, toggleTheme, setLightTheme, setDarkTheme, followSystemTheme } = useTheme();

  const sizeClasses = {
    sm: 'w-8 h-8',
    md: 'w-10 h-10',
    lg: 'w-12 h-12',
  };

  const iconSizes = {
    sm: 16,
    md: 20,
    lg: 24,
  };

  // 简约按钮模式 - 只显示切换按钮
  if (variant === 'button') {
    return (
      <button
        onClick={toggleTheme}
        className={`
          relative inline-flex items-center justify-center rounded-lg
          transition-all duration-300 ease-out
          ${sizeClasses[size]}
          ${isDark 
            ? 'bg-slate-800 hover:bg-slate-700 text-amber-400 shadow-lg shadow-amber-500/20' 
            : 'bg-white hover:bg-gray-100 text-blue-600 shadow-md shadow-blue-500/10 border border-gray-200'
          }
          hover:scale-110 active:scale-95
          ${className}
        `}
        title={isDark ? '切换到浅色模式' : '切换到深色模式'}
        aria-label={isDark ? '切换到浅色模式' : '切换到深色模式'}
      >
        <div className="relative">
          {/* 太阳图标 */}
          <Sun 
            size={iconSizes[size]} 
            className={`
              absolute inset-0 transition-all duration-500
              ${isDark ? 'rotate-90 opacity-0 scale-0' : 'rotate-0 opacity-100 scale-100'}
            `}
          />
          {/* 月亮图标 */}
          <Moon 
            size={iconSizes[size]} 
            className={`
              transition-all duration-500
              ${isDark ? 'rotate-0 opacity-100 scale-100' : '-rotate-90 opacity-0 scale-0'}
            `}
          />
        </div>
      </button>
    );
  }

  // 下拉菜单模式 - 显示三个选项
  if (variant === 'dropdown') {
    return (
      <div className={`relative group ${className}`}>
        <button
          className={`
            inline-flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium
            transition-all duration-200
            ${isDark 
              ? 'bg-slate-800 text-slate-200 hover:bg-slate-700' 
              : 'bg-white text-gray-700 hover:bg-gray-100 border border-gray-200'
            }
          `}
        >
          {isDark ? <Moon size={16} /> : <Sun size={16} />}
          <span>{isDark ? '深色模式' : '浅色模式'}</span>
        </button>
        
        {/* 下拉菜单 */}
        <div className="
          absolute right-0 mt-2 w-40 py-2 rounded-lg shadow-xl
          bg-white dark:bg-slate-800 
          border border-gray-200 dark:border-slate-700
          opacity-0 invisible group-hover:opacity-100 group-hover:visible
          transition-all duration-200 z-50
        ">
          <button
            onClick={setLightTheme}
            className={`
              w-full px-4 py-2 text-left text-sm flex items-center gap-2
              transition-colors duration-150
              ${theme === 'light' 
                ? 'bg-blue-50 dark:bg-slate-700 text-blue-600 dark:text-blue-400' 
                : 'text-gray-700 dark:text-slate-300 hover:bg-gray-100 dark:hover:bg-slate-700'
              }
            `}
          >
            <Sun size={16} />
            浅色模式
          </button>
          <button
            onClick={setDarkTheme}
            className={`
              w-full px-4 py-2 text-left text-sm flex items-center gap-2
              transition-colors duration-150
              ${theme === 'dark' 
                ? 'bg-blue-50 dark:bg-slate-700 text-blue-600 dark:text-blue-400' 
                : 'text-gray-700 dark:text-slate-300 hover:bg-gray-100 dark:hover:bg-slate-700'
              }
            `}
          >
            <Moon size={16} />
            深色模式
          </button>
          <button
            onClick={followSystemTheme}
            className={`
              w-full px-4 py-2 text-left text-sm flex items-center gap-2
              transition-colors duration-150
              ${!localStorage.getItem('wenxi-theme')
                ? 'bg-blue-50 dark:bg-slate-700 text-blue-600 dark:text-blue-400' 
                : 'text-gray-700 dark:text-slate-300 hover:bg-gray-100 dark:hover:bg-slate-700'
              }
            `}
          >
            <Monitor size={16} />
            跟随系统
          </button>
        </div>
      </div>
    );
  }

  // 分段控制模式
  if (variant === 'segmented') {
    return (
      <div className={`
        inline-flex rounded-lg p-1 gap-1
        ${isDark ? 'bg-slate-800' : 'bg-gray-100'}
        ${className}
      `}>
        <button
          onClick={setLightTheme}
          className={`
            flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm font-medium
            transition-all duration-200
            ${theme === 'light'
              ? 'bg-white dark:bg-slate-600 text-blue-600 dark:text-blue-400 shadow-sm'
              : 'text-gray-600 dark:text-slate-400 hover:text-gray-900 dark:hover:text-slate-200'
            }
          `}
        >
          <Sun size={14} />
          浅色
        </button>
        <button
          onClick={setDarkTheme}
          className={`
            flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm font-medium
            transition-all duration-200
            ${theme === 'dark'
              ? 'bg-white dark:bg-slate-600 text-blue-600 dark:text-blue-400 shadow-sm'
              : 'text-gray-600 dark:text-slate-400 hover:text-gray-900 dark:hover:text-slate-200'
            }
          `}
        >
          <Moon size={14} />
          深色
        </button>
        <button
          onClick={followSystemTheme}
          className={`
            flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm font-medium
            transition-all duration-200
            ${!localStorage.getItem('wenxi-theme')
              ? 'bg-white dark:bg-slate-600 text-blue-600 dark:text-blue-400 shadow-sm'
              : 'text-gray-600 dark:text-slate-400 hover:text-gray-900 dark:hover:text-slate-200'
            }
          `}
        >
          <Monitor size={14} />
          系统
        </button>
      </div>
    );
  }

  return null;
}
