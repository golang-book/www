package main

import (
	"os"
	"strconv"
)

func getVersion(name string) string {
	fi, err := os.Stat(name)
	if err != nil {
		return ""
	}
	return strconv.FormatInt(fi.ModTime().Unix(), 36)
}
