-- AI 提问模板：分类表 + 模板表 + 初始数据
-- 请先手动选择数据库：USE `site.navigation`;

CREATE TABLE IF NOT EXISTS `t_prompt_category` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `title` VARCHAR(100) NOT NULL COMMENT '分类名称',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序，越小越靠前',
  `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '软删除 0否 1是',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_title` (`title`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI模板分类';

CREATE TABLE IF NOT EXISTS `t_prompt_item` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '模板ID',
  `category_id` INT NOT NULL COMMENT '所属分类ID',
  `name` VARCHAR(100) NOT NULL COMMENT '模板名称',
  `content` LONGTEXT NOT NULL COMMENT '模板正文（保留换行）',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序，越小越靠前',
  `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '软删除 0否 1是',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_category_sort` (`category_id`, `sort_order`),
  CONSTRAINT `fk_prompt_item_category`
    FOREIGN KEY (`category_id`) REFERENCES `t_prompt_category` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI提问模板';

-- 初始数据（对应原 prompts.js：常用模板 / 开发 / 代码审查）
INSERT INTO `t_prompt_category` (`title`, `sort_order`)
SELECT '常用模板', 1
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM `t_prompt_category` WHERE `title` = '常用模板' AND `is_deleted` = 0
);

SET @cat_id = (
  SELECT `id` FROM `t_prompt_category`
  WHERE `title` = '常用模板' AND `is_deleted` = 0
  ORDER BY `id` ASC
  LIMIT 1
);

INSERT INTO `t_prompt_item` (`category_id`, `name`, `content`, `sort_order`)
SELECT @cat_id, '开发',
'请根据需求实现代码，要求：

1. 先给出实现思路（简洁）
2. 再给出完整可运行代码
3. 关键逻辑加必要注释
4. 说明用法与注意事项
5. 不要省略关键代码，用完整实现代替伪代码

技术栈：【如 C# / ASP.NET / Vue3 / SQL Server】
编码规范：【如现有项目风格 / 命名约定】

需求：
【在此描述要做什么】

相关上下文（可选）：
【现有类 / 接口 / 表结构 / 代码片段】
', 1
FROM DUAL
WHERE @cat_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1 FROM `t_prompt_item`
    WHERE `category_id` = @cat_id AND `name` = '开发' AND `is_deleted` = 0
  );

INSERT INTO `t_prompt_item` (`category_id`, `name`, `content`, `sort_order`)
SELECT @cat_id, '代码审查',
'请对下面代码做 Code Review，按严重级别输出：

- Blocker（必须改）
- Major（建议改）
- Minor（可选优化）
- Nit（风格/命名）

每条包含：位置、问题、原因、修改建议（最好带示例代码）。
最后给一个总体评价：能否合并，以及合并前必须完成的事项。

代码：
```
【在此粘贴】
```

额外背景（可选）：
【业务场景 / 性能要求 / 截止日期】
', 2
FROM DUAL
WHERE @cat_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1 FROM `t_prompt_item`
    WHERE `category_id` = @cat_id AND `name` = '代码审查' AND `is_deleted` = 0
  );
