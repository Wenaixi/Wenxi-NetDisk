/**
 * Wenxi网盘 - 主控制面板
 * 作者：Wenxi
 * 功能：展示用户文件统计和操作入口
 */

import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import FileUpload from './FileUpload';
import FileList from './FileList';
import SmartSearch, { SearchResults } from './SmartSearch';
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
  const [isSearching, setIsSearching] = useState(false);
  const [searchQuery, setSearchQuery] = useState(''); // Wenxi - 当前搜索词

  const fetchFiles = async () => {
    try {
      console.log('Wenxi - 开始获取文件列表...');
      const token = localStorage.getItem('token');
      
      const { getBaseURL } = await import('../utils/apiConfig');
      const response = await axios.get(`${getBaseURL()}/api/files/list`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      console.log('Wenxi - 文件列表获取成功:', response.data);
      
      setFiles(response.data);
    } catch (error) {
      console.error('Wenxi - 获取文件列表失败:', error);
      if (error.response?.status === 401) {
        console.log('Wenxi - 认证失败，需要重新登录');
        logout();
      }
    } finally {
      setLoading(false);
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

  const handleUploadSuccess = () =>> {
    setShowUpload(false);
    fetchFiles();
  };

  // Wenxi - 处理智能搜索结果
  const handleSearchResults = (results) => {
    setSearchResults(results.results || []);
    setIsSearching(false);
    setSearchQuery(results.query || '');
  };

  // Wenxi - 清除搜索
  const handleClearSearch = () => {
    setSearchResults([]);
    setIsSearching(false);
    setSearchQuery('');
    fetchFiles();
  };

  const handleDownload = async (file) => {
    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`/api/files/${file.id}/share`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        if (response.status === 401) {
          alert('登录已过期，请重新登录');
          window.location.href = '/login';
          return;
        }
        throw new Error('获取分享链接失败');
      }

      const data = await response.json();
      const shareUrl = `${window.location.origin}${data.share_url}`;
      window.location.href = shareUrl;
      
    } catch (error) {
      console.error('Wenxi - 文件下载错误:', error);
      alert('下载失败，请重试');
    }
  };

  const handleShare = async (file) => {
    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`/api/files/${file.id}/share`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        if (response.status === 401) {
          alert('登录已过期，请重新登录');
          window.location.href = '/login';
          return;
        }
        throw new Error('分享失败');
      }

      const data = await response.json();
      const shareUrl = `${window.location.origin}${data.share_url}`;
      
      try {
        await navigator.clipboard.writeText(shareUrl);
        alert('分享链接已复制到剪贴板');
      } catch (clipboardError) {
        prompt('分享链接已生成，请复制：', shareUrl);
      }
    } catch (error) {
      console.error('Wenxi - 文件分享错误:', error);
      alert('分享失败，请重试');
    }
  };

  const handleDelete = async (fileId) => {
    if (!confirm('确定要删除这个文件吗？')) return;

    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`/api/files/${fileId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (response.ok) {
        // 从搜索结果中移除
        setSearchResults(prev => prev.filter(f => f.id !== fileId));
        fetchFiles();
      } else if (response.status === 401) {
        alert('登录已过期，请重新登录');
        window.location.href = '/login';
      } else {
        alert('删除失败，请重试');
      }
    } catch (error) {
      console.error('Wenxi - 删除文件错误:', error);
      alert('删除失败，请检查网络连接');
    }
  };

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

        {/* 智能搜索功能 - Wenxi全文搜索 */}
        <div className="px-4 py-6 sm:px-0">
          <div className="bg-white shadow rounded-lg p-4 mb-6">
            <SmartSearch 
              onSearchResults={handleSearchResults}
              onClearSearch={handleClearSearch}
            />
          </div>
        </div>

        {/* 文件列表或搜索结果 */}
        <div className="px-4 py-6 sm:px-0">
          {loading ? (
            <div className="text-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
              <p className="mt-4 text-gray-600">加载中...</p>
            </div>
          ) : searchResults.length > 0 || isSearching ? (
            <>
              {searchQuery && (
                <div className="mb-4 text-sm text-gray-600">
                  搜索 "{searchQuery}" 的结果: {searchResults.length} 个文件
                </div>
              )}
              <SearchResults 
                results={searchResults}
                onDownload={handleDownload}
                onShare={handleShare}
                onDelete={handleDelete}
                formatFileSize={formatFileSize}
                searchQuery={searchQuery}
              />
            </>
          ) : (
            <FileList 
              files={files} 
              onRefresh={() => fetchFiles()}
              formatFileSize={formatFileSize}
            />
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

      {/* 底部开源信息 - Wenxi网盘开源声明 */}
      <footer className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 z-50">
        <div className="max-w-7xl mx-auto py-1 px-4 sm:px-6 lg:px-8">
          <div className="text-center text-[7px] text-gray-500 leading-none">
            <p>
              本项目已在 <a href="https://github.com/Wenaixi/Wenxi-NetDisk/" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:text-blue-800">GitHub</a> 用 <a href="https://opensource.org/licenses/MIT" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:text-blue-800">MIT协议</a> 全面开源
            </p>
            <p>
              作者：<span className="font-semibold">Wenxi</span> | 版本号：<span className="font-semibold">v1.1.2</span>
            </p>
            <p>
              本项目将会在未来不断优化改进，为您提供更好的体验
            </p>
            <p>
              联系方式：<a href="mailto:121645025@qq.com" className="text-blue-600 hover:text-blue-800">121645025@qq.com</a>
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
}