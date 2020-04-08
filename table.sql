-- 创建数据库
CREATE DATABASE IF NOT EXISTS `file_store_server`;

-- 数据库表
CREATE TABLE IF NOT EXISTS `files`
(
    `id`          int(11)       NOT NULL AUTO_INCREMENT,
    `name`        varchar(255)  NOT NULL DEFAULT '' COMMENT '文件名',
    `sha1`        char(40)      NOT NULL DEFAULT '' COMMENT '文件sha1值',
    `size`        bigint(20)             DEFAULT '0' COMMENT '文件大小',
    `path`        varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储路径',
    `status`      tinyint                DEFAULT NULL COMMENT '状态「可用｜禁用｜已删除等状态的标示」',
    `created_at`  datetime               DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  datetime               DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后时间',
    `file_extras` text COMMENT '拓展字段',
    PRIMARY KEY (`id`),
    UNIQUE KEY `index_unique_key_sha1` (`sha1`),
    KEY `index_key_status` (`status`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

CREATE TABLE IF NOT EXISTS `users`
(
    `id`              int(11)      NOT NULL AUTO_INCREMENT,
    `name`            varchar(64)  NOT NULL DEFAULT '' COMMENT '用户名',
    `password`        varchar(255) NOT NULL DEFAULT '' COMMENT '用户密码',
    `email`           varchar(64)  NOT NULL DEFAULT '' COMMENT '邮箱',
    `phone`           varchar(20)           DEFAULT '' COMMENT '手机号',
    `email_validated` tinyint(1)            DEFAULT '0' COMMENT '邮箱是否验证',
    `phone_validated` tinyint(1)            DEFAULT '0' COMMENT '手机号码是否验证',
    `profile`         text COMMENT '其他属性',
    `status`          tinyint(1)   NOT NULL DEFAULT '0' COMMENT '账户状态「启用｜禁用｜锁定｜标记」',
    `last_active`     datetime              DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间',
    `created_at`      datetime              DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
    `sign_up_at`      datetime              DEFAULT NULL COMMENT '最后登录时间',
    PRIMARY KEY (`id`),
    KEY `index_status_key` (`status`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

CREATE TABLE IF NOT EXISTS `user_tokens`
(
    `id`    int(11)     NOT NULL AUTO_INCREMENT,
    `name`  varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
    `token` char(40)    NOT NULL DEFAULT '' COMMENT '用户校验token',
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_index_name` (`name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;