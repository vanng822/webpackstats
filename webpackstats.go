package webpackstats

import (
	"encoding/json"
	"html/template"
	"os"
)

type WebpackStats struct {
	Status  string                  `json:"status"`
	Chunks  map[string][]ChunkEntry `json:"chunks"`
	Error   string                  `json:"error"`
	Message string                  `json:"message"`
	File    string                  `json:"file"`
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

// WebpackUrlFuncMap ..
func WebpackUrlFuncMap(filename string) template.FuncMap {
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
