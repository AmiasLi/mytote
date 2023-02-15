## DBM

A cron job to backup MySQL databases using Xtrabackup.

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