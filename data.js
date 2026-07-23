/**
 * 站点数据 — 直接编辑本文件即可维护，无需数据库。
 *
 * 结构说明：
 * - categories: 分类列表，每个分类显示为一个区块
 * - title: 区块标题
 * - sites: 该分类下的站点
 *   - name: 站点名称（卡片上显示）
 *   - url: 网址
 *   - username: 账号（可选，留空字符串即可）
 *   - password: 密码（可选）
 *   - note: 备注（可选）
 */
window.SITE_DATA = {
  categories: [
    {
      title: "应用环境",
      sites: [
        {
          name: "开发环境",
          url: "http://192.168.77.1:8088/login",
          username: "admin",
          password: "1q2w3E*",
          note: "开发环境登录",
        },
        {
          name: "测试环境",
          url: "http://192.168.77.3:8081/login",
          username: "admin",
          password: "1q2w3E*",
          note: "测试环境",
        },
        {
          name: "开发 Swagger",
          url: "https://192.168.77.1:44324/swagger/index.html",
          username: "",
          password: "",
          note: "",
        },
        {
          name: "测试 Swagger",
          url: "https://192.168.77.3:44324/swagger/index.html",
          username: "",
          password: "",
          note: "",
        },
        
        
      ],
    },
   
    {
      title: "服务器",
      sites: [
        {
          name: "开发服务器内网",
          url: "192.168.77.1",
          username: "administrator",
          password: "Lms@1304",
          note: "192.168.77.1",
        },
        {
          name: "测试服务器内网",
          url: "192.168.77.3",
          username: "administrator",
          password: "lms,123",
          note: "192.168.77.3",
        },
      ],
    },
    {
      title: "数据库",
      sites: [
        {
          name: "开发数据库",
          url: "192.168.77.1",
          username: "sa",
          password: "Lms@1304",
          note: "开发库 SQL Server",
        },
        {
          name: "测试数据库",
          url: "192.168.77.3",
          username: "sa",
          password: "Lms@1304",
          note: "测试库 SQL Server",
        },
      ],
    },
  ],
};