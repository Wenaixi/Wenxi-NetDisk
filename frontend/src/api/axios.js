/**
 * Wenxi网盘 - Axios配置
 * 作者：Wenxi
 * 功能：配置HTTP客户端，设置基础URL、超时、拦截器等
 */

import axios from 'axios';

// 创建axios实例
const api = axios.create({
  baseURL: 'http://localhost:3008',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 检查token有效性
const checkTokenValidity = async () => {
  const token = localStorage.getItem('token');
  if (!token) return false;
  
  try {
    await api.get('/api/auth/validate-token');
    return true;
  } catch (error) {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      return false;
    }
    return true;
  }
};

// 页面刷新时检查token
if (typeof window !== 'undefined' && localStorage.getItem('token')) {
  checkTokenValidity();
}

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    console.log('Wenxi - API请求:', config.method.toUpperCase(), config.url);
    return config;
  },
  (error) => {
    console.error('Wenxi - 请求错误:', error);
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    console.log('Wenxi - API响应:', response.status, response.config.url);
    return response;
  },
  (error) => {
    console.error('Wenxi - 响应错误:', error.response?.status, error.message);
    
    // 处理401错误（未授权）
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      // 只在非登录页面跳转
      if (!window.location.pathname.includes('/login')) {
        window.location.href = '/login';
      }
    }
    
    return Promise.reject(error);
  }
);

export default api;