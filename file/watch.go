package file

import (
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/golang/glog"
)

type action func(filename string)

type Watch struct {
	file      string
	addAction action
	modAction action
	delAction action
}

func NewWatch(filename string, add, mod, del action) *Watch {
	return &Watch{
		file:      filename,
		addAction: add,
		modAction: mod,
		delAction: del,
	}
}

func (w *Watch) Watch() error {
	configFile := filepath.Clean(w.file)
	configDir, _ := filepath.Split(configFile)
	realConfigFile, err := filepath.EvalSymlinks(w.file)
	if err != nil {
		return fmt.Errorf("EvalSymlinks error:%v", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("NewWatcher error:%v", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok { // 'Events' channel is closed
					glog.Error("Events' channel is closed")
					return
				}
				currentConfigFile, _ := filepath.EvalSymlinks(w.file)
				// we only care about the config file with the following cases:
				// 1 - if the config file was modified or created
				// 2 - if the real path to the config file changed (eg: k8s ConfigMap replacement)
				const writeOrCreateMask = fsnotify.Write | fsnotify.Create
				if (filepath.Clean(event.Name) == configFile &&
					event.Op&writeOrCreateMask != 0) ||
					(currentConfigFile != "" && currentConfigFile != realConfigFile) {
					realConfigFile = currentConfigFile
				} else if filepath.Clean(event.Name) == configFile &&
					event.Op&fsnotify.Remove != 0 {
					return
				}

			case err, ok := <-watcher.Errors:
				if ok { // 'Errors' channel is not closed
					glog.Printf("watcher error: %v\n", err)
				}
				return
			}
		}
	}()
	watcher.Add(configDir)
	return nil
}
