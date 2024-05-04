package webpackstats

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"time"
)

type WebpackStats struct {
	Status  string              `json:"status"`
	Chunks  map[string][]string `json:"chunks"`
	Assets  map[string]Asset    `json:"assets"`
	Error   string              `json:"error"`
	Message string              `json:"message"`
	File    string              `json:"file"`
}

// GetUrl returns the public path of an chunk
func (wps *WebpackStats) GetUrl(name string) string {
	if asset, ok := wps.Assets[name]; ok {
		if asset.PublicPath != "" {
			return asset.PublicPath
		}
		return asset.Name
	}
	return ""
}

// general output
type Asset struct {
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
	webp, err := parse(file)
	if err != nil {
		panic(err)
	}
	return webp
}

func parse(file io.Reader) (*WebpackStats, error) {
	decoder := json.NewDecoder(file)
	webp := WebpackStats{}
	err := decoder.Decode(&webp)
	if err != nil {
		return nil, err
	}
	return &webp, nil
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

var templateFuncs = template.FuncMap{
	"webpackUrl": func(name string) string {
		wps := Get()
		if wps == nil {
			return ""
		}
		chunks, ok := wps.Chunks[name]

		if !ok {
			return ""
		}

		if len(chunks) != 1 {
			panic("We have more then one chunk, please use webpackUrls or webpackScripts")
		}

		return wps.GetUrl(chunks[0])
	},
	"webpackUrls": func(name string) []string {
		result := []string{}
		wps := Get()
		if wps == nil {
			return result
		}
		chunks, ok := wps.Chunks[name]
		if !ok {
			return result
		}

		for _, chunk := range chunks {
			result = append(result, wps.GetUrl(chunk))
		}
		return result
	},
	"webpackScripts": func(name string) template.HTML {
		wps := Get()
		if wps == nil {
			return ""
		}
		chunks, ok := wps.Chunks[name]
		if !ok {
			return ""
		}
		result := ""
		for _, chunk := range chunks {
			result = fmt.Sprintf("%s<script src=\"%s\" type=\"text/javascript\"></script>", result, wps.GetUrl(chunk))
		}
		return template.HTML(result)
	},
}

// WebpackURL is a FuncMap for template function
// to map the resource url into the resource hash url
func WebpackURL(file io.Reader) template.FuncMap {
	webp, err := parse(file)
	if err != nil {
		panic(err)
	}
	if webp.Status == "error" {
		panic(webp)
	}
	Set(webp)
	return templateFuncs
}

// WebpackURLFuncMap a FuncMap with template function webpackUrl(name string) string
func WebpackURLFuncMap(filename string) template.FuncMap {
	go LoadStats(filename)
	return templateFuncs
}
