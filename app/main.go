package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Response struct {
	Headers map[string]string `yaml:"headers"`
	Body    interface{}       `yaml:"body"`
	Status  int               `yaml:"status"`
}

type Mock struct {
	Method       string    `yaml:"method"`
	Path         string    `yaml:"path"`
	Response     Response  `yaml:"response"`
}

const templatesDir = "./templates"

func loadMocks() []Mock {
	var all []Mock

	files, err := filepath.Glob(filepath.Join(templatesDir, "*.yaml"))
	if err != nil {
		slog.Warn("glob failed", "err", err)
		return all
	}
	ymlFiles, err := filepath.Glob(filepath.Join(templatesDir, "*.yml"))
	if err != nil {
		slog.Warn("glob failed", "err", err)
		return all
	}
	files = append(files, ymlFiles...)
	sort.Strings(files)

	if len(files) == 0 {
		slog.Warn("no templates found", "dir", templatesDir)
		return all
	}

	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			slog.Warn("read failed", "file", f, "err", err)
			continue
		}
		var mocks []Mock
		if err := yaml.Unmarshal(data, &mocks); err != nil {
			slog.Warn("yaml parse failed", "file", f, "err", err)
			continue
		}
		all = append(all, mocks...)
	}
	return all
}

// matchPath: exact match, or wildcard '*' matching exactly one segment.
func matchPath(pattern, path string) (exact bool, ok bool) {
	if pattern == path {
		return true, true
	}
	if !strings.Contains(pattern, "*") {
		return false, false
	}
	ps := strings.Split(pattern, "/")
	xs := strings.Split(path, "/")
	if len(ps) != len(xs) {
		return false, false
	}
	for i := range ps {
		if ps[i] == "*" {
			if xs[i] == "" {
				return false, false
			}
			continue
		}
		if ps[i] != xs[i] {
			return false, false
		}
	}
	return false, true
}

func findMock(mocks []Mock, method, path string) *Mock {
	// pass 1: exact
	for i := range mocks {
		if !strings.EqualFold(mocks[i].Method, method) {
			continue
		}
		if exact, ok := matchPath(mocks[i].Path, path); ok && exact {
			return &mocks[i]
		}
	}
	// pass 2: wildcard
	for i := range mocks {
		if !strings.EqualFold(mocks[i].Method, method) {
			continue
		}
		if exact, ok := matchPath(mocks[i].Path, path); ok && !exact {
			return &mocks[i]
		}
	}
	return nil
}

func writeBody(w http.ResponseWriter, m *Mock) {
	if m.Response.Body == nil {
		return
	}
	switch v := m.Response.Body.(type) {
	case string:
		fmt.Fprint(w, v)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			slog.Error("marshal failed", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	mocks := loadMocks()
	m := findMock(mocks, r.Method, r.URL.Path)

	if m == nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "hello world")
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "match", false, "status", 200)
		return
	}

	code := m.Response.Status
	if code == 0 {
		code = http.StatusOK
	}

	hasCT := false
	for k, v := range m.Response.Headers {
		if strings.EqualFold(k, "Content-Type") {
			hasCT = true
		}
		w.Header().Set(k, v)
	}
	if !hasCT {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(code)
	writeBody(w, m)
	slog.Info("request", "method", r.Method, "path", r.URL.Path, "match", true, "status", code)
}

func main() {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "80"
	}
	addr := ":" + port

	http.HandleFunc("/", handler)
	slog.Info("starting mock server", "addr", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}