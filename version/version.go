package version

// version vars are set in during build.

var (
	// Version contains the current version in SemVer format.
	Version string

	// Branch is the name of the branch referenced by HEAD.
	Branch string

	// Revision contains the hash of the latest commit on Branch.
	Revision string

	// BuildTime is the compiled build time.
	BuildTime string

	// GoVersion contains the the version of the go that performed the build.
	GoVersion string
)
