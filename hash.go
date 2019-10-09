package main

import (
	"crypto/sha256"
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
			if m.Err.Error() == "permission denied" {
				return nil, nil
			}
		}
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Println(err)
		return nil, err
	}

	return h.Sum(nil), nil
}
