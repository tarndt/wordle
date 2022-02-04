package conf

import (
	"flag"
	"os"
)

const (
	DefDictPath    = "/usr/share/dict/words"
	DefLetterCount = 5
)

func MustGetArgs() (dictPath string, letterc uint) {
	var help bool
	flag.UintVar(&letterc, "letters", DefLetterCount, "Number of letters in word")
	flag.StringVar(&dictPath, "dict-file", DefDictPath, "Path of dictionary file")
	flag.BoolVar(&help, "help", false, "Show usage and exit")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	return
}
