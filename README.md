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

This tool is early in its life, suggestions or constructive criticism are welcome.

## Screencast

<img src="https://raw.githubusercontent.com/pwood/fdbexplorer/main/demo.gif" width="800px" alt="Demonstration of fdbexplorer" title="Demonstration of fdbexplorer" />

## Install

Download the latest release for your operating system and architecture from: [https://github.com/pwood/fdbexplorer/releases](https://github.com/pwood/fdbexplorer/releases)

Build are available for:

* `linux/amd64`
* `darwin/amd64`

`darwin/arm64` (Mac M1/M2) is waiting for GitHub to make Mac M1 runners available, Apple have now released an `arm64`
version of FoundationDB.

Alternatively you may build yourself on a machine that has the `foundationdb-clients` installed by:

`go build`

## Usage

`fdbexplorer` is primarily configured by command line parameters, with no command line parameters by default it will 
attempt to connect to FoundationDB, see below.

```shell
# ./fdbexplorer --help
fdbexplorer (devel) (rev-3df6e48-dirty)

Usage of ./fdbexplorer:
  -cluster-file string
    	Location of FoundationDB cluster file, environment variable FDB_CLUSTER_FILE also obeyed. (default "/etc/foundationdb/fdb.cluster")
  -http-address string
    	Host and port number for http server to listen on, using 0.0.0.0 for all interface bind. (default "127.0.0.1:8080")
  -http-enable status json
    	If the http output should be enabled, making the status json output available on /status/json.
  -input-file string
    	Location of an output of 'status json' to explore, will not connect to FoundationDB.
  -url string
    	URL to fetch status json from periodically.
```

### Connect to FoundationDB

If you are using a build that was compiled with the FoundationDB shared library (linux/amd64) `fdbexplorer` will be able
to connect to your FoundationDB cluster directly.

`fdbexplorer` will use the following search path to find a valid cluster file:
 * `-cluster-file` command line argument
 * `FDB_CLUSTER_FILE` environment variable
 * `/etc/foundationdb/fdb.cluster`

### Read a copy of `status json`

If your FDB cluster is remote or isolated, you may capture a content of `status json` (e.g. `fdbcli --exec 'status json' > status.json`) 
and have `fdbexplorer` read that.

> `fdbexplorer -input-file status.json`

### Provide/Read from a HTTP endpoint

For convenience `fdbexplorer` will also act as a **simple unauthenticated** HTTP server sharing out the status json.

> `fdbexplorer -http-enable -http-address 0.0.0.0:8888`

It can then be read by using the following:

> `fdbexplorer -url http://<internal ip>:8888/status/json`

You do not have to use `fdbexplorer` to publish the contents of `status json`, however the endpoint you provided must
return a `200` and a `Content-Type` of `application/json`.

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