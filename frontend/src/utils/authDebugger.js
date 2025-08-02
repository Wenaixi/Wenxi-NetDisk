/**
 * Wenxi网盘 - 认证调试工具
 * 作者：Wenxi
 * 功能：调试和修复认证相关问题
 */

class AuthDebugger {
  constructor() {
    this.debugMode = localStorage.getItem('debug_auth') === 'true';
  }

  log(message, ...args) {
    if (this.debugMode) {
      console.log(`[Wenxi Auth Debug] ${message}`, ...args);
    }
  }

  error(message, ...args) {
    console.error(`[Wenxi Auth Error] ${message}`, ...args);
  }

  // 检查浏览器存储
  checkStorage() {
    const token = localStorage.getItem('token');
    const user = localStorage.getItem('user');
    
    this.log('存储检查:', {
      token: token ? `存在 (${token.length} chars)` : '不存在',
      user: user ? `存在 (${user.length} chars)` : '不存在'
    });

    return { token, user };
  }

  // 检查网络请求
  checkNetwork() {
    const baseURL = window.location.origin.includes('5173') 
      ? 'http://localhost:3008' 
      : window.location.origin;
    
    this.log('网络检查:', {
      origin: window.location.origin,
      baseURL: baseURL,
      apiEndpoint: `${baseURL}/api/auth/me`
    });

    return baseURL;
  }

  // 测试认证API
  async testAuthAPI() {
    const token = localStorage.getItem('token');
    if (!token) {
      this.log('没有token，跳过API测试');
      return null;
    }

    try {
      const response = await fetch('/api/auth/me', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      this.log('API测试结果:', {
        status: response.status,
        statusText: response.statusText,
        ok: response.ok
      });

      if (response.ok) {
        const data = await response.json();
        this.log('用户数据:', data);
        return data;
      } else {
        this.error('API测试失败:', response.status, response.statusText);
        return null;
      }
    } catch (error) {
      this.error('API测试错误:', error);
      return null;
    }
  }

  // 完整调试报告
  async generateReport() {
    this.log('=== Wenxi网盘认证调试报告 ===');
    
    const storage = this.checkStorage();
    const network = this.checkNetwork();
    const userData = await this.testAuthAPI();

    const report = {
      timestamp: new Date().toISOString(),
      storage,
      network,
      userData,
      userAgent: navigator.userAgent,
      localStorageAvailable: !!window.localStorage,
      cookiesEnabled: navigator.cookieEnabled
    };

    this.log('调试报告:', report);
    return report;
  }

  // 清除认证数据
  clearAuth() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    delete axios.defaults.headers.common['Authorization'];
    this.log('认证数据已清除');
  }

  // 设置调试模式
  setDebugMode(enabled) {
    if (enabled) {
      localStorage.setItem('debug_auth', 'true');
      this.debugMode = true;
    } else {
      localStorage.removeItem('debug_auth');
      this.debugMode = false;
    }
  }
}

// 创建全局调试器实例
window.authDebugger = new AuthDebugger();

// 快捷调试命令
window.debugAuth = () => window.authDebugger.generateReport();
window.clearAuth = () => window.authDebugger.clearAuth();

// 在控制台输出使用说明
console.log(`
Wenxi网盘认证调试工具已加载！

使用方法：
1. debugAuth() - 生成完整认证调试报告
2. clearAuth() - 清除所有认证数据
3. 启用调试模式：localStorage.setItem('debug_auth', 'true')
4. 禁用调试模式：localStorage.removeItem('debug_auth')

检查浏览器：
- 打开开发者工具 (F12)
- 查看Application/Storage中的Local Storage
- 检查Network标签中的/auth/me请求
`);