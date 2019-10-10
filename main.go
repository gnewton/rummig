package main

import (
	//"fmt"

	//"io/ioutil"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	initFlags()

	switch Mode(*mode) {
	case Index:
		index()
	case Compare:
	}

}
