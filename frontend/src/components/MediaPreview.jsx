/**
 * Wenxi网盘 - 媒体预览组件
 * 作者：Wenxi
 * 功能：支持图片、视频、音频、PDF等文件的预览
 * 支持：放大、缩小、下载、关闭操作
 */

import React, { useState, useRef, useEffect } from 'react';
import { 
  X, 
  ZoomIn, 
  ZoomOut, 
  Download, 
  ChevronLeft, 
  ChevronRight,
  FileText,
  FileImage,
  FileVideo,
  FileAudio,
  File
} from 'lucide-react';

// 支持预览的文件类型配置
export const SUPPORTED_PREVIEW_TYPES = {
  // 图片类型
  image: {
    extensions: ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg', 'ico'],
    mimeTypes: ['image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/bmp', 'image/svg+xml', 'image/x-icon'],
    icon: FileImage,
    label: '图片'
  },
  // 视频类型
  video: {
    extensions: ['mp4', 'webm', 'ogg', 'mov', 'mkv', 'avi'],
    mimeTypes: ['video/mp4', 'video/webm', 'video/ogg', 'video/quicktime', 'video/x-matroska', 'video/x-msvideo'],
    icon: FileVideo,
    label: '视频'
  },
  // 音频类型
  audio: {
    extensions: ['mp3', 'wav', 'ogg', 'flac', 'aac', 'm4a', 'wma'],
    mimeTypes: ['audio/mpeg', 'audio/wav', 'audio/ogg', 'audio/flac', 'audio/aac', 'audio/mp4', 'audio/x-ms-wma'],
    icon: FileAudio,
    label: '音频'
  },
  // PDF类型
  pdf: {
    extensions: ['pdf'],
    mimeTypes: ['application/pdf'],
    icon: FileText,
    label: 'PDF文档'
  }
};

/**
 * 获取文件扩展名
 */
export const getFileExtension = (filename) => {
  if (!filename) return '';
  const parts = filename.split('.');
  return parts.length > 1 ? parts.pop().toLowerCase() : '';
};

/**
 * 判断文件是否可预览
 */
export const isPreviewable = (file) => {
  if (!file) return { canPreview: false, type: null };
  
  const ext = getFileExtension(file.original_filename || file.filename);
  const mimeType = file.mime_type || '';
  
  for (const [type, config] of Object.entries(SUPPORTED_PREVIEW_TYPES)) {
    // 检查扩展名
    if (config.extensions.includes(ext)) {
      return { canPreview: true, type };
    }
    // 检查MIME类型
    if (config.mimeTypes.some(mt => mimeType.toLowerCase().includes(mt))) {
      return { canPreview: true, type };
    }
  }
  
  return { canPreview: false, type: null };
};

/**
 * 获取文件预览URL
 */
export const getPreviewUrl = (file) => {
  if (!file) return '';
  const token = localStorage.getItem('token');
  const { getBaseURL } = await import('../utils/apiConfig');
  return `${getBaseURL()}/api/files/download/${file.id}?token=${token}`;
};

/**
 * 媒体预览组件
 */
export default function MediaPreview({ 
  file, 
  files = [], 
  currentIndex = 0, 
  isOpen, 
  onClose, 
  onNavigate 
}) {
  const [scale, setScale] = useState(1);
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [isDragging, setIsDragging] = useState(false);
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 });
  const [isLoading, setIsLoading] = useState(true);
  const [previewUrl, setPreviewUrl] = useState('');
  const [error, setError] = useState(null);
  const containerRef = useRef(null);

  const previewInfo = isPreviewable(file);
  const fileType = previewInfo.type;

  // 获取预览URL
  useEffect(() => {
    if (file && isOpen) {
      setIsLoading(true);
      setError(null);
      setScale(1);
      setPosition({ x: 0, y: 0 });
      
      const loadPreviewUrl = async () => {
        try {
          const token = localStorage.getItem('token');
          const { getBaseURL } = await import('../utils/apiConfig');
          const url = `${getBaseURL()}/api/files/download/${file.id}?token=${token}`;
          setPreviewUrl(url);
        } catch (err) {
          setError('加载预览失败');
          setIsLoading(false);
        }
      };
      
      loadPreviewUrl();
    }
  }, [file, isOpen]);

  // 键盘快捷键支持
  useEffect(() => {
    if (!isOpen) return;

    const handleKeyDown = (e) => {
      switch (e.key) {
        case 'Escape':
          onClose();
          break;
        case 'ArrowLeft':
          handlePrevious();
          break;
        case 'ArrowRight':
          handleNext();
          break;
        case '+':
        case '=':
          handleZoomIn();
          break;
        case '-':
          handleZoomOut();
          break;
        case '0':
          handleResetZoom();
          break;
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, currentIndex, files.length]);

  // 阻止默认拖拽行为
  useEffect(() => {
    const preventDefault = (e) => e.preventDefault();
    if (isOpen) {
      document.addEventListener('dragstart', preventDefault);
      return () => document.removeEventListener('dragstart', preventDefault);
    }
  }, [isOpen]);

  const handleZoomIn = () => {
    setScale(prev => Math.min(prev * 1.2, 5));
  };

  const handleZoomOut = () => {
    setScale(prev => {
      const newScale = Math.max(prev / 1.2, 0.2);
      if (newScale <= 0.5) {
        setPosition({ x: 0, y: 0 });
      }
      return newScale;
    });
  };

  const handleResetZoom = () => {
    setScale(1);
    setPosition({ x: 0, y: 0 });
  };

  const handleMouseDown = (e) => {
    if (scale > 0.5 && fileType === 'image') {
      setIsDragging(true);
      setDragStart({ x: e.clientX - position.x, y: e.clientY - position.y });
    }
  };

  const handleMouseMove = (e) => {
    if (isDragging && scale > 0.5) {
      setPosition({
        x: e.clientX - dragStart.x,
        y: e.clientY - dragStart.y
      });
    }
  };

  const handleMouseUp = () => {
    setIsDragging(false);
  };

  const handlePrevious = () => {
    if (currentIndex > 0 && onNavigate) {
      onNavigate(currentIndex - 1);
    }
  };

  const handleNext = () => {
    if (currentIndex < files.length - 1 && onNavigate) {
      onNavigate(currentIndex + 1);
    }
  };

  const handleDownload = () => {
    if (file) {
      const link = document.createElement('a');
      link.href = previewUrl;
      link.download = file.original_filename;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    }
  };

  const handleWheel = (e) => {
    e.preventDefault();
    if (e.deltaY < 0) {
      handleZoomIn();
    } else {
      handleZoomOut();
    }
  };

  const handleLoad = () => {
    setIsLoading(false);
  };

  const handleError = () => {
    setIsLoading(false);
    setError('文件加载失败');
  };

  // 获取文件类型图标
  const FileIcon = previewInfo.canPreview 
    ? SUPPORTED_PREVIEW_TYPES[fileType]?.icon || File
    : File;

  if (!isOpen || !file) return null;

  return (
    <div 
      className="fixed inset-0 z-50 flex items-center justify-center"
      onClick={(e) => e.target === e.currentTarget && onClose()}
    >
      {/* 背景遮罩 */}
      <div className="absolute inset-0 bg-black/90 backdrop-blur-sm transition-opacity" />

      {/* 关闭按钮 */}
      <button
        onClick={onClose}
        className="absolute top-4 right-4 z-50 p-2 rounded-full bg-white/10 hover:bg-white/20 text-white transition-all duration-200 hover:scale-110"
        title="关闭 (Esc)"
      >
        <X className="h-6 w-6" />
      </button>

      {/* 文件信息栏 */}
      <div className="absolute top-0 left-0 right-0 z-40 bg-gradient-to-b from-black/70 to-transparent px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <FileIcon className="h-6 w-6 text-blue-400" />
            <div>
              <h3 className="text-white font-medium truncate max-w-md">
                {file.original_filename}
              </h3>
              <p className="text-white/60 text-sm">
                {SUPPORTED_PREVIEW_TYPES[fileType]?.label || '文件'} • {(file.file_size / 1024 / 1024).toFixed(2)} MB
              </p>
            </div>
          </div>
          
          {/* 缩放控制 */}
          {fileType === 'image' && (
            <div className="flex items-center space-x-2 bg-black/40 rounded-lg px-3 py-1.5">
              <span className="text-white/60 text-sm min-w-[60px] text-center">
                {Math.round(scale * 100)}%
              </span>
              <button
                onClick={handleZoomOut}
                className="p-1.5 rounded hover:bg-white/20 text-white transition-colors"
                title="缩小 (-)"
              >
                <ZoomOut className="h-4 w-4" />
              </button>
              <button
                onClick={handleResetZoom}
                className="p-1.5 rounded hover:bg-white/20 text-white transition-colors"
                title="重置 (0)"
              >
                <span className="text-xs font-medium">1:1</span>
              </button>
              <button
                onClick={handleZoomIn}
                className="p-1.5 rounded hover:bg-white/20 text-white transition-colors"
                title="放大 (+)"
              >
                <ZoomIn className="h-4 w-4" />
              </button>
            </div>
          )}
        </div>
      </div>

      {/* 导航按钮 */}
      {files.length > 1 && (
        <>
          <button
            onClick={handlePrevious}
            disabled={currentIndex === 0}
            className={`absolute left-4 z-40 p-3 rounded-full bg-white/10 hover:bg-white/20 text-white transition-all duration-200 ${
              currentIndex === 0 ? 'opacity-30 cursor-not-allowed' : 'hover:scale-110'
            }`}
            title="上一个 (←)"
          >
            <ChevronLeft className="h-6 w-6" />
          </button>
          <button
            onClick={handleNext}
            disabled={currentIndex === files.length - 1}
            className={`absolute right-4 z-40 p-3 rounded-full bg-white/10 hover:bg-white/20 text-white transition-all duration-200 ${
              currentIndex === files.length - 1 ? 'opacity-30 cursor-not-allowed' : 'hover:scale-110'
            }`}
            title="下一个 (→)"
          >
            <ChevronRight className="h-6 w-6" />
          </button>

          {/* 文件计数器 */}
          <div className="absolute bottom-20 left-1/2 transform -translate-x-1/2 z-40 bg-black/60 rounded-full px-4 py-1.5">
            <span className="text-white/80 text-sm">
              {currentIndex + 1} / {files.length}
            </span>
          </div>
        </>
      )}

      {/* 预览内容区域 */}
      <div 
        ref={containerRef}
        className="relative z-30 w-full h-full flex items-center justify-center overflow-hidden"
        onMouseDown={handleMouseDown}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        onMouseLeave={handleMouseUp}
        onWheel={handleWheel}
      >
        {/* 加载状态 */}
        {isLoading && (
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
          </div>
        )}

        {/* 错误状态 */}
        {error && (
          <div className="text-center">
            <FileIcon className="h-16 w-16 text-white/30 mx-auto mb-4" />
            <p className="text-white/60">{error}</p>
          </div>
        )}

        {/* 图片预览 */}
        {fileType === 'image' && previewUrl && (
          <img
            src={previewUrl}
            alt={file.original_filename}
            className={`max-w-[90%] max-h-[85vh] object-contain transition-transform duration-100 ${
              isDragging ? 'cursor-grabbing' : scale > 0.5 ? 'cursor-grab' : 'cursor-default'
            }`}
            style={{
              transform: `translate(${position.x}px, ${position.y}px) scale(${scale})`,
              opacity: isLoading ? 0 : 1
            }}
            onLoad={handleLoad}
            onError={handleError}
            draggable={false}
          />
        )}

        {/* 视频预览 */}
        {fileType === 'video' && previewUrl && (
          <video
            src={previewUrl}
            controls
            className="max-w-[90%] max-h-[85vh] rounded-lg shadow-2xl"
            onLoadedData={handleLoad}
            onError={handleError}
            style={{ opacity: isLoading ? 0 : 1 }}
          >
            您的浏览器不支持视频播放
          </video>
        )}

        {/* 音频预览 */}
        {fileType === 'audio' && previewUrl && (
          <div className="bg-white/10 backdrop-blur-md rounded-2xl p-8 max-w-md w-full mx-4">
            <FileAudio className="h-16 w-16 text-blue-400 mx-auto mb-6" />
            <audio
              src={previewUrl}
              controls
              className="w-full"
              onLoadedData={handleLoad}
              onError={handleError}
            >
              您的浏览器不支持音频播放
            </audio>
          </div>
        )}

        {/* PDF预览 */}
        {fileType === 'pdf' && previewUrl && (
          <div className="w-[90%] h-[85vh] bg-white rounded-lg overflow-hidden shadow-2xl">
            <iframe
              src={`${previewUrl}#toolbar=1&navpanes=1`}
              className="w-full h-full"
              title={file.original_filename}
              onLoad={handleLoad}
              onError={handleError}
            />
          </div>
        )}
      </div>

      {/* 底部工具栏 */}
      <div className="absolute bottom-0 left-0 right-0 z-40 bg-gradient-to-t from-black/70 to-transparent px-6 py-4">
        <div className="flex items-center justify-center space-x-4">
          <button
            onClick={handleDownload}
            className="flex items-center space-x-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-all duration-200 hover:scale-105"
            title="下载文件"
          >
            <Download className="h-4 w-4" />
            <span className="text-sm font-medium">下载</span>
          </button>
        </div>
        
        {/* 快捷键提示 */}
        <div className="mt-3 text-center text-white/40 text-xs">
          快捷键: ← → 切换 • +/- 缩放 • 0 重置 • Esc 关闭
        </div>
      </div>
    </div>
  );
}
