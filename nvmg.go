//    Copyright 2018 Yoshi Yamaguchi
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/blang/semver"
	"github.com/mholt/archiver"
	pe "github.com/pkg/errors"
)

const (
	// Version follows semver.
	Version = "0.1.0"

	// NodeDistributionURL is the URL of Node.js distribution list.
	NodeDistributionURL = "https://nodejs.org/dist/"

	// NodeIndexURL is the URL of index.json
	NodeIndexURL = "https://nodejs.org/dist/index.json"
)

type NVMGError struct {
	ErrorString string
}

func (ne *NVMGError) Error() string {
	return ne.ErrorString
}

// NVMG is the struct to express the command `nvmg`
type NVMG struct {
	// ioout and ioerr are the targets for stdout and stderr caused by nvmg command.
	ioout, ioerr io.Writer
	// args should hold os.Args[1:]
	args []string
	// mainFlags is the root FlagSet of nvmg command.
	mainFlags   *flag.FlagSet
	versionFlag *bool
	helpFlag    *bool
	Home        string
}

// NewNVMG returns a new instance of NVMG with the initialization of parsing arguments.
func NewNVMG(args []string, home string) (*NVMG, error) {
	nvmg := &NVMG{
		args:  args,
		ioout: os.Stdout,
		ioerr: os.Stderr,
		Home:  home,
	}

	flags := flag.NewFlagSet("nvmgFlags", flag.ExitOnError)
	flags.SetOutput(nvmg.ioerr)
	flags.Usage = func() {
		fmt.Fprintf(nvmg.ioerr, "%v\n", helpMessage)
	}

	nvmg.versionFlag = flags.Bool("version", false, "Print out the latest released version of nvmg.")
	nvmg.helpFlag = flags.Bool("help", false, "Show this message.")

	if err := flags.Parse(args[1:]); err != nil {
		return nil, pe.Wrap(err, fmt.Sprintf("Could not parse the argument: %v", args[1:]))
	}
	nvmg.mainFlags = flags
	return nvmg, nil
}

func (n *NVMG) printfOut(s string) (int, error) {
	return fmt.Fprintf(n.ioout, "%v\n", s)
}

// Run executes the command.
func (n *NVMG) Run() error {
	if n.mainFlags == nil {
		return fmt.Errorf("nvmg instance is not initialized")
	}

	switch {
	case *n.versionFlag:
		n.printVersion()
		return nil
	case *n.helpFlag:
		n.printHelp()
		return nil
	}

	subcommand := n.mainFlags.Arg(0)
	if subcommand == "" {
		n.printHelp()
		return nil
	}

	switch subcommand {
	case "install":
		return n.RunInstall()
	case "uninstall", "remove", "delete":
		n.printfOut("uninstall") // TODO: replace here to actual command.
	case "use":
	case "exec":
	case "run":
	case "current":
	case "ls":
	case "ls-remote":
	case "version":
	case "version-remote":
	case "deactivate":
	case "alias":
	case "unalias":
	case "reinstall-packages":
	case "unload":
	case "which":
	case "help":
	default:
	}

	return nil
}

// printVersion outputs the version of nvmg itself.
func (n *NVMG) printVersion() {
	fmt.Printf("nvmg version: %v\n", Version)
}

// printHelp just print the usage of nvmg.
func (n *NVMG) printHelp() {
	fmt.Printf("%v\n", helpMessage)
	os.Exit(0)
}

// RunInstall parses the arguments for `install` subcommand and runs it accordingly.
func (n *NVMG) RunInstall() error {
	if len(n.args) < 2 {
		return fmt.Errorf("not enough arguments: %v", n.args)
	}
	flags := flag.NewFlagSet("installFlags", flag.ExitOnError)
	ltsFlag := flags.Bool("lts", false, "Refer to the Long-term support version for aliases.")
	flags.Parse(n.args[1:])
	if flags.NArg() < 1 {
		return fmt.Errorf("not enough arguments for install: %v", flags.Args())
	}
	_ = ltsFlag // TODO: implement LTS context switch.
	ver, err := n.expandVersionNumber(flags.Arg(1))
	if err != nil {
		return err
	}
	return n.Install(ver)
}

// expandVersionNumber checks if the version number is valid and return
func (n *NVMG) expandVersionNumber(ver string) (string, error) {
	if strings.HasPrefix(ver, "v") {
		ver = ver[1:]
	}
	if ver == "stable" {
		// TODO: implement here.
	}
	v, err := semver.Parse(ver)
	if err != nil {
		return "", pe.Wrapf(err, "invalid version number: %v", ver)
	}
	return fmt.Sprintf("v%v", v.String()), nil
}

// Install fetch pre-build binary from the distribution and expand the compressed file in temp dir
// and move the directory into configured directory.
func (n *NVMG) Install(ver string) error {
	filename := nodeBinaryArchiveName(ver)
	dirname := ver
	u, err := url.Parse(NodeDistributionURL)
	if err != nil {
		return err
	}
	p, err := url.Parse(path.Join("./", dirname, filename))
	if err != nil {
		return err
	}
	target := u.ResolveReference(p)
	n.printfOut(target.String())
	resp, err := http.Get(target.String())
	if err != nil {
		return err
	}
	downloaded := path.Join(os.TempDir(), filename)
	file, err := os.Create(downloaded)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return pe.Wrapf(err, "failed to download file: %v", filename)
	}
	destDir := path.Join(n.Home, ver)
	return expandArchiveFile(downloaded, destDir)
}

func expandArchiveFile(filename, dest string) error {
	var a archiver.Archiver
	switch {
	case strings.HasSuffix(filename, ".tar.gz"), strings.HasSuffix(filename, ".tgz"):
		a = archiver.TarGz
	case strings.HasSuffix(filename, ".tar.xz"), strings.HasSuffix(filename, ".txz"):
		a = archiver.TarXZ
	case strings.HasSuffix(filename, ".tar.bz2"), strings.HasSuffix(filename, ".tbz"):
		a = archiver.TarBz2
	case strings.HasSuffix(filename, ".zip"):
		a = archiver.Zip
	}

	tempDir, err := ioutil.TempDir("", "nvmg")
	defer os.RemoveAll(tempDir)
	if err != nil {
		return pe.Wrap(err, "couldn't craete tempdir")
	}
	if err := a.Open(filename, tempDir); err != nil {
		return pe.Wrapf(err, "couldn't open archive file into temp directory: %v", filename)
	}
	files, err := ioutil.ReadDir(tempDir)
	if err != nil {
		return pe.Wrapf(err, "couldn't read temp directory for expand: %v", tempDir)
	}
	if len(files) == 1 && files[0].IsDir() {
		tempDir = path.Join(tempDir, files[0].Name())
	}
	targetFiles, err := ioutil.ReadDir(tempDir)
	if err != nil {
		return pe.Cause(err)
	}
	if err := os.MkdirAll(dest, os.FileMode(0755)); err != nil {
		return pe.Cause(err)
	}
	for _, f := range targetFiles {
		s := path.Join(tempDir, f.Name())
		t := path.Join(dest, f.Name())
		if err := os.Rename(s, t); err != nil {
			return err
		}
	}
	return nil
}

// nodeBinaryArchinveName generates the filename of archive file uploaded on the distribution page.
// The CPU architecture name and OS platform name are listed here:
//    https://go.googlesource.com/go/+/master/src/go/build/syslist.go
func nodeBinaryArchiveName(ver string) string {
	var platform, arch, ext string
	switch runtime.GOOS {
	case "linux":
		platform = "linux"
		ext = "tar.gz"
	case "darwin":
		platform = "darwin"
		ext = "tar.gz"
	case "windows":
		platform = "win"
		ext = "zip"
	case "solaris":
		platform = "sunos"
		ext = "tar.gz"
	default:
		platform = "linux"
		ext = "tar.gz"
	}

	// TODO: there's no easy way to get ARM version from runtime, so it requires some way to
	// embed build target ARM version. This should be achieved in the same method as runtime.GOOS.
	// (ref. https://go.googlesource.com/go/+/master/src/go/build/syslist.go)
	switch runtime.GOARCH {
	case "386":
		arch = "x86"
	case "amd64":
		arch = "x64"
	case "arm64":
		arch = "arm64"
	case "ppc64":
		arch = "ppc64"
	case "ppc64le":
		arch = "ppc64le"
	case "s390x":
		arch = "s390x"
	default:
		arch = "x64"
	}

	return fmt.Sprintf("node-%v-%v-%v.%v", ver, platform, arch, ext)
}
