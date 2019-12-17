package watcher

// AppConfig contains the argument information for delve, watcher and main package
type AppConfig struct {
	DelveArgs []string `yaml:"delve-args,omitempty"`
	Package   struct {
		Path string   `yaml:"path"`
		Args []string `yaml:"args,omitempty"`
	} `yaml:"package"`
	Watcher struct {
		Run         string `yaml:"run"`
		Watch       string `yaml:"watch"`
		WatchVendor string `yaml:"watch-vendor"`
	} `yaml:"watcher,omitempty"`
}

// VersionInfo contains all the information regarding the version
type VersionInfo struct {
	GoVersion string
	Version   string
	Commit    string
	BuildTime string
}
