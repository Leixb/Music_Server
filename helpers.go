package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"syscall"
)

var r = regexp.MustCompile(".*(?P<type>album|artist|playlist|profile|track)/(?P<Code>[0-9]+)")

func getDirContent(path string, getfiles bool) []string {
	_files, err := ioutil.ReadDir(MusicDir + "/" + path)
	if err != nil {
		fmt.Println("Cannot read folder")
		return nil
	}

	var files []string

	for _, f := range _files {
		if f.IsDir() || getfiles {
			files = append(files, f.Name())
		}
	}

	return files
}

func addFullPath(path string, files []string) []string {
	for i := range files {
		files[i] = path + "/" + files[i]
	}
	return files
}

func getValues(url string) (string, string) {
	matches := r.FindStringSubmatch(url)

	if len(matches) == 3 {
		return matches[1], matches[2]
	} else {
		return "", ""
	}
}

func download(qType, qId string) error {
	link := "https://wwww.deezer.com/" + qType + "/" + qId

	cmd := exec.Command("/usr/bin/smloadr", "-q", "FLAC", "-p", MusicDir, "-u", link)

	log.Println("Downloading", link)

	err := cmd.Run()

	if err == nil {
		return nil
	}

	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			if status.ExitStatus() == 1 {
				return nil
			}
		}
	}
	return err
}
