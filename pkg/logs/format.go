package logs

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yi-cloud/rest-server/common"
	"path/filepath"
)

type NewFormatter struct{}

func (m *NewFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	if entry.Buffer != nil {
		b =
			entry.Buffer
	}

	timestamp := entry.Time.Format(common.TimeFormat)
	var newLog string

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog =
			fmt.Sprintf("[%s] [%s] [%s:%d] %s\n",
				timestamp, entry.Level, fName, entry.Caller.Line, entry.Message)
	} else {
		newLog =
			fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}
