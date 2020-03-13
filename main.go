package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/myfantasy/mffc/compress"
)

func main() {

	fdec := flag.String("fd", "", "decode file name")

	fenc := flag.String("fe", "", "encode file name")

	level := flag.Int("l", 9, "compress level")

	flag.Parse()

	c, err := compress.DeflateCompressorCreate(*level)
	if err != nil {
		panic(err)
	}

	if *fdec != "" {

		data, err := ioutil.ReadFile(*fdec)
		if err != nil {
			panic(err)
		}

		res, err := c.Restore(data)
		if err != nil {
			panic(err)
		}

		os.Stdout.Write(res)

	}

	if *fenc != "" {

		data, err := ioutil.ReadFile(*fenc)
		if err != nil {
			panic(err)
		}

		res, err := c.Compress(data)
		if err != nil {
			panic(err)
		}

		os.Stdout.Write(res)

	}
}
