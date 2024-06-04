package main

import (
	"os"

	"github.com/sleepdeprecation/tfmgr/internal/downloader"
)

func main() {
	dl := downloader.New()
	release, err := dl.GetRelease("latest")
	if err != nil {
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = dl.Download(release, cwd)
	if err != nil {
		panic(err)
	}

}
