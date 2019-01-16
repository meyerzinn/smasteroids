package smasteroids

var version string

func Version() string {
	if version == "" {
		return "unknown"
	}
	return version
}
