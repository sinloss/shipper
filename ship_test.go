package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sinloss/shipper/shipped"
)

func TestMain(m *testing.M) {
	os.Args = []string{"shipper", "-p", "shipped", "helloworld", "shipped/helloworld.go", "*o", "--", "*.bar"}
	main()
	m.Run()

}

func TestFacility(t *testing.T) {
	testRestore("hello", "bonjour", t)
	testRestore("world/foo.bar", "foooobaaar", t)
}

func testRestore(name string, dest string, t *testing.T) {
	shipped.A.RestoreAs(name, dest)
	defer os.Remove(dest)
	check(dest, filepath.Join("helloworld", name), t)
}

func check(target string, origin string, t *testing.T) {
	ori, err := ioutil.ReadFile(origin)
	if err != nil {
		t.Errorf("%s not found", origin)
		return
	}

	tar, err := ioutil.ReadFile(target)
	if err != nil {
		t.Errorf("%s not successfully restored", target)
	}

	if len(ori) != len(tar) {
		goto DIFFERENT
	}
	for i, b := range ori {
		if b != tar[i] {
			goto DIFFERENT
		}
	}
	return

DIFFERENT:
	t.Errorf("contents of %s should be identical with %s", target, origin)
}
