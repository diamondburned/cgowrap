package cgowrap

import (
	"log"
	"os"
	"path/filepath"
)

var rootWorkDir string

// WorkDir returns a working directory for cgowrap.
func WorkDir(tail ...string) string {
	return initWork(true, tail...)
}

// WorkFile returns a working file for cgowrap.
func WorkFile(tail ...string) string {
	return initWork(false, tail...)
}

func initWork(isDir bool, tail ...string) string {
	if rootWorkDir == "" {
		tmp, err := os.UserCacheDir()
		if err != nil {
			tmp = os.TempDir()
		}
		rootWorkDir = filepath.Join(tmp, "cgowrap")
	}

	dir := rootWorkDir
	if len(tail) > 0 {
		dir = filepath.Join(dir, filepath.Join(tail...))
	}

	mkdir := dir
	if !isDir {
		mkdir = filepath.Dir(dir)
	}
	if err := os.MkdirAll(mkdir, os.ModePerm); err != nil {
		log.Fatalln("cannot make temp working dir:", err)
	}

	return dir
}
