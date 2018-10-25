# GoEmbedFS

GoEmbedFS is a very simple command line utility that generates Go language code
with embedded static files.

Features of generated files:

- compatibility with `http.FileSystem`
- helper functions for accessing and reading files
- embedded files deduplication
- gzip support with configurable minimal space savings
- go vet and golint clean


## Installation

Make sure that you have installed [Go](https://golang.org/) binary.

```sh
$ go get -u resenje.org/goembedfs/...
```

Executable `goembedfs` binary should be in your `$GOPATH/bin` directory.


## Usage

```
goembedfs [options...] package_name file...
```

Arguments `package_name` and `file` are required. You can provide multiple
arguments for `file` and if `file` is a directory, all files will be added
recursively.

Available options are:

- `-o` to specify output filename. By default the output will be printed to `STDOUT`.
- `-w` to change the working directory before the file is generated.
- `--tags` a comma-delimited list of build tags.

```
goembedfs mypackage index.html img.png robots.txt > mypackage.go
```

### All files from a directory

In this case all files from `assets` directory are embedded.

```
goembedfs mypackage assets > mypackage.go
```


### Output filename

When `STDOUT` redirect is not available, like in `//go:generate`, you can use
output filename option.

```
goembedfs -o mypackage.go mypackage index.html img.png robots.txt
```


### Working directory

It is common to keep all file assets in a specific directly of your project
and not to have its name in the generated file. That can be
accomplished with `-w` option.

```
goembedfs -o mypackage.go -w assets .
```

All file paths in `mypackage.go` will not contain `assets` prefix in path names.


### Build tags

To exclude generated files for some type of build, `--tags` option is available.
It can contain a comma delimited list of tags for each line of `// +build` directive.


### Gzip

Generated data may be gzip compressed with option `--gzip` and to specify minimal
space savings (reduction in size relative to the uncompressed size in percentage)
for data to be saved as gzip with option `--min-gzip-space-savings`.

Default minimal space savings is 5%.

```
goembedfs -o mypackage.go --gzip
```

To set minimal space savings to 30%:

```
goembedfs -o mypackage.go --gzip --min-gzip-space-savings 30
```

It is possible that compressed data is slightly longer then the uncompressed due to
gzip metadata. In this case data for that file will be saved in uncompressed form.


## Inspiration

There is a number of similar project already available and that are used as inspiration
and examples for building this simple alternative. Most notable of them are
https://github.com/jteeuwen/go-bindata and https://github.com/benbjohnson/genesis.
