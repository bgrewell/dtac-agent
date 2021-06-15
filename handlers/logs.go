package handlers

import (
	"encoding/json"
	"errors"
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
	Extra map[string]interface{} `json:"extra"`
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

	// TODO: Maybe figure out a way to aggregate repeated log entries so
	//       that the log isn't spammed with a bunch of lines about a closed
	//       udp port or something like that. Instead showing the line along
	//		 with a count of how many of those lines there were. It could
	//		 either be done until a different line was found or grouped by time
	logs := Logs{}
	logstr := string(logbytes)
	logLines := strings.Split(logstr, "\n")
	for _, line := range logLines {
		var ll LogLine
		err = json.Unmarshal([]byte(line), &ll)
		if err == nil {
			// do a second unmarshal to capture any extra fields
			json.Unmarshal([]byte(line), &ll.Extra)
			// Remove the main fields from extra to avoid the duplication
			delete(ll.Extra, "file")
			delete(ll.Extra, "func")
			delete(ll.Extra, "level")
			delete(ll.Extra, "msg")
			delete(ll.Extra, "time")
			logs.Append(&ll)
		}
	}

	WriteResponseJSON(c, time.Now().Sub(start), &logs)
}

func GetLogsStreamHandler(c *gin.Context) {
	WriteErrorResponseJSON(c, errors.New("this method is not implemented"))
}
