package webpackstats

import (
	"encoding/json"
	"html/template"
	"os"
	"time"
)

type WebpackStats struct {
	Status  string                  `json:"status"`
	Chunks  map[string][]ChunkEntry `json:"chunks"`
	Error   string                  `json:"error"`
	Message string                  `json:"message"`
	File    string                  `json:"file"`
}

// GetUrl returns the public path of an chunk
func (wps *WebpackStats) GetUrl(name string) string {
	if entry, ok := wps.Chunks[name]; ok {
		if len(entry) == 1 {
			if entry[0].PublicPath != "" {
				return entry[0].PublicPath
			}
			return entry[0].Name
		}
	}
	return ""
}

// general output
type ChunkEntry struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	PublicPath string `json:"publicPath"`
}

func load(filename string) *WebpackStats {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	webp := WebpackStats{}
	err = decoder.Decode(&webp)
	if err != nil {
		panic(err)
	}
	return &webp
}

var (
	_webpackStats *WebpackStats
)

// Set webpack stats
func Set(wps *WebpackStats) {
	_webpackStats = wps
}

// Get returning current webpack stats
func Get() *WebpackStats {
	return _webpackStats
}

// LoadStats fetch stats async, if build process going on it will wait
func LoadStats(filename string) {
	var webp *WebpackStats
	// this is for development
	// the file should be built before building go
	ticker := time.Tick(60 * time.Second)
readLoop:
	for {
		select {
		case <-ticker:
			panic("Timeout reading webpack stats file")
		default:
			webp = load(filename)
			if webp.Status == "done" {
				break readLoop
			}
			if webp.Status == "error" {
				panic(webp)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
	Set(webp)
}

// WebpackURLFuncMap a FuncMap with template function webpackUrl(name string) string
func WebpackURLFuncMap(filename string) template.FuncMap {
	go LoadStats(filename)
	var templateFuncs = template.FuncMap{
		"webpackUrl": func(name string) string {
			wps := Get()
			if wps == nil {
				return ""
			}
			return wps.GetUrl(name)
		},
	}
	return templateFuncs
}
