# Developing

## FoundationDB Client Library

To compile `fdbexplorer` on your local machine you will need to install `libfdb_c` for your operating system and 
architecture, otherwise it will not compile. 

This will need to be 7.3.x at present, as the Go bindings used for development are 7.3. This is available as either the
raw library `.so/.dylib` or a `.pkg/.deb/.rpm` from [https://github.com/apple/foundationdb/releases/](https://github.com/apple/foundationdb/releases/).

`input/libfdb/fdb.go` uses a `//go:build` directive to only compile if CGO is enabled, and a known official distribution
for FoundationDB. Be aware at present there is no official distribution of `libfdb_c` for the linux `arm64` platform.