package ddcp

import (
	"flag"
	"log"
)

const (
	DEFAULT_CHUNK_SIZE = 1024 * 1024 * 100
)

type DdcpParams struct {
	source     string
	dest       string
	preserve   bool
	chunk_size int64
}

func ParseFlag() (params *DdcpParams) {
	params = &DdcpParams{}

	flag.StringVar(&params.source, "s", "", "source")
	flag.StringVar(&params.dest, "d", "", "dest")
	flag.BoolVar(&params.preserve, "p", false, "preserve")
	flag.Int64Var(&params.chunk_size, "n", DEFAULT_CHUNK_SIZE, "chunk size")
	flag.Parse()

	if params.source == "" {
		log.Fatal("'-s' is required")
	}

	if params.dest == "" {
		log.Fatal("'-d' is required")
	}

	return
}
