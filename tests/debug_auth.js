/**
 * Wenxi网盘 - 认证调试脚本
 * 作者：Wenxi
 * 功能：调试认证流程和token持久化问题
 */

// 浏览器端调试代码
const debugAuth = {
    // 检查localStorage中的token
    checkLocalStorage() {
        const token = localStorage.getItem('token');
        console.log('LocalStorage token:', token ? '存在' : '不存在');
        if (token) {
            console.log('Token长度:', token.length);
            console.log('Token内容:', token.substring(0, 50) + '...');
        }
        return token;
    },

    // 检查axios默认头
    checkAxiosHeaders() {
        console.log('Axios默认头:', axios.defaults.headers.common);
        console.log('Authorization头:', axios.defaults.headers.common['Authorization']);
    },

    // 模拟页面刷新
    simulateRefresh() {
        console.log('=== 模拟页面刷新 ===');
        
        // 1. 检查localStorage
        const token = this.checkLocalStorage();
        
        // 2. 检查axios头
        this.checkAxiosHeaders();
        
        // 3. 尝试调用/me接口
        if (token) {
            axios.get('/api/auth/me')
                .then(response => {
                    console.log('/me接口成功:', response.data);
                })
                .catch(error => {
                    console.error('/me接口失败:', {
                        status: error.response?.status,
                        message: error.message,
                        response: error.response?.data
                    });
                });
        } else {
            console.log('没有token，跳过/me接口调用');
        }
    },

    // 检查浏览器开发者工具
    checkDevTools() {
        console.log('=== 浏览器调试检查 ===');
        console.log('1. 打开开发者工具 (F12)');
        console.log('2. 切换到Application/Storage标签');
        console.log('3. 查看Local Storage -> http://localhost:5173');
        console.log('4. 检查是否有token键值对');
        console.log('5. 切换到Network标签');
        console.log('6. 刷新页面，查看/auth/me请求');
        console.log('7. 检查请求头是否包含Authorization: Bearer <token>');
    }
};

// 在浏览器控制台运行以下代码
console.log('Wenxi网盘认证调试工具已加载');
debugAuth.simulateRefresh();