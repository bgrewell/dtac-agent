package httprouting

import (
	"crypto/sha1"
	"fmt"
	"github.com/BGrewell/system-api/configuration"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
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

func AddCustomHandlers(c *configuration.Config, r *gin.Engine) {
	customHandlerMap = make(map[string]*CustomHandlerSettings)
	for _, entry := range c.Custom {
		for _, value := range entry {
			route := fmt.Sprintf("/custom/%s", value.Name)
			if value.Route != "" {
				route = value.Route
			}
			//TODO: Check for existing key in map, if found warn user that the existing route is being shadowed
			customHandlerMap[route] = &CustomHandlerSettings{
				Value:    value.Source.Value,
				Type:     value.Source.Type,
				Settings: value.Blocking,
			}
			switch value.Source.Type {
			case "file":
				r.GET(route, CustomFileHandler)
			case "net":
				r.GET(route, CustomNetHandler)
			case "cmd":
				r.GET(route, CustomCmdHandler)
			default:
				// todo: need to push logging throughout the code
				log.Printf("unrecognized custom route source: %v. skipping", value.Source.Type)
			}
		}
	}
}

//TODO: Needs to be updated for new configuration parameters
func CustomFileHandler(c *gin.Context) {
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

	start := time.Now()
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

		if time.Now().Sub(start).Milliseconds() > int64(entry.Settings.TimeoutMs) {
			body = string(contents)
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"file": value, "contents": body, "stale": stale})
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
