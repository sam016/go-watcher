// Package watcher is a command line tool inspired by fresh (https://github.com/pilu/fresh) and used
// for watching .go file changes, and restarting the app in case of an update/delete/add operation.
// After you installed it, you can run your apps with their default parameters as:
// watcher -c config -p 7000 -h localhost
package watcher

import (
	"errors"
	"log"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

// Debugger listens for the change events and depending on that kills
// the obsolete process, and runs a new one
type Debugger struct {
	start chan string
	done  chan struct{}
	cmd   *exec.Cmd
}

// NewDebugger creates a new Debugger instance and returns its pointer
func NewDebugger() *Debugger {
	return &Debugger{
		start: make(chan string),
		done:  make(chan struct{}),
	}
}

// Debug initializes debugger with given parameters.
func (r *Debugger) Debug(appConfig *AppConfig) {
	for pkg := range r.start {

		color.Green("Running %s...\n", appConfig.Watcher.Run)

		args := []string{"debug", pkg}
		args = append(args, appConfig.DelveArgs...)
		args = append(args, "--")
		args = append(args, appConfig.Package.Args...)

		cmd, err := runCommand("dlv", args...)
		if err != nil {
			log.Printf("Could not run the go binary: %s \n", err)
			r.kill(cmd)

			continue
		}

		r.cmd = cmd

		go func(cmd *exec.Cmd) {
			if err := cmd.Wait(); err != nil {
				log.Printf("debug process interrupted: %s \n", err)
				r.kill(cmd)
			}
		}(r.cmd)
	}
}

// Restart kills the process, removes the old binary and
// restarts the new process
func (r *Debugger) restart(pkg string) {
	r.kill(r.cmd)

	r.start <- pkg
}

func (r *Debugger) kill(cmd *exec.Cmd) {
	if cmd != nil {
		var success bool = false

		pidDebugBin, _ := getPidByName("/tmp/__debug_bin")
		pidCmd := string(cmd.Process.Pid)

		// cmd.Process.Kill()

		if pidDebugBin != "" {
			log.Println("Killing debug processes:", pidDebugBin, pidCmd)

			outKill, err := exec.Command("kill", pidDebugBin, pidCmd).Output()

			if err != nil {
				log.Println("kill:err", err, string(outKill))
			} else {
				success = true
			}
		} else {
			log.Println("Failed to get the PID")
		}

		if !success {
			log.Println("Couldn't kill the process")
		}
	}
}

// Close closes the current debugger
func (r *Debugger) Close() {
	close(r.start)
	r.kill(r.cmd)
	close(r.done)
}

// Wait waits until the next msg/changes
func (r *Debugger) Wait() {
	<-r.done
}

func getPidByName(name string) (string, error) {
	pids, err := exec.Command("ps", "ax").Output()

	if err != nil {
		return "", err
	}

	lines := strings.Split(string(pids), "\n")

	for _, item := range lines {
		if strings.Contains(item, name) {
			pid := strings.Split(strings.TrimLeft(item, " "), " ")[0]

			return pid, nil
		}
	}

	return "", errors.New("Process not found")
}
