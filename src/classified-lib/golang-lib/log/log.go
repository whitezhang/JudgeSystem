/* log.go - encapsulation for log4go    */
/*
modification history
--------------------
2014/3/7, by Zhang Miao, create
2015/1/14, by Li Bingyi, modify, do not use "[%D %T] [%L] (%S) %M" as default log format,
    use log4go.LogFormat which could be set by log4go lib instead.
2014/3/7, by Wei He, modify, Add another WarnLogger specialized in writing to wf.log

*/
/*
DESCRIPTION
log: encapsulation for log4go

Usage:
    import "www.baidu.com/golang-lib/log"

    // Two log files will be generated in ./log:
    // test.log, and test.wf.log(for log > warn)
    // The log will rotate, and there is support for backup count
    log.Init("test", "INFO", "./log", true, "midnight", 5)

    log.Logger.Warn("warn msg")
    log.Logger.Info("info msg")

    // it is required, to work around bug of log4go
    time.Sleep(100 * time.Millisecond)
*/
package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

import "code.google.com/p/log4go"

/* global logger    */
var Logger log4go.Logger
var WarnLogger log4go.Logger
var initialized bool = false

/* logDirCreate(): check and create dir if nonexist   */
func logDirCreate(logDir string) error {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		/* create directory */
		err = os.MkdirAll(logDir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

/* filenameGen(): generate filename    */
func filenameGen(progName, logDir string, isErrLog bool) string {
	/* remove the last '/'  */
	strings.TrimSuffix(logDir, "/")

	var fileName string
	if isErrLog {
		/* for log file of warning, error, critical  */
		fileName = filepath.Join(logDir, progName+".wf.log")
	} else {
		/* for log file of all log  */
		fileName = filepath.Join(logDir, progName+".log")
	}

	return fileName
}

/* convert level in string to log4go level  */
func stringToLevel(str string) log4go.LevelType {
	var level log4go.LevelType

	str = strings.ToUpper(str)

	switch str {
	case "DEBUG":
		level = log4go.DEBUG
	case "TRACE":
		level = log4go.TRACE
	case "INFO":
		level = log4go.INFO
	case "WARNING":
		level = log4go.WARNING
	case "ERROR":
		level = log4go.ERROR
	case "CRITICAL":
		level = log4go.CRITICAL
	default:
		level = log4go.INFO
	}
	return level
}

/*
* Init - initialize log lib
*
* PARAMS:
*   - progName: program name. Name of log file will be progName.log
*   - levelStr: "DEBUG", "TRACE", "INFO", "WARNING", "ERROR", "CRITICAL"
*   - logDir: directory for log. It will be created if noexist
*   - hasStdOut: whether to have stdout output
*   - when:
*       "M", minute
*       "H", hour
*       "D", day
*       "MIDNIGHT", roll over at midnight
*   - backupCount: If backupCount is > 0, when rollover is done, no more than
*       backupCount files are kept - the oldest ones are deleted.
*
* RETURNS:
*   nil, if succeed
*   error, if fail
 */
func Init(progName string, levelStr string, logDir string,
	hasStdOut bool, when string, backupCount int) error {
	if initialized {
		return errors.New("Initialized Already")
	}

	/* check when   */
	if !log4go.WhenIsValid(when) {
		return fmt.Errorf("invalid value of when: %s", when)
	}

	/* check, and create dir if nonexist    */
	if err := logDirCreate(logDir); err != nil {
		log4go.Error("Init(), in logDirCreate(%s)", logDir)
		return err
	}

	/* convert level from string to log4go level    */
	level := stringToLevel(levelStr)

	/* create logger    */
	Logger = make(log4go.Logger)
	WarnLogger = make(log4go.Logger)

	/* create writer for stdout */
	if hasStdOut {
		Logger.AddFilter("stdout", level, log4go.NewConsoleLogWriter())
		WarnLogger.AddFilter("stdout", level, log4go.NewConsoleLogWriter())
	}

	/* create file writer for all log   */
	fileName := filenameGen(progName, logDir, false)
	logWriter := log4go.NewTimeFileLogWriter(fileName, when, backupCount)
	if logWriter == nil {
		return fmt.Errorf("error in log4go.NewTimeFileLogWriter(%s)", fileName)
	}
	logWriter.SetFormat(log4go.LogFormat)
	Logger.AddFilter("log", level, logWriter)

	/* create file writer for warning and fatal log */
	fileNameWf := filenameGen(progName, logDir, true)
	logWriter = log4go.NewTimeFileLogWriter(fileNameWf, when, backupCount)
	if logWriter == nil {
		return fmt.Errorf("error in log4go.NewTimeFileLogWriter(%s)", fileNameWf)
	}
	logWriter.SetFormat(log4go.LogFormat)
	WarnLogger.AddFilter("log_wf", log4go.WARNING, logWriter)

	initialized = true
	return nil
}

/*
* InitWithLogSvr - initialize log lib with remote log server
*
* PARAMS:
*   - progName: program name.
*   - levelStr: "DEBUG", "TRACE", "INFO", "WARNING", "ERROR", "CRITICAL"
*   - loggerName: logger name
*   - network: using "udp" or "unixgram"
*   - svrAddr: remote unix sock address for all logger
*   - svrAddrWf: remote unix sock address for warn/fatal logger
*                If svrAddrWf is empty string, no warn/fatal logger will be created.
*   - hasStdOut: whether to have stdout output
*
* RETURNS:
*   nil, if succeed
*   error, if fail
 */
func InitWithLogSvr(progName string, levelStr string, loggerName string,
	network string, svrAddr string, svrAddrWf string,
	hasStdOut bool) error {
	if initialized {
		return errors.New("Initialized Already")
	}

	/* convert level from string to log4go level    */
	level := stringToLevel(levelStr)

	/* create logger    */
	Logger = make(log4go.Logger)

	/* create writer for stdout */
	if hasStdOut {
		Logger.AddFilter("stdout", level, log4go.NewConsoleLogWriter())
	}

	/* create file writer for all log   */
	name := fmt.Sprintf("%s_%s", progName, loggerName)

	logWriter := log4go.NewPacketWriter(name, network, svrAddr, log4go.LogFormat)
	if logWriter == nil {
		return fmt.Errorf("error in log4go.NewPacketWriter(%s)", name)
	}
	Logger.AddFilter("log", level, logWriter)

	if len(svrAddrWf) > 0 {
		/* create file writer for warning and fatal log */
		logWriterWf := log4go.NewPacketWriter(name+".wf", network, svrAddrWf, log4go.LogFormat)
		if logWriterWf == nil {
			return fmt.Errorf("error in log4go.NewPacketWriter(%s, %s)",
				name, svrAddr)
		}
		Logger.AddFilter("log_wf", log4go.WARNING, logWriterWf)
	}

	initialized = true
	return nil
}
