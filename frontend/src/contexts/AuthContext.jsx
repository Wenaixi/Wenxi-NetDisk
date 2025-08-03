/**
 * Wenxi网盘 - 认证上下文
 * 作者：Wenxi
 * 功能：管理用户认证状态，提供全局认证功能
 */

import React, { createContext, useContext, useState, useEffect } from 'react';
import axios from 'axios';

// 导入API配置工具
import { getBaseURL } from '../utils/apiConfig';

// 配置axios实例
const api = axios.create({
  baseURL: getBaseURL(),
  timeout: 10000,
});

// 配置请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

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
    console.log('Wenxi - 页面刷新，检查token:', token ? '存在' : '不存在');
    
    if (token) {
      fetchUser();
    } else {
      console.log('Wenxi - 没有找到token，设置为未认证');
      setLoading(false);
    }
  }, []);

  const fetchUser = async (retryCount = 0) => {
    const maxRetries = 2;
    const baseDelay = 1000; // 1秒基础延迟
    
    try {
      console.log('Wenxi - 开始验证用户身份...');
      const response = await api.get('/api/auth/me', {
        timeout: 5000, // 5秒超时
        headers: {
          'Cache-Control': 'no-cache',
        }
      });
      console.log('Wenxi - 用户身份验证成功:', response.data);
      setUser(response.data);
      setIsAuthenticated(true);
    } catch (error) {
      console.error('Wenxi - 用户身份验证失败:', {
        status: error.response?.status,
        message: error.message,
        response: error.response?.data,
        code: error.code,
        isNetworkError: !error.response
      });
      
      // 网络错误或CORS问题：重试机制
      if (!error.response && retryCount < maxRetries) {
        const delay = baseDelay * Math.pow(2, retryCount); // 指数退避
        console.log(`Wenxi - 网络错误，${delay}ms后重试(${retryCount + 1}/${maxRetries})`);
        
        setTimeout(() => {
          fetchUser(retryCount + 1);
        }, delay);
        return; // 不执行finally的setLoading(false)
      }
      
      // 只有明确认证失败(401)或重试耗尽才清除token
      if (error.response?.status === 401 || retryCount >= maxRetries) {
        console.log('Wenxi - 认证失败或重试耗尽，清除token');
        localStorage.removeItem('token');
        delete axios.defaults.headers.common['Authorization'];
      } else {
        console.log('Wenxi - 临时错误，保留token');
      }
    } finally {
      if (retryCount === 0 || retryCount >= maxRetries) {
        setLoading(false);
      }
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

      const response = await api.post('/api/auth/login', formData, {
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
      const response = await api.post('/api/auth/register', {
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