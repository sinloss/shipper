# Shipper ships binary files as a go source file

It is espically used for shipping dynamic libraries and executables which
needs to get restored to the local file system.

After intallation via `go get -u github.com/sinloss/shipper`, you could use
`shipper` to do the shipping.

The usage is as follows:
```
Usage: shipper [options] <dir> <dest-file> [includes-without-gzip] [-- <includes-with-gzip>]
  Includes are comma seperated file paths in `dir`, wildcards are supportted. If there are
  none comma seperated file paths given, all the files in `dir` will be included without gzip
  by default
Options:
  -p string
        Specify the package name for the generated go file (default "main")
  -t string
        Specify the build tags for the generated go file
  -v string
        Specify the variable name of map containing all the embeded files (default "A")
```


# It's simple yet works

Just by browsing these source files you'll find it very straight-forward
and simple. It just reads file contents and write them as a simple go 
source file.

# Gzip / UnGzip supported

Yes, supported.

# The generated file is simple

As simple as a map declaration. As all the functionality lies within the
*github.com/sinloss/shipper* package, it's easy to extend it's feature.