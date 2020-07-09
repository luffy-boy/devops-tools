CREATE DATABASE IF NOT EXISTS tools;

use tools;

DEFAULT CHARACTER SET utf8;

CREATE TABLE `app_common_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `conf_type` int(11) NOT NULL DEFAULT '0' COMMENT '配置类型',
  `data` text NOT NULL COMMENT '配置，json结构',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '有效状态 1.有效 0.失效',
  `ctime` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `utime` int(11) NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='通用配置表(单条配置)';

CREATE TABLE `app_admin`(
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `username` varchar(50) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(32) NOT NULL DEFAULT '' COMMENT '密码',
  `real_name` varchar(10) NOT NULL DEFAULT '' COMMENT '真实姓名',
  `salt` varchar(10) NOT NULL DEFAULT '' COMMENT '加密盐',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '有效状态 1.有效 0.失效',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除 1.已删除 0.未删除',
  `birthday` int NOT NULL DEFAULT '0' COMMENT '出生日期',
  `phone` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号码',
  `email` varchar(30) NOT NULL DEFAULT '' COMMENT '电子邮箱',
  `sex` tinyint NOT NULL DEFAULT '0' COMMENT '性别 1.男  2.女 3.不详',
  `role_id` int NOT NULL  DEFAULT '0' COMMENT '角色id',
  `route_ids` varchar(200) NOT NULL DEFAULT '' COMMENT '路由id,扩展路由权限',
  `login_ip_addr`  varchar(20) NOT NULL DEFAULT '' COMMENT '登录ip地址',
  `login_time`     int NOT NULL DEFAULT '0' COMMENT '登录时间',
  `introduction`   varchar(200)  NOT NULL DEFAULT '' COMMENT '个性签名'
  `ctime` int NOT NULL DEFAULT '0' COMMENT '创建时间',
  `utime` int NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='后台管理员表';
-- 创建索引
ALTER TABLE `tools`.`app_admin` ADD UNIQUE INDEX `ix_app_admin_username`(`username`) COMMENT '账户';
-- 插入数据
INSERT INTO  `tools`.`app_admin` (username, password, real_name, salt, status, is_delete, birthday, phone, sex, role_id, route_ids, ctime, utime) VALUES  ( "luffy", "b7f81c35d147b459870c082b6ca67761", "luffy", "Se3u", 1, 0, 0, "", 1, 0, "", 1588988259, 0);

CREATE TABLE `app_roles`(
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `role` varchar(10) NOT NULL DEFAULT '' COMMENT '角色',
  `role_name` varchar(10) NOT NULL DEFAULT '' COMMENT '角色姓名',
  `parent_id` int(11) UNSIGNED NOT NULL DEFAULT '0' COMMENT '上级id',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '有效状态 1.有效 0.失效',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除 1.已删除 0.未删除',
  `route_ids` text NOT NULL COMMENT '路由id集合',
  `ctime` int NOT NULL DEFAULT '0' COMMENT '创建时间',
  `utime` int NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='角色配置表';

-- 插入数据
INSERT INTO `tools`.`app_roles` (role_name, parent_id, status, is_delete, route_ids, ctime, utime)  VALUES  ("研发组", 0, 1, 0, "",1588988259,0);

CREATE TABLE `app_routes`(
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `parent_id` int(11) UNSIGNED NOT NULL DEFAULT '0' COMMENT '上级id',
  `route_name` varchar(20) NOT NULL DEFAULT '' COMMENT '路由名称',
  `route` varchar(50) NOT NULL DEFAULT '' COMMENT '路由url',
  `request` varchar(10) NOT NULL DEFAULT '' COMMENT '请求方式',
  `path` varchar(30) NOT NULL DEFAULT '' COMMENT '路径',
  `component` varchar(30) NOT NULL DEFAULT '' COMMENT '页面url',
  `name` varchar(20) NOT NULL DEFAULT '' COMMENT '路由标识',
  `redirect` varchar(40) NOT NULL DEFAULT '' COMMENT '重定向地址',
  `hidden` tinyint NOT NULL DEFAULT '0' COMMENT '是否隐藏',
  `icon` varchar(15) NOT NULL DEFAULT '' COMMENT '图标',
  `extra` varchar(50) NOT NULL DEFAULT '' COMMENT '额外参数',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '有效状态 1.有效 0.失效',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除 1.已删除 0.未删除',
  `ctime` int NOT NULL DEFAULT '0' COMMENT '创建时间',
  `utime` int NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='路由配置表';

CREATE TABLE `app_task_servers`(
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `group_id` int NOT NULL DEFAULT '0' COMMENT '分组id',
  `server_name` varchar(30) NOT NULL DEFAULT '' COMMENT '服务器名称',
  `connection_type` tinyint NOT NULL DEFAULT '0' COMMENT '连接类型 连接类型 0:SSH;1:Telnet;',
  `type` tinyint NOT NULL DEFAULT '0' COMMENT '登录类型：0-密码登录，1-私钥登录',
  `server_account` varchar(50) NOT NULL DEFAULT '' COMMENT '登录账号',
  `password` varchar(50) NOT NULL DEFAULT '0' COMMENT '登录密码',
  `server_outer_ip` varchar(20) NOT NULL DEFAULT '' COMMENT '外网IP',
  `server_ip` varchar(20) NOT NULL DEFAULT '0' COMMENT '内网ip',
  `port` smallint NOT NULL DEFAULT '0' COMMENT '端口',
  `private_key_src` varchar(200) NOT NULL DEFAULT '' COMMENT '私钥路径',
  `public_key_src` varchar(200) NOT NULL DEFAULT '' COMMENT '公钥路径',
  `detail` varchar(50) NOT NULL DEFAULT '' COMMENT '备注',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：1-正常',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除 1.已删除 0.未删除',
  `ctime` int NOT NULL DEFAULT '0' COMMENT '创建时间',
  `utime` int NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='服务器列表';

CREATE TABLE `app_task`(
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `group_id` int NOT NULL DEFAULT '0' COMMENT '分组id',
  `server_ids` varchar(100) NOT NULL DEFAULT '' COMMENT '调用服务器id集合',
  `run_type` varchar(100) NOT NULL DEFAULT '' COMMENT '执行策略：0-同时执行，1-轮询执行',
  `task_name` varchar(30) NOT NULL DEFAULT '' COMMENT '任务类型',
  `description` varchar(50) NOT NULL DEFAULT '' COMMENT '任务描述',
  `cron_spec` varchar(30) NOT NULL DEFAULT '' COMMENT '任务表达式',
  `concurrent` tinyint NOT NULL DEFAULT '0' COMMENT '是否并发执行  1.是',
  `command` varchar(300) NOT NULL DEFAULT '0' COMMENT '执行命令',
  `timeout` int NOT NULL DEFAULT '0' COMMENT '超时时间',
  `job_entry_id` int NOT NULL DEFAULT '0' COMMENT 'job id',
  `execute_times` int NOT NULL DEFAULT '0' COMMENT '累计执行次数',
  `prev_time` int NOT NULL DEFAULT '0' COMMENT '上次执行时间',
  `is_notify` tinyint NOT NULL DEFAULT '0' COMMENT '是否通知 1.通知',
  `notify_type` tinyint NOT NULL DEFAULT '0' COMMENT '通知类型',
  `notify_tpl_id` int NOT NULL DEFAULT '0' COMMENT '通知末班id',
  `notify_user_ids` varchar(500) NOT NULL DEFAULT '0' COMMENT '通知用户集合',
  `create_id` int NOT NULL DEFAULT '0' COMMENT '通知用户集合',
  `update_id` int NOT NULL DEFAULT '0' COMMENT '修改者id',
  `is_audit` tinyint NOT NULL DEFAULT '0' COMMENT '审核：0.待审核 1.已通过 2.已拒绝',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0.未开始 1.已开始',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除 1.已删除 0.未删除',
  `ctime` int NOT NULL DEFAULT '0' COMMENT '创建时间',
  `utime` int NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='任务配置表';

CREATE TABLE `app_notify_tpl`(
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `tpl_name` varchar(20) NOT NULL DEFAULT '0' COMMENT '模板名称',
  `notify_type` tinyint NOT NULL DEFAULT '0' COMMENT '模板通知类型',
  `tpl_data` varchar(500) NOT NULL DEFAULT '0' COMMENT '模板数据',
  `create_id` int NOT NULL DEFAULT '0' COMMENT '通知用户集合',
  `update_id` int NOT NULL DEFAULT '0' COMMENT '修改者id',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0.未开始 1.已开始',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除 1.已删除 0.未删除',
  `ctime` int NOT NULL DEFAULT '0' COMMENT '创建时间',
  `utime` int NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='消息通知模板表';