package cli

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func versionString() string {
	return Version + " (" + Commit + ", " + Date + ")"
}
