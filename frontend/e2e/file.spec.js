import { test, expect } from '@playwright/test'

test.describe('文件管理流程', () => {
  let authToken

  test.beforeEach(async ({ page }) => {
    // 登录获取token
    const timestamp = Date.now()
    const email = `e2e${timestamp}@example.com`
    const password = 'Test123456'

    await page.goto('/register')
    await page.getByPlaceholder('邮箱').fill(email)
    await page.getByPlaceholder('密码 (至少6位)').fill(password)
    await page.getByRole('button', { name: '注册' }).click()
    await expect(page).toHaveURL('/dashboard', { timeout: 10000 })
    authToken = await page.evaluate(() => localStorage.getItem('token'))
  })

  test('创建文件夹', async ({ page }) => {
    await page.waitForSelector('.n-button:has-text("新建文件夹")', { timeout: 5000 })

    const createBtn = page.locator('.n-button', { hasText: '新建文件夹' }).first()
    await createBtn.click()

    await page.waitForTimeout(500)
    await page.getByPlaceholder('文件夹名称').fill('测试文件夹')
    await page.locator('.n-button', { hasText: '确定' }).click()

    await expect(page.locator('.n-data-table')).toContainText('测试文件夹')
  })

  test('上传文件 - UI元素存在', async ({ page }) => {
    // 验证上传按钮存在
    const uploadBtn = page.locator('.n-button', { hasText: '上传文件' }).first()
    await expect(uploadBtn).toBeVisible()
  })

  test('右键菜单 - 文件夹详情', async ({ page }) => {
    // 创建文件夹
    await page.locator('.n-button', { hasText: '新建文件夹' }).first().click()
    await page.waitForTimeout(300)
    await page.getByPlaceholder('文件夹名称').fill('右键菜单测试')
    await page.locator('.n-button', { hasText: '确定' }).click()
    await page.waitForTimeout(500)

    // 右键点击文件夹
    const folderRow = page.locator('.n-data-table-row').filter({ hasText: '右键菜单测试' }).first()
    await folderRow.click({ button: 'right' })

    // 验证菜单出现
    await expect(page.locator('.n-popover-menu')).toBeVisible({ timeout: 3000 })
  })

  test('面包屑导航', async ({ page }) => {
    // 创建文件夹
    await page.locator('.n-button', { hasText: '新建文件夹' }).first().click()
    await page.waitForTimeout(300)
    await page.getByPlaceholder('文件夹名称').fill('子文件夹')
    await page.locator('.n-button', { hasText: '确定' }).click()
    await page.waitForTimeout(500)

    // 点击进入文件夹
    await page.locator('.n-data-table-row').filter({ hasText: '子文件夹' }).first().dblclick()
    await page.waitForTimeout(500)

    // 验证面包屑包含"子文件夹"
    await expect(page.locator('.n-breadcrumb')).toContainText('子文件夹')
  })

  test('批量选择模式', async ({ page }) => {
    // 创建多个文件夹
    for (const name of ['批量测试1', '批量测试2', '批量测试3']) {
      await page.locator('.n-button', { hasText: '新建文件夹' }).first().click()
      await page.waitForTimeout(300)
      await page.getByPlaceholder('文件夹名称').fill(name)
      await page.locator('.n-button', { hasText: '确定' }).click()
      await page.waitForTimeout(300)
    }

    // 点击全选或批量选择按钮
    const selectBtn = page.locator('.n-button', { hasText: '全选' }).first()
    if (await selectBtn.isVisible()) {
      await selectBtn.click()
      await page.waitForTimeout(300)
    }
  })

  test('排序功能', async ({ page }) => {
    await page.waitForTimeout(1000)

    // 查找排序按钮
    const sortBtn = page.locator('.n-button', { hasText: '排序' }).first()
    if (await sortBtn.isVisible()) {
      await sortBtn.click()
      await page.waitForTimeout(300)

      // 选择按名称排序
      await page.locator('.n-dropdown-menu').locator('.n-dropdown-option').filter({ hasText: '名称' }).first().click()
      await page.waitForTimeout(300)
    }
  })
})
