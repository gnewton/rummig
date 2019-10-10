package main

import (
	//	"crypto/sha256"
	"crypto/md5"
	"errors"
	"io"
	"log"
	"os"
)

func hash(filename string) ([]byte, error) {
	if len(filename) == 0 {
		return nil, errors.New("Filename is empty")
	}

	f, err := os.Open(filename)

	if err != nil {
		log.Println(err)
		m, ok := err.(*os.PathError)
		if ok {
			e := m.Err.Error()
			if e == "permission denied" || e == "no data available" {
				return nil, nil
			}
		}
		return nil, err
	}
	defer f.Close()

	//h := sha256.New()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		m, ok := err.(*os.PathError)
		if ok {
			if m.Err.Error() == "no data available" {
				return nil, nil
			}
		}
		log.Println(err)
		return nil, err
	}

	return h.Sum(nil), nil
}
