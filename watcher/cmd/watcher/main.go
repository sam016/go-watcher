package main

import (
	"os"

	"github.com/sam016/go-watcher/watcher"
)

var (
	goversion string
	version   string
	commitID  string
	buildTime string
)

func main() {
	vinfo := watcher.VersionInfo{
		Version:   version,
		GoVersion: goversion,
		Commit:    commitID,
		BuildTime: buildTime,
	}

	config := watcher.ParseArgs(os.Args, vinfo)

	if config == nil {
		return
	}

	w := watcher.MustRegisterWatcher(config)

	r := watcher.NewRunner()

	// wait for changes and run the binary with given config
	go r.Run(config)
	b := watcher.NewBuilder(w, r)

	// build given package
	go b.Build(config)

	// listen for further changes
	go w.Watch()

	r.Wait()
}
