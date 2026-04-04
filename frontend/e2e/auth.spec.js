import { test, expect } from '@playwright/test'

test.describe('认证流程', () => {
  test('用户注册 - 成功', async ({ page }) => {
    const timestamp = Date.now()
    const email = `test${timestamp}@example.com`
    const password = 'Test123456'

    await page.goto('/register')

    await page.getByPlaceholder('邮箱').fill(email)
    await page.getByPlaceholder('密码 (至少6位)').fill(password)
    await page.getByRole('button', { name: '注册' }).click()

    await expect(page).toHaveURL('/dashboard', { timeout: 10000 })
  })

  test('用户登录 - 成功', async ({ page }) => {
    await page.goto('/login')

    await page.getByPlaceholder('邮箱').fill('test@example.com')
    await page.getByPlaceholder('密码').fill('Test123456')
    await page.getByRole('button', { name: '登录' }).click()

    await expect(page).toHaveURL('/dashboard', { timeout: 10000 })
  })

  test('用户登录 - 错误密码', async ({ page }) => {
    await page.goto('/login')

    await page.getByPlaceholder('邮箱').fill('test@example.com')
    await page.getByPlaceholder('密码').fill('wrongpassword')
    await page.getByRole('button', { name: '登录' }).click()

    await expect(page.locator('.n-notification')).toBeVisible({ timeout: 5000 })
  })

  test('导航守卫 - 未登录访问仪表盘重定向', async ({ page }) => {
    await page.goto('/dashboard')

    await expect(page).toHaveURL(/\/login/)
  })
})
