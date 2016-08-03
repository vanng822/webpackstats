package webpackstats

import (
	"encoding/json"
	"html/template"
	"os"
)

type webpackStats struct {
	Status  string                  `json:"status"`
	Chunks  map[string][]chunkEntry `json:"chunks"`
	Error   string                  `json:"error"`
	Message string                  `json:"message"`
}

// general output
type chunkEntry struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	PublicPath string `json:"publicPath"`
}

func load(filename string) *webpackStats {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	webp := webpackStats{}
	err = decoder.Decode(&webp)
	if err != nil {
		panic(err)
	}
	return &webp
}

// WebpackStats ..
func WebpackStats(filename string) template.FuncMap {
	webp := load(filename)
	if webp.Status != "done" {
		panic(webp)
	}
	var templateFuncs = template.FuncMap{
		"webpackUrl": func(name string) string {
			if entry, ok := webp.Chunks[name]; ok {
				if len(entry) == 1 {
					if entry[0].PublicPath != "" {
						return entry[0].PublicPath
					}
					return entry[0].Name
				}
			}
			return ""
		},
	}
	return templateFuncs
}
