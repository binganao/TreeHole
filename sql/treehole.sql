/*
 Navicat Premium Data Transfer

 Source Server         : dmsj
 Source Server Type    : MySQL
 Source Server Version : 50734
 Source Host           : localhost:8889
 Source Schema         : treehole

 Target Server Type    : MySQL
 Target Server Version : 50734
 File Encoding         : 65001

 Date: 23/12/2021 23:30:03
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `enc_user_name` text NOT NULL,
  `enc_pass_word` text NOT NULL,
  `enc_secret` text,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=59 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
