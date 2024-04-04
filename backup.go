package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	sourceDir        string
	destinationDir   string
	share            bool
	filesToSend      string
	previousLogs     bool
	logFile          *os.File
	previousLogPath  string
	sharedBackupName string
)

func main() {
	flag.StringVar(&sourceDir, "source", "", "Source directory to share")
	flag.StringVar(&destinationDir, "destination", "", "Destination directory on remote server")
	flag.BoolVar(&share, "share", false, "Share the backed-up directory")
	flag.StringVar(&filesToSend, "files-to-send", "", "Comma-separated list of files to send")
	flag.BoolVar(&previousLogs, "previous-logs", false, "Include previous backup logs in shared backup")
	flag.Parse()

	if sourceDir == "" || destinationDir == "" {
		fmt.Println("Usage: backup -source [source directory] -destination [destination directory] [-share] [-files-to-send=file1,file2] [-previous-logs]")
		os.Exit(1)
	}

	// Initialize logger
	initLogger(destinationDir)

	if share {
		err := shareBackup(sourceDir, destinationDir)
		if err != nil {
			fmt.Printf("Error sharing backup: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Backup shared successfully")
	} else {
		fmt.Println("Use -share flag to share the backed-up directory")
	}
}

func shareBackup(sourceDir, destinationDir string) error {
	// Construct shared backup name
	sharedBackupName = fmt.Sprintf("backup_%s", time.Now().Format("20060102150405"))

	// Create destination directory if it doesn't exist
	if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
		os.MkdirAll(destinationDir, os.ModePerm)
	}

	// Create shared backup directory
	sharedBackupDir := filepath.Join(destinationDir, sharedBackupName)
	if err := os.Mkdir(sharedBackupDir, os.ModePerm); err != nil {
		return err
	}

	// Perform backup
	err := performBackup(sourceDir, sharedBackupDir)
	if err != nil {
		return err
	}

	// If files-to-send flag is provided, only transfer specified files
	if filesToSend != "" {
		filesList := strings.Split(filesToSend, ",")
		for _, file := range filesList {
			srcFile := filepath.Join(sourceDir, file)
			destFile := filepath.Join(sharedBackupDir, file)
			if _, err := os.Stat(srcFile); err == nil {
				if err := copyFile(srcFile, destFile); err != nil {
					return err
				}
			}
		}
	}

	// If previous-logs flag is provided, include previous backup logs in shared backup
	if previousLogs {
		if previousLogPath != "" {
			if err := copyFile(previousLogPath, filepath.Join(sharedBackupDir, "previous_logs.log")); err != nil {
				return err
			}
		}
	}

	return nil
}

func performBackup(sourceDir, destinationDir string) error {
	// Walk through the source directory and copy files/directories to the destination directory
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		// Handle errors
		if err != nil {
			return err
		}

		// Construct destination path
		dest := filepath.Join(destinationDir, path[len(sourceDir):])

		// If directory, create it in the destination directory
		if info.IsDir() {
			return os.MkdirAll(dest, os.ModePerm)
		}

		// If file, copy it to the destination directory
		if err := copyFile(path, dest); err != nil {
			return err
		}

		return nil
	})
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}

func initLogger(destinationDir string) {
	// Open log file in destination directory
	logFilePath := filepath.Join(destinationDir, "backup.log")
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Set previousLogPath
	previousLogPath = logFilePath
}
