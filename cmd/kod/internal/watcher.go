package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/samber/lo"
)

type Watcher interface {
	Add(string) error
	Events() chan fsnotify.Event
	Errors() chan error
}

type watcher struct {
	w *fsnotify.Watcher
}

func (w *watcher) Add(name string) error {
	return w.w.Add(name)
}

func (w *watcher) Events() chan fsnotify.Event {
	return w.w.Events
}

func (w *watcher) Errors() chan error {
	return w.w.Errors
}

// Watch watches the directory and calls the callback function when a file is modified.
func Watch(watcher Watcher, dir string, callback func()) {

	lo.Must0(filepath.Walk(dir, func(path string, info os.FileInfo, _ error) error {
		if info != nil && info.IsDir() {
			if isHiddenDirectory(path) {
				return filepath.SkipDir
			}

			err := watcher.Add(path)
			if err != nil {
				return filepath.SkipDir
			}

			fmt.Println("watching", path)
		}

		return nil
	}))

	stop := make(chan struct{}, 1)
	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events():
				if !ok {
					stop <- struct{}{}
					return
				}

				if !validEvent(event) {
					continue
				}

				if strings.HasPrefix(filepath.Base(event.Name), "kod_gen") {
					continue
				}

				if !strings.HasSuffix(filepath.Base(event.Name), ".go") {
					continue
				}

				log.Println("modified file:", event.Name)
				callback()
			case err, ok := <-watcher.Errors():
				if !ok {
					stop <- struct{}{}
					return
				}
				log.Println("error:", err)
				stop <- struct{}{}
			}
		}
	}()

	// Block main goroutine forever.
	<-stop
}

func isHiddenDirectory(path string) bool {
	return len(path) > 1 && strings.HasPrefix(path, ".") && filepath.Base(path) != ".."
}

func validEvent(ev fsnotify.Event) bool {
	return ev.Op&fsnotify.Create == fsnotify.Create ||
		ev.Op&fsnotify.Write == fsnotify.Write ||
		ev.Op&fsnotify.Remove == fsnotify.Remove
}
