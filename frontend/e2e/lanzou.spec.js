import { test, expect } from '@playwright/test'

test.describe('蓝奏云功能', () => {
  test.beforeEach(async ({ page }) => {
    // 登录
    const timestamp = Date.now()
    const email = `lanzou${timestamp}@example.com`
    const password = 'Test123456'

    await page.goto('/register')
    await page.getByPlaceholder('邮箱').fill(email)
    await page.getByPlaceholder('密码 (至少6位)').fill(password)
    await page.getByRole('button', { name: '注册' }).click()
    await expect(page).toHaveURL('/dashboard', { timeout: 10000 })
  })

  test('蓝奏云设置页面访问', async ({ page }) => {
    await page.goto('/settings')

    await expect(page.locator('text=蓝奏云')).toBeVisible({ timeout: 5000 })
  })

  test('蓝奏云浏览器页面访问', async ({ page }) => {
    await page.goto('/lanzou')

    await expect(page.locator('text=蓝奏云浏览器')).toBeVisible({ timeout: 5000 })
  })

  test('分享解析页面访问', async ({ page }) => {
    await page.goto('/share-parse')

    await expect(page.locator('text=链接解析')).toBeVisible({ timeout: 5000 })
  })

  test('同步资源页面访问', async ({ page }) => {
    await page.goto('/sync')

    await expect(page.locator('text=同步')).toBeVisible({ timeout: 5000 })
  })

  test('回收站页面访问', async ({ page }) => {
    await page.goto('/recycle')

    await expect(page.locator('text=回收站')).toBeVisible({ timeout: 5000 })
  })
})
