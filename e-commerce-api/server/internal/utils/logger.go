package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

// Color codes for console output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[37m"
	ColorWhite  = "\033[97m"
)

// Logger represents the application logger
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	warnLogger  *log.Logger
	logFile     *os.File
	errorFile   *os.File
	useColor    bool
}

var globalLogger *Logger

// pintu gerbang app
func InitLogger(logDir string, useColor bool) error {
	// Create log directory if not exists
	if err := validateLogPath(logDir); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	if err := os.MkdirAll(logDir, 0750); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create archive directory
	archiveDir := filepath.Join(logDir, "archive")
	if err := os.MkdirAll(archiveDir, 0750); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	// Open log files
	// #nosec G304 -- logDir comes from internal configuration, not user input
	logFile, err := os.OpenFile(
		filepath.Join(logDir, "app.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0600,
	)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
    // #nosec G304 -- logDir comes from internal configuration, not user input
	errorFile, err := os.OpenFile(
		filepath.Join(logDir, "error.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0600,
	)
	if err != nil {
		if closeErr := logFile.Close(); closeErr != nil {
			log.Printf("warning: failed to close log file:%v",closeErr)
		}
		return fmt.Errorf("failed to open error log file: %w", err)
	}

	// Create multi-writers (console + file)
	infoWriter := io.MultiWriter(os.Stdout, logFile)
	errorWriter := io.MultiWriter(os.Stderr, errorFile, logFile)

	// Initialize loggers
	globalLogger = &Logger{
		infoLogger:  log.New(infoWriter, "", 0),
		errorLogger: log.New(errorWriter, "", 0),
		debugLogger: log.New(infoWriter, "", 0),
		warnLogger:  log.New(infoWriter, "", 0),
		logFile:     logFile,
		errorFile:   errorFile,
		useColor:    useColor,
	}

	return nil
}
// NEW path validation to prevent directory
func validateLogPath(logDir string) error {
	// 1. Clean the path (remove .., //, etc)
	cleanPath := filepath.Clean(logDir)
	
	// 2. Reject absolute paths (security: should be relative)
	if filepath.IsAbs(cleanPath) {
		return fmt.Errorf("log directory must be a relative path, got: %s", cleanPath)
	}
	
	// 3. Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected in log directory: %s", cleanPath)
	}
	
	// 4. Ensure path doesn't escape working directory
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}
	
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	
	// Path must be inside current working directory
	if !strings.HasPrefix(absPath, cwd) {
		return fmt.Errorf("log directory outside working directory: %s", cleanPath)
	}
	
	return nil
}

// CloseLogger closes the log files
func CloseLogger() {
	if globalLogger != nil {
		//  Handle errors from Close()
		if globalLogger.logFile != nil {
			if err := globalLogger.logFile.Close(); err != nil {
				log.Printf("Warning: failed to close log file: %v", err)
			}
		}
		if globalLogger.errorFile != nil {
			if err := globalLogger.errorFile.Close(); err != nil {
				log.Printf("Warning: failed to close error file: %v", err)
			}
		}
	}
}

// formatMessage formats the log message with timestamp, level, and caller info
func (l *Logger) formatMessage(level LogLevel, message string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	// Get caller info (file:line)
	_, file, line, ok := runtime.Caller(3)
	caller := "unknown"
	if ok {
		caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	// Format without color (for file)
	plainMessage := fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, level, caller, message)

	// Add color for console if enabled
	if l.useColor {
		var color string
		switch level {
		case DEBUG:
			color = ColorGray
		case INFO:
			color = ColorGreen
		case WARN:
			color = ColorYellow
		case ERROR:
			color = ColorRed
		case FATAL:
			color = ColorPurple
		default:
			color = ColorWhite
		}
		return fmt.Sprintf("%s[%s] [%s%s%s] [%s] %s%s", 
			ColorCyan, timestamp, color, level, ColorReset, caller, message, ColorReset)
	}

	return plainMessage
}

// Debug logs a debug message
func Debug(message string, args ...interface{}) {
	if globalLogger == nil {
		return
	}
	msg := fmt.Sprintf(message, args...)
	globalLogger.debugLogger.Println(globalLogger.formatMessage(DEBUG, msg))
}

// Info logs an info message
func Info(message string, args ...interface{}) {
	if globalLogger == nil {
		return
	}
	msg := fmt.Sprintf(message, args...)
	globalLogger.infoLogger.Println(globalLogger.formatMessage(INFO, msg))
}

// Warn logs a warning message
func Warn(message string, args ...interface{}) {
	if globalLogger == nil {
		return
	}
	msg := fmt.Sprintf(message, args...)
	globalLogger.warnLogger.Println(globalLogger.formatMessage(WARN, msg))
}

// Error logs an error message
func Error(message string, args ...interface{}) {
	if globalLogger == nil {
		return
	}
	msg := fmt.Sprintf(message, args...)
	globalLogger.errorLogger.Println(globalLogger.formatMessage(ERROR, msg))
}

// Fatal logs a fatal message and exits
func Fatal(message string, args ...interface{}) {
	if globalLogger == nil {
		log.Fatalf(message, args...)
		return
	}
	msg := fmt.Sprintf(message, args...)
	globalLogger.errorLogger.Println(globalLogger.formatMessage(FATAL, msg))
	CloseLogger()
	os.Exit(1)
}

// LogRequest logs HTTP request information
func LogRequest(method, path string, statusCode int, duration time.Duration, clientIP string) {
	if globalLogger == nil {
		return
	}

	var color string
	if globalLogger.useColor {
		switch {
		case statusCode >= 500:
			color = ColorRed
		case statusCode >= 400:
			color = ColorYellow
		case statusCode >= 300:
			color = ColorCyan
		case statusCode >= 200:
			color = ColorGreen
		default:
			color = ColorWhite
		}
	}

	message := fmt.Sprintf("%s%s %s %d%s - %v - %s",
		color, method, path, statusCode, ColorReset, duration, clientIP)

	globalLogger.infoLogger.Println(globalLogger.formatMessage(INFO, message))
}

// LogError logs error with stack trace
func LogError(err error, context string) {
	if globalLogger == nil || err == nil {
		return
	}

	// Get stack trace
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	stackTrace := string(buf[:n])

	message := fmt.Sprintf("%s: %v\nStack Trace:\n%s", context, err, stackTrace)
	globalLogger.errorLogger.Println(globalLogger.formatMessage(ERROR, message))
}

// RotateLogs rotates log files (can be called daily via cron)
func RotateLogs(logDir string) error {
//  validation path before rotation
	if err := validateLogPath(logDir); err != nil {
		return fmt.Errorf("invalid log directory: %w", err)
	}


	timestamp := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	archiveDir := filepath.Join(logDir, "archive")

	// Rotate app.log
	appLog := filepath.Join(logDir, "app.log")
	if _, err := os.Stat(appLog); err == nil {
		archiveFile := filepath.Join(archiveDir, fmt.Sprintf("app-%s.log", timestamp))
		if err := os.Rename(appLog, archiveFile); err != nil {
			return fmt.Errorf("failed to rotate app.log: %w", err)
		}
	}

	// Rotate error.log
	errorLog := filepath.Join(logDir, "error.log")
	if _, err := os.Stat(errorLog); err == nil {
		archiveFile := filepath.Join(archiveDir, fmt.Sprintf("error-%s.log", timestamp))
		if err := os.Rename(errorLog, archiveFile); err != nil {
			return fmt.Errorf("failed to rotate error.log: %w", err)
		}
	}

	return nil
}

// CleanOldLogs removes log files older than specified days
func CleanOldLogs(logDir string, daysToKeep int) error {
// Validation path
	if err := validateLogPath(logDir); err != nil {
		return fmt.Errorf("invalid log directory: %w", err)
	}
	

	archiveDir := filepath.Join(logDir, "archive")
	cutoffDate := time.Now().AddDate(0, 0, -daysToKeep)

	files, err := os.ReadDir(archiveDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoffDate) {
			filePath := filepath.Join(archiveDir, file.Name())
			if err := os.Remove(filePath); err != nil {
				Error("Failed to remove old log file: %s, error: %v", filePath, err)
			} else {
				Info("Removed old log file: %s", filePath)
			}
		}
	}

	return nil
}