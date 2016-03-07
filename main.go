package main

import (
	"ddcp"
	"log"
	"os"
	"runtime"
)

func init() {
	log.SetFlags(0)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	max_procs := os.Getenv("GOMAXPROCS")

	if max_procs == "" {
		cpus := runtime.NumCPU()
		runtime.GOMAXPROCS(cpus)
	}

	params := ddcp.ParseFlag()
	err := ddcp.Ddcp(params)

	if err != nil {
		log.Fatal(err)
	}
}
