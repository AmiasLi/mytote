# send logs to mysql server
create database mysql_backup_info;
use mysql_backup_info;
CREATE TABLE `backup_logs`
(
    `id`                 int          NOT NULL AUTO_INCREMENT,
    `host_name`          varchar(255) NOT NULL,
    `ip`                 varchar(128) NOT NULL,
    `port`               int          NOT NULL DEFAULT '3306',
    `backup_status`      varchar(10)           DEFAULT NULL,
    `start_time`         datetime              DEFAULT NULL,
    `end_time`           datetime              DEFAULT NULL,
    `duration`           varchar(256)          DEFAULT NULL,
    `backup_date`        datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `backup_type`        varchar(10)  NOT NULL,
    `backup_file`        varchar(255)          DEFAULT NULL,
    `backup_size`        int                   DEFAULT '0',
    `backup_file_status` smallint              DEFAULT '1' COMMENT '1: exist,0: dropped',
    `file_drop_time`     datetime              DEFAULT NULL,
    `err_message`        varchar(532)          DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_backup_date` (`backup_date`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;