package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sinloss/shipper/shipper"
)

var (
	t *string
	p *string
	v *string
)

func init() {
	t = flag.String("t", "", "Specify the build tags for the generated go file")
	p = flag.String("p", "main", "Specify the package name for the generated go file")
	v = flag.String("v", "A", "Specify the variable name of map containing all the embeded files")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <dir> <dest-file> [includes-without-gzip] [-- <includes-with-gzip>] \n",
			filepath.Base(os.Args[0]))
		fmt.Println("  Includes are comma seperated file paths in `dir`, wildcards are supportted. If there are")
		fmt.Println("  none comma seperated file paths given, all the files in `dir` will be included without gzip")
		fmt.Println("  by default")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
}

func parse() (shipper.Meta, string) {
	flag.Parse()
	positional := flag.Args()

	l := len(positional)
	if l < 2 {
		log.Fatalf("expecting at least 2 arguments yet got %d", l)
	}

	meta := shipper.Meta{Tags: *t, Package: *p, VarName: *v}
	meta.Dir = positional[0]
	destfile := positional[1]

	var includes, gzincludes []string
	if l >= 3 {
		gzindex := 2
		if positional[gzindex] != "--" {
			includes = strings.Split(positional[2], ",")
			gzindex++
		}

		if l >= gzindex+1 {
			positional = positional[gzindex:]
			l = len(positional)
			if l != 2 {
				log.Fatal("the includes-with-gzip should be like a `--` followed by comma seperated" +
					" file paths, yet got " + fmt.Sprint(positional))
			}
			gzincludes = strings.Split(positional[1], ",")
		}
	}

	for _, include := range includes {
		err := meta.Including(include, false)
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, include := range gzincludes {
		err := meta.Including(include, true)
		if err != nil {
			log.Fatal(err)
		}
	}
	// if no include pattern is given use the default pattern
	if len(meta.Includes) == 0 {
		err := meta.Including("*", false)
		if err != nil {
			log.Fatal(err)
		}
	}

	return meta, destfile
}

func main() {
	meta, dest := parse()
	err := shipper.Ship(meta, dest)
	if err != nil {
		log.Fatal(err)
	}
}
