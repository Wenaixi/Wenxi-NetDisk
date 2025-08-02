/**
 * Wenxi网盘 - 认证上下文
 * 作者：Wenxi
 * 功能：管理用户认证状态，提供全局认证功能
 */

import React, { createContext, useContext, useState, useEffect } from 'react';
import axios from 'axios';

const AuthContext = createContext();

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth必须在AuthProvider中使用');
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      fetchUser();
    } else {
      setLoading(false);
    }
  }, []);

  const fetchUser = async () => {
    try {
      const response = await axios.get('/api/auth/me');
      setUser(response.data);
      setIsAuthenticated(true);
    } catch (error) {
      localStorage.removeItem('token');
      delete axios.defaults.headers.common['Authorization'];
    } finally {
      setLoading(false);
    }
  };

  const login = async (username, password, rememberMe = false) => {
    try {
      console.log('Wenxi - 开始登录请求，用户名:', username, '记住我:', rememberMe);
      
      const formData = new FormData();
      formData.append('username', username);
      formData.append('password', password);
      if (rememberMe) {
        formData.append('client_id', 'remember_me');
      }

      const response = await axios.post('/api/auth/login', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        timeout: 10000, // 10秒超时
      });

      console.log('Wenxi - 登录响应:', response);
      
      const { access_token } = response.data;
      localStorage.setItem('token', access_token);
      axios.defaults.headers.common['Authorization'] = `Bearer ${access_token}`;
      
      await fetchUser();
      console.log('Wenxi - 登录成功');
      return { success: true };
    } catch (error) {
      console.error('Wenxi - 登录错误详情:', {
        message: error.message,
        response: error.response?.data,
        status: error.response?.status,
        config: error.config
      });
      
      let errorMessage = '登录失败';
      if (error.code === 'ECONNREFUSED') {
        errorMessage = '无法连接到服务器，请检查网络连接';
      } else if (error.response?.status === 401) {
        errorMessage = error.response?.data?.detail || '用户名或密码错误';
      } else if (error.response?.status >= 500) {
        errorMessage = '服务器错误，请稍后重试';
      } else if (error.message?.includes('timeout')) {
        errorMessage = '请求超时，请检查网络连接';
      }
      
      return { 
        success: false, 
        error: errorMessage
      };
    }
  };

  const register = async (username, email, password) => {
    try {
      const response = await axios.post('/api/auth/register', {
        username,
        email,
        password,
      });
      
      return { success: true, data: response.data };
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.detail || '注册失败' 
      };
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    delete axios.defaults.headers.common['Authorization'];
    setUser(null);
    setIsAuthenticated(false);
  };

  const value = {
    user,
    isAuthenticated,
    loading,
    login,
    register,
    logout,
    setUser,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};