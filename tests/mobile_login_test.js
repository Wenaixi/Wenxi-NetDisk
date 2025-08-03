/**
 * Wenxi网盘 - 移动端登录测试
 * 作者：Wenxi
 * 功能：测试移动端登录功能是否正常
 */

// 移动端登录测试函数
function testMobileLogin() {
  console.log('Wenxi - 开始移动端登录测试...');
  
  // 模拟移动端环境
  const originalLocation = window.location;
  
  // 测试1：开发环境（localhost）
  Object.defineProperty(window, 'location', {
    value: {
      hostname: 'localhost',
      port: '5173',
      protocol: 'http:',
      href: 'http://localhost:5173'
    },
    writable: true
  });
  
  // 动态导入测试
  import('../src/utils/apiConfig.js').then(({ getBaseURL }) => {
    const devURL = getBaseURL();
    console.log('Wenxi - 开发环境测试:', devURL === 'http://localhost:3008' ? '✅ 通过' : '❌ 失败');
    
    // 测试2：移动端环境（IP地址）
    Object.defineProperty(window, 'location', {
      value: {
        hostname: '192.168.1.100',
        port: '5173',
        protocol: 'http:',
        href: 'http://192.168.1.100:5173'
      },
      writable: true
    });
    
    const mobileURL = getBaseURL();
    console.log('Wenxi - 移动端环境测试:', mobileURL === 'http://192.168.1.100:5173' ? '✅ 通过' : '❌ 失败');
    
    // 测试3：生产环境（无端口）
    Object.defineProperty(window, 'location', {
      value: {
        hostname: 'wenxi-disk.com',
        port: '',
        protocol: 'https:',
        href: 'https://wenxi-disk.com'
      },
      writable: true
    });
    
    const prodURL = getBaseURL();
    console.log('Wenxi - 生产环境测试:', prodURL === 'https://wenxi-disk.com' ? '✅ 通过' : '❌ 失败');
    
    // 恢复原始location
    window.location = originalLocation;
    
    console.log('Wenxi - 移动端登录测试完成！');
  }).catch(err => {
    console.error('Wenxi - 测试失败:', err);
  });
}

// 自动执行测试
if (typeof window !== 'undefined') {
  testMobileLogin();
}

// 导出测试函数供其他模块使用
export { testMobileLogin };