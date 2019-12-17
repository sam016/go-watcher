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

	wtcher := watcher.MustRegisterWatcher(config)

	debugger := watcher.NewDebugger()

	// wait for changes and run the binary with given config
	go debugger.Debug(config)
	builder := watcher.NewBuilder(wtcher, debugger)

	// build given package
	go builder.Build(config)

	// listen for further changes
	go wtcher.Watch()

	debugger.Wait()
}
