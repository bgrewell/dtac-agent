package handlers

import (
	"crypto/sha1"
	"fmt"
	"github.com/BGrewell/system-agent/configuration"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	customHandlerMap map[string]*CustomHandlerSettings
)

type CustomHandlerSettings struct {
	Value      string
	Type       string
	Settings   *configuration.BlockingEntry
	LastHash   string
	LastAccess string
}

func init() {
	customHandlerMap = make(map[string]*CustomHandlerSettings)
}

func AddCustomHandler(route string, settings *CustomHandlerSettings) {
	customHandlerMap[route] = settings
}

//TODO: Needs to be updated for new configuration parameters
func CustomFileHandler(c *gin.Context) {
	start := time.Now()
	// Get path value to map to a file entry
	key := c.Request.URL.Path
	//key := filepath.Base(c.Request.URL.Path)
	value := ""
	var entry *CustomHandlerSettings
	var ok bool
	if entry, ok = customHandlerMap[key]; ok {
		value = entry.Value
	}

	// Check for file changes
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "time": time.Now().Format(time.RFC3339Nano)})
		return
	}
	defer watcher.Close()
	watcher.Add(value)

	watchStart := time.Now()
	body := ""
	stale := true
	for {
		// Check to see if contents changed
		changed, contents, h, _ := contentsChanged(value, entry.LastHash) //todo: ignoring err here could lead to silent failures

		// If they changed then store the change and break out of the for loop
		if changed {
			stale = false
			body = string(contents)
			entry.LastHash = h
			break
		}

		// If they didn't then setup a fsnotify watcher to watch for the file to change and if it does then check again
		select {
		case <-watcher.Events:
			// if the event fires break out of the select and the for loop will check for file changes
			break
		case err := <-watcher.Errors:
			log.Printf("fsnotify error: %v", err)
		case <-time.After(time.Duration(entry.Settings.TimeoutMs) * time.Millisecond):
			break
		}

		if time.Now().Sub(watchStart).Milliseconds() > int64(entry.Settings.TimeoutMs) {
			body = string(contents)
			break
		}
	}
	WriteResponseJSON(c, time.Since(start), gin.H{"file": value, "contents": body, "stale": stale})
}

func contentsChanged(filepath string, previousHash string) (changed bool, contents []byte, hashVal string, err error) {
	contents, err = ioutil.ReadFile(filepath)
	if err != nil {
		return false, nil, "", err
	}
	h := sha1.New()
	h.Write(contents)
	hashVal = fmt.Sprintf("%x", h.Sum(nil))
	return previousHash != hashVal, contents, hashVal, nil
}

func CustomNetHandler(c *gin.Context) {
	// Get path value to map to a file entry
	key := c.Request.URL.Path
	//key := filepath.Base(c.Request.URL.Path)
	value := ""
	var entry *CustomHandlerSettings
	var ok bool
	if entry, ok = customHandlerMap[key]; ok {
		value = entry.Value
	}
	fmt.Println(value)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "this method has not been implemented", "time": time.Now().Format(time.RFC3339Nano)})
}

func CustomCmdHandler(c *gin.Context) {
	start := time.Now()
	// Get path value to map to a file entry
	key := c.Request.URL.Path
	//key := filepath.Base(c.Request.URL.Path)
	value := ""
	var entry *CustomHandlerSettings
	var ok bool
	if entry, ok = customHandlerMap[key]; ok {
		value = entry.Value
	}
	stdout, stderr, err := execute.ExecuteCmdEx(value)
	stdout = strings.Trim(stdout, "\r\n")
	stderr = strings.Trim(stderr, "\r\n")
	WriteResponseJSON(c, time.Since(start), gin.H{"stdout": stdout, "stderr": stderr, "err": err})
}
