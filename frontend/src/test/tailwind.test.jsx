/**
 * Wenxi网盘 - Tailwind CSS验证测试
 * 作者：Wenxi
 * 功能：验证Tailwind CSS配置是否正确工作
 * 说明：测试Tailwind类名是否正确应用样式
 */

import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import '../index.css';

describe('Tailwind CSS 配置验证', () => {
  test('Tailwind基础类名应该正确应用', () => {
    const { container } = render(
      <div className="bg-primary-500 text-white p-4 rounded-lg">
        <h1 className="text-2xl font-bold">Wenxi网盘测试</h1>
        <p className="text-sm opacity-75">Tailwind CSS测试内容</p>
      </div>
    );
    
    const div = container.firstChild;
    expect(div).toHaveClass('bg-primary-500');
    expect(div).toHaveClass('text-white');
    expect(div).toHaveClass('p-4');
    expect(div).toHaveClass('rounded-lg');
  });

  test('自定义颜色应该可用', () => {
    render(
      <div className="bg-wenxi-500 text-wenxi-50">
        <span data-testid="color-test">Wenxi品牌色测试</span>
      </div>
    );
    
    const element = screen.getByTestId('color-test');
    expect(element).toBeInTheDocument();
  });

  test('响应式类名应该可用', () => {
    const { container } = render(
      <div className="hidden md:block lg:flex">
        响应式测试内容
      </div>
    );
    
    const div = container.firstChild;
    expect(div).toHaveClass('hidden');
    expect(div).toHaveClass('md:block');
    expect(div).toHaveClass('lg:flex');
  });

  test('动画类名应该可用', () => {
    const { container } = render(
      <div className="animate-fade-in">
        动画测试内容
      </div>
    );
    
    const div = container.firstChild;
    expect(div).toHaveClass('animate-fade-in');
  });
});