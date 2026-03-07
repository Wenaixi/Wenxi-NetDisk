/**
 * Wenxi网盘 - API配置工具 (Docker环境专用)
 * 作者：Wenxi
 * 功能：Docker容器内部使用，直接访问后端服务
 */

/**
 * 获取动态API基础URL
 * Docker环境：使用相对路径，通过nginx代理到后端
 * @returns {string} API基础URL
 */
export const getBaseURL = () => {
  // Docker环境下，前端通过nginx代理访问后端
  // 使用相对路径，让nginx处理代理
  return '';
};

/**
 * 获取完整的API端点URL
 * @param {string} endpoint - API端点路径（如 '/api/auth/login'）
 * @returns {string} 完整的API URL
 */
export const getApiEndpoint = (endpoint) => {
  // Docker环境下直接返回endpoint，nginx会代理到后端
  return endpoint;
};

/**
 * 检查当前是否为开发环境
 * @returns {boolean} 是否为开发环境
 */
export const isDevEnvironment = () => {
  return false;
};

/**
 * 获取WebSocket地址（用于未来实时功能）
 * @returns {string} WebSocket URL
 */
export const getWebSocketURL = () => {
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
  const hostname = window.location.hostname;
  const port = window.location.port;
  
  return `${protocol}://${hostname}${port ? ':' + port : ''}`;
};
