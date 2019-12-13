package shipper

import (
	"errors"
	"os"
	"path/filepath"
)

// Content represents the file's content
type Content struct {
	Gziped bool
	Bytes  []byte
}

// Assets maps a file's name to its content
type Assets map[string]Content

// Restore restores the underlying contents to the current working directory
// with its original name
func (as Assets) Restore(names ...string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	for _, name := range names {
		err := as.RestoreAs(name, filepath.Join(wd, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// RestoreAs restores the underlying contents to the given dest path
func (as Assets) RestoreAs(name string, dest string) error {
	content := as[name]
	if content.Bytes != nil {
		// check directory
		dir := filepath.Dir(dest)
		if stat, err := os.Stat(dir); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			os.MkdirAll(dir, 0666)
		} else if !stat.IsDir() {
			return errors.New("a same name non-folder file exists")
		}

		f, err := os.Create(dest)
		if err != nil {
			return err
		}

		data := content.Bytes
		if content.Gziped {
			data, err = UnGzip(content.Bytes)
			if err != nil {
				return err
			}
		}
		_, err = f.Write(data)
		return err
	}
	return errors.New("could not find contents mapped to the given filename " + name)
}
