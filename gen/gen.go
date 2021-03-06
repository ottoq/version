package main

////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"text/template"
)

////////////////////////////////////////////////////////////////////////////////

const (
	genFile = "version_gen.go"
)

////////////////////////////////////////////////////////////////////////////////

var (
	basePkg, runDir string

	// Absolute path to the file being generated
	genPath string

	// prop holds VCS populated information about the parent package
	prop *Properties

	// Version (and other) extraction
	vcsGit = &struct {
		name       string
		cmd        string   // binary to involve VCS
		dirtyCmd   []string // commands to check if a VCS repo is in a dirty state
		versionCmd []string // commands to return the build version identifier from the VCS's view of the application
	}{
		cmd:      "git",
		dirtyCmd: []string{"status", "-z", "--porcelain"},
		// Acceptable tag names must start with "v"
		versionCmd: []string{"describe", "--long", "--dirty", "--abbrev=10", "--tags", "--match=v", "--always"},
	}
)

////////////////////////////////////////////////////////////////////////////////

const versionTmpl = `
// Package version IS AUTO GENERATED by [version/gen/gen.go].
// This package is intended to be used as a method to inject build information
// during the compiling of a go project either from a developer's desk or from
// one of our CI tools. DO NOT HAND EDIT THIS FILE!
package version

var (
	// ID is the build id from our build pipeline (DEV if this is a local build).
	ID = "{{ .ID }}"

	// Description is a build description.
	Description = "{{ .Description }}"

	// Hostname is the machine hostname that ran the "go generate" step.
	Hostname = "{{ .Hostname }}"

	// Runtime is the go runtime version used in compilation.
	Runtime = "{{ .Runtime }}"
)

// String returns all version information
func String() string {
	s := "ID:          " + ID + "\n"
	s += "Description: " + Description + "\n"
	s += "Hostname:    " + Hostname + "\n"
	s += "Go Runtime:  " + Runtime + "\n"
	return s
}
`

////////////////////////////////////////////////////////////////////////////////

// Properties contain the variables generated
type Properties struct {
	ID          string // Builds unique, "DEV" if not a production build
	Description string // Description, commit hash and relevant info
	Hostname    string // Hostname of the machine building this version file
	Runtime     string // Runtime info, go version etc
}

// makeVersion returns a Properties instance based on the code repository.
func makeVersion() error {
	// describe repository by calling VCS command
	cmd := exec.Command(vcsGit.cmd, vcsGit.versionCmd...)

	// we must run the command inside the checked out project directory
	cmd.Dir = runDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	version := strings.TrimSpace(string(out))

	// get this machine's local hostname
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	prop = &Properties{
		ID:          "DEV",
		Description: version,
		Hostname:    hostname,
		Runtime:     runtime.Version(),
	}

	// inject Gitlab CI variables
	id, hasID := os.LookupEnv("CI_PIPELINE_ID")
	bid, hasBID := os.LookupEnv("CI_BUILD_ID")
	if hasID && hasBID {
		prop.ID = fmt.Sprintf("Pipeline %s (Build %s)", id, bid)
	}

	return nil
}

// generate figures out interesting things about the build platform, version
// and other such properties and returns the contents for a would-be `version.go'
// file.
func generate() (string, error) {
	log.Printf("generating %s version package\n", basePkg)

	t, err := template.New("").Parse(versionTmpl)
	if err != nil {
		return "", err
	}

	b := new(bytes.Buffer)
	if err := t.Execute(b, prop); err != nil {
		return "", err
	}

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		return string(b.Bytes()), err
	}

	return string(formatted), nil
}

////////////////////////////////////////////////////////////////////////////////

// verifyPath confirms generated code will be created under the correct path
// one directory above the gen directory
func verifyPath() {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		log.Fatal("failed to get determine current path")
	}

	goPath, ok := os.LookupEnv("GOPATH")
	if !ok {
		log.Fatalln("$GOPATH is not set in the os env")
	}

	runDir = path.Join(goPath, "src", basePkg)
	genPath = path.Join(path.Dir(filename), "../", genFile)
}

func write(data, path string) error {
	log.Printf("writing %s\n", genPath)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write([]byte(data)); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	var ok bool
	basePkg, ok = os.LookupEnv("GOPACKAGE")
	if !ok {
		log.Fatalln("$GOPACKAGE is not set in the os env")
	}

	err := makeVersion()
	if err != nil {
		log.Fatalln("issue determine version:\n" + err.Error())
	}

	verifyPath()

	gen, err := generate()
	if err != nil {
		log.Fatalf("generate error: %s\n%s\n", err.Error(), gen)
	}

	if err := write(gen, genPath); err != nil {
		log.Fatalf("write error: %s\n", err.Error())
	}
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	log.SetFlags(0)
	log.SetPrefix("VERSION::")
}

////////////////////////////////////////////////////////////////////////////////
