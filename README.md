# Mock generator

It generates the mock implementing the interface, which can be used during development and testing

## Installation

Use [`go get`](https://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies) to install and update:

```sh
$ go get -u github.com/AntonParaskiv/mockGen/...
```

## Usage

From the commandline, `mockGen` can generate Mock-package for Interface-package. It creates a directory with a mock in the same parent directory with an interface.

```sh
$ mockGen INTERFACE_PACKAGE_PATH
```
