# Backup CLI Tool

The Backup CLI Tool is a command-line utility written in Go that allows you to create and share backups of directories easily. It provides features such as creating backups, encrypting files, sharing backups through a data transfer protocol, and more.

## Installation



Once you have Go installed, you can install the Backup CLI Tool using the following command:

```bash
go install github.com/Rushabh528/SNTRecruitmentTask


Commands:

1. backup -source [source directory] -destination [destination directory] [flags] #For backup of a directory
2. backup -share -source [source directory] -destination [destination directory] [flags] #To share a backed-up directory through a data transfer protocol


Flags used:

-source: Specifies the source directory to be backed up.
-destination: Specifies the destination directory for the backup.
-encrypt: Encrypts files while backing up.
-recursive: Recursively encrypts subdirectories.
-selective-encrypt: Selectively encrypts files based on file extensions.
-share: Shares a backed-up directory through a data transfer protocol.
-files-to-send: Specifies which files to send when sharing the backup.
-previous-logs: Includes previous backup logs in the shared backup.
-root-dir: Defines the root directory for the backup.
-logger-format: Specifies the format of the logger file in the backup directory.


Example: go run backup.go -share -source "C:\sourcePath" -destination "C:\destinationPath"  -files-to-send=file1.txt,file2.txt -previous-logs





