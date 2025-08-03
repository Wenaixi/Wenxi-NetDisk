/**
 * Wenxi网盘 - API统一配置
 * 作者：Wenxi
 * 功能：提供统一的API客户端配置，处理认证和错误
 */

import axios from 'axios';

// 导入API配置工具
import { getBaseURL } from '../utils/apiConfig';

// 创建统一的axios实例
const api = axios.create({
  baseURL: getBaseURL(),
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

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