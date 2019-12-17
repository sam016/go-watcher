package watcher

// AppConfig contains the argument information for delve, watcher and main package
type AppConfig struct {
	DelveArgs   map[string]string `yaml:delve`
	PackageArgs map[string]string `yaml:package`
	WatcherArgs map[string]string `yaml:watcher`
}

// VersionInfo contains all the information regarding the version
type VersionInfo struct {
	GoVersion string
	Version   string
	Commit    string
	BuildTime string
}
