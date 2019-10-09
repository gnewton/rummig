package main

import (
	//"fmt"

	//"io/ioutil"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dbi, err := dbInit("dbfile")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Println("Closing db")
		if dbi.tx != nil {
			err := commit(dbi.tx)
			if err != nil {
				log.Fatal(err)
			}
		}
		if dbi.db != nil {
			err := dbi.db.Close()
			if err != nil {
				log.Fatal(err)
			}
		}

	}()

	err = Walk("/home/gnewton/",
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
					log.Println(path)
					h, err := hash(path)
					if err != nil {
						log.Println(err)
						return err
					}
					log.Printf("%s  %x", path, h)
					if h != nil {
						err = store(dbi.hashes, h)
						if err != nil {
							log.Println(err)
							return err
						}
					}
				}
				return nil
			}
			return nil
		})
	//store()

}
