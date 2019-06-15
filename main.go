package main

import (
	"flag"

	"github.com/awgh/marvin/marvin"
)

func main() {

	var configDir, chainFile, mcflyFile string
	flag.StringVar(&configDir, "confdir", "config", "Config Directory")
	flag.StringVar(&chainFile, "chain", "markov.chain", "Markov Chain File")
	flag.StringVar(&mcflyFile, "mcfly", "mcfly.chain", "McFly Chain File")
	flag.Parse()

	marvin.Run(configDir, chainFile, mcflyFile)
}
