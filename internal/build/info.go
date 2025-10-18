package build

import "fmt"

// AppBuildInfo is a base application build info
type AppBuildInfo struct {
	Version string `json:"version"`
	Date    string `json:"date"`
	Commit  string `json:"commit"`
}

// NewBuildInfo create new build info instance
func NewBuildInfo(version, date, commit string) *AppBuildInfo {
	var info AppBuildInfo

	if version != "" {
		info.Version = version
	}
	if date != "" {
		info.Date = date
	}
	if commit != "" {
		info.Commit = commit
	}

	return &info
}

// Print prints application build info to stdout
func (b *AppBuildInfo) Print() {
	fmt.Printf(
		"Build version: %s\n"+
			"Build date: %s\n"+
			"Build commit: %s\n",
		b.Version,
		b.Date,
		b.Commit)
}
