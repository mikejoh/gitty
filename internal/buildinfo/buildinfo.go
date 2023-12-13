package buildinfo

import "fmt"

var (
	Version = "0.0.0"
	Name    string
)

type BuildInfo struct {
	Name    string
	Version string
}

func (bi BuildInfo) String() string {
	return fmt.Sprintf("%s version: %s", bi.Name, bi.Version)
}

func Get() BuildInfo {
	return BuildInfo{
		Version: Version,
		Name:    Name,
	}
}
