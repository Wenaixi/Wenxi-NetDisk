/**
 * Wenxi超级高速网盘 - 登录组件
 * 作者：Wenxi
 * 功能：用户登录和注册界面
 */

import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import { UploadCloud, User, Mail, Lock, Eye, EyeOff } from 'lucide-react';

export default function Login() {
  const [isLogin, setIsLogin] = useState(true);
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
  });
  const [rememberMe, setRememberMe] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const { login, register } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      // 检查网络连接
      if (!navigator.onLine) {
        setError('网络连接断开，请检查网络连接');
        setLoading(false);
        return;
      }

      let result;
      if (isLogin) {
        result = await login(formData.username, formData.password, rememberMe);
        if (result.success) {
          console.log('Wenxi - 登录成功，准备跳转到dashboard');
          navigate('/dashboard');
        } else {
          console.log('Wenxi - 登录失败:', result.error);
          setError(result.error);
        }
      } else {
        result = await register(formData.username, formData.email, formData.password);
        if (result.success) {
          console.log('Wenxi - 注册成功');
          // 注册成功后自动切换到登录模式
          setIsLogin(true);
          setFormData({
            username: formData.username,
            email: '',
            password: ''
          });
          setError('注册成功！请登录');
        } else {
          console.log('Wenxi - 注册失败:', result.error);
          setError(result.error);
        }
      }
    } catch (err) {
      console.error('Wenxi - 未预期的错误:', err);
      setError('操作失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <div className="mx-auto h-12 w-12 flex items-center justify-center rounded-full bg-blue-600">
            <UploadCloud className="h-8 w-8 text-white" />
          </div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Wenxi网盘
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            {isLogin ? '登录您的账户' : '创建新账户'}
          </p>
        </div>

        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="rounded-md shadow-sm -space-y-px">
            <div>
              <label htmlFor="username" className="sr-only">用户名</label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <User className="h-5 w-5 text-gray-400" />
                </div>
                <input
                  id="username"
                  name="username"
                  type="text"
                  required
                  className="appearance-none rounded-none relative block w-full px-3 py-2 pl-10 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
                  placeholder="用户名"
                  value={formData.username}
                  onChange={handleChange}
                />
              </div>
            </div>

            {!isLogin && (
              <div>
                <label htmlFor="email" className="sr-only">邮箱</label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Mail className="h-5 w-5 text-gray-400" />
                  </div>
                  <input
                    id="email"
                    name="email"
                    type="email"
                    required={!isLogin}
                    className="appearance-none rounded-none relative block w-full px-3 py-2 pl-10 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
                    placeholder="邮箱地址"
                    value={formData.email}
                    onChange={handleChange}
                  />
                </div>
              </div>
            )}

            <div>
              <label htmlFor="password" className="sr-only">密码</label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <Lock className="h-5 w-5 text-gray-400" />
                </div>
                <input
                  id="password"
                  name="password"
                  type={showPassword ? 'text' : 'password'}
                  required
                  className="appearance-none rounded-none relative block w-full px-3 py-2 pl-10 pr-10 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
                  placeholder="密码"
                  value={formData.password}
                  onChange={handleChange}
                />
                <div className="absolute inset-y-0 right-0 pr-3 flex items-center">
                  <button
                    type="button"
                    className="text-gray-400 hover:text-gray-600"
                    onClick={() => setShowPassword(!showPassword)}
                  >
                    {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                  </button>
                </div>
              </div>
            </div>
          </div>

          {error && (
            <div className="text-red-600 text-sm text-center">{error}</div>
          )}

          {isLogin && (
            <div className="flex items-center">
              <input
                id="remember-me"
                name="remember-me"
                type="checkbox"
                checked={rememberMe}
                onChange={(e) => setRememberMe(e.target.checked)}
                className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
              />
              <label htmlFor="remember-me" className="ml-2 block text-sm text-gray-900">
                记住我（72小时内免登录）
              </label>
            </div>
          )}

          <div>
            <button
              type="submit"
              disabled={loading}
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
            >
              {loading ? '处理中...' : (isLogin ? '登录' : '注册')}
            </button>
          </div>

          <div className="text-center">
            <button
              type="button"
              className="text-blue-600 hover:text-blue-500 text-sm"
              onClick={() => {
                setIsLogin(!isLogin);
                setError('');
              }}
            >
              {isLogin ? '没有账户？立即注册' : '已有账户？立即登录'}
            </button>
          </div>
        </form>

        {/* 底部开源信息 - Wenxi网盘开源声明 */}
        <footer className="absolute bottom-0 left-0 right-0 bg-white/80 backdrop-blur-sm border-t border-gray-200">
          <div className="max-w-7xl mx-auto py-4 px-4 sm:px-6 lg:px-8">
            <div className="text-center text-xs text-gray-500 space-y-1">
              <p>
                本项目已在 <a href="https://github.com/Wenaixi/Wenxi-Network-Disk" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:text-blue-800">GitHub</a> 用 <a href="https://opensource.org/licenses/MIT" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:text-blue-800">MIT协议</a> 全面开源
              </p>
              <p>
                作者：<span className="font-semibold">Wenxi</span> | 版本号：<span className="font-semibold">v1.1.1</span>
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
    </div>
  );
}