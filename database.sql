# ************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 8.0.27)
# Database: demo
# Generation Time: 2021-12-08 09:25:06 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table list
# ------------------------------------------------------------

DROP TABLE IF EXISTS `list`;

CREATE TABLE `list` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` varchar(255) NOT NULL,
  `title` varchar(100) NOT NULL COMMENT '标题',
  `type` varchar(100) NOT NULL DEFAULT '',
  `hide` tinyint(1) NOT NULL DEFAULT '0',
  `color` varchar(100) NOT NULL DEFAULT '',
  `create_id` int NOT NULL DEFAULT '0' COMMENT '创建人',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改日期',
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

LOCK TABLES `list` WRITE;
/*!40000 ALTER TABLE `list` DISABLE KEYS */;

INSERT INTO `list` (`id`, `uid`, `title`, `type`, `hide`, `color`, `create_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
	(1,'1e3bc83e58ad3d1bffc52ea2f8dca499','这是一个测试任务','',0,'1',0,'2021-12-08 14:28:13','2021-12-08 14:28:13',NULL),
	(2,'ff6218b87c65fc8515a3ec95e38d6f78','测试任务','任务清单',1,'212321',0,'2021-12-08 14:29:13','2021-12-08 15:16:22',NULL);

/*!40000 ALTER TABLE `list` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table todo
# ------------------------------------------------------------

DROP TABLE IF EXISTS `todo`;

CREATE TABLE `todo` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `title` varchar(100) NOT NULL COMMENT '标题',
  `type` varchar(10) NOT NULL DEFAULT '' COMMENT '类型',
  `content` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '内容',
  `list_id` int NOT NULL DEFAULT '0' COMMENT '分类ID',
  `parent_id` int NOT NULL DEFAULT '0' COMMENT '父级任务ID',
  `user_id` int NOT NULL DEFAULT '0' COMMENT '用户ID',
  `priority` tinyint(1) NOT NULL DEFAULT '0' COMMENT '优先级',
  `status` varchar(10) NOT NULL DEFAULT '' COMMENT '状态',
  `top` int NOT NULL DEFAULT '0' COMMENT '置顶',
  `deadline` timestamp NULL DEFAULT NULL COMMENT '截止日期',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改日期',
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;



# Dump of table user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(100) NOT NULL COMMENT '用户名',
  `email` varchar(100) NOT NULL DEFAULT '' COMMENT '邮箱',
  `phone` varchar(256) DEFAULT '' COMMENT '电话',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改日期',
  `deleted_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `remember_token` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;

INSERT INTO `user` (`id`, `name`, `email`, `phone`, `created_at`, `updated_at`, `deleted_at`, `remember_token`)
VALUES
	(3,'Wenjunzheng','wenjunzheng@sohu-inc.com','','2021-10-25 14:46:04','2021-11-24 11:08:04',NULL,''),
	(4,'Ansme','ansme@ansme.me','118','2021-11-18 10:20:15','2021-11-18 10:20:15',NULL,''),
	(5,'abc','ansme@ansme.me','118','2021-11-18 10:21:19','2021-11-18 11:41:54','2021-11-18 11:41:54',''),
	(7,'ansme','ansme','118','2021-11-18 11:41:54','2021-11-18 11:41:54',NULL,'');

/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table user_auth
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_auth`;

CREATE TABLE `user_auth` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `email` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `auth` varchar(100) NOT NULL DEFAULT '' COMMENT '密码',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改日期',
  `deleted_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `remember_token` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

LOCK TABLES `user_auth` WRITE;
/*!40000 ALTER TABLE `user_auth` DISABLE KEYS */;

INSERT INTO `user_auth` (`id`, `email`, `auth`, `created_at`, `updated_at`, `deleted_at`, `remember_token`)
VALUES
	(1,'wenjunzheng@sohu-inc.com','123','2021-10-25 14:46:04','2021-10-25 14:46:04',NULL,'');

/*!40000 ALTER TABLE `user_auth` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;