package handlers

import (
	"encoding/json"
	. "github.com/BGrewell/system-api/common"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"runtime"
	"strings"
	"time"
)

type Logs struct {
	Lines []*LogLine
}

func (l *Logs) Append(line *LogLine) {
	if l.Lines == nil {
		l.Lines = make([]*LogLine, 0)
	}
	l.Lines = append(l.Lines, line)
}

type LogLine struct {
	File string `json:"file"`
	Func string `json:"func"`
	Level string `json:"level"`
	Msg string `json:"msg"`
	Time string `json:"time"`
}

func GetLogsHandler(c *gin.Context) {
	start := time.Now()
	filename := "/var/log/system-apid/system-apid.log"
	if runtime.GOOS == "windows" {
		filename = "C:\\Logs\\system-apid.log"
	}

	logbytes, err := ioutil.ReadFile(filename)
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}

	logs := Logs{}
	logstr := string(logbytes)
	logLines := strings.Split(logstr, "\n")
	for _, line := range logLines {
		var ll LogLine
		err = json.Unmarshal([]byte(line), &ll)
		if err == nil {
			logs.Append(&ll)
		}
	}

	WriteResponseJSON(c, time.Now().Sub(start), &logs)
}
