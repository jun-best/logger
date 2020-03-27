package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	RotationTypeSizeStr  = "size"
	RotationTypeTimeStr  = "time"
	RotationTypeSize     = 1
	RotationTypeTime     = 2
	RotaitionTimeMinStr  = "MIN"
	RotaitionTimeHourStr = "HOUR"
	RotaitionTimeDayStr  = "DAY"
	RotaitionTimeMin     = 1
	RotaitionTimeHour    = 2
	RotaitionTimeDay     = 3
	DefaultRotation      = "size"
	DefaultSize          = "10m"
	RotationTimeLayout   = "2006-01-02_15:04:05"
)

type FileManager struct {
	file             *os.File
	storeToFile      bool
	fileName         string
	rotationType     int
	rotationSize     int64
	rotationTime     time.Duration
	rotationTimeType int
	recvLogChan      chan string
}

func loadFile(config FileConfig) (*FileManager, error) {
	if !config.PrintToFile {
		return &FileManager{storeToFile: false}, nil
	}
	file := &FileManager{recvLogChan: make(chan string, MaxChannelBuffer),
		storeToFile: true}
	name := config.Name
	if name == "" {
		name = DefaultName
	}
	path := strings.TrimSuffix(config.Path, "/")
	if path == "" {
		path = DefaultPath
	}
	//mkdir log path
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil && err != os.ErrExist {
			return nil, err
		}
	}
	file.fileName = path + "/" + name
	//rotation
	rotationTypeStr := strings.ToLower(config.RotationType)
	switch rotationTypeStr {
	case RotationTypeSizeStr:
		file.rotationType = RotationTypeSize
		if strings.Contains(config.RotationSize, "K") {
			value, _ := strconv.Atoi(strings.TrimSuffix(config.RotationSize, "K"))
			file.rotationSize = int64(value * 1 << 10)
		} else if strings.Contains(config.RotationSize, "M") {
			value, _ := strconv.Atoi(strings.TrimSuffix(config.RotationSize, "M"))
			file.rotationSize = int64(value * 1 << 20)
		} else {
			value, _ := strconv.Atoi(config.RotationSize)
			file.rotationSize = int64(value)
		}
		//default is time
	default:
		file.rotationType = RotationTypeTime
		timeStr := strings.ToUpper(config.RotationTime)
		switch timeStr {
		case RotaitionTimeMinStr:
			file.rotationTime = time.Minute
			file.rotationTimeType = RotaitionTimeMin
		case RotaitionTimeHourStr:
			file.rotationTime = time.Hour
			file.rotationTimeType = RotaitionTimeHour
		default:
			//default is day
			file.rotationTime = time.Hour * 24
			file.rotationTimeType = RotaitionTimeDay
		}
	}
	return file, nil
}
func (f *FileManager) writeToFile() {
	defer func() {
		fmt.Println("success to exit logger")
		wg.Done()
	}()
	var err error
	f.file, err = os.OpenFile(f.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		fmt.Println("fail to create", err.Error())
		return
	}
	var timer *time.Ticker
	if f.rotationType == RotationTypeSize {
		//100毫秒检查一次
		timer = time.NewTicker(time.Millisecond * 100)
	} else {
		timer = time.NewTicker(f.rotationTime)
	}

	for {
		select {
		case str, alive := <-f.recvLogChan:
			if !alive {
				return
			}
			f.file.WriteString(str)
		case <-timer.C:
			if f.rotationType == RotationTypeSize {
				//check file size
				fileInfo, _ := os.Stat(f.fileName)
				if fileInfo.Size() > f.rotationSize {
					//do rotation
					file := doRotation(f.fileName)
					f.file.Close()
					f.file = file
				}
			} else {
				//do rotation
				file := doRotation(f.fileName)
				f.file.Close()
				f.file = file
			}
		}

	}
}

func doRotation(fileName string) *os.File {
	newName := fileName + "_" + time.Now().Format(RotationTimeLayout)
	os.Rename(fileName, newName)
	fmt.Println("rotation:", fileName, newName)
	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	return file
}
