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

const helpMessage = `               Node Version Manager in Go

Note: <version> refers to any version-like string nvmg understands. This includes:
  - full or partial version numbers, starting with an optional "v" (0.10, v0.1.2, v1)
  - default (built-in) aliases: node, stable, unstable, iojs, system
  - custom aliases you define with 'nvmg alias foo'

Usage:
  nvmg help                                  Show this message
  nvmg --version                             Print out the latest released version of nvmg
  nvmg install [-s] <version>                Download and install a <version>, [-s] from source. Uses .nvmgrc if available
    --reinstall-packages-from=<version>     When installing, reinstall packages installed in <node|iojs|node version number>
  nvmg uninstall <version>                   Uninstall a version
  nvmg use [--silent] <version>              Modify PATH to use <version>. Uses .nvmgrc if available
  nvmg exec [--silent] <version> [<command>] Run <command> on <version>. Uses .nvmgrc if available
  nvmg run [--silent] <version> [<args>]     Run 'node' on <version> with <args> as arguments. Uses .nvmgrc if available
  nvmg current                               Display currently activated version
  nvmg ls                                    List installed versions
  nvmg ls <version>                          List versions matching a given description
  nvmg ls-remote                             List remote versions available for install
  nvmg version <version>                     Resolve the given description to a single local version
  nvmg version-remote <version>              Resolve the given description to a single remote version
  nvmg deactivate                            Undo effects of 'nvmg' on current shell
  nvmg alias [<pattern>]                     Show all aliases beginning with <pattern>
  nvmg alias <name> <version>                Set an alias named <name> pointing to <version>
  nvmg unalias <name>                        Deletes the alias named <name>
  nvmg reinstall-packages <version>          Reinstall global 'npm' packages contained in <version> to current version
  nvmg unload                                Unload 'nvmg' from shell
  nvmg which [<version>]                     Display path to installed node version. Uses .nvmgrc if available

Example:
  nvmg install v0.10.32                  Install a specific version number
  nvmg use 0.10                          Use the latest available 0.10.x release
  nvmg run 0.10.32 app.js                Run app.js using node v0.10.32
  nvmg exec 0.10.32 node app.js          Run 'node app.js' with the PATH pointing to node v0.10.32
  nvmg alias default 0.10.32             Set default node version on a shell

Note:
  to remove, delete, or uninstall nvmg - just remove the '$NVMG_DIR' folder (usually '~/.nvmg')
`
