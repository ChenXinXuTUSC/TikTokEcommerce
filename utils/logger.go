package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

type TtLogger struct {
	cnslLogger *log.Logger
	fileLogger *log.Logger
	cnslPrefix string
	filePrefix string
	lvl    int
}

func (l *TtLogger) Println(v ...any) {
	fileline := " " + GetFileLine()
	msg := append(v, fileline)
	l.cnslLogger.Print(msg...)
	if l.fileLogger != nil {
		l.fileLogger.Print(msg...)
	}
}

func (l *TtLogger) Printf(format string, v ...any) {
	fileline := GetFileLine()
	msg := append(v, fileline)
	format += " %s"
	l.cnslLogger.Printf(format, msg...)
	if l.fileLogger != nil {
		l.fileLogger.Printf(format, msg...)
	}
}

func GetFileLine() string {
	// 获取调用者的文件和行号
	// skip = 2：跳过 logWithCaller 和调用 logWithCaller 的函数
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	return fmt.Sprintf("%s:%d", file, line)
}




func init() {
	for i := range loggers {
		loggers[i].cnslLogger = log.New(os.Stderr, loggers[i].cnslPrefix+" ", log.LstdFlags)
	}
}

const (
	dbugClrFmt = "\033[1;34m[%s]\033[0m"
	infoClrFmt = "\033[1;32m[%s]\033[0m"
	warnClrFmt = "\033[1;33m[%s]\033[0m"
	erroClrFmt = "\033[1;31m[%s]\033[0m"
	fatlClrFmt = "\033[1;35m[%s]\033[0m"
)

var (
	Dbug = &TtLogger{
		cnslPrefix: fmt.Sprintf(dbugClrFmt, "DBUG"),
		filePrefix: "[DBUG]",
		lvl:    0,
	}
	Info = &TtLogger{
		cnslPrefix: fmt.Sprintf(infoClrFmt, "INFO"),
		filePrefix: "[INFO]",
		lvl:    1,
	}
	Warn = &TtLogger{
		cnslPrefix: fmt.Sprintf(warnClrFmt, "WARN"),
		filePrefix: "[WARN]",
		lvl:    2,
	}
	Erro = &TtLogger{
		cnslPrefix: fmt.Sprintf(erroClrFmt, "ERRO"),
		filePrefix: "[ERRO]",
		lvl:    3,
	}
	Fatl = &TtLogger{
		cnslPrefix: fmt.Sprintf(fatlClrFmt, "FATL"),
		filePrefix: "[FATL]",
		lvl:    4,
	}
	loggers = []*TtLogger{
		Dbug,
		Info,
		Warn,
		Erro,
		Fatl,
	}
)

var mu sync.Mutex

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, TtLogger := range loggers {
		if TtLogger.lvl < level {
			TtLogger.cnslLogger.SetOutput(io.Discard)
			if TtLogger.fileLogger != nil {
				TtLogger.fileLogger.SetOutput(io.Discard)
			}
		}
	}
}

func SetPrefix(prefix string) {
	mu.Lock()
	defer mu.Unlock()

	for _, TtLogger := range loggers {
		TtLogger.cnslLogger.SetPrefix(TtLogger.cnslPrefix + " " + prefix + " ")
		if TtLogger.fileLogger != nil {
			TtLogger.fileLogger.SetPrefix(TtLogger.filePrefix + " " + prefix + " ")
		}
	}
}

func SetOutputFile(path string) error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	for i:= range loggers {
		if loggers[i].fileLogger == nil {
			loggers[i].fileLogger = log.New(io.Writer(file), loggers[i].filePrefix+" ", log.LstdFlags)
		} else {
			loggers[i].fileLogger.SetOutput(io.Writer(file))
		}
	}

	return nil
}
