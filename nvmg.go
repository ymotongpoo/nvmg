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
)

const (
	// Version follows semver.
	Version = "0.1.0"
)

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
		n.printfOut("install") // TODO: replace here to actual command.
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

func (n *NVMG) printHelp() {
	fmt.Printf("%v\n", helpMessage)
	os.Exit(0)
}
