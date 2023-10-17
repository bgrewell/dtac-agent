package version

import (
	"encoding/json"
	"fmt"
	"time"
)

var (
	date    = time.Now().Format("2006-01-02 15:04:05")
	rev     = "DEBUG"
	branch  = "DEBUG"
	version = "DEBUG"
)

// Info is the struct for the version information
type Info struct {
	Version string `json:"version" yaml:"version"`
	Date    string `json:"date" yaml:"date"`
	Branch  string `json:"branch" yaml:"branch"`
	Rev     string `json:"rev" yaml:"rev"`
}

// String returns the version information as a string
func (i Info) String() string {
	return fmt.Sprintf("version: %s date: %s branch: %s rev: %s", i.Version, i.Date, i.Branch, i.Rev)
}

// JSON returns the version information as a JSON string
func (i Info) JSON(pretty bool) string {
	var j []byte
	var err error
	if pretty {
		j, err = json.MarshalIndent(i, "", "  ")
	} else {
		j, err = json.Marshal(i)
	}
	if err != nil {
		return "Error marshaling version info"
	}
	return string(j)
}

// Current returns the current version information
func Current() *Info {
	v := Info{
		Version: version,
		Date:    date,
		Branch:  branch,
		Rev:     rev,
	}
	return &v
}
