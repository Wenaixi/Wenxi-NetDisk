/**
 * Wenxi网盘 - API配置工具
 * 作者：Wenxi
 * 功能：统一管理API地址配置，支持移动端和开发环境
 */

/**
 * 获取动态API基础URL
 * 支持开发环境（localhost:3008）和生产环境（当前主机）
 * @returns {string} API基础URL
 */
export const getBaseURL = () => {
  // Wenxi - 动态获取API地址，支持移动端访问
  const hostname = window.location.hostname;
  const port = window.location.port;
  
  // 如果是开发环境（localhost或5173端口），使用localhost:3008
  if (hostname === 'localhost' || hostname === '127.0.0.1' || port === '5173') {
    return 'http://localhost:3008';
  }
  
  // 生产环境或移动端访问，使用当前主机地址
  const protocol = window.location.protocol;
  return `${protocol}//${hostname}${port ? ':' + port : ''}`;
};

/**
 * 获取完整的API端点URL
 * @param {string} endpoint - API端点路径（如 '/api/auth/login'）
 * @returns {string} 完整的API URL
 */
export const getApiEndpoint = (endpoint) => {
  return `${getBaseURL()}${endpoint}`;
};

/**
 * 检查当前是否为开发环境
 * @returns {boolean} 是否为开发环境
 */
export const isDevEnvironment = () => {
  const hostname = window.location.hostname;
  const port = window.location.port;
  return hostname === 'localhost' || hostname === '127.0.0.1' || port === '5173';
};

/**
 * 获取WebSocket地址（用于未来实时功能）
 * @returns {string} WebSocket URL
 */
export const getWebSocketURL = () => {
  const baseURL = getBaseURL();
  const protocol = baseURL.startsWith('https') ? 'wss' : 'ws';
  const hostname = window.location.hostname;
  const port = window.location.port;
  
  return `${protocol}://${hostname}${port ? ':' + port : ''}`;
};