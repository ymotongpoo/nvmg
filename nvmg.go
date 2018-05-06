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
)

const (
	// Version follows semver.
	Version = "0.1.0"

	// NodeDistributionURL is the URL of Node.js distribution list.
	NodeDistributionURL = "https://nodejs.org/dist/"
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
)

// NVMG is the struct to express the command `nvmg`
type NVMG struct {
	ioout, ioerr io.Writer
	flags        *flag.FlagSet
	versionFlag  *bool
	helpFlag     *bool
	subcommand   string
}

// NewNVMG returns a new instance of NVMG with the initialization of parsing arguments.
func NewNVMG(args []string) (*NVMG, ErrorStatus) {
	nvmg := &NVMG{
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
	if len(flags.Args()) > 0 {
		nvmg.subcommand = flags.Arg(0)
	}

	nvmg.flags = flags

	return nvmg, ExitStatusOK
}

func (n *NVMG) printfOut(s string) (int, error) {
	return fmt.Fprintf(n.ioout, "%v\n", s)
}

// Run executes the command.
func (n *NVMG) Run() ErrorStatus {
	if n.flags == nil {
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

	if n.subcommand == "" {
		n.printHelp()
		return ExitStatusOK
	}

	switch n.subcommand {
	case "install":
		n.Install()
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

// Install fetch pre-build binary from the distribution and expand the compressed file in temp dir
// and move the directory into configured directory.
func (n *NVMG) Install(ver string) ErrorStatus {
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
