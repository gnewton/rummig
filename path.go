// This file is directly copied from: https://golang.org/src/path/filepath/path.go
// so its license and copyright apply to this file
// The only difference is that readDirNames() does not return a sorted directly

package main

import (
	//"errors"
	"os"
	"path/filepath"
)

func Walk(root string, walkFn filepath.WalkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		err = walkFn(root, nil, err)
	} else {
		err = walk(root, info, walkFn)
	}
	if err == filepath.SkipDir {
		return nil
	}
	return err
}

func walk(path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	// NEW
	// filter by FileMode (default reject ModeDevice, ModeNamedPipe, ModeSocket, ModeIrregular
	if !info.IsDir() {
		return walkFn(path, info, nil)
	}

	names, err := readDirNames(path)

	err1 := walkFn(path, info, err)

	// If err != nil, walk can't walk into this directory.

	// err1 != nil means walkFn want walk to skip this directory or stop walking.

	// Therefore, if one of err and err1 isn't nil, walk will return.

	if err != nil || err1 != nil {
		// The caller's behavior is controlled by the return value, which is decided
		// by walkFn. walkFn may ignore err and return nil.
		// If walkFn returns SkipDir, it will be handled by the caller.
		// So walk should return whatever walkFn returns.
		return err1
	}

	for _, name := range names {
		filename := filepath.Join(path, name)
		fileInfo, err := os.Lstat(filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir {
				return err
			}
		} else {
			err = walk(filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != filepath.SkipDir {
					return err
				}
			}
		}
	}
	return nil
}

func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)

	if err != nil {
		return nil, err
	}

	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	return names, nil
}
