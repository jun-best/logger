package logger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	//日志级别map
	level map[string]int
	//配置
	timeLayout string
	file       *FileManager
}

const MaxChannelBuffer = 1000

var (
	wg sync.WaitGroup
)

const (
	DefaultTimeLayout = "2006-01-02 15:04:05 CST"
	DefaultName       = "logger"
	DefaultLevel      = "INFO"
	DefaultPath       = "log"
	LevelInfo         = "INFO"
	LevelErr          = "ERR"
	LevelDebug        = "DEBUG"
)
const (
	LevelErrInt   = 1
	LevelInfoInt  = 2
	LevelDebugInt = 3
)

var levelMap = map[string]int{
	LevelInfo:  LevelInfoInt,
	LevelErr:   LevelErrInt,
	LevelDebug: LevelDebugInt,
}

func loadConf(fileName string) (*Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	conf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(conf, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
func Init(confFile string) (*Logger, error) {
	config, err := loadConf(confFile)
	if err != nil {
		return nil, err
	}
	l := &Logger{level: make(map[string]int),
		timeLayout: config.Format.TimeLayout}
	//check input
	if l.timeLayout == "" {
		l.timeLayout = DefaultTimeLayout
	}
	//get level
	l.setLevel(config.Format.Level)
	//file
	l.file, err = loadFile(config.File)
	if err != nil {
		return nil, err
	}
	if l.file.storeToFile {
		wg.Add(1)
		go l.file.writeToFile()
	}

	fmt.Println("success to init:", l.file.fileName)
	return l, err
}
func (l *Logger) setLevel(levelStr string) {
	if levelStr == "" {
		//默认为INFO
		levelStr = LevelInfo
	}
	level := levelMap[strings.ToUpper(levelStr)]
	if level >= LevelErrInt {
		l.level[LevelErr] = 1
	}
	if level >= LevelInfoInt {
		l.level[LevelInfo] = 1
	}
	if level >= LevelDebugInt {
		l.level[LevelDebug] = 1
	}
}
func (l *Logger) Close() {
	close(l.file.recvLogChan)
	wg.Wait()
}

func (l *Logger) Info(format string, args ...interface{}) {
	if _, ok := l.level[LevelInfo]; !ok {
		return
	}
	l.Print(format, LevelInfo, args...)
}
func (l *Logger) Error(format string, args ...interface{}) {
	if _, ok := l.level[LevelErr]; !ok {
		return
	}
	l.Print(format, LevelErr, args...)
}
func (l *Logger) Debug(format string, args ...interface{}) {
	if _, ok := l.level[LevelDebug]; !ok {
		return
	}
	l.Print(format, LevelDebug, args...)
}
func (l *Logger) Print(format, level string, args ...interface{}) {
	if ok := strings.Contains(format, "\n"); !ok {
		format = format + "\n"
	}
	_, file, line, _ := runtime.Caller(1)
	prefix := fmt.Sprintf("[%v] [%v] [%v:%v] %v",
		time.Now().Format(l.timeLayout), level, file, line, format)
	logStr := fmt.Sprintf(prefix, args...)
	fmt.Printf(logStr)
	if l.file.storeToFile {
		l.file.recvLogChan <- logStr
	}
}
