-- ============================================================
-- 任务调度系统数据库初始化脚本
-- 生成日期: 2024-01-01
-- 说明: 此脚本创建任务调度系统所需的全部表结构
-- ============================================================

-- 如果数据库不存在则创建
CREATE DATABASE IF NOT EXISTS task_scheduler
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE task_scheduler;

-- ============================================================
-- 用户表
-- 存储系统登录用户信息
-- ============================================================
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
    `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID，主键',
    `username`     VARCHAR(64)     NOT NULL                COMMENT '用户名，唯一',
    `password`     VARCHAR(255)    NOT NULL                COMMENT '密码（加密存储）',
    `display_name` VARCHAR(128)    DEFAULT NULL            COMMENT '显示名称',
    `token`        VARCHAR(255)    DEFAULT NULL            COMMENT '登录令牌',
    `status`       TINYINT         NOT NULL DEFAULT 1      COMMENT '状态: 1-启用, 0-禁用',
    `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    KEY `idx_token` (`token`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT ='系统用户表';

-- 插入默认管理员用户 (密码: admin123)
INSERT INTO `users` (`username`, `password`, `display_name`, `token`)
VALUES ('admin', 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855', '系统管理员', 'admin-token-000000000000');

-- ============================================================
-- 任务定义表
-- 存储调度任务的定义配置
-- ============================================================
DROP TABLE IF EXISTS `tasks`;
CREATE TABLE `tasks` (
    `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '任务ID，主键',
    `name`          VARCHAR(128)    NOT NULL                COMMENT '任务名称',
    `description`   TEXT            DEFAULT NULL            COMMENT '任务描述',
    `cron_expression` VARCHAR(64)   NOT NULL                COMMENT 'Cron表达式，如: 0/5 * * * * ?',
    `command`       TEXT            NOT NULL                COMMENT '要执行的命令或脚本',
    `timeout`       INT             NOT NULL DEFAULT 300    COMMENT '超时时间（秒）',
    `retry_count`   INT             NOT NULL DEFAULT 0      COMMENT '失败重试次数',
    `retry_interval` INT            NOT NULL DEFAULT 60     COMMENT '重试间隔（秒）',
    `status`        TINYINT         NOT NULL DEFAULT 1      COMMENT '状态: 1-启用, 0-禁用',
    `created_by`    BIGINT UNSIGNED DEFAULT NULL            COMMENT '创建人ID',
    `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_name` (`name`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT ='任务定义表';

-- ============================================================
-- 任务执行日志表
-- 存储每次任务执行的详细记录
-- ============================================================
DROP TABLE IF EXISTS `execution_logs`;
CREATE TABLE `execution_logs` (
    `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '执行日志ID，主键',
    `task_id`      BIGINT UNSIGNED NOT NULL                COMMENT '关联的任务ID',
    `task_name`    VARCHAR(128)    NOT NULL                COMMENT '任务名称（冗余存储）',
    `trigger_type` VARCHAR(32)     NOT NULL                COMMENT '触发类型: manual-手动, cron-定时',
    `status`       VARCHAR(32)     NOT NULL DEFAULT 'pending' COMMENT '执行状态: pending-等待中, running-执行中, success-成功, failed-失败, timeout-超时, retrying-重试中',
    `start_time`   DATETIME(3)     DEFAULT NULL            COMMENT '开始执行时间',
    `end_time`     DATETIME(3)     DEFAULT NULL            COMMENT '执行结束时间',
    `duration`     BIGINT          DEFAULT NULL            COMMENT '执行耗时（毫秒）',
    `retry_attempts` INT           NOT NULL DEFAULT 0      COMMENT '实际重试次数',
    `output`       TEXT            DEFAULT NULL            COMMENT '执行输出',
    `error_message` TEXT           DEFAULT NULL            COMMENT '错误信息',
    `params`       TEXT            DEFAULT NULL            COMMENT '业务参数(JSON)',
    `node_id`      VARCHAR(64)     DEFAULT NULL            COMMENT '执行节点ID',
    `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_task_id` (`task_id`),
    KEY `idx_status` (`status`),
    KEY `idx_start_time` (`start_time`),
    KEY `idx_task_status` (`task_id`, `status`),
    CONSTRAINT `fk_execution_task` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT ='任务执行日志表';

-- ============================================================
-- 节点状态表
-- 存储调度节点的状态信息
-- ============================================================
DROP TABLE IF EXISTS `nodes`;
CREATE TABLE `nodes` (
    `id`            VARCHAR(64)  NOT NULL                COMMENT '节点ID（UUID）',
    `name`          VARCHAR(128) NOT NULL                COMMENT '节点名称',
    `ip_address`    VARCHAR(64)  NOT NULL                COMMENT '节点IP地址',
    `status`        VARCHAR(32)  NOT NULL DEFAULT 'offline' COMMENT '节点状态: online-在线, offline-离线, busy-繁忙',
    `last_heartbeat` DATETIME   DEFAULT NULL             COMMENT '最后心跳时间',
    `cpu_usage`     DECIMAL(5,2) DEFAULT NULL             COMMENT 'CPU使用率(%)',
    `memory_usage`  DECIMAL(5,2) DEFAULT NULL             COMMENT '内存使用率(%)',
    `created_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`),
    KEY `idx_last_heartbeat` (`last_heartbeat`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT ='调度节点状态表';
