# fdbexplorer

[![license](https://img.shields.io/github/license/pwood/fdbexplorer.svg)](https://github.com/pwood/fdbexplorer/blob/master/LICENSE)
[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg)](https://github.com/RichardLitt/standard-readme)

> A TUI tool for exploring the status of FoundationDB clusters.

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Background

The output of `status details` from fdbcli is incomplete, and the JSON file produced by `status json` is difficult
to consume at a glance. `fdbexplorer` is a took in the same vane as `k9s` to make viewing an FoundationDB cluster
easier for operators at a glance.

**This tool is very early in development, it still needs more refactoring to make it maintainable - plus more functionality needs to be added.**

## Install

Download the latest release for your operating system and architecture from: [https://github.com/pwood/fdbexplorer/releases](https://github.com/pwood/fdbexplorer/releases)

Alternatively you may build yourself on a machine that has the `foundationdb-clients` installed by:

`go build`

## Usage

`fdbexplorer` is primarily configured by command line parameters. If `FDB_CLUSTER_FILE` environment variable is set, 
then it will attempt to find the cluster file at that location automatically (i.e. no command line arguments are required).

```shell
# ./fdbexplorer --help
Usage of ./fdbexplorer:
  -cluster-file string
    	Location of FoundationDB cluster file. (default "/etc/foundationdb/fdb.cluster")
  -input-file string
    	Location of an output of 'status json' to explore, will not connect to FoundationDB.
  -interval duration
    	Interval for polling FoundationDB for status. (default 10s)
```

## Maintainers

[@pwood](https://github.com/pwood)

## Contributing

Feel free to dive in! [Open an issue](https://github.com/pwood/fdbexplorer/issues/new) or submit PRs.

This project follows the [Contributor Covenant](https://www.contributor-covenant.org/version/1/4/code-of-conduct/) Code
of Conduct.

## License

Copyright 2022 Peter Wood & Contributors

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "
AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.