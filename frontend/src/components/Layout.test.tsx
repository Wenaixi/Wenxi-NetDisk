/**
 * @文件：Layout.test.tsx
 * @作者：Wenxi
 * @功能：Layout组件的单元测试
 * @创建时间：2025
 */

import React from 'react';
import { describe, test, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import Layout from './Layout';

// 简单的测试用例，确保组件能够正常渲染
describe('Layout Component', () => {
  test('renders without crashing', () => {
    render(
      <MemoryRouter>
        <Layout>
          <div>测试内容</div>
        </Layout>
      </MemoryRouter>
    );
    
    expect(screen.getByText('Wenxi网盘')).toBeInTheDocument();
    expect(screen.getByText('测试内容')).toBeInTheDocument();
  });

  test('displays navigation menu items', () => {
    render(
      <MemoryRouter>
        <Layout>
          <div>测试内容</div>
        </Layout>
      </MemoryRouter>
    );
    
    expect(screen.getByText('文件管理')).toBeInTheDocument();
    expect(screen.getByText('上传文件')).toBeInTheDocument();
  });
});