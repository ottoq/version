// Package version IS AUTO GENERATED by [version/gen/gen.go].
// This copy of the package is intended to be placeholder that is overwritten by go generate
package version

////////////////////////////////////////////////////////////////////////////////

var (
	// ID is the build id from our build pipeline (DEV if this is a local build).
	ID = "DEV"

	// Description is a build description.
	Description = "DEV"

	// Hostname is the machine hostname that ran the "go generate" step.
	Hostname = ""

	// Runtime is the go runtime version used in compilation.
	Runtime = ""
)

////////////////////////////////////////////////////////////////////////////////

// String returns all version information
func String() string {
	return "DEVELOPER BUILD\n"
}

////////////////////////////////////////////////////////////////////////////////
