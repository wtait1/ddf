package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/labstack/gommon/color"
)

// Version to be injected at build time
var Version string

var opts struct {
	Dir     string `short:"d" long:"dir" default:"." description:"Directory to search for duplicate files"`
	Quiet   bool   `short:"q" long:"quiet" description:"Only show names of duplicate files (default when piping output)"`
	Version bool   `short:"v" long:"version" description:"version info"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	_, err := flags.Parse(&opts)
	check(err)

	if opts.Version {
		fmt.Printf("ddf %s\n", Version)
		return
	}

	items, err := ioutil.ReadDir(opts.Dir)
	check(err)

	fileMap := map[string][]os.FileInfo{}
	for _, i := range items {
		if !i.IsDir() {
			hash := getFileHash(i.Name())
			if fileList, ok := fileMap[hash]; ok {
				fileMap[hash] = append(fileList, i)
			} else {
				fileMap[hash] = []os.FileInfo{i}
			}
		}
	}

	repeats := 0
	for _, fileList := range fileMap {
		for ind, file := range fileList {
			if ind == 0 {
				if beQuiet() {
					continue
				}
				fmt.Println(color.Green(" " + file.Name()))
			} else {
				prefix := "> "
				if beQuiet() {
					prefix = ""
				}
				fmt.Println(color.Red(prefix + file.Name()))
				repeats++
			}
		}
	}
	if repeats > 0 && !beQuiet() {
		plural := ""
		if repeats > 1 {
			plural = "s"
		}
		fmt.Printf("\nFound %d repeat%s\n", repeats, plural)
	}
}

func getFileHash(file string) string {
	f, err := os.Open(file)
	check(err)
	defer f.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, f)
	check(err)

	return hex.EncodeToString(hash.Sum(nil))
}

func beQuiet() bool {
	return opts.Quiet || outputIsPiped()
}

func outputIsPiped() bool {
	fi, _ := os.Stdout.Stat()
	return fi.Mode()&os.ModeCharDevice == 0
}
