package watcher

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// Binary name used for built package
const binaryName = "watcher"

func (appConfig *AppConfig) packagePath() string {
	run := appConfig.Watcher.Run
	if run != "" {
		return run
	}

	return "."
}

// generateBinaryName generates a new binary name for each rebuild, for preventing any sorts of conflicts
func (appConfig *AppConfig) generateBinaryName() string {
	rand.Seed(time.Now().UnixNano())
	randName := rand.Int31n(999999)
	packageName := strings.Replace(appConfig.packagePath(), "/", "-", -1)

	return fmt.Sprintf("%s-%s-%d", generateBinaryPrefix(), packageName, randName)
}

func generateBinaryPrefix() string {
	path := os.Getenv("GOPATH")
	if path != "" {
		return fmt.Sprintf("%s/bin/%s", path, binaryName)
	}

	return path
}

func (appConfig *AppConfig) clean() bool {

	if appConfig.Watcher.Run == "" {
		log.Fatalln("Watchers `run` arg not set")
	}

	if appConfig.Watcher.Watch == "" {
		log.Fatalln("Watchers `watch` arg not set")
	}

	return true
}

// runCommand runs the command with given name and arguments. It copies the
// logs to standard output
func runCommand(name string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return cmd, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return cmd, err
	}

	if err := cmd.Start(); err != nil {
		return cmd, err
	}

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	return cmd, nil
}

// ParseArgs extracts the application parameters from args and returns
// Params instance with separated watcher and application parameters
func ParseArgs(args []string, vinfo VersionInfo) *AppConfig {

	var appConfig AppConfig

	isConfigLoaded := false

	// remove the command argument
	args = args[1:len(args)]

	for i := 0; i < len(args); i++ {
		arg := stripDash(args[i])

		if arg == "version" {
			fmt.Println("go-version=", vinfo.GoVersion)
			fmt.Println("version=", vinfo.Version)
			fmt.Println("commit=", vinfo.Commit)
			fmt.Println("built-at=", vinfo.BuildTime)
			return nil
		}

		if arg == "f" {
			configFileName := args[i+1]

			yamlFile, err := ioutil.ReadFile(configFileName)
			if err != nil {
				fmt.Printf("Error reading YAML file: %s\n", err)
				return nil
			}

			err = yaml.Unmarshal(yamlFile, &appConfig)
			if err != nil {
				fmt.Printf("Error parsing YAML file: %s\n", err)
				return nil
			}

			fmt.Printf("%#v\n", appConfig)

			isConfigLoaded = true
		}
	}

	if !isConfigLoaded || !appConfig.clean() {
		return nil
	}

	return &appConfig
}

// stripDash removes the both single and double dash chars and returns
// the actual parameter name
func stripDash(arg string) string {
	if len(arg) > 1 {
		if arg[1] == '-' {
			return arg[2:]
		} else if arg[0] == '-' {
			return arg[1:]
		}
	}

	return arg
}

func existIn(search string, in []string) bool {
	for i := range in {
		if search == in[i] {
			return true
		}
	}

	return false
}

func removeFile(fileName string) {
	if fileName != "" {
		cmd := exec.Command("rm", fileName)
		cmd.Run()
		cmd.Wait()
	}
}
