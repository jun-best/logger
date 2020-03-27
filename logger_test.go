package logger

import (
	"fmt"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	type TestTable struct {
		confFile  string
		formatStr string
		args      []interface{}
		wantStr   string
	}
	tests := []TestTable{
		{confFile: "logger.json", formatStr: "just info test", wantStr: "just test"},
	}

	for k, test := range tests {
		fmt.Println("test:", k)
		l, err := Init(test.confFile)
		if err != nil {
			fmt.Println("fail to init", err.Error())
			continue
		}
		l.Info("test-info")
		l.Error("test-err")
		l.Debug("test-debug")
		time.Sleep(time.Second * 65)
		l.Close()
	}
}
