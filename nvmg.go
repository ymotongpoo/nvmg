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
	"os"
	"runtime"
	"strings"

	"github.com/blang/semver"
)

const (
	// Version follows semver.
	Version = "0.1.0"

	// NodeDistributionURL is the URL of Node.js distribution list.
	NodeDistributionURL = "https://nodejs.org/dist/"

	// NodeIndexURL is the URL of index.json
	NodeIndexURL = "https://nodejs.org/dist/index.json"
)

// ErrorStatus is the type to express the error status within nvmg command.
type ErrorStatus int

const (
	// ExitStatusOK is the status where parsing arguments went successful.
	ExitStatusOK ErrorStatus = iota

	// ExitStatusError is the status where parsing arguments went failure.
	ExitStatusError

	// ExitStatusNotInitialized is the status where NVMG is called without the initialization.
	ExitStatusNotInitialized

	// ExitStatusVersionNotFound is the status where invalid version number is specified.
	ExitStatusVersionNotFound
)

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
}

// NewNVMG returns a new instance of NVMG with the initialization of parsing arguments.
func NewNVMG(args []string) (*NVMG, ErrorStatus) {
	nvmg := &NVMG{
		args:  args,
		ioout: os.Stdout,
		ioerr: os.Stderr,
	}

	flags := flag.NewFlagSet("nvmgFlags", flag.ExitOnError)
	flags.SetOutput(nvmg.ioerr)
	flags.Usage = func() {
		fmt.Fprintf(nvmg.ioerr, "%v\n", helpMessage)
	}

	nvmg.versionFlag = flags.Bool("version", false, "Print out the latest released version of nvmg.")
	nvmg.helpFlag = flags.Bool("help", false, "Show this message.")

	if err := flags.Parse(args[1:]); err != nil {
		return nil, ExitStatusError
	}
	nvmg.mainFlags = flags
	return nvmg, ExitStatusOK
}

func (n *NVMG) printfOut(s string) (int, error) {
	return fmt.Fprintf(n.ioout, "%v\n", s)
}

// Run executes the command.
func (n *NVMG) Run() ErrorStatus {
	if n.mainFlags == nil {
		return ExitStatusNotInitialized
	}

	switch {
	case *n.versionFlag:
		n.printVersion()
		return ExitStatusOK
	case *n.helpFlag:
		n.printHelp()
		return ExitStatusOK
	}

	subcommand := n.mainFlags.Arg(0)
	if subcommand == "" {
		n.printHelp()
		return ExitStatusOK
	}

	switch subcommand {
	case "install":
		n.RunInstall()
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

	return ExitStatusOK
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
func (n *NVMG) RunInstall() ErrorStatus {
	if len(n.args) < 2 {
		return ExitStatusError
	}
	flags := flag.NewFlagSet("installFlags", flag.ExitOnError)
	ltsFlag := flags.Bool("lts", false, "Refer to the Long-term support version for aliases.")
	flags.Parse(n.args[1:])
	if flags.NArg() < 1 {
		return ExitStatusError
	}
	_ = ltsFlag // TODO: implement LTS context switch.
	ver, errStatus := n.expandVersionNumber(flags.Arg(0))
	if errStatus != ExitStatusOK {
		return errStatus
	}
	return n.Install(ver)
}

// expandVersionNumber checks if the version number is valid and return
func (n *NVMG) expandVersionNumber(ver string) (string, ErrorStatus) {
	if strings.HasPrefix(ver, "v") {
		ver = ver[1:]
	}
	if ver == "stable" {
		// TODO: implement here.
	}
	v, err := semver.Parse(ver)
	if err != nil {
		fmt.Fprintf(n.ioerr, "%v\n", err)
		return "", ExitStatusVersionNotFound
	}
	return fmt.Sprintf("v%v", v.String()), ExitStatusOK
}

// Install fetch pre-build binary from the distribution and expand the compressed file in temp dir
// and move the directory into configured directory.
func (n *NVMG) Install(ver string) ErrorStatus {
	filename := nodeBinaryArchiveName(ver)
	n.printfOut(filename)
	return ExitStatusOK
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
