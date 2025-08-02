/**
 * Wenxi网盘 - 文件列表组件
 * 作者：Wenxi
 * 功能：展示用户文件列表，支持排序、搜索、批量操作
 */

import React, { useState } from 'react';
import { Download, Trash2, Share2, Copy, Check } from 'lucide-react';

export default function FileList({ files, onRefresh, formatFileSize }) {
  const [copiedId, setCopiedId] = useState(null);

  const [downloadingId, setDownloadingId] = useState(null);

  const handleDownload = async (file) => {
    setDownloadingId(file.id);
    
    try {
      const token = localStorage.getItem('token');
      const downloadUrl = `/api/files/download/${file.id}`;
      const fullDownloadUrl = `${window.location.origin}${downloadUrl}?token=${token}&download=1`;
      
      // Wenxi优化：使用window.open直接下载，避免iframe被拦截
      const newWindow = window.open(fullDownloadUrl, '_blank');
      
      // 如果弹窗被拦截，尝试直接下载
      if (!newWindow) {
        const link = document.createElement('a');
        link.href = fullDownloadUrl;
        link.download = file.original_filename;
        link.style.display = 'none';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      }
      
    } catch (error) {
      console.error('Wenxi - 文件下载错误:', error);
      alert('下载失败，请重试');
    } finally {
      setDownloadingId(null);
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
        onRefresh();
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
        setCopiedId(file.id);
        setTimeout(() => setCopiedId(null), 2000);
      } catch (clipboardError) {
        // 如果剪贴板API失败，显示分享链接
        prompt('分享链接已生成，请复制：', shareUrl);
      }
    } catch (error) {
      console.error('Wenxi - 文件分享错误:', error);
      alert('分享失败，请重试');
    }
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString('zh-CN');
  };

  if (files.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="mx-auto h-12 w-12 text-gray-400">
          <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
          </svg>
        </div>
        <h3 className="mt-2 text-sm font-medium text-gray-900">没有文件</h3>
        <p className="mt-1 text-sm text-gray-500">开始上传您的第一个文件</p>
      </div>
    );
  }

  return (
    <div className="bg-white shadow overflow-hidden sm:rounded-md">
      <ul className="divide-y divide-gray-200">
        {files.map((file) => (
          <li key={file.id} className="px-4 py-4 sm:px-6">
            <div className="flex items-center justify-between">
              <div className="flex-1 min-w-0">
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    <div className="h-8 w-8 rounded-full bg-gray-100 flex items-center justify-center">
                      <span className="text-sm font-medium text-gray-600">
                        {file.original_filename.split('.').pop()?.toUpperCase()}
                      </span>
                    </div>
                  </div>
                  <div className="ml-4">
                    <p className="text-sm font-medium text-gray-900 truncate">
                      {file.original_filename}
                    </p>
                    <p className="text-sm text-gray-500">
                      {formatFileSize(file.file_size)} • {formatDate(file.created_at)}
                    </p>
                  </div>
                </div>
              </div>
              <div className="ml-4 flex-shrink-0 flex space-x-2">
                <button
                  onClick={() => handleDownload(file)}
                  disabled={downloadingId === file.id}
                  className={`p-2 text-gray-400 hover:text-gray-600 ${downloadingId === file.id ? 'opacity-50 cursor-not-allowed' : ''}`}
                  title={downloadingId === file.id ? '下载中...' : '下载'}
                >
                  {downloadingId === file.id ? (
                    <svg className="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                  ) : (
                    <Download className="h-4 w-4" />
                  )}
                </button>
                <button
                  onClick={() => handleShare(file)}
                  className="p-2 text-gray-400 hover:text-gray-600"
                  title="分享"
                >
                  {copiedId === file.id ? (
                    <Check className="h-4 w-4 text-green-600" />
                  ) : (
                    <Share2 className="h-4 w-4" />
                  )}
                </button>
                <button
                  onClick={() => handleDelete(file.id)}
                  className="p-2 text-gray-400 hover:text-red-600"
                  title="删除"
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              </div>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}