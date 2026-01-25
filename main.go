package main

import (
	"fmt"
	"embed"
	"encoding/json"
	"io/fs"
	"libreArca/internal/config"
	"net/http"
	"strings"
)

//go:embed ui/dist/*
var distFS embed.FS

func main() {
	config.LoadConfig("./config.yml")

	mux := http.NewServeMux()

	distSub, _ := fs.Sub(distFS, "ui/dist")

	mux.HandleFunc("/config.js", ServeConfigJS)
	
	fileServer := http.FileServer(http.FS(distSub))
	mux.Handle("/", SpaHandler(distSub, fileServer))
	http.ListenAndServe(":8088", mux)
}

func SpaHandler(staticFS fs.FS, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		cleanPath := strings.TrimPrefix(path, "/")
		f, err := staticFS.Open(cleanPath)

		if err != nil {
			r.URL.Path = "/"
			next.ServeHTTP(w, r)
			return
		}

		f.Close()

		next.ServeHTTP(w, r)
	})
}

func ServeConfigJS(w http.ResponseWriter, r *http.Request) {
	jsonBytes, _ := json.Marshal(config.AppConfig)

	jsCode := fmt.Sprintf("window.__WIKI_CONFIG__ = %s;", string(jsonBytes))

	w.Header().Set("Content-Type", "application/javascript")
	w.Write([]byte(jsCode))
}
