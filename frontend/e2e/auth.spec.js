import { test, expect } from '@playwright/test'

test.describe('认证流程', () => {
  test('用户注册 - 成功', async ({ page }) => {
    const timestamp = Date.now()
    const email = `test${timestamp}@example.com`
    const password = 'Test123456'

    await page.goto('/register')
    await page.waitForLoadState('networkidle')

    await page.locator('input[type="text"], input[placeholder*="邮箱"]').first().fill(email)
    await page.locator('input[type="password"]').nth(0).fill(password)
    await page.locator('input[type="password"]').nth(1).fill(password)
    await page.locator('button:has-text("注册")').click()

    await expect(page).toHaveURL('/dashboard', { timeout: 15000 })
  })

  test('用户登录 - 成功', async ({ page }) => {
    await page.goto('/login')
    await page.waitForLoadState('networkidle')

    await page.locator('input[placeholder*="邮箱"]').first().fill('test@example.com')
    await page.locator('input[placeholder*="密码"]').first().fill('Test123456')
    await page.locator('button:has-text("登录")').click()

    await expect(page).toHaveURL('/dashboard', { timeout: 15000 })
  })

  test('用户登录 - 错误密码', async ({ page }) => {
    await page.goto('/login')
    await page.waitForLoadState('networkidle')

    await page.locator('input[placeholder*="邮箱"]').first().fill('test@example.com')
    await page.locator('input[placeholder*="密码"]').first().fill('wrongpassword')
    await page.locator('button:has-text("登录")').click()

    await expect(page.locator('.n-notification, .n-message, [role="alert"]')).toBeVisible({ timeout: 8000 })
  })

  test('导航守卫 - 未登录访问仪表盘重定向', async ({ page }) => {
    await page.goto('/dashboard')
    await page.waitForLoadState('networkidle')

    await expect(page).toHaveURL(/\/login/)
  })
})
