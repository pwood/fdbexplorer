# Developing

## FoundationDB Client Library

To compile `fdbexplorer` on your local machine you will need to install `libfdb_c` for your operating system and 
architecture, otherwise it will not compile. This is available as either the `.so` or a `.pkg/.deb/.rpm` from 
[https://github.com/apple/foundationdb/releases/](https://github.com/apple/foundationdb/releases/).

`data_fdb.go` uses a `//go:build` directive to only compile if CGO is enabled, and your machine is linux or darwin on
`amd64`. This is because at present there is no official distribution of `libfdb_c` for any `arm64` platform.

## Developing on an Apple Mac M1/M2 (darwin/arm64)

You may develop without direct FoundationDB integration and rely on the `--input-file` functionality against a static 
`status json` dump.

Alternatively, to work around the above and permit development on Apple Mac M1/M2 processors a `build.Dockerfile` is
provided. Assuming you are running Docker, you can build and start this container - Docker will then use QEMU to emulate `amd64`. This 
environment is slow, but does function well enough to build and test.

```shell
# Build
docker build -f build.Dockerfile . -t golang-with-fdb
# Run
docker run -it -v ~/workspace:/workspace golang-with-fdb
```

*Take note to mount your workspace directory correctly.*