/**
 * @文件：vitest.config.js
 * @作者：Wenxi
 * @功能：Vitest测试框架配置
 * @创建时间：2025
 */

import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: './src/test/setup.js',
  },
});