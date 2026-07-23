
CREATE TABLE IF NOT EXISTS `t_user` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `account` VARCHAR(50) NOT NULL COMMENT '登录账号',
  `password` VARCHAR(255) NOT NULL COMMENT '密码哈希',
  `salt` VARCHAR(256) NULL COMMENT '盐值',
  `name` VARCHAR(50) NOT NULL COMMENT '姓名',
  `en_name` VARCHAR(100) NULL COMMENT '英文名',
  `create_time` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_user_id` INT NULL COMMENT '创建人ID',
  `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '软删除 0否 1是',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_account` (`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
