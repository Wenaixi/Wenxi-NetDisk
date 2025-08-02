/**
 * Wenxi网盘 - 用户设置组件
 * 作者：Wenxi
 * 功能：用户个人信息管理、密码修改、存储配额管理
 */

import React, { useState } from 'react';
import axios from 'axios';
import { X, User, Mail, Lock, Eye, EyeOff, AlertTriangle, Trash2 } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';

export default function UserSettings({ isOpen, onClose, currentUser, onUserUpdate }) {
  const { logout } = useAuth();
  const [activeTab, setActiveTab] = useState('username');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  // 表单状态
  const [username, setUsername] = useState(currentUser?.username || '');
  const [email, setEmail] = useState(currentUser?.email || '');
  const [password, setPassword] = useState('');
  const [newEmail, setNewEmail] = useState('');
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  
  // 注销账户状态
  const [deleteUsername, setDeleteUsername] = useState('');
  const [deletePassword, setDeletePassword] = useState('');
  const [deleteEmail, setDeleteEmail] = useState('');
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  // 密码可见性管理 - 为每个密码字段单独管理
  const [showFields, setShowFields] = useState({
    emailPassword: false,
    oldPassword: false,
    newPassword: false,
    confirmPassword: false,
    deletePassword: false
  });

  const toggleShowField = (field) => {
    setShowFields(prev => ({ ...prev, [field]: !prev[field] }));
  };

  const handleUsernameUpdate = async () => {
    if (!username.trim()) {
      setError('用户名不能为空');
      return;
    }

    setLoading(true);
    setError('');
    setSuccess('');

    try {
      await axios.put('/api/auth/update-username', { username });
      setSuccess('用户名修改成功');
      onUserUpdate({ ...currentUser, username });
      setTimeout(() => {
        setSuccess('');
      }, 3000);
    } catch (err) {
      setError(err.response?.data?.detail || '修改失败');
    } finally {
      setLoading(false);
    }
  };

  const handleEmailUpdate = async () => {
    if (!newEmail.trim()) {
      setError('邮箱不能为空');
      return;
    }

    if (!password) {
      setError('请输入密码验证');
      return;
    }

    setLoading(true);
    setError('');
    setSuccess('');

    try {
      await axios.put('/api/auth/update-email', { email: newEmail, password });
      setSuccess('邮箱修改成功，请重新登录');
      onUserUpdate({ ...currentUser, email: newEmail });
      setEmail(newEmail);
      setNewEmail('');
      setPassword('');
      
      // 2秒后自动登出
      setTimeout(() => {
        onClose();
        logout();
      }, 2000);
    } catch (err) {
      setError(err.response?.data?.detail || '修改失败');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteAccount = async () => {
    if (!deleteUsername.trim() || !deletePassword || !deleteEmail.trim()) {
      setError('请填写所有必填信息');
      return;
    }

    setLoading(true);
    setError('');
    setSuccess('');

    try {
      await axios.delete('/api/auth/delete-account', {
        data: {
          username: deleteUsername,
          password: deletePassword,
          email: deleteEmail
        }
      });
      
      setSuccess('账户已成功删除，正在退出...');
      
      // 2秒后自动登出
      setTimeout(() => {
        onClose();
        logout();
      }, 2000);
    } catch (err) {
      setError(err.response?.data?.detail || '注销失败');
    } finally {
      setLoading(false);
    }
  };

  const handlePasswordUpdate = async () => {
    if (!oldPassword || !newPassword) {
      setError('请输入原密码和新密码');
      return;
    }

    if (newPassword !== confirmPassword) {
      setError('两次输入的新密码不一致');
      return;
    }

    if (newPassword.length < 6) {
      setError('新密码至少需要6位');
      return;
    }

    setLoading(true);
    setError('');
    setSuccess('');

    try {
      const token = localStorage.getItem('token');
      await axios.put('http://localhost:3008/api/auth/update-password', {
        old_password: oldPassword,
        new_password: newPassword
      }, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      setSuccess('密码修改成功，请重新登录');
      setOldPassword('');
      setNewPassword('');
      setConfirmPassword('');
      
      // 2秒后自动登出
      setTimeout(() => {
        onClose();
        logout();
      }, 2000);
    } catch (err) {
      setError(err.response?.data?.detail || '修改失败');
    } finally {
      setLoading(false);
    }
  };

  const tabs = [
    { id: 'username', label: '修改用户名', icon: User },
    { id: 'email', label: '修改邮箱', icon: Mail },
    { id: 'password', label: '修改密码', icon: Lock },
    { id: 'delete', label: '注销账户', icon: Trash2 },
  ];

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-full max-w-xl mx-4">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold text-gray-900">用户设置</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        {/* 标签页 */}
        <div className="flex space-x-1 mb-4 border-b">
          {tabs.map((tab) => {
            const Icon = tab.icon;
            return (
              <button
                key={tab.id}
                onClick={() => {
                  setActiveTab(tab.id);
                  setError('');
                  setSuccess('');
                  setShowDeleteConfirm(false);
                }}
                className={`flex items-center space-x-2 px-3 py-2 text-sm font-medium rounded-t-lg ${
                  activeTab === tab.id
                    ? 'border-b-2 border-blue-500 text-blue-600'
                    : 'text-gray-500 hover:text-gray-700'
                }`}
              >
                <Icon className="h-4 w-4" />
                <span>{tab.label}</span>
              </button>
            );
          })}
        </div>

        {/* 错误和成功消息 */}
        {error && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
            <p className="text-sm text-red-600">{error}</p>
          </div>
        )}
        {success && (
          <div className="mb-4 p-3 bg-green-50 border border-green-200 rounded-md">
            <p className="text-sm text-green-600">{success}</p>
          </div>
        )}

        {/* 表单内容 */}
        <div className="space-y-4">
          {/* 修改用户名 */}
          {activeTab === 'username' && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                新用户名
              </label>
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="请输入新用户名"
              />
              <button
                onClick={handleUsernameUpdate}
                disabled={loading}
                className="mt-3 w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:opacity-50"
              >
                {loading ? '修改中...' : '修改用户名'}
              </button>
            </div>
          )}

          {/* 修改邮箱 */}
          {activeTab === 'email' && (
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  当前邮箱: {email}
                </label>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  新邮箱
                </label>
                <input
                  type="email"
                  value={newEmail}
                  onChange={(e) => setNewEmail(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="请输入新邮箱"
                />
              </div>
              <div className="relative">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  密码验证
                </label>
                <input
                  type={showFields.emailPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full px-3 py-2 pr-10 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="请输入密码验证身份"
                />
                <button
                  type="button"
                  onClick={() => toggleShowField('emailPassword')}
                  className="absolute right-2 top-9 text-gray-400 hover:text-gray-600"
                >
                  {showFields.emailPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              </div>
              <button
                onClick={handleEmailUpdate}
                disabled={loading}
                className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:opacity-50"
              >
                {loading ? '修改中...' : '修改邮箱'}
              </button>
            </div>
          )}

          {/* 修改密码 */}
          {activeTab === 'password' && (
            <div className="space-y-4">
              <div className="relative">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  原密码
                </label>
                <input
                  type={showFields.oldPassword ? 'text' : 'password'}
                  value={oldPassword}
                  onChange={(e) => setOldPassword(e.target.value)}
                  className="w-full px-3 py-2 pr-10 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="请输入原密码"
                />
                <button
                  type="button"
                  onClick={() => toggleShowField('oldPassword')}
                  className="absolute right-2 top-9 text-gray-400 hover:text-gray-600"
                >
                  {showFields.oldPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              </div>
              <div className="relative">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  新密码
                </label>
                <input
                  type={showFields.newPassword ? 'text' : 'password'}
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  className="w-full px-3 py-2 pr-10 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="请输入新密码"
                />
                <button
                  type="button"
                  onClick={() => toggleShowField('newPassword')}
                  className="absolute right-2 top-9 text-gray-400 hover:text-gray-600"
                >
                  {showFields.newPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              </div>
              <div className="relative">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  确认新密码
                </label>
                <input
                  type={showFields.confirmPassword ? 'text' : 'password'}
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  className="w-full px-3 py-2 pr-10 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="请再次输入新密码"
                />
                <button
                  type="button"
                  onClick={() => toggleShowField('confirmPassword')}
                  className="absolute right-2 top-9 text-gray-400 hover:text-gray-600"
                >
                  {showFields.confirmPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              </div>
              <button
                onClick={handlePasswordUpdate}
                disabled={loading}
                className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:opacity-50"
              >
                {loading ? '修改中...' : '修改密码'}
              </button>
            </div>
          )}

          {/* 注销账户 - 现在作为独立标签页 */}
          {activeTab === 'delete' && (
            <div className="space-y-4">
              <div className="p-4 bg-red-50 border border-red-200 rounded-md">
                <div className="flex items-start">
                  <AlertTriangle className="h-5 w-5 text-red-400 mt-0.5" />
                  <div className="ml-3">
                    <h3 className="text-sm font-medium text-red-800">
                      危险操作：注销账户
                    </h3>
                    <p className="mt-1 text-sm text-red-700">
                      此操作将永久删除您的账户和所有数据，包括上传的文件，此操作无法撤销。
                    </p>
                  </div>
                </div>
              </div>

              {!showDeleteConfirm ? (
                <button
                  onClick={() => setShowDeleteConfirm(true)}
                  disabled={loading}
                  className="w-full bg-red-600 text-white py-2 px-4 rounded-md hover:bg-red-700 disabled:opacity-50"
                >
                  我理解风险，继续注销
                </button>
              ) : (
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      用户名确认
                    </label>
                    <input
                      type="text"
                      value={deleteUsername}
                      onChange={(e) => setDeleteUsername(e.target.value)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-red-500 focus:border-transparent"
                      placeholder={`请输入您的用户名: ${currentUser?.username}`}
                    />
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      邮箱确认
                    </label>
                    <input
                      type="email"
                      value={deleteEmail}
                      onChange={(e) => setDeleteEmail(e.target.value)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-red-500 focus:border-transparent"
                      placeholder={`请输入您的邮箱: ${currentUser?.email}`}
                    />
                  </div>
                  
                  <div className="relative">
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      密码确认
                    </label>
                    <input
                      type={showFields.deletePassword ? 'text' : 'password'}
                      value={deletePassword}
                      onChange={(e) => setDeletePassword(e.target.value)}
                      className="w-full px-3 py-2 pr-10 border border-gray-300 rounded-md focus:ring-2 focus:ring-red-500 focus:border-transparent"
                      placeholder="请输入您的密码"
                    />
                    <button
                      type="button"
                      onClick={() => toggleShowField('deletePassword')}
                      className="absolute right-2 top-9 text-gray-400 hover:text-gray-600"
                    >
                      {showFields.deletePassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                    </button>
                  </div>
                  
                  <div className="flex space-x-3">
                    <button
                      onClick={handleDeleteAccount}
                      disabled={loading}
                      className="flex-1 bg-red-600 text-white py-2 px-4 rounded-md hover:bg-red-700 disabled:opacity-50"
                    >
                      {loading ? '注销中...' : '确认注销'}
                    </button>
                    <button
                      onClick={() => {
                        setShowDeleteConfirm(false);
                        setDeleteUsername('');
                        setDeletePassword('');
                        setDeleteEmail('');
                      }}
                      disabled={loading}
                      className="flex-1 bg-gray-300 text-gray-700 py-2 px-4 rounded-md hover:bg-gray-400 disabled:opacity-50"
                    >
                      取消
                    </button>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}