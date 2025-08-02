/**
 * @文件：setup.js
 * @作者：Wenxi
 * @功能：测试环境初始化配置
 * @创建时间：2025
 */

import '@testing-library/jest-dom';
import { vi } from 'vitest';

// Mock matchMedia for antd components
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(), // deprecated
    removeListener: vi.fn(), // deprecated
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});