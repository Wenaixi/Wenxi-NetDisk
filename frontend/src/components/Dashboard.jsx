/**
 * Wenxi网盘 - 主控制面板
 * 作者：Wenxi
 * 功能：展示用户文件统计和操作入口
 */

import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import FileUpload from './FileUpload';
import FileList from './FileList';
import UserSettings from './UserSettings';
import { UploadCloud, LogOut, User, HardDrive, Settings } from 'lucide-react';
import axios from 'axios';

export default function Dashboard() {
  const { user, logout, setUser } = useAuth();
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showUpload, setShowUpload] = useState(false);
  const [showSettings, setShowSettings] = useState(false);

  const [searchTerm, setSearchTerm] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [allFiles, setAllFiles] = useState([]); // Wenxi - 本地文件缓存
  const [isSearching, setIsSearching] = useState(false); // Wenxi - 搜索状态

  // Wenxi - 前端本地搜索优化
  const performLocalSearch = (query, fileList) => {
    if (!query.trim()) return fileList;
    
    const searchTerm = query.toLowerCase();
    return fileList.filter(file => 
      file.original_filename.toLowerCase().includes(searchTerm) ||
      file.filename.toLowerCase().includes(searchTerm)
    );
  };

  const fetchFiles = async (searchQuery = '') => {
    try {
      // Wenxi - 如果有本地缓存且搜索词不为空，优先使用本地搜索
      if (allFiles.length > 0 && searchQuery) {
        setIsSearching(true);
        const localResults = performLocalSearch(searchQuery, allFiles);
        setFiles(localResults);
        setIsSearching(false);
        return;
      }
      
      // Wenxi - 否则从后端获取
      console.log('Wenxi - 开始获取文件列表...');
      const response = await axios.get('/api/files/list', {
        params: searchQuery ? { search: searchQuery } : {}
      });
      console.log('Wenxi - 文件列表获取成功:', response.data);
      
      // Wenxi - 更新本地缓存
      setAllFiles(response.data);
      setFiles(response.data);
      setSearchResults(response.data);
    } catch (error) {
      console.error('Wenxi - 获取文件列表失败:', error);
      if (error.response?.status === 401) {
        console.log('Wenxi - 认证失败，需要重新登录');
        logout();
      }
    } finally {
      setLoading(false);
      setIsSearching(false);
    }
  };

  useEffect(() => {
    console.log('Wenxi - Dashboard组件挂载，用户状态:', user);
    if (user) {
      fetchFiles();
    } else {
      console.log('Wenxi - 用户未认证，跳转到登录页');
      setLoading(false);
    }
  }, [user]);

  const handleUploadSuccess = () => {
    setShowUpload(false);
    fetchFiles(searchTerm);
  };

  // Wenxi - 防抖定时器引用
  const [searchTimeout, setSearchTimeout] = useState(null);

  const handleSearch = (e) => {
    e.preventDefault();
    fetchFiles(searchTerm);
  };

  const handleSearchChange = (e) => {
    const value = e.target.value;
    setSearchTerm(value);
    
    // Wenxi - 清除之前的定时器
    if (searchTimeout) {
      clearTimeout(searchTimeout);
    }
    
    // Wenxi - 实时搜索：空值时立即获取全部文件，否则延迟300ms搜索（优化响应速度）
    if (value === '') {
      fetchFiles();
    } else {
      const timeoutId = setTimeout(() => {
        fetchFiles(value);
      }, 100); // 从500ms优化到100ms，提升响应速度
      setSearchTimeout(timeoutId);
    }
  };

  // Wenxi - 组件卸载时清除定时器
  useEffect(() => {
    return () => {
      if (searchTimeout) {
        clearTimeout(searchTimeout);
      }
    };
  }, [searchTimeout]);

  const calculateTotalSize = () => {
    return files.reduce((total, file) => total + file.file_size, 0);
  };

  const formatFileSize = (bytes) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* 顶部导航 */}
      <nav className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <UploadCloud className="h-8 w-8 text-blue-600" />
              <span className="ml-2 text-xl font-bold text-gray-900">Wenxi网盘</span>
            </div>
            
            <div className="flex items-center space-x-4">
              <button
                onClick={() => setShowSettings(true)}
                className="flex items-center text-sm text-gray-600 hover:text-gray-800 cursor-pointer"
              >
                <div className="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white font-semibold mr-2">
                  {user?.username?.charAt(0)?.toUpperCase() || 'U'}
                </div>
                <span className="mr-1">{user?.username}</span>
                <Settings className="h-4 w-4" />
              </button>
              <button
                onClick={logout}
                className="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-gray-500 hover:text-gray-700 focus:outline-none"
              >
                <LogOut className="h-4 w-4 mr-1" />
                退出
              </button>
            </div>
          </div>
        </div>
      </nav>

      {/* 主要内容 */}
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        {/* 统计卡片 */}
        <div className="px-4 py-6 sm:px-0">
          <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
            <div className="bg-white overflow-hidden shadow rounded-lg">
              <div className="p-5">
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    <UploadCloud className="h-6 w-6 text-gray-400" />
                  </div>
                  <div className="ml-5 w-0 flex-1">
                    <dl>
                      <dt className="text-sm font-medium text-gray-500 truncate">文件总数</dt>
                      <dd className="text-lg font-medium text-gray-900">{files.length}</dd>
                    </dl>
                  </div>
                </div>
              </div>
            </div>

            <div className="bg-white overflow-hidden shadow rounded-lg">
              <div className="p-5">
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    <HardDrive className="h-6 w-6 text-gray-400" />
                  </div>
                  <div className="ml-5 w-0 flex-1">
                    <dl>
                      <dt className="text-sm font-medium text-gray-500 truncate">已用空间</dt>
                      <dd className="text-lg font-medium text-gray-900">{formatFileSize(calculateTotalSize())}</dd>
                    </dl>
                  </div>
                </div>
              </div>
            </div>

            <div className="bg-white overflow-hidden shadow rounded-lg">
              <div className="p-5">
                <button
                  onClick={() => setShowUpload(true)}
                  className="w-full flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                >
                  <UploadCloud className="h-5 w-5 mr-2" /> 上传文件
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* 搜索功能 - Wenxi实时搜索优化 */}
        <div className="px-4 py-6 sm:px-0">
          <div className="bg-white shadow rounded-lg p-4 mb-6">
            <div className="flex gap-2">
              <div className="flex-1 relative">
                <input
                  type="text"
                  placeholder="搜索"
                  value={searchTerm}
                  onChange={handleSearchChange}
                  className="w-full px-4 py-2 pl-10 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
                />
                <svg className="absolute left-3 top-2.5 h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
              {searchTerm && (
                <button
                  type="button"
                  onClick={() => {
                    setSearchTerm('');
                    fetchFiles();
                  }}
                  className="px-4 py-2 text-gray-600 bg-gray-100 rounded-md hover:bg-gray-200 focus:outline-none transition-colors duration-200"
                >
                  清除
                </button>
              )}
            </div>
            {searchTerm && (
              <div className="mt-2 text-xs text-gray-500 flex items-center">
                {isSearching ? (
                  <>
                    <div className="animate-spin rounded-full h-3 w-3 border-b-2 border-blue-600 mr-2"></div>
                    正在搜索...
                  </>
                ) : (
                  <>
                    <div className="h-3 w-3 bg-green-500 rounded-full mr-2"></div>
                    本地搜索完成: {files.length} 个结果
                  </>
                )}
              </div>
            )}
          </div>
        </div>

        {/* 文件列表 */}
        <div className="px-4 py-6 sm:px-0">
          {loading ? (
            <div className="text-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
              <p className="mt-4 text-gray-600">加载中...</p>
            </div>
          ) : (
            <>
              {searchTerm && (
                <div className="mb-4 text-sm text-gray-600">
                  搜索 "{searchTerm}" 的结果: {files.length} 个文件
                </div>
              )}
              <FileList 
                files={files} 
                onRefresh={() => fetchFiles(searchTerm)}
                formatFileSize={formatFileSize}
              />
            </>
          )}
        </div>

        {/* 上传弹窗 */}
        {showUpload && (
          <FileUpload 
            onClose={() => setShowUpload(false)} 
            onSuccess={handleUploadSuccess}
          />
        )}

        {/* 用户设置弹窗 */}
        {showSettings && (
          <UserSettings
            isOpen={showSettings}
            onClose={() => setShowSettings(false)}
            currentUser={user}
            onUserUpdate={setUser}
          />
        )}
      </main>
    </div>
  );
}