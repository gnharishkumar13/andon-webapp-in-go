package view

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/user/andon-webapp-in-go/src/config"
)

var (
	m          sync.RWMutex
	templates  = map[string]*template.Template{}
	timestamps = map[string]time.Time{}
)

// PipelineBase contains fields that are expected to be present for
// shared templates, such as _layout.html.
type PipelineBase struct {
	Title string
}

func init() {
	setupViews()
	if config.Get().Environment == config.Environment {
		go watchForChanges()
	}
}

func setupViews() {
	log.Println("Loading view templates")
	viewRoot := config.Get().ViewRoot
	layout, err := template.ParseFiles(path.Join(viewRoot, "_layout.gohtml"))
	if err != nil {
		log.Printf("could not parse _layout.html: %v", err)
		return
	}
	_, err = layout.ParseFiles(
		path.Join(viewRoot, "_header.gohtml"),
		path.Join(viewRoot, "_footer.gohtml"),
	)
	if err != nil {
		log.Printf("could not parse _header.gohtml or _footer.gohtml: %v", err)
	}
	viewFIs, err := ioutil.ReadDir(path.Join(viewRoot, "content"))
	if err != nil {
		log.Printf("could not open view content directory: %v", err)
	}
	registerFunctions(layout)
	for _, fi := range viewFIs {
		f, err := os.Open(path.Join(viewRoot, "content", fi.Name()))
		if err != nil {
			log.Printf("failed to read content template %q: %v", fi.Name(), err)
		}
		content, err := ioutil.ReadAll(f)
		f.Close()
		if err != nil {
			log.Printf("failed to read content from template %q: %v", fi.Name(), err)
		}
		tmpl := template.Must(layout.Clone())
		_, err = tmpl.Parse(string(content))
		if err != nil {
			log.Printf("failed to parse template %q: %v", fi.Name(), err)
		}
		templates[strings.TrimSuffix(fi.Name(), ".gohtml")] = tmpl
	}
}

func watchForChanges() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		filesChanged := false
		filepath.Walk(config.Get().ViewRoot, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if strings.HasSuffix(info.Name(), ".gohtml") {
				t, ok := timestamps[info.Name()]
				if !ok || info.ModTime().After(t) {
					filesChanged = true
					timestamps[info.Name()] = info.ModTime()
					return nil
				}
			}
			return nil
		})
		if filesChanged {
			log.Print("Updated view detected, reloading")
			m.Lock()
			setupViews()
			m.Unlock()
		}
	}
}

// Get returns the view template stored with the provided key
// or an error if no template is found
func Get(key string) (*template.Template, error) {
	m.RLock()
	defer m.RUnlock()
	t, ok := templates[key]
	if !ok {
		return nil, fmt.Errorf("cannot find template with key %q", key)
	}
	return t, nil
}
