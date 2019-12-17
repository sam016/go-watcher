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

	params := watcher.ParseArgs(os.Args, vinfo)

	if params == nil {
		return
	}

	w := watcher.MustRegisterWatcher(params)

	r := watcher.NewRunner()

	// wait for build and run the binary with given params
	go r.Run(params)
	b := watcher.NewBuilder(w, r)

	// build given package
	go b.Build(params)

	// listen for further changes
	go w.Watch()

	r.Wait()
}
