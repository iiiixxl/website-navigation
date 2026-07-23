# Site.Navigation

个人站点导航页，纯静态 HTML，无需数据库和构建工具。用于集中管理常用环境地址、服务器与数据库账号，以及 AI 提问模板。

## 功能

### 站点导航（`/`）

- 按分类分块展示站点卡片
- 有账号/密码的卡片：点击展开查看网址与凭据，支持一键复制
- 无账号/密码的卡片：点击直接打开（仅 `http://` / `https://` 开头地址）
- 纯 IP 或非网页地址：只展示文本，不会误跳转
- 顶部搜索：按名称或网址过滤
- 分类模块可折叠，折叠状态自动保存在浏览器本地

### AI 提问模板（`/AI`）

- 左侧：分类 + 模板列表
- 右侧：原文展示（换行、空格、缩进原样保留）+ 一键复制
- 支持搜索与分类折叠

## 文件说明

```
Site.Navigation/
├── index.html        # 站点导航页
├── data.js           # 站点数据（本地维护，已加入 .gitignore）
├── AI/
│   ├── index.html    # AI 提问模板页（访问路径 /AI）
│   └── prompts.js    # 模板数据（本地维护，已加入 .gitignore）
├── web.config        # IIS 部署配置（可选，已忽略）
└── .gitignore
```

## 快速开始

1. 克隆或下载项目
2. 创建 `data.js`、`AI/prompts.js`（参考下方数据格式）
3. 用浏览器打开 `index.html`，或通过 IIS / 任意静态服务器访问
4. AI 页面地址：`/AI` 或 `/AI/`（IIS 会自动打开目录下的 `index.html`）

> `data.js`、`AI/prompts.js` 默认不提交到 Git，首次使用需自行创建。

## 维护数据

编辑 `data.js` 中的 `SITE_DATA`，保存后刷新页面即可生效。

```javascript
window.SITE_DATA = {
  categories: [
    {
      title: "应用环境",       // 分类标题（可点击折叠）
      sites: [
        {
          name: "开发环境",    // 卡片显示名称
          url: "http://192.168.1.1:8081/login",
          username: "admin",   // 可选，无则留 ""
          password: "******",  // 可选，无则留 ""
          note: "开发环境登录", // 可选备注
        },
      ],
    },
  ],
};
```

### 卡片行为说明

| 类型 | 条件 | 行为 |
|------|------|------|
| 快捷入口 | 无 `username` 且 无 `password`，URL 以 `http(s)://` 开头 | 点击卡片直接打开链接 |
| 凭据卡片 | 有 `username` 或 `password` | 点击卡片/眼睛按钮展开，可复制账号密码 |
| 纯地址 | URL 不以 `http(s)://` 开头（如 IP） | 仅展示，不可跳转 |

同一分类内，无密码的快捷入口卡片高度为有密码卡片的 1/3，一列最多叠放 3 个。

## 维护 AI 模板

编辑 `AI/prompts.js` 中的 `PROMPT_DATA`。`content` 建议用反引号写多行，展示时会**原样保留**换行与空格。

```javascript
window.PROMPT_DATA = {
  categories: [
    {
      title: "编程开发",
      items: [
        {
          name: "代码审查",
          content: `请帮我审查下面这段代码：

1. 潜在 bug
2. 可读性

代码：
（粘贴）
`,
        },
      ],
    },
  ],
};
```

## 部署

### 本地使用

双击 `index.html` 即可。复制功能在 `file://` 协议下会自动降级为兼容方案。

### IIS 部署

将目录发布到 IIS 站点根目录，可按需配置 `web.config`（该文件已在 `.gitignore` 中忽略）。

## Git 说明

以下文件不会提交到仓库：

- `data.js` — 含账号密码
- `AI/prompts.js` — 个人 AI 模板
- `*.txt` — 文本文件
- `*.config` / `web.config` — 配置文件

若文件曾被提交过，需执行以下命令停止跟踪：

```bash
git rm --cached data.js
git rm --cached AI/prompts.js
```

## 安全提示

- `data.js` 以明文存储密码，仅适合个人本机或内网私有环境使用
- 切勿将含真实凭据的 `data.js` 推送到公开仓库
- 建议仅在内网访问，或配合访问控制使用
