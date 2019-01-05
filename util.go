package main

import "fmt"

func getVersion(name string) string {
	return fmt.Sprint(fileVersions[name])
}
