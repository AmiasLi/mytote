## mytote

A cron tool to back up MySQL using xtrabackup, and manage back upped files.
However, xtrabackup is not included in this tool and should be installed separately.

```bash
[user@node ~]$ sudo ./mytote
mytote is backup Cron for your MySQL using xtrabackup, and manage backup file.

Usage:
  mytote [command]

Available Commands:
  backup      Start the backup service
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  start       start the backup immediately
  version     Print the version number of mytote

Flags:
  -c, --config string   default ./config.yml
  -h, --help            help for mytote

Use "mytote [command] --help" for more information about a command.
```

## Start backup immediately
```bash
[user@node ~]$ sudo ./mytote start
```
## Start backup cron job
```bash
[user@node ~]$ sudo ./mytote backup --config config.yaml &
```

## Configuration
```yaml
# config.yml
# Path: config.yml

server:
  host: 127.0.0.1
  port: 3306
  socket:

  #  only support backup on localhost
  user: dba
  password: 123456GHY%

  # at present, only support backup with xtrabackup
  # backup_method: mysqldump
  backup_method: xtrabackup

  # Only support full backup at current version
  backup_type: full
  compress: 1
  compress_threads: 3

  #  If the backup_method is "incr", the day of full backup will be specified
  #  0-6, 0 is Sunday
  #  If not specified, the full backup will be performed Sunday
  backup_full_weekday: 1
  backup_hour: 13
  backup_minute: 54

  #  If the backup fails, it will be retried after the specified time
  retry_duration: 15
  retry_times: 1

  # keep days of backup
  backup_retain: 7

  #  if null will use current directory
  backup_dir: /data/backup_data_backup

  backup_log: /data/backup_data/backup.log

  #  Spare a space for the disk
  #  If the disk space is less than the reserve space, the backup will not be performed
  #  Data unit is GB
  #  Current Disk size DataBase size is going to be backup
  reserve_space: 5

#record backup log to mysql
mysql_log:
  host: 192.168.3.62
  port: 3306
  user: dba
  password: 123456@Aa
  db: mysql_backup_info
  table: backup_logs

#record backup log to mongodb
#mongodb_log:
#  host: 127.0.0.1
#  port: 27017
#  user: dba
#  password: 123456@Aa
#  db: test
#  Collection: backup_logs
```

In current version, only support backup with xtrabackup, 
and only support full backup.

In current version, only support send logs to mysql.
the table structure is as follows:
```sql
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
```
