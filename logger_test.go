package logger

import (
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	l, err := Init("logger.json")
	if err != nil {
		t.Error("fail to init", err.Error())
		return
	}
	l.Info("test-info")
	l.Error("test-err")
	l.Debug("test-debug")
	time.Sleep(time.Second * 65)
	l.Close()
	t.Log("exit test")
}
