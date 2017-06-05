package version

var (
	version string
	build   string
)

// Version provide the version number given at compile time. If not set,
// the version number is "0.0.0"
func Version() string {
	if version == "" {
		return "0.0.0"
	}
	return version
}

// Build provide a build number givent at compile time (git rev-parse HEAD).
// If not set, the build number is "000000"
func Build() string {
	if build == "" {
		return "000000"
	}
	return build
}
