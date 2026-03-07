/**
 * Wenxi网盘 - 文件夹管理组件
 * 作者：Wenxi
 * 功能：文件夹树形导航、面包屑、新建文件夹、拖拽移动
 */

import React, { useState, useEffect } from 'react';
import { 
  Folder, FolderOpen, File, ChevronRight, ChevronDown, 
  Home, Plus, MoreVertical, Edit2, Trash2, Move, X
} from 'lucide-react';

// 文件夹树形组件
export function FolderTree({ 
  folders, 
  currentFolderId, 
  onSelectFolder, 
  onToggleExpand, 
  expandedFolders,
  level = 0 
}) {
  return (
    <div className="space-y-1">
      {folders.map(folder => (
        <div key={folder.id}>
          <div
            className={`flex items-center py-1.5 px-2 rounded-lg cursor-pointer transition-all duration-200 ${
              currentFolderId === folder.id
                ? 'bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-300'
                : 'hover:bg-gray-100 dark:hover:bg-slate-800 text-gray-700 dark:text-slate-300'
            }`}
            style={{ paddingLeft: `${level * 16 + 8}px` }}
            onClick={() => onSelectFolder(folder.id)}
          >
            <button
              onClick={(e) => {
                e.stopPropagation();
                onToggleExpand(folder.id);
              }}
              className="mr-1 p-0.5 hover:bg-gray-200 dark:hover:bg-slate-700 rounded"
            >
              {folder.children?.length > 0 && (
                expandedFolders.has(folder.id) 
                  ? <ChevronDown className="h-3.5 w-3.5" />
                  : <ChevronRight className="h-3.5 w-3.5" />
              )}
            </button>
            
            {expandedFolders.has(folder.id) 
              ? <FolderOpen className="h-4 w-4 mr-2 text-yellow-500" />
              : <Folder className="h-4 w-4 mr-2 text-yellow-500" />
            }
            
            <span className="text-sm truncate">{folder.name}</span>
            
            {folder.children?.length > 0 && (
              <span className="ml-1.5 text-xs text-gray-400 dark:text-slate-500">
                ({folder.children.length})
              </span>
            )}
          </div>
          
          {expandedFolders.has(folder.id) && folder.children?.length > 0 && (
            <FolderTree
              folders={folder.children}
              currentFolderId={currentFolderId}
              onSelectFolder={onSelectFolder}
              onToggleExpand={onToggleExpand}
              expandedFolders={expandedFolders}
              level={level + 1}
            />
          )}
        </div>
      ))}
    </div>
  );
}

// 面包屑导航组件
export function Breadcrumb({ breadcrumbs, onNavigate }) {
  return (
    <nav className="flex items-center space-x-1 text-sm text-gray-600 dark:text-slate-400">
      <button
        onClick={() => onNavigate(null)}
        className="flex items-center hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
      >
        <Home className="h-4 w-4" />
        <span className="ml-1">根目录</span>
      </button>
      
      {breadcrumbs.map((crumb, index) => (
        <React.Fragment key={crumb.id}>
          <ChevronRight className="h-4 w-4 text-gray-400" />
          <button
            onClick={() => onNavigate(crumb.id)}
            className={`hover:text-blue-600 dark:hover:text-blue-400 transition-colors ${
              index === breadcrumbs.length - 1 ? 'font-medium text-gray-900 dark:text-slate-200' : ''
            }`}
          >
            {crumb.name}
          </button>
        </React.Fragment>
      ))}
    </nav>
  );
}

// 新建文件夹对话框
export function CreateFolderDialog({ isOpen, onClose, onCreate, parentId }) {
  const [folderName, setFolderName] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (isOpen) {
      setFolderName('');
    }
  }, [isOpen]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!folderName.trim()) return;
    
    setIsLoading(true);
    try {
      await onCreate(folderName.trim(), parentId);
      onClose();
    } catch (error) {
      console.error('创建文件夹失败:', error);
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white dark:bg-slate-800 rounded-lg shadow-xl w-96 p-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-medium text-gray-900 dark:text-white">
            新建文件夹
          </h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
          >
            <X className="h-5 w-5" />
          </button>
        </div>
        
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            value={folderName}
            onChange={(e) => setFolderName(e.target.value)}
            placeholder="文件夹名称"
            className="w-full px-3 py-2 border border-gray-300 dark:border-slate-600 rounded-lg 
                     bg-white dark:bg-slate-700 text-gray-900 dark:text-white
                     focus:ring-2 focus:ring-blue-500 focus:border-transparent
                     placeholder-gray-400 dark:placeholder-slate-500"
            autoFocus
          />
          
          <div className="flex justify-end space-x-3 mt-4">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm text-gray-600 dark:text-slate-400 hover:text-gray-800 dark:hover:text-slate-200"
            >
              取消
            </button>
            <button
              type="submit"
              disabled={!folderName.trim() || isLoading}
              className="px-4 py-2 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700 
                       disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {isLoading ? '创建中...' : '创建'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

// 文件夹操作菜单
export function FolderActions({ folder, onRename, onDelete, onMove }) {
  const [isOpen, setIsOpen] = useState(false);
  const menuRef = React.useRef(null);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (menuRef.current && !menuRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  return (
    <div className="relative" ref={menuRef}>
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="p-1.5 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 rounded-lg hover:bg-gray-100 dark:hover:bg-slate-700"
      >
        <MoreVertical className="h-4 w-4" />
      </button>
      
      {isOpen && (
        <div className="absolute right-0 mt-1 w-32 bg-white dark:bg-slate-800 rounded-lg shadow-lg border border-gray-200 dark:border-slate-700 py-1 z-10">
          <button
            onClick={() => { onRename(folder); setIsOpen(false); }}
            className="w-full px-3 py-2 text-left text-sm text-gray-700 dark:text-slate-300 hover:bg-gray-100 dark:hover:bg-slate-700 flex items-center"
          >
            <Edit2 className="h-3.5 w-3.5 mr-2" />
            重命名
          </button>
          <button
            onClick={() => { onMove(folder); setIsOpen(false); }}
            className="w-full px-3 py-2 text-left text-sm text-gray-700 dark:text-slate-300 hover:bg-gray-100 dark:hover:bg-slate-700 flex items-center"
          >
            <Move className="h-3.5 w-3.5 mr-2" />
            移动
          </button>
          <button
            onClick={() => { onDelete(folder); setIsOpen(false); }}
            className="w-full px-3 py-2 text-left text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 flex items-center"
          >
            <Trash2 className="h-3.5 w-3.5 mr-2" />
            删除
          </button>
        </div>
      )}
    </div>
  );
}

// 文件夹列表视图（图标模式）
export function FolderGrid({ 
  folders, 
  onSelectFolder, 
  onRename, 
  onDelete, 
  onMove,
  onDropFile 
}) {
  const [draggedOver, setDraggedOver] = useState(null);

  const handleDragOver = (e, folderId) => {
    e.preventDefault();
    setDraggedOver(folderId);
  };

  const handleDragLeave = () => {
    setDraggedOver(null);
  };

  const handleDrop = (e, folderId) => {
    e.preventDefault();
    setDraggedOver(null);
    
    const fileId = e.dataTransfer.getData('fileId');
    if (fileId && onDropFile) {
      onDropFile(parseInt(fileId), folderId);
    }
  };

  if (folders.length === 0) {
    return null;
  }

  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 mb-6">
      {folders.map(folder => (
        <div
          key={folder.id}
          onClick={() => onSelectFolder(folder.id)}
          onDragOver={(e) => handleDragOver(e, folder.id)}
          onDragLeave={handleDragLeave}
          onDrop={(e) => handleDrop(e, folder.id)}
          className={`relative group p-4 rounded-xl border-2 cursor-pointer transition-all duration-200 ${
            draggedOver === folder.id
              ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
              : 'border-transparent hover:bg-gray-50 dark:hover:bg-slate-800 bg-white dark:bg-slate-900'
          }`}
        >
          <div className="flex flex-col items-center">
            <div className="relative">
              <Folder className="h-12 w-12 text-yellow-500 mb-2" />
              {draggedOver === folder.id && (
                <div className="absolute inset-0 flex items-center justify-center">
                  <div className="bg-blue-500 text-white rounded-full p-1">
                    <Plus className="h-4 w-4" />
                  </div>
                </div>
              )}
            </div>
            <span className="text-sm text-center text-gray-700 dark:text-slate-300 truncate w-full">
              {folder.name}
            </span>
            <span className="text-xs text-gray-400 dark:text-slate-500 mt-1">
              {folder.file_count || 0} 个文件
            </span>
          </div>
          
          <div 
            className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity"
            onClick={(e) => e.stopPropagation()}
          >
            <FolderActions
              folder={folder}
              onRename={onRename}
              onDelete={onDelete}
              onMove={onMove}
            />
          </div>
        </div>
      ))}
    </div>
  );
}

// 重命名文件夹对话框
export function RenameFolderDialog({ isOpen, onClose, onRename, folder }) {
  const [folderName, setFolderName] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (isOpen && folder) {
      setFolderName(folder.name);
    }
  }, [isOpen, folder]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!folderName.trim() || folderName === folder?.name) {
      onClose();
      return;
    }
    
    setIsLoading(true);
    try {
      await onRename(folder.id, folderName.trim());
      onClose();
    } catch (error) {
      console.error('重命名文件夹失败:', error);
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen || !folder) return null;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white dark:bg-slate-800 rounded-lg shadow-xl w-96 p-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-medium text-gray-900 dark:text-white">
            重命名文件夹
          </h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
          >
            <X className="h-5 w-5" />
          </button>
        </div>
        
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            value={folderName}
            onChange={(e) => setFolderName(e.target.value)}
            placeholder="文件夹名称"
            className="w-full px-3 py-2 border border-gray-300 dark:border-slate-600 rounded-lg 
                     bg-white dark:bg-slate-700 text-gray-900 dark:text-white
                     focus:ring-2 focus:ring-blue-500 focus:border-transparent
                     placeholder-gray-400 dark:placeholder-slate-500"
            autoFocus
          />
          
          <div className="flex justify-end space-x-3 mt-4">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm text-gray-600 dark:text-slate-400 hover:text-gray-800 dark:hover:text-slate-200"
            >
              取消
            </button>
            <button
              type="submit"
              disabled={!folderName.trim() || isLoading}
              className="px-4 py-2 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700 
                       disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {isLoading ? '重命名中...' : '重命名'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

// 移动文件夹对话框
export function MoveFolderDialog({ isOpen, onClose, onMove, folder, allFolders }) {
  const [selectedFolderId, setSelectedFolderId] = useState(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (isOpen) {
      setSelectedFolderId(null);
    }
  }, [isOpen]);

  const handleSubmit = async () => {
    if (selectedFolderId === folder?.id) {
      onClose();
      return;
    }
    
    setIsLoading(true);
    try {
      await onMove(folder.id, selectedFolderId);
      onClose();
    } catch (error) {
      console.error('移动文件夹失败:', error);
    } finally {
      setIsLoading(false);
    }
  };

  // 递归渲染文件夹树（排除当前文件夹及其子文件夹）
  const renderFolderOptions = (folders, level = 0) => {
    return folders.map(f => {
      if (f.id === folder?.id) return null; // 排除当前文件夹
      
      return (
        <React.Fragment key={f.id}>
          <button
            onClick={() => setSelectedFolderId(f.id)}
            className={`w-full text-left px-3 py-2 text-sm rounded-lg transition-colors ${
              selectedFolderId === f.id
                ? 'bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-300'
                : 'hover:bg-gray-100 dark:hover:bg-slate-700 text-gray-700 dark:text-slate-300'
            }`}
            style={{ paddingLeft: `${level * 16 + 12}px` }}
          >
            <Folder className="h-4 w-4 inline mr-2 text-yellow-500" />
            {f.name}
          </button>
          {f.children?.length > 0 && renderFolderOptions(f.children, level + 1)}
        </React.Fragment>
      );
    });
  };

  if (!isOpen || !folder) return null;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white dark:bg-slate-800 rounded-lg shadow-xl w-96 p-6 max-h-[80vh] flex flex-col">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-medium text-gray-900 dark:text-white">
            移动文件夹
          </h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
          >
            <X className="h-5 w-5" />
          </button>
        </div>
        
        <div className="mb-3 text-sm text-gray-600 dark:text-slate-400">
          将 "{folder.name}" 移动到:
        </div>
        
        <div className="flex-1 overflow-y-auto space-y-1 max-h-64">
          <button
            onClick={() => setSelectedFolderId(null)}
            className={`w-full text-left px-3 py-2 text-sm rounded-lg transition-colors ${
              selectedFolderId === null
                ? 'bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-300'
                : 'hover:bg-gray-100 dark:hover:bg-slate-700 text-gray-700 dark:text-slate-300'
            }`}
          >
            <Home className="h-4 w-4 inline mr-2 text-blue-500" />
            根目录
          </button>
          {renderFolderOptions(allFolders)}
        </div>
        
        <div className="flex justify-end space-x-3 mt-4 pt-4 border-t border-gray-200 dark:border-slate-700">
          <button
            type="button"
            onClick={onClose}
            className="px-4 py-2 text-sm text-gray-600 dark:text-slate-400 hover:text-gray-800 dark:hover:text-slate-200"
          >
            取消
          </button>
          <button
            onClick={handleSubmit}
            disabled={isLoading}
            className="px-4 py-2 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700 
                     disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {isLoading ? '移动中...' : '移动'}
          </button>
        </div>
      </div>
    </div>
  );
}

// 侧边栏文件夹导航
export function FolderSidebar({ 
  folderTree, 
  currentFolderId, 
  onSelectFolder, 
  onCreateFolder,
  isExpanded = true 
}) {
  const [expandedFolders, setExpandedFolders] = useState(new Set());

  const handleToggleExpand = (folderId) => {
    setExpandedFolders(prev => {
      const newSet = new Set(prev);
      if (newSet.has(folderId)) {
        newSet.delete(folderId);
      } else {
        newSet.add(folderId);
      }
      return newSet;
    });
  };

  if (!isExpanded) return null;

  return (
    <div className="w-64 bg-white dark:bg-slate-900 border-r border-gray-200 dark:border-slate-800 flex flex-col">
      <div className="p-4 border-b border-gray-200 dark:border-slate-800">
        <button
          onClick={() => onCreateFolder()}
          className="w-full flex items-center justify-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          <Plus className="h-4 w-4 mr-2" />
          新建文件夹
        </button>
      </div>
      
      <div className="flex-1 overflow-y-auto p-2">
        <div
          className={`flex items-center py-2 px-3 rounded-lg cursor-pointer mb-2 ${
            currentFolderId === null
              ? 'bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-300'
              : 'hover:bg-gray-100 dark:hover:bg-slate-800 text-gray-700 dark:text-slate-300'
          }`}
          onClick={() => onSelectFolder(null)}
        >
          <Home className="h-4 w-4 mr-2 text-blue-500" />
          <span className="text-sm font-medium">根目录</span>
        </div>
        
        {folderTree.length > 0 && (
          <div className="mt-2">
            <div className="text-xs font-semibold text-gray-400 dark:text-slate-500 uppercase tracking-wider px-3 mb-2">
              我的文件夹
            </div>
            <FolderTree
              folders={folderTree}
              currentFolderId={currentFolderId}
              onSelectFolder={onSelectFolder}
              onToggleExpand={handleToggleExpand}
              expandedFolders={expandedFolders}
            />
          </div>
        )}
      </div>
    </div>
  );
}

export default {
  FolderTree,
  Breadcrumb,
  CreateFolderDialog,
  FolderActions,
  FolderGrid,
  RenameFolderDialog,
  MoveFolderDialog,
  FolderSidebar
};
