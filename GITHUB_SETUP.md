# GitHub绑定指南 - Wenxi NetDisk

## 步骤1: 创建GitHub仓库

1. 打开 [GitHub.com](https://github.com) 并登录你的账户
2. 点击右上角的 "+" 图标，选择 "New repository"
3. 填写仓库信息：
   - Repository name: `Wenxi-NetDisk` (或你喜欢的名字)
   - Description: `一个功能完整的网盘系统，支持分片上传、断点续传和权限管理`
   - 选择 Public 或 Private
   - 不要勾选 "Initialize this repository with a README"
4. 点击 "Create repository"

## 步骤2: 关联本地仓库到GitHub

在PowerShell中运行以下命令：

```powershell
# 关联远程仓库 (替换为你的仓库URL)
git remote add origin https://github.com/Wenaixi/Wenxi-NetDisk.git

# 验证远程仓库
git remote -v

# 推送代码到GitHub
git push -u origin master
```

如果遇到分支名称问题，使用：
```powershell
git branch -M master
git push -u origin master
```

## 步骤3: 验证绑定成功

```powershell
# 查看远程分支
git branch -r

# 拉取最新代码测试
git pull origin master
```

## 可选：SSH密钥配置(推荐)

为了避免每次推送都输入密码，建议配置SSH密钥：

1. 生成SSH密钥：
```powershell
ssh-keygen -t rsa -b 4096 -C "你的邮箱@example.com"
```

2. 添加SSH密钥到ssh-agent：
```powershell
# 启动ssh-agent
Get-Service ssh-agent | Set-Service -StartupType Manual
Start-Service ssh-agent

# 添加私钥
ssh-add ~/.ssh/id_rsa
```

3. 复制公钥到GitHub：
```powershell
# 查看公钥内容
cat ~/.ssh/id_rsa.pub
```

4. 在GitHub设置中添加SSH密钥

5. 使用SSH地址关联远程仓库：
```powershell
git remote set-url origin git@github.com:你的用户名/Wenxi-NetDisk.git
```

## 常用Git命令

```powershell
# 查看状态
git status

# 添加修改
git add .

# 提交修改
git commit -m "描述你的修改"

# 推送到GitHub
git push origin master

# 拉取最新代码
git pull origin master
```

## 注意事项

- 确保.gitignore文件已正确配置，避免上传敏感文件
- 定期检查并清理临时文件
- 提交前运行测试确保代码质量

## 技术支持

如果遇到任何问题，请检查：
1. 网络连接是否正常
2. GitHub账户权限是否正确
3. 本地Git配置是否正确 (`git config --list`)

---
*Author: Wenxi*
*Created: $(Get-Date -Format 'yyyy-MM-dd')*