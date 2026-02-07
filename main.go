package main

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"libreArca/internal/config"
)

//go:embed all:ui/dist
var distFS embed.FS

//go:embed internal/defaults/config.yml
//go:embed internal/defaults/skins/default/layout.html
//go:embed internal/defaults/skins/default/style.css
var defaultsFS embed.FS

const (
	configPath    = "./config.yml"
	skinsDir      = "./skins"
	defaultSkin   = "default"
	manifestPath  = "ui/dist/.vite/manifest.json"
	defaultLayout = "internal/defaults/skins/default/layout.html"
	defaultStyle  = "internal/defaults/skins/default/style.css"
)

type manifestEntry struct {
	File    string   `json:"file"`
	Css     []string `json:"css"`
	IsEntry bool     `json:"isEntry"`
	Src     string   `json:"src"`
}

type layoutData struct {
	AppTitle      string
	ScriptPath    string
	StylePaths    []string
	SkinStylePath string
}

func main() {
	if err := ensureDefaults(); err != nil {
		log.Fatalf("init defaults: %v", err)
	}

	if err := config.LoadConfig(configPath); err != nil {
		log.Fatalf("load config: %v", err)
	}

	manifest, err := loadManifest(distFS)
	if err != nil {
		log.Fatalf("load manifest: %v", err)
	}

	scriptPath, stylePaths, err := entryPaths(manifest)
	if err != nil {
		log.Fatalf("resolve entry paths: %v", err)
	}

	skinName := config.AppConfig.Skin
	if skinName == "" {
		skinName = defaultSkin
	}
	skinDir := filepath.Join(skinsDir, skinName)
	layoutPath := filepath.Join(skinDir, "layout.html")

	mux := http.NewServeMux()

	assetsSub, err := fs.Sub(distFS, "ui/dist/assets")
	if err != nil {
		log.Fatalf("sub assets: %v", err)
	}

	mux.HandleFunc("/config.js", ServeConfigJS)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsSub))))
	mux.Handle("/skin-assets/", http.StripPrefix("/skin-assets/", http.FileServer(http.Dir(skinDir))))
	mux.HandleFunc("/", RootHandler(layoutPath, layoutData{
		AppTitle:      config.AppConfig.WikiName,
		ScriptPath:    scriptPath,
		StylePaths:    stylePaths,
		SkinStylePath: "/skin-assets/style.css",
	}))

	log.Println("listening on :8088")
	if err := http.ListenAndServe(":8088", mux); err != nil {
		log.Fatalf("listen: %v", err)
	}
}

func RootHandler(layoutPath string, data layoutData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(layoutPath)
		if err != nil {
			http.Error(w, "layout missing", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "template error", http.StatusInternalServerError)
			return
		}
	}
}

func ServeConfigJS(w http.ResponseWriter, r *http.Request) {
	jsonBytes, _ := json.Marshal(config.AppConfig)

	jsCode := fmt.Sprintf("window.__WIKI_CONFIG__ = %s;", string(jsonBytes))

	w.Header().Set("Content-Type", "application/javascript")
	w.Write([]byte(jsCode))
}

func ensureDefaults() error {
	if _, err := os.Stat(configPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		data, readErr := defaultsFS.ReadFile("internal/defaults/config.yml")
		if readErr != nil {
			return readErr
		}
		if writeErr := os.WriteFile(configPath, data, 0644); writeErr != nil {
			return writeErr
		}
	}

	if err := os.MkdirAll(skinsDir, 0755); err != nil {
		return err
	}

	defaultSkinDir := filepath.Join(skinsDir, defaultSkin)
	if err := os.MkdirAll(defaultSkinDir, 0755); err != nil {
		return err
	}

	if err := writeDefaultIfMissing(defaultSkinDir, "layout.html", defaultLayout); err != nil {
		return err
	}
	if err := writeDefaultIfMissing(defaultSkinDir, "style.css", defaultStyle); err != nil {
		return err
	}

	return nil
}

func writeDefaultIfMissing(dir, filename, embeddedPath string) error {
	filePath := filepath.Join(dir, filename)
	if _, err := os.Stat(filePath); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	data, err := defaultsFS.ReadFile(embeddedPath)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func loadManifest(dist embed.FS) (map[string]manifestEntry, error) {
	data, err := dist.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	var manifest map[string]manifestEntry
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return manifest, nil
}

func entryPaths(manifest map[string]manifestEntry) (string, []string, error) {
	entry, ok := manifest["src/main.ts"]
	if !ok {
		for _, candidate := range manifest {
			if candidate.IsEntry {
				entry = candidate
				ok = true
				break
			}
		}
	}

	if !ok {
		return "", nil, fmt.Errorf("entry not found in manifest")
	}

	if entry.File == "" {
		return "", nil, fmt.Errorf("entry file missing in manifest")
	}

	scriptPath := "/" + strings.TrimPrefix(entry.File, "/")
	stylePaths := make([]string, 0, len(entry.Css))
	for _, cssFile := range entry.Css {
		if cssFile == "" {
			continue
		}
		stylePaths = append(stylePaths, "/"+strings.TrimPrefix(cssFile, "/"))
	}

	return scriptPath, stylePaths, nil
}
