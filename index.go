package main

import (
	"log"
	"os"
)

func index() error {
	dbi, err := dbInit(*DBFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {

		if dbi.tx != nil {
			err := commit(dbi.tx)
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Println("Closing db")
		if dbi.db != nil {
			err := dbi.db.Close()
			if err != nil {
				log.Fatal(err)
			}
		}

	}()
	totalCounter := 0
	counter := 0
	err = Walk(*RootDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				//log.Println(err)
				m, ok := err.(*os.PathError)
				if ok {
					//log.Println("z++++++", m.Op, ok)
					//log.Println("++++++", m.Path, ok)
					//log.Println("++++++", m.Err, ok)
					if m.Err.Error() == "permission denied" || m.Err.Error() == "no such file or directory" {
						// do not try and store
						return nil
					}
				}
				return err
				//log.Println("qwerty ", err)
				//return nil
			}
			mode := info.Mode()
			if mode.IsRegular() {
				filesize := info.Size()
				// No empty files
				if filesize > 0 {
					err := storeSize(dbi.filesizes, filesize)
					if err != nil {
						log.Println(err)
						return err
					}
					log.Println(counter, "/", totalCounter, " ", path)
					h, err := hash(path)
					if err != nil {
						log.Println(err)
						return err
					}
					log.Printf("%s  %x", path, h)
					if h == nil {
						return nil
					}
					err = store(dbi.hashes, h)
					if err != nil {
						log.Println(err)
						return err
					}
					counter++

					if counter > *TxSize {
						err = pushTx(dbi)
						if err != nil {
							return err
						}
						totalCounter += counter
						counter = 0
					}
				}
				return nil
			}
			return nil
		})
	//store()
	return nil
}
