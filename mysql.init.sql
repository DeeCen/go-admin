-- Adminer 4.7.7 MySQL dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `goadmin_menu`;
CREATE TABLE `goadmin_menu` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `parentId` int(11) unsigned NOT NULL DEFAULT '0',
  `type` tinyint(4) unsigned NOT NULL DEFAULT '0',
  `order` int(11) unsigned NOT NULL DEFAULT '0',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `uri` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `header` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `pluginName` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `uuid` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `createAt` int(10) unsigned NOT NULL DEFAULT '0',
  `updateAt` int(10) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `goadmin_menu` (`id`, `parentId`, `type`, `order`, `title`, `icon`, `uri`, `header`, `pluginName`, `uuid`, `createAt`, `updateAt`) VALUES
(1,    0,    1,    3,    'Admin',    'fa-tasks',    '',    '',    '',    '',    0,    0),
(2,    1,    1,    3,    'Users',    'fa-users',    '/info/manager',    '',    '',    '',    0,    1669193050),
(3,    1,    1,    4,    'Roles',    'fa-user',    '/info/role',    '',    '',    '',    0,    1669193053),
(4,    1,    1,    6,    '权限',    'fa-bars',    '/menu',    '',    '',    '',    0,    1669193055);

DROP TABLE IF EXISTS `goadmin_role`;
CREATE TABLE `goadmin_role` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `slug` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `createAt` int(10) unsigned NOT NULL DEFAULT '0',
  `updateAt` int(10) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `goadmin_role` (`id`, `name`, `slug`, `createAt`, `updateAt`) VALUES
(1,    'Administrator',    '',    1669183736,    1669183736);

DROP TABLE IF EXISTS `goadmin_role_menu`;
CREATE TABLE `goadmin_role_menu` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `roleId` int(11) unsigned NOT NULL DEFAULT '0',
  `menuId` int(11) unsigned NOT NULL DEFAULT '0',
  `createAt` int(11) unsigned NOT NULL DEFAULT '0',
  `updateAt` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `roleId_menuId` (`roleId`,`menuId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `goadmin_role_menu` (`id`, `roleId`, `menuId`, `createAt`, `updateAt`) VALUES
(1,    1,    1,    1669183736,    0),
(2,    1,    2,    1669183736,    0),
(3,    1,    3,    1669183736,    0),
(4,    1,    4,    1669183736,    0);

DROP TABLE IF EXISTS `goadmin_role_user`;
CREATE TABLE `goadmin_role_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `roleId` int(11) unsigned NOT NULL DEFAULT '0',
  `userId` int(11) unsigned NOT NULL DEFAULT '0',
  `createAt` int(11) unsigned NOT NULL DEFAULT '0',
  `updateAt` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `roleId_userId` (`roleId`,`userId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `goadmin_role_user` (`id`, `roleId`, `userId`, `createAt`, `updateAt`) VALUES
(1,    1,    1,    0,    0),
(7,    3,    3,    0,    0);

DROP TABLE IF EXISTS `goadmin_site`;
CREATE TABLE `goadmin_site` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `value` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `description` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `state` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `createAt` int(11) NOT NULL DEFAULT '0',
  `updateAt` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `goadmin_user`;
CREATE TABLE `goadmin_user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `rememberToken` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `createAt` int(10) unsigned NOT NULL DEFAULT '0',
  `updateAt` int(10) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `goadmin_user` (`id`, `username`, `password`, `name`, `avatar`, `rememberToken`, `createAt`, `updateAt`) VALUES
(1,    'admin',    '$2a$10$8QLBZIctzA9F.jAG1sOHh.xjOWRrPS.GJ3gap/Bq0xiL8r8XZN5b6',    'admin',    '',    'tlNcBVK9AvfYH7WEnwB1RKvocJu8FfRy4um3DJtwdHuJy0dwFsLOgAc0xUfh',    1600000000,    1600000000);

-- admin/admin
-- 2022-11-25 06:34:42
