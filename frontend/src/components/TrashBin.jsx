/**
 * Wenxi网盘 - 垃圾桶组件
 * 作者：Wenxi
 * 功能：管理已删除文件，支持恢复、永久删除、清空垃圾桶
 */

import React, { useState, useEffect } from 'react';
import { Trash2, RotateCcw, AlertTriangle, X, Trash, Clock, FileText } from 'lucide-react';

export default function TrashBin({ isOpen, onClose, onRefresh }) {
  const [trashFiles, setTrashFiles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({ file_count: 0, total_size: 0, expiring_soon: 0 });
  const [showConfirmDialog, setShowConfirmDialog] = useState(false);
  const [confirmAction, setConfirmAction] = useState(null);
  const [confirmFile, setConfirmFile] = useState(null);
  const [restoringId, setRestoringId] = useState(null);
  const [deletingId, setDeletingId] = useState(null);

  const fetchTrashFiles = async () => {
    try {
      const token = localStorage.getItem('token');
      const { getBaseURL } = await import('../utils/apiConfig');
      const baseURL = getBaseURL();
      
      const response = await fetch(`${baseURL}/api/trash/`, {
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
        throw new Error('获取垃圾桶列表失败');
      }

      const data = await response.json();
      setTrashFiles(data);
    } catch (error) {
      console.error('Wenxi - 获取垃圾桶列表失败:', error);
      alert('获取垃圾桶列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchTrashStats = async () => {
    try {
      const token = localStorage.getItem('token');
      const { getBaseURL } = await import('../utils/apiConfig');
      const baseURL = getBaseURL();
      
      const response = await fetch(`${baseURL}/api/trash/stats`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (response.ok) {
        const data = await response.json();
        setStats(data);
      }
    } catch (error) {
      console.error('Wenxi - 获取垃圾桶统计失败:', error);
    }
  };

  useEffect(() => {
    if (isOpen) {
      fetchTrashFiles();
      fetchTrashStats();
    }
  }, [isOpen]);

  const handleRestore = async (fileId) => {
    setRestoringId(fileId);
    try {
      const token = localStorage.getItem('token');
      const { getBaseURL } = await import('../utils/apiConfig');
      const baseURL = getBaseURL();
      
      const response = await fetch(`${baseURL}/api/trash/${fileId}/restore`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (response.ok) {
        setTrashFiles(prev => prev.filter(f => f.id !== fileId));
        fetchTrashStats();
        if (onRefresh) onRefresh();
        // 显示成功提示
        showToast('文件恢复成功');
      } else if (response.status === 401) {
        alert('登录已过期，请重新登录');
        window.location.href = '/login';
      } else {
        const error = await response.json();
        alert(error.detail || '恢复文件失败');
      }
    } catch (error) {
      console.error('Wenxi - 恢复文件失败:', error);
      alert('恢复文件失败');
    } finally {
      setRestoringId(null);
    }
  };

  const handlePermanentDelete = async (fileId) => {
    setDeletingId(fileId);
    try {
      const token = localStorage.getItem('token');
      const { getBaseURL } = await import('../utils/apiConfig');
      const baseURL = getBaseURL();
      
      const response = await fetch(`${baseURL}/api/trash/${fileId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (response.ok) {
        setTrashFiles(prev => prev.filter(f => f.id !== fileId));
        fetchTrashStats();
        showToast('文件已永久删除');
      } else if (response.status === 401) {
        alert('登录已过期，请重新登录');
        window.location.href = '/login';
      } else {
        const error = await response.json();
        alert(error.detail || '删除文件失败');
      }
    } catch (error) {
      console.error('Wenxi - 永久删除文件失败:', error);
      alert('永久删除文件失败');
    } finally {
      setDeletingId(null);
    }
  };

  const handleEmptyTrash = async () => {
    try {
      const token = localStorage.getItem('token');
      const { getBaseURL } = await import('../utils/apiConfig');
      const baseURL = getBaseURL();
      
      const response = await fetch(`${baseURL}/api/trash/`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (response.ok) {
        const data = await response.json();
        setTrashFiles([]);
        fetchTrashStats();
        if (onRefresh) onRefresh();
        showToast(`已清空垃圾桶，删除 ${data.deleted_count} 个文件`);
      } else if (response.status === 401) {
        alert('登录已过期，请重新登录');
        window.location.href = '/login';
      } else {
        const error = await response.json();
        alert(error.detail || '清空垃圾桶失败');
      }
    } catch (error) {
      console.error('Wenxi - 清空垃圾桶失败:', error);
      alert('清空垃圾桶失败');
    }
  };

  const showToast = (message) => {
    // 简单的toast提示实现
    const toast = document.createElement('div');
    toast.className = 'fixed bottom-4 right-4 bg-green-500 text-white px-6 py-3 rounded-lg shadow-lg z-50 animate-fade-in';
    toast.textContent = message;
    document.body.appendChild(toast);
    setTimeout(() => {
      toast.remove();
    }, 3000);
  };

  const openConfirmDialog = (action, file = null) => {
    setConfirmAction(action);
    setConfirmFile(file);
    setShowConfirmDialog(true);
  };

  const executeConfirmedAction = () => {
    setShowConfirmDialog(false);
    if (confirmAction === 'delete' && confirmFile) {
      handlePermanentDelete(confirmFile.id);
    } else if (confirmAction === 'empty') {
      handleEmptyTrash();
    }
  };

  const formatFileSize = (bytes) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString('zh-CN');
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] flex flex-col">
        {/* 头部 */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <div className="flex items-center">
            <Trash2 className="h-6 w-6 text-red-600 mr-3" />
            <div>
              <h2 className="text-xl font-bold text-gray-900">垃圾桶</h2>
              <p className="text-sm text-gray-500">
                {stats.file_count} 个文件 · {formatFileSize(stats.total_size)}
                {stats.expiring_soon > 0 && (
                  <span className="ml-2 text-orange-600">
                    ({stats.expiring_soon} 个文件即将过期)
                  </span>
                )}
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-400 hover:text-gray-600 transition-colors"
          >
            <X className="h-6 w-6" />
          </button>
        </div>

        {/* 内容区域 */}
        <div className="flex-1 overflow-y-auto p-6">
          {loading ? (
            <div className="text-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
              <p className="mt-4 text-gray-600">加载中...</p>
            </div>
          ) : trashFiles.length === 0 ? (
            <div className="text-center py-12">
              <Trash2 className="h-16 w-16 text-gray-300 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">垃圾桶是空的</h3>
              <p className="text-gray-500">删除的文件会在这里显示，保留30天</p>
            </div>
          ) : (
            <div className="space-y-4">
              {trashFiles.map((file) => (
                <div
                  key={file.id}
                  className="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <div className="flex items-center flex-1 min-w-0">
                    <div className="flex-shrink-0">
                      <div className="h-10 w-10 rounded-full bg-red-100 flex items-center justify-center">
                        <FileText className="h-5 w-5 text-red-600" />
                      </div>
                    </div>
                    <div className="ml-4 min-w-0">
                      <p className="text-sm font-medium text-gray-900 truncate">
                        {file.original_filename}
                      </p>
                      <div className="flex items-center text-xs text-gray-500 mt-1 space-x-3">
                        <span>{formatFileSize(file.file_size)}</span>
                        <span className="flex items-center">
                          <Clock className="h-3 w-3 mr-1" />
                          {formatDate(file.deleted_at)}
                        </span>
                        {file.days_until_deletion <= 7 && (
                          <span className={`flex items-center ${
                            file.days_until_deletion <= 3 ? 'text-red-600' : 'text-orange-600'
                          }`}>
                            <AlertTriangle className="h-3 w-3 mr-1" />
                            {file.days_until_deletion === 0 
                              ? '今天过期' 
                              : `剩余 ${file.days_until_deletion} 天`}
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2 ml-4">
                    <button
                      onClick={() => handleRestore(file.id)}
                      disabled={restoringId === file.id}
                      className="flex items-center px-3 py-1.5 text-sm font-medium text-blue-600 bg-blue-50 rounded-md hover:bg-blue-100 transition-colors disabled:opacity-50"
                    >
                      {restoringId === file.id ? (
                        <>
                          <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600 mr-2"></div>
                          恢复中...
                        </>
                      ) : (
                        <>
                          <RotateCcw className="h-4 w-4 mr-1" />
                          恢复
                        </>
                      )}
                    </button>
                    <button
                      onClick={() => openConfirmDialog('delete', file)}
                      disabled={deletingId === file.id}
                      className="flex items-center px-3 py-1.5 text-sm font-medium text-red-600 bg-red-50 rounded-md hover:bg-red-100 transition-colors disabled:opacity-50"
                    >
                      {deletingId === file.id ? (
                        <>
                          <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-red-600 mr-2"></div>
                          删除中...
                        </>
                      ) : (
                        <>
                          <Trash className="h-4 w-4 mr-1" />
                          删除
                        </>
                      )}
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* 底部操作栏 */}
        {trashFiles.length > 0 && (
          <div className="p-6 border-t border-gray-200 bg-gray-50 rounded-b-lg">
            <div className="flex items-center justify-between">
              <div className="text-sm text-gray-500">
                <AlertTriangle className="h-4 w-4 inline mr-1" />
                文件将在 30 天后自动永久删除
              </div>
              <button
                onClick={() => openConfirmDialog('empty')}
                className="flex items-center px-4 py-2 text-sm font-medium text-red-600 bg-white border border-red-300 rounded-md hover:bg-red-50 transition-colors"
              >
                <Trash2 className="h-4 w-4 mr-2" />
                清空垃圾桶
              </button>
            </div>
          </div>
        )}
      </div>

      {/* 确认对话框 */}
      {showConfirmDialog && (
        <div className="fixed inset-0 bg-black bg-opacity-50 z-60 flex items-center justify-center p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-md w-full p-6">
            <div className="flex items-center mb-4">
              <AlertTriangle className="h-6 w-6 text-red-600 mr-3" />
              <h3 className="text-lg font-bold text-gray-900">
                {confirmAction === 'delete' ? '确认永久删除？' : '确认清空垃圾桶？'}
              </h3>
            </div>
            <p className="text-gray-600 mb-6">
              {confirmAction === 'delete' 
                ? `确定要永久删除 "${confirmFile?.original_filename}" 吗？此操作不可恢复。`
                : '确定要清空垃圾桶吗？所有文件将被永久删除，此操作不可恢复。'
              }
            </p>
            <div className="flex justify-end space-x-3">
              <button
                onClick={() => setShowConfirmDialog(false)}
                className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 transition-colors"
              >
                取消
              </button>
              <button
                onClick={executeConfirmedAction}
                className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700 transition-colors"
              >
                确认删除
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
