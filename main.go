package main

import (
	"embed"
	"strings"
	"io/fs"
	"net/http"
)

//go:embed ui/dist/*
var distFS embed.FS

func main() {
	mux := http.NewServeMux()

	distSub, _ := fs.Sub(distFS, "ui/dist")
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
