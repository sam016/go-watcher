package watcher

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
)

// Builder composes of both debugger and watcher. Whenever watcher gets notified, builder starts a build process, and forces the debugger to restart
type Builder struct {
	debugger *Debugger
	watcher  *Watcher
}

// NewBuilder constructs the Builder instance
func NewBuilder(w *Watcher, r *Debugger) *Builder {
	return &Builder{watcher: w, debugger: r}
}

// Build listens watch events from Watcher and sends messages to Debugger
// when new changes are built.
func (b *Builder) Build(appConfig *AppConfig) {
	go b.registerSignalHandler()
	go func() {
		// used for triggering the first build
		b.watcher.update <- struct{}{}
	}()

	for range b.watcher.Wait() {
		pkg := appConfig.packagePath()

		log.Println("Change(s) detected")
		color.Cyan("Starting debugger for %s...\n", pkg)

		// and start the new process
		b.debugger.restart(pkg)
	}
}

func (b *Builder) registerSignalHandler() {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signals
	b.watcher.Close()
	b.debugger.Close()
}
