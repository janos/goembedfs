// Copyright (c) 2018, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found s the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"resenje.org/goembedfs"
)

var (
	cli = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	output              = cli.String("o", "", "Output filename.")
	cwd                 = cli.String("w", "", "Change working directory.")
	tags                = cli.String("tags", "", "Comma-delimited list of build tags.")
	gzip                = cli.Bool("gzip", false, "Compress data with GZip.")
	minGzipSpaceSavings = cli.Float64("min-gzip-space-savings", 5, "Minimal reduction in size relative to the uncompressed size in percentage. Default 5.")
	strip               = cli.Int("strip", 0, "Remove the specified number of leading path elements.")
	help                = cli.Bool("h", false, "Show program usage.")
)

func main() {
	cli.Usage = func() {
		fmt.Fprintf(os.Stderr, `USAGE

%s [options...] package_name file...

Generates a single Go source file with package name package_name and
embedded files provied as file arguments. If a file is a directory,
all files from it will be added recursively.

OPTIONS
		`, os.Args[0])
		cli.PrintDefaults()
	}

	cli.Parse(os.Args[1:])

	if *help {
		cli.Usage()
		return
	}

	if cli.NArg() < 2 {
		cli.Usage()
		return
	}

	if *cwd != "" {
		handleError(os.Chdir(*cwd), "chdir")
	}

	args := cli.Args()

	var paths []string
	for _, arg := range args[1:] {
		a, err := expand(arg)
		handleError(err, "")
		paths = append(paths, a...)
	}

	var w io.Writer = os.Stdout
	if *output != "" {
		dir := filepath.Dir(*output)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0777)
			handleError(err, "create directories")
		}
		f, err := os.Create(*output)
		handleError(err, "output")
		defer f.Close()
		w = f
	}

	generator := goembedfs.New(w, args[0], goembedfs.WithTags(strings.Split(*tags, ",")...), goembedfs.WithGzip(*gzip), goembedfs.WithMinGzipSpaceSavings(*minGzipSpaceSavings))

	for _, path := range paths {
		fi, err := os.Stat(path)
		handleError(err, "stat")

		data, err := ioutil.ReadFile(path)
		handleError(err, "read file")

		path := filepath.Clean(filepath.ToSlash(path))
		if *strip > 0 {
			elements := strings.Split(path, string(filepath.Separator))
			if len(elements) > *strip {
				path = filepath.Join(elements[*strip:]...)
			}
		}
		err = generator.AddFile(
			path,
			data,
			fi.ModTime(),
		)
		handleError(err, "add")
	}

	err := generator.WriteFooter()
	handleError(err, "write footer")
}

// expand converts path into a list of all files within path.
// If path is a file then path is returned.
func expand(path string) ([]string, error) {
	if fi, err := os.Stat(path); err != nil {
		return nil, err
	} else if !fi.IsDir() {
		return []string{path}, nil
	}

	// Read files in directory.
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Iterate over files and expand.
	expanded := make([]string, 0, len(fis))
	for _, fi := range fis {
		a, err := expand(filepath.Join(path, fi.Name()))
		if err != nil {
			return nil, err
		}
		expanded = append(expanded, a...)
	}
	return expanded, nil
}

func handleError(err error, msg string) {
	if err == nil {
		return
	}
	if msg == "" {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Fprintf(os.Stderr, msg+": %v\n", err)
	}
	os.Exit(2)
}
