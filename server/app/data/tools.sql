CREATE DATABASE IF NOT EXISTS tools;

use tools;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for app_admin
-- ----------------------------
DROP TABLE IF EXISTS `app_admin`;
CREATE TABLE `app_admin`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `username` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `real_name` varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '真实姓名',
  `salt` varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '加密盐',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '有效状态 1.有效 0.失效',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除 1.已删除 0.未删除',
  `birthday` int(11) NOT NULL DEFAULT 0 COMMENT '出生日期',
  `phone` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '手机号码',
  `email` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '电子邮箱',
  `sex` tinyint(4) NOT NULL DEFAULT 0 COMMENT '性别 1.男  2.女 3.不详',
  `role_id` int(11) NOT NULL DEFAULT 0 COMMENT '角色id',
  `route_ids` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '路由id,扩展路由权限',
  `login_ip_addr` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '登录ip地址',
  `login_time` int(11) NOT NULL DEFAULT 0 COMMENT '登录时间',
  `ctime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `utime` int(11) NOT NULL DEFAULT 0 COMMENT '修改时间',
  `introduction` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '个人签名',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `ix_app_admin_username`(`username`) USING BTREE COMMENT '账户'
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '后台管理员表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of app_admin
-- ----------------------------
INSERT INTO `app_admin` VALUES (1, 'admin', 'b7f81c35d147b459870c082b6ca67761', 'admin', 'Se3u', 1, 0, 0, '', '', 1, 0, '', '127.0.0.1', 1594553805, 1588988259, 0, '');
-- ----------------------------
-- Table structure for app_notify_tpl
-- ----------------------------
DROP TABLE IF EXISTS `app_notify_tpl`;
CREATE TABLE `app_notify_tpl`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `tpl_name` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '0' COMMENT '模板名称',
  `notify_type` tinyint(4) NOT NULL DEFAULT 0 COMMENT '模板通知类型',
  `tpl_data` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '0' COMMENT '模板数据',
  `create_id` int(11) NOT NULL DEFAULT 0 COMMENT '通知用户集合',
  `update_id` int(11) NOT NULL DEFAULT 0 COMMENT '修改者id',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态：0.未开始 1.已开始',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除 1.已删除 0.未删除',
  `ctime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `utime` int(11) NOT NULL DEFAULT 0 COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '消息通知模板表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of app_notify_tpl
-- ----------------------------
INSERT INTO `app_notify_tpl` VALUES (1, '定时任务模板', 1, '定时任务：{{task_name.DATA}}出错\n报错内容：{{err_msg.DATA}}\n执行时间：{{run_time.DATA}}', 1, 1, 1, 0, 1593960033, 1594101467);

-- ----------------------------
-- Table structure for app_roles
-- ----------------------------
DROP TABLE IF EXISTS `app_roles`;
CREATE TABLE `app_roles`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `role_name` varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '角色姓名',
  `parent_id` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '上级id',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '有效状态 1.有效 0.失效',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除 1.已删除 0.未删除',
  `route_ids` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '路由id集合',
  `ctime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `utime` int(11) NOT NULL DEFAULT 0 COMMENT '修改时间',
  `role` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '角色',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '角色配置表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of app_roles
-- ----------------------------
INSERT INTO `app_roles` VALUES (1, '研发组1', 0, 1, 0, '1,6,7,8,24,23,9,25,26,10,11,12,20,2,27,28,41,40,3,4,5,13,14,29,30,15,16,22,17,33,31,32,34,18,19,21,35,39,36,37,38', 1588988259, 1594294760, 'admin');
INSERT INTO `app_roles` VALUES (3, '研发小弟', 1, 1, 0, '', 1592379895, 0, 'admin_dd');
INSERT INTO `app_roles` VALUES (4, '研发弟中弟1', 3, 1, 0, '', 1592380095, 1594551646, 'admin_dd_dd');
INSERT INTO `app_roles` VALUES (5, '运营', 0, 1, 0, '', 1592380842, 1594551788, 'MD');
INSERT INTO `app_roles` VALUES (6, '研发弟中弟中弟', 4, 1, 1, '', 1592380932, 1592471939, 'admin_dd_dd_dd');
INSERT INTO `app_roles` VALUES (7, '测试', 0, 1, 1, '', 1592468254, 1592471945, 'ces');

-- ----------------------------
-- Table structure for app_routes
-- ----------------------------
DROP TABLE IF EXISTS `app_routes`;
CREATE TABLE `app_routes`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `parent_id` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '上级id',
  `route_name` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '路由名称',
  `route` varchar(50) CHARACTER SET utf16 COLLATE utf16_general_ci NOT NULL DEFAULT '' COMMENT '路由url',
  `request` varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '请求方式',
  `component` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '页面url',
  `path` varbinary(30) NOT NULL DEFAULT '' COMMENT '路径',
  `name` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '路由标识',
  `redirect` varchar(40) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '重定向地址',
  `hidden` tinyint(4) NOT NULL DEFAULT 1 COMMENT '是否隐藏',
  `icon` varchar(15) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '图标',
  `extra` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '额外参数',
  `is_route` tinyint(4) NOT NULL DEFAULT 0 COMMENT '前端路由',
  `sort` smallint(6) NOT NULL DEFAULT 50 COMMENT '排序',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '有效状态 1.有效 0.失效',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除 1.已删除 0.未删除',
  `ctime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `utime` int(11) NOT NULL DEFAULT 0 COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 42 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '路由配置表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of app_routes
-- ----------------------------
INSERT INTO `app_routes` VALUES (1, 0, '权限管理', 'auth', '', 'layout', 0x61757468, 'Auth', '/auth/admin/list', 0, 'lock', '', 1, 1, 1, 0, 0, 1594551254);
INSERT INTO `app_routes` VALUES (2, 1, '路由管理', 'auth/route', '', 'auth/index', 0x726F757465, 'Route', '/auth/route/list', 0, '', '', 1, 4, 1, 0, 0, 1594293210);
INSERT INTO `app_routes` VALUES (3, 2, '路由信息', 'auth/route/detail', 'GET', '', '', '', '', 0, '', '', 0, 50, 1, 0, 1592272269, 1594196441);
INSERT INTO `app_routes` VALUES (4, 2, '路由列表', 'auth/route/list', 'GET', 'auth/route/list', 0x6C697374, 'RouteList', '', 0, '', '', 1, 50, 1, 0, 1592272269, 0);
INSERT INTO `app_routes` VALUES (5, 2, '路由列表（全部）', 'auth/route/all', 'GET', '', '', '', '', 0, '', '', 0, 50, 1, 0, 1592272269, 1592900878);
INSERT INTO `app_routes` VALUES (6, 1, '管理员', 'auth/admin', '', 'auth/index', 0x61646D696E, 'Admin', '/auth/admin/list', 0, '', '', 1, 2, 1, 0, 1592272289, 1594113117);
INSERT INTO `app_routes` VALUES (7, 6, '管理员列表', 'auth/admin/list', 'GET', 'auth/admin/list', 0x6C697374, 'AdminList', '', 0, '', '', 1, 50, 1, 0, 1592273445, 1594113207);
INSERT INTO `app_routes` VALUES (8, 6, '管理员信息', 'auth/admin/detail', 'GET', '', '', '', '', 1, '', '', 0, 50, 1, 0, 1592274848, 1594196459);
INSERT INTO `app_routes` VALUES (9, 1, '角色管理', '1', '1', 'auth/index', 0x726F6C65, 'Role', '/auth/role/list', 0, '', '', 1, 3, 1, 0, 1592277285, 1594181897);
INSERT INTO `app_routes` VALUES (10, 9, '角色列表', 'auth/role/all', 'GET', '', '', '', '', 0, '', '', 0, 50, 1, 0, 1592277384, 1592900917);
INSERT INTO `app_routes` VALUES (11, 9, '角色信息', 'auth/role/detail', '', '', '', '', '', 0, '', '', 0, 50, 1, 0, 1592277440, 0);
INSERT INTO `app_routes` VALUES (12, 9, '角色列表', 'auth/role/list', '', 'auth/role/list', 0x6C697374, 'RoleList', '', 0, '', '', 1, 50, 1, 0, 0, 1594182011);
INSERT INTO `app_routes` VALUES (13, 0, '运维', 'devops', '', 'layout', 0x6465766F7073, 'devops', '/devops/servers/list', 0, 'servers', '', 1, 5, 1, 0, 1592385492, 1594194689);
INSERT INTO `app_routes` VALUES (14, 13, '服务器', 'devops/servers', '', 'devops/index', 0x73657276657273, 'servers', '/devops/servers/list', 0, '', '', 1, 6, 1, 0, 1592385492, 1592385500);
INSERT INTO `app_routes` VALUES (15, 14, '服务器列表', 'devops/servers/list', '', 'devops/servers/list', 0x6C697374, 'serversList', '', 0, '', '', 1, 50, 1, 0, 1592389641, 1594194740);
INSERT INTO `app_routes` VALUES (16, 14, '服务器信息', 'devops/servers/detail', '', '', '', '', '', 1, '', '', 0, 50, 1, 0, 1592389672, 1592389734);
INSERT INTO `app_routes` VALUES (17, 13, '定时任务', 'devops/task', '', 'devops/index', 0x7461736B, 'task', '/devops/task/list', 0, '', '', 1, 7, 1, 0, 1592389723, 1594195038);
INSERT INTO `app_routes` VALUES (18, 17, '任务列表', 'devops/task/list', '', 'devops/task/list', 0x6C697374, 'taskList', '', 0, '', '', 1, 50, 1, 0, 1592389767, 0);
INSERT INTO `app_routes` VALUES (19, 17, '任务信息', 'devops/task/detail', '', '', '', '', '', 0, '', '', 0, 50, 1, 0, 1592389783, 0);
INSERT INTO `app_routes` VALUES (20, 9, '角色权限设置', 'auth/role/route_edit', '', 'auth/role/route_edit', 0x726F7574655F65646974, '', '', 1, '', ':id', 1, 50, 1, 0, 1592901326, 0);
INSERT INTO `app_routes` VALUES (21, 17, '任务分组', 'devops/task/group_list', '', '', '', '', '', 0, '', '', 0, 50, 1, 0, 1592901513, 0);
INSERT INTO `app_routes` VALUES (22, 14, '服务器分组', 'devops/servers/group_list', '', '', '', '', '', 0, '', '', 0, 50, 1, 0, 1592901539, 0);
INSERT INTO `app_routes` VALUES (23, 6, '管理员新增', 'auth/admin/add', 'PUT', 'auth/admin/detail', 0x616464, 'AdminAdd', '', 1, '', '', 1, 50, 1, 0, 1594177906, 0);
INSERT INTO `app_routes` VALUES (24, 6, '管理员修改', 'auth/admin/edit', 'POST', 'auth/admin/detail', 0x65646974, 'AdminEdit', '', 1, '', ':id', 1, 50, 1, 0, 1594177998, 1594180413);
INSERT INTO `app_routes` VALUES (25, 9, '角色新增', 'auth/role/add', 'PUT', 'auth/role/detail', 0x616464, 'RoleAdd', '', 1, '', '', 1, 50, 1, 0, 1594182089, 0);
INSERT INTO `app_routes` VALUES (26, 9, '角色修改', 'auth/role/edit', 'POST', 'auth/role/detail', 0x65646974, 'RoleEdit', '', 1, '', ':id', 1, 50, 1, 0, 1594182136, 1594182158);
INSERT INTO `app_routes` VALUES (27, 2, '路由新增', 'auth/route/add', '', 'auth/route/detail', 0x616464, 'RouteAdd', '', 1, '', '', 1, 50, 1, 0, 1594182357, 0);
INSERT INTO `app_routes` VALUES (28, 2, '路由修改', 'auth/route/edit', '', 'auth/route/detail', 0x65646974, 'RouteEdit', '', 1, '', ':id', 1, 50, 1, 0, 1594194593, 0);
INSERT INTO `app_routes` VALUES (29, 14, '服务器新增', 'devops/servers/add', '', 'devops/servers/detail', 0x616464, 'serversAdd', '', 1, '', '', 1, 50, 1, 0, 1594194901, 0);
INSERT INTO `app_routes` VALUES (30, 14, '服务器修改', 'devops/servers/edit', '', 'devops/servers/detail', 0x65646974, '', '', 1, '', ':id', 1, 50, 1, 0, 1594194951, 0);
INSERT INTO `app_routes` VALUES (31, 17, '定时任新增', 'devops/task/add', '', 'devops/task/detail', 0x6465766F70732F7461736B2F616464, 'taskAdd', '', 1, '', '', 1, 50, 1, 0, 1594195138, 1594195191);
INSERT INTO `app_routes` VALUES (32, 17, '定时任务修改', 'devops/task/edit', '', 'devops/task/detail', 0x65646974, 'taskEdit', '', 1, '', ':id', 1, 50, 1, 0, 1594195236, 0);
INSERT INTO `app_routes` VALUES (33, 17, '任务执行日志', 'devops/task/log', '', 'devops/task/log', 0x6C6F67, 'taskLog', '', 1, '', ':id', 1, 50, 1, 0, 1594195289, 0);
INSERT INTO `app_routes` VALUES (34, 17, '任务日志信息', 'devops/task/log_detail', '', 'devops/task/log_detail', 0x6C6F675F64657461696C, 'taskLogDetail', '', 0, '', ':id', 0, 50, 1, 0, 1594195434, 0);
INSERT INTO `app_routes` VALUES (35, 13, '消息模板', 'devops/notify', '', 'devops/index', 0x6E6F74696679, 'notify', '/devops/notify/list', 0, '', '', 1, 8, 1, 0, 1594195503, 0);
INSERT INTO `app_routes` VALUES (36, 35, '消息模板', 'devops/notify/list', '', 'devops/notify/list', 0x6C697374, 'notifyList', '', 0, '', '', 1, 50, 1, 0, 1594195561, 0);
INSERT INTO `app_routes` VALUES (37, 35, '消息模板新增', 'devops/notify/add', '', 'devops/notify/detail', 0x616464, 'notifyAdd', '', 1, '', '', 1, 50, 1, 0, 1594195619, 0);
INSERT INTO `app_routes` VALUES (38, 35, '模板消息修改', 'devops/notify/edit', '', 'devops/notify/detail', 0x65646974, 'notifyEdit', '', 1, '', ':id', 1, 50, 1, 0, 1594195656, 0);
INSERT INTO `app_routes` VALUES (39, 35, '消息模板信息', 'devops/notify/detail', '', '', '', '', '', 0, '', '', 0, 50, 1, 0, 1594195722, 1594195762);
INSERT INTO `app_routes` VALUES (40, 2, '路由排序设置', 'auth/route/edit_sort', '', '', '', '', '', 1, '', '', 0, 50, 1, 0, 1594292858, 1594293032);
INSERT INTO `app_routes` VALUES (41, 2, '前端菜单', 'auth/route/menu', '', '', '', '', '', 0, '', '', 0, 50, 1, 0, 1594294746, 0);

-- ----------------------------
-- Table structure for app_task
-- ----------------------------
DROP TABLE IF EXISTS `app_task`;
CREATE TABLE `app_task`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `group_id` int(11) NOT NULL DEFAULT 0 COMMENT '分组id',
  `server_ids` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '调用服务器id集合',
  `run_type` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '执行策略：0-同时执行，1-轮询执行',
  `task_name` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '任务类型',
  `description` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '任务描述',
  `cron_spec` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '任务表达式',
  `concurrent` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否并发执行  1.是',
  `command` varchar(300) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '0' COMMENT '执行命令',
  `timeout` int(11) NOT NULL DEFAULT 0 COMMENT '超时时间',
  `job_entry_id` int(11) NOT NULL DEFAULT 0 COMMENT 'job id',
  `execute_times` int(11) NOT NULL DEFAULT 0 COMMENT '累计执行次数',
  `prev_time` int(11) NOT NULL DEFAULT 0 COMMENT '上次执行时间',
  `is_notify` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否通知 1.通知',
  `notify_type` tinyint(4) NOT NULL DEFAULT 0 COMMENT '通知类型',
  `notify_tpl_id` int(11) NOT NULL DEFAULT 0 COMMENT '通知模板id',
  `notify_user_ids` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '0' COMMENT '通知用户集合',
  `create_id` int(11) NOT NULL DEFAULT 0 COMMENT '创建者id',
  `update_id` int(11) NOT NULL DEFAULT 0 COMMENT '修改者id',
  `is_audit` tinyint(4) NOT NULL DEFAULT 0 COMMENT '审核：0.待审核 1.已通过 2.已拒绝',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态：0.未开始 1.已开始',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除 1.已删除 0.未删除',
  `ctime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `utime` int(11) NOT NULL DEFAULT 0 COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '任务配置表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of app_task
-- ----------------------------
INSERT INTO `app_task` VALUES (1, 1, '1', '1', '测试任务', '测试定时任务', '0 */1 * * * *', 0, 'date \"+%Y-%m-%d %H:%M:%S\" >> /usr/local/www/cron.log', 60, 0, 3, 1593627180, 0, 0, 0, '', 1, 0, 1, 0, 0, 1593399461, 0);
INSERT INTO `app_task` VALUES (2, 2, '0', '1', '测试删除', '333', '0/1 * * * * *', 0, 'echo 111', 60, 0, 0, 0, 0, 0, 0, '', 1, 0, 2, 0, 1, 1593502301, 0);
INSERT INTO `app_task` VALUES (3, 1, '1', '1', 'ccc', 'eee', '0 */1 * * * *', 0, 'll', 60, 0, 2, 1593627180, 1, 1, 1, '3,2,1', 1, 0, 1, 0, 0, 1593502453, 0);

-- ----------------------------
-- Table structure for app_task_servers
-- ----------------------------
DROP TABLE IF EXISTS `app_task_servers`;
CREATE TABLE `app_task_servers`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `group_id` int(11) NOT NULL DEFAULT 0 COMMENT '分组id',
  `server_name` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '服务器名称',
  `connection_type` tinyint(4) NOT NULL DEFAULT 0 COMMENT '连接类型 连接类型 0:SSH;1:Telnet;',
  `type` tinyint(4) NOT NULL DEFAULT 0 COMMENT '登录类型：0-密码登录，1-私钥登录',
  `server_account` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '登录账号',
  `password` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '0' COMMENT '登录密码',
  `server_outer_ip` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '外网IP',
  `server_ip` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '0' COMMENT '内网ip',
  `port` smallint(6) NOT NULL DEFAULT 0 COMMENT '端口',
  `private_key_src` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '私钥路径',
  `public_key_src` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '公钥路径',
  `detail` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态：1-正常',
  `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否删除 1.已删除 0.未删除',
  `ctime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
  `utime` int(11) NOT NULL DEFAULT 0 COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '服务器列表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of app_task_servers
-- ----------------------------
INSERT INTO `app_task_servers` VALUES (1, 1, '测试服务器', 0, 0, 'root', 'root', '', '127.0.0.1', 22, '', '', '测试专用', 1, 0, 1592901915, 0);

SET FOREIGN_KEY_CHECKS = 1;
