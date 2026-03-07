/**
 * Wenxi网盘 - 智能搜索组件
 * 作者：Wenxi
 * 功能：提供全文搜索、文件类型过滤、搜索结果高亮显示
 */

import React, { useState, useEffect, useCallback } from 'react';
import { Search, Filter, X, FileText, FileImage, FileVideo, FileAudio, FileCode, FileArchive, FileSpreadsheet, FilePresentation, Eye } from 'lucide-react';
import { isPreviewable } from './MediaPreview';

// 文件类型图标映射
const fileTypeIcons = {
  text: FileText,
  document: FileText,
  spreadsheet: FileSpreadsheet,
  presentation: FilePresentation,
  image: FileImage,
  video: FileVideo,
  audio: FileAudio,
  archive: FileArchive,
  code: FileCode,
  other: FileText,
};

// 文件类型标签映射
const fileTypeLabels = {
  text: '文本文档',
  document: '办公文档',
  spreadsheet: '电子表格',
  presentation: '演示文稿',
  image: '图片',
  video: '视频',
  audio: '音频',
  archive: '压缩文件',
  code: '代码文件',
  other: '其他',
};

export default function SmartSearch({ onSearchResults, onClearSearch }) {
  const [query, setQuery] = useState('');
  const [fileType, setFileType] = useState('');
  const [isSearching, setIsSearching] = useState(false);
  const [fileTypes, setFileTypes] = useState([]);
  const [showFilters, setShowFilters] = useState(false);
  const [lastSearch, setLastSearch] = useState(null);
  const [searchStats, setSearchStats] = useState(null);

  // 获取文件类型列表
  useEffect(() => {
    fetchFileTypes();
    fetchSearchStats();
  }, []);

  const fetchFileTypes = async () => {
    try {
      const token = localStorage.getItem('token');
      const response = await fetch('/api/search/file-types', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      if (response.ok) {
        const data = await response.json();
        setFileTypes(data);
      }
    } catch (error) {
      console.error('Wenxi - 获取文件类型失败:', error);
    }
  };

  const fetchSearchStats = async () => {
    try {
      const token = localStorage.getItem('token');
      const response = await fetch('/api/search/stats', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      if (response.ok) {
        const data = await response.json();
        setSearchStats(data);
      }
    } catch (error) {
      console.error('Wenxi - 获取搜索统计失败:', error);
    }
  };

  // 防抖搜索
  const debouncedSearch = useCallback(
    debounce(async (searchQuery, searchFileType) => {
      if (!searchQuery.trim()) {
        onClearSearch();
        return;
      }

      setIsSearching(true);
      
      try {
        const token = localStorage.getItem('token');
        const params = new URLSearchParams({
          query: searchQuery,
          limit: '20',
          offset: '0'
        });
        
        if (searchFileType) {
          params.append('file_type', searchFileType);
        }
        
        const response = await fetch(`/api/search?${params}`, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        if (!response.ok) {
          if (response.status === 401) {
            alert('登录已过期，请重新登录');
            window.location.href = '/login';
            return;
          }
          throw new Error('搜索失败');
        }

        const data = await response.json();
        setLastSearch(data);
        onSearchResults(data);
        
        // 记录搜索性能
        if (data.search_time_ms > 200) {
          console.warn(`Wenxi - 搜索响应较慢: ${data.search_time_ms}ms`);
        } else {
          console.log(`Wenxi - 搜索完成: ${data.search_time_ms}ms, 找到 ${data.total} 个结果`);
        }
        
      } catch (error) {
        console.error('Wenxi - 搜索错误:', error);
      } finally {
        setIsSearching(false);
      }
    }, 300),
    [onSearchResults, onClearSearch]
  );

  // 监听查询变化
  useEffect(() => {
    debouncedSearch(query, fileType);
  }, [query, fileType, debouncedSearch]);

  const handleClear = () => {
    setQuery('');
    setFileType('');
    onClearSearch();
  };

  const highlightText = (text, highlightQuery) => {
    if (!highlightQuery || !text) return text;
    
    const parts = text.split(new RegExp(`(${highlightQuery})`, 'gi'));
    return parts.map((part, i) => 
      part.toLowerCase() === highlightQuery.toLowerCase() ? 
        <mark key={i} className="bg-yellow-200 px-1 rounded">{part}</mark> : part
    );
  };

  const getFileIcon = (fileType) => {
    const Icon = fileTypeIcons[fileType] || fileTypeIcons.other;
    return <Icon className="h-5 w-5" />;
  };

  return (
    <div className="space-y-4">
      {/* 搜索栏 */}
      <div className="relative">
        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <Search className={`h-5 w-5 ${isSearching ? 'animate-pulse text-blue-500' : 'text-gray-400'}`} />
        </div>
        
        <input
          type="text"
          className="block w-full pl-10 pr-20 py-3 border border-gray-300 rounded-lg leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm transition-all"
          placeholder="搜索文件名、描述..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        
        <div className="absolute inset-y-0 right-0 flex items-center pr-2 space-x-1">
          {query && (
            <button
              onClick={handleClear}
              className="p-1 text-gray-400 hover:text-gray-600 rounded-full hover:bg-gray-100"
            >
              <X className="h-4 w-4" />
            </button>
          )}
          
          <button
            onClick={() => setShowFilters(!showFilters)}
            className={`p-2 rounded-md transition-colors ${
              showFilters || fileType ? 'text-blue-600 bg-blue-50' : 'text-gray-400 hover:text-gray-600'
            }`}
            title="过滤器"
          >
            <Filter className="h-4 w-4" />
          </button>
        </div>
      </div>

      {/* 过滤器面板 */}
      {showFilters && (
        <div className="bg-gray-50 rounded-lg p-4 border border-gray-200">
          <div className="flex items-center justify-between mb-3">
            <h4 className="text-sm font-medium text-gray-700">文件类型过滤</h4>
            {fileType && (
              <button
                onClick={() => setFileType('')}
                className="text-xs text-blue-600 hover:text-blue-800"
              >
                清除筛选
              </button>
            )}
          </div>
          
          <div className="flex flex-wrap gap-2">
            {fileTypes.map((type) => (
              <button
                key={type.value}
                onClick={() => setFileType(fileType === type.value ? '' : type.value)}
                className={`inline-flex items-center px-3 py-1.5 rounded-full text-xs font-medium transition-colors ${
                  fileType === type.value
                    ? 'bg-blue-100 text-blue-800 border border-blue-300'
                    : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50'
                }`}
              >
                {React.createElement(fileTypeIcons[type.value] || FileText, {
                  className: "h-3 w-3 mr-1"
                })}
                {type.label}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* 搜索统计 */}
      {searchStats && (
        <div className="flex items-center justify-between text-xs text-gray-500">
          <span>已索引 {searchStats.total_documents} 个文件</span>
          {lastSearch && query && (
            <span>
              找到 {lastSearch.total} 个结果 
              <span className={`ml-1 ${lastSearch.search_time_ms > 200 ? 'text-red-500' : 'text-green-500'}`}>
                ({lastSearch.search_time_ms}ms)
              </span>
            </span>
          )}
        </div>
      )}
    </div>
  );
}

// 防抖函数
function debounce(func, wait) {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}

// 搜索结果展示组件
export function SearchResults({ results, onDownload, onShare, onDelete, onPreview, formatFileSize, searchQuery }) {
  const highlightText = (text) => {
    if (!searchQuery || !text) return text;
    
    const parts = text.split(new RegExp(`(${searchQuery})`, 'gi'));
    return parts.map((part, i) => 
      part.toLowerCase() === searchQuery.toLowerCase() ? 
        <mark key={i} className="bg-yellow-200 px-1 rounded">{part}</mark> : part
    );
  };

  const getFileIcon = (fileType) => {
    const Icon = fileTypeIcons[fileType] || fileTypeIcons.other;
    return Icon;
  };

  if (results.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="mx-auto h-12 w-12 text-gray-400">
          <Search className="h-12 w-12" />
        </div>
        <h3 className="mt-2 text-sm font-medium text-gray-900">没有找到匹配的文件</h3>
        <p className="mt-1 text-sm text-gray-500">尝试使用不同的关键词</p>
      </div>
    );
  }

  return (
    <div className="bg-white shadow overflow-hidden sm:rounded-md">
      <ul className="divide-y divide-gray-200">
        {results.map((file, index) => {
          const Icon = getFileIcon(file.file_type);
          const previewInfo = isPreviewable(file);
          return (
            <li key={file.id} className="px-4 py-4 sm:px-6 hover:bg-gray-50 transition-colors">
              <div className="flex items-center justify-between">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      <div className={`h-10 w-10 rounded-lg flex items-center justify-center ${
                        file.file_type === 'image' ? 'bg-purple-100 text-purple-600' :
                        file.file_type === 'video' ? 'bg-red-100 text-red-600' :
                        file.file_type === 'audio' ? 'bg-green-100 text-green-600' :
                        file.file_type === 'code' ? 'bg-gray-100 text-gray-600' :
                        'bg-blue-100 text-blue-600'
                      }`}>
                        <Icon className="h-5 w-5" />
                      </div>
                    </div>
                    <div className="ml-4 flex-1 min-w-0">
                      <button
                        onClick={() => previewInfo.canPreview && onPreview && onPreview(file, index)}
                        className={`text-sm font-medium truncate text-left transition-colors ${
                          previewInfo.canPreview
                            ? 'text-blue-600 hover:text-blue-800 cursor-pointer'
                            : 'text-gray-900 cursor-default'
                        }`}
                        title={previewInfo.canPreview ? '点击预览' : file.original_filename}
                      >
                        {highlightText(file.original_filename)}
                      </button>
                      <div className="flex items-center mt-1 space-x-3 text-xs text-gray-500">
                        <span>{formatFileSize(file.file_size)}</span>
                        <span>•</span>
                        <span>{new Date(file.created_at).toLocaleString('zh-CN')}</span>
                        {file.file_type && (
                          <>
                            <span>•</span>
                            <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">
                              {fileTypeLabels[file.file_type] || file.file_type}
                            </span>
                          </>
                        )}
                        {previewInfo.canPreview && (
                          <>
                            <span>•</span>
                            <span className="inline-flex items-center text-xs text-blue-500 font-medium">
                              <Eye className="h-3 w-3 mr-0.5" />
                              可预览
                            </span>
                          </>
                        )}
                      </div>
                      {file.description && (
                        <p className="mt-1 text-xs text-gray-500 truncate">
                          {highlightText(file.description)}
                        </p>
                      )}
                      {file.highlights?.content && (
                        <p className="mt-1 text-xs text-gray-600 bg-gray-50 p-2 rounded">
                          {highlightText(file.highlights.content)}
                        </p>
                      )}
                    </div>
                  </div>
                </div>
                <div className="ml-4 flex-shrink-0 flex space-x-2">
                  <button
                    onClick={() => onDownload(file)}
                    className="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
                    title="下载"
                  >
                    <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                    </svg>
                  </button>
                  <button
                    onClick={() => onShare(file)}
                    className="p-2 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-md transition-colors"
                    title="分享"
                  >
                    <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
                    </svg>
                  </button>
                  <button
                    onClick={() => onDelete(file.id)}
                    className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors"
                    title="删除"
                  >
                    <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </div>
            </li>
          );
        })}
      </ul>
    </div>
  );
}