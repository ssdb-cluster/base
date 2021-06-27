package util

import (
	"os"
	"path/filepath"
)

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func FileSize(path string) int64 {
	st, err := os.Stat(path)
	if err != nil {
		return -1
	}
	return st.Size()
}

func Dirname(path string) string {
	return filepath.Dir(path)
}

func Basename(path string) string {
	return filepath.Base(path)
}

// panic if error occurs
func FsyncDir(dir string) {
	fp, err := os.OpenFile(dir, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	if err := fp.Sync(); err != nil {
		panic(err)
	}
	if err := fp.Close(); err != nil {
		panic(err)
	}
}
