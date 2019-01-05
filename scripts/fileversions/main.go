package main

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	f, err := os.Create("fileversions.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintln(f, "package main")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "var fileVersions = map[string]int{")

	filepath.Walk("public", func(p string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}

		bs, err := ioutil.ReadFile(p)
		if err != nil {
			panic(err)
		}

		h := fnv.New32a()
		h.Write(bs)
		sum := h.Sum32()

		fmt.Fprintf(f, "\t\"%s\": %d,\n", p, sum)

		return nil
	})

	fmt.Fprintln(f, "}")
}
