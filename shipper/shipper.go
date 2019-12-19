package shipper

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/sinloss/shipper/wildcard"
)

// Include carries the useful data of the included files
type Include struct {
	Filename string
	Wc       *wildcard.FA
	Gziped   bool
}

// Meta carries the metadata for templates and shipping process
type Meta struct {
	Tags     string
	Package  string
	VarName  string
	Dir      string    // ship from
	Includes []Include // including
}

// Shipped moulds the shipped go file's content
var shipped = template.New("shipped").Funcs(template.FuncMap{
	"cap": func(s string) (string, error) {
		runes := []rune(s)
		cap := runes[0]

		if unicode.IsLetter(cap) {
			runes[0] = unicode.ToUpper(cap)
			return string(runes), nil
		}
		return "", errors.New("the given string must start with a letter")
	},
})

// Fore moulds the fore part of the shipped go file
var fore = template.Must(shipped.New("fore").Parse(`// Code generated by shipper; DO NOT EDIT.

{{with .Tags}}// +build {{.}}

{{end}}package {{.Package}}

import (
	"github.com/sinloss/shipper/shipper"
)

// {{cap .VarName}} is the Asset
var {{cap .VarName}} = &shipper.Assets{
`))

// EntryStart moulds the start part of an asset entry
var entryStart = template.Must(
	shipped.New("entryStart").Parse(`
	"{{.Filename}}": shipper.Content{
		Gziped: {{.Gziped}},
		Bytes:  []byte("`))

// EntryStart moulds the end part of an asset entry
var entryEnd = template.Must(
	shipped.New("entryEnd").Parse(`"),
	},`))

// Aft moulds the aft part of the shipped go file
var aft = template.Must(shipped.New("Aft").Parse(`
}`))

func traverse(root string, dir string, callback func(string, string, string)) error {
	d, err := ioutil.ReadDir(path.Join(root, dir))
	if err != nil {
		return err
	}
	for _, fi := range d {
		if fi.IsDir() {
			traverse(root, path.Join(dir, fi.Name()), callback)
		} else {
			callback(root, dir, fi.Name())
		}
	}
	return nil
}

// Including adds a suit of include to the includes array
func (meta *Meta) Including(filename string, gziped bool) error {
	if filename == "" {
		return nil
	}

	fa, err := wildcard.Compile([]rune(
		strings.ReplaceAll(
			filepath.Join(meta.Dir, filename), "\\", "\\\\")))
	if err != nil {
		return err
	}

	meta.Includes = append(meta.Includes,
		Include{Filename: filename, Wc: fa, Gziped: gziped})
	return nil
}

// Ship ships the given set of files to a destfile
func Ship(meta Meta, destfile string) error {
	// check meta validity
	if len(meta.Includes) == 0 {
		return errors.New("not a single file is included")
	}
	if meta.Package == "" {
		return errors.New("empty package")
	}
	if meta.VarName == "" {
		return errors.New("empty variable name")
	}
	if filepath.Ext(destfile) != ".go" {
		return errors.New("destfile should be a go file")
	}

	// check dir validity
	stat, err := os.Stat(meta.Dir)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return errors.New("invalid directory")
	}

	// check dest file's directory
	if err := ckdir(filepath.Dir(destfile)); err != nil {
		return err
	}
	// create output file
	dest, err := os.Create(destfile)
	if err != nil {
		return err
	}
	defer dest.Close()
	wo := &w{dest}

	fore.Execute(dest, meta)

	buf := make([]byte, 1048576)
	traverse(meta.Dir, "", func(root string, dir string, filename string) {
		for _, include := range meta.Includes {
			// check file path
			fullpath := filepath.Join(root, dir, filename)
			if !include.Wc.Search([]rune(fullpath), true).AllMatching() {
				continue
			}
			// open file
			f, err := os.Open(fullpath)
			if err != nil {
				continue
			}
			defer f.Close()
			// write entry
			for {
				if n, err := f.Read(buf); err == io.EOF {
					break
				} else if err != nil {
					break
				} else {
					entryStart.Execute(dest, Include{Filename: path.Join(dir, filename), Gziped: include.Gziped})
					if include.Gziped {
						Gzip(wo, buf[:n])
					} else {
						wo.Write(buf[:n])
					}
					entryEnd.Execute(dest, nil)
				}
			}
		}
	})

	aft.Execute(dest, nil)
	return nil
}
