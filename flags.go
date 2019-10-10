package main

import (
	flag "github.com/spf13/pflag"
	"strings"
)

type Mode int

const (
	Index Mode = iota
	Show
	Compare
)

var TxSize *int = flag.IntP("txsize", "t", 1000, "Transaction size")
var DBFile *string = flag.StringP("dbfile", "f", "rummig.db", "Database filename")
var RootDir *string = flag.StringP("dir", "d", ".", "Root dir")

var mode *int = flag.IntP("mode", "m", int(Index), "Mode")

var ignoreDirsFlag *string = flag.StringP("ignore", "i", "", "List of dirs to ignore, comma separated")
var ignoreDirs []string

func initFlags() {
	flag.Parse()

	ignoreDirs = strings.Split(*ignoreDirsFlag, ",")
}
