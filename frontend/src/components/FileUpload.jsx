/**
 * Wenxi网盘 - 文件上传组件
 * 作者：Wenxi
 * 功能：提供拖拽上传、进度显示、批量上传等功能
 */

import React, { useState, useRef } from 'react';
import { X, UploadCloud, FileText } from 'lucide-react';

export default function FileUpload({ onClose, onSuccess }) {
  const [files, setFiles] = useState([]);
  const [uploading, setUploading] = useState(false);
  const [dragOver, setDragOver] = useState(false);
  const fileInputRef = useRef(null);

  const handleFileSelect = (selectedFiles) => {
    const newFiles = Array.from(selectedFiles).map(file => ({
      file,
      id: Math.random().toString(36).substr(2, 9),
      progress: 0,
      status: 'pending'
    }));
    setFiles(prev => [...prev, ...newFiles]);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    setDragOver(false);
    handleFileSelect(e.dataTransfer.files);
  };

  const handleDragOver = (e) => {
    e.preventDefault();
    setDragOver(true);
  };

  const handleDragLeave = (e) => {
    e.preventDefault();
    setDragOver(false);
  };

  const uploadFile = async (fileObj) => {
    const file = fileObj.file;
    const token = localStorage.getItem('token');
    
    // 设置文件状态为上传中
    setFiles(prev => prev.map(f => 
      f.id === fileObj.id 
        ? { ...f, status: 'uploading' }
        : f
    ));
    
    // Wenxi上传策略：16MB分块传输
    const CHUNK_SIZE = 16 * 1024 * 1024; // 16MB分块（提升8倍速度）
    const isLargeFile = file.size > CHUNK_SIZE; // 大于16MB使用分块
    
    try {
      if (!isLargeFile) {
        // 小文件直接上传
        return await uploadSingleFile(fileObj, token);
      } else {
        // 大文件分块上传
        return await uploadChunkedFile(fileObj, CHUNK_SIZE, token);
      }
    } catch (error) {
      setFiles(prev => prev.map(f => 
        f.id === fileObj.id 
          ? { ...f, status: 'error' }
          : f
      ));
      console.error('Wenxi - 文件上传错误:', error);
      return false;
    }
  };

  const uploadSingleFile = async (fileObj, token) => {
    const formData = new FormData();
    formData.append('file', fileObj.file);

    // 使用 XMLHttpRequest 来支持进度监听
    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest();
      
      // 监听上传进度
      xhr.upload.addEventListener('progress', (event) => {
        if (event.lengthComputable) {
          const progress = Math.round((event.loaded / event.total) * 100);
          updateProgress(fileObj.id, progress);
        }
      });

      // 监听上传完成
      xhr.addEventListener('load', () => {
        if (xhr.status === 200) {
          setFiles(prev => prev.map(f => 
            f.id === fileObj.id 
              ? { ...f, status: 'completed', progress: 100 }
              : f
          ));
          resolve(true);
        } else if (xhr.status === 401) {
          reject(new Error('认证失败，请重新登录'));
        } else {
          reject(new Error('上传失败'));
        }
      });

      // 监听错误
      xhr.addEventListener('error', () => {
        reject(new Error('上传失败'));
      });

      // 设置请求头并发送
      xhr.open('POST', '/api/files/upload');
      xhr.setRequestHeader('Authorization', `Bearer ${token}`);
      xhr.send(formData);
    });
  };

  const uploadChunkedFile = async (fileObj, chunkSize, token) => {
    const file = fileObj.file;
    const totalChunks = Math.ceil(file.size / chunkSize);
    const fileHash = await calculateFileHash(file);
    
    // 检查已上传的分块
    const uploadedChunks = await checkUploadedChunks(fileHash, token);
    
    for (let i = 0; i < totalChunks; i++) {
      if (uploadedChunks.includes(i)) {
        // 跳过已上传的分块
        updateProgress(fileObj.id, (i + 1) / totalChunks * 100);
        continue;
      }
      
      const start = i * chunkSize;
      const end = Math.min(start + chunkSize, file.size);
      const chunk = file.slice(start, end);
      
      const chunkFormData = new FormData();
      chunkFormData.append('chunk', chunk);
      chunkFormData.append('chunk_index', i);
      chunkFormData.append('total_chunks', totalChunks);
      chunkFormData.append('file_name', file.name);
      chunkFormData.append('file_hash', fileHash);
      chunkFormData.append('chunk_hash', await calculateFileHash(chunk));
      
      // 使用 XMLHttpRequest 来支持分块上传的进度监听
      await new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        
        // 监听分块上传进度
        xhr.upload.addEventListener('progress', (event) => {
          if (event.lengthComputable) {
            // 计算整体进度：已上传分块 + 当前分块进度
            const baseProgress = i / totalChunks * 100;
            const chunkProgress = (event.loaded / event.total) * (100 / totalChunks);
            const totalProgress = Math.round(baseProgress + chunkProgress);
            updateProgress(fileObj.id, totalProgress);
          }
        });
        
        xhr.addEventListener('load', () => {
          if (xhr.status === 200) {
            updateProgress(fileObj.id, Math.round((i + 1) / totalChunks * 100));
            resolve();
          } else {
            reject(new Error(`分块 ${i + 1} 上传失败`));
          }
        });
        
        xhr.addEventListener('error', () => {
          reject(new Error(`分块 ${i + 1} 上传失败`));
        });
        
        xhr.open('POST', '/api/files/upload/chunk');
        xhr.setRequestHeader('Authorization', `Bearer ${token}`);
        xhr.send(chunkFormData);
      });
    }
    
    // 合并分块
    const mergeFormData = new FormData();
    mergeFormData.append('file_name', file.name);
    mergeFormData.append('file_hash', fileHash);
    mergeFormData.append('total_chunks', totalChunks);
    
    const mergeResponse = await fetch('/api/files/upload/merge', {
      method: 'POST',
      body: mergeFormData,
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (mergeResponse.ok) {
      setFiles(prev => prev.map(f => 
        f.id === fileObj.id 
          ? { ...f, status: 'completed', progress: 100 }
          : f
      ));
      return true;
    } else {
      throw new Error('合并分块失败');
    }
  };

  const checkUploadedChunks = async (fileHash, token) => {
    try {
      const response = await fetch(`/api/files/upload/check?file_hash=${fileHash}`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      if (response.ok) {
        const data = await response.json();
        return data.uploaded_chunks || [];
      }
      return [];
    } catch (error) {
      return [];
    }
  };

  const calculateFileHash = async (file) => {
    const buffer = await file.arrayBuffer();
    const hashBuffer = await crypto.subtle.digest('SHA-256', buffer);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
  };

  const updateProgress = (fileId, progress) => {
    setFiles(prev => prev.map(f => 
      f.id === fileId 
        ? { ...f, progress: Math.round(progress) }
        : f
    ));
  };

  const handleUpload = async (e) => {
    if (e) e.preventDefault(); // Wenxi：阻止默认行为防止页面刷新
    if (files.length === 0) return;

    setUploading(true);
    const results = await Promise.all(files.map(uploadFile));
    setUploading(false);

    if (results.every(r => r)) {
      onSuccess();
    }
  };

  const removeFile = (id) => {
    setFiles(prev => prev.filter(f => f.id !== id));
  };

  const formatFileSize = (bytes) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-full max-w-2xl max-h-[80vh] overflow-hidden">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold text-gray-900">上传文件</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <X className="h-6 w-6" />
          </button>
        </div>

        <div
          className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
            dragOver ? 'border-blue-500 bg-blue-50' : 'border-gray-300'
          }`}
          onDrop={handleDrop}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
        >
          <UploadCloud className="mx-auto h-12 w-12 text-gray-400" />
          <p className="mt-2 text-sm text-gray-600">
            拖拽文件到此处或{' '}
            <button
              onClick={() => fileInputRef.current?.click()}
              className="text-blue-600 hover:text-blue-500"
            >
              选择文件
            </button>
          </p>
          <p className="mt-1 text-xs text-gray-500">
            支持图片、文档、视频、音频、压缩包等所有格式
          </p>
          <input
            ref={fileInputRef}
            type="file"
            multiple
            accept="*/*"
            className="hidden"
            onChange={(e) => {
              e.preventDefault(); // Wenxi：防止移动端选择文件后刷新
              handleFileSelect(e.target.files);
            }}
          />
        </div>

        {files.length > 0 && (
          <div className="mt-4 space-y-2 max-h-64 overflow-y-auto">
            {files.map((fileObj) => (
              <div key={fileObj.id} className="flex items-center justify-between p-3 bg-gray-50 rounded">
                <div className="flex items-center flex-1">
                  <FileText className="h-5 w-5 text-gray-400 mr-3" />
                  <div className="flex-1">
                    <p className="text-sm font-medium text-gray-900">{fileObj.file.name}</p>
                    <p className="text-xs text-gray-500">{formatFileSize(fileObj.file.size)}</p>
                    
                    {/* 进度条 */}
                    {fileObj.progress > 0 && fileObj.progress < 100 && (
                      <div className="mt-2">
                        <div className="flex justify-between text-xs text-gray-600 mb-1">
                          <span>上传进度</span>
                          <span>{fileObj.progress}%</span>
                        </div>
                        <div className="w-full bg-gray-200 rounded-full h-2">
                          <div 
                            className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                            style={{ width: `${fileObj.progress}%` }}
                          ></div>
                        </div>
                      </div>
                    )}
                  </div>
                </div>
                <div className="flex items-center space-x-2">
                  {fileObj.status === 'uploading' && (
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
                  )}
                  {fileObj.status === 'completed' && (
                    <span className="text-xs text-green-600">✓</span>
                  )}
                  {fileObj.status === 'error' && (
                    <span className="text-xs text-red-600">✗</span>
                  )}
                  <button
                    onClick={() => removeFile(fileObj.id)}
                    className="text-gray-400 hover:text-gray-600"
                  >
                    <X className="h-4 w-4" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}

        <div className="mt-6 flex justify-end space-x-3">
          <button
            onClick={onClose}
            className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
          >
            取消
          </button>          <button
            onClick={(e) => handleUpload(e)}
            disabled={files.length === 0 || uploading}
            className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 disabled:opacity-50"
            type="button"  // Wenxi：明确指定按钮类型防止表单提交
          >
            {uploading ? '上传中...' : '开始上传'}
          </button>
        </div>
      </div>
    </div>
  );
}