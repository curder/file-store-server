-- 创建数据库
CREATE DATABASE IF NOT EXISTS `file_store_server`;

-- 数据库表
CREATE TABLE `files` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '文件名',
  `sha1` char(40) NOT NULL DEFAULT '' COMMENT '文件sha1值',
  `size` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `path` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储路径',
  `status` datetime DEFAULT NULL COMMENT '状态「可用｜禁用｜已删除等状态的标示」',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后时间',
  `file_extras` text COMMENT '拓展字段',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_unique_key_sha1` (`sha1`),
  KEY `index_key_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
